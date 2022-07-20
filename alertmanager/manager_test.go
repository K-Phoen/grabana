package alertmanager

import (
	"testing"

	"github.com/K-Phoen/sdk"
	"github.com/stretchr/testify/require"
)

func TestDefaultContactPoint(t *testing.T) {
	req := require.New(t)

	manager := New(
		ContactPoints(
			ContactPoint("team-b"),
			ContactPoint("team-a"),
		),
		DefaultContactPoint("team-a"),
	)

	req.Equal("team-a", manager.builder.Config.Route.Receiver)
}

func TestDefaultGroupBy(t *testing.T) {
	req := require.New(t)

	manager := New(
		DefaultGroupBys("priority", "service"),
	)

	req.ElementsMatch([]string{"priority", "service"}, manager.builder.Config.Route.GroupBy)
}

func TestDefaultContactPointCanBeImplicit(t *testing.T) {
	req := require.New(t)

	manager := New(
		ContactPoints(
			ContactPoint("team-b"),
			ContactPoint("team-a"),
		),
	)

	req.Equal("team-b", manager.builder.Config.Route.Receiver)
}

func TestContactPoints(t *testing.T) {
	req := require.New(t)

	manager := New(
		ContactPoints(
			ContactPoint("team-b"),
			ContactPoint("team-a"),
		),
	)

	req.Len(manager.builder.Config.Receivers, 2)
}

func TestRouting(t *testing.T) {
	req := require.New(t)

	manager := New(
		ContactPoints(
			ContactPoint("team-a"),
		),
		Routing(
			Policy("team-a", TagEq("owner", "team-a")),
		),
	)

	req.Len(manager.builder.Config.Route.Routes, 1)
}

func TestTemplates(t *testing.T) {
	req := require.New(t)

	manager := New(
		Templates(map[string]string{
			"custom_template": "lala",
		}),
	)

	req.Equal(sdk.MessageTemplate{
		"custom_template": "lala",
	}, manager.builder.TemplateFiles)
}

func TestMarshalJSON(t *testing.T) {
	req := require.New(t)

	manager := New(
		ContactPoints(
			ContactPoint("team-a"),
		),
		Routing(
			Policy("team-a", TagEq("owner", "team-a")),
		),
	)

	_, errJSON := manager.MarshalJSON()
	_, errJSONIndent := manager.MarshalIndentJSON()

	req.NoError(errJSON)
	req.NoError(errJSONIndent)
}
