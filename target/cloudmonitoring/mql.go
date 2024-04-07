package cloudmonitoring

import "github.com/K-Phoen/sdk"

const (
	GraphPeriodDisabled = "disabled"
	GraphPeriodAuto     = "auto"
)

type MQLOption func(*MQL)

func MQLAliasBy(alias string) MQLOption {
	return func(m *MQL) {
		m.target.AliasBy = alias
	}
}

func GraphPeriod(graphPeriod string) MQLOption {
	return func(m *MQL) {
		m.target.TimeSeriesQuery.GraphPeriod = graphPeriod
	}
}

type MQL struct {
	target *sdk.Target
}

// NewMQL returns a target builder making an MQL query.
func NewMQL(projectName, query string, options ...MQLOption) *MQL {
	mql := &MQL{
		target: &sdk.Target{
			QueryType: "timeSeriesQuery",
			TimeSeriesQuery: &sdk.StackdriverTimeSeriesQuery{
				ProjectName: projectName,
				Query:       query,
			},
		},
	}

	for _, opt := range options {
		opt(mql)
	}

	return mql
}

func (m *MQL) Target() *sdk.Target { return m.target }
