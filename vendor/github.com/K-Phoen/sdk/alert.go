package sdk

/*
   Copyright 2016 Alexander I.Grafov <grafov@gmail.com>
   Copyright 2016-2019 The Grafana SDK authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

	   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

   ॐ तारे तुत्तारे तुरे स्व
*/

type Alert struct {
	Name     string      `json:"name"`
	Interval string      `json:"interval"`
	Rules    []AlertRule `json:"rules"`
}

type AlertRule struct {
	For          string            `json:"for"`
	GrafanaAlert *GrafanaAlert     `json:"grafana_alert,omitempty"`
	Annotations  map[string]string `json:"annotations,omitempty"`
	Labels       map[string]string `json:"labels,omitempty"`
}

type GrafanaAlert struct {
	Title               string       `json:"title"`
	Condition           string       `json:"condition"`
	NoDataState         string       `json:"no_data_state"`
	ExecutionErrorState string       `json:"exec_err_state,omitempty"`
	Data                []AlertQuery `json:"data"`
}

type AlertQuery struct {
	RefID             string                  `json:"refId"`
	QueryType         string                  `json:"queryType"`
	RelativeTimeRange *AlertRelativeTimeRange `json:"relativeTimeRange,omitempty"`
	DatasourceUID     string                  `json:"datasourceUid"`
	Model             AlertModel              `json:"model"`
}

type AlertRelativeTimeRange struct {
	From int `json:"from"` // seconds
	To   int `json:"to"`   // seconds
}

type AlertModel struct {
	RefID        string             `json:"refId,omitempty"`
	QueryType    string             `json:"queryType,omitempty"`
	Type         string             `json:"type,omitempty"`
	Expr         string             `json:"expr,omitempty"`
	Format       string             `json:"format,omitempty"`
	LegendFormat string             `json:"legendFormat,omitempty"`
	Datasource   AlertDatasourceRef `json:"datasource"`
	Interval     string             `json:"interval,omitempty"`
	IntervalMs   int                `json:"intervalMs,omitempty"`
	Hide         *bool              `json:"hide,omitempty"`
	Conditions   []AlertCondition   `json:"conditions,omitempty"`

	// For Graphite
	Target string `json:"target,omitempty"`

	// For Stackdriver
	MetricQuery *StackdriverAlertQuery `json:"metricQuery,omitempty"`
}

type StackdriverAlertQuery struct {
	AlignOptions       []StackdriverAlignOptions `json:"alignOptions,omitempty"`
	AliasBy            string                    `json:"aliasBy,omitempty"`
	MetricType         string                    `json:"metricType,omitempty"`
	MetricKind         string                    `json:"metricKind,omitempty"`
	Filters            []string                  `json:"filters,omitempty"`
	AlignmentPeriod    string                    `json:"alignmentPeriod,omitempty"`
	CrossSeriesReducer string                    `json:"crossSeriesReducer,omitempty"`
	PerSeriesAligner   string                    `json:"perSeriesAligner,omitempty"`
	ValueType          string                    `json:"valueType,omitempty"`
	Preprocessor       string                    `json:"preprocessor,omitempty"`
	GroupBys           []string                  `json:"groupBys,omitempty"`
}

type AlertDatasourceRef struct {
	UID  string `json:"uid"`
	Type string `json:"type"`
}

type AlertCondition struct {
	Type      string                 `json:"type,omitempty"`
	Evaluator AlertEvaluator         `json:"evaluator,omitempty"`
	Operator  AlertOperator          `json:"operator,omitempty"`
	Query     AlertConditionQueryRef `json:"query,omitempty"`
	Reducer   AlertReducer           `json:"reducer,omitempty"`
}
type AlertConditionQueryRef struct {
	Params []string `json:"params,omitempty"`
}
type AlertEvaluator struct {
	Params []float64 `json:"params,omitempty"`
	Type   string    `json:"type,omitempty"`
}

type AlertOperator struct {
	Type string `json:"type,omitempty"`
}

type AlertReducer struct {
	Params []string `json:"params,omitempty"`
	Type   string   `json:"type,omitempty"`
}
