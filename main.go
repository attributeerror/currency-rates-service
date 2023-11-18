package main

import (
	"fmt"
	"os"
	"time"

	"github.com/attributeerror/currency-rates-service/database"
	"github.com/attributeerror/currency-rates-service/handlers/convert_handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := loadDotEnv(); err != nil {
		panic(fmt.Errorf("error whilst loading .env file: %w", err))
	}

	db, err := initDatabase()
	if err != nil {
		panic(fmt.Errorf("error whilst initialising database: %w", err))
	}

	engine := gin.New()
	engine.Use(gin.Recovery())

	convert_handlers.InitialiseRoutes(engine, db)

	port, _ := loadenvvar("PORT", false)
	if port == nil {
		engine.Run(":80")
	} else {
		engine.Run(fmt.Sprintf(":%s", *port))
	}
}

func loadDotEnv() error {
	err := godotenv.Load()
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		return nil
	}

	return nil
}

func initDatabase() (database.TursoDatabase, error) {
	tursoPrimaryUrl, err := loadenvvar("TURSO_URL", true)
	if err != nil {
		return nil, err
	}
	tursoAuthToken, err := loadenvvar("TURSO_AUTH_TOKEN", true)
	if err != nil {
		return nil, err
	}
	tursoDbName, err := loadenvvar("TURSO_DB_NAME", true)
	if err != nil {
		return nil, err
	}

	db, err := database.InitTursoDatabase(
		database.WithPrimaryUrl(*tursoPrimaryUrl),
		database.WithAuthToken(*tursoAuthToken),
		database.WithDbName(*tursoDbName),
		database.WithSyncInterval(1*time.Minute),
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func loadenvvar(key string, required bool) (*string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return &value, nil
	} else if required {
		return nil, fmt.Errorf("required env var not set: %s", key)
	}

	return nil, nil
}
