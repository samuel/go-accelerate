package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math"
	"os"

	"github.com/samuel/go-accelerate/accel"
)

var (
	flagLog2n        = flag.Int("log2n", 10, "log2n of number of samples for FFT (2^log2n samples)")
	flagSampleFormat = flag.String("sample.format", "8uc", "Sample format")
	flagSampleRate   = flag.Float64("sample.rate", 0.0, "Sample rate")
	flagScale        = flag.Float64("scale", 0.0, "Scale for the magnitude (default is 0.0 which means to use scaleRatio)")
	flagScaleLinear  = flag.Bool("scale.linear", false, "use a linear scale (default is log)")
	flagScaleRatio   = flag.Float64("scale.ratio", 0.5, "Ratio of max magnitude to use as scale (if scale is 0.0)")
	flagMaxHeight    = flag.Int("maxHeight", 480, "Max height of image.")
	flagHeight       = flag.Int("height", 0, "Height of output image (default is 0 meaning to make it up to maxHeight or out of samples)")
	flagWidth        = flag.Int("width", 640, "Width of output image")
	flagWindow       = flag.String("window", "hanning", "Window function (hanning, hamming, blackman, triangle)")
)

var windowFuncs = map[string]func([]float32){
	"hanning": func(output []float32) {
		accel.HannWindow(output, accel.WindowFlagHannNorm)
	},
	"hamming": func(output []float32) {
		accel.HammWindow(output, 0)
	},
	"blackman": func(output []float32) {
		accel.BlkmanWindow(output, 0)
	},
	"triangle": func(output []float32) {
		for n := 0; n < len(output); n++ {
			output[n] = float32(1 - math.Abs((float64(n)-float64(len(output)-1)/2.0)/(float64(len(output)+1)/2.0)))
		}
	},
}

var gradient = [13]color.RGBA{
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x20, 0xff},
	color.RGBA{0x00, 0x00, 0x30, 0xff},
	color.RGBA{0x00, 0x00, 0x50, 0xff},
	color.RGBA{0x00, 0x00, 0x91, 0xff},
	color.RGBA{0x1e, 0x90, 0xff, 0xff},
	color.RGBA{0xff, 0xff, 0x00, 0xff},
	color.RGBA{0xfe, 0x6d, 0x16, 0xff},
	color.RGBA{0xff, 0x00, 0x00, 0xff},
	color.RGBA{0xc6, 0x00, 0x00, 0xff},
	color.RGBA{0x9f, 0x00, 0x00, 0xff},
	color.RGBA{0x75, 0x00, 0x00, 0xff},
	color.RGBA{0x4a, 0x00, 0x00, 0xff},
}

func colorForValue(value float32) color.RGBA {
	if value < 0.0 {
		value = 0.0
	} else if value >= 1.0 {
		value = 1.0
	}

	colorF := value * float32(len(gradient)+1)
	alpha := colorF - float32(math.Floor(float64(colorF)))

	color1I := int(colorF)
	if color1I >= len(gradient) {
		return gradient[len(gradient)-1]
	}
	color2I := color1I + 1
	if color2I >= len(gradient) {
		color2I = len(gradient) - 1
	}
	color1 := gradient[color1I]
	color2 := gradient[color2I]
	return color.RGBA{
		uint8(int(color1.R) + int(float32(int(color2.R)-int(color1.R))*alpha)),
		uint8(int(color1.G) + int(float32(int(color2.G)-int(color1.G))*alpha)),
		uint8(int(color1.B) + int(float32(int(color2.B)-int(color1.B))*alpha)),
		255,
	}
}

func usage() {
	fmt.Println("syntax: fft [options] <input file.samples> <output file.png>")
	flag.PrintDefaults()
	fmt.Printf("\nSample formats:\n")
	for name, maker := range sampleFormats {
		fmt.Printf("  %s: %s\n", name, maker.Description())
	}
	os.Exit(1)
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 2 {
		usage()
	}

	sampler := sampleFormats[*flagSampleFormat]
	if sampler == nil {
		println("ERROR: unknown sample format", *flagSampleFormat)
		println()
		usage()
	}

	log2n := *flagLog2n
	nSamples := 1 << uint(log2n)
	radix := accel.FFTRadix2
	height := *flagHeight
	width := *flagWidth
	maxHeight := *flagMaxHeight
	scale := float32(*flagScale)
	inpath := flag.Arg(0)
	outpath := flag.Arg(1)
	windowFunc := windowFuncs[*flagWindow]
	if windowFunc == nil && *flagWindow != "" && *flagWindow != "none" {
		log.Fatalf("Unknown window function %s", *flagWindow)
	}
	// pre-calculate window
	var window []float32
	if windowFunc != nil {
		window = make([]float32, nSamples)
		windowFunc(window)
	}

	file, err := os.Open(inpath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	if height == 0 {
		size, err := file.Seek(0, 2)
		if err == nil {
			_, err = file.Seek(0, 0)
		}
		if err != nil {
			log.Fatal(err)
		}
		height = int(size / int64(nSamples*sampler.SampleSize()))
		if height == 0 {
			height = 1
		}
		if height > maxHeight {
			height = maxHeight
		}
	}

	fft, err := accel.CreateFFTSetup(log2n, radix)
	if err != nil {
		log.Fatal(err)
	}
	defer fft.Destroy()

	data := accel.DSPSplitComplex{
		Real: make([]float32, nSamples),
		Imag: make([]float32, nSamples),
	}
	// Align the buffers to 16-byte boundaries to allow SIMD to do its thing
	// var data accel.DSPSplitComplex
	// sampleBuf := make([]float32, (nSamples+4)*2)
	// data.Real = sampleBuf[:nSamples]
	// data.Imag = sampleBuf[nSamples : nSamples*2]
	// println(uintptr(unsafe.Pointer(&data.Imag[0])) & 0xf)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	buf := make([]byte, sampler.SampleSize()*nSamples)
	for y := 0; y < height; y++ {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		n /= sampler.SampleSize()
		sampler.Transform(buf[:n*sampler.SampleSize()], data)
		if n != nSamples {
			// Zero out any samples that weren't filled
			accel.Vclr(data.Real[n:], 1)
			accel.Vclr(data.Imag[n:], 1)
		}

		if windowFunc != nil {
			accel.Vmul(data.Real, 1, window, 1, data.Real, 1)
			accel.Vmul(data.Imag, 1, window, 1, data.Imag, 1)
		}

		fft.Zip(data, 1, log2n, accel.FFTDirectionForward)
		data.Real[0] = 0
		data.Imag[0] = 0
		accel.Zvabs(data, 1, data.Real, 1)
		maxM := accel.Maxv(data.Real, 1)
		if !*flagScaleLinear {
			accel.Vsdiv(data.Real, 1, maxM, data.Real, 1)
			accel.Vvlog10f(data.Real, data.Real)
			if scale == 0.0 {
				// mean := accel.Meanv(data.Real, 1)
				// scale = 1 / (mean * float32(*flagScaleRatio))
				scale = float32(*flagScaleRatio)
			}
			accel.Vsmsa(data.Real, 1, scale, 1.0, data.Real, 1)
		} else {
			if scale == 0.0 {
				scale = 1 / (maxM * float32(*flagScaleRatio))
			}
			accel.Vsmsa(data.Real, 1, scale, 0.0, data.Real, 1)
		}

		dx := nSamples / width
		if dx == 0 {
			dx = width / nSamples
			// TODO
		} else {
			x2 := width/2 - 1
			yoff := y * img.Stride
			for x := 0; x < width; x++ {
				sum := float32(0)
				n := 0
				for j := x * dx; j < x*dx+dx && j < nSamples; j++ {
					sum += data.Real[j]
					n++
				}
				if math.IsInf(float64(sum), 0) {
					sum = 1.0
				} else if math.IsNaN(float64(sum)) {
					sum = 0.0
				}
				sum /= float32(n)
				c := colorForValue(sum)
				off := yoff + x2*4
				img.Pix[off] = c.R
				img.Pix[off+1] = c.G
				img.Pix[off+2] = c.B
				img.Pix[off+3] = c.A
				x2 = (x2 + 1) % width
			}
		}
	}

	// Draw X scale marks
	for y := height - 8; y < height; y++ {
		off := y * img.Stride
		for x := 0; x < width; x++ {
			img.Pix[off+x*4] = 0
			img.Pix[off+x*4+1] = 0
			img.Pix[off+x*4+2] = 0
		}
		xoff := (width / 2) * 4
		img.Pix[off+xoff] = 0
		img.Pix[off+xoff+1] = 255
		img.Pix[off+xoff+2] = 0
		img.Pix[off+xoff+3] = 255
		for i := 4; i < 32; i = i * 2 {
			doff := 4 * width / i
			xoff = 4*width/2 - doff
			img.Pix[off+xoff] = 255
			img.Pix[off+xoff+1] = 255
			img.Pix[off+xoff+2] = 255
			img.Pix[off+xoff+3] = 255
			xoff += doff * 2
			img.Pix[off+xoff] = 255
			img.Pix[off+xoff+1] = 255
			img.Pix[off+xoff+2] = 255
			img.Pix[off+xoff+3] = 255
		}
	}
	sampleRate := *flagSampleRate
	if sampleRate != 0.0 {
		fmt.Printf("Sampler rate: %f Hz\n", sampleRate)
		fmt.Printf("Sampler rate f/2: %f Hz\n", sampleRate/2.0)
		fmt.Printf("Sampler rate f/4: %f Hz\n", sampleRate/4.0)
		fmt.Printf("Sampler rate f/8: %f Hz\n", sampleRate/8.0)
	}

	outFile, err := os.Create(outpath)
	if err != nil {
		log.Fatal(err)
	}
	if err := png.Encode(outFile, img); err != nil {
		log.Fatal(err)
	}
	outFile.Close()
}
