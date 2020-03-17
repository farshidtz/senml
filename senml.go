// Package senml is an implementation of RFC8428 - Sensor Measurement Lists (SenML):
// https://tools.ietf.org/html/rfc8428
//
// The sub-package codec provides various encoding/decoding functions for senml.Pack:
// https://godoc.org/github.com/farshidtz/senml/codec
package senml

import (
	"fmt"
	"regexp"
	"time"
)

// DefaultBaseVersion is the default version of the SenML data model based on:
// https://tools.ietf.org/html/rfc8428
const DefaultBaseVersion = 10

// Pack is a SenML Pack which is one or more SenML Records in an array structure.
type Pack []Record

// Record is a SenML Record which is one measurement or configuration instance in time presented using the SenML data model.
type Record struct {
	XMLName *bool `json:"_,omitempty" xml:"senml"`

	// BaseName is a string that is prepended to the names found in the entries.
	BaseName string `json:"bn,omitempty"  xml:"bn,attr,omitempty" cbor:"-2,keyasint,omitempty"`
	// BaseTime is added to the time found in an entry.
	BaseTime float64 `json:"bt,omitempty"  xml:"bt,attr,omitempty" cbor:"-3,keyasint,omitempty"`
	// BaseUnit is assumed for all entries, unless otherwise indicated.
	BaseUnit string `json:"bu,omitempty"  xml:"bu,attr,omitempty" cbor:"-4,keyasint,omitempty"`
	// BaseVersion is a positive integer and defaults to 10 if not present.
	BaseVersion *int `json:"bver,omitempty"  xml:"bver,attr,omitempty" cbor:"-1,keyasint,omitempty"`
	// BaseValue is added to the value found in an entry, similar to BaseTime.
	BaseValue *float64 `json:"bv,omitempty"  xml:"bv,attr,omitempty" cbor:"-5,keyasint,omitempty"`
	// BaseSum is added to the sum found in an entry, similar to BaseTime.
	BaseSum *float64 `json:"bs,omitempty"  xml:"bs,attr,omitempty" cbor:"-6,keyasint,omitempty"`

	// Name of the sensor or parameter.
	Name string `json:"n,omitempty"  xml:"n,attr,omitempty" cbor:"0,keyasint,omitempty"`
	// Unit for a measurement value.
	Unit string `json:"u,omitempty"  xml:"u,attr,omitempty" cbor:"1,keyasint,omitempty"`
	// Time in seconds when the value was recorded.
	Time float64 `json:"t,omitempty"  xml:"t,attr,omitempty" cbor:"6,keyasint,omitempty"`
	// UpdateTime is the maximum seconds before there is an updated reading for a measurement.
	UpdateTime float64 `json:"ut,omitempty"  xml:"ut,attr,omitempty" cbor:"7,keyasint,omitempty"`

	// Value is the float value of the entry.
	Value *float64 `json:"v,omitempty"  xml:"v,attr,omitempty" cbor:"2,keyasint,omitempty"`
	// StringValue is the string value of the entry.
	StringValue string `json:"vs,omitempty"  xml:"vs,attr,omitempty" cbor:"3,keyasint,omitempty"`
	// DataValue is a base64-encoded string value of the entry with the URL-safe alphabet.
	DataValue string `json:"vd,omitempty"  xml:"vd,attr,omitempty" cbor:"8,keyasint,omitempty"`
	// BoolValue is the boolean value of the entry.
	BoolValue *bool `json:"vb,omitempty"  xml:"vb,attr,omitempty" cbor:"4,keyasint,omitempty"`
	// Sum is the integrated sum of the float values over time.
	Sum *float64 `json:"s,omitempty"  xml:"s,attr,omitempty" cbor:"5,keyasint,omitempty"`
}

// Normalize converts the SenML Pack to to the resolved format according to:
// https://tools.ietf.org/html/rfc8428#section-4.6
//
// Normalize must be called on a validated pack only.
func (p Pack) Normalize() {
	var bname string
	var btime float64
	var bunit string
	var bver int
	var bvalue float64
	var bsum float64

	var now = float64(time.Now().UnixNano()) / 1000000000
	const pivot = 268435456 // rfc8428: values less than 2**28 represent time relative to the current time.
	var r *Record
	for i := range p {
		r = &p[i]

		// Time
		if r.BaseTime != 0 {
			btime = r.BaseTime
			r.BaseTime = 0
		}
		r.Time = btime + r.Time
		if r.Time < pivot {
			// convert to absolute time
			r.Time = now + r.Time
		}

		// Version
		if r.BaseVersion == nil && bver != 0 {
			r.BaseVersion = &bver
		} else if r.BaseVersion != nil {
			if *r.BaseVersion == DefaultBaseVersion {
				r.BaseVersion = nil
			} else {
				bver = *r.BaseVersion
			}
		}

		// Value
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

		// Sum
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

		// Unit
		if len(r.BaseUnit) > 0 {
			bunit = r.BaseUnit
			r.BaseUnit = ""
		}
		if len(r.Unit) == 0 {
			r.Unit = bunit
		}

		// Name
		if len(r.BaseName) > 0 {
			bname = r.BaseName
			r.BaseName = ""
		}
		r.Name = bname + r.Name
	}

	return
}

// Clone returns a deep copy of the SenML Pack
func (p Pack) Clone() (clone Pack) {
	cloneBool := func(b *bool) *bool {
		if b != nil {
			clone := new(bool)
			*clone = *b
			return clone
		}
		return nil
	}
	cloneInt := func(i *int) *int {
		if i != nil {
			clone := new(int)
			*clone = *i
			return clone
		}
		return nil
	}
	cloneFloat64 := func(f *float64) *float64 {
		if f != nil {
			clone := new(float64)
			*clone = *f
			return clone
		}
		return nil
	}
	clone = make(Pack, len(p))
	for i := range p {
		clone[i] = Record{
			XMLName:     cloneBool(p[i].XMLName),
			BaseName:    p[i].BaseName,
			BaseTime:    p[i].BaseTime,
			BaseUnit:    p[i].BaseUnit,
			BaseVersion: cloneInt(p[i].BaseVersion),
			BaseValue:   cloneFloat64(p[i].BaseValue),
			BaseSum:     cloneFloat64(p[i].BaseSum),
			Name:        p[i].Name,
			Unit:        p[i].Unit,
			Time:        p[i].Time,
			UpdateTime:  p[i].UpdateTime,
			Value:       cloneFloat64(p[i].Value),
			StringValue: p[i].StringValue,
			DataValue:   p[i].DataValue,
			BoolValue:   cloneBool(p[i].BoolValue),
			Sum:         cloneFloat64(p[i].Sum),
		}
	}

	return clone
}

// Validate tests if the SenML Pack is valid
func (p Pack) Validate() error {
	var bname string
	var bver = -1

	for _, r := range p {
		// validate version
		if r.BaseVersion == nil {
			if bver == -1 {
				bver = DefaultBaseVersion
			}
		} else {
			if *r.BaseVersion < 0 {
				fmt.Errorf("negative base version")
			}
			if bver == -1 {
				bver = *r.BaseVersion
			} else {
				return fmt.Errorf("unallowed version change")
			}
		}

		// validate name
		if len(r.BaseName) > 0 {
			bname = r.BaseName
		}
		name := bname + r.Name
		err := ValidateName(name)
		if err != nil {
			return err
		}

		// validate values
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
		return fmt.Errorf("invalid regex for name validation: %s", err)
	}
	if !validName.MatchString(name) {
		return fmt.Errorf("invalid name: must begin with alphanumeric and contain alphanumeric or one of - : . / _")
	}
	return nil
}
