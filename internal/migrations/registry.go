package migrations

import "github.com/exgamer/gosdk-db-core/pkg/migration"

// All возвращает все миграции проекта в порядке применения.
// Генерируется и обновляется командой `codegen migration add` — не редактировать руками.
func All() []*migration.Migration {
	return []*migration.Migration{
		migration20260824050406(),
		migration20260824050450(),
	}
}
