package main

import (
	"encoding/binary"
	"io"
	"math"

	"github.com/samuel/go-accelerate"
)

type samplerMaker func() sampler

type sampler interface {
	Read(rd io.Reader, data accel.DSPSplitComplex) error
	Description() string
	SampleSize() int // Return the size of a sample in bytes
}

var sampleFormats = map[string]samplerMaker{
	"8uc":    newComplexU8Sampler,
	"le16s":  newRealS16Sampler,
	"le16sc": newComplexS16Sampler,
	"le32fc": newComplexF32Sampler,
	"le64fc": newComplexF64Sampler,
}

type complexU8Sampler struct {
	buf []byte
}

func newComplexU8Sampler() sampler {
	return &complexU8Sampler{}
}

func (smp *complexU8Sampler) Read(rd io.Reader, data accel.DSPSplitComplex) error {
	if smp.buf == nil && len(smp.buf) < len(data.Real)*2 {
		smp.buf = make([]byte, len(data.Real)*2)
	}
	n, err := rd.Read(smp.buf[:len(data.Real)*2])
	if err != nil {
		return err
	}
	for i := 0; i < len(data.Real); i++ {
		j := i * 2
		if j < n {
			data.Real[i] = float32(smp.buf[j]) - 128.0
			j++
		} else {
			data.Real[i] = 0.0
		}
		if j < n {
			data.Imag[i] = float32(smp.buf[j]) - 128.0
		} else {
			data.Imag[i] = 0.0
		}
	}
	return nil
}

func (smp *complexU8Sampler) SampleSize() int {
	return 2
}

func (smp *complexU8Sampler) Description() string {
	return "8-bit unsigned complex (interleaved)"
}

type complexS16Sampler struct {
	buf []byte
}

func newComplexS16Sampler() sampler {
	return &complexS16Sampler{}
}

func (smp *complexS16Sampler) Read(rd io.Reader, data accel.DSPSplitComplex) error {
	if smp.buf == nil && len(smp.buf) < len(data.Real)*4 {
		smp.buf = make([]byte, len(data.Real)*4)
	}
	n, err := rd.Read(smp.buf[:len(data.Real)*4])
	if err != nil {
		return err
	}
	for i := 0; i < len(data.Real); i++ {
		j := i * 4
		if j < n-1 {
			data.Real[i] = float32(int16(int(smp.buf[j]) | (int(smp.buf[j+1]) << 8)))
			j += 2
		} else {
			data.Real[i] = 0.0
		}
		if j < n-1 {
			data.Imag[i] = float32(int16(int(smp.buf[j]) | (int(smp.buf[j+1]) << 8)))
		} else {
			data.Imag[i] = 0.0
		}
	}
	return nil
}

func (smp *complexS16Sampler) SampleSize() int {
	return 4
}

func (smp *complexS16Sampler) Description() string {
	return "Little-endian 16-bit signed complex (interleaved)"
}

type complexF32Sampler struct {
	buf []byte
}

func newComplexF32Sampler() sampler {
	return &complexF32Sampler{}
}

func (smp *complexF32Sampler) Read(rd io.Reader, data accel.DSPSplitComplex) error {
	if smp.buf == nil && len(smp.buf) < len(data.Real)*8 {
		smp.buf = make([]byte, len(data.Real)*8)
	}
	n, err := rd.Read(smp.buf[:len(data.Real)*8])
	if err != nil {
		return err
	}
	for i := 0; i < len(data.Real); i++ {
		j := i * 8
		if j <= n-4 {
			data.Real[i] = math.Float32frombits(binary.LittleEndian.Uint32(smp.buf[j : j+4]))
			j += 4
		} else {
			data.Real[i] = 0.0
		}
		if j <= n-4 {
			data.Imag[i] = math.Float32frombits(binary.LittleEndian.Uint32(smp.buf[j : j+4]))
		} else {
			data.Imag[i] = 0.0
		}
	}
	return nil
}

func (smp *complexF32Sampler) SampleSize() int {
	return 8
}

func (smp *complexF32Sampler) Description() string {
	return "Little-endian 32-bit float complex (interleaved)"
}

type complexF64Sampler struct {
	buf []byte
}

func newComplexF64Sampler() sampler {
	return &complexF64Sampler{}
}

func (smp *complexF64Sampler) Read(rd io.Reader, data accel.DSPSplitComplex) error {
	if smp.buf == nil && len(smp.buf) < len(data.Real)*16 {
		smp.buf = make([]byte, len(data.Real)*16)
	}
	n, err := rd.Read(smp.buf[:len(data.Real)*16])
	if err != nil {
		return err
	}
	for i := 0; i < len(data.Real); i++ {
		j := i * 16
		if j <= n-8 {
			data.Real[i] = float32(math.Float64frombits(binary.LittleEndian.Uint64(smp.buf[j : j+8])))
			j += 8
		} else {
			data.Real[i] = 0.0
		}
		if j <= n-8 {
			data.Imag[i] = float32(math.Float64frombits(binary.LittleEndian.Uint64(smp.buf[j : j+8])))
		} else {
			data.Imag[i] = 0.0
		}
	}
	return nil
}

func (smp *complexF64Sampler) SampleSize() int {
	return 16
}

func (smp *complexF64Sampler) Description() string {
	return "Little-endian 64-bit float complex (interleaved)"
}

type realS16Sampler struct {
	buf []byte
}

func newRealS16Sampler() sampler {
	return &realS16Sampler{}
}

// TODO: this is skipping one channel of a stereo pair. should get in "stride" from command like arguments
func (smp *realS16Sampler) Read(rd io.Reader, data accel.DSPSplitComplex) error {
	if smp.buf == nil && len(smp.buf) < len(data.Real)*4 {
		smp.buf = make([]byte, len(data.Real)*4)
	}
	n, err := rd.Read(smp.buf[:len(data.Real)*4])
	if err != nil {
		return err
	}
	for i := 0; i < len(data.Real); i++ {
		j := i * 4
		if j < n-1 {
			data.Real[i] = float32(int16(int(smp.buf[j]) | (int(smp.buf[j+1]) << 8)))
			j += 2
		} else {
			data.Real[i] = 0.0
		}
		data.Imag[i] = 0.0
	}
	return nil
}

func (smp *realS16Sampler) SampleSize() int {
	return 4
}

func (smp *realS16Sampler) Description() string {
	return "Little-endian 16-bit signed real"
}
