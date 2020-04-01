package codec

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/farshidtz/senml/v2"
)

// Decode:
// protoc --decode=senml_protobuf.Message senml.proto < senml_protobuf_test_output.bin
// Decode raw:
// protoc --decode_raw < senml_protobuf_test_output.bin

const (
	protobufHexBytesString = "0a490a0664657631323311f6285c8fc2d546c01a046465674320053a0474656d7042046465674349000000000000f0bf510000000000002440599a999999991936407900000000000000000a183a04726f6f6d49000000000000f0bf62076b69746368656e0a0b3a04646174616a036162630a063a026f6b7001"
)

func TestEncodeProtobuf(t *testing.T) {

	dataOut, err := EncodeProtobuf(referencePack())
	if err != nil {
		t.Fatalf("Encoding error: %s", err)
	}
	dataOutHex := hex.EncodeToString(dataOut)

	if dataOutHex != protobufHexBytesString {
		err = ioutil.WriteFile("./senml_protobuf_test_output.bin", dataOut, 0644)
		if err != nil {
			t.Fatalf("Error writing encoded message to file: %s", err)
		}

		t.Logf("Expected (hex):\n%v", protobufHexBytesString)
		t.Fatalf("Got (hex):\n%v", dataOutHex)
	}
}

func TestDecodeProtobuf(t *testing.T) {

	t.Run("compare fields", func(t *testing.T) {

		cborBytes, err := hex.DecodeString(protobufHexBytesString)
		if err != nil {
			t.Fatalf("Error decoding test value: %s", err)
		}

		pack, err := DecodeProtobuf(cborBytes)
		if err != nil {
			t.Fatalf("Error decoding: %s", err)
		}

		if err := compareFields(pack, referencePack()); err != nil {
			t.Fatalf("Error matching records: %s", err)
		}
	})

}

// EXAMPLES

func ExampleEncodeProtobuf() {
	v := 23.1
	var p senml.Pack = []senml.Record{
		{Value: &v, Unit: "Cel", Name: "urn:dev:ow:10e2073a01080063"},
	}

	dataOut, err := EncodeProtobuf(p)
	if err != nil {
		panic(err) // handle the error
	}
	fmt.Printf("%v", dataOut)
	// Output: [10 43 58 27 117 114 110 58 100 101 118 58 111 119 58 49 48 101 50 48 55 51 97 48 49 48 56 48 48 54 51 66 3 67 101 108 89 154 153 153 153 153 25 55 64]
}

// Diagnose the hex output:
// https://protogen.marcgravell.com/decode
func ExampleEncodeProtobuf_hex() {
	v := 23.1
	var p senml.Pack = []senml.Record{
		{Value: &v, Unit: "Cel", Name: "urn:dev:ow:10e2073a01080063"},
	}

	dataOut, err := EncodeProtobuf(p)
	if err != nil {
		panic(err) // handle the error
	}
	fmt.Printf(hex.EncodeToString(dataOut))
	// Output: 0a2b3a1b75726e3a6465763a6f773a31306532303733613031303830303633420343656c599a99999999193740
}
