package cloudmonitoring

import "github.com/K-Phoen/sdk"

// TimeSeriesOption represents an option that can be used to configure a timeseries query.
type TimeSeriesOption func(target *TimeSeries)

func CrossSeriesReducer(reducer Reducer) TimeSeriesOption {
	return func(target *TimeSeries) {
		target.target.TimeSeriesList.CrossSeriesReducer = string(reducer)
	}
}

func PerSeriesAligner(aligner Aligner, alignmentPeriod string) TimeSeriesOption {
	return func(target *TimeSeries) {
		target.target.TimeSeriesList.PerSeriesAligner = string(aligner)
		target.target.TimeSeriesList.AlignmentPeriod = alignmentPeriod
	}
}

func GroupBy(field string) TimeSeriesOption {
	return func(target *TimeSeries) {
		target.target.TimeSeriesList.GroupBys = append(
			target.target.TimeSeriesList.GroupBys,
			field,
		)
	}
}

func Filter(field string, op FilterOperator, value string) TimeSeriesOption {
	return func(target *TimeSeries) {
		target.target.TimeSeriesList.Filters = append(
			target.target.TimeSeriesList.Filters,
			"AND", field, string(op), value,
		)
	}
}

func View(view string) TimeSeriesOption {
	return func(target *TimeSeries) {
		target.target.TimeSeriesList.View = view
	}
}

func Title(title string) TimeSeriesOption {
	return func(target *TimeSeries) {
		target.target.TimeSeriesList.Title = title
	}
}

func SecondaryCrossSeriesReducer(reducer Reducer) TimeSeriesOption {
	return func(target *TimeSeries) {
		target.target.TimeSeriesList.SecondaryCrossSeriesReducer = string(reducer)
	}
}

func SecondaryPerSeriesAligner(aligner Aligner, alignmentPeriod string) TimeSeriesOption {
	return func(target *TimeSeries) {
		target.target.TimeSeriesList.SecondaryPerSeriesAligner = string(aligner)
		target.target.TimeSeriesList.SecondaryAlignmentPeriod = alignmentPeriod
	}
}

func SecondaryGroupBy(field string) TimeSeriesOption {
	return func(target *TimeSeries) {
		target.target.TimeSeriesList.SecondaryGroupBys = append(
			target.target.TimeSeriesList.SecondaryGroupBys,
			field,
		)
	}
}

func Preprocessor(preprocessor PreprocessorMethod) TimeSeriesOption {
	return func(target *TimeSeries) {
		target.target.TimeSeriesList.Preprocessor = string(preprocessor)
	}
}

func TimeSeriesAliasBy(alias string) TimeSeriesOption {
	return func(target *TimeSeries) {
		target.target.AliasBy = alias
	}
}

// TimeSeries represents a google cloud monitoring query.
type TimeSeries struct {
	target *sdk.Target
}

// NewTimeSeries returns a target builder making a classic timeseries query.
func NewTimeSeries(projectName, metricType string, options ...TimeSeriesOption) *TimeSeries {
	cloudMonitoring := &TimeSeries{
		target: &sdk.Target{
			QueryType: "timeSeriesList",
			TimeSeriesList: &sdk.GCMTimeSeriesList{
				ProjectName: projectName,
				Filters: []string{
					"metric.type", string(FilterOperatorEqual), metricType,
				},
			},
		},
	}

	for _, opt := range options {
		opt(cloudMonitoring)
	}

	return cloudMonitoring
}

// Target implements the Target interface
func (t *TimeSeries) Target() *sdk.Target { return t.target }

// AlertModel implements the AlertModel interface
func (t *TimeSeries) AlertModel() sdk.AlertModel {
	return sdk.AlertModel{
		QueryType:      t.target.QueryType,
		TimeSeriesList: t.target.TimeSeriesList,
	}
}
