package migrations

import (
	"github.com/exgamer/gosdk-db-core/pkg/migration"
	"gorm.io/gorm"
)

func migration20260824050450() *migration.Migration {
	return &migration.Migration{
		ID: "20260824050450_handbook_city_seed_initial_cities",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`INSERT INTO city (name, status) VALUES
				('Astana', 1),
				('Almaty', 1),
				('Shymkent', 1)`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`DELETE FROM city WHERE name IN ('Astana', 'Almaty', 'Shymkent')`).Error
		},
	}
}
