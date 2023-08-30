package pubsub

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.infratographer.com/x/events"
	"go.infratographer.com/x/viperx"
)

// Config defines the configuration for setting up events.
type Config struct {
	events.Config      `mapstructure:",squash"`
	Topics             []string
	MaxProcessAttempts uint64
}

// MustViperFlags sets the cobra flags and viper config for events.
func MustViperFlags(v *viper.Viper, flags *pflag.FlagSet, appName string) {
	events.MustViperFlags(v, flags, appName)

	flags.StringSlice("events-topics", []string{}, "event topics to subscribe to")
	viperx.MustBindFlag(v, "events.topics", flags.Lookup("events-topics"))

	flags.Uint64("events-max-process-attempts", 0, "maximum number of times an event may be processed")
	viperx.MustBindFlag(v, "events.maxprocessattempts", flags.Lookup("events-max-process-attempts"))
}
