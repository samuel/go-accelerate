package accel

// #include <Accelerate/Accelerate.h>
import "C"

import (
	"errors"
	"runtime"
)

var ErrFailedToCreateFFTSetup = errors.New("accel: failed to create FFT setup")

type FFTSetup struct {
	cFFTSetup C.FFTSetup
}
type FFTRadix C.FFTRadix
type FFTDirection C.FFTDirection

type DSPSplitComplex struct {
	Real []float32
	Imag []float32
}

const (
	FFTRadix2 FFTRadix = C.kFFTRadix2
	FFTRadix3 FFTRadix = C.kFFTRadix3
	FFTRadix5 FFTRadix = C.kFFTRadix5

	FFTDirectionForward FFTDirection = C.kFFTDirection_Forward
	FFTDirectionInverse FFTDirection = C.kFFTDirection_Inverse
)

func destroyFFTSetup(fftSetup *FFTSetup) {
	if fftSetup != nil && fftSetup.cFFTSetup != nil {
		C.vDSP_destroy_fftsetup(fftSetup.cFFTSetup)
		fftSetup.cFFTSetup = nil
	}
}

func CreateFFTSetup(log2n int, radix FFTRadix) (*FFTSetup, error) {
	fftSetup := C.vDSP_create_fftsetup(C.vDSP_Length(log2n), C.FFTRadix(radix))
	if fftSetup == nil {
		return nil, ErrFailedToCreateFFTSetup
	}
	setup := &FFTSetup{fftSetup}
	runtime.SetFinalizer(setup, destroyFFTSetup)
	return setup, nil
}

func (fs *FFTSetup) Destroy() {
	destroyFFTSetup(fs)
}

// Computes an in-place single-precision complex discrete Fourier transform of the
// input/output vector signal, either from the time domain to the frequency domain
// (forward) or from the frequency domain to the time domain (inverse).
func (fs *FFTSetup) Zip(ioData DSPSplitComplex, stride, log2n int, direction FFTDirection) {
	var splitComplex C.DSPSplitComplex
	splitComplex.realp = (*C.float)(&ioData.Real[0])
	splitComplex.imagp = (*C.float)(&ioData.Imag[0])
	C.vDSP_fft_zip(fs.cFFTSetup, &splitComplex, C.vDSP_Stride(stride), C.vDSP_Length(log2n), C.FFTDirection(direction))
}
