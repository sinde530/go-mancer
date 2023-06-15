package handler

import (
	"net/http"
	"strconv"
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

	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	// group.CreatedByUID = user.UID
	group.CreatedByUID = user.UID + "-" + timestamp
	group.CreatedByUsername = user.Username
	group.ID = "" // You should generate the group ID here

	group.Members = []string{user.Email}
	group.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	group.UpdatedAt = group.CreatedAt

	err = db.SaveGroup(&group)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Failed to save group"})
		return
	}

	userFromDB, err := db.GetUserByEmail(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrive user"})
		return
	}

	userFromDB.Groups = append(userFromDB.Groups, group.CreatedByUID)
	err = db.UpdateUser(userFromDB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
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

// func HandleGetGroups(c *gin.Context) {
// 	groups, err := db.SendGroups()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve groups"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"groups": groups})
// }

func HandleGetGroups(c *gin.Context) {
	// 토큰에서 사용자 정보를 가져옴.
	_, claims, err := token.VerifyToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
		return
	}

	user := claims.User

	// 사용자의 이메일로 DB에서 사용자를 찾음.
	userFromDB, err := db.GetUserByEmail(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrive user"})
		return
	}

	// 사용자가 속한 그룹의 ID를 가져옴.
	userGroups := userFromDB.Groups

	// DB에서 모든 그룹을 가져옴.
	groups, err := db.SendGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve groups"})
		return
	}

	// 사용자가 속한 그룹만 필터링.
	var filteredGroups []*model.Group
	for _, group := range groups {
		if contains(userGroups, group.CreatedByUID) {
			filteredGroups = append(filteredGroups, group)
		}
	}

	c.JSON(http.StatusOK, gin.H{"groups": filteredGroups})
}

// contains 함수는 주어진 문자열 슬라이스에 특정 문자열이 포함되어 있는지 확인.
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
