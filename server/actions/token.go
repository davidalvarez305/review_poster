package actions

import (
	"time"

	"github.com/davidalvarez305/review_poster/server/database"
	"github.com/davidalvarez305/review_poster/server/models"
	"github.com/google/uuid"
)

func generateToken(userId int) (models.Token, error) {

	// Initialize & Generate Token
	token := models.Token{
		UUID:      uuid.New().String(),
		CreatedAt: time.Now().Unix(),
		UserID:    userId,
	}

	err := database.DB.Save(&token).First(&token).Error

	if err != nil {
		return token, err
	}

	return token, nil
}

func GetToken(uuid string) (models.Token, error) {
	var token models.Token
	err := database.DB.Where("uuid = ?", uuid).First(&token).Error

	if err != nil {
		return token, err
	}

	return token, nil
}

func DeleteToken(token models.Token) error {
	return database.DB.Delete(&token).Error
}
