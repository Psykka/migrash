package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create config file for migrash",
	Run: func(cmd *cobra.Command, args []string) {
		f, err := os.Create(".migrashrc")

		if err != nil {
			panic(err)
		}

		defer f.Close()

		initialConfig := `# example config file

MIGRATION_DIR=migrations
# DATABASE_URL=$DATABASE # from env
DATABASE_URL=mysql://root:password@localhost:3306/migrash
MIGRATION_TABLE=_migrash_migrations`

		_, err = f.WriteString(initialConfig)

		if err != nil {
			panic(err)
		}

		fmt.Println("Initial config file created!")
	},
}
