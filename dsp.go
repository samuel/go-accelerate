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

type WindowFlag int

const (
	WindowFlagHannDenorm WindowFlag = C.vDSP_HANN_DENORM // creates a denormalized window.
	WindowFlagHannNorm   WindowFlag = C.vDSP_HANN_NORM   // creates a normalized window.
	WindowFlagHalfWindow WindowFlag = C.vDSP_HALF_WINDOW // creates only the first (N+1)/2 points.
)

type DBFlag int

const (
	DBFlagPower     DBFlag = 0
	DBFlagAmplitude DBFlag = 1
)

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

// Computes an out-of-place single-precision complex discrete Fourier transform of the input vector, either from the time domain to the frequency domain (forward) or from the frequency domain to the time domain (inverse).
func (fs *FFTSetup) Zop(input DSPSplitComplex, inputStride int, output DSPSplitComplex, outputStride int, log2n int, direction FFTDirection) {
	var inC C.DSPSplitComplex
	inC.realp = (*C.float)(&input.Real[0])
	inC.imagp = (*C.float)(&input.Imag[0])
	var outC C.DSPSplitComplex
	outC.realp = (*C.float)(&output.Real[0])
	outC.imagp = (*C.float)(&output.Imag[0])
	C.vDSP_fft_zop(fs.cFFTSetup, &inC, C.vDSP_Stride(inputStride), &outC, C.vDSP_Stride(outputStride), C.vDSP_Length(log2n), C.FFTDirection(direction))
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
func Ztoc(src DSPSplitComplex, srcStride int, dest []complex64, destStride int) {
	var splitComplex C.DSPSplitComplex
	splitComplex.realp = (*C.float)(&src.Real[0])
	splitComplex.imagp = (*C.float)(&src.Imag[0])
	C.vDSP_ztoc(&splitComplex, C.vDSP_Stride(srcStride), (*C.DSPComplex)(unsafe.Pointer(&dest[0])), C.vDSP_Stride(destStride), C.vDSP_Length(2*len(dest)/destStride))
}

// Copies the contents of a split complex vector Z to an interleaved complex vector C; single precision.
func Ztoc_float(src DSPSplitComplex, srcStride int, dest []float32, destStride int) {
	var splitComplex C.DSPSplitComplex
	splitComplex.realp = (*C.float)(&src.Real[0])
	splitComplex.imagp = (*C.float)(&src.Imag[0])
	C.vDSP_ztoc(&splitComplex, C.vDSP_Stride(srcStride), (*C.DSPComplex)(unsafe.Pointer(&dest[0])), C.vDSP_Stride(destStride), C.vDSP_Length(len(dest)/destStride))
}

// Copies the contents of a split complex vector Z to an interleaved complex vector C; single precision. Operate on a byte buffer which contains complex64
func Ztoc_byte(src DSPSplitComplex, srcStride int, dest []byte, destStride int) {
	var splitComplex C.DSPSplitComplex
	splitComplex.realp = (*C.float)(&src.Real[0])
	splitComplex.imagp = (*C.float)(&src.Imag[0])
	C.vDSP_ztoc(&splitComplex, C.vDSP_Stride(srcStride), (*C.DSPComplex)(unsafe.Pointer(&dest[0])), C.vDSP_Stride(destStride), C.vDSP_Length(len(dest)/4/destStride))
}

// Vector clear; single precision.
func Vclr(vec []float32, stride int) {
	C.vDSP_vclr((*C.float)(&vec[0]), (C.vDSP_Stride)(stride), (C.vDSP_Length)(len(vec)/stride))
}

// Vector fill; single precision.
func Vfill(value float32, output []float32, stride int) {
	C.vDSP_vfill((*C.float)(&value), (*C.float)(&output[0]), (C.vDSP_Stride)(stride), (C.vDSP_Length)(len(output)/stride))
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

// Complex vector absolute values; single precision.
func Zvabs(input DSPSplitComplex, inputStride int, output []float32, outputStride int) {
	var in C.DSPSplitComplex
	in.realp = (*C.float)(&input.Real[0])
	in.imagp = (*C.float)(&input.Imag[0])
	C.vDSP_zvabs(&in, C.vDSP_Stride(inputStride), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), C.vDSP_Length(len(output)/outputStride))
}

// Vector maximum value; single precision.
func Maxv(input []float32, stride int) float32 {
	var out C.float
	C.vDSP_maxv((*C.float)(&input[0]), C.vDSP_Stride(stride), &out, C.vDSP_Length(len(input)/stride))
	return float32(out)
}

// Multiplies vector A by vector B and leaves the result in vector C; single precision.
func Vmul(input1 []float32, stride1 int, input2 []float32, stride2 int, output []float32, outputStride int) {
	C.vDSP_vmul((*C.float)(&input1[0]), C.vDSP_Stride(stride1), (*C.float)(&input2[0]), C.vDSP_Stride(stride2), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), C.vDSP_Length(len(output)/outputStride))
}

// Vector scalar add; single precision.
func Vsadd(input []float32, inputStride int, add float32, output []float32, outputStride int) {
	C.vDSP_vsadd((*C.float)(&input[0]), C.vDSP_Stride(inputStride), (*C.float)(&add), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), C.vDSP_Length(len(output)/outputStride))
}

// Vector scalar divide; single precision.
func Vsdiv(input []float32, inputStride int, divisor float32, output []float32, outputStride int) {
	C.vDSP_vsdiv((*C.float)(&input[0]), C.vDSP_Stride(inputStride), (*C.float)(&divisor), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), C.vDSP_Length(len(output)/outputStride))
}

// Vector negative values; single precision.
func Vneg(input []float32, inputStride int, output []float32, outputStride int) {
	C.vDSP_vneg((*C.float)(&input[0]), C.vDSP_Stride(inputStride), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), C.vDSP_Length(len(output)/outputStride))
}

// Vector convert power or amplitude to decibels; single precision.
// α * log10(input(n)/zeroReference) [α is 20 if Amplitude, or 10 if F is Power]
func Vdbcon(input []float32, inputStride int, zeroReference float32, output []float32, outputStride int, flag DBFlag) {
	C.vDSP_vdbcon((*C.float)(&input[0]), C.vDSP_Stride(inputStride), (*C.float)(&zeroReference), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), C.vDSP_Length(len(output)/outputStride), C.uint(flag))
}

// Vector mean value; single precision.
func Meanv(input []float32, stride int) float32 {
	var mean C.float
	C.vDSP_meanv((*C.float)(&input[0]), C.vDSP_Stride(stride), &mean, C.vDSP_Length(len(input)/stride))
	return float32(mean)
}

// Creates a single-precision Hanning window.
func HannWindow(output []float32, flag WindowFlag) {
	C.vDSP_hann_window((*C.float)(&output[0]), C.vDSP_Length(len(output)), C.int(flag))
}

// Creates a single-precision Hamming window.
func HammWindow(output []float32, flag WindowFlag) {
	C.vDSP_hamm_window((*C.float)(&output[0]), C.vDSP_Length(len(output)), C.int(flag))
}

// Creates a single-precision Blackman window.
func BlkmanWindow(output []float32, flag WindowFlag) {
	C.vDSP_blkman_window((*C.float)(&output[0]), C.vDSP_Length(len(output)), C.int(flag))
}
