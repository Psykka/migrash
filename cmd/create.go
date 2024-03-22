package cmd

import (
	"errors"
	"fmt"
	"migrash/internal/config"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create [migration name]",
	Short: "Create a new migration file and a directory if it does not exist",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires migration name")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.ParseConfig()

		if err != nil {
			panic(err)
		}

		if err := os.MkdirAll(config.MigrationDir, os.ModePerm); err != nil {
			panic(err)
		}

		timestamp := time.Now().Format("20060102150405")

		dir := fmt.Sprintf("%s/%s-%s", config.MigrationDir, timestamp, args[0])
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			panic(err)
		}

		files := []string{"up", "down"}

		for _, file := range files {
			filePath := fmt.Sprintf("%s/%s.sql", dir, file)

			createdFile, err := os.Create(filePath)

			if err != nil {
				panic(err)
			}

			defer createdFile.Close()

			createdFile.WriteString(fmt.Sprintf("-- %s\n", strings.ToUpper(file)))

			defer fmt.Printf("Created %s\n", filePath)
		}

		fmt.Printf("Migration %s created!\n", args[0])
	},
}
