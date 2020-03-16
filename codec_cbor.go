package senml

import (
	"github.com/fxamacker/cbor/v2"
)

// EncodeCBOR serializes the SenML pack into CBOR bytes
func (p Pack) EncodeCBOR() ([]byte, error) {

	return cbor.Marshal(p)
}

// DecodeCBOR takes a SenML pack in CBOR bytes and decodes it into a Pack
func DecodeCBOR(b []byte) (Pack, error) {
	var p Pack

	err := cbor.Unmarshal(b, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
