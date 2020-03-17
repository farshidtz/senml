package codec

import (
	"fmt"

	"github.com/farshidtz/senml/v2"
)

func referencePack(absoluteTime ...bool) senml.Pack {
	var btime float64
	if len(absoluteTime) == 1 && absoluteTime[0] {
		btime = 946684800
	} else {
		btime = -45.67
	}
	bver := 5
	value := 22.1
	sum := 0.0
	vb := true
	return senml.Pack{
		{BaseName: "dev123",
			BaseTime:    btime,
			BaseUnit:    "degC",
			BaseVersion: &bver,
			Value:       &value, Unit: "degC", Name: "temp", Time: -1.0, UpdateTime: 10.0, Sum: &sum},
		{StringValue: "kitchen", Name: "room", Time: -1.0},
		{DataValue: "abc", Name: "data"},
		{BoolValue: &vb, Name: "ok"},
	}
}

func compareFields(pack senml.Pack, ref senml.Pack) error {
	if len(pack) != len(ref) {
		return fmt.Errorf("number of records differ between the given pack (%d) and reference (%d)", len(pack), len(ref))
	}
	type pair struct {
		got      interface{}
		expected interface{}
	}
	pairs := make(map[string]pair)
	for i := range pack {
		pairs["XMLName"] = pair{pack[i].XMLName, ref[i].XMLName}
		// Base Name
		pairs["BaseName"] = pair{pack[i].BaseName, ref[i].BaseName}
		// Base Time
		pairs["BaseTime"] = pair{pack[i].BaseTime, ref[i].BaseTime}
		// Base Unit
		pairs["BaseUnit"] = pair{pack[i].BaseUnit, ref[i].BaseUnit}
		// Base Value
		if ref[i].BaseValue != nil {
			pairs["BaseValue"] = pair{pack[i].BaseValue, ref[i].BaseValue}
		}
		// Base Sum
		if ref[i].BaseSum != nil {
			pairs["BaseSum"] = pair{*pack[i].BaseSum, *ref[i].BaseSum}
		}
		// Base Version
		if pack[i].BaseVersion != nil {
			pairs["BaseVersion"] = pair{*pack[i].BaseVersion, *ref[i].BaseVersion}
		}
		// Name
		pairs["Name"] = pair{pack[i].Name, ref[i].Name}
		// Unit
		pairs["Unit"] = pair{pack[i].Unit, ref[i].Unit}
		// Value
		if ref[i].Value != nil {
			pairs["Value"] = pair{*pack[i].Value, *ref[i].Value}
		}
		// String Value
		pairs["StringValue"] = pair{pack[i].StringValue, ref[i].StringValue}
		// Boolean Value
		if ref[i].BoolValue != nil {
			pairs["BoolValue"] = pair{*pack[i].BoolValue, *ref[i].BoolValue}
		}
		// Data Value
		pairs["DataValue"] = pair{pack[i].DataValue, ref[i].DataValue}
		// Sum
		if ref[i].Sum != nil {
			pairs["Sum"] = pair{*pack[i].Sum, *ref[i].Sum}
		}
		// Time
		pairs["Time"] = pair{pack[i].Time, ref[i].Time}
		// Update Time
		pairs["UpdateTime"] = pair{pack[i].UpdateTime, ref[i].UpdateTime}

		// compare values
		for fieldName, p := range pairs {
			if p.got != p.expected {
				return fmt.Errorf("assertion failed for %s: Got: '%v' instead of: '%v'", fieldName, p.got, p.expected)
			}
		}
	}
	return nil
}
