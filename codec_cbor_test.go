package senml

import (
	"encoding/hex"
	"fmt"
	"testing"
)

const (
	// Visit http://cbor.me/ for converting from/to diagnostic notation
	// Diagnostic format for the following CBOR hex string:
	//`[
	//  {-2:"dev123",-3:-45.67,-4:"degC",-1:5,0:"temp",1:"degC",6:-1,7:10,2:22.1,5:0},
	//  {0:"room",6:-1,3:"kitchen"},
	//  {0:"data",8:"abc"},
	//  {0:"ok",4:true}
	//]`
	cborHexBytesString = "84aa216664657631323322fbc046d5c28f5c28f62364646567432005006474656d7001646465674306fbbff000000000000007fb402400000000000002fb403619999999999a05fb0000000000000000a30064726f6f6d06fbbff000000000000003676b69746368656ea20064646174610863616263a200626f6b04f5"
)

func TestGenerateHexReference(t *testing.T) {
	ref := referencePack()
	dataOut, err := ref.EncodeCBOR()
	if err != nil {
		t.Fatalf("Encoding error: %s", err)
	}
	fmt.Printf("Hex string for CBOR reference: %s\n", hex.EncodeToString(dataOut))
}

func TestEncodeCBOR(t *testing.T) {
	ref := referencePack()
	dataOut, err := ref.EncodeCBOR()
	if err != nil {
		t.Fatalf("Encoding error: %s", err)
	}
	dataOutHex := hex.EncodeToString(dataOut)

	if dataOutHex != cborHexBytesString {
		t.Logf("Expected (hex):\n%v", cborHexBytesString)
		t.Fatalf("Got (hex):\n%v", dataOutHex)
	}
}

func TestDecodeCBOR(t *testing.T) {

	t.Run("compare fields", func(t *testing.T) {
		type pair struct {
			got      interface{}
			expected interface{}
		}

		cborBytes, err := hex.DecodeString(cborHexBytesString)
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
