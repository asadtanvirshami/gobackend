package initializers

import "your-app/models"

func SyncDB () {	
	DB.AutoMigrate(&models.User{})
}