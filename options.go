package compass

import tea "charm.land/bubbletea/v2"

type OptionFunc func(opts *Options)

type Options struct {
	// AutoExitOnEmpty quits the when the navigation stack is empty.
	AutoQuitOnEmpty bool

	// FallbackView is the view shown when the navigation stack is empty.
	FallbackView tea.View
}

func WithOptions(options Options) OptionFunc {
	return func(opts *Options) {
		*opts = options
	}
}

func WithAutoQuitOnEmpty(enabled bool) OptionFunc {
	return func(opts *Options) {
		opts.AutoQuitOnEmpty = enabled
	}
}

func WithFallbackView(view tea.View) OptionFunc {
	return func(opts *Options) {
		opts.FallbackView = view
	}
}
