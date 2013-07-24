package accel

import (
	"math"
	"testing"
)

const maxFloatDiffErr = 1e-4

func almostEqual32(a, b, err float32) bool {
	return almostEqual64(float64(a), float64(b), float64(err))
}

func almostEqual64(a, b, err float64) bool {
	return math.Abs(a-b) < err
}

func TestVfltu8(t *testing.T) {
	input := []byte{4, 127, 250, 190}
	output := make([]float32, 8)

	for i := 0; i < len(output); i++ {
		output[i] = float32(math.NaN())
	}

	Vfltu8(input, 1, output, 1)
	for i, x := range input {
		expected := float32(x)
		if expected != output[i] {
			t.Errorf("Vfltu8 strides == 1 : output %f != expected %f for index %d", output[i], expected, i)
		}
	}

	for i := 0; i < len(output); i++ {
		output[i] = float32(math.NaN())
	}

	Vfltu8(input, 2, output, 1)
	for i := 0; i < len(input); i += 2 {
		expected := float32(input[i])
		if expected != output[i/2] {
			t.Errorf("Vfltu8 in stride = 2 out stride = 1 : output %f != expected %f for index %d", output[i/2], expected, i)
		}
	}
	for i := len(input) / 2; i < len(output); i++ {
		if !math.IsNaN(float64(output[i])) {
			t.Errorf("Vfltu8 wrote too far for input stride 2 (%d=%f)", i, output[i])
		}
	}
}

func TestZtoc(t *testing.T) {
	split := DSPSplitComplex{
		Real: []float32{2.0, 4.0, 8.0, 16.0, 32.0, 64.0, 128.0, 256.0},
		Imag: []float32{1.0, 3.0, 7.0, 15.0, 31.0, 63.0, 127.0, 255.0},
	}
	n := len(split.Real)
	output := make([]complex64, n*2) // extra padding to check for overflow
	for i := 0; i < len(output); i++ {
		output[i] = complex(float32(math.NaN()), float32(math.NaN()))
	}
	Ztoc(split, 1, output[:n], 2)
	for i := 0; i < n; i++ {
		expected := complex(split.Real[i], split.Imag[i])
		if output[i] != expected {
			t.Errorf("failed for strides of 1 index %d : output %+v != expected %+v", i, output[i*2], expected)
		}
	}
}

func TestZtoc_float(t *testing.T) {
	split := DSPSplitComplex{
		Real: []float32{2.0, 4.0, 8.0, 16.0, 32.0, 64.0, 128.0, 256.0},
		Imag: []float32{1.0, 3.0, 7.0, 15.0, 31.0, 63.0, 127.0, 255.0},
	}
	n := len(split.Real)
	output := make([]float32, n*4) // extra padding to check for overflow
	for i := 0; i < len(output); i++ {
		output[i] = float32(math.NaN())
	}
	Ztoc_float(split, 1, output[:n*2], 2)
	for i := 0; i < n; i++ {
		out := complex(output[i*2], output[i*2+1])
		expected := complex(split.Real[i], split.Imag[i])
		if out != expected {
			t.Errorf("failed for strides of 1 index %d : output %+v != expected %+v", i, out, expected)
		}
	}
}

// func lowPassReal(samples DSPSplitComplex, fast, slow int) DSPSplitComplex {
// 	i2 := 0
// 	var nowLPR complex64 = complex(0, 0)
// 	prevLPRIndex := 0
// 	fastSlowRatio := float32(fast) / float32(slow)
// 	for i := 0; i < len(samples.Real); i++ {
// 		nowLPR += complex(samples.Real[i], samples.Imag[i])
// 		prevLPRIndex += slow
// 		if prevLPRIndex < fast {
// 			continue
// 		}
// 		samples.Real[i2] = real(nowLPR) / fastSlowRatio
// 		samples.Imag[i2] = imag(nowLPR) / fastSlowRatio
// 		prevLPRIndex -= fast
// 		nowLPR = 0
// 		i2++
// 	}
// 	return samples
// }

// func TestZrdesamp(t *testing.T) {
// 	src := DSPSplitComplex{
// 		Real: []float32{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096},
// 		Imag: []float32{1, 3, 6, 12, 24, 48, 96, 192, 384, 768, 1536, 3072, 6144},
// 	}
// 	dst := DSPSplitComplex{make([]float32, len(src.Real)/2+1), make([]float32, len(src.Imag)/2+1)}
// 	scale := float32(1.0 / 2.0)
// 	Zrdesamp(src, 1, []float32{scale, scale, scale, scale, scale, scale, scale}, dst, 2)
// 	t.Logf("%+v", dst)
// 	src = lowPassReal(src, 3, 2)
// 	t.Logf("%+v", src)
// }

// Vector scalar divide; single precision.
func TestVsdiv(t *testing.T) {
	input := []float32{1.0, 2.0, 3.0, 4.0, 5.0}
	output := make([]float32, len(input))
	Vsdiv(input, 1, 3.0, output, 1)
	for i := 0; i < len(output); i++ {
		expected := input[i] / 3.0
		if !almostEqual32(output[i], expected, maxFloatDiffErr) {
			t.Errorf("Expected %f/3.0 to return %f instead of %f", input[i], expected, output[i])
		}
	}
}
