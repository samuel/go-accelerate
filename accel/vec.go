package accel

// #include <Accelerate/Accelerate.h>
import "C"

// For each single-precision array element, sets y to the base 10 logarithm of x.
func Vvlog10f(output []float32, input []float32) {
	var n C.int = C.int(len(output))
	C.vvlog10f((*C.float)(&output[0]), (*C.float)(&input[0]), &n)
}
