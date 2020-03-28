package decoder

import (
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/text"
)

type DashboardText struct {
	Title    string
	Span     float32 `yaml:",omitempty"`
	Height   string  `yaml:",omitempty"`
	HTML     string  `yaml:",omitempty"`
	Markdown string  `yaml:",omitempty"`
}

func (textPanel DashboardText) toOption() row.Option {
	opts := []text.Option{}

	if textPanel.Span != 0 {
		opts = append(opts, text.Span(textPanel.Span))
	}
	if textPanel.Height != "" {
		opts = append(opts, text.Height(textPanel.Height))
	}
	if textPanel.Markdown != "" {
		opts = append(opts, text.Markdown(textPanel.Markdown))
	}
	if textPanel.HTML != "" {
		opts = append(opts, text.HTML(textPanel.HTML))
	}

	return row.WithText(textPanel.Title, opts...)
}
