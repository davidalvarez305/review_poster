package actions

import (
	"time"

	"github.com/davidalvarez305/review_poster/cms/server/models"
	"github.com/davidalvarez305/review_poster/cms/server/server"
	"github.com/google/uuid"
)

func GenerateToken() (models.Token, error) {

	// Initialize & Generate Token
	token := models.Token{
		UUID:      uuid.New().String(),
		CreatedAt: time.Now().Unix(),
	}

	err := server.DB.Save(&token).First(&token).Error

	if err != nil {
		return token, err
	}

	return token, nil
}

func GetToken(uuid string, userId int) (models.Token, error) {
	var token models.Token
	err := server.DB.Where("uuid = ? AND user_id = ?", uuid, userId).First(&token).Error

	if err != nil {
		return token, err
	}

	return token, nil
}

func DeleteToken(token models.Token) error {
	return server.DB.Delete(&token).Error
}
