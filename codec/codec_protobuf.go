package codec

import (
	senmlprotobuf "github.com/farshidtz/senml-protobuf/go"
	"github.com/farshidtz/senml/v2"
	"github.com/golang/protobuf/proto"
)

// ConvertToProtobufMessage converts senml.Pack to senmlprotobuf.Message
func ConvertToProtobufMessage(p senml.Pack) senmlprotobuf.Message {
	var message senmlprotobuf.Message
	message.Pack = make([]*senmlprotobuf.Record, len(p))
	for i := range p {
		var r senmlprotobuf.Record
		// BaseName
		r.BaseName = p[i].BaseName
		// BaseTime
		r.BaseTime = p[i].BaseTime
		// BaseUnit
		r.BaseUnit = p[i].BaseUnit
		// BaseVersion
		if p[i].BaseVersion != nil {
			r.BaseVersion = int32(*p[i].BaseVersion)
		}
		// BaseValue
		if p[i].BaseValue != nil {
			r.BaseValueOptional = &senmlprotobuf.Record_BaseValue{*p[i].BaseValue}
		}
		// BaseSum
		if p[i].BaseSum != nil {
			r.BaseSumOptional = &senmlprotobuf.Record_BaseSum{*p[i].BaseSum}
		}
		//
		// Name
		r.Name = p[i].Name
		// Unit
		r.Unit = p[i].Unit
		// Time
		r.Time = p[i].Time
		// UpdateTime
		r.UpdateTime = p[i].UpdateTime
		//
		// Value, BoolValue, StringValue, DataValue
		if p[i].Value != nil {
			r.ValueOneof = &senmlprotobuf.Record_Value{*p[i].Value}
		} else if p[i].BoolValue != nil {
			r.ValueOneof = &senmlprotobuf.Record_BoolValue{*p[i].BoolValue}
		} else if p[i].StringValue != "" {
			r.ValueOneof = &senmlprotobuf.Record_StringValue{p[i].StringValue}
		} else if p[i].DataValue != "" {
			r.ValueOneof = &senmlprotobuf.Record_DataValue{p[i].DataValue}
		}
		// Sum
		if p[i].Sum != nil {
			r.SumOptional = &senmlprotobuf.Record_Sum{*p[i].Sum}
		}
		message.Pack[i] = &r
	}
	return message
}

// EncodeProtobuf serializes the SenML pack into Protobuf bytes. The options are ignored.
func EncodeProtobuf(p senml.Pack, _ ...Option) ([]byte, error) {
	message := ConvertToProtobufMessage(p)
	return proto.Marshal(&message)
}

// ConvertFromProtobufMessage coverts senmlprotobuf.Message to senml.Pack
func ConvertFromProtobufMessage(message senmlprotobuf.Message) senml.Pack {
	var p = make(senml.Pack, len(message.Pack))
	for i := range message.Pack {
		// BaseName
		p[i].BaseName = message.Pack[i].BaseName
		// BaseTime
		p[i].BaseTime = message.Pack[i].BaseTime
		// BaseUnit
		p[i].BaseUnit = message.Pack[i].BaseUnit
		// BaseVersion
		if message.Pack[i].BaseVersion != 0 {
			bver := int(message.Pack[i].BaseVersion)
			p[i].BaseVersion = &bver
		}
		// BaseValue
		if v, ok := message.Pack[i].BaseValueOptional.(*senmlprotobuf.Record_BaseValue); ok {
			p[i].BaseValue = &v.BaseValue
		}
		// BaseSum
		if v, ok := message.Pack[i].BaseSumOptional.(*senmlprotobuf.Record_BaseSum); ok {
			p[i].BaseSum = &v.BaseSum
		}

		// Name
		p[i].Name = message.Pack[i].Name
		// Unit
		p[i].Unit = message.Pack[i].Unit
		// Time
		p[i].Time = message.Pack[i].Time
		// UpdateTime
		p[i].UpdateTime = message.Pack[i].UpdateTime
		//
		// Value, BoolValue, StringValue, DataValue
		switch message.Pack[i].ValueOneof.(type) {
		case *senmlprotobuf.Record_Value:
			p[i].Value = &message.Pack[i].ValueOneof.(*senmlprotobuf.Record_Value).Value
		case *senmlprotobuf.Record_BoolValue:
			p[i].BoolValue = &message.Pack[i].ValueOneof.(*senmlprotobuf.Record_BoolValue).BoolValue
		case *senmlprotobuf.Record_StringValue:
			p[i].StringValue = message.Pack[i].ValueOneof.(*senmlprotobuf.Record_StringValue).StringValue
		case *senmlprotobuf.Record_DataValue:
			p[i].DataValue = message.Pack[i].ValueOneof.(*senmlprotobuf.Record_DataValue).DataValue
		}
		// Sum
		if v, ok := message.Pack[i].SumOptional.(*senmlprotobuf.Record_Sum); ok {
			p[i].Sum = &v.Sum
		}
	}
	return p
}

// DecodeProtobuf takes a SenML pack in Protobuf bytes and decodes it into a Pack. The options are ignored.
func DecodeProtobuf(b []byte, _ ...Option) (senml.Pack, error) {
	var message senmlprotobuf.Message
	err := proto.Unmarshal(b, &message)
	if err != nil {
		return nil, err
	}

	return ConvertFromProtobufMessage(message), nil
}
