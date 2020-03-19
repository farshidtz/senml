package codec

import (
	"fmt"

	"github.com/farshidtz/senml/v2"
)

func ExampleSetPrettyPrint() {
	var p senml.Pack = []senml.Record{
		{Time: 946684700, Name: "lamp/brightness", StringValue: "100", Unit: senml.UnitLumen},
		{Time: 946684800, Name: "lamp/brightness", StringValue: "500", Unit: senml.UnitLumen},
	}

	dataOut, err := EncodeJSON(p, SetPrettyPrint)
	if err != nil {
		panic(err) // handle the error
	}
	fmt.Printf("%s", dataOut)
	// Output:
	// [
	//   {"n":"lamp/brightness","u":"lm","t":946684700,"vs":"100"},
	//   {"n":"lamp/brightness","u":"lm","t":946684800,"vs":"500"}
	// ]
}

func ExampleSetDefaultHeader() {
	var p senml.Pack = []senml.Record{
		{Time: 946684700, Name: "lamp/brightness", StringValue: "100", Unit: senml.UnitLumen},
		{Time: 946684800, Name: "lamp/brightness", StringValue: "500", Unit: senml.UnitLumen},
	}

	dataOut, err := EncodeCSV(p, SetDefaultHeader)
	if err != nil {
		panic(err) // handle the error
	}
	fmt.Printf("%s", dataOut)
	// Output:
	// Time,Update Time,Name,Unit,Value,String Value,Boolean Value,Data Value,Sum
	// 946684700,0,lamp/brightness,lm,,100,,,
	// 946684800,0,lamp/brightness,lm,,500,,,
}
