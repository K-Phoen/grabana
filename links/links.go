package links

import (
	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure a panel link.
type Option func(link *Link)

// Link represents a panel link.
type Link struct {
	Builder sdk.Link
}

// New creates a new logs panel.
func New(title string, url string, options ...Option) Link {
	link := &Link{Builder: sdk.Link{
		Title: title,
		URL:   &url,
	}}

	for _, opt := range options {
		opt(link)
	}

	return *link
}

// OpenBlank configures the link to open in a new tab.
func OpenBlank() Option {
	return func(link *Link) {
		yep := true
		link.Builder.TargetBlank = &yep
	}
}
