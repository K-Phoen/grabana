package golang

import (
	"github.com/K-Phoen/sdk"
	"github.com/dave/jennifer/jen"
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
		jen.Lit(target.Expr),
	}

	if target.LegendFormat != "" {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/target/prometheus", "Legend").Call(jen.Lit(target.LegendFormat)),
		)
	}
	if target.RefID != "" {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/target/prometheus", "Ref").Call(jen.Lit(target.RefID)),
		)
	}
	if target.Hide {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/target/prometheus", "Hide").Call(),
		)
	}
	if target.Instant {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/target/prometheus", "Instant").Call(),
		)
	}
	if target.IntervalFactor != 0 {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/target/prometheus", "IntervalFactor").Call(jen.Lit(target.IntervalFactor)),
		)
	}

	formatConstName := "FormatTimeSeries"
	switch target.Format {
	case "table":
		formatConstName = "FormatTable"
	case "heatmap":
		formatConstName = "FormatHeatmap"
	case "time_series":
		formatConstName = "FormatTimeSeries"
	default:
		encoder.logger.Warn("unhandled prometheus target format '%s': using 'time_series' instead", zap.String("format", target.Format))

	}

	// only emit code if the default isn't used
	if formatConstName != "FormatTimeSeries" {
		settings = append(
			settings,
			jen.Qual(packageImportPath+"/target/prometheus", "Format").Call(
				jen.Qual(packageImportPath+"/target/prometheus", formatConstName),
			),
		)
	}

	return jen.Qual(packageImportPath+"/"+grabanaPackage, "WithPrometheusTarget").Call(
		settings...,
	)
}
