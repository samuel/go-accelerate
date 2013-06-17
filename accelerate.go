package accel

// #include <Accelerate/Accelerate.h>
// #cgo LDFLAGS: -framework Accelerate
import "C"

// Reorders the channels in an ARGB8888 image.
//
// permuteMap:
//     An array of four 8-bit integers with the values 0, 1, 2, and 3,
//     in some order. Each value specifies a plane from the source image
//     that should be copied to that plane in the destination image. 0
//     denotes the alpha channel, 1 the red channel, 2 the green channel,
//     and 3 the blue channel. The following figure shows the result of
//     using a permute map shows values are (0, 3, 2, 1). The data in the
//     alpha and green channels remain the same, but the data in the source
//     red channel maps to the destination blue channel while the data in
//     the source blue channel maps to the destination red channel.
func VImagePermuteChannels_ARGB8888(src, dst *VImageBuffer, permuteMap [4]uint8, flags VImageFlag) error {
	srcC := src.toC()
	dstC := dst.toC()
	return toError(C.vImagePermuteChannels_ARGB8888(&srcC, &dstC, (*C.uint8_t)(&permuteMap[0]), C.vImage_Flags(flags)))
}
