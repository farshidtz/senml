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
		dataOut, err := new(JSONCoder).Encode(referencePack())
		if err != nil {
			t.Fatalf("Encoding error: %s", err)
		}

		if string(dataOut) != jsonStringMinified {
			t.Logf("Expected:\n'%s'", jsonStringMinified)
			t.Fatalf("Got:\n'%s'", dataOut)
		}
	})

	t.Run("pretty", func(t *testing.T) {
		dataOut, err := new(JSONCoder).Encode(referencePack(), senml.PrettyPrint(true))
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

		pack, err := new(JSONCoder).Decode([]byte(jsonStringMinified))
		if err != nil {
			t.Fatalf("Error decoding: %s", err)
		}

		if err := compareFields(pack, referencePack()); err != nil {
			t.Fatalf("Error matching records: %s", err)
		}

	})

	t.Run("invalid object", func(t *testing.T) {
		data := []byte(" foo ")
		_, err := new(JSONCoder).Decode(data)
		if err == nil {
			t.Fatalf("No error for invalid object")
		}
	})

	t.Run("no pack", func(t *testing.T) {
		data := []byte(`{"n":"hi"}`)
		_, err := new(JSONCoder).Decode(data)
		if err == nil {
			t.Fatalf("No error for record out of pack")
		}
	})

	t.Run("empty pack", func(t *testing.T) {
		data := []byte(`[]`)
		_, err := new(JSONCoder).Decode(data)
		if err != nil {
			t.Fatalf("Error for valid, empty pack")
		}
	})
}
