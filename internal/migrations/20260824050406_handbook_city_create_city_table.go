package migrations

import (
	"github.com/exgamer/gosdk-db-core/pkg/migration"
	"gorm.io/gorm"
)

func migration20260824050406() *migration.Migration {
	return &migration.Migration{
		ID: "20260824050406_handbook_city_create_city_table",
		Migrate: func(tx *gorm.DB) error {
			return tx.Exec(`CREATE TABLE city (
				id SERIAL PRIMARY KEY,
				name VARCHAR(100) NOT NULL,
				status INT NOT NULL DEFAULT 0
			)`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Exec(`DROP TABLE city`).Error
		},
	}
}
