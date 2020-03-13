package senml

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

	ref := referencePack()

	t.Run("minified", func(t *testing.T) {
		dataOut, err := ref.EncodeXML(false)
		if err != nil {
			t.Fatalf("Encoding error: %s", err)
		}

		if string(dataOut) != xmlStringMinified {
			t.Logf("Expected:\n'%s'", xmlStringMinified)
			t.Fatalf("Got:\n'%s'", dataOut)
		}
	})

	t.Run("pretty", func(t *testing.T) {
		dataOut, err := ref.EncodeXML(true)
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
		type pair struct {
			got      interface{}
			expected interface{}
		}
		ref := referencePack()

		pack, err := DecodeXML([]byte(xmlStringMinified))
		if err != nil {
			t.Fatalf("Error decoding: %s", err)
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
					t.Logf("Assertion failed for %s:", fieldName)
					t.Fatalf("Got: '%v' instead of: '%v'", p.got, p.expected)
				}
			}
		}
	})

}
