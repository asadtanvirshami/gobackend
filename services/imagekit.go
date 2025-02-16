package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/imagekit-developer/imagekit-go"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
)

var ik *imagekit.ImageKit

// Initialize ImageKit in init()
func init() {
	privateKey := "private_BkdC3TtLFKR7bdHuaCgR6T565WQ="
	publicKey := "public_JCxCGw8zqXg0IfqdHmyRNCbb2HM="
	endpoint := "https://ik.imagekit.io/epbtkdzri1"

	fmt.Println("Private Key:", privateKey)
	fmt.Println("Public Key:", publicKey)
	fmt.Println("Endpoint URL:", endpoint)

	if privateKey == "" || publicKey == "" || endpoint == "" {
		fmt.Println("❌ Missing ImageKit environment variables")
	}

	ik = imagekit.NewFromParams(imagekit.NewParams{
		PrivateKey:  privateKey,
		PublicKey:   publicKey,
		UrlEndpoint: endpoint,
	})

	if ik == nil {
		fmt.Println("Failed to initialize ImageKit")
	}
}

// UploadImage uploads a Base64-encoded image to ImageKit and returns the URL
func UploadImage(imagePath, fileName, folder string) (string, string, error) {
	// Read file bytes
	fileData, err := os.ReadFile(imagePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read file: %v", err)
	}

	// Convert to Base64
	base64Data := base64.StdEncoding.EncodeToString(fileData)
	base64Image := fmt.Sprintf("data:image/jpeg;base64,%s", base64Data)

	// Create a context for API request
	ctx := context.Background()

	// Upload file to ImageKit
	resp, err := ik.Uploader.Upload(ctx, base64Image, uploader.UploadParam{
		FileName: fileName,
		Folder:   folder, // Ensure your folder exists in ImageKit
	})

	if err != nil {
		return "", "", fmt.Errorf("failed to upload file: %v", err)
	}

	// Now we can safely delete the file
	defer func() {
		if err := os.Remove(imagePath); err != nil {
			fmt.Println("Error removing file:", err)
		} else {
			fmt.Println("File successfully removed:", imagePath)
		}
	}()

	// Return uploaded file URL and file ID
	return resp.Data.Url, resp.Data.FileId, nil
}
