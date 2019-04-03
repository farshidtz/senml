package senml

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ugorji/go/codec"
)

type Format int

const (
	JSON Format = 1 + iota
	XML
	CBOR
	CSV
	MPACK
	LINEP
	JSONLINE
)

type OutputOptions struct {
	PrettyPrint bool
	Topic       string
}

// Pack is a SenML Pack:
//	One or more SenML Records in an array structure.
type Pack []Record

// Record is a SenML Record:
//	One measurement or configuration instance in time presented using the SenML data model.
type Record struct {
	XMLName *bool `json:"_,omitempty" xml:"senml"`

	BaseName    string  `json:"bn,omitempty"  xml:"bn,attr,omitempty"`
	BaseTime    float64 `json:"bt,omitempty"  xml:"bt,attr,omitempty"`
	BaseUnit    string  `json:"bu,omitempty"  xml:"bu,attr,omitempty"`
	BaseVersion int     `json:"bver,omitempty"  xml:"bver,attr,omitempty"`

	Name       string  `json:"n,omitempty"  xml:"n,attr,omitempty"`
	Unit       string  `json:"u,omitempty"  xml:"u,attr,omitempty"`
	Time       float64 `json:"t,omitempty"  xml:"t,attr,omitempty"`
	UpdateTime float64 `json:"ut,omitempty"  xml:"ut,attr,omitempty"`

	Value       *float64 `json:"v,omitempty"  xml:"v,attr,omitempty"`
	StringValue string   `json:"vs,omitempty"  xml:"vs,attr,omitempty"`
	DataValue   string   `json:"vd,omitempty"  xml:"vd,attr,omitempty"`
	BoolValue   *bool    `json:"vb,omitempty"  xml:"vb,attr,omitempty"`

	Sum *float64 `json:"s,omitempty"  xml:"s,attr,omitempty"`
}

type xmlPack struct {
	Pack
	XMLName *bool  `xml:"sensml"`
	XMLNS   string `xml:"xmlns,attr"`
}

// Decode takes a SenML message in the given format and parses it and decodes it
//	into the returned SenML record.
func Decode(msg []byte, format Format) (Pack, error) {
	var p Pack
	var err error

	switch {
	case format == JSON:
		// parse the input JSON stream
		err = json.Unmarshal(msg, &p)
		if err != nil {
			//fmt.Println("error parsing JSON SenML Stream: ", err)
			//fmt.Println("msg=", msg)
			return p, err
		}

	case format == JSONLINE:
		// parse the input JSON line
		lines := strings.Split(string(msg), "\n")
		for _, line := range lines {
			r := new(Record)
			if len(line) > 5 {
				err = json.Unmarshal([]byte(line), r)
				if err != nil {
					//fmt.Println("error parsing JSON SenML Line: ", err)
					return p, err
				}
				p = append(p, *r)
			}
		}

	case format == XML:
		// parse the input XML
		var temp xmlPack
		err = xml.Unmarshal(msg, &temp)
		if err != nil {
			//fmt.Println("error parsing XML SenML", err)
			return nil, err
		}
		p = temp.Pack

	case format == CBOR:
		// parse the input CBOR
		var cborHandle codec.Handle = new(codec.CborHandle)
		var decoder *codec.Decoder = codec.NewDecoderBytes(msg, cborHandle)
		err = decoder.Decode(&p)
		if err != nil {
			//fmt.Println("error parsing CBOR SenML", err)
			return p, err
		}

	case format == MPACK:
		// parse the input MPACK
		// spec for MessagePack is at https://github.com/msgpack/msgpack/
		var mpackHandle codec.Handle = new(codec.MsgpackHandle)
		var decoder *codec.Decoder = codec.NewDecoderBytes(msg, mpackHandle)
		err = decoder.Decode(&p)
		if err != nil {
			//fmt.Println("error parsing MPACK SenML", err)
			return p, err
		}

	}

	if err := p.Validate(); err != nil {
		return p, fmt.Errorf("invalid SenML Pack: %s", err)
	}

	return p, nil
}

// Encode takes a SenML record, and encodes it using the given format.
func (p Pack) Encode(format Format, options OutputOptions) ([]byte, error) {
	var data []byte
	var err error

	if options.Topic == "" {
		options.Topic = "senml"
	}

	switch {

	case format == JSON:
		// ouput JSON version
		if options.PrettyPrint {
			// data, err = json.MarshalIndent(s.Records, "", "  ")
			var lines string
			lines += fmt.Sprintf("[\n  ")
			for i, r := range p {
				if i != 0 {
					lines += ",\n  "
				}
				recData, err := json.Marshal(r)
				if err != nil {
					//fmt.Println("error encoding JSON SenML", err)
					return nil, err
				}
				lines += fmt.Sprintf("%s", recData)
			}
			lines += fmt.Sprintf("\n]\n")
			data = []byte(lines)
		} else {
			return json.Marshal(p)
		}

	case format == XML:
		xmlPack := xmlPack{
			Pack:  p,
			XMLNS: "urn:ietf:params:xml:ns:senml",
		}
		// output a XML version
		if options.PrettyPrint {
			data, err = xml.MarshalIndent(&xmlPack, "", "  ")
		} else {
			data, err = xml.Marshal(&xmlPack)
		}
		if err != nil {
			//fmt.Println("error encoding XML SenML", err)
			return nil, err
		}

	case format == CSV:
		// output a CSV version
		var lines string
		for _, r := range p {
			if r.Value != nil {
				// TODO - replace sprintf with bytes.Buffer
				lines += fmt.Sprintf("%s,", r.Name)
				// excell time in days since 1900, unix seconds since 1970
				// ( 1970 is 25569 days after 1900 )
				lines += fmt.Sprintf("%f,", (r.Time/(24.0*3600.0))+25569.0)
				lines += fmt.Sprintf("%f", *r.Value)
				if len(r.Unit) > 0 {
					lines += fmt.Sprintf(",%s", r.Unit)
				}
				lines += fmt.Sprintf("\r\n")
			}
		}
		data = []byte(lines)
		if err != nil {
			//fmt.Println("error encoding CSV SenML", err)
			return nil, err
		}

	case format == CBOR:
		// output a CBOR version
		var cborHandle codec.Handle = new(codec.CborHandle)
		var encoder *codec.Encoder = codec.NewEncoderBytes(&data, cborHandle)
		err = encoder.Encode(p)
		if err != nil {
			//fmt.Println("error encoding CBOR SenML", err)
			return nil, err
		}

	case format == MPACK:
		// output a MPACK version
		var mpackHandle codec.Handle = new(codec.MsgpackHandle)
		var encoder *codec.Encoder = codec.NewEncoderBytes(&data, mpackHandle)
		err = encoder.Encode(p)
		if err != nil {
			//fmt.Println("error encoding MPACK SenML", err)
			return nil, err
		}

	case format == LINEP:
		// output Line Protocol
		var buf bytes.Buffer
		for _, r := range p {
			if r.Value != nil {
				buf.WriteString(options.Topic)
				buf.WriteString(",n=")
				buf.WriteString(r.Name)
				buf.WriteString(",u=")
				buf.WriteString(r.Unit)
				buf.WriteString(" v=")
				buf.WriteString(strconv.FormatFloat(*r.Value, 'f', -1, 64))
				buf.WriteString(" ")
				buf.WriteString(strconv.FormatInt(int64(r.Time*1.0e9), 10))
				buf.WriteString("\n")
			}
		}
		data = buf.Bytes()

	case format == JSONLINE:
		// output Line Protocol
		var buf bytes.Buffer
		for _, r := range p {
			if r.Value != nil {
				data, err = json.Marshal(r)
				if err != nil {
					//fmt.Println("error encoding JSONLINE SenML", err)
					return nil, err
				}
				buf.Write(data)
				buf.WriteString("\n")
			}
		}
		data = buf.Bytes()
	}

	return data, nil
}

// Normalize removes all the base items and expands records to have items that include
//	what previously in base items. Converts relative times to absolute times.
func (p Pack) Normalize() Pack {
	var bname string = ""
	var btime float64 = 0
	var bunit string = ""
	var ver = 5
	var ret Pack

	var totalRecords int = 0
	for _, r := range p {
		if (r.Value != nil) || (len(r.StringValue) > 0) || (len(r.DataValue) > 0) || (r.BoolValue != nil) {
			totalRecords += 1
		}
	}

	ret = make([]Record, totalRecords)
	var numRecords = 0

	for _, r := range p {
		if r.BaseTime != 0 {
			btime = r.BaseTime
		}
		if r.BaseVersion != 0 {
			ver = r.BaseVersion
		}
		if len(r.BaseUnit) > 0 {
			bunit = r.BaseUnit
		}
		if len(r.BaseName) > 0 {
			bname = r.BaseName
		}
		r.BaseTime = 0
		r.BaseUnit = ""
		r.BaseName = ""
		r.Name = bname + r.Name
		r.Time = btime + r.Time
		if len(r.Unit) == 0 {
			r.Unit = bunit
		}
		r.BaseVersion = ver

		if r.Time <= 0 {
			// convert to absolute time
			var now float64 = float64(time.Now().UnixNano()) / 1000000000
			r.Time = now + r.Time
		}

		if (r.Value != nil) || (len(r.StringValue) > 0) || (len(r.DataValue) > 0) || (r.BoolValue != nil) {
			ret[numRecords] = r
			numRecords += 1
		}
	}

	return ret
}

// Validate test if SenML is valid
func (p Pack) Validate() error {
	var bname string = ""
	var bver = -1

	for _, r := range p {

		// Check version is same for all records
		if bver == -1 {
			// set the bver the first time it is seen
			if r.BaseVersion != 0 {
				bver = r.BaseVersion
			}
		} else {
			if r.BaseVersion != 0 {
				// next time a version in seen, check it has not changed
				if r.BaseVersion != bver {
					return fmt.Errorf("unallowed version change")
				}
			}
		}

		// Check name
		if len(r.BaseName) > 0 {
			bname = r.BaseName
		}
		name := bname + r.Name
		err := ValidateSenmlName(name)

		if err != nil {
			return err
		}

		valueCount := 0
		if r.Value != nil {
			valueCount = valueCount + 1
		}
		if r.BoolValue != nil {
			valueCount = valueCount + 1
		}
		if len(r.DataValue) > 0 {
			valueCount = valueCount + 1
		}
		if len(r.StringValue) > 0 {
			valueCount = valueCount + 1
		}
		if valueCount > 1 {
			return fmt.Errorf("too many values")
		}
		if r.Sum != nil {
			valueCount = valueCount + 1
		}
		if valueCount < 1 {
			return fmt.Errorf("no value or sum")
		}

		// Check if name is known Mandatory To Understand
		//for k :=  r {
		// 	fmt.Println( "key=" , k  )
		//         if k[ len(k)-1 ] == '_' {
		//         	fmt.Println("unknown MTU in record")
		//		return false
		//        }
		// }
	}

	return nil
}

// ValidateName validates the SenML name
func ValidateName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("empty name")
	}
	validName, err := regexp.Compile(`^[a-zA-Z0-9]+[a-zA-Z0-9-:./_]*$`)
	if err != nil {
		fmt.Println(err)
	}
	if !validName.MatchString(name) {
		return fmt.Errorf("invalid name: must begin with alphanumeric and contain alphanumeric or one of - : . / _")
	}
	return nil
}
