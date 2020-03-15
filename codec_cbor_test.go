package senml

import (
	"encoding/base64"
	"reflect"
	"testing"
)

const (
	cborBase64EncodedString = `hKpiYm5mZGV2MTIzYmJ0+8BG1cKPXCj2YmJ1ZGRlZ0NkYnZlcgVhbmR0ZW1wYXP7AAAAAAAAAABhdPu/8AAAAAAAAGF1ZGRlZ0NidXT7QCQAAAAAAABhdvtANhmZmZmZmqNhbmRyb29tYXT7v/AAAAAAAABidnNna2l0Y2hlbqJhbmRkYXRhYnZkY2FiY6JhbmJva2J2YvU=`
)

func TestEncodeCBOR(t *testing.T) {
	cborBytes, err := base64.StdEncoding.DecodeString(cborBase64EncodedString)
	if err != nil {
		t.Fatalf("Error decoding test value: %s", err)
	}

	ref := referencePack()

	dataOut, err := ref.EncodeCBOR()
	if err != nil {
		t.Fatalf("Encoding error: %s", err)
	}

	if !reflect.DeepEqual(dataOut, cborBytes) {
		decoded, err := base64.StdEncoding.DecodeString(cborBase64EncodedString)
		if err != nil {
			t.Fatalf("Error decoding test value: %s", err)
		}
		t.Logf("Expected:\n%v", decoded)
		t.Fatalf("Got:\n%v", dataOut)
	}
}

func TestDecodeCBOR(t *testing.T) {

	t.Run("compare fields", func(t *testing.T) {
		type pair struct {
			got      interface{}
			expected interface{}
		}

		cborBytes, err := base64.StdEncoding.DecodeString(cborBase64EncodedString)
		if err != nil {
			t.Fatalf("Error decoding test value: %s", err)
		}

		pack, err := DecodeCBOR(cborBytes)
		if err != nil {
			t.Fatalf("Error decoding: %s", err)
		}

		ref := referencePack()
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
