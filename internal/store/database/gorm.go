package database

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/guilherme-or/go-api-template/internal/store/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func getLogger(debug bool) logger.Interface {
	if debug {
		return logger.NewSlogLogger(slog.Default(), logger.Config{
			LogLevel: logger.Info,
		})
	}

	return nil
}

func NewGORMPostgresConnection(dsn string, debug bool) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: getLogger(debug),
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		return db, err
	}

	db.Create(&models.User{
		Default: models.Default{
			ID: uuid.New(),
		},
		Email:    "guioliroc@gmail.com",
		Name:     "Guilherme Rocha",
		Password: "$2a$12$ASolmc9y3MtypLAlQ.wAOeXBiW5BYYXtsDGFNTExN6Cw8PKhldNyO", // G29aj~@#
		Role:     models.RoleAdmin,
	})

	db.Create(&models.User{
		Default: models.Default{
			ID: uuid.New(),
		},
		Email:    "tester@hotmail.com",
		Name:     "Tester McTestface",
		Password: "$2a$12$ASolmc9y3MtypLAlQ.wAOeXBiW5BYYXtsDGFNTExN6Cw8PKhldNyO", // G29aj~@#
		Role:     models.RoleUser,
	})

	return db, nil
}
