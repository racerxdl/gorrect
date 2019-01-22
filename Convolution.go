package gorrect

import "github.com/racerxdl/gorrect/correctwrap"

type ConvolutionCoder struct {
	r    int
	k    int
	poly []uint16

	cc correctwrap.Correct_convolutional
}

// MakeConvolutionCoder creates a new Convolution Decoder / Encoder for r => rate, k => order and specified polys
func MakeConvolutionCoder(r, k int, poly []uint16) *ConvolutionCoder {
	if len(poly) != r {
		panic("The number of polys should match the rate")
	}

	return &ConvolutionCoder{
		r:    r,
		k:    k,
		poly: poly,
		cc:   correctwrap.Correct_convolutional_create(int64(r), int64(k), &poly[0]),
	}
}

// EncodedSize return number of encoded bits given dataLength bytes of data
func (cc *ConvolutionCoder) EncodedSize(dataLength int) int {
	return int(correctwrap.Correct_convolutional_encode_len(cc.cc, int64(dataLength)))
}

// DecodedSize return number of decoded bytes given numBits input data
func (cc *ConvolutionCoder) DecodedSize(numBits int) int {
	return numBits / (8 * cc.r)
}

// EncodeSoft encodes the byte array to "soft" symbols (each output byte as one bit, 0 for 0 and 255 for 1)
func (cc *ConvolutionCoder) EncodeSoft(data []byte) (output []byte) {
	encoded := cc.Encode(data)
	output = make([]byte, cc.EncodedSize(len(data)))
	bl := len(encoded)

	for i := 0; i < bl && i*8 < len(output); i++ {
		d := encoded[i]
		for z := 7; z >= 0; z-- {
			output[i*8+(7-z)] = 0
			if d&(1<<uint(z)) == 0 {
				output[i*8+(7-z)] = 255
			}
		}
	}

	return output
}

// Encode encodes the byte array (each output byte has 8 encoded bits)
func (cc *ConvolutionCoder) Encode(data []byte) (output []byte) {
	frameBits := cc.EncodedSize(len(data))
	bl := frameBits/8 + 1

	if frameBits%8 == 0 {
		bl -= 1
	}

	output = make([]byte, bl)

	correctwrap.Correct_convolutional_encode(cc.cc, &data[0], int64(len(data)), &output[0])

	return output
}

// Decode decodes a byte array containing 8 hard symbols per byte
func (cc *ConvolutionCoder) Decode(data []byte) (output []byte) {
	frameBits := int64(len(data)) * 8
	output = make([]byte, int(frameBits)/(8*cc.r))
	correctwrap.Correct_convolutional_decode(cc.cc, &data[0], frameBits, &output[0])

	return output
}

// DecodeSoft decodes a byte array containing one bit per byte as soft symbol (0 as 0, 1 as 255)
func (cc *ConvolutionCoder) DecodeSoft(data []byte) (output []byte) {
	frameBits := int64(len(data))
	output = make([]byte, int(frameBits)/(8*cc.r))
	correctwrap.Correct_convolutional_decode_soft(cc.cc, &data[0], frameBits, &output[0])

	return output
}

// Close cleans the Convolution Coder native resources
func (cc *ConvolutionCoder) Close() {
	correctwrap.Correct_convolutional_destroy(cc.cc)
}
