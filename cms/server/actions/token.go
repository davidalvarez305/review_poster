package actions

import (
	"time"

	"github.com/davidalvarez305/review_poster/cms/server/database"
	"github.com/davidalvarez305/review_poster/cms/server/models"
	"github.com/google/uuid"
)

func GenerateToken() (models.Token, error) {

	// Initialize & Generate Token
	token := models.Token{
		UUID:      uuid.New().String(),
		CreatedAt: time.Now().Unix(),
	}

	err := database.DB.Save(&token).First(&token).Error

	if err != nil {
		return token, err
	}

	return token, nil
}

func GetToken(uuid string, userId int) (models.Token, error) {
	var token models.Token
	err := database.DB.Where("uuid = ? AND user_id = ?", uuid, userId).First(&token).Error

	if err != nil {
		return token, err
	}

	return token, nil
}

func DeleteToken(token models.Token) error {
	return database.DB.Delete(&token).Error
}
