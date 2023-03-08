package actions

import (
	"github.com/davidalvarez305/review_poster/cms/server/database"
	"github.com/davidalvarez305/review_poster/cms/server/models"
)

type Synonym struct {
	*models.Synonym
}

type Synonyms []*models.Synonym

// Create a single record and return the inserted record.
func (synonym *Synonym) CreateSynonym() error {
	query := database.DB.Save(&synonym).First(&synonym)
	return query.Error
}

// Create multiple records without returning anything.
func (synonyms *Synonyms) CreateSynonyms() error {
	query := database.DB.Create(&synonyms)
	return query.Error
}

// Update a single record that belongs a user.
func (synonym *Synonym) UpdateSynonym(userId, wordId string) error {

	err := database.DB.Where("word_id = ? AND user_id = ?", wordId, userId).Find(&synonym).Error

	if err != nil {
		return err
	}

	err = database.DB.Save(&synonym).First(&synonym).Error

	return err
}

// Updates multiple records and returns DB values.
func (synonyms *Synonyms) UpdateSynonyms(word, userId string) error {
	err := database.DB.Save(&synonyms).Error

	if err != nil {
		return err
	}

	err = synonyms.GetSynonymsByWord(word, userId)
	return err
}

// Select records from DB by a word string.
func (synonyms *Synonyms) GetSynonymsByWord(word, userId string) error {
	return database.DB.Where("user_id = ?", userId).Preload("Word", "name = ?", word).Find(*synonyms).Error
}

// Delete a slice of records without returning any values.
func (synonyms *Synonyms) DeleteSynonyms(s []int, word, userId string) error {
	err := database.DB.Delete(&models.Synonym{}, s).Error

	if err != nil {
		return err
	}

	return synonyms.GetSynonymsByWord(word, userId)
}

// Delete without returning anything.
func (synonyms *Synonyms) SimpleDelete() error {
	res := database.DB.Delete(&synonyms)
	return res.Error
}

// Takes structs from the client & deletes them. Does not return records from DB.
func (synonyms Synonyms) DeleteBulkSynonyms(clientSnonyms, existingSynonyms *Synonyms) error {
	for _, existingSynonym := range *existingSynonyms {
		found := false
		for _, clientSnonym := range *clientSnonyms {
			if existingSynonym.Synonym == clientSnonym.Synonym {
				found = true
			}
		}
		if !found {
			synonyms = append(synonyms, existingSynonym)
		}
	}

	if len(synonyms) > 0 {
		err := synonyms.SimpleDelete()

		if err != nil {
			return err
		}
	}

	return nil
}

// Take structs from client and creates them. Does not return any records.
func (synonyms Synonyms) AddBulkSynonyms(clientSnonyms, existingSynonyms *Synonyms) error {
	for _, clientSnonym := range *clientSnonyms {
		found := false
		for _, existingSynonym := range *existingSynonyms {
			if clientSnonym.Synonym == existingSynonym.Synonym {
				found = true
			}
		}
		if !found {
			synonyms = append(synonyms, clientSnonym)
		}
	}

	if len(synonyms) > 0 {
		err := synonyms.CreateSynonyms()

		if err != nil {
			return err
		}
	}

	return nil
}
