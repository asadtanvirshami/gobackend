package initializers

import "your-app/models"

func SyncDB() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Category{})
	DB.AutoMigrate(&models.Community{})
	DB.AutoMigrate(&models.CommunityEvent{})
	DB.AutoMigrate(&models.CommunityCourse{})
	DB.AutoMigrate(&models.CommunityLeaderboard{})
	DB.AutoMigrate(&models.CommunityPoll{})
	DB.AutoMigrate(&models.CommunityPollOption{})
	DB.AutoMigrate(&models.CommunityUser{})
	DB.AutoMigrate(&models.CommunityPost{})
}
