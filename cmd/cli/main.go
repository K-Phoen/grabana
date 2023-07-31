package main

import (
	"fmt"
	"os"

	"github.com/K-Phoen/grabana/cmd/cli/cmd"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var version = "SNAPSHOT"

func main() {
	logger, err := createLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not create logger: %s", err)
		os.Exit(1)
	}

	root := &cobra.Command{Use: "grabana"}
	root.Version = version
	root.SilenceUsage = true

	root.AddCommand(cmd.Apply())
	root.AddCommand(cmd.Validate())
	root.AddCommand(cmd.SelfUpdate(version))
	root.AddCommand(cmd.Render())
	root.AddCommand(cmd.ConvertGo(logger))

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func createLogger() (*zap.Logger, error) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableStacktrace: true,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "console",
		EncoderConfig:    encoderCfg,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return cfg.Build()
}
