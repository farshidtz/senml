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

// Format is the SenML encoding/decoding format
type Format int

// Encoding/Decoding constants
const (
	JSON Format = 1 + iota
	XML
	CBOR
	CSV
	MPACK
	LINEP
	JSONLINE
)

const DEFAULT_BASE_VERSION = 10

// OutputOptions are encoding options
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

	BaseName string  `json:"bn,omitempty"  xml:"bn,attr,omitempty"`
	BaseTime float64 `json:"bt,omitempty"  xml:"bt,attr,omitempty"`
	BaseUnit string  `json:"bu,omitempty"  xml:"bu,attr,omitempty"`
	// BaseVersion is an optional positive integer and defaults to 10 if not present.
	BaseVersion *int `json:"bver,omitempty"  xml:"bver,attr,omitempty"`
	// BaseValue is added to the value found in an entry, similar to BaseTime.
	BaseValue *float64 `json:"bv,omitempty"  xml:"bv,attr,omitempty"`
	// BaseSum is added to the sum found in an entry, similar to BaseTime.
	BaseSum *float64 `json:"bs,omitempty"  xml:"bs,attr,omitempty"`

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

// Decode takes a SenML message in the given format, parses it and decodes it
//	into the returned SenML pack.
func Decode(msg []byte, format Format) (Pack, error) {
	var p Pack
	var err error

	switch {
	case format == JSON:
		// parse the input JSON object
		err = json.Unmarshal(msg, &p)
		if err != nil {
			return p, err
		}

	case format == JSONLINE:
		// parse the input JSON lines
		lines := strings.Split(string(msg), "\n")
		for _, line := range lines {
			r := new(Record)
			if len(line) > 5 {
				err = json.Unmarshal([]byte(line), r)
				if err != nil {
					return p, fmt.Errorf("error parsing JSON line: %s", err)
				}
				p = append(p, *r)
			}
		}

	case format == XML:
		// parse the input XML
		var temp xmlPack
		err = xml.Unmarshal(msg, &temp)
		if err != nil {
			return nil, err
		}
		p = temp.Pack

	case format == CBOR:
		// parse the input CBOR
		var cborHandle codec.Handle = new(codec.CborHandle)
		var decoder *codec.Decoder = codec.NewDecoderBytes(msg, cborHandle)
		err = decoder.Decode(&p)
		if err != nil {
			return p, fmt.Errorf("error parsing CBOR: %s", err)
		}

	case format == MPACK:
		// parse the input MPACK
		// spec for MessagePack is at https://github.com/msgpack/msgpack/
		var mpackHandle codec.Handle = new(codec.MsgpackHandle)
		var decoder *codec.Decoder = codec.NewDecoderBytes(msg, mpackHandle)
		err = decoder.Decode(&p)
		if err != nil {
			return p, fmt.Errorf("error parsing MPACK: %s", err)
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
		// output JSON version
		if options.PrettyPrint {
			var lines string
			lines += fmt.Sprintf("[\n  ")
			for i, r := range p {
				if i != 0 {
					lines += ",\n  "
				}
				recData, err := json.Marshal(r)
				if err != nil {
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
			return nil, err
		}

	case format == CSV:
		// normalize first to add base values to record values
		normalized := p.Normalize()
		// output a CSV version
		// format: name,excel-time,value(,unit)
		var lines string
		for _, r := range normalized {
			if r.Value != nil {
				// TODO - replace sprintf with bytes.Buffer
				lines += fmt.Sprintf("%s,%f,%f", r.Name, r.Time, *r.Value)
				if len(r.Unit) > 0 {
					lines += fmt.Sprintf(",%s", r.Unit)
				}
				lines += fmt.Sprintf("\r\n")
			}
		}
		data = []byte(lines)

	case format == CBOR:
		// output a CBOR version
		var cborHandle codec.Handle = new(codec.CborHandle)
		var encoder *codec.Encoder = codec.NewEncoderBytes(&data, cborHandle)
		err = encoder.Encode(p)
		if err != nil {
			return nil, fmt.Errorf("error encoding CBOR: %s", err)
		}

	case format == MPACK:
		// output a MPACK version
		var mpackHandle codec.Handle = new(codec.MsgpackHandle)
		var encoder *codec.Encoder = codec.NewEncoderBytes(&data, mpackHandle)
		err = encoder.Encode(p)
		if err != nil {
			return nil, fmt.Errorf("error encoding MPACK: %s", err)
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
					return nil, fmt.Errorf("error encoding JSON line: %s", err)
				}
				buf.Write(data)
				buf.WriteString("\n")
			}
		}
		data = buf.Bytes()
	}

	return data, nil
}

// Normalize removes all the base items adds them to corresponding record fields. It converts relative times to absolute times.
// Normalize must be called on a validated pack only. The Decode function performs this validation internally.
// A SenML Record is referred to as "resolved" if it does not contain
//   any base values, i.e., labels starting with the character "b", except
//   for Base Version fields (see below), and has no relative times.  To
//   resolve the Records, the applicable base values of the SenML Pack (if
//   any) are applied to the Record.  That is, for the base values in the
//   Record or before the Record in the Pack, Name and Base Name are
//   concatenated, the Base Time is added to the time of the Record, the
//   Base Unit is applied to the Record if it did not contain a Unit, etc.
//   In addition, the Records need to be in chronological order in the
//   Pack.
//  Source: https://tools.ietf.org/html/rfc8428#section-4.6
func (p Pack) Normalize() Pack {
	var bname string
	var btime float64
	var bunit string
	var bver int
	var bvalue float64
	var bsum float64

	ret := make(Pack, len(p))

	var now = float64(time.Now().UnixNano()) / 1000000000
	const pivot = 268435456 // rfc8428: values less than 2**28 represent time relative to the current time.
	for i, r := range p {
		if r.BaseTime != 0 {
			btime = r.BaseTime
		}
		// RFC8428: The Base Version field MUST NOT be present in resolved Records if the
		//   SenML version defined in this document is used; otherwise, it MUST be
		//   present in all the resolved SenML Records.
		if r.BaseVersion == nil && bver != 0 {
			r.BaseVersion = &bver
		} else if r.BaseVersion != nil {
			if *r.BaseVersion == DEFAULT_BASE_VERSION {
				r.BaseVersion = nil
			} else {
				bver = *r.BaseVersion
			}
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

		if r.Time < pivot {
			// convert to absolute time
			r.Time = now + r.Time
		}

		if r.BaseValue != nil {
			bvalue = *r.BaseValue
			r.BaseValue = nil
		}
		if bvalue != 0 {
			if r.Value == nil {
				r.Value = new(float64)
			}
			*r.Value += bvalue
		}

		if r.BaseSum != nil {
			bsum = *r.BaseSum
			r.BaseSum = nil
		}
		if bsum != 0 {
			if r.Sum == nil {
				r.Sum = new(float64)
			}
			*r.Sum += bsum
		}

		ret[i] = r
	}

	return ret
}

// Validate tests if SenML is valid
func (p Pack) Validate() error {
	var bname string
	var bver = DEFAULT_BASE_VERSION

	for _, r := range p {

		// Check version is same for all records
		if bver == DEFAULT_BASE_VERSION {
			// set the bver the first time it is seen
			if r.BaseVersion != nil && *r.BaseVersion != DEFAULT_BASE_VERSION {
				bver = *r.BaseVersion
			}
		} else {
			if r.BaseVersion != nil && *r.BaseVersion != DEFAULT_BASE_VERSION {
				// next time a version in seen, check it has not changed
				if *r.BaseVersion != bver {
					return fmt.Errorf("unallowed version change")
				}
			}
		}

		// Check name
		if len(r.BaseName) > 0 {
			bname = r.BaseName
		}
		name := bname + r.Name
		err := ValidateName(name)
		if err != nil {
			return err
		}

		floatValueCount := 0
		nonFloatValueCount := 0
		if r.Value != nil || r.BaseValue != nil {
			floatValueCount++
		}
		if r.BoolValue != nil {
			nonFloatValueCount++
		}
		if len(r.DataValue) > 0 {
			nonFloatValueCount++
		}
		if len(r.StringValue) > 0 {
			nonFloatValueCount++
		}

		if floatValueCount+nonFloatValueCount > 1 {
			return fmt.Errorf("too many values in single record")
		}

		if nonFloatValueCount == 1 {
			if r.Sum != nil || r.BaseSum != nil {
				return fmt.Errorf("sum together with non-float value in a single record")
			}
		} else {
			if floatValueCount == 0 && r.Sum == nil && r.BaseSum == nil {
				return fmt.Errorf("no value or sum")
			}
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
