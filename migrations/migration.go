package main

import (
	"basic-go/infra"
	"basic-go/models"
)

func main() {
	infra.Initialize()
	db := infra.SetupDB()

	if err := db.AutoMigrate(&models.Item{}); err != nil {
		panic("Failed to migrate database: " + err.Error())
	}
}
