package accel

// #include <Accelerate/Accelerate.h>
import "C"

// Performs nonpremultiplied alpha compositing of two ARGB8888 images, placing the result in a destination buffer.
func VImageAlphaBlend_ARGB8888(srcTop, srcBottom, dst *VImageBuffer, flags VImageFlag) error {
	srcTopC := srcTop.toC()
	srcBottomC := srcBottom.toC()
	dstC := dst.toC()
	return toError(C.vImageAlphaBlend_ARGB8888(&srcTopC, &srcBottomC, &dstC, C.vImage_Flags(flags)))
}

// Performs premultiplied alpha compositing of two ARGB8888 images, using a single alpha value for the whole image and placing the result in a destination buffer.
func VImagePremultipliedConstAlphaBlend_ARGB8888(srcTop *VImageBuffer, constAlpha uint8, srcBottom, dst *VImageBuffer, flags VImageFlag) error {
	srcTopC := srcTop.toC()
	srcBottomC := srcBottom.toC()
	dstC := dst.toC()
	return toError(C.vImagePremultipliedConstAlphaBlend_ARGB8888(&srcTopC, C.Pixel_8(constAlpha), &srcBottomC, &dstC, C.vImage_Flags(flags)))
}
