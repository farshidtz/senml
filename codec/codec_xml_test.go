package codec

import (
	"fmt"
	"testing"

	"github.com/farshidtz/senml/v2"
)

const (
	xmlStringMinified = `<sensml xmlns="urn:ietf:params:xml:ns:senml"><senml bn="dev123" bt="-45.67" bu="degC" bver="5" n="temp" u="degC" t="-1" ut="10" v="22.1" s="0"></senml><senml n="room" t="-1" vs="kitchen"></senml><senml n="data" vd="abc"></senml><senml n="ok" vb="true"></senml></sensml>`

	xmlStringPretty = `<sensml xmlns="urn:ietf:params:xml:ns:senml">
  <senml bn="dev123" bt="-45.67" bu="degC" bver="5" n="temp" u="degC" t="-1" ut="10" v="22.1" s="0"></senml>
  <senml n="room" t="-1" vs="kitchen"></senml>
  <senml n="data" vd="abc"></senml>
  <senml n="ok" vb="true"></senml>
</sensml>`
)

func TestEncodeXML(t *testing.T) {

	t.Run("minified", func(t *testing.T) {
		dataOut, err := EncodeXML(referencePack())
		if err != nil {
			t.Fatalf("Encoding error: %s", err)
		}

		if string(dataOut) != xmlStringMinified {
			t.Logf("Expected:\n'%s'", xmlStringMinified)
			t.Fatalf("Got:\n'%s'", dataOut)
		}
	})

	t.Run("pretty", func(t *testing.T) {
		dataOut, err := EncodeXML(referencePack(), PrettyPrint)
		if err != nil {
			t.Fatalf("Encoding error: %s", err)
		}

		if string(dataOut) != xmlStringPretty {
			t.Logf("Expected:\n'%s'", xmlStringPretty)
			t.Fatalf("Got:\n'%s'", dataOut)
		}
	})
}

func TestDecodeXML(t *testing.T) {

	t.Run("compare fields", func(t *testing.T) {

		pack, err := DecodeXML([]byte(xmlStringMinified))
		if err != nil {
			t.Fatalf("Error decoding: %s", err)
		}

		if err := compareFields(pack, referencePack()); err != nil {
			t.Fatalf("Error matching records: %s", err)
		}
	})

}

// EXAMPLES

func ExampleEncodeXML() {
	var pack senml.Pack = []senml.Record{
		{Time: 1276020000, Name: "room1/temp_label", StringValue: "hot"},
		{Time: 1276020100, Name: "room1/temp_label", StringValue: "cool"},
	}

	// encode to Pretty XML
	xmlBytes, err := EncodeXML(pack, PrettyPrint)
	if err != nil {
		panic(err) // handle the error
	}
	fmt.Printf("%s\n", xmlBytes)
	// Output:
	// <sensml xmlns="urn:ietf:params:xml:ns:senml">
	//   <senml n="room1/temp_label" t="1.27602e+09" vs="hot"></senml>
	//   <senml n="room1/temp_label" t="1.2760201e+09" vs="cool"></senml>
	// </sensml>
}

func ExampleDecodeXML() {
	input := `<sensml xmlns="urn:ietf:params:xml:ns:senml"><senml bn="dev123" bt="-45.67" bu="degC" bver="5" n="temp" u="degC" t="-1" ut="10" v="22.1" s="0"></senml><senml n="room" t="-1" vs="kitchen"></senml><senml n="data" vd="abc"></senml><senml n="ok" vb="true"></senml></sensml>`

	// decode XML
	pack, err := DecodeXML([]byte(input))
	if err != nil {
		panic(err) // handle the error
	}

	// validate the SenML Pack
	err = pack.Validate()
	if err != nil {
		panic(err) // handle the error
	}
}
