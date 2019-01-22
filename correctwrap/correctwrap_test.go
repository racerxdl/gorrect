package correctwrap

import (
	"github.com/racerxdl/gorrect/Codes"
	"math/rand"
	"testing"
)

func TestReedSolomon(t *testing.T) {
	blockLength := 255
	minDistance := 32
	messageLength := blockLength - minDistance

	rs := Correct_reed_solomon_create(Codes.ReedSolomonPrimitivePolynomialCCSDS, 112, 11, int64(minDistance))

	data := make([]byte, messageLength)
	originalData := make([]byte, messageLength)
	encoded := make([]byte, blockLength)

	for i := 0; i < messageLength; i++ {
		data[i] = byte(rand.Int31() & 0xFF)
	}

	copy(originalData, data)

	Correct_reed_solomon_encode(rs, &data[0], int64(messageLength), &encoded[0])

	for i := 0; i < 8; i++ {
		v := rand.Int31n(int32(messageLength))
		encoded[v] = 0x00 // Add Error
	}

	data = make([]byte, messageLength)

	Correct_reed_solomon_decode(rs, &encoded[0], int64(blockLength), &data[0])

	for i := 0; i < messageLength; i++ {
		if data[i] != originalData[i] {
			t.Fatalf("Error on position %d. Expected %d got %d", i, originalData[i], data[i])
			t.FailNow()
		}
	}
	Correct_reed_solomon_destroy(rs)
}
