package stat

import (
	"fmt"

	"github.com/K-Phoen/grabana/errors"
	"github.com/K-Phoen/grabana/links"
	"github.com/K-Phoen/grabana/scheme"
	"github.com/K-Phoen/grabana/target/graphite"
	"github.com/K-Phoen/grabana/target/influxdb"
	"github.com/K-Phoen/grabana/target/prometheus"
	"github.com/K-Phoen/grabana/target/stackdriver"
	"github.com/K-Phoen/sdk"
)

// Option represents an option that can be used to configure a stat panel.
type Option func(stat *Stat) error

type ThresholdStep struct {
	Color string
	Value *float64
}

// TextMode controls if name and value is displayed or just name.
type TextMode string

const (
	TextAuto         TextMode = "auto"
	TextValue        TextMode = "value"
	TextName         TextMode = "name"
	TextValueAndName TextMode = "value_and_name"
	TextNone         TextMode = "none"
)

// OrientationMode controls the layout.
type OrientationMode string

const (
	OrientationAuto       OrientationMode = ""
	OrientationHorizontal OrientationMode = "horizontal"
	OrientationVertical   OrientationMode = "vertical"
)

// ReductionType lets you set the function that your entire query is reduced into a
// single value with.
type ReductionType int

const (
	// Min displays the smallest value of the series.
	Min ReductionType = iota
	// Max displays the largest value of the series.
	Max
	// Avg displays the average of the series.
	Avg

	// First displays the first value of the series.
	First
	// FirstNonNull displays the first non-null value of the series.
	FirstNonNull
	// Last displays the last value of the series.
	Last
	// LastNonNull displays the last non-null value of the series.
	LastNonNull

	// Total displays the sum of values in the series.
	Total
	// Count displays the number of value in the series.
	Count
	// Range displays the difference between the minimum and maximum values.
	Range
)

// Stat represents a stat panel.
type Stat struct {
	Builder *sdk.Panel
}

// New creates a new stat panel.
func New(title string, options ...Option) (*Stat, error) {
	panel := &Stat{Builder: sdk.NewStat(title)}

	panel.Builder.IsNew = false

	for _, opt := range append(defaults(), options...) {
		if err := opt(panel); err != nil {
			return nil, err
		}
	}

	return panel, nil
}

func defaults() []Option {
	return []Option{
		Text(TextValue),
		ColorValue(),
		Orientation(OrientationVertical),
		Span(6),
		ValueType(Last),
		NoValue("N/A"),
		ColorScheme(scheme.ThresholdsValue(scheme.Last)),
	}
}

// Links adds links to be displayed on this panel.
func Links(panelLinks ...links.Link) Option {
	return func(stat *Stat) error {
		stat.Builder.Links = make([]sdk.Link, 0, len(panelLinks))

		for _, link := range panelLinks {
			stat.Builder.Links = append(stat.Builder.Links, link.Builder)
		}

		return nil
	}
}

// DataSource sets the data source to be used by the panel.
func DataSource(source string) Option {
	return func(stat *Stat) error {
		stat.Builder.Datasource = &sdk.DatasourceRef{LegacyName: source}

		return nil
	}
}

// WithPrometheusTarget adds a prometheus query to the graph.
func WithPrometheusTarget(query string, options ...prometheus.Option) Option {
	target := prometheus.New(query, options...)

	return func(stat *Stat) error {
		stat.Builder.AddTarget(&sdk.Target{
			RefID:          target.Ref,
			Hide:           target.Hidden,
			Expr:           target.Expr,
			IntervalFactor: target.IntervalFactor,
			Interval:       target.Interval,
			Step:           target.Step,
			LegendFormat:   target.LegendFormat,
			Instant:        target.Instant,
			Format:         target.Format,
		})

		return nil
	}
}

// WithGraphiteTarget adds a Graphite target to the graph.
func WithGraphiteTarget(query string, options ...graphite.Option) Option {
	target := graphite.New(query, options...)

	return func(stat *Stat) error {
		stat.Builder.AddTarget(target.Builder)

		return nil
	}
}

// WithInfluxDBTarget adds an InfluxDB target to the graph.
func WithInfluxDBTarget(query string, options ...influxdb.Option) Option {
	target := influxdb.New(query, options...)

	return func(stat *Stat) error {
		stat.Builder.AddTarget(target.Builder)

		return nil
	}
}

// WithStackdriverTarget adds a stackdriver query to the graph.
func WithStackdriverTarget(target *stackdriver.Stackdriver) Option {
	return func(stat *Stat) error {
		stat.Builder.AddTarget(target.Builder)

		return nil
	}
}

// Span sets the width of the panel, in grid units. Should be a positive
// number between 1 and 12. Example: 6.
func Span(span float32) Option {
	return func(stat *Stat) error {
		if span < 1 || span > 12 {
			return fmt.Errorf("span must be between 1 and 12: %w", errors.ErrInvalidArgument)
		}

		stat.Builder.Span = span

		return nil
	}
}

// Height sets the height of the panel, in pixels. Example: "400px".
func Height(height string) Option {
	return func(stat *Stat) error {
		stat.Builder.Height = &height

		return nil
	}
}

// Description annotates the current visualization with a human-readable description.
func Description(content string) Option {
	return func(stat *Stat) error {
		stat.Builder.Description = &content

		return nil
	}
}

// Transparent makes the background transparent.
func Transparent() Option {
	return func(stat *Stat) error {
		stat.Builder.Transparent = true

		return nil
	}
}

// Unit sets the unit of the data displayed on this axis.
func Unit(unit string) Option {
	return func(stat *Stat) error {
		stat.Builder.StatPanel.FieldConfig.Defaults.Unit = unit

		return nil
	}
}

// Decimals sets the number of decimals that should be displayed.
func Decimals(count int) Option {
	return func(stat *Stat) error {
		if count < 0 {
			return fmt.Errorf("decimals must be greater than 0: %w", errors.ErrInvalidArgument)
		}

		stat.Builder.StatPanel.FieldConfig.Defaults.Decimals = &count

		return nil
	}
}

// SparkLine displays the spark line summary of the series in addition to the
// stat.
func SparkLine() Option {
	return func(stat *Stat) error {
		stat.Builder.StatPanel.Options.GraphMode = "area"

		return nil
	}
}

// SparkLineYMin defines the smallest value expected on the Y axis of the spark line.
func SparkLineYMin(value float64) Option {
	return func(stat *Stat) error {
		stat.Builder.StatPanel.FieldConfig.Defaults.Min = &value

		return nil
	}
}

// SparkLineYMax defines the largest value expected on the Y axis of the spark line.
func SparkLineYMax(value float64) Option {
	return func(stat *Stat) error {
		stat.Builder.StatPanel.FieldConfig.Defaults.Max = &value

		return nil
	}
}

// ValueType configures how the series will be reduced to a single value.
func ValueType(valueType ReductionType) Option {
	return func(stat *Stat) error {
		var valType string

		switch valueType {
		case First:
			valType = "first"
		case FirstNonNull:
			valType = "firstNotNull"
		case Last:
			valType = "last"
		case LastNonNull:
			valType = "lastNotNull"

		case Min:
			valType = "min"
		case Max:
			valType = "max"
		case Avg:
			valType = "mean"

		case Count:
			valType = "count"
		case Total:
			valType = "sum"
		case Range:
			valType = "range"

		default:
			return fmt.Errorf("unknown value type: %w", errors.ErrInvalidArgument)
		}
		stat.Builder.StatPanel.Options.ReduceOptions.Calcs = []string{valType}

		return nil
	}
}

// ValueFontSize sets the font size used to display the value.
func ValueFontSize(size int) Option {
	return func(stat *Stat) error {
		stat.Builder.StatPanel.Options.Text.ValueSize = size

		return nil
	}
}

// TitleFontSize sets the font size used to display the title.
func TitleFontSize(size int) Option {
	return func(stat *Stat) error {
		stat.Builder.StatPanel.Options.Text.TitleSize = size

		return nil
	}
}

// ColorNone will not color the value or the background.
func ColorNone() Option {
	return func(stat *Stat) error {
		stat.Builder.StatPanel.Options.ColorMode = "none"

		return nil
	}
}

// ColorValue will show the threshold's colors on the value itself.
func ColorValue() Option {
	return func(stat *Stat) error {
		stat.Builder.StatPanel.Options.ColorMode = "value"

		return nil
	}
}

// ColorBackground will show the threshold's colors in the background.
func ColorBackground() Option {
	return func(stat *Stat) error {
		stat.Builder.StatPanel.Options.ColorMode = "background"

		return nil
	}
}

// AbsoluteThresholds changes the background and value colors dynamically within the
// panel, depending on the value. The threshold is defined by a series of steps
// values which, each having a value and an associated color.
func AbsoluteThresholds(steps []ThresholdStep) Option {
	return func(stat *Stat) error {
		sdkSteps := make([]sdk.ThresholdStep, 0, len(steps))
		for _, step := range steps {
			sdkSteps = append(sdkSteps, sdk.ThresholdStep{
				Color: step.Color,
				Value: step.Value,
			})
		}

		stat.Builder.StatPanel.FieldConfig.Defaults.Thresholds = sdk.Thresholds{
			Mode:  "absolute",
			Steps: sdkSteps,
		}

		return nil
	}
}

// RelativeThresholds changes the background and value colors dynamically within the
// panel, depending on the value. The threshold is defined by a series of steps
// values which, each having a value defined as a percentage and an associated color.
func RelativeThresholds(steps []ThresholdStep) Option {
	return func(stat *Stat) error {
		sdkSteps := make([]sdk.ThresholdStep, 0, len(steps))
		for _, step := range steps {
			sdkSteps = append(sdkSteps, sdk.ThresholdStep{
				Color: step.Color,
				Value: step.Value,
			})
		}

		stat.Builder.StatPanel.FieldConfig.Defaults.Thresholds = sdk.Thresholds{
			Mode:  "percentage",
			Steps: sdkSteps,
		}

		return nil
	}
}

type StatValueMapping struct {
	sdk.ValueMapping
}

type ValueMappingBuilder struct {
	options       func(*StatValueMapping)
	resultOptions []ValueMappingResultOption
	resultKey     string
	result        struct {
		text  *string
		color *string
	}
}

func (vmb *ValueMappingBuilder) build() *StatValueMapping {
	vm := &StatValueMapping{}
	vmb.options(vm)
	vm.Options[vmb.resultKey] = &sdk.ValueMappingResult{}
	for _, resultOption := range vmb.resultOptions {
		resultOption(vm)
	}
	return vm
}

type ValueMappingResultOption func(*StatValueMapping)

func (vmb *ValueMappingBuilder) Text(text string) *ValueMappingBuilder {
	resultFunc := func(mapping *StatValueMapping) {
		result, ok := mapping.Options[vmb.resultKey].(*sdk.ValueMappingResult)
		if ok {
			result.Text = &text
		}
	}
	vmb.resultOptions = append(vmb.resultOptions, resultFunc)
	return vmb
}

func (vmb *ValueMappingBuilder) Color(color string) *ValueMappingBuilder {
	resultFunc := func(mapping *StatValueMapping) {
		result, ok := mapping.Options[vmb.resultKey].(*sdk.ValueMappingResult)
		if ok {
			result.Color = &color
		}
	}
	vmb.resultOptions = append(vmb.resultOptions, resultFunc)
	return vmb
}

func ValueMappings(options ...*ValueMappingBuilder) Option {
	return func(stat *Stat) error {
		mappings := make([]sdk.ValueMapping, 0, len(options))

		for idx, option := range options {
			mapping := option.build()
			result, ok := mapping.Options[option.resultKey].(*sdk.ValueMappingResult)
			if ok {
				result.Index = &idx
			}
			mappings = append(mappings, mapping.ValueMapping)
		}
		stat.Builder.StatPanel.FieldConfig.Defaults.ValueMappings = mappings
		return nil
	}
}

// RangeMapping allows to translate a range of values from the summary stat into explicit
// text or color.
func RangeMapping(from float64, to float64) *ValueMappingBuilder {
	builder := &ValueMappingBuilder{}
	builder.resultKey = "result"
	builder.options = func(vm *StatValueMapping) {
		vm.MappingType = string(sdk.MappingTypeRange)
		vm.Options = map[string]interface{}{
			sdk.MappingOptionFrom: from,
			sdk.MappingOptionTo:   to,
		}
	}
	return builder
}

// RegexMapping allows to translate the value of the summary stat matching a regexp into explicit
// text or color.
func RegexMapping(pattern string) *ValueMappingBuilder {
	builder := &ValueMappingBuilder{}
	builder.resultKey = "result"
	builder.options = func(vm *StatValueMapping) {
		vm.MappingType = string(sdk.MappingTypeRegex)
		vm.Options = map[string]interface{}{
			sdk.MappingOptionPattern: pattern,
		}
	}
	return builder
}

// ValueMapping allows to translate the value of the summary stat into explicit
// text or color, using a list of mapping functions.
func ValueMapping(match string) *ValueMappingBuilder {
	builder := &ValueMappingBuilder{}
	builder.resultKey = match
	builder.options = func(vm *StatValueMapping) {
		vm.MappingType = string(sdk.MappingTypeValue)
		vm.Options = map[string]interface{}{}
	}
	return builder
}

type SpecialValue string

const SpecialNull = SpecialValue("null")
const SpecialNullOrNaN = SpecialValue("null+nan")
const SpecialNaN = SpecialValue("nan")
const SpecialTrue = SpecialValue("true")
const SpecialFalse = SpecialValue("false")
const SpecialEmpty = SpecialValue("empty")

// SpecialMapping allows to translate special values of the summary stat into explicit
// text or color.
func SpecialMapping(value SpecialValue) *ValueMappingBuilder {
	builder := &ValueMappingBuilder{}
	builder.resultKey = "result"
	valueStr := string(value)
	builder.options = func(vm *StatValueMapping) {
		vm.MappingType = string(sdk.MappingTypeSpecial)
		vm.Options = map[string]interface{}{
			sdk.MappingOptionMatch: valueStr,
		}
	}
	return builder
}

// Repeat configures repeating a panel for a variable
func Repeat(repeat string) Option {
	return func(stat *Stat) error {
		stat.Builder.Repeat = &repeat

		return nil
	}
}

// Text indicates if name and value is displayed or just name.
func Text(mode TextMode) Option {
	return func(stat *Stat) error {
		stat.Builder.StatPanel.Options.TextMode = string(mode)

		return nil
	}
}

// Orientation changes the orientation of the layout.
func Orientation(mode OrientationMode) Option {
	return func(stat *Stat) error {
		stat.Builder.StatPanel.Options.Orientation = string(mode)

		return nil
	}
}

// ColorScheme configures the color scheme.
func ColorScheme(options ...scheme.Option) Option {
	return func(stat *Stat) error {
		scheme.New(&stat.Builder.StatPanel.FieldConfig, options...)

		return nil
	}
}

// NoValue defines what to show when there is no value.
func NoValue(text string) Option {
	return func(stat *Stat) error {
		stat.Builder.StatPanel.FieldConfig.Defaults.NoValue = text

		return nil
	}
}

func float64Ptr(input float64) *float64 {
	return &input
}
