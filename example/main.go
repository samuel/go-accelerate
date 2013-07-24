package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"math"
	"os"
	"time"

	"github.com/samuel/go-accelerate/accel"
)

func writeImage(img image.Image, name string) error {
	w, err := os.Create(name)
	if err != nil {
		return err
	}
	if err := jpeg.Encode(w, img, nil); err != nil {
		return err
	}
	w.Close()
	return nil
}

func CalcGaussian1D(stddev float64, radius int) []float64 {
	scale := 1.0 / (math.Sqrt(2*math.Pi) * stddev)
	d2 := 2.0 * stddev * stddev
	size := 2*radius + 1
	out := make([]float64, size)
	sum := float64(0)
	for i := 0; i < size; i++ {
		x := i - radius
		v := scale * math.Exp(-float64(x*x)/d2)
		out[i] = v
		sum += v
	}
	for i := 0; i < size; i++ {
		out[i] = out[i] / sum
	}
	return out
}

func CalcGaussian1Di16(stddev float64, radius, scale int) []int16 {
	k := CalcGaussian1D(stddev, radius)
	k16 := make([]int16, len(k))
	for i, v := range k {
		k16[i] = int16(v * float64(scale))
	}
	return k16
}

func main() {
	rd, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(rd)
	if err != nil {
		log.Fatal(err)
	}
	rd.Close()

	src, format := accel.VImageBufferFromImage(img)
	switch format {
	case "argb8888":
		// The format we expect
	case "rgba8888":
		if err := accel.VImagePermuteChannels_ARGB8888(src, src, [4]uint8{3, 0, 1, 2}, accel.VImageFlagNoFlags); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Unsupported format %s", format)
	}

	scale := 1024
	kernel := CalcGaussian1Di16(4.5, 31, scale)

	dst := accel.CreateVImageBuffer(img.Bounds().Dx(), img.Bounds().Dy(), 4, 0)
	t := time.Now()
	if err := accel.VImageConvolve_ARGB8888(src, dst, nil, 0, 0, kernel, 1, len(kernel), scale, [4]uint8{}, accel.VImageFlagEdgeExtend); err != nil {
		log.Fatal(err)
	}
	if err := accel.VImageConvolve_ARGB8888(dst, src, nil, 0, 0, kernel, len(kernel), 1, scale, [4]uint8{}, accel.VImageFlagEdgeExtend); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Convolution: %d ms\n", time.Since(t).Nanoseconds()/1e6)
	t = time.Now()
	if err := accel.VImagePermuteChannels_ARGB8888(src, dst, [4]uint8{1, 2, 3, 0}, accel.VImageFlagNoFlags); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Permutation: %d ms\n", time.Since(t).Nanoseconds()/1e6)
	if err := writeImage(dst.ToRGBA(), "out.jpg"); err != nil {
		log.Fatal(err)
	}
}
