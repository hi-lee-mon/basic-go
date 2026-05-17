package main

import (
	"basic-go/infra"
	"basic-go/src/models"
)

func main() {
	infra.Initialize()
	db := infra.SetupDB()

	if err := db.AutoMigrate(&models.Item{}, &models.User{}); err != nil {
		panic("Failed to migrate database: " + err.Error())
	}
}
