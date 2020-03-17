# SenML: Sensor Measurement Lists

[![GoDoc](https://godoc.org/github.com/farshidtz/senml?status.svg)](https://godoc.org/github.com/farshidtz/senml)
[![Build Status](https://travis-ci.org/farshidtz/senml.svg)](https://travis-ci.org/farshidtz/senml)
[![Go Report Card](https://goreportcard.com/badge/github.com/farshidtz/senml)](https://goreportcard.com/report/github.com/farshidtz/senml)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/farshidtz/senml?sort=semver&label=stable)
![GitHub tag (latest SemVer pre-release)](https://img.shields.io/github/v/tag/farshidtz/senml?include_prereleases&sort=semver&label=pre)


SenML package is an implementation of [RFC8428](https://tools.ietf.org/html/rfc8428) - Sensor Measurement Lists (SenML) in Go.

It provides fully compliant data model and functionalities for:

* Validation of various SenML fields
* [Normalization](https://tools.ietf.org/html/rfc8428#section-4.6)
* [SenML Units](https://tools.ietf.org/html/rfc8428#section-12.1)
* [SenML Media Types](https://tools.ietf.org/html/rfc8428#section-12.3)
* Encoding/Decoding (codec package)
    * [JSON](https://tools.ietf.org/html/rfc8428#section-5)
    * [XML](https://tools.ietf.org/html/rfc8428#section-7)
    * [CBOR](https://tools.ietf.org/html/rfc8428#section-6)
    * CSV (custom)
      

## Install
```
go get github.com/farshidtz/senml/v2
```

## Usage
```go
package main

import (
	"fmt"
	"github.com/farshidtz/senml/v2"
	"github.com/farshidtz/senml/v2/codec"
)

func main() {
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
```
