package accel

// #include <Accelerate/Accelerate.h>
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
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

// Converts an array of signed 8-bit integers to single-precision floating-point values.
func Vflt8(src []int8, dest []float32, srcStride, destStride int) {
	C.vDSP_vflt8((*C.char)(&src[0]), C.vDSP_Stride(srcStride), (*C.float)(&dest[0]), C.vDSP_Stride(destStride), C.vDSP_Length(len(src)/srcStride))
}

// Converts an array of signed 8-bit integers to single-precision floating-point values.
func Vflt8_byte(src []byte, dest []float32, srcStride, destStride int) {
	C.vDSP_vflt8((*C.char)(unsafe.Pointer(&src[0])), C.vDSP_Stride(srcStride), (*C.float)(&dest[0]), C.vDSP_Stride(destStride), C.vDSP_Length(len(src)/srcStride))
}

// Converts an array of unsigned 8-bit integers to single-precision floating-point values.
func Vfltu8(src []byte, dest []float32, srcStride, destStride int) {
	C.vDSP_vfltu8((*C.uchar)(&src[0]), C.vDSP_Stride(srcStride), (*C.float)(&dest[0]), C.vDSP_Stride(destStride), C.vDSP_Length(len(src)/srcStride))
}

// Vector convert double-precision to single-precision.
func Vdpsp(src []float64, dest []float32, srcStride, destStride int) {
	C.vDSP_vdpsp((*C.double)(&src[0]), C.vDSP_Stride(srcStride), (*C.float)(&dest[0]), C.vDSP_Stride(destStride), C.vDSP_Length(len(src)/srcStride))
}

// Vector convert double-precision to single-precision. Operate on a byte buffer which contains float64
func Vdpsp_byte(src []byte, dest []float32, srcStride, destStride int) {
	n := len(src) / 8
	if n > 0 {
		C.vDSP_vdpsp((*C.double)(unsafe.Pointer(&src[0])), C.vDSP_Stride(srcStride), (*C.float)(&dest[0]), C.vDSP_Stride(destStride), C.vDSP_Length(n/srcStride))
	}
}

// Copies the contents of an interleaved complex vector C to a split complex vector Z; single precision.
func Ctoz(src []complex64, dest DSPSplitComplex, srcStride, destStride int) {
	var splitComplex C.DSPSplitComplex
	splitComplex.realp = (*C.float)(&dest.Real[0])
	splitComplex.imagp = (*C.float)(&dest.Imag[0])
	C.vDSP_ctoz((*C.DSPComplex)(unsafe.Pointer(&src[0])), C.vDSP_Stride(srcStride), &splitComplex, C.vDSP_Stride(destStride), C.vDSP_Length(2*len(src)/srcStride))
}

// Copies the contents of an interleaved complex vector C to a split complex vector Z; single precision. Operate on a byte buffer which contains complex64
func Ctoz_byte(src []byte, dest DSPSplitComplex, srcStride, destStride int) {
	n := len(src) / 8
	var splitComplex C.DSPSplitComplex
	splitComplex.realp = (*C.float)(&dest.Real[0])
	splitComplex.imagp = (*C.float)(&dest.Imag[0])
	C.vDSP_ctoz((*C.DSPComplex)(unsafe.Pointer(&src[0])), C.vDSP_Stride(srcStride), &splitComplex, C.vDSP_Stride(destStride), C.vDSP_Length(2*n/srcStride))
}

// Copies the contents of a split complex vector Z to an interleaved complex vector C; single precision.
func Ztoc(src DSPSplitComplex, dest []complex64, srcStride, destStride int) {
	var splitComplex C.DSPSplitComplex
	splitComplex.realp = (*C.float)(&src.Real[0])
	splitComplex.imagp = (*C.float)(&src.Imag[0])
	C.vDSP_ztoc(&splitComplex, C.vDSP_Stride(srcStride), (*C.DSPComplex)(unsafe.Pointer(&dest[0])), C.vDSP_Stride(destStride), C.vDSP_Length(2*len(dest)/destStride))
}

// Copies the contents of a split complex vector Z to an interleaved complex vector C; single precision.
func Ztoc_float(src DSPSplitComplex, dest []float32, srcStride, destStride int) {
	var splitComplex C.DSPSplitComplex
	splitComplex.realp = (*C.float)(&src.Real[0])
	splitComplex.imagp = (*C.float)(&src.Imag[0])
	C.vDSP_ztoc(&splitComplex, C.vDSP_Stride(srcStride), (*C.DSPComplex)(unsafe.Pointer(&dest[0])), C.vDSP_Stride(destStride), C.vDSP_Length(len(dest)/destStride))
}

// Copies the contents of a split complex vector Z to an interleaved complex vector C; single precision. Operate on a byte buffer which contains complex64
func Ztoc_byte(src DSPSplitComplex, dest []byte, srcStride, destStride int) {
	var splitComplex C.DSPSplitComplex
	splitComplex.realp = (*C.float)(&src.Real[0])
	splitComplex.imagp = (*C.float)(&src.Imag[0])
	C.vDSP_ztoc(&splitComplex, C.vDSP_Stride(srcStride), (*C.DSPComplex)(unsafe.Pointer(&dest[0])), C.vDSP_Stride(destStride), C.vDSP_Length(len(dest)/4/destStride))
}

// Vector clear; single precision.
func Vclr(vec []float32, stride int) {
	C.vDSP_vclr((*C.float)(&vec[0]), (C.vDSP_Stride)(stride), (C.vDSP_Length)(len(vec)/stride))
}

// Convolution with decimation; single precision.
func Desamp(input []float32, desamplingFactor int, coeff []float32, output []float32) {
	C.vDSP_desamp((*C.float)(&input[0]), C.vDSP_Stride(desamplingFactor), (*C.float)(&coeff[0]), (*C.float)(&output[0]), C.vDSP_Length(len(output)), C.vDSP_Length(len(coeff)))
}

// Complex-real downsample with anti-aliasing; single precision.
func Zrdesamp(src DSPSplitComplex, decimationFactor int, coefficients []float32, dest DSPSplitComplex) {
	var srcC C.DSPSplitComplex
	srcC.realp = (*C.float)(&src.Real[0])
	srcC.imagp = (*C.float)(&src.Imag[0])
	var dstC C.DSPSplitComplex
	dstC.realp = (*C.float)(&dest.Real[0])
	dstC.imagp = (*C.float)(&dest.Imag[0])
	C.vDSP_zrdesamp(&srcC, C.vDSP_Stride(decimationFactor), (*C.float)(&coefficients[0]), &dstC, C.vDSP_Length(len(dest.Real)), C.vDSP_Length(len(coefficients)))
}

// Complex vector phase; single precision.
func Zvphas(input DSPSplitComplex, inputStride int, output []float32, outputStride int) {
	var srcC C.DSPSplitComplex
	srcC.realp = (*C.float)(&input.Real[0])
	srcC.imagp = (*C.float)(&input.Imag[0])
	C.vDSP_zvphas(&srcC, C.vDSP_Stride(inputStride), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), C.vDSP_Length(len(output)))
}

// Calculates the conjugate dot product (or inner dot product) of complex vectors A and B and leave the result in complex vector C; single precision.
func Zidotpr(input1 DSPSplitComplex, stride1 int, input2 DSPSplitComplex, stride2 int, result DSPSplitComplex) {
	var in1 C.DSPSplitComplex
	in1.realp = (*C.float)(&input1.Real[0])
	in1.imagp = (*C.float)(&input1.Imag[0])
	var in2 C.DSPSplitComplex
	in2.realp = (*C.float)(&input2.Real[0])
	in2.imagp = (*C.float)(&input2.Imag[0])
	var res C.DSPSplitComplex
	res.realp = (*C.float)(&result.Real[0])
	res.imagp = (*C.float)(&result.Imag[0])
	C.vDSP_zidotpr(&in1, C.vDSP_Stride(stride1), &in2, C.vDSP_Stride(stride2), &res, C.vDSP_Length(len(result.Real)))
}

// Complex vector conjugate and multiply; single precision.
func Zvcmul(input1 DSPSplitComplex, stride1 int, input2 DSPSplitComplex, stride2 int, result DSPSplitComplex, resultStride int) {
	var in1 C.DSPSplitComplex
	in1.realp = (*C.float)(&input1.Real[0])
	in1.imagp = (*C.float)(&input1.Imag[0])
	var in2 C.DSPSplitComplex
	in2.realp = (*C.float)(&input2.Real[0])
	in2.imagp = (*C.float)(&input2.Imag[0])
	var res C.DSPSplitComplex
	res.realp = (*C.float)(&result.Real[0])
	res.imagp = (*C.float)(&result.Imag[0])
	C.vDSP_zvcmul(&in1, C.vDSP_Stride(stride1), &in2, C.vDSP_Stride(stride2), &res, C.vDSP_Stride(resultStride), C.vDSP_Length(len(result.Real)/resultStride))
}

// Converts an array of single-precision floating-point values to signed 16-bit integer values, rounding towards zero.
func Vfix16(input []float32, inputStride int, output []int16, outputStride int) {
	C.vDSP_vfix16((*C.float)(&input[0]), C.vDSP_Stride(inputStride), (*C.short)(&output[0]), C.vDSP_Stride(outputStride), C.vDSP_Length(len(output)/outputStride))
}

// Converts an array of single-precision floating-point values to signed 16-bit integer values, rounding towards zero. Output to a byte stream.
func Vfix16_byte(input []float32, inputStride int, output []byte, outputStride int) {
	C.vDSP_vfix16((*C.float)(&input[0]), C.vDSP_Stride(inputStride), (*C.short)(unsafe.Pointer(&output[0])), C.vDSP_Stride(outputStride), C.vDSP_Length(len(output)/2/outputStride))
}

// Vector scalar multiply and scalar add; single precision.
func Vsmsa(input []float32, inputStride int, mult, add float32, output []float32, outputStride int) {
	C.vDSP_vsmsa((*C.float)(&input[0]), C.vDSP_Stride(inputStride), (*C.float)(&mult), (*C.float)(&add), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), C.vDSP_Length(len(output)/outputStride))
}
