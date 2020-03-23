package codec

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/farshidtz/senml/v2"
)

// DefaultCSVHeader is the default (currently fixed) CSV header
const DefaultCSVHeader = "Time,Update Time,Name,Unit,Value,String Value,Boolean Value,Data Value,Sum"

var WriteCSV Writer = func(p senml.Pack, w io.Writer, options ...Option) error {
	o := &codecOptions{
		header: false,
	}
	for _, opt := range options {
		opt(o)
	}

	csvWriter := csv.NewWriter(w)

	if o.header {
		err := csvWriter.Write(strings.Split(DefaultCSVHeader, ","))
		if err != nil {
			return err
		}
	}

	// normalize first to add base values to row values
	p.Normalize()

	for i := range p {
		row := make([]string, 9)
		row[0] = strconv.FormatFloat(p[i].Time, 'f', -1, 64)
		row[1] = strconv.FormatFloat(p[i].UpdateTime, 'f', -1, 64)
		row[2] = p[i].Name
		row[3] = p[i].Unit
		if p[i].Value != nil {
			row[4] = strconv.FormatFloat(*p[i].Value, 'f', -1, 64)
		}
		row[5] = p[i].StringValue
		if p[i].BoolValue != nil {
			row[6] = fmt.Sprintf("%t", *p[i].BoolValue)
		}
		row[7] = p[i].DataValue
		if p[i].Sum != nil {
			row[8] = strconv.FormatFloat(*p[i].Sum, 'f', -1, 64)
		}

		err := csvWriter.Write(row)
		if err != nil {
			return err
		}
		csvWriter.Flush()
	}
	if err := csvWriter.Error(); err != nil {
		return err
	}
	return nil
}

// EncodeCSV serializes the SenML pack into CSV bytes
var EncodeCSV Encoder = func(p senml.Pack, options ...Option) ([]byte, error) {

	var buf bytes.Buffer
	err := WriteCSV(p, &buf, options...)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

var ReadCSV Reader = func(r io.Reader, options ...Option) (senml.Pack, error) {
	o := &codecOptions{
		header: false,
	}
	for _, opt := range options {
		opt(o)
	}

	csvReader := csv.NewReader(r)

	if o.header {
		row, err := csvReader.Read()
		if err == io.EOF {
			return nil, fmt.Errorf("missing header or no input")
		}
		if err != nil {
			return nil, err
		}
		if joined := strings.Join(row, ","); joined != DefaultCSVHeader {
			return nil, fmt.Errorf("unexpected header: %s. Expected: %s", joined, DefaultCSVHeader)
		}
	}

	var p senml.Pack
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		var record senml.Record
		// Time
		record.Time, err = strconv.ParseFloat(row[0], 10)
		if err != nil {
			return nil, err
		}
		// Update Time
		record.UpdateTime, err = strconv.ParseFloat(row[1], 10)
		if err != nil {
			return nil, err
		}
		// Name
		record.Name = row[2]
		// Unit
		record.Unit = row[3]
		// Value
		if row[4] != "" {
			value, err := strconv.ParseFloat(row[4], 10)
			if err != nil {
				return nil, err
			}
			record.Value = &value
		}
		// String Value
		record.StringValue = row[5]
		// Boolean Value
		if row[6] != "" {
			boolValue, err := strconv.ParseBool(row[6])
			if err != nil {
				return nil, err
			}
			record.BoolValue = &boolValue
		}
		// Data Value
		record.DataValue = row[7]
		// Sum
		if row[8] != "" {
			sum, err := strconv.ParseFloat(row[8], 10)
			if err != nil {
				return nil, err
			}
			record.Sum = &sum
		}

		p = append(p, record)
	}

	return p, nil
}

// DecodeCSV takes a SenML pack in CSV bytes and decodes it into a Pack
var DecodeCSV Decoder = func(b []byte, options ...Option) (senml.Pack, error) {

	p, err := ReadCSV(bytes.NewReader(b), options...)
	if err != nil {
		return nil, err
	}

	return p, nil
}
