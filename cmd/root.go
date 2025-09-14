package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "url-shortener",
	Short: "URL shortener service",
	Long:  `CLI для управления сервисом url-shortener (запуск сервера, миграции базы данных и др).`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
