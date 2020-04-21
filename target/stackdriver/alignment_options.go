package stackdriver

import (
	"github.com/grafana-tools/sdk"
)

func defaultAlignmentOpts() sdk.StackdriverAlignOptions {
	return sdk.StackdriverAlignOptions{
		Expanded: true,
		Label:    "Alignment options",
		Options: []sdk.StackdriverAlignOption{
			{
				Label:       "delta",
				MetricKinds: []string{"CUMULATIVE", "DELTA"},
				Text:        "delta",
				Value:       "ALIGN_DELTA",
				ValueTypes:  []string{"INT64", "DOUBLE", "MONEY", "DISTRIBUTION"},
			},
			{
				Label:       "rate",
				MetricKinds: []string{"CUMULATIVE", "DELTA"},
				Text:        "rate",
				Value:       "ALIGN_RATE",
				ValueTypes:  []string{"INT64", "DOUBLE", "MONEY"},
			},
			{
				Label:       "min",
				MetricKinds: []string{"GAUGE", "DELTA"},
				Text:        "min",
				Value:       "ALIGN_MIN",
				ValueTypes:  []string{"INT64", "DOUBLE", "MONEY"},
			},
			{
				Label:       "max",
				MetricKinds: []string{"GAUGE", "DELTA"},
				Text:        "max",
				Value:       "ALIGN_MAX",
				ValueTypes:  []string{"INT64", "DOUBLE", "MONEY"},
			},
			{
				Label:       "mean",
				MetricKinds: []string{"GAUGE", "DELTA"},
				Text:        "mean",
				Value:       "ALIGN_MEAN",
				ValueTypes:  []string{"INT64", "DOUBLE", "MONEY"},
			},
			{
				Label:       "count",
				MetricKinds: []string{"GAUGE", "DELTA"},
				Text:        "count",
				Value:       "ALIGN_COUNT",
				ValueTypes:  []string{"INT64", "DOUBLE", "MONEY", "BOOL"},
			},
			{
				Label:       "sum",
				MetricKinds: []string{"GAUGE", "DELTA"},
				Text:        "sum",
				Value:       "ALIGN_SUM",
				ValueTypes:  []string{"INT64", "DOUBLE", "MONEY", "DISTRIBUTION"},
			},
			{
				Label:       "stddev",
				MetricKinds: []string{"GAUGE", "DELTA"},
				Text:        "stddev",
				Value:       "ALIGN_STDDEV",
				ValueTypes:  []string{"INT64", "DOUBLE", "MONEY"},
			},
			{
				Label:       "percent change",
				MetricKinds: []string{"GAUGE", "DELTA"},
				Text:        "percent change",
				Value:       "ALIGN_PERCENT_CHANGE",
				ValueTypes:  []string{"INT64", "DOUBLE", "MONEY"},
			},
		},
	}
}
