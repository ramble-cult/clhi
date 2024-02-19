/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/ramble-cult/clhi/cmd"
	"github.com/spf13/viper"
)

func main() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".clhi") // Register config file name (no extension)
	viper.SetConfigType("yaml")  // Look for specific type
	viper.ReadInConfig()
	viper.WatchConfig()

	cmd.Execute()
}
