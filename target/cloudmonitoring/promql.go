package cloudmonitoring

import (
	"github.com/K-Phoen/sdk"
)

// PromQLOption represents an option that can be used to configure a promQL query.
type PromQLOption func(*PromQL)

func MinStep(step string) PromQLOption {
	return func(q *PromQL) {
		q.target.PromQLQuery.Step = step
	}
}

// PromQL represents a google cloud monitoring query.
type PromQL struct {
	target *sdk.Target
}

// NewPromQL returns a target builder making a PromQL query.
func NewPromQL(projectName, expr string, options ...PromQLOption) *PromQL {
	promQL := &PromQL{
		target: &sdk.Target{
			QueryType: "promQL",
			// For some reason I can't explain, Grafana seems to require TimeSeriesQuery to be set
			// when we're making a promQL query.
			TimeSeriesQuery: &sdk.StackdriverTimeSeriesQuery{
				ProjectName: projectName,
			},
			PromQLQuery: &sdk.StackdriverPromQLQuery{
				ProjectName: projectName,
				Expr:        expr,
				Step:        "10s",
			},
		},
	}

	for _, opt := range options {
		opt(promQL)
	}

	return promQL
}

func (p *PromQL) Target() *sdk.Target { return p.target }

func (p *PromQL) AlertModel() sdk.AlertModel {
	return sdk.AlertModel{
		QueryType:   p.target.QueryType,
		PromQLQuery: p.target.PromQLQuery,
	}
}
