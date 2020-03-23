// Package codec provides various encoding/decoding functions for SenML Packs:
// http://github.com/farshidtz/senml
package codec

import (
	"fmt"
	"io"

	"github.com/farshidtz/senml/v2"
)

// Encoder is the encoding function type
type Encoder func(p senml.Pack, options ...Option) ([]byte, error)

// Decoder is the decoding function type
type Decoder func(b []byte, options ...Option) (senml.Pack, error)

// Reader is the reading function type
type Reader func(r io.Reader, options ...Option) (senml.Pack, error)

// Write is the writing function type
type Writer func(p senml.Pack, w io.Writer, options ...Option) error

func Encode(mediaType string, p senml.Pack, options ...Option) ([]byte, error) {
	switch mediaType {
	case senml.MediaTypeSenmlJSON:
		return EncodeJSON(p, options...)
	case senml.MediaTypeSenmlXML:
		return EncodeXML(p, options...)
	case senml.MediaTypeSenmlCBOR:
		return EncodeCBOR(p)
	case senml.MediaTypeSenmlCSV:
		return EncodeCSV(p, options...)
	}
	return nil, fmt.Errorf("unsupported media type: %s", mediaType)
}

func Decode(mediaType string, b []byte, options ...Option) (senml.Pack, error) {
	switch mediaType {
	case senml.MediaTypeSenmlJSON:
		return DecodeJSON(b)
	case senml.MediaTypeSenmlXML:
		return DecodeXML(b)
	case senml.MediaTypeSenmlCBOR:
		return DecodeCBOR(b)
	case senml.MediaTypeSenmlCSV:
		return DecodeCSV(b, options...)
	}
	return nil, fmt.Errorf("unsupported media type: %s", mediaType)
}
