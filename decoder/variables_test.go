package decoder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidHideValues(t *testing.T) {
	validValues := []string{"", "label", "variable"}

	for _, hideValue := range validValues {
		tc := hideValue

		t.Run(tc, func(t *testing.T) {
			req := require.New(t)

			interval := VariableInterval{Hide: tc}
			_, err := interval.toOption()
			req.NoError(err)

			custom := VariableCustom{Hide: tc}
			_, err = custom.toOption()
			req.NoError(err)

			query := VariableQuery{Hide: tc}
			_, err = query.toOption()
			req.NoError(err)

			constant := VariableConst{Hide: tc}
			_, err = constant.toOption()
			req.NoError(err)

			datasource := VariableDatasource{Hide: tc}
			_, err = datasource.toOption()
			req.NoError(err)
		})
	}
}

func TestInvalidHideValue(t *testing.T) {
	req := require.New(t)
	invalidValue := "invalid"

	interval := VariableInterval{Hide: invalidValue}
	_, err := interval.toOption()
	req.Equal(ErrInvalidHideValue, err)

	custom := VariableCustom{Hide: invalidValue}
	_, err = custom.toOption()
	req.Equal(ErrInvalidHideValue, err)

	query := VariableQuery{Hide: invalidValue}
	_, err = query.toOption()
	req.Equal(ErrInvalidHideValue, err)

	constant := VariableConst{Hide: invalidValue}
	_, err = constant.toOption()
	req.Equal(ErrInvalidHideValue, err)

	datasource := VariableDatasource{Hide: invalidValue}
	_, err = datasource.toOption()
	req.Equal(ErrInvalidHideValue, err)
}
