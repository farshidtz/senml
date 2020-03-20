package senml

import (
	"encoding/json"
	"math"
	"testing"
	"time"
)

func referencePack(absoluteTime ...bool) Pack {
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
	return Pack{
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

func referencePackFloats() Pack {
	bver := 5
	value, value2 := 22.1, 30.0
	sum, sum2 := 100.0, 200.0
	return Pack{
		{BaseName: "dev123",
			BaseTime:    -45.67,
			BaseUnit:    "degC",
			BaseVersion: &bver,
			Value:       &value, Unit: "degC", Name: "temp", Time: -1.0, UpdateTime: 10.0, Sum: &sum},
		{Value: &value2, Time: 1.0, Sum: &sum2},
	}
}

func stringifyPack(p Pack) string {
	b, _ := json.Marshal(&p)
	return string(b)
}

func TestNormalize(t *testing.T) {

	t.Run("Positive relative time", func(t *testing.T) {
		p := referencePack()
		p[0].BaseTime = 1000
		p.Normalize()
		now := float64(time.Now().UnixNano()) / 1000000000
		expected := now + 1000
		if math.Abs(p[0].Time-expected) > 5 { // fail if difference is more than 5s
			t.Fatalf("Time is not absolute. Got %f instead of %f", p[0].Time, expected)
		}
	})

	t.Run("Negative relative time", func(t *testing.T) {
		// negative relative time
		p := referencePack()
		p[0].BaseTime = -1000
		p.Normalize()
		now := float64(time.Now().UnixNano()) / 1000000000
		expected := now - 1000
		if math.Abs(p[0].Time-expected) > 5 { // fail if difference is more than 5s
			t.Fatalf("Time is not absolute. Got %f instead of %f", p[0].Time, expected)
		}
	})

	t.Run("Absolute time", func(t *testing.T) {
		p := referencePack(true)
		base := p[0].BaseTime
		p.Normalize()

		ref := referencePack()
		for i := range p {
			expected := ref[i].Time + base
			if p[i].Time != expected {
				t.Fatalf("Time is not absolute. Got %f instead of %f", p[i].Time, expected)
			}
		}
	})

	t.Run("Default base version in first record", func(t *testing.T) {
		p := referencePack()
		*p[0].BaseVersion = DefaultBaseVersion
		p.Normalize()
		for i, r := range p {
			if r.BaseVersion != nil {
				t.Errorf("Default base version was not omitted in record %d: %s", i, stringifyPack(p))
			}
		}
		if t.Failed() {
			t.FailNow()
		}
	})

	t.Run("Default base version in second record", func(t *testing.T) {
		p := referencePack()
		p[0].BaseVersion = nil
		p[1].BaseVersion = new(int)
		*p[1].BaseVersion = DefaultBaseVersion
		p.Normalize()
		for i, r := range p {
			if r.BaseVersion != nil {
				t.Errorf("Default base version was not omitted in record %d: %s", i, stringifyPack(p))
			}
		}
		if t.Failed() {
			t.FailNow()
		}
	})

	t.Run("No base version", func(t *testing.T) {
		p := referencePack()
		p[0].BaseVersion = nil
		p.Normalize()
		for i, r := range p {
			if r.BaseVersion != nil {
				t.Errorf("Default base version was not omitted in record %d: %s", i, stringifyPack(p))
			}
		}
		if t.Failed() {
			t.FailNow()
		}
	})

	t.Run("Floats pack with base value", func(t *testing.T) {
		p := referencePackFloats()
		p[0].Value = nil
		p[1].Value = nil
		p[0].BaseValue = new(float64)
		*p[0].BaseValue = 10
		p.Normalize()
		if *p[0].Value != 10 && *p[1].Value != 10 {
			t.Fatalf("Base value was not added to value in records. Got values: %f, %f", *p[0].Value, *p[1].Value)
		}
		if p[0].BaseValue != nil {
			t.Fatalf("Base value was not removed from record in pack: %s", stringifyPack(p))
		}
	})

	t.Run("Floats pack with base value and values", func(t *testing.T) {
		p := referencePackFloats()
		p[0].BaseValue = new(float64)
		*p[0].BaseValue = 10
		p.Normalize()
		if *p[0].Value != 32.1 {
			t.Fatalf("Base value was not added to value in first record. Got value: %f", *p[0].Value)
		}
		if p[0].BaseValue != nil {
			t.Fatalf("Base value was not removed from record in pack: %s", stringifyPack(p))
		}
		if *p[1].Value != 40 {
			t.Fatalf("Base value was not added to value in second record. Got value: %f", *p[1].Value)
		}
	})

	t.Run("Floats pack with base sum", func(t *testing.T) {
		p := referencePackFloats()
		p[0].Sum = nil
		p[1].Sum = nil
		p[0].BaseSum = new(float64)
		*p[0].BaseSum = 10
		p.Normalize()
		if *p[0].Sum != 10 && *p[1].Sum != 10 {
			t.Fatalf("Base sum was not added to sum in records. Got sums: %f, %f", *p[0].Sum, *p[1].Sum)
		}
		if p[0].BaseSum != nil {
			t.Fatalf("Base sum was not removed from record in pack: %s", stringifyPack(p))
		}
	})

	t.Run("Floats pack with base sum and sums", func(t *testing.T) {
		p := referencePackFloats()
		p[0].BaseSum = new(float64)
		*p[0].BaseSum = 10
		p.Normalize()
		if *p[0].Sum != 110 {
			t.Fatalf("Base sum was not added to sum in first record. Got sum: %f", *p[0].Sum)
		}
		if p[0].BaseSum != nil {
			t.Fatalf("Base sum was not removed from record in pack: %s", stringifyPack(p))
		}
		if *p[1].Sum != 210 {
			t.Fatalf("Base sum was not added to sum in second record. Got sum: %f", *p[1].Sum)
		}
	})

}

func TestClone(t *testing.T) {
	p := referencePack()
	p[0].XMLName = new(bool)

	c := p.Clone()

	*p[0].XMLName = true
	*p[0].Value = 123.456
	*p[0].BaseVersion = 123
	p[0].Time = 123
	p[0].StringValue = "changed"

	if *p[0].XMLName == *c[0].XMLName ||
		*p[0].Value == *c[0].Value ||
		*p[0].BaseVersion == *c[0].BaseVersion ||
		p[0].Time == c[0].Time ||
		p[0].StringValue == c[0].StringValue {
		t.Fatalf("Clone is changed after changing the original pack.")
	}
}

func TestValidate(t *testing.T) {

	t.Run("no value", func(t *testing.T) {
		pack := Pack{
			{Name: "dev"},
		}
		err := pack.Validate()
		if err == nil {
			t.Fatalf("No error for record with no value: %s", stringifyPack(pack))
		}
	})

	t.Run("numeric name", func(t *testing.T) {
		value := 1.0
		pack := Pack{
			{Name: "3a", Value: &value},
		}
		err := pack.Validate()
		if err != nil {
			t.Fatalf("Error decoding record with numeric name: %s", err)
		}
	})

	t.Run("bad numeric name", func(t *testing.T) {
		value := 1.0
		pack := Pack{
			{Name: "-3a", Value: &value},
		}
		err := pack.Validate()
		if err == nil {
			t.Fatalf("No error for bad numeric name in: %s", stringifyPack(pack))
		}
	})

	t.Run("weird name", func(t *testing.T) {
		value := 1.0
		pack := Pack{
			{Name: "Az3-:./_", Value: &value},
		}
		err := pack.Validate()
		if err != nil {
			t.Fatalf("Error decoding record with valid name: %s", err)
		}
	})

	t.Run("bad weird name", func(t *testing.T) {
		value := 1.0
		pack := Pack{
			{Name: "A;b", Value: &value},
		}
		err := pack.Validate()
		if err == nil {
			t.Fatalf("No error for invalid name in: %s", stringifyPack(pack))
		}
	})

	t.Run("weird base name", func(t *testing.T) {
		value := 1.0
		pack := Pack{
			{BaseName: "Az3-:./_", Name: "/b", Value: &value},
		}
		err := pack.Validate()
		if err != nil {
			t.Fatalf("Error decoding record with valid base name: %s", err)
		}
	})

	t.Run("bad numeric base name", func(t *testing.T) {
		value := 1.0
		pack := Pack{
			{BaseName: "/room", Name: "/dev", Value: &value},
		}
		err := pack.Validate()
		if err == nil {
			t.Fatalf("No error for invalid numeric base name in: %s", stringifyPack(pack))
		}
		//
		pack[0].BaseName = "room#3"
		err = pack.Validate()
		if err == nil {
			t.Fatalf("No error for invalid numeric base name in: %s", stringifyPack(pack))
		}
	})

	t.Run("multiple values in record", func(t *testing.T) {
		value := 1.0
		pack := Pack{
			{Name: "dev", Value: &value, StringValue: "on"},
		}
		err := pack.Validate()
		if err == nil {
			t.Fatalf("No error for multi-valued record in pack: %s", stringifyPack(pack))
		}
	})

	t.Run("base value with other types in pack", func(t *testing.T) {
		bval := 1.0
		pack := Pack{
			{Name: "dev", BaseValue: &bval, StringValue: "on"},
		}
		err := pack.Validate()
		if err == nil {
			t.Fatalf("No error for record with base value (float) and another non-float value in pack: %s", stringifyPack(pack))
		}
	})

	t.Run("multiple base versions in pack", func(t *testing.T) {
		value := 1.0
		bver := 5
		pack := Pack{
			{Name: "dev", Value: &value, BaseVersion: nil},
			{Name: "dev", Value: &value, BaseVersion: &bver},
		}
		err := pack.Validate()
		if err == nil {
			t.Fatalf("No error for pack with no version followed by custom version: %s", stringifyPack(pack))
		}
		//
		bver_default := DefaultBaseVersion
		pack = Pack{
			{Name: "dev", Value: &value, BaseVersion: nil},
			{Name: "dev", Value: &value, BaseVersion: &bver_default},
		}
		err = pack.Validate()
		if err == nil {
			t.Fatalf("No error for pack with no version followed by default version: %s", stringifyPack(pack))
		}
		pack = Pack{
			{Name: "dev", Value: &value, BaseVersion: &bver_default},
			{Name: "dev", Value: &value, BaseVersion: &bver},
		}
		err = pack.Validate()
		if err == nil {
			t.Fatalf("No error for pack with default followed by custom version: %s", stringifyPack(pack))
		}
		//
		pack = Pack{
			{Name: "dev", Value: &value, BaseVersion: &bver},
			{Name: "dev", Value: &value, BaseVersion: &bver_default},
		}
		err = pack.Validate()
		if err == nil {
			t.Fatalf("No error for pack with custom followed by default version: %s", stringifyPack(pack))
		}
	})

	t.Run("custom base version", func(t *testing.T) {
		value := 1.0
		bver := 5
		pack := Pack{
			{Name: "dev", Value: &value, BaseVersion: &bver},
			{Name: "dev", Value: &value, BaseVersion: nil},
		}
		err := pack.Validate()
		if err != nil {
			t.Fatalf("Error for pack with custom followed by no version: %s", err)
		}
	})

	t.Run("base value only", func(t *testing.T) {
		bval := 1.0
		pack := Pack{
			{Name: "dev", BaseValue: &bval},
		}
		err := pack.Validate()
		if err != nil {
			t.Fatalf("Error for pack with base value: %s", err)
		}
	})

	t.Run("sum only", func(t *testing.T) {
		sum := 1.0
		pack := Pack{
			{Name: "dev", Sum: &sum},
		}
		err := pack.Validate()
		if err != nil {
			t.Fatalf("Error for pack with base sum: %s", err)
		}
	})

	t.Run("base sum only", func(t *testing.T) {
		bsum := 1.0
		pack := Pack{
			{Name: "dev", BaseSum: &bsum},
		}
		err := pack.Validate()
		if err != nil {
			t.Fatalf("Error for pack with base sum: %s", err)
		}
		//
		sum := 10.0
		pack = Pack{
			{Name: "dev", BaseSum: &bsum},
			{Name: "dev", Sum: &sum},
		}
		err = pack.Validate()
		if err != nil {
			t.Fatalf("Error for pack with base sum and sum: %s", err)
		}
	})
}

func TestValidateName(t *testing.T) {
	t.Run("valid names", func(t *testing.T) {
		names := []string{"Aa-:./_", "urn:dev:ow:10e2073a", "http://example.com"}
		for _, name := range names {
			err := ValidateName(name)
			if err != nil {
				t.Fatalf("Error for valid name: %s: %s", name, err)
			}
		}
	})

	t.Run("invalid names", func(t *testing.T) {
		names := []string{"-A", ":A", ".A", "/A", "_A",
			"~A", "!A", "@A", "#A", "$A", "%A", "^A", "&A", "*A", "(A", "+A", "=A", " A", " ", "A "}
		for _, name := range names {
			err := ValidateName(name)
			if err == nil {
				t.Fatalf("No error for invalid name: %s", name)
			}
		}
	})
}
