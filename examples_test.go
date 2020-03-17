package senml_test

import (
	"fmt"

	"github.com/farshidtz/senml/v2"
	"github.com/farshidtz/senml/v2/codec"
)

func Example_DecodeJSON() {
	input := `[{"bn":"room1/temp","u":"Cel","t":1276020076.305,"v":23.5},{"u":"Cel","t":1276020091.305,"v":23.6}]`

	// decode JSON
	pack, err := codec.DecodeJSON([]byte(input))
	if err != nil {
		panic(err) // handle the error
	}

	// validate the SenML Pack
	err = pack.Validate()
	if err != nil {
		panic(err) // handle the error
	}
}

func Example_EncodeJSON() {
	v := 23.1
	var p senml.Pack = []senml.Record{
		{Value: &v, Unit: "Cel", Name: "urn:dev:ow:10e2073a01080063"},
	}

	dataOut, err := codec.EncodeJSON(p)
	if err != nil {
		panic(err) // handle the error
	}
	fmt.Printf("%s", dataOut)
	// Output: [{"n":"urn:dev:ow:10e2073a01080063","u":"Cel","v":23.1}]
}

func ExamplePack_Normalize() {
	input := `[{"bn":"room1/temp","u":"Cel","t":1276020076.305,"v":23.5},{"u":"Cel","t":1276020091.305,"v":23.6}]`

	// decode JSON
	pack, err := codec.DecodeJSON([]byte(input))
	if err != nil {
		panic(err) // handle the error
	}

	pack.Normalize()

	dataOut, err := codec.EncodeJSON(pack, codec.PrettyPrint)
	if err != nil {
		panic(err) // handle the error
	}
	fmt.Printf("%s", dataOut)
	// Output:
	// [
	//   {"n":"room1/temp","u":"Cel","t":1276020076.305,"v":23.5},
	//   {"n":"room1/temp","u":"Cel","t":1276020091.305,"v":23.6}
	// ]
}

func ExamplePack_Validate() {
	input := `[{"bn":"room1/ temp","t":1270000040,"v":23.5},{"t":1270000050,"v":23.6}]`

	// decode JSON
	pack, err := codec.DecodeJSON([]byte(input))
	if err != nil {
		panic(err) // handle the error
	}

	// validate the SenML Pack
	err = pack.Validate()
	if err != nil {
		fmt.Println(err) // handle the error
	}
	// Output: invalid name: must begin with alphanumeric and contain alphanumeric or one of - : . / _
}

func ExamplePack_Validate2() {
	input := `[{"bn":"room1/temp","t":1270000050,"v":23.6,"vs":"cool"}]`

	// decode JSON
	pack, err := codec.DecodeJSON([]byte(input))
	if err != nil {
		panic(err) // handle the error
	}

	// validate the SenML Pack
	err = pack.Validate()
	if err != nil {
		fmt.Println(err) // handle the error
	}
	// Output: too many values in single record
}

func Example_EncodeXML() {
	var pack senml.Pack = []senml.Record{
		{Time: 1276020000, Name: "room1/temp_label", StringValue: "hot"},
		{Time: 1276020100, Name: "room1/temp_label", StringValue: "cool"},
	}

	// encode to Pretty XML
	xmlBytes, err := codec.EncodeXML(pack, codec.PrettyPrint)
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

func Example_EncodeCSV() {
	var pack senml.Pack = []senml.Record{
		{Time: 1276020000, Name: "room1/temp_label", StringValue: "hot"},
		{Time: 1276020100, Name: "room1/temp_label", StringValue: "cool"},
	}

	// encode to CSV (format: name,excel-time,value,unit)
	csvBytes, err := codec.EncodeCSV(pack, codec.WithHeader)
	if err != nil {
		panic(err) // handle the error
	}
	fmt.Printf("%s\n", csvBytes)
	// Output:
	// Time,Update Time,Name,Unit,Value,String Value,Boolean Value,Data Value,Sum
	// 1276020000,0,room1/temp_label,,,hot,,,
	// 1276020100,0,room1/temp_label,,,cool,,,
}
