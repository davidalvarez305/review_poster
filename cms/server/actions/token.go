package actions

import (
	"time"

	"github.com/davidalvarez305/review_poster/cms/server/database"
	"github.com/davidalvarez305/review_poster/cms/server/models"
	"github.com/google/uuid"
)

type Token struct {
	*models.Token
}

func (token *Token) GenerateToken() error {

	// Initialize & Generate Token
	token.Token = &models.Token{
		UUID:      uuid.New().String(),
		CreatedAt: time.Now().Unix(),
	}

	return database.DB.Save(&token).First(&token).Error
}

func (token *Token) GetToken(uuid string, userId int) error {
	result := database.DB.Where("uuid = ? AND user_id = ?", uuid, userId).First(&token)
	return result.Error
}

func (token *Token) DeleteToken() error {
	result := database.DB.Delete(&token)
	return result.Error
}
