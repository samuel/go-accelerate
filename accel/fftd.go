package accel

// #include <Accelerate/Accelerate.h>
import "C"

import "runtime"

type FFTSetupD struct {
	cFFTSetupD C.FFTSetupD
}

func CreateFFTSetupD(log2n int, radix FFTRadix) (*FFTSetupD, error) {
	fftSetup := C.vDSP_create_fftsetupD(C.vDSP_Length(log2n), C.FFTRadix(radix))
	if fftSetup == nil {
		return nil, ErrFailedToCreateFFTSetup
	}
	setup := &FFTSetupD{fftSetup}
	runtime.SetFinalizer(setup, destroyFFTSetupD)
	return setup, nil
}

func (fs *FFTSetupD) Destroy() {
	destroyFFTSetupD(fs)
}

func destroyFFTSetupD(fftSetup *FFTSetupD) {
	if fftSetup != nil && fftSetup.cFFTSetupD != nil {
		C.vDSP_destroy_fftsetupD(fftSetup.cFFTSetupD)
		fftSetup.cFFTSetupD = nil
	}
}

// Zrip computess an in-place single-precision real discrete Fourier transform of the
// input/output vector signal, either from the time domain to the frequency domain
// (forward) or from the frequency domain to the time domain (inverse).
func (fs *FFTSetupD) Zrip(ioData DSPDoubleSplitComplex, stride, log2n int, direction FFTDirection) {
	var splitComplex C.DSPDoubleSplitComplex
	splitComplex.realp = (*C.double)(&ioData.Real[0])
	splitComplex.imagp = (*C.double)(&ioData.Imag[0])
	C.vDSP_fft_zripD(fs.cFFTSetupD, &splitComplex, C.vDSP_Stride(stride), C.vDSP_Length(log2n), C.FFTDirection(direction))
}

// Zop computes an out-of-place single-precision real discrete Fourier transform of the
// input vector, either from the time domain to the frequency domain (forward) or from the
// frequency domain to the time domain (inverse).
func (fs *FFTSetupD) Zrop(input DSPDoubleSplitComplex, inputStride int, output DSPDoubleSplitComplex, outputStride int, log2n int, direction FFTDirection) {
	var inC C.DSPDoubleSplitComplex
	inC.realp = (*C.double)(&input.Real[0])
	inC.imagp = (*C.double)(&input.Imag[0])
	var outC C.DSPDoubleSplitComplex
	outC.realp = (*C.double)(&output.Real[0])
	outC.imagp = (*C.double)(&output.Imag[0])
	C.vDSP_fft_zopD(fs.cFFTSetupD, &inC, C.vDSP_Stride(inputStride), &outC, C.vDSP_Stride(outputStride), C.vDSP_Length(log2n), C.FFTDirection(direction))
}

// Zip computess an in-place single-precision complex discrete Fourier transform of the
// input/output vector signal, either from the time domain to the frequency domain
// (forward) or from the frequency domain to the time domain (inverse).
func (fs *FFTSetupD) Zip(ioData DSPDoubleSplitComplex, stride, log2n int, direction FFTDirection) {
	var splitComplex C.DSPDoubleSplitComplex
	splitComplex.realp = (*C.double)(&ioData.Real[0])
	splitComplex.imagp = (*C.double)(&ioData.Imag[0])
	C.vDSP_fft_zipD(fs.cFFTSetupD, &splitComplex, C.vDSP_Stride(stride), C.vDSP_Length(log2n), C.FFTDirection(direction))
}

// Zop computes an out-of-place single-precision complex discrete Fourier transform of the
// input vector, either from the time domain to the frequency domain (forward) or from the
// frequency domain to the time domain (inverse).
func (fs *FFTSetupD) Zop(input DSPDoubleSplitComplex, inputStride int, output DSPDoubleSplitComplex, outputStride int, log2n int, direction FFTDirection) {
	var inC C.DSPDoubleSplitComplex
	inC.realp = (*C.double)(&input.Real[0])
	inC.imagp = (*C.double)(&input.Imag[0])
	var outC C.DSPDoubleSplitComplex
	outC.realp = (*C.double)(&output.Real[0])
	outC.imagp = (*C.double)(&output.Imag[0])
	C.vDSP_fft_zopD(fs.cFFTSetupD, &inC, C.vDSP_Stride(inputStride), &outC, C.vDSP_Stride(outputStride), C.vDSP_Length(log2n), C.FFTDirection(direction))
}
