package mysql

import (
	"database/sql"
	root "xiaoyun/pkg"

	"github.com/DavidHuie/gomigrate"
)

// Migrate 数据库迁移
type Migrate struct {
	migrator *gomigrate.Migrator
}

// NewMigrate 生成新的数据库迁移对象
func NewMigrate(db *sql.DB, migrator string, logger gomigrate.Logger) (*Migrate, error) {

	var customErr root.Error
	customErr.Op = "mysql.NewMigrate"

	goMirator, err := gomigrate.NewMigratorWithLogger(db, gomigrate.Mysql{}, migrator, logger)
	if err != nil {
		return nil, err
	}

	var migrate Migrate
	migrate.migrator = goMirator

	return &migrate, nil

}

// Up 升级所以可用版本
func (m *Migrate) Up() error {

	return m.migrator.Migrate()

}

// Down 降级到最近一个可用版本
func (m *Migrate) Down() error {
	return m.migrator.Rollback()
}
