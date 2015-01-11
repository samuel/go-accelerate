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

const (
	FFTRadix2 FFTRadix = C.kFFTRadix2
	FFTRadix3 FFTRadix = C.kFFTRadix3
	FFTRadix5 FFTRadix = C.kFFTRadix5

	FFTDirectionForward FFTDirection = C.kFFTDirection_Forward
	FFTDirectionInverse FFTDirection = C.kFFTDirection_Inverse
)

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

func destroyFFTSetup(fftSetup *FFTSetup) {
	if fftSetup != nil && fftSetup.cFFTSetup != nil {
		C.vDSP_destroy_fftsetup(fftSetup.cFFTSetup)
		fftSetup.cFFTSetup = nil
	}
}

// Zrip computess an in-place single-precision real discrete Fourier transform of the
// input/output vector signal, either from the time domain to the frequency domain
// (forward) or from the frequency domain to the time domain (inverse).
func (fs *FFTSetup) Zrip(ioData DSPSplitComplex, stride, log2n int, direction FFTDirection) {
	var splitComplex C.DSPSplitComplex
	splitComplex.realp = (*C.float)(&ioData.Real[0])
	splitComplex.imagp = (*C.float)(&ioData.Imag[0])
	C.vDSP_fft_zrip(fs.cFFTSetup, &splitComplex, C.vDSP_Stride(stride), C.vDSP_Length(log2n), C.FFTDirection(direction))
}

// Zop computes an out-of-place single-precision real discrete Fourier transform of the
// input vector, either from the time domain to the frequency domain (forward) or from the
// frequency domain to the time domain (inverse).
func (fs *FFTSetup) Zrop(input DSPSplitComplex, inputStride int, output DSPSplitComplex, outputStride int, log2n int, direction FFTDirection) {
	var inC C.DSPSplitComplex
	inC.realp = (*C.float)(&input.Real[0])
	inC.imagp = (*C.float)(&input.Imag[0])
	var outC C.DSPSplitComplex
	outC.realp = (*C.float)(&output.Real[0])
	outC.imagp = (*C.float)(&output.Imag[0])
	C.vDSP_fft_zop(fs.cFFTSetup, &inC, C.vDSP_Stride(inputStride), &outC, C.vDSP_Stride(outputStride), C.vDSP_Length(log2n), C.FFTDirection(direction))
}

// Zip computess an in-place single-precision complex discrete Fourier transform of the
// input/output vector signal, either from the time domain to the frequency domain
// (forward) or from the frequency domain to the time domain (inverse).
func (fs *FFTSetup) Zip(ioData DSPSplitComplex, stride, log2n int, direction FFTDirection) {
	var splitComplex C.DSPSplitComplex
	splitComplex.realp = (*C.float)(&ioData.Real[0])
	splitComplex.imagp = (*C.float)(&ioData.Imag[0])
	C.vDSP_fft_zip(fs.cFFTSetup, &splitComplex, C.vDSP_Stride(stride), C.vDSP_Length(log2n), C.FFTDirection(direction))
}

// Zop computes an out-of-place single-precision complex discrete Fourier transform of the
// input vector, either from the time domain to the frequency domain (forward) or from the
// frequency domain to the time domain (inverse).
func (fs *FFTSetup) Zop(input DSPSplitComplex, inputStride int, output DSPSplitComplex, outputStride int, log2n int, direction FFTDirection) {
	var inC C.DSPSplitComplex
	inC.realp = (*C.float)(&input.Real[0])
	inC.imagp = (*C.float)(&input.Imag[0])
	var outC C.DSPSplitComplex
	outC.realp = (*C.float)(&output.Real[0])
	outC.imagp = (*C.float)(&output.Imag[0])
	C.vDSP_fft_zop(fs.cFFTSetup, &inC, C.vDSP_Stride(inputStride), &outC, C.vDSP_Stride(outputStride), C.vDSP_Length(log2n), C.FFTDirection(direction))
}
