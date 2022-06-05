package decoder

import (
	"github.com/K-Phoen/grabana/row"
	"github.com/K-Phoen/grabana/text"
)

type DashboardText struct {
	Title       string
	Description string              `yaml:",omitempty"`
	Span        float32             `yaml:",omitempty"`
	Height      string              `yaml:",omitempty"`
	Transparent bool                `yaml:",omitempty"`
	Links       DashboardPanelLinks `yaml:",omitempty"`
	HTML        string              `yaml:",omitempty"`
	Markdown    string              `yaml:",omitempty"`
}

func (textPanel DashboardText) toOption() row.Option {
	opts := []text.Option{}

	if textPanel.Description != "" {
		opts = append(opts, text.Description(textPanel.Description))
	}
	if textPanel.Span != 0 {
		opts = append(opts, text.Span(textPanel.Span))
	}
	if len(textPanel.Links) != 0 {
		opts = append(opts, text.Links(textPanel.Links.toModel()...))
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
	if textPanel.Transparent {
		opts = append(opts, text.Transparent())
	}

	return row.WithText(textPanel.Title, opts...)
}
