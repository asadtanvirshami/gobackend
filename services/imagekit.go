package services

import (
	"context"
	"fmt"
	"os"

	"github.com/imagekit-developer/imagekit-go"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
)

// Global ImageKit client instance
var ik *imagekit.ImageKit

// Initialize ImageKit in init()
func init() {
	ik = imagekit.NewFromParams(imagekit.NewParams{
		PrivateKey:  os.Getenv("IMAGEKIT_PRIVATE_KEY"),
		PublicKey:   os.Getenv("IMAGEKIT_PUBLIC_KEY"),
		UrlEndpoint: os.Getenv("IMAGEKIT_ENDPOINT_URL"),
	})

	if ik == nil {
		fmt.Println("Failed to initialize ImageKit")
	}
}

// UploadToImageKit uploads a file to ImageKit and returns the URL
func UploadToImageKit(filePath string) (string, error) {
	// Read file as bytes
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// Ensure ImageKit is initialized
	if ik == nil {
		return "", fmt.Errorf("ImageKit client is not initialized")
	}

	// Upload file to ImageKit
	resp, err := ik.Uploader.Upload(context.Background(), fileData, uploader.UploadParam{
		FileName: "myimage.jpg",
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload to ImageKit: %v", err)
	}

	return resp.Data.Url, nil
}
