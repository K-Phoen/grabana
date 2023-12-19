package decoder

import (
	"fmt"

	"github.com/K-Phoen/sdk"
)

func intPtr(input int) *int {
	return &input
}

func float64Ptr(input float64) *float64 {
	return &input
}

func parsePanelRepeatDirection(input string) (sdk.RepeatDirection, error) {
	switch input {
	case "vertical":
		return sdk.RepeatDirectionVertical, nil
	case "horizontal":
		return sdk.RepeatDirectionHorizontal, nil
	}

	return "", fmt.Errorf("invalid repeat direction '%s'", input)
}
