package accel

import (
	"math"
	"testing"
)

func TestVvlog10f(t *testing.T) {
	input := []float32{0.1, 0.8, 1.0, 2.5, 10.0}
	output := make([]float32, len(input))
	Vvlog10f(output, input)
	for i := 0; i < len(output); i++ {
		expected := float32(math.Log10(float64(input[i])))
		if output[i] != expected {
			t.Errorf("Expected log(%f) to return %f instead of %f", input[i], expected, output[i])
		}
	}
	// Test in-place
	Vvlog10f(input, input)
	for i := 0; i < len(input); i++ {
		if output[i] != input[i] {
			t.Errorf("In-place: expected %f instead of %f", output[i], input[i])
		}
	}
}
