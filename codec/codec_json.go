package codec

import (
	"bytes"
	"encoding/json"

	"github.com/farshidtz/senml/v2"
)

type JSONCoder struct {
	senml.Options
}

// EncodeJSON serializes the SenML pack into JSON bytes
func (c JSONCoder) Encode(p senml.Pack, options ...senml.Option) ([]byte, error) {
	for _, opt := range options {
		opt(&c.Options)
	}

	if c.PrettyPrint {
		var buf bytes.Buffer
		buf.WriteString("[\n  ")
		for i, r := range p {
			if i != 0 {
				buf.WriteString(",\n  ")
			}
			recData, err := json.Marshal(r)
			if err != nil {
				return nil, err
			}
			buf.Write(recData)
		}
		buf.WriteString("\n]\n")
		return buf.Bytes(), nil
	}

	return json.Marshal(p)
}

// DecodeJSON takes a SenML pack in JSON bytes and decodes it into a Pack
func (JSONCoder) Decode(b []byte, options ...senml.Option) (senml.Pack, error) {
	var p senml.Pack

	err := json.Unmarshal(b, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
