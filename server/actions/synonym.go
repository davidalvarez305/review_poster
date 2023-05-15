package actions

import (
	"github.com/davidalvarez305/review_poster/server/database"
	"github.com/davidalvarez305/review_poster/server/models"
)

// Select records from DB by a word string.
func GetSynonymsByWord(word, userId string) ([]models.Synonym, error) {
	var synonyms []models.Synonym

	err := database.DB.Preload("Word").Joins("INNER JOIN word ON word.id = synonym.word_id").Where("word.name = ? AND word.user_id = ?", word, userId).Find(&synonyms).Error

	if err != nil {
		return synonyms, err
	}

	return synonyms, nil
}

// Takes structs from the client & deletes them. Does not return records from DB.
func DeleteBulkSynonyms(clientSnonyms, existingSynonyms []models.Synonym) error {
	var synonyms []models.Synonym

	for _, existingSynonym := range existingSynonyms {
		found := false
		for _, clientSnonym := range clientSnonyms {
			if existingSynonym.Synonym == clientSnonym.Synonym {
				found = true
			}
		}
		if !found {
			synonyms = append(synonyms, existingSynonym)
		}
	}

	if len(synonyms) > 0 {
		err := database.DB.Delete(&synonyms).Error

		if err != nil {
			return err
		}
	}

	return nil
}

// Take structs from client and creates them. Does not return any records.
func AddBulkSynonyms(clientSnonyms, existingSynonyms []models.Synonym) error {
	var synonyms []models.Synonym

	for _, clientSnonym := range clientSnonyms {
		found := false
		for _, existingSynonym := range existingSynonyms {
			if clientSnonym.Synonym == existingSynonym.Synonym {
				found = true
			}
		}
		if !found {
			synonyms = append(synonyms, clientSnonym)
		}
	}

	if len(synonyms) > 0 {
		err := database.DB.Create(&synonyms).Error

		if err != nil {
			return err
		}
	}

	return nil
}
