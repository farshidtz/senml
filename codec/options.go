package codec

type codecOptions struct {
	prettyPrint bool
	header      bool
}

// SetPrettyPrint enables indentation for JSON and XML encoding
func SetPrettyPrint(o *codecOptions) {
	o.prettyPrint = true
}

// SetDefaultHeader enables header for CSV encoding/decoding
func SetDefaultHeader(o *codecOptions) {
	o.header = true
}

// TODO set custom CSV header
//func SetHeader(header string) Option {
//	return func(o *codecOptions) {
//		o.header = header
//	}
//}

// Option is the function type for setting codec options
type Option func(*codecOptions)
