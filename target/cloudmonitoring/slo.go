package cloudmonitoring

import "github.com/K-Phoen/sdk"

// SLOOption represents an option that can be used to configure an SLO query.
type SLOOption func(*SLO)

func SLOPerSeriesAligner(aligner Aligner, alignmentPeriod string) SLOOption {
	return func(s *SLO) {
		s.target.SLOQuery.PerSeriesAligner = string(aligner)
		s.target.SLOQuery.AlignmentPeriod = alignmentPeriod
	}
}

func SLOAliasBy(alias string) SLOOption {
	return func(s *SLO) {
		s.target.SLOQuery.AliasBy = alias
	}
}

func SelectorName(name string) SLOOption {
	return func(s *SLO) {
		s.target.SLOQuery.SelectorName = name
	}
}

func ServiceRef(id, name string) SLOOption {
	return func(s *SLO) {
		s.target.SLOQuery.ServiceID = id
		s.target.SLOQuery.ServiceName = name
	}
}

func SLORef(id, name string) SLOOption {
	return func(s *SLO) {
		s.target.SLOQuery.SLOID = id
		s.target.SLOQuery.SLOName = name
	}
}

func LookbackPeriod(period string) SLOOption {
	return func(s *SLO) {
		s.target.SLOQuery.LookbackPeriod = period
	}
}

type SLO struct {
	target *sdk.Target
}

func NewSLO(projectName string, options ...SLOOption) *SLO {
	slo := &SLO{
		target: &sdk.Target{
			QueryType: "slo",
			SLOQuery: &sdk.StackdriverSLOQuery{
				ProjectName: projectName,
			},
		},
	}

	for _, opt := range options {
		opt(slo)
	}

	return slo
}

func (s *SLO) Target() *sdk.Target { return s.target }
func (s *SLO) AlertModel() sdk.AlertModel {
	return sdk.AlertModel{
		QueryType: s.target.QueryType,
		SLOQuery:  s.target.SLOQuery,
	}
}
