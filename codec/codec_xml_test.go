package codec

import (
	"testing"
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
		dataOut, err := EncodeXML(referencePack(), false)
		if err != nil {
			t.Fatalf("Encoding error: %s", err)
		}

		if string(dataOut) != xmlStringMinified {
			t.Logf("Expected:\n'%s'", xmlStringMinified)
			t.Fatalf("Got:\n'%s'", dataOut)
		}
	})

	t.Run("pretty", func(t *testing.T) {
		dataOut, err := EncodeXML(referencePack(), true)
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
