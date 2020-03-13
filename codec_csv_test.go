package senml

import (
	"testing"
)

const (
	// CSV should be in normalized form
	csvString = `946684799,dev123temp,degC,22.1,,,,0,10
946684799,dev123room,degC,,kitchen,,,,0
946684800,dev123data,degC,,,,abc,,0
946684800,dev123ok,degC,,,true,,,0
`

	csvStringWithHeader = `Time,Name,Unit,Value,String Value,Boolean Value,Data Value,Sum,Update Time
` + csvString
)

func TestEncodeCSV(t *testing.T) {

	ref := referencePack(true)

	t.Run("without header", func(t *testing.T) {
		dataOut, err := ref.EncodeCSV(false)
		if err != nil {
			t.Fatalf("Encoding error: %s", err)
		}

		if string(dataOut) != csvString {
			t.Logf("Expected:\n'%s'", csvString)
			t.Fatalf("Got:\n'%s'", dataOut)
		}
	})

	t.Run("with header", func(t *testing.T) {
		dataOut, err := ref.EncodeCSV(true)
		if err != nil {
			t.Fatalf("Encoding error: %s", err)
		}

		if string(dataOut) != csvStringWithHeader {
			t.Logf("Expected:\n'%s'", csvStringWithHeader)
			t.Fatalf("Got:\n'%s'", dataOut)
		}
	})
}

func TestDecodeCSV(t *testing.T) {

	t.Run("compare fields", func(t *testing.T) {
		type pair struct {
			got      interface{}
			expected interface{}
		}
		ref := referencePack(true)
		ref.Normalize()

		pack, err := DecodeCSV([]byte(csvString), false)
		if err != nil {
			t.Fatalf("Error decoding: %s", err)
		}

		pairs := make(map[string]pair)
		for i := range pack {
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
			if pack[i].Sum != nil {
				pairs["Sum"] = pair{*pack[i].Sum, *ref[i].Sum}
			}
			// compare values
			for fieldName, p := range pairs {
				if p.got != p.expected {
					t.Logf("Assertion failed for %s", fieldName)
					t.Fatalf("Got: '%v' instead of: '%v'", p.got, p.expected)
				}
			}
		}
	})

	t.Run("wrong header", func(t *testing.T) {
		data := []byte("Bad,Time,Name,Unit,Value,String Value,Boolean Value,Data Value,Sum,Update Time")
		_, err := DecodeCSV(data, true)
		if err == nil {
			t.Fatalf("No error for wrong header")
		}
	})

}
