package accel

// #include <Accelerate/Accelerate.h>
import "C"

// Vvlog10f performs a log base 10 on every value in input and
// writes the result into output.
func Vvlog10f(output, input []float32) {
	n := C.int(len(output))
	C.vvlog10f((*C.float)(&output[0]), (*C.float)(&input[0]), &n)
}
