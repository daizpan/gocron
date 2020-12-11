package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// NewRootCmd represents the base command when called without any subcommands
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "gocron",
		Version: "0.1.0",
		Short:   "Cron like scheduler",
	}
	cobra.OnInitialize(initConfig)
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/gocron.yaml)")

	cmd.AddCommand(NewRunCmd())
	return cmd
}

func Execute() {
	cmd := NewRootCmd()
	cmd.SetOutput(os.Stdout)
	cmd.SetErr(os.Stderr)
	if err := cmd.Execute(); err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".gocron" (without extension).
		viper.SetConfigName("gocron")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		log.Println(err)
	}
}
