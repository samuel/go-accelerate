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
	// "unsafe"

	"github.com/samuel/go-accelerate"
)

var (
	flagLog2n        = flag.Int("log2n", 10, "log2n of number of samples for FFT (2^log2n samples)")
	flagSampleFormat = flag.String("sample.format", "8uc", "Sample format")
	flagSampleRate   = flag.Float64("sample.rate", 0.0, "Sample rate")
	flagScale        = flag.Float64("scale", 0.0, "Scale for the magnitude (default is 0.0 which means to use scaleRatio)")
	flagScaleRatio   = flag.Float64("scaleRatio", 0.5, "Ratio of max magnitude to use as scale (if scale is 0.0)")
	flagMaxHeight    = flag.Int("maxHeight", 480, "Max height of image.")
	flagHeight       = flag.Int("height", 0, "Height of output image (default is 0 meaning to make it up to maxHeight or out of samples)")
	flagWidth        = flag.Int("width", 640, "Width of output image")
)

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

// var colorIndex = (int)((_powerSpectrum[i] + _contrast * 50.0 / 25.0) * _gradientPixels.Length / byte.MaxValue);
// colorIndex = Math.Max(colorIndex, 0);
// colorIndex = Math.Min(colorIndex, _gradientPixels.Length - 1);

// *ptr++ = _gradientPixels[colorIndex];

// private void BuildGradientVector()
// {
// if (_gradientPixels == null || _gradientPixels.Length != ClientRectangle.Height - AxisMargin)
// {
//     _gradientPixels = new int[ClientRectangle.Height - AxisMargin - 1];
// }
// for (var i = 0; i < _gradientPixels.Length; i++)
// {
//     _gradientPixels[_gradientPixels.Length - i - 1] = _buffer.GetPixel(
// ClientRectangle.Width - AxisMargin / 2, i + AxisMargin / 2 + 1).ToArgb();
// }
// }

func usage() {
	fmt.Println("syntax: fft [options] <input file.samples> <output file.png>")
	flag.PrintDefaults()
	fmt.Printf("\nSample formats:\n")
	for name, maker := range sampleFormats {
		fmt.Printf("  %s: %s\n", name, maker().Description())
	}
	os.Exit(1)
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 2 {
		usage()
	}

	sampleMaker := sampleFormats[*flagSampleFormat]
	if sampleMaker == nil {
		println("ERROR: unknown sample format", *flagSampleFormat)
		println()
		usage()
	}
	sampler := sampleMaker()

	log2n := *flagLog2n
	nSamples := 1 << uint(log2n)
	radix := accel.FFTRadix2
	height := *flagHeight
	width := *flagWidth
	maxHeight := *flagMaxHeight
	scale := float32(*flagScale)
	inpath := flag.Arg(0)
	outpath := flag.Arg(1)

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

	for y := 0; y < height; y++ {
		if err := sampler.Read(file, data); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		fft.Zip(data, 1, log2n, accel.FFTDirectionForward)

		maxM := float32(0.0)
		for n := 0; n < nSamples; n++ {
			magnitude := float32(math.Sqrt(float64(data.Real[n]*data.Real[n] + data.Imag[n]*data.Imag[n])))
			if magnitude > maxM {
				maxM = magnitude
			}
			data.Real[n] = magnitude
		}
		if scale == 0.0 {
			scale = 1 / (maxM * float32(*flagScaleRatio))
		}

		dx := nSamples / width
		if dx == 0 {
			dx = width / nSamples
		} else {
			off := y * img.Stride
			for x := 0; x < width; x++ {
				sum := float32(0)
				n := 0
				for j := x * dx; j < x*dx+dx && j < nSamples; j++ {
					sum += data.Real[(j+nSamples/2)%nSamples]
					n++
				}
				if math.IsInf(float64(sum), 0) {
					sum = 1.0
				} else if math.IsNaN(float64(sum)) {
					sum = 0.0
				}
				sum /= float32(n)
				c := colorForValue(sum * scale)
				img.Pix[off] = c.R
				img.Pix[off+1] = c.G
				img.Pix[off+2] = c.B
				img.Pix[off+3] = c.A
				off += 4
			}
		}
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
