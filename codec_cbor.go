package senml

import (
	"github.com/ugorji/go/codec"
)

// EncodeCBOR serializes the SenML pack into CBOR bytes
func (p Pack) EncodeCBOR() ([]byte, error) {
	var b []byte

	// TODO change to 1 liner?
	var cborHandle codec.Handle = new(codec.CborHandle)
	var encoder *codec.Encoder = codec.NewEncoderBytes(&b, cborHandle)
	err := encoder.Encode(p)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// DecodeCBOR takes a SenML pack in CBOR bytes and decodes it into a Pack
func DecodeCBOR(b []byte) (Pack, error) {
	var p Pack

	// TODO change to 1 liner?
	var cborHandle codec.Handle = new(codec.CborHandle)
	var decoder *codec.Decoder = codec.NewDecoderBytes(b, cborHandle)
	err := decoder.Decode(&p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
