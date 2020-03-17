package senml_test

import (
	"fmt"

	"github.com/farshidtz/senml/v2/codec"
)

func Example_decodeValidateNormalizeEncode() {
	input := `[{"bn":"room1/temp","u":"Cel","t":1276020076,"v":23.5},{"u":"Cel","t":1276020091,"v":23.6}]`

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

	// normalize the SenML Pack
	pack.Normalize()

	// encode the normalized SenML Pack to JSON
	dataOut, err := codec.EncodeJSON(pack, codec.PrettyPrint)
	if err != nil {
		panic(err) // handle the error
	}
	fmt.Printf("%s", dataOut)
	// Output:
	// [
	//   {"n":"room1/temp","u":"Cel","t":1276020076,"v":23.5},
	//   {"n":"room1/temp","u":"Cel","t":1276020091,"v":23.6}
	// ]
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
