package senml

type Options struct {
	PrettyPrint bool
	Header      bool
}

func PrettyPrint(enable bool) Option {
	return func(f *Options) {
		f.PrettyPrint = enable
	}
}

func Header(enable bool) Option {
	return func(f *Options) {
		f.Header = enable
	}
}

type Option func(*Options)

type Coder interface {
	Encode(Pack, ...Option) ([]byte, error)
	Decode([]byte, ...Option) (Pack, error)
}
