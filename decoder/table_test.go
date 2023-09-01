package decoder

import (
	"testing"

	"github.com/K-Phoen/grabana/target/cloudwatch"
	"github.com/stretchr/testify/assert"
)

func TestTableCloudwatchTarget(t *testing.T) {
	panel := DashboardTable{
		Title: "Table target test",
		Targets: []Target{
			{
				Cloudwatch: &CloudwatchTarget{
					QueryParams: cloudwatch.QueryParams{
						Dimensions: map[string]string{
							"Name": "test",
						},
					},
				},
			},
		},
	}

	option, err := panel.target(panel.Targets[0])
	assert.NoError(t, err)
	assert.NotNil(t, option)
}
