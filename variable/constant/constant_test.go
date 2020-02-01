package constant

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConstantsCanBeCreated(t *testing.T) {
	req := require.New(t)

	panel := New("percentile")

	req.Equal("percentile", panel.Builder.Name)
	req.Equal("percentile", panel.Builder.Label)
	req.Equal("constant", panel.Builder.Type)
}

func TestLabelCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("const", Label("Constant"))

	req.Equal("const", panel.Builder.Name)
	req.Equal("Constant", panel.Builder.Label)
}

func TestValuesCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("const", Values(map[string]string{
		"90th": "90",
		"95th": "95",
		"99th": "99",
	}))

	labels := make([]string, 0, 3)
	values := make([]string, 0, 3)

	for _, opt := range panel.Builder.Options {
		labels = append(labels, opt.Text)
		values = append(values, opt.Value)
	}

	req.Len(values, 3)
	req.Equal([]string{"90th", "95th", "99th"}, labels)
	req.Equal([]string{"90", "95", "99"}, values)
}

func TestDefaultValueCanBeSet(t *testing.T) {
	req := require.New(t)

	panel := New("const", Default("99th"))

	req.Equal("99th", panel.Builder.Current.Text)
}

func TestLabelCanBeHidden(t *testing.T) {
	req := require.New(t)

	panel := New("", HideLabel())

	req.Equal(uint8(1), panel.Builder.Hide)
}

func TestVariableCanBeHidden(t *testing.T) {
	req := require.New(t)

	panel := New("custom", Hide())

	req.Equal(uint8(2), panel.Builder.Hide)
}
