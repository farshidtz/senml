package codec

import (
	"github.com/farshidtz/senml/v2"
	"github.com/fxamacker/cbor/v2"
)

// EncodeCBOR serializes the SenML pack into CBOR bytes
var EncodeCBOR Encoder = func(p senml.Pack, options ...Option) ([]byte, error) {

	return cbor.Marshal(p)
}

// DecodeCBOR takes a SenML pack in CBOR bytes and decodes it into a Pack
var DecodeCBOR Decoder = func(b []byte, options ...Option) (senml.Pack, error) {
	var p senml.Pack

	err := cbor.Unmarshal(b, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
