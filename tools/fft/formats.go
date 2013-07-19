package main

import (
	"encoding/binary"
	"math"

	"github.com/samuel/go-accelerate"
)

type sampler interface {
	Transform(buf []byte, data accel.DSPSplitComplex)
	Description() string
	SampleSize() int // Return the size of a sample in bytes
}

var sampleFormats = map[string]sampler{
	"8uc":    complexU8Sampler(1),
	"le16s":  realLES16Sampler(1),
	"le16sc": complexS16Sampler(1),
	"32fc":   complexF32Sampler(1),
	"le32fc": complexLEF32Sampler(1),
	"64fc":   complexF64Sampler(1),
	"le64fc": complexLEF64Sampler(1),
}

type complexU8Sampler int

func (smp complexU8Sampler) Transform(buf []byte, data accel.DSPSplitComplex) {
	// for i := 0; i < len(data.Real); i++ {
	// 	j := i * 2
	// 	data.Real[i] = float32(buf[j]) - 128.0
	// 	data.Imag[i] = float32(buf[j+1]) - 128.0
	// }
	accel.Vfltu8(buf, data.Real, 2, 1)
	accel.Vfltu8(buf[1:], data.Imag, 2, 1)
	accel.Vsadd(data.Real, 1, -128.0, data.Real, 1)
	accel.Vsadd(data.Imag, 1, -128.0, data.Imag, 1)
}

func (smp complexU8Sampler) SampleSize() int {
	return 2
}

func (smp complexU8Sampler) Description() string {
	return "8-bit unsigned complex (interleaved)"
}

type complexS16Sampler int

func (smp complexS16Sampler) Transform(buf []byte, data accel.DSPSplitComplex) {
	for i := 0; i < len(data.Real); i++ {
		j := i * 4
		data.Real[i] = float32(int16(int(buf[j]) | (int(buf[j+1]) << 8)))
		data.Imag[i] = float32(int16(int(buf[j+2]) | (int(buf[j+3]) << 8)))
	}
}

func (smp complexS16Sampler) SampleSize() int {
	return 4
}

func (smp complexS16Sampler) Description() string {
	return "Little-endian 16-bit signed (byte-128) complex (interleaved)"
}

type complexLEF32Sampler int

func (smp complexLEF32Sampler) Transform(buf []byte, data accel.DSPSplitComplex) {
	for i := 0; i < len(data.Real); i++ {
		j := i * 8
		data.Real[i] = math.Float32frombits(binary.LittleEndian.Uint32(buf[j : j+4]))
		data.Imag[i] = math.Float32frombits(binary.LittleEndian.Uint32(buf[j+4 : j+8]))
	}
}

func (smp complexLEF32Sampler) SampleSize() int {
	return 8
}

func (smp complexLEF32Sampler) Description() string {
	return "Little-endian 32-bit float complex (interleaved)"
}

type complexF32Sampler int

func (smp complexF32Sampler) Transform(buf []byte, data accel.DSPSplitComplex) {
	accel.Ctoz_byte(buf, data, 2, 1)
}

func (smp complexF32Sampler) SampleSize() int {
	return 8
}

func (smp complexF32Sampler) Description() string {
	return "Native-endian 32-bit float complex (interleaved)"
}

type complexLEF64Sampler int

func (smp complexLEF64Sampler) Transform(buf []byte, data accel.DSPSplitComplex) {
	for i := 0; i < len(data.Real); i++ {
		j := i * 16
		data.Real[i] = float32(math.Float64frombits(binary.LittleEndian.Uint64(buf[j : j+8])))
		data.Imag[i] = float32(math.Float64frombits(binary.LittleEndian.Uint64(buf[j+8 : j+16])))
	}
}

func (smp complexLEF64Sampler) SampleSize() int {
	return 16
}

func (smp complexLEF64Sampler) Description() string {
	return "Little-endian 64-bit float complex (interleaved)"
}

type complexF64Sampler int

func (smp complexF64Sampler) Transform(buf []byte, data accel.DSPSplitComplex) {
	accel.Vdpsp_byte(buf, data.Real, 2, 1)
	accel.Vdpsp_byte(buf[8:], data.Imag, 2, 1)
}

func (smp complexF64Sampler) SampleSize() int {
	return 16
}

func (smp complexF64Sampler) Description() string {
	return "Native-endian 64-bit float complex (interleaved)"
}

type realLES16Sampler int

// TODO: this is skipping one channel of a stereo pair. should get in "stride" from command like arguments
func (smp realLES16Sampler) Transform(buf []byte, data accel.DSPSplitComplex) {
	n := len(buf)
	for i := 0; i < len(data.Real); i++ {
		j := i * 4
		if j < n-1 {
			data.Real[i] = float32(int16(int(buf[j]) | (int(buf[j+1]) << 8)))
			j += 2
		} else {
			data.Real[i] = 0.0
		}
		data.Imag[i] = 0.0
	}
}

func (smp realLES16Sampler) SampleSize() int {
	return 4
}

func (smp realLES16Sampler) Description() string {
	return "Little-endian 16-bit signed real"
}
