package console

import (
	postgres "github.com/exgamer/gosdk-postgres-core/pkg/app"
	"github.com/spf13/cobra"
)

// migrateCmd — группа команд управления миграциями схемы БД (см.
// MIGRATIONS.md). PostgresKernel сам миграции не трогает — весь накат и
// откат идёт через Migrator, отдельно от Init/DI-жизненного цикла кернела.
func migrateCmd(m *postgres.Migrator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "manage DB schema migrations",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "up",
		Short: "apply all pending migrations",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Up()
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "down",
		Short: "rollback the last applied migration",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.RollbackLast()
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "down-to <id>",
		Short: "rollback everything after <id> (exclusive)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.RollbackTo(args[0])
		},
	})

	return cmd
}
