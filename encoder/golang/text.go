package golang

import (
	"github.com/K-Phoen/jennifer/jen"
	"github.com/K-Phoen/sdk"
)

func (encoder *Encoder) encodeText(panel sdk.Panel) jen.Code {
	settings := encoder.encodeCommonPanelProperties(panel, "text")

	if textMode(panel.TextPanel) == "markdown" {
		settings = append(settings,
			textQual("Markdown").Call(lit(textContent(panel.TextPanel))),
		)
	} else {
		settings = append(settings,
			textQual("HTML").Call(lit(textContent(panel.TextPanel))),
		)
	}

	return qual("row", "WithText").MultiLineCall(
		settings...,
	)
}

func textContent(panel *sdk.TextPanel) string {
	if panel.Options.Content != "" {
		return panel.Options.Content
	}

	return panel.Content
}

func textMode(panel *sdk.TextPanel) string {
	if panel.Options.Mode != "" {
		return panel.Options.Mode
	}

	return panel.Mode
}

func textQual(name string) *jen.Statement {
	return qual("text", name)
}
