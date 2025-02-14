package config

import (
	"os"
)

var (
	SMTP_HOST     = os.Getenv("SMTP_HOST")
	SMTP_PORT     = os.Getenv("SMTP_PORT")
	SMTP_USER     = os.Getenv("SMTP_USER")
	SMTP_PASSWORD = os.Getenv("SMTP_PASSWORD")
	IMAGEKIT_URL  = os.Getenv("IMAGEKIT_URL")
	IMAGEKIT_KEY  = os.Getenv("IMAGEKIT_KEY")
)
