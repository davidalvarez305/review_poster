package actions

import (
	"time"

	"github.com/davidalvarez305/content_go/server/database"
	"github.com/davidalvarez305/content_go/server/models"
	"github.com/google/uuid"
)

type Token struct {
	*models.Token
}

func (token *Token) GenerateToken(user *Users) error {

	// Create UUID for Token
	uuid := uuid.New().String()

	// Initialize & Generate Token
	t := models.Token{
		UUID:      uuid,
		UserID:    user.ID,
		CreatedAt: time.Now().Unix(),
	}

	// Assign token to struct
	token.Token = &t

	result := database.DB.Save(&token).First(&token)

	return result.Error
}

func (token *Token) GetToken(uuid string, userId uint) error {
	result := database.DB.Where("uuid = ? AND user_id = ?", uuid, userId).First(&token)
	return result.Error
}

func (token *Token) DeleteToken() error {
	result := database.DB.Delete(&token)
	return result.Error
}
