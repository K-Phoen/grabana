package fields

import (
	"testing"

	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
)

func TestByName(t *testing.T) {
	req := require.New(t)

	overrideCfg := &sdk.FieldConfigOverride{}
	ByName("some-name")(overrideCfg)

	req.Equal("byName", overrideCfg.Matcher.ID)
	req.Equal("some-name", overrideCfg.Matcher.Options)
}

func TestByQuery(t *testing.T) {
	req := require.New(t)

	overrideCfg := &sdk.FieldConfigOverride{}
	ByQuery("A")(overrideCfg)

	req.Equal("byFrameRefID", overrideCfg.Matcher.ID)
	req.Equal("A", overrideCfg.Matcher.Options)
}
