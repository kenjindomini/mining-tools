/*
Package miningtools contains the various supported CLI commands for mining-tools
Copyright © 2020 Keith Olenchak <kenjin.domini@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package miningtools

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mining-tools",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initLog)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (Default: $HOME/mining-tools.yaml)")
	rootCmd.PersistentFlags().String("log", "", "log file (Default:  $HOME/mining-tools.log)")
	viper.BindPFlag("miningtools.logging.file", rootCmd.PersistentFlags().Lookup("log"))
	rootCmd.PersistentFlags().Uint32("log-level", 4, "Sets global log level 5=Debug 0=Panic/virtually silent (Default: 4)")
	viper.BindPFlag("miningtools.logging.level", rootCmd.PersistentFlags().Lookup("log-level"))
	rootCmd.PersistentFlags().String("timeseriesDB", "127.0.0.1:9009", "timeseriesDB address (Default:  127.0.0.1:9009)")
	viper.BindPFlag("miningtools.timeseriesDB.address", rootCmd.PersistentFlags().Lookup("timeseriesDB"))
	rootCmd.PersistentFlags().String("timeseriesProtocol", "InfluxDB", "timeseriesDB protocol, supports InfluxDB (Default:  InfluxDB)")
	viper.BindPFlag("miningtools.timeseriesDB.protocol", rootCmd.PersistentFlags().Lookup("timeseriesProtocol"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

		// Search config in home directory with name ".mining-tools" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName("mining-tools.yml")
		viper.SetConfigType("yaml")

		fmt.Println(fmt.Sprintf("Reading config file %s\\%s, if it exists", home, ".mining-tools"))
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initLog() {
	logFile := viper.GetString("miningtools.logging.file")
	logLevel := log.Level(viper.GetUint32("miningtools.logging.level"))
	if logFile == "" {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		logFile = home + "\\mining-tools.log"
	}

	log.SetLevel(logLevel)

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
}
