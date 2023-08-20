package types

// Contains the list of annotations that are associated with the dashboard.
// Annotations are used to overlay event markers and overlay event tags on graphs.
// Grafana comes with a native annotation store and the ability to add annotation events directly from the graph panel or via the HTTP API.
// See https://grafana.com/docs/grafana/latest/dashboards/build-dashboards/annotate-visualizations/
type AnnotationContainer struct {
	// List of annotations
	List []AnnotationQuery `json:"list,omitempty"`
}

type AnnotationPanelFilter struct {
	// Should the specified panels be included or excluded
	Exclude *bool `json:"exclude,omitempty"`
	// Panel IDs that should be included or excluded
	Ids []uint8 `json:"ids"`
}

// TODO docs
// FROM: AnnotationQuery in grafana-data/src/types/annotations.ts
type AnnotationQuery struct {
	// Name of annotation.
	Name string `json:"name"`
	// Datasource where the annotations data is
	Datasource DataSourceRef `json:"datasource"`
	// When enabled the annotation query is issued with every dashboard refresh
	Enable bool `json:"enable"`
	// Annotation queries can be toggled on or off at the top of the dashboard.
	// When hide is true, the toggle is not shown in the dashboard.
	Hide *bool `json:"hide,omitempty"`
	// Color to use for the annotation event markers
	IconColor string `json:"iconColor"`
	// Filters to apply when fetching annotations
	Filter *AnnotationPanelFilter `json:"filter,omitempty"`
	// TODO.. this should just be a normal query target
	Target *AnnotationTarget `json:"target,omitempty"`
	// TODO -- this should not exist here, it is based on the --grafana-- datasource
	Type *string `json:"type,omitempty"`
}

// TODO: this should be a regular DataQuery that depends on the selected dashboard
// these match the properties of the "grafana" datasource that is default in most dashboards
type AnnotationTarget struct {
	// Only required/valid for the grafana datasource...
	// but code+tests is already depending on it so hard to change
	Limit int64 `json:"limit"`
	// Only required/valid for the grafana datasource...
	// but code+tests is already depending on it so hard to change
	MatchAny bool `json:"matchAny"`
	// Only required/valid for the grafana datasource...
	// but code+tests is already depending on it so hard to change
	Tags []string `json:"tags"`
	// Only required/valid for the grafana datasource...
	// but code+tests is already depending on it so hard to change
	Type string `json:"type"`
}

// This is a dashboard.
type Dashboard struct {
	// Unique numeric identifier for the dashboard.
	// `id` is internal to a specific Grafana instance. `uid` should be used to identify a dashboard across Grafana instances.
	Id *int64 `json:"id,omitempty"`
	// Unique dashboard identifier that can be generated by anyone. string (8-40)
	Uid *string `json:"uid,omitempty"`
	// Title of dashboard.
	Title *string `json:"title,omitempty"`
	// Description of dashboard.
	Description *string `json:"description,omitempty"`
	// This property should only be used in dashboards defined by plugins.  It is a quick check
	// to see if the version has changed since the last time.
	Revision *int64 `json:"revision,omitempty"`
	// ID of a dashboard imported from the https://grafana.com/grafana/dashboards/ portal
	GnetId *string `json:"gnetId,omitempty"`
	// Tags associated with dashboard.
	Tags []string `json:"tags,omitempty"`
	// Theme of dashboard.
	// Default value: dark.
	Style DashboardStyle `json:"style"`
	// Timezone of dashboard. Accepted values are IANA TZDB zone ID or "browser" or "utc".
	Timezone *string `json:"timezone,omitempty"`
	// Whether a dashboard is editable or not.
	Editable bool `json:"editable"`
	// Configuration of dashboard cursor sync behavior.
	// Accepted values are 0 (sync turned off), 1 (shared crosshair), 2 (shared crosshair and tooltip).
	GraphTooltip DashboardCursorSync `json:"graphTooltip"`
	// Time range for dashboard.
	// Accepted values are relative time strings like {from: 'now-6h', to: 'now'} or absolute time strings like {from: '2020-07-10T08:00:00.000Z', to: '2020-07-10T14:00:00.000Z'}.
	Time struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"time,omitempty"`
	// Configuration of the time picker shown at the top of a dashboard.
	Timepicker *TimePicker `json:"timepicker,omitempty"`
	// The month that the fiscal year starts on.  0 = January, 11 = December
	FiscalYearStartMonth *uint8 `json:"fiscalYearStartMonth,omitempty"`
	// When set to true, the dashboard will redraw panels at an interval matching the pixel width.
	// This will keep data "moving left" regardless of the query refresh rate. This setting helps
	// avoid dashboards presenting stale live data
	LiveNow *bool `json:"liveNow,omitempty"`
	// Day when the week starts. Expressed by the name of the day in lowercase, e.g. "monday".
	WeekStart *string `json:"weekStart,omitempty"`
	// Refresh rate of dashboard. Represented via interval string, e.g. "5s", "1m", "1h", "1d".
	Refresh *StringOrBool `json:"refresh,omitempty"`
	// Version of the JSON schema, incremented each time a Grafana update brings
	// changes to said schema.
	SchemaVersion uint16 `json:"schemaVersion"`
	// Version of the dashboard, incremented each time the dashboard is updated.
	Version *uint32 `json:"version,omitempty"`
	// List of dashboard panels
	Panels []RowPanel `json:"panels,omitempty"`
	// Configured template variables
	Templating *DashboardTemplating `json:"templating,omitempty"`
	// Contains the list of annotations that are associated with the dashboard.
	// Annotations are used to overlay event markers and overlay event tags on graphs.
	// Grafana comes with a native annotation store and the ability to add annotation events directly from the graph panel or via the HTTP API.
	// See https://grafana.com/docs/grafana/latest/dashboards/build-dashboards/annotate-visualizations/
	Annotations *AnnotationContainer `json:"annotations,omitempty"`
	// Links with references to other dashboards or external websites.
	Links []DashboardLink `json:"links,omitempty"`
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
	// Title to display with the link
	Title string `json:"title"`
	// Link type. Accepted values are dashboards (to refer to another dashboard) and link (to refer to an external resource)
	Type DashboardLinkType `json:"type"`
	// Icon name to be displayed with the link
	Icon string `json:"icon"`
	// Tooltip to display when the user hovers their mouse over it
	Tooltip string `json:"tooltip"`
	// Link URL. Only required/valid if the type is link
	Url string `json:"url"`
	// List of tags to limit the linked dashboards. If empty, all dashboards will be displayed. Only valid if the type is dashboards
	Tags []string `json:"tags"`
	// If true, all dashboards links will be displayed in a dropdown. If false, all dashboards links will be displayed side by side. Only valid if the type is dashboards
	AsDropdown bool `json:"asDropdown"`
	// If true, the link will be opened in a new tab
	TargetBlank bool `json:"targetBlank"`
	// If true, includes current template variables values in the link as query params
	IncludeVars bool `json:"includeVars"`
	// If true, includes current time range in the link as query params
	KeepTime bool `json:"keepTime"`
}

// Dashboard Link type. Accepted values are dashboards (to refer to another dashboard) and link (to refer to an external resource)
type DashboardLinkType string

const (
	Link       DashboardLinkType = "link"
	Dashboards DashboardLinkType = "dashboards"
)

type DashboardStyle string

const (
	Light DashboardStyle = "light"
	Dark  DashboardStyle = "dark"
)

type DashboardTemplating struct {
	// List of configured template variables with their saved values along with some other metadata
	List []VariableModel `json:"list,omitempty"`
}

// Ref to a DataSource instance
type DataSourceRef struct {
	// The plugin type-id
	Type *string `json:"type,omitempty"`
	// Specific datasource instance
	Uid *string `json:"uid,omitempty"`
}

// Transformations allow to manipulate data returned by a query before the system applies a visualization.
// Using transformations you can: rename fields, join time series data, perform mathematical operations across queries,
// use the output of one transformation as the input to another transformation, etc.
type DataTransformerConfig struct {
	// Unique identifier of transformer
	Id string `json:"id"`
	// Disabled transformations are skipped
	Disabled *bool `json:"disabled,omitempty"`
	// Optional frame matcher. When missing it will be applied to all results
	Filter *MatcherConfig `json:"filter,omitempty"`
	// Options to be passed to the transformer
	// Valid options depend on the transformer id
	Options any `json:"options"`
}

type DynamicConfigValue struct {
	Id    string `json:"id"`
	Value any    `json:"value,omitempty"`
}

// Map a field to a color.
type FieldColor struct {
	// The main color scheme mode.
	Mode FieldColorModeId `json:"mode"`
	// The fixed color value for fixed or shades color modes.
	FixedColor *string `json:"fixedColor,omitempty"`
	// Some visualizations need to know how to assign a series color from by value color schemes.
	SeriesBy *FieldColorSeriesByMode `json:"seriesBy,omitempty"`
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
	// The display value for this field.  This supports template variables blank is auto
	DisplayName *string `json:"displayName,omitempty"`
	// This can be used by data sources that return and explicit naming structure for values and labels
	// When this property is configured, this value is used rather than the default naming strategy.
	DisplayNameFromDS *string `json:"displayNameFromDS,omitempty"`
	// Human readable field metadata
	Description *string `json:"description,omitempty"`
	// An explicit path to the field in the datasource.  When the frame meta includes a path,
	// This will default to `${frame.meta.path}/${field.name}
	//
	// When defined, this value can be used as an identifier within the datasource scope, and
	// may be used to update the results
	Path *string `json:"path,omitempty"`
	// True if data source can write a value to the path. Auth/authz are supported separately
	Writeable *bool `json:"writeable,omitempty"`
	// True if data source field supports ad-hoc filters
	Filterable *bool `json:"filterable,omitempty"`
	// Unit a field should use. The unit you select is applied to all fields except time.
	// You can use the units ID availables in Grafana or a custom unit.
	// Available units in Grafana: https://github.com/grafana/grafana/blob/main/packages/grafana-data/src/valueFormats/categories.ts
	// As custom unit, you can use the following formats:
	// `suffix:<suffix>` for custom unit that should go after value.
	// `prefix:<prefix>` for custom unit that should go before value.
	// `time:<format>` For custom date time formats type for example `time:YYYY-MM-DD`.
	// `si:<base scale><unit characters>` for custom SI units. For example: `si: mF`. This one is a bit more advanced as you can specify both a unit and the source data scale. So if your source data is represented as milli (thousands of) something prefix the unit with that SI scale character.
	// `count:<unit>` for a custom count unit.
	// `currency:<unit>` for custom a currency unit.
	Unit *string `json:"unit,omitempty"`
	// Specify the number of decimals Grafana includes in the rendered value.
	// If you leave this field blank, Grafana automatically truncates the number of decimals based on the value.
	// For example 1.1234 will display as 1.12 and 100.456 will display as 100.
	// To display all decimals, set the unit to `String`.
	Decimals *float64 `json:"decimals,omitempty"`
	// The minimum value used in percentage threshold calculations. Leave blank for auto calculation based on all series and fields.
	Min *float64 `json:"min,omitempty"`
	// The maximum value used in percentage threshold calculations. Leave blank for auto calculation based on all series and fields.
	Max *float64 `json:"max,omitempty"`
	// Convert input values into a display string
	Mappings []ValueMapOrRangeMapOrRegexMapOrSpecialValueMap `json:"mappings,omitempty"`
	// Map numeric values to states
	Thresholds *ThresholdsConfig `json:"thresholds,omitempty"`
	// Panel color configuration
	Color *FieldColor `json:"color,omitempty"`
	// The behavior when clicking on a result
	Links []any `json:"links,omitempty"`
	// Alternative to empty string
	NoValue *string `json:"noValue,omitempty"`
	// custom is specified by the FieldConfig field
	// in panel plugin schemas.
	Custom any `json:"custom,omitempty"`
}

// The data model used in Grafana, namely the data frame, is a columnar-oriented table structure that unifies both time series and table query results.
// Each column within this structure is called a field. A field can represent a single time series or table column.
// Field options allow you to change how the data is displayed in your visualizations.
type FieldConfigSource struct {
	// Defaults are the options applied to all fields.
	Defaults FieldConfig `json:"defaults"`
	// Overrides are the options applied to specific fields overriding the defaults.
	Overrides []FieldConfigSourceOverride `json:"overrides"`
}

type FieldConfigSourceOverride struct {
	Matcher    MatcherConfig        `json:"matcher"`
	Properties []DynamicConfigValue `json:"properties"`
}

// Position and dimensions of a panel in the grid
type GridPos struct {
	// Panel height. The height is the number of rows from the top edge of the panel.
	H uint32 `json:"h"`
	// Panel width. The width is the number of columns from the left edge of the panel.
	W uint32 `json:"w"`
	// Panel x. The x coordinate is the number of columns from the left edge of the grid
	X uint32 `json:"x"`
	// Panel y. The y coordinate is the number of rows from the top edge of the grid
	Y uint32 `json:"y"`
	// Whether the panel is fixed within the grid. If true, the panel will not be affected by other panels' interactions
	Static *bool `json:"static,omitempty"`
}

// A library panel is a reusable panel that you can use in any dashboard.
// When you make a change to a library panel, that change propagates to all instances of where the panel is used.
// Library panels streamline reuse of panels across multiple dashboards.
type LibraryPanelRef struct {
	// Library panel name
	Name string `json:"name"`
	// Library panel uid
	Uid string `json:"uid"`
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
	// The matcher id. This is used to find the matcher implementation from registry.
	Id string `json:"id"`
	// The matcher options. This is specific to the matcher implementation.
	Options any `json:"options,omitempty"`
}

// Dashboard panels are the basic visualization building blocks.
type Panel struct {
	// The panel plugin type id. This is used to find the plugin to display the panel.
	Type string `json:"type"`
	// Unique identifier of the panel. Generated by Grafana when creating a new panel. It must be unique within a dashboard, but not globally.
	Id *uint32 `json:"id,omitempty"`
	// The version of the plugin that is used for this panel. This is used to find the plugin to display the panel and to migrate old panel configs.
	PluginVersion *string `json:"pluginVersion,omitempty"`
	// Tags for the panel.
	Tags []string `json:"tags,omitempty"`
	// Depends on the panel plugin. See the plugin documentation for details.
	Targets []Target `json:"targets,omitempty"`
	// Panel title.
	Title *string `json:"title,omitempty"`
	// Panel description.
	Description *string `json:"description,omitempty"`
	// Whether to display the panel without a background.
	Transparent bool `json:"transparent"`
	// The datasource used in all targets.
	Datasource *DataSourceRef `json:"datasource,omitempty"`
	// Grid position.
	GridPos *GridPos `json:"gridPos,omitempty"`
	// Panel links.
	Links []DashboardLink `json:"links,omitempty"`
	// Name of template variable to repeat for.
	Repeat *string `json:"repeat,omitempty"`
	// Direction to repeat in if 'repeat' is set.
	// `h` for horizontal, `v` for vertical.
	RepeatDirection *PanelRepeatDirection `json:"repeatDirection,omitempty"`
	// Id of the repeating panel.
	RepeatPanelId *int64 `json:"repeatPanelId,omitempty"`
	// The maximum number of data points that the panel queries are retrieving.
	MaxDataPoints *float64 `json:"maxDataPoints,omitempty"`
	// List of transformations that are applied to the panel data before rendering.
	// When there are multiple transformations, Grafana applies them in the order they are listed.
	// Each transformation creates a result set that then passes on to the next transformation in the processing pipeline.
	Transformations []DataTransformerConfig `json:"transformations"`
	// The min time interval setting defines a lower limit for the $__interval and $__interval_ms variables.
	// This value must be formatted as a number followed by a valid time
	// identifier like: "40s", "3d", etc.
	// See: https://grafana.com/docs/grafana/latest/panels-visualizations/query-transform-data/#query-options
	Interval *string `json:"interval,omitempty"`
	// Overrides the relative time range for individual panels,
	// which causes them to be different than what is selected in
	// the dashboard time picker in the top-right corner of the dashboard. You can use this to show metrics from different
	// time periods or days on the same dashboard.
	// The value is formatted as time operation like: `now-5m` (Last 5 minutes), `now/d` (the day so far),
	// `now-5d/d`(Last 5 days), `now/w` (This week so far), `now-2y/y` (Last 2 years).
	// Note: Panel time overrides have no effect when the dashboard’s time range is absolute.
	// See: https://grafana.com/docs/grafana/latest/panels-visualizations/query-transform-data/#query-options
	TimeFrom *string `json:"timeFrom,omitempty"`
	// Overrides the time range for individual panels by shifting its start and end relative to the time picker.
	// For example, you can shift the time range for the panel to be two hours earlier than the dashboard time picker setting `2h`.
	// Note: Panel time overrides have no effect when the dashboard’s time range is absolute.
	// See: https://grafana.com/docs/grafana/latest/panels-visualizations/query-transform-data/#query-options
	TimeShift *string `json:"timeShift,omitempty"`
	// Dynamically load the panel
	LibraryPanel *LibraryPanelRef `json:"libraryPanel,omitempty"`
	// It depends on the panel plugin. They are specified by the Options field in panel plugin schemas.
	Options any `json:"options"`
	// Field options allow you to change how the data is displayed in your visualizations.
	FieldConfig FieldConfigSource `json:"fieldConfig"`
}

type PanelRepeatDirection string

const (
	Horizontal PanelRepeatDirection = "h"
	Vertical   PanelRepeatDirection = "v"
)

// Maps numerical ranges to a display text and color.
// For example, if a value is within a certain range, you can configure a range value mapping to display Low or High rather than the number.
type RangeMap struct {
	Type string `json:"type"`
	// Range to match against and the result to apply when the value is within the range
	Options struct {
		// Min value of the range. It can be null which means -Infinity
		From *float64 `json:"from"`
		// Max value of the range. It can be null which means +Infinity
		To *float64 `json:"to"`
		// Config to apply when the value is within the range
		Result ValueMappingResult `json:"result"`
	} `json:"options"`
}

// Maps regular expressions to replacement text and a color.
// For example, if a value is www.example.com, you can configure a regex value mapping so that Grafana displays www and truncates the domain.
type RegexMap struct {
	Type string `json:"type"`
	// Regular expression to match against and the result to apply when the value matches the regex
	Options struct {
		// Regular expression to match against
		Pattern string `json:"pattern"`
		// Config to apply when the value matches the regex
		Result ValueMappingResult `json:"result"`
	} `json:"options"`
}

// Row panel
type RowPanel struct {
	// The panel type
	Type string `json:"type"`
	// Whether this row should be collapsed or not.
	Collapsed bool `json:"collapsed"`
	// Row title
	Title *string `json:"title,omitempty"`
	// Name of default datasource for the row
	Datasource *DataSourceRef `json:"datasource,omitempty"`
	// Row grid position
	GridPos *GridPos `json:"gridPos,omitempty"`
	// Unique identifier of the panel. Generated by Grafana when creating a new panel. It must be unique within a dashboard, but not globally.
	Id uint32 `json:"id"`
	// List of panels in the row
	Panels []Panel `json:"panels"`
	// Name of template variable to repeat for.
	Repeat *string `json:"repeat,omitempty"`
}

// Maps special values like Null, NaN (not a number), and boolean values like true and false to a display text and color.
// See SpecialValueMatch to see the list of special values.
// For example, you can configure a special value mapping so that null values appear as N/A.
type SpecialValueMap struct {
	Type    string `json:"type"`
	Options struct {
		// Special value to match against
		Match SpecialValueMatch `json:"match"`
		// Config to apply when the value matches the special value
		Result ValueMappingResult `json:"result"`
	} `json:"options"`
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

type StringOrArray struct {
	ValString *string  `json:"ValString,omitempty"`
	ValArray  []string `json:"ValArray,omitempty"`
}

type StringOrBool struct {
	ValString *string `json:"ValString,omitempty"`
	ValBool   *bool   `json:"ValBool,omitempty"`
}

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
	// Value represents a specified metric for the threshold, which triggers a visual change in the dashboard when this value is met or exceeded.
	// Nulls currently appear here when serializing -Infinity to JSON.
	Value *float64 `json:"value"`
	// Color represents the color of the visual change that will occur in the dashboard when the threshold value is met or exceeded.
	Color string `json:"color"`
}

// Thresholds configuration for the panel
type ThresholdsConfig struct {
	// Thresholds mode.
	Mode ThresholdsMode `json:"mode"`
	// Must be sorted by 'value', first value is always -Infinity
	Steps []Threshold `json:"steps"`
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
	// Whether timepicker is visible or not.
	Hidden bool `json:"hidden"`
	// Interval options available in the refresh picker dropdown.
	Refresh_intervals []string `json:"refresh_intervals"`
	// Whether timepicker is collapsed or not. Has no effect on provisioned dashboard.
	Collapse bool `json:"collapse"`
	// Whether timepicker is enabled or not. Has no effect on provisioned dashboard.
	Enable bool `json:"enable"`
	// Selectable options available in the time picker dropdown. Has no effect on provisioned dashboard.
	Time_options []string `json:"time_options"`
}

// Maps text values to a color or different display text and color.
// For example, you can configure a value mapping so that all instances of the value 10 appear as Perfection! rather than the number.
type ValueMap struct {
	Type string `json:"type"`
	// Map with <value_to_match>: ValueMappingResult. For example: { "10": { text: "Perfection!", color: "green" } }
	Options map[string]ValueMappingResult `json:"options"`
}

type ValueMapOrRangeMapOrRegexMapOrSpecialValueMap struct {
	ValValueMap        *ValueMap        `json:"ValValueMap,omitempty"`
	ValRangeMap        *RangeMap        `json:"ValRangeMap,omitempty"`
	ValRegexMap        *RegexMap        `json:"ValRegexMap,omitempty"`
	ValSpecialValueMap *SpecialValueMap `json:"ValSpecialValueMap,omitempty"`
}

// Result used as replacement with text and color when the value matches
type ValueMappingResult struct {
	// Text to display when the value matches
	Text *string `json:"text,omitempty"`
	// Text to use when the value matches
	Color *string `json:"color,omitempty"`
	// Icon to display when the value matches. Only specific visualizations.
	Icon *string `json:"icon,omitempty"`
	// Position in the mapping array. Only used internally.
	Index *int32 `json:"index,omitempty"`
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
	// Unique numeric identifier for the variable.
	Id string `json:"id"`
	// Type of variable
	Type VariableType `json:"type"`
	// Name of variable
	Name string `json:"name"`
	// Optional display name
	Label *string `json:"label,omitempty"`
	// Visibility configuration for the variable
	Hide VariableHide `json:"hide"`
	// Whether the variable value should be managed by URL query params or not
	SkipUrlSync bool `json:"skipUrlSync"`
	// Description of variable. It can be defined but `null`.
	Description *string `json:"description,omitempty"`
	// Query used to fetch values for a variable
	Query any `json:"query,omitempty"`
	// Data source used to fetch values for a variable. It can be defined but `null`.
	Datasource *DataSourceRef `json:"datasource,omitempty"`
	// Format to use while fetching all values from data source, eg: wildcard, glob, regex, pipe, etc.
	AllFormat *string `json:"allFormat,omitempty"`
	// Shows current selected variable text/value on the dashboard
	Current *VariableOption `json:"current,omitempty"`
	// Whether multiple values can be selected or not from variable value list
	Multi *bool `json:"multi,omitempty"`
	// Options that can be selected for a variable.
	Options []VariableOption `json:"options,omitempty"`
	Refresh *VariableRefresh `json:"refresh,omitempty"`
}

// Option to be selected in a variable.
type VariableOption struct {
	// Whether the option is selected or not
	Selected *bool `json:"selected,omitempty"`
	// Text to be displayed for the option
	Text StringOrArray `json:"text"`
	// Value of the option
	Value StringOrArray `json:"value"`
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
