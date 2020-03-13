package senml

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// CSVHeader is the fixed header to support records with different value types
const CSVHeader = "Time,Name,Unit,Value,String Value,Boolean Value,Data Value,Sum,Update Time"

func (p Pack) WriteCSV(w io.Writer, header bool) error {

	csvWriter := csv.NewWriter(w)

	if header {
		err := csvWriter.Write(strings.Split(CSVHeader, ","))
		if err != nil {
			return err
		}
	}

	// normalize first to add base values to row values
	p.Normalize()

	for i := range p {
		row := make([]string, 9)
		row[0] = strconv.FormatFloat(p[i].Time, 'f', -1, 64)
		row[1] = p[i].Name
		row[2] = p[i].Unit
		if p[i].Value != nil {
			row[3] = strconv.FormatFloat(*p[i].Value, 'f', -1, 64)
		}
		row[4] = p[i].StringValue
		if p[i].BoolValue != nil {
			row[5] = fmt.Sprintf("%t", *p[i].BoolValue)
		}
		row[6] = p[i].DataValue
		if p[i].Sum != nil {
			row[7] = strconv.FormatFloat(*p[i].Sum, 'f', -1, 64)
		}
		row[8] = strconv.FormatFloat(p[i].UpdateTime, 'f', -1, 64)

		err := csvWriter.Write(row)
		if err != nil {
			return err
		}
	}
	csvWriter.Flush() // TODO flush during the iterations?
	if err := csvWriter.Error(); err != nil {
		return err
	}
	return nil
}

// EncodeCSV serializes the SenML pack into CSV bytes
func (p Pack) EncodeCSV(header bool) ([]byte, error) {

	var buf bytes.Buffer
	err := p.WriteCSV(&buf, header)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ReadCSV(r io.Reader, header bool) (Pack, error) {
	csvReader := csv.NewReader(r)

	if header {
		row, err := csvReader.Read()
		if err == io.EOF {
			return nil, fmt.Errorf("missing header or no input")
		}
		if err != nil {
			return nil, err
		}
		if joined := strings.Join(row, ","); joined != CSVHeader {
			return nil, fmt.Errorf("unexpected header: %s. Expected: %s", joined, CSVHeader)
		}
	}

	var p Pack
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		var record Record
		// Time
		record.Time, err = strconv.ParseFloat(row[0], 10)
		if err != nil {
			return nil, err
		}
		// Name
		record.Name = row[1]
		// Unit
		record.Unit = row[2]
		// Value
		if row[3] != "" {
			value, err := strconv.ParseFloat(row[3], 10)
			if err != nil {
				return nil, err
			}
			record.Value = &value
		}
		// String Value
		record.StringValue = row[4]
		// Boolean Value
		if row[5] != "" {
			boolValue, err := strconv.ParseBool(row[5])
			if err != nil {
				return nil, err
			}
			record.BoolValue = &boolValue
		}
		// Data Value
		record.DataValue = row[6]
		// Sum
		if row[7] != "" {
			sum, err := strconv.ParseFloat(row[7], 10)
			if err != nil {
				return nil, err
			}
			record.Sum = &sum
		}
		// Update Time
		record.UpdateTime, err = strconv.ParseFloat(row[8], 10)
		if err != nil {
			return nil, err
		}

		p = append(p, record)
	}

	return p, nil
}

// DecodeCSV takes a SenML pack in CSV bytes and decodes it into a Pack
func DecodeCSV(b []byte, header bool) (Pack, error) {

	p, err := ReadCSV(bytes.NewReader(b), header)
	if err != nil {
		return nil, err
	}

	return p, nil
}
