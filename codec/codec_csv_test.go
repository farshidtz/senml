package codec

import (
	"testing"
)

const (
	// CSV should be in normalized form
	csvString = `946684799,10,dev123temp,degC,22.1,,,,0,
946684799,0,dev123room,degC,,kitchen,,,,
946684800,0,dev123data,degC,,,,abc,,
946684800,0,dev123ok,degC,,,true,,,
`

	csvStringWithHeader = `Time,Update Time,Name,Unit,Value,String Value,Boolean Value,Data Value,Sum
` + csvString
)

func TestEncodeCSV(t *testing.T) {

	t.Run("without header", func(t *testing.T) {
		dataOut, err := EncodeCSV(referencePack(true))
		if err != nil {
			t.Fatalf("Encoding error: %s", err)
		}

		if string(dataOut) != csvString {
			t.Logf("Expected:\n'%s'", csvString)
			t.Fatalf("Got:\n'%s'", dataOut)
		}
	})

	t.Run("with header", func(t *testing.T) {
		dataOut, err := EncodeCSV(referencePack(true), WithHeader)
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

		pack, err := DecodeCSV([]byte(csvString))
		if err != nil {
			t.Fatalf("Error decoding: %s", err)
		}

		ref := referencePack(true)
		ref.Normalize()

		if err := compareFields(pack, ref); err != nil {
			t.Fatalf("Error matching records: %s", err)
		}
	})

	t.Run("wrong header", func(t *testing.T) {
		data := []byte("Bad,Time,Name,Unit,Value,String Value,Boolean Value,Data Value,Sum,Update Time")
		_, err := DecodeCSV(data, WithHeader)
		if err == nil {
			t.Fatalf("No error for wrong header")
		}
	})

}
