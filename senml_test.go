package senml

import (
	"encoding/base64"
	"fmt"
	"math"
	"time"

	"testing"
)

func ExampleEncode1() {
	v := 23.1
	var p Pack = []Record{
		{Value: &v, Unit: "Cel", Name: "urn:dev:ow:10e2073a01080063"},
	}

	dataOut, err := p.Encode(JSON, OutputOptions{})
	if err != nil {
		fmt.Println("Encode of SenML failed")
	} else {
		fmt.Println(string(dataOut))
	}
	// Output: [{"n":"urn:dev:ow:10e2073a01080063","u":"Cel","v":23.1}]
}

func ExampleEncode2() {
	v1 := 23.5
	v2 := 23.6
	var p Pack = []Record{
		{Value: &v1, Unit: "Cel", BaseName: "urn:dev:ow:10e2073a01080063", Time: 1.276020076305e+09},
		{Value: &v2, Unit: "Cel", Time: 1.276020091305e+09},
	}

	dataOut, err := p.Encode(JSON, OutputOptions{})
	if err != nil {
		fmt.Println("Encode of SenML failed")
	} else {
		fmt.Println(string(dataOut))
	}
	// Output: [{"bn":"urn:dev:ow:10e2073a01080063","u":"Cel","t":1276020076.305,"v":23.5},{"u":"Cel","t":1276020091.305,"v":23.6}]
}

type TestVector struct {
	testDecode bool
	format     Format
	label      string
	binary     bool
	value      string
}

var testVectors = []TestVector{
	{true, JSON, "JSON", false, "W3siYm4iOiJkZXYxMjMiLCJidCI6LTQ1LjY3LCJidSI6ImRlZ0MiLCJidmVyIjo1LCJuIjoidGVtcCIsInUiOiJkZWdDIiwidCI6LTEsInV0IjoxMCwidiI6MjIuMSwicyI6MH0seyJuIjoicm9vbSIsInQiOi0xLCJ2cyI6ImtpdGNoZW4ifSx7Im4iOiJkYXRhIiwidmQiOiJhYmMifSx7Im4iOiJvayIsInZiIjp0cnVlfV0="},
	{true, CBOR, "CBOR", true, "hKpiYm5mZGV2MTIzYmJ0+8BG1cKPXCj2YmJ1ZGRlZ0NkYnZlcgVhbmR0ZW1wYXP7AAAAAAAAAABhdPu/8AAAAAAAAGF1ZGRlZ0NidXT7QCQAAAAAAABhdvtANhmZmZmZmqNhbmRyb29tYXT7v/AAAAAAAABidnNna2l0Y2hlbqJhbmRkYXRhYnZkY2FiY6JhbmJva2J2YvU="},
	{true, XML, "XML", false, "PHNlbnNtbCB4bWxucz0idXJuOmlldGY6cGFyYW1zOnhtbDpuczpzZW5tbCI+PHNlbm1sIGJuPSJkZXYxMjMiIGJ0PSItNDUuNjciIGJ1PSJkZWdDIiBidmVyPSI1IiBuPSJ0ZW1wIiB1PSJkZWdDIiB0PSItMSIgdXQ9IjEwIiB2PSIyMi4xIiBzPSIwIj48L3Nlbm1sPjxzZW5tbCBuPSJyb29tIiB0PSItMSIgdnM9ImtpdGNoZW4iPjwvc2VubWw+PHNlbm1sIG49ImRhdGEiIHZkPSJhYmMiPjwvc2VubWw+PHNlbm1sIG49Im9rIiB2Yj0idHJ1ZSI+PC9zZW5tbD48L3NlbnNtbD4="},
	{false, CSV, "CSV", false, "ZGV2MTIzdGVtcCw5NDY2ODQ3OTkuMDAwMDAwLDIyLjEwMDAwMCxkZWdDDQo="},
	{true, MPACK, "MPACK", true, "lIqiYm6mZGV2MTIzomJ0y8BG1cKPXCj2omJ1pGRlZ0OkYnZlcgWhbqR0ZW1woXPLAAAAAAAAAAChdMu/8AAAAAAAAKF1pGRlZ0OidXTLQCQAAAAAAAChdstANhmZmZmZmoOhbqRyb29toXTLv/AAAAAAAACidnOna2l0Y2hlboKhbqRkYXRhonZko2FiY4KhbqJva6J2YsM="},
	{false, LINEP, "LINEP", false, "Zmx1ZmZ5U2VubWwsbj10ZW1wLHU9ZGVnQyB2PTIyLjEgLTEwMDAwMDAwMDAK"},
}

func referencePack() Pack {
	value := 22.1
	sum := 0.0
	vb := true
	return Pack{
		{BaseName: "dev123",
			BaseTime:    -45.67,
			BaseUnit:    "degC",
			BaseVersion: 5,
			Value:       &value, Unit: "degC", Name: "temp", Time: -1.0, UpdateTime: 10.0, Sum: &sum},
		{StringValue: "kitchen", Name: "room", Time: -1.0},
		{DataValue: "abc", Name: "data"},
		{BoolValue: &vb, Name: "ok"},
	}
}

func TestEncode(t *testing.T) {

	options := OutputOptions{Topic: "fluffySenml", PrettyPrint: false}
	for _, vector := range testVectors {
		ref := referencePack()
		if vector.label == "CSV" {
			// change to an absolute time: https://tools.ietf.org/html/rfc8428#section-4.5.3
			ref[0].BaseTime = 946684800
		}
		dataOut, err := ref.Encode(vector.format, options)
		if err != nil {
			t.Fatalf("Encoding error: %s", err)
		}

		if base64.StdEncoding.EncodeToString(dataOut) != vector.value {
			t.Errorf("Assertion failed for encoded %s:", vector.label)
			t.Logf("Got (encoded): %s", base64.StdEncoding.EncodeToString(dataOut))
			if !vector.binary {
				t.Logf("Got:\n'%s'", dataOut)
				decoded, err := base64.StdEncoding.DecodeString(vector.value)
				if err != nil {
					t.Fatalf("Error decoding test value: %s", err)
				}
				t.Fatalf("Expected:\n'%s'", decoded)
			} else {
				t.Logf("Got:\n%v", dataOut)
				decoded, err := base64.StdEncoding.DecodeString(vector.value)
				if err != nil {
					t.Fatalf("Error decoding test value: %s", err)
				}
				t.Fatalf("Expected:\n%v", decoded)
			}
		}
	}

}

func TestDecode(t *testing.T) {
	type pair struct {
		got      interface{}
		expected interface{}
	}
	ref := referencePack()
	for _, vector := range testVectors {
		if vector.testDecode {
			data, err := base64.StdEncoding.DecodeString(vector.value)
			if err != nil {
				t.Fatalf("Error decoding test value for %s: %s", vector.label, err)
			}

			pack, err := Decode(data, vector.format)
			if err != nil {
				t.Fatalf("Error decoding %s: %s", vector.label, err)
			}

			pairs := make(map[string]pair)
			for i := range pack {
				pairs["XMLName"] = pair{pack[i].XMLName, ref[i].XMLName}
				pairs["BaseName"] = pair{pack[i].BaseName, ref[i].BaseName}
				pairs["BaseTime"] = pair{pack[i].BaseTime, ref[i].BaseTime}
				pairs["BaseUnit"] = pair{pack[i].BaseUnit, ref[i].BaseUnit}
				pairs["BaseVersion"] = pair{pack[i].BaseVersion, ref[i].BaseVersion}
				pairs["Name"] = pair{pack[i].Name, ref[i].Name}
				pairs["Unit"] = pair{pack[i].Unit, ref[i].Unit}
				pairs["Time"] = pair{pack[i].Time, ref[i].Time}
				pairs["UpdateTime"] = pair{pack[i].UpdateTime, ref[i].UpdateTime}
				pairs["StringValue"] = pair{pack[i].StringValue, ref[i].StringValue}
				pairs["DataValue"] = pair{pack[i].DataValue, ref[i].DataValue}
				// pointers
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
						t.Logf("Assertion failed for %s in encoded %s:", fieldName, vector.label)
						t.Fatalf("Got: '%v' instead of: '%v'", p.got, p.expected)
					}
				}
			}
		}
	}
}

func TestNormalize(t *testing.T) {
	ref := referencePack()

	// positive relative time
	ref[0].BaseTime = 1000
	normalized := ref.Normalize()
	now := float64(time.Now().UnixNano()) / 1000000000
	expected := now + ref[0].BaseTime
	if math.Abs(normalized[0].Time-expected) > 5 { // fail if difference is more than 5s
		t.Fatalf("Time is not relative. Got %f instead of %f", ref[0].Time, expected)
	}

	// negative relative time
	ref[0].BaseTime = -1000
	normalized = ref.Normalize()
	now = float64(time.Now().UnixNano()) / 1000000000
	expected = now + ref[0].BaseTime
	if math.Abs(normalized[0].Time-expected) > 5 { // fail if difference is more than 5s
		t.Fatalf("Time is not relative. Got %f instead of %f", ref[0].Time, expected)
	}

	// absolute time
	ref[0].BaseTime = 946684800.123
	normalized = ref.Normalize()
	dataOut, err := normalized.Encode(JSON, OutputOptions{PrettyPrint: true})
	if err != nil {
		t.Fatalf("Error encoding: %s", err)
	}

	testValue := "WwogIHsiYnZlciI6NSwibiI6ImRldjEyM3RlbXAiLCJ1IjoiZGVnQyIsInQiOjk0NjY4NDc5OS4xMjMsInV0IjoxMCwidiI6MjIuMSwicyI6MH0sCiAgeyJidmVyIjo1LCJuIjoiZGV2MTIzcm9vbSIsInUiOiJkZWdDIiwidCI6OTQ2Njg0Nzk5LjEyMywidnMiOiJraXRjaGVuIn0sCiAgeyJidmVyIjo1LCJuIjoiZGV2MTIzZGF0YSIsInUiOiJkZWdDIiwidCI6OTQ2Njg0ODAwLjEyMywidmQiOiJhYmMifSwKICB7ImJ2ZXIiOjUsIm4iOiJkZXYxMjNvayIsInUiOiJkZWdDIiwidCI6OTQ2Njg0ODAwLjEyMywidmIiOnRydWV9Cl0K"
	if base64.StdEncoding.EncodeToString(dataOut) != testValue {
		t.Logf("Got (encoded): %s", base64.StdEncoding.EncodeToString(dataOut))
		t.Errorf("Assertion failed for normalized pack. Got:\n'%s'", dataOut)
		decoded, err := base64.StdEncoding.DecodeString(testValue)
		if err != nil {
			t.Fatalf("Error decoding test value: %s", err)
		}
		t.Errorf("Expected:\n'%s'", decoded)
	}
}

func TestBadInput1(t *testing.T) {
	data := []byte(" foo ")
	_, err := Decode(data, JSON)
	if err == nil {
		t.Fail()
	}
}

func TestBadInput2(t *testing.T) {
	data := []byte(" { \"n\":\"hi\" } ")
	_, err := Decode(data, JSON)
	if err == nil {
		t.Fail()
	}
}

func TestBadInputNoValue(t *testing.T) {
	data := []byte("  [ { \"n\":\"hi\" } ] ")
	_, err := Decode(data, JSON)
	if err == nil {
		t.Fail()
	}
}

func TestInputNumericName(t *testing.T) {
	data := []byte("  [ { \"n\":\"3a\", \"v\":1.0 } ] ")
	_, err := Decode(data, JSON)
	if err != nil {
		t.Fail()
	}
}

func TestBadInputNumericName(t *testing.T) {
	data := []byte("  [ { \"n\":\"-3b\", \"v\":1.0 } ] ")
	_, err := Decode(data, JSON)
	if err == nil {
		t.Fail()
	}
}

func TestInputWeirdName(t *testing.T) {
	data := []byte("  [ { \"n\":\"Az3-:./_\", \"v\":1.0 } ] ")
	_, err := Decode(data, JSON)
	if err != nil {
		t.Fail()
	}
}

func TestBadInputWeirdName(t *testing.T) {
	data := []byte("  [ { \"n\":\"A;b\", \"v\":1.0 } ] ")
	_, err := Decode(data, JSON)
	if err == nil {
		t.Fail()
	}
}

func TestInputWeirdBaseName(t *testing.T) {
	data := []byte("[ { \"bn\": \"a\" , \"n\":\"/b\" , \"v\":1.0} ] ")
	_, err := Decode(data, JSON)
	if err != nil {
		t.Fail()
	}
}

func TestBadInputNumericBaseName(t *testing.T) {
	data := []byte("[ { \"bn\": \"/3h\" , \"n\":\"i\" , \"v\":1.0} ] ")
	_, err := Decode(data, JSON)
	if err == nil {
		t.Fail()
	}
	data = []byte("[ { \"bn\": \"3h#\" , \"n\":\"i\" , \"v\":1.0} ] ")
	_, err = Decode(data, JSON)
	if err == nil {
		t.Fail()
	}
}

// TODO add
//func TestBadInputUnknownMtuField(t *testing.T) {
//	data := []byte("[ { \"n\":\"hi\", \"v\":1.0, \"mtu_\":1.0  } ] ")
//	_ , err := Decode(data, JSON)
//	if err == nil {
//		t.Fail()
//	}
//}

func TestInputSumOnly(t *testing.T) {
	data := []byte("[ { \"n\":\"a\", \"s\":1.0 } ] ")
	_, err := Decode(data, JSON)
	if err != nil {
		t.Fail()
	}
}

func TestInputBoolean(t *testing.T) {
	data := []byte("[ { \"n\":\"a\", \"vd\": \"aGkgCg\" } ] ")
	_, err := Decode(data, JSON)
	if err != nil {
		t.Fail()
	}
}

func TestInputData(t *testing.T) {
	data := []byte("  [ { \"n\":\"a\", \"vb\": true } ] ")
	_, err := Decode(data, JSON)
	if err != nil {
		t.Fail()
	}
}

func TestInputString(t *testing.T) {
	data := []byte("  [ { \"n\":\"a\", \"vs\": \"Hi\" } ] ")
	_, err := Decode(data, JSON)
	if err != nil {
		t.Fail()
	}
}
