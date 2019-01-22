package main

import (
	"github.com/racerxdl/gorrect/correctwrap"
)

type ReedSolomon struct {
	blockSize int
	dataSize  int
	poly      uint16

	rs        correctwrap.Correct_reed_solomon
	canBeUsed bool
}

func MakeReedSolomon(blockSize, dataSize, minDistance int, poly uint16) *ReedSolomon {
	if minDistance == -1 {
		minDistance = blockSize - dataSize
	}
	return &ReedSolomon{
		blockSize: blockSize,
		dataSize:  dataSize,
		poly:      poly,
		rs:        correctwrap.Correct_reed_solomon_create(poly, 1, 1, int64(minDistance)),
		canBeUsed: true,
	}
}

func (rs *ReedSolomon) Encode(data []byte) (outData []byte) {
	outData = make([]byte, rs.blockSize)
	correctwrap.Correct_reed_solomon_encode(rs.rs, &data[0], int64(len(data)), &outData[0])
	return outData
}

func (rs *ReedSolomon) Decode(data []byte) (outData []byte, errors int) {
	if !rs.canBeUsed {
		panic("Reed Solomon instance used after Close!!!")
	}
	outData = make([]byte, rs.dataSize)
	correctwrap.Correct_reed_solomon_decode(rs.rs, &data[0], int64(len(data)), &outData[0])

	errors = 0

	for i := 0; i < rs.dataSize; i++ {
		if outData[i] != data[i] {
			errors++
		}
	}

	return outData, errors
}

func (rs *ReedSolomon) Close() {
	correctwrap.Correct_reed_solomon_destroy(rs.rs)
	rs.canBeUsed = false
}
