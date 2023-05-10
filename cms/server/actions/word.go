package actions

import (
	"github.com/davidalvarez305/review_poster/cms/server/database"
	"github.com/davidalvarez305/review_poster/cms/server/models"
)

func DeleteWord(word models.Word, userId, word_id string) error {
	return database.DB.Where("user_id = ? AND id = ?", userId, word_id).Delete(&word).Error
}
