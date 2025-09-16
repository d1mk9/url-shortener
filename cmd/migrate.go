package cmd

import (
	"fmt"
	"log"
	"url-shortener/pkg/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"

	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
)

const migrationsDir = "migrations"

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Управление миграциями базы данных",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Применить все миграции вверх",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.MustLoad()

		db, err := goose.OpenDBWithDriver("postgres", cfg.PostgresDSN())
		if err != nil {
			log.Fatalf("failed to open DB: %v", err)
		}
		defer db.Close()

		if err := goose.Up(db, migrationsDir); err != nil {
			log.Fatalf("goose up: %v", err)
		}
		fmt.Println("✅ Миграции применены")
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Откатить одну миграцию вниз",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.MustLoad()

		db, err := goose.OpenDBWithDriver("postgres", cfg.PostgresDSN())
		if err != nil {
			log.Fatalf("failed to open DB: %v", err)
		}
		defer db.Close()

		if err := goose.Down(db, migrationsDir); err != nil {
			log.Fatalf("goose down: %v", err)
		}
		fmt.Println("✅ Миграция откатилась")
	},
}

func init() {
	migrateCmd.AddCommand(migrateUpCmd, migrateDownCmd)
	rootCmd.AddCommand(migrateCmd)
}
