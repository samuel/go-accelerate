package accel

// #include <Accelerate/Accelerate.h>
import "C"

// Calculates histograms for each channel of an ARGB8888 image.
func VImageHistogramCalculation_ARGB8888(src *VImageBuffer, flags VImageFlag) ([4][]int, error) {
	srcC := src.toC()
	var hist [4][256]C.vImagePixelCount
	var histPtrs [4]*C.vImagePixelCount
	for i := 0; i < 4; i++ {
		histPtrs[i] = &hist[i][0]
	}
	if err := toError(C.vImageHistogramCalculation_ARGB8888(&srcC, &histPtrs[0], C.vImage_Flags(flags))); err != nil {
		return [4][]int{}, err
	}
	var outHist [4][]int
	for i, h := range hist {
		outHist[i] = make([]int, 256)
		for j, count := range h {
			outHist[i][j] = int(count)
		}
	}
	return outHist, nil
}

// Calculates a histogram for a Planar8 image.
func VImageHistogramCalculation_Planar8(src *VImageBuffer, flags VImageFlag) ([]int, error) {
	srcC := src.toC()
	var hist [256]C.vImagePixelCount
	if err := toError(C.vImageHistogramCalculation_Planar8(&srcC, &hist[0], C.vImage_Flags(flags))); err != nil {
		return nil, err
	}
	outHist := make([]int, 256)
	for i, count := range hist {
		outHist[i] = int(count)
	}
	return outHist, nil
}
