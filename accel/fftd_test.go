package accel

import "testing"

func BenchmarkFFTDoubleZip10Radix2(b *testing.B) {
	fft, err := CreateFFTSetupD(10, FFTRadix2)
	if err != nil {
		b.Fatal(err)
	}
	samples := DSPDoubleSplitComplex{
		Real: make([]float64, 1024),
		Imag: make([]float64, 1024),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fft.Zip(samples, 1, 10, FFTDirectionForward)
	}
}
