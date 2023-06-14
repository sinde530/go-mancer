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

	// 사용자 인증 토큰을 추출
	accessToken := c.GetHeader("Authorization")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing access token"})
		return
	}

	// 토큰을 검증하고 사용자 정보를 가져옴
	_, claims, err := token.VerifyToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
		return
	}

	// 사용자 정보를 그룹 데이터에 추가
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
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file"})
		return
	}

	filename := generateUniqueFilename(file.Filename)

	err = c.SaveUploadedFile(file, "uploads/"+filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	presignedURL := generatePresignedURL(filename)

	c.JSON(http.StatusOK, gin.H{"url": presignedURL})
}

func generateUniqueFilename(originalFilename string) string {
	return time.Now().Format("20060102150405") + "_" + originalFilename
}

func generatePresignedURL(filename string) string {
	return "http://localhost:8080/uploads/" + filename
}

func HandleGetGroups(c *gin.Context) {
	groups, err := db.SendGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve groups"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"groups": groups})
}
