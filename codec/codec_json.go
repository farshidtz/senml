package codec

import (
	"bytes"
	"encoding/json"

	"github.com/farshidtz/senml/v2"
)

// EncodeJSON serializes the SenML pack into JSON bytes
func EncodeJSON(p senml.Pack, options ...Option) ([]byte, error) {
	o := &Options{
		PrettyPrint: false,
	}
	for _, opt := range options {
		opt(o)
	}

	if o.PrettyPrint {
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
func DecodeJSON(b []byte) (senml.Pack, error) {
	var p senml.Pack

	err := json.Unmarshal(b, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
