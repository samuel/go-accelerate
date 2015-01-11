package accel

// #include <Accelerate/Accelerate.h>
import "C"

import "unsafe"

type DSPSplitComplex struct {
	Real []float32
	Imag []float32
}

type DSPDoubleSplitComplex struct {
	Real []float64
	Imag []float64
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

func minLen(size ...int) C.vDSP_Length {
	min := size[0]
	for i := 1; i < len(size); i++ {
		if size[i] < min {
			min = size[i]
		}
	}
	return C.vDSP_Length(min)
}

// Vflt8 converts an array of signed 8-bit integers to single-precision floating-point values.
func Vflt8(input []int8, inputStride int, output []float32, outputStride int) {
	C.vDSP_vflt8((*C.char)(&input[0]), C.vDSP_Stride(inputStride), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), minLen(len(input)/inputStride, len(output)/outputStride))
}

// Vflt8_byte converts an array of signed 8-bit integers to single-precision floating-point values.
func Vflt8_byte(input []byte, inputStride int, output []float32, outputStride int) {
	C.vDSP_vflt8((*C.char)(unsafe.Pointer(&input[0])), C.vDSP_Stride(inputStride), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), minLen(len(input)/inputStride, len(output)/outputStride))
}

// Vfltu8 converts an array of unsigned 8-bit integers to single-precision floating-point values.
func Vfltu8(input []byte, inputStride int, output []float32, outputStride int) {
	C.vDSP_vfltu8((*C.uchar)(&input[0]), C.vDSP_Stride(inputStride), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), minLen(len(input)/inputStride, len(output)/outputStride))
}

// Vflt16 converts an array of signed 16-bit integers to single-precision floating-point values.
func Vflt16(input []int16, inputStride int, output []float32, outputStride int) {
	C.vDSP_vflt16((*C.short)(&input[0]), C.vDSP_Stride(inputStride), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), minLen(len(input)/inputStride, len(output)/outputStride))
}

// Vflt16_byte converts an array of signed 16-bit integers to single-precision floating-point values.
func Vflt16_byte(input []byte, inputStride int, output []float32, outputStride int) {
	C.vDSP_vflt16((*C.short)(unsafe.Pointer(&input[0])), C.vDSP_Stride(inputStride), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), minLen(len(input)/(2*inputStride), len(output)/outputStride))
}

// Vflt32 converts an array of signed 32-bit integers to single-precision floating-point values.
func Vflt32(input []int32, inputStride int, output []float32, outputStride int) {
	C.vDSP_vflt32((*C.int)(&input[0]), C.vDSP_Stride(inputStride), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), minLen(len(input)/inputStride, len(output)/outputStride))
}

// Vflt32_byte converts an array of signed 16-bit integers to single-precision floating-point values.
func Vflt32_byte(input []byte, inputStride int, output []float32, outputStride int) {
	C.vDSP_vflt32((*C.int)(unsafe.Pointer(&input[0])), C.vDSP_Stride(inputStride), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), minLen(len(input)/(4*inputStride), len(output)/outputStride))
}

// Vdpsp convert a double-precision vector to single-precision.
func Vdpsp(input []float64, inputStride int, output []float32, outputStride int) {
	C.vDSP_vdpsp((*C.double)(&input[0]), C.vDSP_Stride(inputStride), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), minLen(len(input)/inputStride, len(output)/outputStride))
}

// Vdpsp_byte converts a double-precision to single-precision.
// Operate on a byte buffer which contains float64
func Vdpsp_byte(input []byte, inputStride int, output []float32, outputStride int) {
	C.vDSP_vdpsp((*C.double)(unsafe.Pointer(&input[0])), C.vDSP_Stride(inputStride), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), minLen(len(input)/8/inputStride, len(output)/outputStride))
}

// Ctoz copies the contents of an interleaved complex vector C to a split complex vector Z; single precision.
func Ctoz(input []complex64, inputStride int, output DSPSplitComplex, outputStride int) {
	var splitComplex C.DSPSplitComplex
	splitComplex.realp = (*C.float)(&output.Real[0])
	splitComplex.imagp = (*C.float)(&output.Imag[0])
	n := 2 * len(output.Real) / outputStride
	if n2 := 2 * len(input) / inputStride; n2 < n {
		n = n2
	}
	C.vDSP_ctoz((*C.DSPComplex)(unsafe.Pointer(&input[0])), C.vDSP_Stride(inputStride), &splitComplex, C.vDSP_Stride(outputStride), C.vDSP_Length(n))
}

func Ctoz_float(input []float32, inputStride int, output DSPSplitComplex, outputStride int) {
	var splitComplex C.DSPSplitComplex
	splitComplex.realp = (*C.float)(&output.Real[0])
	splitComplex.imagp = (*C.float)(&output.Imag[0])
	n := 2 * len(output.Real) / outputStride
	if n2 := len(input) / inputStride; n2 < n {
		n = n2
	}
	C.vDSP_ctoz((*C.DSPComplex)(unsafe.Pointer(&input[0])), C.vDSP_Stride(inputStride), &splitComplex, C.vDSP_Stride(outputStride), C.vDSP_Length(n))
}

// Ctoz_byte copies the contents of an interleaved complex vector C to a split complex vector Z; single precision. Operate on a byte buffer which contains complex64
func Ctoz_byte(input []byte, inputStride int, output DSPSplitComplex, outputStride int) {
	var splitComplex C.DSPSplitComplex
	splitComplex.realp = (*C.float)(&output.Real[0])
	splitComplex.imagp = (*C.float)(&output.Imag[0])
	n := 2 * len(output.Real) / outputStride
	if n2 := len(input) / (4 * inputStride); n2 < n {
		n = n2
	}
	C.vDSP_ctoz((*C.DSPComplex)(unsafe.Pointer(&input[0])), C.vDSP_Stride(inputStride), &splitComplex, C.vDSP_Stride(outputStride), C.vDSP_Length(n))
}

// Ztoc copies the contents of a split complex vector Z to an interleaved complex vector C; single precision.
func Ztoc(input DSPSplitComplex, inputStride int, output []complex64, outputStride int) {
	var splitComplex C.DSPSplitComplex
	splitComplex.realp = (*C.float)(&input.Real[0])
	splitComplex.imagp = (*C.float)(&input.Imag[0])
	C.vDSP_ztoc(&splitComplex, C.vDSP_Stride(inputStride), (*C.DSPComplex)(unsafe.Pointer(&output[0])), C.vDSP_Stride(outputStride), C.vDSP_Length(2*len(output)/outputStride))
}

// Ztoc_float copies the contents of a split complex vector Z to an interleaved complex vector C; single precision.
func Ztoc_float(input DSPSplitComplex, inputStride int, output []float32, outputStride int) {
	var splitComplex C.DSPSplitComplex
	splitComplex.realp = (*C.float)(&input.Real[0])
	splitComplex.imagp = (*C.float)(&input.Imag[0])
	C.vDSP_ztoc(&splitComplex, C.vDSP_Stride(inputStride), (*C.DSPComplex)(unsafe.Pointer(&output[0])), C.vDSP_Stride(outputStride), C.vDSP_Length(len(output)/outputStride))
}

// Ztoc_byte copies the contents of a split complex vector Z to an interleaved complex vector C; single precision. Operate on a byte buffer which contains complex64
func Ztoc_byte(input DSPSplitComplex, inputStride int, output []byte, outputStride int) {
	var splitComplex C.DSPSplitComplex
	splitComplex.realp = (*C.float)(&input.Real[0])
	splitComplex.imagp = (*C.float)(&input.Imag[0])
	C.vDSP_ztoc(&splitComplex, C.vDSP_Stride(inputStride), (*C.DSPComplex)(unsafe.Pointer(&output[0])), C.vDSP_Stride(outputStride), C.vDSP_Length(len(output)/4/outputStride))
}

// Vclr clears the provided vector
func Vclr(vec []float32, stride int) {
	C.vDSP_vclr((*C.float)(&vec[0]), (C.vDSP_Stride)(stride), (C.vDSP_Length)(len(vec)/stride))
}

// Vfill fills the output vector with the provided value.
func Vfill(value float32, output []float32, stride int) {
	C.vDSP_vfill((*C.float)(&value), (*C.float)(&output[0]), (C.vDSP_Stride)(stride), (C.vDSP_Length)(len(output)/stride))
}

// Vclip clips the input vector using the given low and high and writes
// the result to the output vector.
func Vclip(input []float32, inputStride int, low, high float32, output []float32, outputStride int) {
	C.vDSP_vclip((*C.float)(&input[0]), C.vDSP_Stride(inputStride), (*C.float)(&low), (*C.float)(&high), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), minLen(len(input)/inputStride, len(output)/outputStride))
}

// Vthr thresholds the input vector writing the the result ot the output vector.
func Vthr(input []float32, inputStride int, low float32, output []float32, outputStride int) {
	C.vDSP_vthr((*C.float)(&input[0]), C.vDSP_Stride(inputStride), (*C.float)(&low), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), minLen(len(input)/inputStride, len(output)/outputStride))
}

// Desamp performs convolution with decimation.
func Desamp(input []float32, desamplingFactor int, coeff []float32, output []float32) {
	C.vDSP_desamp((*C.float)(&input[0]), C.vDSP_Stride(desamplingFactor), (*C.float)(&coeff[0]), (*C.float)(&output[0]), C.vDSP_Length(len(output)), C.vDSP_Length(len(coeff)))
}

// Zrdesamp performs a complex-real downsample with anti-aliasing.
func Zrdesamp(input DSPSplitComplex, decimationFactor int, coefficients []float32, output DSPSplitComplex) {
	var srcC C.DSPSplitComplex
	srcC.realp = (*C.float)(&input.Real[0])
	srcC.imagp = (*C.float)(&input.Imag[0])
	var dstC C.DSPSplitComplex
	dstC.realp = (*C.float)(&output.Real[0])
	dstC.imagp = (*C.float)(&output.Imag[0])
	C.vDSP_zrdesamp(&srcC, C.vDSP_Stride(decimationFactor), (*C.float)(&coefficients[0]), &dstC, C.vDSP_Length(len(output.Real)), C.vDSP_Length(len(coefficients)))
}

// Zvphas calculates the complex vector phase.
func Zvphas(input DSPSplitComplex, inputStride int, output []float32, outputStride int) {
	var srcC C.DSPSplitComplex
	srcC.realp = (*C.float)(&input.Real[0])
	srcC.imagp = (*C.float)(&input.Imag[0])
	C.vDSP_zvphas(&srcC, C.vDSP_Stride(inputStride), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), C.vDSP_Length(len(output)))
}

// Zidotpr calculates the conjugate dot product (or inner dot product) of complex vectors A and B and leave the result in complex vector C; single precision.
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

// Zvcmul performs a complex vector conjugate and multiply.
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
	C.vDSP_zvcmul(&in1, C.vDSP_Stride(stride1), &in2, C.vDSP_Stride(stride2), &res, C.vDSP_Stride(resultStride), minLen(len(input1.Real)/stride1, len(input2.Real)/stride2, len(result.Real)/resultStride))
}

// Vfix16 converts an array of single-precision floating-point values to signed 16-bit integer values, rounding towards zero.
func Vfix16(input []float32, inputStride int, output []int16, outputStride int) {
	C.vDSP_vfix16((*C.float)(&input[0]), C.vDSP_Stride(inputStride), (*C.short)(&output[0]), C.vDSP_Stride(outputStride), minLen(len(input)/inputStride, len(output)/outputStride))
}

// Vfix16_byte converts an array of single-precision floating-point values to signed 16-bit integer values, rounding towards zero. Output to a byte stream.
func Vfix16_byte(input []float32, inputStride int, output []byte, outputStride int) {
	C.vDSP_vfix16((*C.float)(&input[0]), C.vDSP_Stride(inputStride), (*C.short)(unsafe.Pointer(&output[0])), C.vDSP_Stride(outputStride), C.vDSP_Length(len(output)/2/outputStride))
}

// Vadd adds two vectors.
func Vadd(input1 []float32, input1Stride int, input2 []float32, input2Stride int, output []float32, outputStride int) {
	C.vDSP_vadd((*C.float)(&input1[0]), C.vDSP_Stride(input1Stride), (*C.float)(&input2[0]), C.vDSP_Stride(input2Stride), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), minLen(len(input1)/input1Stride, len(input2)/input2Stride, len(output)/outputStride))
}

// Vsmsa is vector scalar multiply and scalar add; single precision.
// output[n] = input[n] * mult + add
func Vsmsa(input []float32, inputStride int, mult, add float32, output []float32, outputStride int) {
	C.vDSP_vsmsa((*C.float)(&input[0]), C.vDSP_Stride(inputStride), (*C.float)(&mult), (*C.float)(&add), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), minLen(len(input)/inputStride, len(output)/outputStride))
}

// Vabs calcultes the absolute value of every value in the provided vector.
func Vabs(input []float32, inputStride int, output []float32, outputStride int) {
	C.vDSP_vabs((*C.float)(&input[0]), C.vDSP_Stride(inputStride), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), minLen(len(input)/inputStride, len(output)/outputStride))
}

// Computes the squared values of vector input and leaves the result in vector result; single precision.
func Vsq(input []float32, inputStride int, output []float32, outputStride int) {
	C.vDSP_vsq((*C.float)(&input[0]), C.vDSP_Stride(inputStride), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), minLen(len(input)/inputStride, len(output)/outputStride))
}

// Zvabs calculates the absolute values of all values in the complex input.
func Zvabs(input DSPSplitComplex, inputStride int, output []float32, outputStride int) {
	var in C.DSPSplitComplex
	in.realp = (*C.float)(&input.Real[0])
	in.imagp = (*C.float)(&input.Imag[0])
	C.vDSP_zvabs(&in, C.vDSP_Stride(inputStride), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), C.vDSP_Length(len(output)/outputStride))
}

// Complex vector absolute values; double precision.
func ZvabsD(input DSPDoubleSplitComplex, inputStride int, output []float64, outputStride int) {
	var in C.DSPDoubleSplitComplex
	in.realp = (*C.double)(&input.Real[0])
	in.imagp = (*C.double)(&input.Imag[0])
	C.vDSP_zvabsD(&in, C.vDSP_Stride(inputStride), (*C.double)(&output[0]), C.vDSP_Stride(outputStride), C.vDSP_Length(len(output)/outputStride))
}

// Vector maximum value; single precision.
func Maxv(input []float32, stride int) float32 {
	var out C.float
	C.vDSP_maxv((*C.float)(&input[0]), C.vDSP_Stride(stride), &out, C.vDSP_Length(len(input)/stride))
	return float32(out)
}

// Vector minimum value; single precision.
func Minv(input []float32, stride int) float32 {
	var out C.float
	C.vDSP_minv((*C.float)(&input[0]), C.vDSP_Stride(stride), &out, C.vDSP_Length(len(input)/stride))
	return float32(out)
}

// Vector sum; single precision.
func Sve(input []float32, inputStride int) float32 {
	var sum C.float
	C.vDSP_sve((*C.float)(&input[0]), C.vDSP_Stride(inputStride), &sum, C.vDSP_Length(len(input)/inputStride))
	return float32(sum)
}

// Vector linear average; single precision.
func Vavlin(input []float32, inputStride int, count float32, output []float32, outputStride int) {
	C.vDSP_vavlin((*C.float)(&input[0]), C.vDSP_Stride(inputStride), (*C.float)(&count), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), C.vDSP_Length(len(output)/outputStride))
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

// Vector sliding window sum; single precision.
func Vswsum(input []float32, inputStride int, output []float32, outputStride, windowLen int) {
	C.vDSP_vswsum((*C.float)(&input[0]), C.vDSP_Stride(inputStride), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), C.vDSP_Length(len(output)/outputStride), C.vDSP_Length(windowLen))
}

// Vector convert power or amplitude to decibels; single precision.
// α * log10(input(n)/zeroReference) [α is 20 if Amplitude (flag=1), or 10 if F is Power (flag=0)]
func Vdbcon(input []float32, inputStride int, zeroReference float32, output []float32, outputStride int, flag DBFlag) {
	C.vDSP_vdbcon((*C.float)(&input[0]), C.vDSP_Stride(inputStride), (*C.float)(&zeroReference), (*C.float)(&output[0]), C.vDSP_Stride(outputStride), C.vDSP_Length(len(output)/outputStride), C.uint(flag))
}

// Meanv returns the mean of the input vector.
func Meanv(input []float32, stride int) float32 {
	var mean C.float
	C.vDSP_meanv((*C.float)(&input[0]), C.vDSP_Stride(stride), &mean, C.vDSP_Length(len(input)/stride))
	return float32(mean)
}

// HannWindow creates a single-precision Hanning window.
func HannWindow(output []float32, flag WindowFlag) {
	C.vDSP_hann_window((*C.float)(&output[0]), C.vDSP_Length(len(output)), C.int(flag))
}

// HammWindow creates a single-precision Hamming window.
func HammWindow(output []float32, flag WindowFlag) {
	C.vDSP_hamm_window((*C.float)(&output[0]), C.vDSP_Length(len(output)), C.int(flag))
}

// BlkmanWindow creates a single-precision Blackman window.
func BlkmanWindow(output []float32, flag WindowFlag) {
	C.vDSP_blkman_window((*C.float)(&output[0]), C.vDSP_Length(len(output)), C.int(flag))
}
