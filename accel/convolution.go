package accel

// #include <Accelerate/Accelerate.h>
import "C"

import (
	"unsafe"
)

// Sharpens an ARGBFFFF image by undoing a previous convolution that blurred the image, such as diffraction effects in a camera lens.
func VImageRichardsonLucyDeConvolve_ARGBFFFF(src, dst *VImageBuffer, tempBuffer []byte, roiX, roiY int, kernel, kernel2 []float32, kernelHeight, kernelWidth, kernelHeight2, kernelWidth2 int, backgroundColor [4]float32, iterationCount int, flags VImageFlag) error {
	var tmpBuf unsafe.Pointer
	if tempBuffer != nil {
		tmpBuf = unsafe.Pointer(&tempBuffer[0])
	}
	var kernel2Ptr *C.float
	if kernel2 != nil {
		kernel2Ptr = (*C.float)(&kernel2[0])
	}
	srcC := src.toC()
	dstC := dst.toC()
	return toError(C.vImageRichardsonLucyDeConvolve_ARGBFFFF(&srcC, &dstC, tmpBuf, C.vImagePixelCount(roiX),
		C.vImagePixelCount(roiY), (*C.float)(&kernel[0]), kernel2Ptr, C.uint32_t(kernelHeight),
		C.uint32_t(kernelWidth), C.uint32_t(kernelHeight2), C.uint32_t(kernelWidth2), (*C.float)(&backgroundColor[0]),
		C.uint32_t(iterationCount), C.vImage_Flags(flags)))
}

// Sharpens an ARGB8888 image by undoing a previous convolution that blurred the image, such as diffraction effects in a camera lens.
func VImageRichardsonLucyDeConvolve_ARGB8888(src, dst *VImageBuffer, tempBuffer []byte, roiX, roiY int, kernel, kernel2 []int16, kernelHeight, kernelWidth, kernelHeight2, kernelWidth2, divisor, divisor2 int, backgroundColor [4]uint8, iterationCount int, flags VImageFlag) error {
	var tmpBuf unsafe.Pointer
	if tempBuffer != nil {
		tmpBuf = unsafe.Pointer(&tempBuffer[0])
	}
	var kernel2Ptr *C.int16_t
	if kernel2 != nil {
		kernel2Ptr = (*C.int16_t)(&kernel2[0])
	}
	srcC := src.toC()
	dstC := dst.toC()
	return toError(C.vImageRichardsonLucyDeConvolve_ARGB8888(&srcC, &dstC, tmpBuf, C.vImagePixelCount(roiX),
		C.vImagePixelCount(roiY), (*C.int16_t)(&kernel[0]), kernel2Ptr, C.uint32_t(kernelHeight),
		C.uint32_t(kernelWidth), C.uint32_t(kernelHeight2), C.uint32_t(kernelWidth2), C.int32_t(divisor),
		C.int32_t(divisor2), (*C.uint8_t)(&backgroundColor[0]), C.uint32_t(iterationCount), C.vImage_Flags(flags)))
}

// Convolves a region of interest within a source image by an M x N kernel, then divides the pixel values by a divisor.
func VImageConvolve_ARGB8888(src, dst *VImageBuffer, tempBuffer []byte, roiX, roiY int, kernel []int16, kernelHeight, kernelWidth, divisor int, backgroundColor [4]uint8, flags VImageFlag) error {
	var tmpBuf unsafe.Pointer
	if tempBuffer != nil {
		tmpBuf = unsafe.Pointer(&tempBuffer[0])
	}
	srcC := src.toC()
	dstC := dst.toC()
	return toError(C.vImageConvolve_ARGB8888(&srcC, &dstC, tmpBuf, C.vImagePixelCount(roiX),
		C.vImagePixelCount(roiY), (*C.int16_t)(&kernel[0]), C.uint32_t(kernelHeight),
		C.uint32_t(kernelWidth), C.int32_t(divisor), (*C.uint8_t)(&backgroundColor[0]), C.vImage_Flags(flags)))
}
