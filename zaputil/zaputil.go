// Copyright (c) 2022 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package zaputil

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// A FlagConfig enables changing the command flags what will be used for zap configuration.
type FlagConfig struct {
	// LevelName specifies the flag to use for the zap level.
	LevelName string
	// EncoderName specifies the flag to use for the zap encoding.
	EncoderName string

	level   zapLevelFlag
	encoder zapEncoderFlag
}

// NewFlagConfig returns a FlagConfig with the default flag names.
func NewFlagConfig() FlagConfig {
	return FlagConfig{
		LevelName:   "zap-log-level",
		EncoderName: "zap-encoder",
	}
}

// Bind will register the flags on the flagset so they can be parsed. If not set, the level will
// default to 'info', and the encoding will default to 'console'.
func (fc *FlagConfig) Bind(fs *flag.FlagSet) {
	fc.level.Level = zap.InfoLevel
	fc.encoder.string = "console"

	fs.Var(&fc.level, fc.LevelName, "Zap level to configure the verbosity of logging.")
	fs.Var(&fc.encoder, fc.EncoderName, "Zap log encoding (one of 'json' or 'console')")
}

// GetConfig returns a Zap configuration based off of the "production" configuration from zap.
// It will have the level and encoder specified in the command flags, and it will use ISO-8601
// timestamps. Note that it is configured with sampling, so it could drop some logs. See
// https://github.com/uber-go/zap/blob/master/FAQ.md#why-sample-application-logs
func (fc *FlagConfig) GetConfig() zap.Config {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		num := int8(l)
		if num > -2 { // These levels are known by zap, as "info", "error", etc.
			enc.AppendString(l.String())
		} else { // Zap doesn't like these levels as much. Format them like "lvl-n"
			enc.AppendString("lvl" + strconv.Itoa(int(num)))
		}
	}

	cfg.Level = zap.NewAtomicLevelAt(fc.level.Level)
	cfg.Encoding = fc.encoder.string

	return cfg
}

// BuildForCtrl returns a zap.Logger built to work well with the controller-runtime.
// Note: using the controller-runtime logger directly will not record the correct
// filename/line - instead, call `WithName` or `WithValues` in your function, and
// use the resulting logger.
func (fc *FlagConfig) BuildForCtrl() (*zap.Logger, error) {
	return fc.GetConfig().Build()
}

// BuildForKlog returns a zap.Logger built from given config, made to work well with klog. The
// FlagSet should be the same one klog is bound to. Usually the given config should be created
// from `FlagConfig.GetConfig()`.
func BuildForKlog(cfg zap.Config, klogFlagSet *flag.FlagSet) (*zap.Logger, error) {
	if cfg.Encoding == "console" {
		// klog already adds a newline, so have zap skip adding another in console mode
		cfg.EncoderConfig.SkipLineEnding = true
	}

	klogV := klogFlagSet.Lookup("v")
	if klogV == nil {
		return nil, fmt.Errorf("no 'v' flag found in given FlagSet")
	}

	klogLevel, err := strconv.Atoi(klogV.Value.String())
	if err != nil {
		return nil, fmt.Errorf("invalid value passed in 'v' flag, couldn't convert to int: %w", err)
	}

	cfg.Level = zap.NewAtomicLevelAt(zapcore.Level(int8(-1 * klogLevel)))

	return cfg.Build()
}

// SyncWithGlogFlags copies the values of all the flags on the command line with the given FlagSet.
// It is intended to help workaround problems with interoperability between glog and klog.
// See https://github.com/kubernetes/klog/blob/v2.40.1/examples/coexist_glog/coexist_glog.go#L18
func SyncWithGlogFlags(klogFlagSet *flag.FlagSet) error {
	var syncErr error

	flag.CommandLine.VisitAll(func(f1 *flag.Flag) {
		f2 := klogFlagSet.Lookup(f1.Name)
		if f2 != nil {
			value := f1.Value.String()
			if err := f2.Value.Set(value); err != nil {
				// skip this one flag - klog somehow doesn't like glog's default
				if f1.Name != "log_backtrace_at" {
					syncErr = fmt.Errorf("unable to sync klog flag: '%s', value '%s'",
						f1.Name, f1.Value.String())
				}
			}
		}
	})

	return syncErr
}

type zapLevelFlag struct {
	zapcore.Level
}

var _ flag.Value = &zapLevelFlag{}

// Set ensures that the level passed by a flag is valid
func (f *zapLevelFlag) Set(val string) error {
	lval := strings.ToLower(val)

	level, err := zapcore.ParseLevel(lval)
	if err != nil { // couldn't parse as text, try parsing as a number
		l, err := strconv.Atoi(lval)
		if err != nil || l < 0 {
			return fmt.Errorf("invalid log level \"%s\"", val)
		}

		level = zapcore.Level(int8(-1 * l))
	}

	f.Level = level

	return nil
}

// Type satisfies the interface for pflag.Value
func (f *zapLevelFlag) Type() string {
	return "level"
}

type zapEncoderFlag struct {
	string
}

var _ flag.Value = &zapEncoderFlag{}

// Set ensures that the encoder passed by a flag is valid
func (f *zapEncoderFlag) Set(val string) error {
	lval := strings.ToLower(val)

	// Use a map in case we ever add a custom encoder.
	validEncoders := map[string]bool{
		"json":    true,
		"console": true,
	}

	if !validEncoders[lval] {
		return fmt.Errorf("invalid encoder type \"%s\"", val)
	}

	f.string = lval

	return nil
}

// Type satisfies the interface for pflag.Value
func (f *zapEncoderFlag) Type() string {
	return "encoder"
}

// String satisfies the interface for flag.Value
func (f *zapEncoderFlag) String() string {
	return f.string
}
