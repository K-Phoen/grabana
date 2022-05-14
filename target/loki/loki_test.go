package loki

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewLokiTargetCanBeCreated(t *testing.T) {
	req := require.New(t)
	query := "{app=\"loki\"}"

	target := New(query)

	req.Equal(query, target.Expr)
}

func TestLegendCanBeConfigured(t *testing.T) {
	req := require.New(t)
	legend := "lala"

	target := New("", Legend(legend))

	req.Equal(legend, target.LegendFormat)
}

func TestTargetCanBeHidden(t *testing.T) {
	req := require.New(t)

	target := New("", Hide())

	req.True(target.Hidden)
}
