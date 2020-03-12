package senml

import (
	"bytes"
	"encoding/json"
)

// EncodeJSON serializes the SenML pack into JSON bytes
func (p Pack) EncodeJSON(pretty bool) ([]byte, error) {

	if pretty {
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
func DecodeJSON(b []byte) (Pack, error) {
	var p Pack

	err := json.Unmarshal(b, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
