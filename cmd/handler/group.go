package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sinde530/go-mancer/cmd/token"
	"github.com/sinde530/go-mancer/db"
	"github.com/sinde530/go-mancer/model"
)

func HandleCreateGroup(c *gin.Context) {
	var group model.Group

	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse request"})
		return
	}

	// 사용자 인증 토큰을 추출합니다.
	accessToken := c.GetHeader("Authorization")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access token"})
		return
	}

	// 토큰을 검증하고 사용자 정보를 가져옵니다.
	_, claims, err := token.VerifyToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
		return
	}

	// 사용자 정보를 그룹 데이터에 추가합니다.
	user := claims.User
	group.CreatedByUID = user.UID
	group.CreatedByUsername = user.Username
	group.ID = "" // You should generate the group ID here

	group.Members = []string{}
	group.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	group.UpdatedAt = group.CreatedAt

	err = db.SaveGroup(&group)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Failed to save group"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"group": group})
}

func HandleUploadURL(c *gin.Context) {
	// Get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file"})
		return
	}

	// Generate a unique filename for the uploaded file
	filename := generateUniqueFilename(file.Filename)

	// Save the file to your desired location
	// For example, you can save it to the "uploads" directory
	err = c.SaveUploadedFile(file, "uploads/"+filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Create the presigned URL for the uploaded file
	// Replace the URL below with your own logic to generate the presigned URL
	presignedURL := generatePresignedURL(filename)

	// Return the presigned URL to the client
	c.JSON(http.StatusOK, gin.H{"url": presignedURL})
}

func generateUniqueFilename(originalFilename string) string {
	// Implement your own logic to generate a unique filename
	// For example, you can use a combination of timestamp and a random string
	// Make sure the generated filename does not collide with existing files
	return time.Now().Format("20060102150405") + "_" + originalFilename
}

func generatePresignedURL(filename string) string {
	// Implement your own logic to generate the presigned URL
	// This URL should be accessible to download the uploaded file
	// For example, you can use a cloud storage service like AWS S3 and generate a presigned URL
	// Make sure the generated URL has appropriate access permissions and expires after a certain time
	return "http://localhost:8080/uploads/" + filename
}
