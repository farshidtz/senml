package senml

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"testing"
	"time"
)

func ExampleEncode1() {
	v := 23.1
	var p Pack = []Record{
		{Value: &v, Unit: "Cel", Name: "urn:dev:ow:10e2073a01080063"},
	}

	dataOut, err := p.Encode(JSON, nil)
	if err != nil {
		fmt.Println("Encode of SenML failed")
	} else {
		fmt.Println(string(dataOut))
	}
	// Output: [{"n":"urn:dev:ow:10e2073a01080063","u":"Cel","v":23.1}]
}

func ExampleEncode2() {
	v1 := 23.5
	v2 := 23.6
	var p Pack = []Record{
		{Value: &v1, Unit: "Cel", BaseName: "urn:dev:ow:10e2073a01080063", Time: 1.276020076305e+09},
		{Value: &v2, Unit: "Cel", Time: 1.276020091305e+09},
	}

	dataOut, err := p.Encode(JSON, nil)
	if err != nil {
		fmt.Println("Encode of SenML failed")
	} else {
		fmt.Println(string(dataOut))
	}
	// Output: [{"bn":"urn:dev:ow:10e2073a01080063","u":"Cel","t":1276020076.305,"v":23.5},{"u":"Cel","t":1276020091.305,"v":23.6}]
}

type TestVector struct {
	testDecode bool
	format     Format
	label      string
	binary     bool
	value      string
}

var testVectors = []TestVector{
	{true, JSON, "JSON", false, "W3siYm4iOiJkZXYxMjMiLCJidCI6LTQ1LjY3LCJidSI6ImRlZ0MiLCJidmVyIjo1LCJuIjoidGVtcCIsInUiOiJkZWdDIiwidCI6LTEsInV0IjoxMCwidiI6MjIuMSwicyI6MH0seyJuIjoicm9vbSIsInQiOi0xLCJ2cyI6ImtpdGNoZW4ifSx7Im4iOiJkYXRhIiwidmQiOiJhYmMifSx7Im4iOiJvayIsInZiIjp0cnVlfV0="},
	{true, CBOR, "CBOR", true, "hKpiYm5mZGV2MTIzYmJ0+8BG1cKPXCj2YmJ1ZGRlZ0NkYnZlcgVhbmR0ZW1wYXP7AAAAAAAAAABhdPu/8AAAAAAAAGF1ZGRlZ0NidXT7QCQAAAAAAABhdvtANhmZmZmZmqNhbmRyb29tYXT7v/AAAAAAAABidnNna2l0Y2hlbqJhbmRkYXRhYnZkY2FiY6JhbmJva2J2YvU="},
	{true, XML, "XML", false, "PHNlbnNtbCB4bWxucz0idXJuOmlldGY6cGFyYW1zOnhtbDpuczpzZW5tbCI+PHNlbm1sIGJuPSJkZXYxMjMiIGJ0PSItNDUuNjciIGJ1PSJkZWdDIiBidmVyPSI1IiBuPSJ0ZW1wIiB1PSJkZWdDIiB0PSItMSIgdXQ9IjEwIiB2PSIyMi4xIiBzPSIwIj48L3Nlbm1sPjxzZW5tbCBuPSJyb29tIiB0PSItMSIgdnM9ImtpdGNoZW4iPjwvc2VubWw+PHNlbm1sIG49ImRhdGEiIHZkPSJhYmMiPjwvc2VubWw+PHNlbm1sIG49Im9rIiB2Yj0idHJ1ZSI+PC9zZW5tbD48L3NlbnNtbD4="},
	{false, CSV, "CSV", false, "ZGV2MTIzdGVtcCw5NDY2ODQ3OTkuMDAwMDAwLDIyLjEwMDAwMCxkZWdDDQo="},
	{true, MPACK, "MPACK", true, "lIqiYm6mZGV2MTIzomJ0y8BG1cKPXCj2omJ1pGRlZ0OkYnZlcgWhbqR0ZW1woXPLAAAAAAAAAAChdMu/8AAAAAAAAKF1pGRlZ0OidXTLQCQAAAAAAAChdstANhmZmZmZmoOhbqRyb29toXTLv/AAAAAAAACidnOna2l0Y2hlboKhbqRkYXRhonZko2FiY4KhbqJva6J2YsM="},
	{false, LINEP, "LINEP", false, "Zmx1ZmZ5U2VubWwsbj10ZW1wLHU9ZGVnQyB2PTIyLjEgLTEwMDAwMDAwMDAK"},
}

func referencePack() Pack {
	bver := 5
	value := 22.1
	sum := 0.0
	vb := true
	return Pack{
		{BaseName: "dev123",
			BaseTime:    -45.67,
			BaseUnit:    "degC",
			BaseVersion: &bver,
			Value:       &value, Unit: "degC", Name: "temp", Time: -1.0, UpdateTime: 10.0, Sum: &sum},
		{StringValue: "kitchen", Name: "room", Time: -1.0},
		{DataValue: "abc", Name: "data"},
		{BoolValue: &vb, Name: "ok"},
	}
}

func referencePackFloats() Pack {
	bver := 5
	value, value2 := 22.1, 30.0
	sum, sum2 := 100.0, 200.0
	return Pack{
		{BaseName: "dev123",
			BaseTime:    -45.67,
			BaseUnit:    "degC",
			BaseVersion: &bver,
			Value:       &value, Unit: "degC", Name: "temp", Time: -1.0, UpdateTime: 10.0, Sum: &sum},
		{Value: &value2, Time: 1.0, Sum: &sum2},
	}
}

func stringifyPack(p Pack) string {
	b, _ := json.Marshal(&p)
	return string(b)
}

func TestEncode(t *testing.T) {

	options := OutputOptions{Topic: "fluffySenml", PrettyPrint: false}
	for _, vector := range testVectors {
		ref := referencePack()
		t.Run(vector.label, func(t *testing.T) {
			if vector.label == "CSV" {
				// change to an absolute time: https://tools.ietf.org/html/rfc8428#section-4.5.3
				ref[0].BaseTime = 946684800
			}
			dataOut, err := ref.Encode(vector.format, &options)
			if err != nil {
				t.Fatalf("Encoding error: %s", err)
			}

			if base64.StdEncoding.EncodeToString(dataOut) != vector.value {
				t.Errorf("Assertion failed for encoded %s:", vector.label)
				t.Logf("Got (encoded): %s", base64.StdEncoding.EncodeToString(dataOut))
				if !vector.binary {
					t.Logf("Got:\n'%s'", dataOut)
					decoded, err := base64.StdEncoding.DecodeString(vector.value)
					if err != nil {
						t.Fatalf("Error decoding test value: %s", err)
					}
					t.Fatalf("Expected:\n'%s'", decoded)
				} else {
					t.Logf("Got:\n%v", dataOut)
					decoded, err := base64.StdEncoding.DecodeString(vector.value)
					if err != nil {
						t.Fatalf("Error decoding test value: %s", err)
					}
					t.Fatalf("Expected:\n%v", decoded)
				}
			}
		})
	}

}

func TestDecode(t *testing.T) {
	// test different serializations
	type pair struct {
		got      interface{}
		expected interface{}
	}
	ref := referencePack()
	for _, vector := range testVectors {
		if vector.testDecode {
			t.Run(vector.label, func(t *testing.T) {
				data, err := base64.StdEncoding.DecodeString(vector.value)
				if err != nil {
					t.Fatalf("Error decoding test value for %s: %s", vector.label, err)
				}

				pack, err := Decode(data, vector.format)
				if err != nil {
					t.Fatalf("Error decoding %s: %s", vector.label, err)
				}

				pairs := make(map[string]pair)
				for i := range pack {
					pairs["XMLName"] = pair{pack[i].XMLName, ref[i].XMLName}
					pairs["BaseName"] = pair{pack[i].BaseName, ref[i].BaseName}
					pairs["BaseTime"] = pair{pack[i].BaseTime, ref[i].BaseTime}
					pairs["BaseUnit"] = pair{pack[i].BaseUnit, ref[i].BaseUnit}
					pairs["Name"] = pair{pack[i].Name, ref[i].Name}
					pairs["Unit"] = pair{pack[i].Unit, ref[i].Unit}
					pairs["Time"] = pair{pack[i].Time, ref[i].Time}
					pairs["UpdateTime"] = pair{pack[i].UpdateTime, ref[i].UpdateTime}
					pairs["StringValue"] = pair{pack[i].StringValue, ref[i].StringValue}
					pairs["DataValue"] = pair{pack[i].DataValue, ref[i].DataValue}
					// pointers
					if pack[i].BaseVersion != nil {
						pairs["Value"] = pair{*pack[i].BaseVersion, *ref[i].BaseVersion}
					}
					if pack[i].Value != nil {
						pairs["Value"] = pair{*pack[i].Value, *ref[i].Value}
					}
					if pack[i].BoolValue != nil {
						pairs["BoolValue"] = pair{*pack[i].BoolValue, *ref[i].BoolValue}
					}
					if pack[i].Sum != nil {
						pairs["Sum"] = pair{*pack[i].Sum, *ref[i].Sum}
					}
					// compare values
					for fieldName, p := range pairs {
						if p.got != p.expected {
							t.Logf("Assertion failed for %s in encoded %s:", fieldName, vector.label)
							t.Fatalf("Got: '%v' instead of: '%v'", p.got, p.expected)
						}
					}
				}
			})
		}
	}

	// test various JSON inputs
	t.Run("JSON invalid object", func(t *testing.T) {
		data := []byte(" foo ")
		_, err := Decode(data, JSON)
		if err == nil {
			t.Fatalf("No error for invalid object")
		}
	})

	t.Run("JSON no pack", func(t *testing.T) {
		data := []byte(`{"n":"hi"}`)
		_, err := Decode(data, JSON)
		if err == nil {
			t.Fatalf("No error for record out of pack")
		}
	})

}

func TestNormalize(t *testing.T) {

	t.Run("Positive relative time", func(t *testing.T) {
		p := referencePack()
		p[0].BaseTime = 1000
		p.Normalize()
		now := float64(time.Now().UnixNano()) / 1000000000
		expected := now + 1000
		if math.Abs(p[0].Time-expected) > 5 { // fail if difference is more than 5s
			t.Fatalf("Time is not absolute. Got %f instead of %f", p[0].Time, expected)
		}
	})

	t.Run("Negative relative time", func(t *testing.T) {
		// negative relative time
		p := referencePack()
		p[0].BaseTime = -1000
		p.Normalize()
		now := float64(time.Now().UnixNano()) / 1000000000
		expected := now - 1000
		if math.Abs(p[0].Time-expected) > 5 { // fail if difference is more than 5s
			t.Fatalf("Time is not absolute. Got %f instead of %f", p[0].Time, expected)
		}
	})

	t.Run("Absolute relative time", func(t *testing.T) {
		p := referencePack()
		p[0].BaseTime = 946684800.123
		p.Normalize()
		dataOut, err := p.Encode(JSON, &OutputOptions{PrettyPrint: true})
		if err != nil {
			t.Fatalf("Error encoding: %s", err)
		}
		testValue := "WwogIHsiYnZlciI6NSwibiI6ImRldjEyM3RlbXAiLCJ1IjoiZGVnQyIsInQiOjk0NjY4NDc5OS4xMjMsInV0IjoxMCwidiI6MjIuMSwicyI6MH0sCiAgeyJidmVyIjo1LCJuIjoiZGV2MTIzcm9vbSIsInUiOiJkZWdDIiwidCI6OTQ2Njg0Nzk5LjEyMywidnMiOiJraXRjaGVuIn0sCiAgeyJidmVyIjo1LCJuIjoiZGV2MTIzZGF0YSIsInUiOiJkZWdDIiwidCI6OTQ2Njg0ODAwLjEyMywidmQiOiJhYmMifSwKICB7ImJ2ZXIiOjUsIm4iOiJkZXYxMjNvayIsInUiOiJkZWdDIiwidCI6OTQ2Njg0ODAwLjEyMywidmIiOnRydWV9Cl0K"
		if base64.StdEncoding.EncodeToString(dataOut) != testValue {
			t.Logf("Got (encoded): %s", base64.StdEncoding.EncodeToString(dataOut))
			t.Errorf("Assertion failed for normalized pack. Got:\n'%s'", dataOut)
			decoded, err := base64.StdEncoding.DecodeString(testValue)
			if err != nil {
				t.Fatalf("Error decoding test value: %s", err)
			}
			t.Errorf("Expected:\n'%s'", decoded)
		}
	})

	t.Run("Default base version in first record", func(t *testing.T) {
		p := referencePack()
		*p[0].BaseVersion = DefaultBaseVersion
		p.Normalize()
		for i, r := range p {
			if r.BaseVersion != nil {
				t.Errorf("Default base version was not omitted in record %d: %s", i, stringifyPack(p))
			}
		}
		if t.Failed() {
			t.FailNow()
		}
	})

	t.Run("Default base version in second record", func(t *testing.T) {
		p := referencePack()
		p[0].BaseVersion = nil
		p[1].BaseVersion = new(int)
		*p[1].BaseVersion = DefaultBaseVersion
		p.Normalize()
		for i, r := range p {
			if r.BaseVersion != nil {
				t.Errorf("Default base version was not omitted in record %d: %s", i, stringifyPack(p))
			}
		}
		if t.Failed() {
			t.FailNow()
		}
	})

	t.Run("No base version", func(t *testing.T) {
		p := referencePack()
		p[0].BaseVersion = nil
		p.Normalize()
		for i, r := range p {
			if r.BaseVersion != nil {
				t.Errorf("Default base version was not omitted in record %d: %s", i, stringifyPack(p))
			}
		}
		if t.Failed() {
			t.FailNow()
		}
	})

	t.Run("Floats pack with base value", func(t *testing.T) {
		p := referencePackFloats()
		p[0].Value = nil
		p[1].Value = nil
		p[0].BaseValue = new(float64)
		*p[0].BaseValue = 10
		p.Normalize()
		if *p[0].Value != 10 && *p[1].Value != 10 {
			t.Fatalf("Base value was not added to value in records. Got values: %f, %f", *p[0].Value, *p[1].Value)
		}
		if p[0].BaseValue != nil {
			t.Fatalf("Base value was not removed from record in pack: %s", stringifyPack(p))
		}
	})

	t.Run("Floats pack with base value and values", func(t *testing.T) {
		p := referencePackFloats()
		p[0].BaseValue = new(float64)
		*p[0].BaseValue = 10
		p.Normalize()
		if *p[0].Value != 32.1 {
			t.Fatalf("Base value was not added to value in first record. Got value: %f", *p[0].Value)
		}
		if p[0].BaseValue != nil {
			t.Fatalf("Base value was not removed from record in pack: %s", stringifyPack(p))
		}
		if *p[1].Value != 40 {
			t.Fatalf("Base value was not added to value in second record. Got value: %f", *p[1].Value)
		}
	})

	t.Run("Floats pack with base sum", func(t *testing.T) {
		p := referencePackFloats()
		p[0].Sum = nil
		p[1].Sum = nil
		p[0].BaseSum = new(float64)
		*p[0].BaseSum = 10
		p.Normalize()
		if *p[0].Sum != 10 && *p[1].Sum != 10 {
			t.Fatalf("Base sum was not added to sum in records. Got sums: %f, %f", *p[0].Sum, *p[1].Sum)
		}
		if p[0].BaseSum != nil {
			t.Fatalf("Base sum was not removed from record in pack: %s", stringifyPack(p))
		}
	})

	t.Run("Floats pack with base sum and sums", func(t *testing.T) {
		p := referencePackFloats()
		p[0].BaseSum = new(float64)
		*p[0].BaseSum = 10
		p.Normalize()
		if *p[0].Sum != 110 {
			t.Fatalf("Base sum was not added to sum in first record. Got sum: %f", *p[0].Sum)
		}
		if p[0].BaseSum != nil {
			t.Fatalf("Base sum was not removed from record in pack: %s", stringifyPack(p))
		}
		if *p[1].Sum != 210 {
			t.Fatalf("Base sum was not added to sum in second record. Got sum: %f", *p[1].Sum)
		}
	})

}

func TestClone(t *testing.T) {
	p := referencePack()
	p[0].XMLName = new(bool)

	c := p.Clone()

	*p[0].XMLName = true
	*p[0].Value = 123.456
	*p[0].BaseVersion = 123
	p[0].Time = 123
	p[0].StringValue = "changed"

	if *p[0].XMLName == *c[0].XMLName ||
		*p[0].Value == *c[0].Value ||
		*p[0].BaseVersion == *c[0].BaseVersion ||
		p[0].Time == c[0].Time ||
		p[0].StringValue == c[0].StringValue {
		t.Fatalf("Clone is changed after changing the original pack.")
	}
}

func TestValidate(t *testing.T) {

	t.Run("no value", func(t *testing.T) {
		pack := Pack{
			{Name: "dev"},
		}
		err := pack.Validate()
		if err == nil {
			t.Fatalf("No error for record with no value: %s", stringifyPack(pack))
		}
	})

	t.Run("numeric name", func(t *testing.T) {
		value := 1.0
		pack := Pack{
			{Name: "3a", Value: &value},
		}
		err := pack.Validate()
		if err != nil {
			t.Fatalf("Error decoding record with numeric name: %s", err)
		}
	})

	t.Run("bad numeric name", func(t *testing.T) {
		value := 1.0
		pack := Pack{
			{Name: "-3a", Value: &value},
		}
		err := pack.Validate()
		if err == nil {
			t.Fatalf("No error for bad numeric name in: %s", stringifyPack(pack))
		}
	})

	t.Run("weird name", func(t *testing.T) {
		value := 1.0
		pack := Pack{
			{Name: "Az3-:./_", Value: &value},
		}
		err := pack.Validate()
		if err != nil {
			t.Fatalf("Error decoding record with valid name: %s", err)
		}
	})

	t.Run("bad weird name", func(t *testing.T) {
		value := 1.0
		pack := Pack{
			{Name: "A;b", Value: &value},
		}
		err := pack.Validate()
		if err == nil {
			t.Fatalf("No error for invalid name in: %s", stringifyPack(pack))
		}
	})

	t.Run("weird base name", func(t *testing.T) {
		value := 1.0
		pack := Pack{
			{BaseName: "Az3-:./_", Name: "/b", Value: &value},
		}
		err := pack.Validate()
		if err != nil {
			t.Fatalf("Error decoding record with valid base name: %s", err)
		}
	})

	t.Run("bad numeric base name", func(t *testing.T) {
		value := 1.0
		pack := Pack{
			{BaseName: "/room", Name: "/dev", Value: &value},
		}
		err := pack.Validate()
		if err == nil {
			t.Fatalf("No error for invalid numeric base name in: %s", stringifyPack(pack))
		}
		//
		pack[0].BaseName = "room#3"
		err = pack.Validate()
		if err == nil {
			t.Fatalf("No error for invalid numeric base name in: %s", stringifyPack(pack))
		}
	})

	t.Run("bad unknown MTU field", func(t *testing.T) {
		t.Skip("TODO")
		//data := []byte(`[{"n":"hi","v":1.0,"mtu_":1.0}]`)
		//_, err := Decode(data, JSON)
		//if err == nil {
		//	t.Fatalf("No error for bad unknown MTU field in: %s", data)
		//}
	})

	t.Run("multiple values in record", func(t *testing.T) {
		value := 1.0
		pack := Pack{
			{Name: "dev", Value: &value, StringValue: "on"},
		}
		err := pack.Validate()
		if err == nil {
			t.Fatalf("No error for multi-valued record in pack: %s", stringifyPack(pack))
		}
	})

	t.Run("base value with other types in pack", func(t *testing.T) {
		bval := 1.0
		pack := Pack{
			{Name: "dev", BaseValue: &bval, StringValue: "on"},
		}
		err := pack.Validate()
		if err == nil {
			t.Fatalf("No error for record with base value (float) and another non-float value in pack: %s", stringifyPack(pack))
		}
	})

	t.Run("multiple base versions in pack", func(t *testing.T) {
		value := 1.0
		bver := 5
		pack := Pack{
			{Name: "dev", Value: &value, BaseVersion: nil},
			{Name: "dev", Value: &value, BaseVersion: &bver},
		}
		err := pack.Validate()
		if err == nil {
			t.Fatalf("No error for pack with no version followed by custom version: %s", stringifyPack(pack))
		}
		//
		bver_default := DefaultBaseVersion
		pack = Pack{
			{Name: "dev", Value: &value, BaseVersion: nil},
			{Name: "dev", Value: &value, BaseVersion: &bver_default},
		}
		err = pack.Validate()
		if err == nil {
			t.Fatalf("No error for pack with no version followed by default version: %s", stringifyPack(pack))
		}
		pack = Pack{
			{Name: "dev", Value: &value, BaseVersion: &bver_default},
			{Name: "dev", Value: &value, BaseVersion: &bver},
		}
		err = pack.Validate()
		if err == nil {
			t.Fatalf("No error for pack with default followed by custom version: %s", stringifyPack(pack))
		}
		//
		pack = Pack{
			{Name: "dev", Value: &value, BaseVersion: &bver},
			{Name: "dev", Value: &value, BaseVersion: &bver_default},
		}
		err = pack.Validate()
		if err == nil {
			t.Fatalf("No error for pack with custom followed by default version: %s", stringifyPack(pack))
		}
	})

	t.Run("custom base version", func(t *testing.T) {
		value := 1.0
		bver := 5
		pack := Pack{
			{Name: "dev", Value: &value, BaseVersion: &bver},
			{Name: "dev", Value: &value, BaseVersion: nil},
		}
		err := pack.Validate()
		if err != nil {
			t.Fatalf("Error for pack with custom followed by no version: %s", err)
		}
	})

	t.Run("base value only", func(t *testing.T) {
		bval := 1.0
		pack := Pack{
			{Name: "dev", BaseValue: &bval},
		}
		err := pack.Validate()
		if err != nil {
			t.Fatalf("Error for pack with base value: %s", err)
		}
	})

	t.Run("sum only", func(t *testing.T) {
		sum := 1.0
		pack := Pack{
			{Name: "dev", Sum: &sum},
		}
		err := pack.Validate()
		if err != nil {
			t.Fatalf("Error for pack with base sum: %s", err)
		}
	})

	t.Run("base sum only", func(t *testing.T) {
		bsum := 1.0
		pack := Pack{
			{Name: "dev", BaseSum: &bsum},
		}
		err := pack.Validate()
		if err != nil {
			t.Fatalf("Error for pack with base sum: %s", err)
		}
		//
		sum := 10.0
		pack = Pack{
			{Name: "dev", BaseSum: &bsum},
			{Name: "dev", Sum: &sum},
		}
		err = pack.Validate()
		if err != nil {
			t.Fatalf("Error for pack with base sum and sum: %s", err)
		}
	})
}

func TestValidateName(t *testing.T) {
	t.Run("valid names", func(t *testing.T) {
		names := []string{"Aa-:./_", "urn:dev:ow:10e2073a", "http://example.com"}
		for _, name := range names {
			err := ValidateName(name)
			if err != nil {
				t.Fatalf("Error for valid name: %s: %s", name, err)
			}
		}
	})

	t.Run("invalid names", func(t *testing.T) {
		names := []string{"-A", ":A", ".A", "/A", "_A",
			"~A", "!A", "@A", "#A", "$A", "%A", "^A", "&A", "*A", "(A", "+A", "=A", " A", " ", "A "}
		for _, name := range names {
			err := ValidateName(name)
			if err == nil {
				t.Fatalf("No error for invalid name: %s", name)
			}
		}
	})
}
