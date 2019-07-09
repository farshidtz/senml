# SenML: Sensor Measurement Lists

[![GoDoc](https://godoc.org/github.com/farshidtz/senml?status.svg)](https://godoc.org/github.com/farshidtz/senml)
[![Build Status](https://travis-ci.org/farshidtz/senml.svg)](https://travis-ci.org/farshidtz/senml)
[![Go Report Card](https://goreportcard.com/badge/github.com/farshidtz/senml)](https://goreportcard.com/report/github.com/farshidtz/senml)

SenML package is an implementation of [RFC8428](https://tools.ietf.org/html/rfc8428) - Sensor Measurement Lists (SenML) in Go.

Note: This library is under development. Only use tagged commits for production.

## Install
```
go get github.com/farshidtz/senml
```

## Usage
```go
package main

import (
	"fmt"
	"github.com/farshidtz/senml"
)

func main() {
	jsonBytes := []byte(`[{"bn":"urn:dev:ow:10e2073a01080063","u":"Cel","t":1276020076.305,"v":23.5},{"u":"Cel","t":1276020091.305,"v":23.6}]`)

	pack, err := senml.Decode(jsonBytes, senml.JSON)
	if err != nil {
		panic(err) // don't panic, handle the error
	}

	senmlBytes, err := pack.Encode(senml.JSON, senml.OutputOptions{PrettyPrint: true})
	if err != nil {
		panic(err) // don't panic, handle the error
	}
	fmt.Printf("%s\n", senmlBytes)
	/* Output:
	[
		{"bn":"urn:dev:ow:10e2073a01080063","u":"Cel","t":1276020076.305,"v":23.5},
		{"u":"Cel","t":1276020091.305,"v":23.6}
	]
	*/

	xmlBytes, err := pack.Encode(senml.XML, senml.OutputOptions{})
	if err != nil {
		panic(err) // don't panic, handle the error
	}
	fmt.Printf("%s\n", xmlBytes)
	/* Output:
	<sensml xmlns="urn:ietf:params:xml:ns:senml"><senml bn="urn:dev:ow:10e2073a01080063" u="Cel" t="1.276020076305e+09" v="23.5"></senml><senml u="Cel" t="1.276020091305e+09" v="23.6"></senml></sensml>
	*/
}

```
