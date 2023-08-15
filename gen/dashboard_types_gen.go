package dashboard

// Contains the list of annotations that are associated with the dashboard.
// Annotations are used to overlay event markers and overlay event tags on graphs.
// Grafana comes with a native annotation store and the ability to add annotation events directly from the graph panel or via the HTTP API.
// See https://grafana.com/docs/grafana/latest/dashboards/build-dashboards/annotate-visualizations/
type AnnotationContainer struct {
	List []AnnotationQuery `json:"list,omitempty"`
}

type AnnotationPanelFilter struct {
	Exclude bool    `json:"exclude,omitempty"`
	Ids     []int64 `json:"ids"`
}

// TODO docs
// FROM: AnnotationQuery in grafana-data/src/types/annotations.ts
type AnnotationQuery struct {
	Name       string                `json:"name"`
	Datasource DataSourceRef         `json:"datasource"`
	Enable     bool                  `json:"enable"`
	Hide       bool                  `json:"hide,omitempty"`
	IconColor  string                `json:"iconColor"`
	Filter     AnnotationPanelFilter `json:"filter,omitempty"`
	Target     AnnotationTarget      `json:"target,omitempty"`
	Type       string                `json:"type,omitempty"`
}

// TODO: this should be a regular DataQuery that depends on the selected dashboard
// these match the properties of the "grafana" datasource that is default in most dashboards
type AnnotationTarget struct {
	Limit    int64    `json:"limit"`
	MatchAny bool     `json:"matchAny"`
	Tags     []string `json:"tags"`
	Type     string   `json:"type"`
}

// 0 for no shared crosshair or tooltip (default).
// 1 for shared crosshair.
// 2 for shared crosshair AND shared tooltip.
type DashboardCursorSync int64

const (
	Off       DashboardCursorSync = 0
	Crosshair DashboardCursorSync = 1
	Tooltip   DashboardCursorSync = 2
)

// Links with references to other dashboards or external resources
type DashboardLink struct {
	Title       string            `json:"title"`
	Type        DashboardLinkType `json:"type"`
	Icon        string            `json:"icon"`
	Tooltip     string            `json:"tooltip"`
	Url         string            `json:"url"`
	Tags        []string          `json:"tags"`
	AsDropdown  bool              `json:"asDropdown"`
	TargetBlank bool              `json:"targetBlank"`
	IncludeVars bool              `json:"includeVars"`
	KeepTime    bool              `json:"keepTime"`
}

// Dashboard Link type. Accepted values are dashboards (to refer to another dashboard) and link (to refer to an external resource)
type DashboardLinkType string

const (
	Link       DashboardLinkType = "link"
	Dashboards DashboardLinkType = "dashboards"
)

type DashboardTemplating struct {
	List []VariableModel `json:"list,omitempty"`
}

// Ref to a DataSource instance
type DataSourceRef struct {
	Type string `json:"type,omitempty"`
	Uid  string `json:"uid,omitempty"`
}

// Transformations allow to manipulate data returned by a query before the system applies a visualization.
// Using transformations you can: rename fields, join time series data, perform mathematical operations across queries,
// use the output of one transformation as the input to another transformation, etc.
type DataTransformerConfig struct {
	Id       string        `json:"id"`
	Disabled bool          `json:"disabled,omitempty"`
	Filter   MatcherConfig `json:"filter,omitempty"`
	Options  any           `json:"options"`
}

// Map a field to a color.
type FieldColor struct {
	Mode       FieldColorModeId       `json:"mode"`
	FixedColor string                 `json:"fixedColor,omitempty"`
	SeriesBy   FieldColorSeriesByMode `json:"seriesBy,omitempty"`
}

// Color mode for a field. You can specify a single color, or select a continuous (gradient) color schemes, based on a value.
// Continuous color interpolates a color using the percentage of a value relative to min and max.
// Accepted values are:
// `thresholds`: From thresholds. Informs Grafana to take the color from the matching threshold
// `palette-classic`: Classic palette. Grafana will assign color by looking up a color in a palette by series index. Useful for Graphs and pie charts and other categorical data visualizations
// `palette-classic-by-name`: Classic palette (by name). Grafana will assign color by looking up a color in a palette by series name. Useful for Graphs and pie charts and other categorical data visualizations
// `continuous-GrYlRd`: ontinuous Green-Yellow-Red palette mode
// `continuous-RdYlGr`: Continuous Red-Yellow-Green palette mode
// `continuous-BlYlRd`: Continuous Blue-Yellow-Red palette mode
// `continuous-YlRd`: Continuous Yellow-Red palette mode
// `continuous-BlPu`: Continuous Blue-Purple palette mode
// `continuous-YlBl`: Continuous Yellow-Blue palette mode
// `continuous-blues`: Continuous Blue palette mode
// `continuous-reds`: Continuous Red palette mode
// `continuous-greens`: Continuous Green palette mode
// `continuous-purples`: Continuous Purple palette mode
// `shades`: Shades of a single color. Specify a single color, useful in an override rule.
// `fixed`: Fixed color mode. Specify a single color, useful in an override rule.
type FieldColorModeId string

const (
	Thresholds           FieldColorModeId = "thresholds"
	PaletteClassic       FieldColorModeId = "palette-classic"
	PaletteClassicByName FieldColorModeId = "palette-classic-by-name"
	ContinuousGrYlRd     FieldColorModeId = "continuous-GrYlRd"
	ContinuousRdYlGr     FieldColorModeId = "continuous-RdYlGr"
	ContinuousBlYlRd     FieldColorModeId = "continuous-BlYlRd"
	ContinuousYlRd       FieldColorModeId = "continuous-YlRd"
	ContinuousBlPu       FieldColorModeId = "continuous-BlPu"
	ContinuousYlBl       FieldColorModeId = "continuous-YlBl"
	ContinuousBlues      FieldColorModeId = "continuous-blues"
	ContinuousReds       FieldColorModeId = "continuous-reds"
	ContinuousGreens     FieldColorModeId = "continuous-greens"
	ContinuousPurples    FieldColorModeId = "continuous-purples"
	Fixed                FieldColorModeId = "fixed"
	Shades               FieldColorModeId = "shades"
)

// Defines how to assign a series color from "by value" color schemes. For example for an aggregated data points like a timeseries, the color can be assigned by the min, max or last value.
type FieldColorSeriesByMode string

const (
	Min  FieldColorSeriesByMode = "min"
	Max  FieldColorSeriesByMode = "max"
	Last FieldColorSeriesByMode = "last"
)

// The data model used in Grafana, namely the data frame, is a columnar-oriented table structure that unifies both time series and table query results.
// Each column within this structure is called a field. A field can represent a single time series or table column.
// Field options allow you to change how the data is displayed in your visualizations.
type FieldConfig struct {
	DisplayName       string           `json:"displayName,omitempty"`
	DisplayNameFromDS string           `json:"displayNameFromDS,omitempty"`
	Description       string           `json:"description,omitempty"`
	Path              string           `json:"path,omitempty"`
	Writeable         bool             `json:"writeable,omitempty"`
	Filterable        bool             `json:"filterable,omitempty"`
	Unit              string           `json:"unit,omitempty"`
	Decimals          float64          `json:"decimals,omitempty"`
	Min               float64          `json:"min,omitempty"`
	Max               float64          `json:"max,omitempty"`
	Mappings          []ValueMapping   `json:"mappings,omitempty"`
	Thresholds        ThresholdsConfig `json:"thresholds,omitempty"`
	Color             FieldColor       `json:"color,omitempty"`
	Links             []any            `json:"links,omitempty"`
	NoValue           string           `json:"noValue,omitempty"`
	Custom            any              `json:"custom,omitempty"`
}

// The data model used in Grafana, namely the data frame, is a columnar-oriented table structure that unifies both time series and table query results.
// Each column within this structure is called a field. A field can represent a single time series or table column.
// Field options allow you to change how the data is displayed in your visualizations.
type FieldConfigSource struct {
	Defaults  FieldConfig `json:"defaults"`
	Overrides []any       `json:"overrides"`
}

// Position and dimensions of a panel in the grid
type GridPos struct {
	H      int64 `json:"h"`
	W      int64 `json:"w"`
	X      int64 `json:"x"`
	Y      int64 `json:"y"`
	Static bool  `json:"static,omitempty"`
}

// A library panel is a reusable panel that you can use in any dashboard.
// When you make a change to a library panel, that change propagates to all instances of where the panel is used.
// Library panels streamline reuse of panels across multiple dashboards.
type LibraryPanelRef struct {
	Name string `json:"name"`
	Uid  string `json:"uid"`
}

// Loading status
// Accepted values are `NotStarted` (the request is not started), `Loading` (waiting for response), `Streaming` (pulling continuous data), `Done` (response received successfully) or `Error` (failed request).
type LoadingState string

const (
	NotStarted LoadingState = "NotStarted"
	Loading    LoadingState = "Loading"
	Streaming  LoadingState = "Streaming"
	Done       LoadingState = "Done"
	Error      LoadingState = "Error"
)

// Supported value mapping types
// `value`: Maps text values to a color or different display text and color. For example, you can configure a value mapping so that all instances of the value 10 appear as Perfection! rather than the number.
// `range`: Maps numerical ranges to a display text and color. For example, if a value is within a certain range, you can configure a range value mapping to display Low or High rather than the number.
// `regex`: Maps regular expressions to replacement text and a color. For example, if a value is www.example.com, you can configure a regex value mapping so that Grafana displays www and truncates the domain.
// `special`: Maps special values like Null, NaN (not a number), and boolean values like true and false to a display text and color. See SpecialValueMatch to see the list of special values. For example, you can configure a special value mapping so that null values appear as N/A.
type MappingType string

const (
	ValueToText  MappingType = "value"
	RangeToText  MappingType = "range"
	RegexToText  MappingType = "regex"
	SpecialValue MappingType = "special"
)

// Matcher is a predicate configuration. Based on the config a set of field(s) or values is filtered in order to apply override / transformation.
// It comes with in id ( to resolve implementation from registry) and a configuration that’s specific to a particular matcher type.
type MatcherConfig struct {
	Id      string `json:"id"`
	Options any    `json:"options,omitempty"`
}

// Dashboard panels are the basic visualization building blocks.
type Panel struct {
	Type            string                  `json:"type"`
	Id              int64                   `json:"id,omitempty"`
	PluginVersion   string                  `json:"pluginVersion,omitempty"`
	Tags            []string                `json:"tags,omitempty"`
	Targets         []Target                `json:"targets,omitempty"`
	Title           string                  `json:"title,omitempty"`
	Description     string                  `json:"description,omitempty"`
	Transparent     bool                    `json:"transparent"`
	Datasource      DataSourceRef           `json:"datasource,omitempty"`
	GridPos         GridPos                 `json:"gridPos,omitempty"`
	Links           []DashboardLink         `json:"links,omitempty"`
	Repeat          string                  `json:"repeat,omitempty"`
	RepeatDirection PanelRepeatDirection    `json:"repeatDirection,omitempty"`
	RepeatPanelId   int64                   `json:"repeatPanelId,omitempty"`
	MaxDataPoints   float64                 `json:"maxDataPoints,omitempty"`
	Transformations []DataTransformerConfig `json:"transformations"`
	Interval        string                  `json:"interval,omitempty"`
	TimeFrom        string                  `json:"timeFrom,omitempty"`
	TimeShift       string                  `json:"timeShift,omitempty"`
	LibraryPanel    LibraryPanelRef         `json:"libraryPanel,omitempty"`
	Options         any                     `json:"options"`
	FieldConfig     FieldConfigSource       `json:"fieldConfig"`
}

type PanelRepeatDirection string

const (
	Horizontal PanelRepeatDirection = "h"
	Vertical   PanelRepeatDirection = "v"
)

// Maps numerical ranges to a display text and color.
// For example, if a value is within a certain range, you can configure a range value mapping to display Low or High rather than the number.
type RangeMap struct {
	Options any `json:"options"`
}

// Maps regular expressions to replacement text and a color.
// For example, if a value is www.example.com, you can configure a regex value mapping so that Grafana displays www and truncates the domain.
type RegexMap struct {
	Options any `json:"options"`
}

// Row panel
type RowPanel struct {
	Type       string        `json:"type"`
	Collapsed  bool          `json:"collapsed"`
	Title      string        `json:"title,omitempty"`
	Datasource DataSourceRef `json:"datasource,omitempty"`
	GridPos    GridPos       `json:"gridPos,omitempty"`
	Id         int64         `json:"id"`
	Panels     []Panel       `json:"panels"`
	Repeat     string        `json:"repeat,omitempty"`
}

// A dashboard snapshot shares an interactive dashboard publicly.
// It is a read-only version of a dashboard, and is not editable.
// It is possible to create a snapshot of a snapshot.
// Grafana strips away all sensitive information from the dashboard.
// Sensitive information stripped: queries (metric, template,annotation) and panel links.
type Snapshot struct {
	Created     string `json:"created"`
	Expires     string `json:"expires"`
	External    bool   `json:"external"`
	ExternalUrl string `json:"externalUrl"`
	Id          int64  `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	OrgId       int64  `json:"orgId"`
	Updated     string `json:"updated"`
	Url         string `json:"url,omitempty"`
	UserId      int64  `json:"userId"`
}

// Maps special values like Null, NaN (not a number), and boolean values like true and false to a display text and color.
// See SpecialValueMatch to see the list of special values.
// For example, you can configure a special value mapping so that null values appear as N/A.
type SpecialValueMap struct {
	Options any `json:"options"`
}

// Special value types supported by the `SpecialValueMap`
type SpecialValueMatch string

const (
	True       SpecialValueMatch = "true"
	False      SpecialValueMatch = "false"
	Null       SpecialValueMatch = "null"
	NaN        SpecialValueMatch = "nan"
	NullAndNan SpecialValueMatch = "null+nan"
	Empty      SpecialValueMatch = "empty"
)

type Style string

const (
	Light Style = "light"
	Dark  Style = "dark"
)

// Schema for panel targets is specified by datasource
// plugins. We use a placeholder definition, which the Go
// schema loader either left open/as-is with the Base
// variant of the Dashboard and Panel families, or filled
// with types derived from plugins in the Instance variant.
// When working directly from CUE, importers can extend this
// type directly to achieve the same effect.
type Target struct {
}

// User-defined value for a metric that triggers visual changes in a panel when this value is met or exceeded
// They are used to conditionally style and color visualizations based on query results , and can be applied to most visualizations.
type Threshold struct {
	Value *float64 `json:"value"`
	Color string   `json:"color"`
}

// Thresholds configuration for the panel
type ThresholdsConfig struct {
	Mode  ThresholdsMode `json:"mode"`
	Steps []Threshold    `json:"steps"`
}

// Thresholds can either be `absolute` (specific number) or `percentage` (relative to min or max, it will be values between 0 and 1).
type ThresholdsMode string

const (
	Absolute   ThresholdsMode = "absolute"
	Percentage ThresholdsMode = "percentage"
)

type TimeInterval struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type TimePicker struct {
	Hidden            bool     `json:"hidden"`
	Refresh_intervals []string `json:"refresh_intervals"`
	Collapse          bool     `json:"collapse"`
	Enable            bool     `json:"enable"`
	Time_options      []string `json:"time_options"`
}

// Maps text values to a color or different display text and color.
// For example, you can configure a value mapping so that all instances of the value 10 appear as Perfection! rather than the number.
type ValueMap struct {
	Options any `json:"options"`
}

// Allow to transform the visual representation of specific data values in a visualization, irrespective of their original units
type ValueMapping struct {
}

// Result used as replacement with text and color when the value matches
type ValueMappingResult struct {
	Text  string `json:"text,omitempty"`
	Color string `json:"color,omitempty"`
	Icon  string `json:"icon,omitempty"`
	Index int64  `json:"index,omitempty"`
}

// Determine if the variable shows on dashboard
// Accepted values are 0 (show label and value), 1 (show value only), 2 (show nothing).
type VariableHide int64

const (
	DontHide     VariableHide = 0
	HideLabel    VariableHide = 1
	HideVariable VariableHide = 2
)

// A variable is a placeholder for a value. You can use variables in metric queries and in panel titles.
type VariableModel struct {
	Id          string           `json:"id"`
	Type        VariableType     `json:"type"`
	Name        string           `json:"name"`
	Label       string           `json:"label,omitempty"`
	Hide        VariableHide     `json:"hide"`
	SkipUrlSync bool             `json:"skipUrlSync"`
	Description string           `json:"description,omitempty"`
	Query       any              `json:"query,omitempty"`
	Datasource  DataSourceRef    `json:"datasource,omitempty"`
	AllFormat   string           `json:"allFormat,omitempty"`
	Current     VariableOption   `json:"current,omitempty"`
	Multi       bool             `json:"multi,omitempty"`
	Options     []VariableOption `json:"options,omitempty"`
	Refresh     VariableRefresh  `json:"refresh,omitempty"`
}

// Option to be selected in a variable.
type VariableOption struct {
	Selected bool          `json:"selected,omitempty"`
	Text     StringOrArray `json:"text"`
	Value    StringOrArray `json:"value"`
}

// Options to config when to refresh a variable
// `0`: Never refresh the variable
// `1`: Queries the data source every time the dashboard loads.
// `2`: Queries the data source when the dashboard time range changes.
type VariableRefresh int64

const (
	Never              VariableRefresh = 0
	OnDashboardLoad    VariableRefresh = 1
	OnTimeRangeChanged VariableRefresh = 2
)

// Sort variable options
// Accepted values are:
// `0`: No sorting
// `1`: Alphabetical ASC
// `2`: Alphabetical DESC
// `3`: Numerical ASC
// `4`: Numerical DESC
// `5`: Alphabetical Case Insensitive ASC
// `6`: Alphabetical Case Insensitive DESC
type VariableSort int64

const (
	Disabled                        VariableSort = 0
	AlphabeticalAsc                 VariableSort = 1
	AlphabeticalDesc                VariableSort = 2
	NumericalAsc                    VariableSort = 3
	NumericalDesc                   VariableSort = 4
	AlphabeticalCaseInsensitiveAsc  VariableSort = 5
	AlphabeticalCaseInsensitiveDesc VariableSort = 6
)

// Dashboard variable type
// `query`: Query-generated list of values such as metric names, server names, sensor IDs, data centers, and so on.
// `adhoc`: Key/value filters that are automatically added to all metric queries for a data source (Prometheus, Loki, InfluxDB, and Elasticsearch only).
// `constant`: 	Define a hidden constant.
// `datasource`: Quickly change the data source for an entire dashboard.
// `interval`: Interval variables represent time spans.
// `textbox`: Display a free text input field with an optional default value.
// `custom`: Define the variable options manually using a comma-separated list.
// `system`: Variables defined by Grafana. See: https://grafana.com/docs/grafana/latest/dashboards/variables/add-template-variables/#global-variables
type VariableType string

const (
	Query      VariableType = "query"
	Adhoc      VariableType = "adhoc"
	Constant   VariableType = "constant"
	Datasource VariableType = "datasource"
	Interval   VariableType = "interval"
	Textbox    VariableType = "textbox"
	Custom     VariableType = "custom"
	System     VariableType = "system"
)

// This is a dashboard.
type Dashboard struct {
	Id                   *int64              `json:"id,omitempty"`
	Uid                  string              `json:"uid,omitempty"`
	Title                string              `json:"title,omitempty"`
	Description          string              `json:"description,omitempty"`
	Revision             int64               `json:"revision,omitempty"`
	GnetId               string              `json:"gnetId,omitempty"`
	Tags                 []string            `json:"tags,omitempty"`
	Style                Style               `json:"style"`
	Timezone             string              `json:"timezone,omitempty"`
	Editable             bool                `json:"editable"`
	GraphTooltip         DashboardCursorSync `json:"graphTooltip"`
	Time                 TimeInterval        `json:"time,omitempty"`
	Timepicker           TimePicker          `json:"timepicker,omitempty"`
	FiscalYearStartMonth int64               `json:"fiscalYearStartMonth,omitempty"`
	LiveNow              bool                `json:"liveNow,omitempty"`
	WeekStart            string              `json:"weekStart,omitempty"`
	Refresh              StringOrBool        `json:"refresh,omitempty"`
	SchemaVersion        int64               `json:"schemaVersion"`
	Version              int64               `json:"version,omitempty"`
	Panels               []RowPanel          `json:"panels,omitempty"`
	Templating           DashboardTemplating `json:"templating,omitempty"`
	Annotations          AnnotationContainer `json:"annotations,omitempty"`
	Links                []DashboardLink     `json:"links,omitempty"`
	Snapshot             Snapshot            `json:"snapshot,omitempty"`
}

type StringOrArray struct {
	ValString *string  `json:"ValString,omitempty"`
	ValArray  []string `json:"ValArray,omitempty"`
}

type StringOrBool struct {
	ValString *string `json:"ValString,omitempty"`
	ValBool   *bool   `json:"ValBool,omitempty"`
}
