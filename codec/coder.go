package codec

type Options struct {
	// PrettyPrint enables indentation of JSON and XML outputs
	PrettyPrint bool
	// WithHeader enables CSV header
	WithHeader bool
}

// PrettyPrint enables indentation of JSON and XML outputs
func PrettyPrint(o *Options) {
	o.PrettyPrint = true
}

// WithHeader enables CSV header
func WithHeader(o *Options) {
	o.WithHeader = true
}

type Option func(*Options)
