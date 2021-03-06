package codec

import (
	"fmt"
	"testing"

	"github.com/farshidtz/senml/v2"
)

const (
	jsonStringMinified = `[{"bn":"dev123","bt":-45.67,"bu":"degC","bver":5,"n":"temp","u":"degC","t":-1,"ut":10,"v":22.1,"s":0},{"n":"room","t":-1,"vs":"kitchen"},{"n":"data","vd":"abc"},{"n":"ok","vb":true}]`

	jsonStringPretty = `[
  {"bn":"dev123","bt":-45.67,"bu":"degC","bver":5,"n":"temp","u":"degC","t":-1,"ut":10,"v":22.1,"s":0},
  {"n":"room","t":-1,"vs":"kitchen"},
  {"n":"data","vd":"abc"},
  {"n":"ok","vb":true}
]
`
)

func TestEncodeJSON(t *testing.T) {

	t.Run("minified", func(t *testing.T) {
		dataOut, err := EncodeJSON(referencePack())
		if err != nil {
			t.Fatalf("Encoding error: %s", err)
		}

		if string(dataOut) != jsonStringMinified {
			t.Logf("Expected:\n'%s'", jsonStringMinified)
			t.Fatalf("Got:\n'%s'", dataOut)
		}
	})

	t.Run("pretty", func(t *testing.T) {
		dataOut, err := EncodeJSON(referencePack(), SetPrettyPrint)
		if err != nil {
			t.Fatalf("Encoding error: %s", err)
		}

		if string(dataOut) != jsonStringPretty {
			t.Logf("Expected:\n'%s'", jsonStringPretty)
			t.Fatalf("Got:\n'%s'", dataOut)
		}
	})
}

func TestDecodeJSON(t *testing.T) {

	t.Run("compare fields", func(t *testing.T) {

		pack, err := DecodeJSON([]byte(jsonStringMinified))
		if err != nil {
			t.Fatalf("Error decoding: %s", err)
		}

		if err := compareFields(pack, referencePack()); err != nil {
			t.Fatalf("Error matching records: %s", err)
		}

	})

	t.Run("invalid object", func(t *testing.T) {
		data := []byte(" foo ")
		_, err := DecodeJSON(data)
		if err == nil {
			t.Fatalf("No error for invalid object")
		}
	})

	t.Run("no pack", func(t *testing.T) {
		data := []byte(`{"n":"hi"}`)
		_, err := DecodeJSON(data)
		if err == nil {
			t.Fatalf("No error for record out of pack")
		}
	})

	t.Run("empty pack", func(t *testing.T) {
		data := []byte(`[]`)
		_, err := DecodeJSON(data)
		if err != nil {
			t.Fatalf("Error for valid, empty pack")
		}
	})
}

// EXAMPLES

func ExampleEncodeJSON() {
	value := 22.1
	var pack senml.Pack = []senml.Record{
		{Time: 1276020000, Name: "air_quality", StringValue: "good", BaseName: "room1/"},
		{Time: 1276020100, Name: "air_quality", StringValue: "excellent"},
		{Time: 1276020100, Name: "temp", Value: &value, Unit: senml.UnitCelsius},
	}

	pack.Normalize() // optional

	dataOut, err := EncodeJSON(pack, SetPrettyPrint)
	if err != nil {
		panic(err) // handle the error
	}
	fmt.Printf("%s", dataOut)
	// Output:
	// [
	//   {"n":"room1/air_quality","t":1276020000,"vs":"good"},
	//   {"n":"room1/air_quality","t":1276020100,"vs":"excellent"},
	//   {"n":"room1/temp","u":"Cel","t":1276020100,"v":22.1}
	// ]
}

func ExampleDecodeJSON() {
	input := `[{"bn":"room1/temp","u":"Cel","t":1276020076.305,"v":23.5},{"u":"Cel","t":1276020091.305,"v":23.6}]`

	// decode JSON
	pack, err := DecodeJSON([]byte(input))
	if err != nil {
		panic(err) // handle the error
	}

	// validate the SenML Pack
	err = pack.Validate()
	if err != nil {
		panic(err) // handle the error
	}
}
