package objects

type printOpFunction func(*printOptions)

type printOptions struct {
	indent     int
	keysOnly   bool
	valuesOnly bool
}

func NewPrintOptions(opOptions []printOpFunction) *printOptions {
	opts := &printOptions{
		indent: 2,
	}

	for _, opt := range opOptions {
		opt(opts)
	}

	return opts
}

func WithIndent(indent uint) printOpFunction {
	return func(opts *printOptions) {
		if indent > (1 << 16) {
			indent = (1 << 16) - 1
		}

		opts.indent = int(indent)
	}
}

func WithoutIndent() printOpFunction {
	return func(opts *printOptions) {
		opts.indent = 0
	}
}

func WithKeysOnly() printOpFunction {
	return func(opts *printOptions) {
		opts.keysOnly = true
	}
}

func WithValuesOnly() printOpFunction {
	return func(opts *printOptions) {
		opts.valuesOnly = true
	}
}
