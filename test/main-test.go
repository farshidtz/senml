package main

import (
	"fmt"
	"github.com/farshidtz/senml"
)

func main() {
	jsonBytes := []byte(`[{"bn":"room1/temp","u":"Cel","t":1276020076.305,"v":23.5},{"u":"Cel","t":1276020091.305,"v":23.6}]`)

	// decode JSON
	pack, err := senml.Decode(jsonBytes, senml.JSON)
	if err != nil {
		panic(err) // don't panic, handle the error
	}

	// encode to pretty JSON
	senmlBytes, err := pack.Encode(senml.JSON, senml.OutputOptions{PrettyPrint: true})
	if err != nil {
		panic(err) // don't panic, handle the error
	}
	fmt.Printf("%s\n", senmlBytes)
	/* Output:
	[
		{"bn":"room1/temp","u":"Cel","t":1276020076.305,"v":23.5},
		{"u":"Cel","t":1276020091.305,"v":23.6}
	]
	*/

	// encode to XML
	xmlBytes, err := pack.Encode(senml.XML, senml.OutputOptions{})
	if err != nil {
		panic(err) // don't panic, handle the error
	}
	fmt.Printf("%s\n", xmlBytes)
	/* Output:
	<sensml xmlns="urn:ietf:params:xml:ns:senml"><senml bn="room1/temp" u="Cel" t="1.276020076305e+09" v="23.5"></senml><senml u="Cel" t="1.276020091305e+09" v="23.6"></senml></sensml>
	*/

	// encode to CSV (format: name,excel-time,value,unit)
	csvBytes, err := pack.Encode(senml.CSV, senml.OutputOptions{})
	if err != nil {
		panic(err) // don't panic, handle the error
	}
	fmt.Printf("%s\n", csvBytes)
	/* Output:
	room1/temp,40337.750883,23.500000,Cel
	room1/temp,40337.751057,23.600000,Cel
	*/
}
