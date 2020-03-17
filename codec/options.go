package codec

type codecOptions struct {
	prettyPrint bool
	withHeader  bool
}

// PrettyPrint enables indentation of JSON and XML outputs
func PrettyPrint(o *codecOptions) {
	o.prettyPrint = true
}

// WithHeader enables CSV header
func WithHeader(o *codecOptions) {
	o.withHeader = true
}

// Option is the function type for setting codec options
type Option func(*codecOptions)
