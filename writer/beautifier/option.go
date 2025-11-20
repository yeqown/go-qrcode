package beautifier

type options struct {
	patterns []patternDesc
}

type Option func(*options)

func WithPattern(pa patternDesc, resourceFile string) Option {
	return func(o *options) {
		if pa == nil || resourceFile == "" {
			return
		}

	}
}
