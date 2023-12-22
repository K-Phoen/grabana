package golang

import (
	"github.com/K-Phoen/jennifer/jen"
	"github.com/K-Phoen/sdk"
	"go.uber.org/zap"
)

func (encoder *Encoder) encodeTargets(targets []sdk.Target, grabanaPackage string) []jen.Code {
	var targetsStmt []jen.Code

	for _, target := range targets {
		encodedTarget := encoder.encodeTarget(target, grabanaPackage)
		if encodedTarget == nil {
			continue
		}

		targetsStmt = append(targetsStmt, encodedTarget)
	}

	return targetsStmt
}

func (encoder *Encoder) encodeTarget(target sdk.Target, grabanaPackage string) jen.Code {
	// looks like a prometheus target
	if target.Expr != "" {
		return encoder.encodePrometheusTarget(target, grabanaPackage)
	}

	/*
		// looks like graphite
		if target.Target != "" {
			return encoder.encodeGraphiteTarget(target)
		}

		// looks like influxdb
		if target.Measurement != "" {
			return encoder.encodeInfluxDBTarget(target)
		}

		// looks like stackdriver
		if target.MetricType != "" {
			return encoder.encodeStackdriverTarget(target)
		}
	*/

	encoder.logger.Warn("unhandled target type: skipped", zap.Any("target", target))

	return nil
}

func (encoder *Encoder) encodePrometheusTarget(target sdk.Target, grabanaPackage string) jen.Code {
	settings := []jen.Code{
		lit(target.Expr),
	}

	if target.LegendFormat != "" {
		settings = append(
			settings,
			qual("target/prometheus", "Legend").Call(lit(target.LegendFormat)),
		)
	}
	if target.RefID != "" {
		settings = append(
			settings,
			qual("target/prometheus", "Ref").Call(lit(target.RefID)),
		)
	}
	if target.Hide {
		settings = append(
			settings,
			qual("target/prometheus", "Hide").Call(),
		)
	}
	if target.Instant {
		settings = append(
			settings,
			qual("target/prometheus", "Instant").Call(),
		)
	}
	if target.IntervalFactor != 0 {
		settings = append(
			settings,
			qual("target/prometheus", "IntervalFactor").Call(lit(target.IntervalFactor)),
		)
	}

	formatConstName := "FormatTimeSeries"
	switch target.Format {
	case "table":
		formatConstName = "FormatTable"
	case "heatmap":
		formatConstName = "FormatHeatmap"
	case "time_series":
	case "":
		formatConstName = "FormatTimeSeries"
	default:
		encoder.logger.Warn("unhandled prometheus target format: using 'time_series' instead", zap.String("format", target.Format))

	}

	// only emit code if the default isn't used
	if formatConstName != "FormatTimeSeries" {
		settings = append(
			settings,
			qual("target/prometheus", "Format").Call(
				qual("target/prometheus", formatConstName),
			),
		)
	}

	return qual(grabanaPackage, "WithPrometheusTarget").MultiLineCall(settings...)
}
