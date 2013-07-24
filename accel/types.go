package accel

// #include <Accelerate/Accelerate.h>
import "C"

import (
	"errors"
	"fmt"
	"image"
	"unsafe"
)

type ErrOther int

func (e ErrOther) Error() string {
	return fmt.Sprintf("accel: error %d", int(e))
}

var (
	// The region of interest, as specified by the srcOffsetToROI_X and
	// srcOffsetToROI_Y parameters and the height and width of the
	// destination buffer, extends beyond the bottom edge or right edge
	// of the source buffer.
	ErrImageRoiLargerThanInputBuffer = errors.New("accel: Image RIO larger than input buffer")
	// Either the kernel height, the kernel width, or both, are even.
	ErrImageInvalidKernelSize = errors.New("accel: Invalid kernel size")
	// The edge style specified is invalid. This usually means that a
	// particular function requires you to set at least one edge option
	// flag (kvImageCopyInPlace, kvImageBackgroundColorFill, or kvImageEdgeExtend),
	// but you did not specify one.
	ErrImageInvalidEdgeStyle = errors.New("accel: Invalid edge style")
	// The srcOffsetToROI_X parameter that specifies the left edge of the
	// region of interest is greater than the width of the source image.
	ErrImageInvalidOffsetX = errors.New("accel: Invalid offset X")
	// The srcOffsetToROI_Y parameter that specifies the top edge of the
	// region of interest is greater than the height of the source image.
	ErrImageInvalidOffsetY = errors.New("accel: Invalid offset Y")
	// An attempt to allocate memory failed.
	ErrImageMemoryAllocationError = errors.New("accel: Memory allocation error")
	// A pointer parameter is NULL and it must not be.
	ErrImageNullPointerArgument = errors.New("accel: Null pointer argument")
	// Invalid parameter.
	ErrImageInvalidParameter = errors.New("accel: Invalid parameter")
	// The function requires the source and destination buffers to have the
	// same height and the same width, but they do not.
	ErrImageBufferSizeMismatch = errors.New("accel: Buffer size mismatch")
	// The flag is not recognized.
	ErrImageUnknownFlagsBit = errors.New("accel: Unknown flag bits")
)

func toError(code C.vImage_Error) error {
	switch code {
	case C.kvImageNoError:
		return nil
	case C.kvImageRoiLargerThanInputBuffer:
		return ErrImageRoiLargerThanInputBuffer
	case C.kvImageInvalidKernelSize:
		return ErrImageInvalidKernelSize
	case C.kvImageInvalidEdgeStyle:
		return ErrImageInvalidEdgeStyle
	case C.kvImageInvalidOffset_X:
		return ErrImageInvalidOffsetX
	case C.kvImageInvalidOffset_Y:
		return ErrImageInvalidOffsetY
	case C.kvImageMemoryAllocationError:
		return ErrImageMemoryAllocationError
	case C.kvImageNullPointerArgument:
		return ErrImageNullPointerArgument
	case C.kvImageInvalidParameter:
		return ErrImageInvalidParameter
	case C.kvImageBufferSizeMismatch:
		return ErrImageBufferSizeMismatch
	case C.kvImageUnknownFlagsBit:
		return ErrImageUnknownFlagsBit
	}
	return ErrOther(int(code))
}

type VImageFlag int

const (
	// Do not set any flags.
	VImageFlagNoFlags VImageFlag = C.kvImageNoFlags
	// Operate on red, green, and blue channels only. When you set this flag,
	// the alpha value is copied from source to destination. You can set this
	// flag only for interleaved image formats.
	VImageFlagLeaveAlphaUnchanged VImageFlag = C.kvImageLeaveAlphaUnchanged
	// Copy the value of the edge pixel in the source to the destination. When
	// you set this flag, and a convolution function is processing an image
	// pixel for which some of the kernel extends beyond the image boundaries,
	// vImage does not computer the convolution. Instead, the value of the
	// particular pixel unchanged; it’s simply copied to the destination image.
	// This flag is valid only for convolution operations. The morphology
	// functions do not use this flag because they do not use pixels outside
	// the image in any of their calculations.
	VImageFlagCopyInPlace VImageFlag = C.kvImageCopyInPlace
	// A background color fill. The associated value is a background color
	// (that is, a pixel value). When you set this flag, vImage assigns the
	// pixel value to all pixels outside the image. You can set this flag
	// for convolution and geometry functions. The morphology functions do
	// not use this flag because they do not use pixels outside the image
	// in any of their calculations.
	VImageFlagBackgroundColorFill VImageFlag = C.kvImageBackgroundColorFill
	// Extend the edges of the image infinitely. When you set this flag,
	// vImage replicates the edges of the image outward. It repeats the top
	// row of the image infinitely above the image, the bottom row infinitely
	// below the image, the first column infinitely to the left of the image,
	// and the last column infinitely to the right. For spaces that are
	// diagonal to the image, vImage uses the value of the corresponding
	// corner pixel. For example, for all pixels that are both above and to
	// the left of the image, the upper-leftmost pixel of the image (the pixel
	// at row 0, column 0) supplies the value. In this way, vImage assigns
	// every pixel location outside the image some value. You can set this
	// flag for convolution and geometry functions. The morphology functions
	// do not use this flag because they do not use pixels outside the image
	// in any of their calculations.
	VImageFlagEdgeExtend VImageFlag = C.kvImageEdgeExtend
	// Do not use vImage internal tiling routines. When you set this flag,
	// vImage turns off internal tiling. Set this flag if you want to perform
	// your own tiling or your own multithreading, or to use the minimum or
	// maximum filters in place.
	VImageFlagDoNotTile VImageFlag = C.kvImageDoNotTile
	// Use a higher quality, slower resampling filter for for geometry
	// operations—shear, scale, rotate, affine transform, and so forth.
	VImageFlagHighQualityResampling VImageFlag = C.kvImageHighQualityResampling
	// Use the part of the kernel that overlaps the image. This flag is valid
	// only for convolution operations. When you set this flag, vImage restricts
	// calculations to the portion of the kernel overlapping the image. It
	// corrects the calculated pixel by first multiplying by the sum of all the
	// kernel elements, then dividing by the sum of the kernel elements that are
	// actually used. This preserves image brightness at the edges.
	//
	// For integer kernels:
	//   real_divisor = divisor * (sum of used kernel elements) / (sum of kernel elements)
	//
	// The morphology functions do not use this flag because they do not use
	// pixels outside the image in any of their calculations.
	// Kernel truncation is not robust for certain kernels. It can ail when any
	// rectangular segment of the kernel that includes the center, and at least
	// one of the corners, sums to zero. You typically see this with emboss or
	// edge detection filters, or other filters that are designed to find the
	// slope of a signal. For those kinds of filters, you should use the
	// kvImageEdgeExtend option instead.
	VImageFlagTruncateKernel VImageFlag = C.kvImageTruncateKernel
	// Get the minimum temporary buffer size for the operation, given the
	// parameters provided. When you set this flag, the function returns the
	// number of bytes required for the temporary buffer. A negative value
	// specifies an error.
	VImageFlagGetTempBufferSize VImageFlag = C.kvImageGetTempBufferSize
)

type VImageBuffer struct {
	Width, Height int
	RowBytes      int
	Data          []byte
}

func (vib *VImageBuffer) toC() C.vImage_Buffer {
	var cv C.vImage_Buffer
	cv.data = unsafe.Pointer(&vib.Data[0])
	cv.width = C.vImagePixelCount(vib.Width)
	cv.height = C.vImagePixelCount(vib.Height)
	cv.rowBytes = C.size_t(vib.RowBytes)
	return cv
}

// Return a VImageBuffer of the given image. The memory may or may not
// be shared depending on the format. Also return the format of the returned
// image. (e.g. argb8888, rgba8888, 8, ...)
func VImageBufferFromImage(img image.Image) (*VImageBuffer, string) {
	switch m := img.(type) {
	case *image.Gray:
		return &VImageBuffer{Width: m.Bounds().Dx(), Height: m.Bounds().Dy(), RowBytes: m.Stride, Data: m.Pix}, "8"
	case *image.RGBA:
		return &VImageBuffer{Width: m.Bounds().Dx(), Height: m.Bounds().Dy(), RowBytes: m.Stride, Data: m.Pix}, "rgba8888"
	}

	b := img.Bounds()
	w := b.Dx()
	h := b.Dy()
	data := make([]byte, w*h*4)
	dataOffset := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := img.At(x+b.Min.X, y+b.Min.Y)
			r, g, b, a := c.RGBA()
			data[dataOffset] = uint8(a >> 8)
			data[dataOffset+1] = uint8(r >> 8)
			data[dataOffset+2] = uint8(g >> 8)
			data[dataOffset+3] = uint8(b >> 8)
			dataOffset += 4
		}
	}
	return &VImageBuffer{Width: w, Height: h, RowBytes: w * 4, Data: data}, "argb8888"
}

// Return an allocated vImage_Buffer of the given dimensions and type.
// rowBytes may be 0 in which case it will be calculated as width*channels.
func CreateVImageBuffer(width, height, channels, rowBytes int) *VImageBuffer {
	if rowBytes <= 0 {
		rowBytes = width * height * channels
	} else if rowBytes < width*height*channels {
		panic("accel: trying to create a buffer with an invalid rowBytes size")
	}
	return &VImageBuffer{
		Data:     make([]byte, rowBytes*height),
		Width:    width,
		Height:   height,
		RowBytes: rowBytes,
	}
}

// Return an instance of *image.RGBA that shares the data with the VImageBuffer.
// No checks are done to guarantee the format matches.
func (b *VImageBuffer) ToRGBA() *image.RGBA {
	return &image.RGBA{
		Pix:    b.Data,
		Stride: b.RowBytes,
		Rect:   image.Rect(0, 0, b.Width, b.Height),
	}
}
