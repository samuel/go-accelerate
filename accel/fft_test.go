package accel

import "testing"

func BenchmarkFFTZip10Radix2(b *testing.B) {
	fft, err := CreateFFTSetup(10, FFTRadix2)
	if err != nil {
		b.Fatal(err)
	}
	samples := DSPSplitComplex{
		Real: make([]float32, 1024),
		Imag: make([]float32, 1024),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fft.Zip(samples, 1, 10, FFTDirectionForward)
	}
}
