// Package cmd is the root of our application
package cmd

import (
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"go.infratographer.com/kubevirt-provider/internal/config"

	"go.infratographer.com/x/loggingx"
	"go.infratographer.com/x/otelx"
	"go.infratographer.com/x/versionx"
	"go.infratographer.com/x/viperx"
)

var (
	cfgFile string
	logger  *zap.SugaredLogger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubevirt-provider",
	Short: "A controller for processing requests related to kubevirt provisioning",
	Long:  "A controller for processing requests related to kubevirt provisioning",
}

var appName = "kubevirtprovider"

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/."+appName+".yaml)")
	viperx.MustBindFlag(viper.GetViper(), "config", rootCmd.PersistentFlags().Lookup("config"))

	// Logging flags
	loggingx.MustViperFlags(viper.GetViper(), rootCmd.PersistentFlags())

	// Register version command
	versionx.RegisterCobraCommand(rootCmd, func() { versionx.PrintVersion(logger) })
	otelx.MustViperFlags(viper.GetViper(), rootCmd.Flags())
}

// Setup migrate command
/*	goosex.RegisterCobraCommand(rootCmd, func() {
	goosex.SetBaseFS(dbm.Migrations)
	goosex.SetDBURI(config.AppConfig.CRDB.URI)
	goosex.SetLogger(logger)
})*/

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// load the config file
		viper.AddConfigPath(home)
		viper.SetConfigName("." + appName)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.SetEnvPrefix("kubevirtprovider")

	// read in environment variables that match
	viper.AutomaticEnv()

	// loads our config.AppConfig struct with the values bound by
	// viper. Then, anywhere we need these values, we can just return to AppConfig
	// instead of performing viper.GetString(...), viper.GetBool(...), etc.
	err := viper.Unmarshal(&config.AppConfig)
	cobra.CheckErr(err)

	// setupLogging()
	logger = loggingx.InitLogger(appName, config.AppConfig.Logging)

	// If a config file is found, read it in.
	err = viper.ReadInConfig()
	if err == nil {
		logger.Infow("using config file",
			"file", viper.ConfigFileUsed(),
		)
	}
}
