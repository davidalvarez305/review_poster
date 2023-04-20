package actions

import (
	"github.com/davidalvarez305/review_poster/cms/server/database"
	"github.com/davidalvarez305/review_poster/cms/server/models"
)

// Create a single record and return the inserted record.
func CreateSynonym(synonym models.Synonym) (models.Synonym, error) {
	var newSynonym models.Synonym

	err := database.DB.Save(&synonym).First(&newSynonym).Error

	if err != nil {
		return newSynonym, err
	}

	return newSynonym, nil
}

// Create multiple records without returning anything.
func CreateSynonyms(synonyms []models.Synonym) error {
	return database.DB.Create(&synonyms).Error
}

// Update a single record that belongs a user.
func UpdateSynonym(synonym models.Synonym, userId, wordId string) error {

	err := database.DB.Where("word_id = ? AND user_id = ?", wordId, userId).Find(&synonym).Error

	if err != nil {
		return err
	}

	err = database.DB.Save(&synonym).First(&synonym).Error

	return err
}

// Updates multiple records and returns DB values.
func UpdateSynonyms(synonyms []models.Synonym, word, userId string) ([]models.Synonym, error) {
	err := database.DB.Save(&synonyms).Error

	if err != nil {
		return synonyms, err
	}

	updatedSynonyms, err := GetSynonymsByWord(word, userId)

	if err != nil {
		return updatedSynonyms, err
	}

	return updatedSynonyms, nil
}

// Select records from DB by a word string.
func GetSynonymsByWord(word, userId string) ([]models.Synonym, error) {
	var synonyms []models.Synonym

	err := database.DB.Where("name = ? ", word).Joins("Word").Find(&synonyms).Error

	if err != nil {
		return synonyms, err
	}

	return synonyms, nil
}

// Delete a slice of records without returning any values.
func DeleteSynonyms(s []int, word, userId string) ([]models.Synonym, error) {
	var synonyms []models.Synonym

	err := database.DB.Delete(&models.Synonym{}, s).Error

	if err != nil {
		return synonyms, err
	}

	synonyms, err = GetSynonymsByWord(word, userId)

	if err != nil {
		return synonyms, err
	}

	return synonyms, nil
}

// Delete without returning anything.
func SimpleDeleteSynonyms(synonyms []models.Synonym) error {
	return database.DB.Delete(&synonyms).Error
}

// Takes structs from the client & deletes them. Does not return records from DB.
func DeleteBulkSynonyms(clientSnonyms, existingSynonyms []models.Synonym) ([]models.Synonym, error) {
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
		err := SimpleDeleteSynonyms(synonyms)

		if err != nil {
			return synonyms, err
		}
	}

	return synonyms, nil
}

// Take structs from client and creates them. Does not return any records.
func AddBulkSynonyms(clientSnonyms, existingSynonyms []models.Synonym) ([]models.Synonym, error) {
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
		err := CreateSynonyms(synonyms)

		if err != nil {
			return synonyms, err
		}
	}

	return synonyms, nil
}
