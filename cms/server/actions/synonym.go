package actions

import (
	"fmt"

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
func (synonym *Synonym) UpdateSynonym(userId string) error {

	// Store what came from the client separately to first make sure that the synonym being updated exists
	syn := synonym

	sql := fmt.Sprintf(`SELECT * FROM synonym WHERE word_id = (
		SELECT id FROM word WHERE name = '%v' AND user_id = '%s'
	)`, synonym.WordID, userId)

	query := database.DB.Raw(sql).Scan(&synonym)

	if query.Error != nil {
		return query.Error
	}

	// After checking that the record exists & belongs to the client user, re-assign and save.
	synonym = syn
	query = database.DB.Save(&synonym).First(&synonym)

	return query.Error
}

// Updates multiple records and returns DB values.
func (synonyms *Synonyms) UpdateSynonyms(word, userId string) error {
	query := database.DB.Save(&synonyms)

	if query.Error != nil {
		return query.Error
	}

	err := synonyms.GetSynonymsByWord(word, userId)
	return err
}

// Select records from DB by a word string.
func (synonyms *Synonyms) GetSynonymsByWord(word, userId string) error {

	sql := fmt.Sprintf(`SELECT * FROM synonym WHERE word_id = (
		SELECT id FROM word WHERE name = '%s' AND user_id = '%s'
	)`, word, userId)

	result := database.DB.Raw(sql).Scan(&synonyms)

	return result.Error
}

// Delete a slice of records without returning any values.
func (synonyms *Synonyms) DeleteSynonyms(s []int, word, userId string) error {
	res := database.DB.Delete(&models.Synonym{}, s)

	if res.Error != nil {
		return res.Error
	}

	sql := fmt.Sprintf(`SELECT * FROM synonym WHERE word_id = (
		SELECT id FROM word WHERE name = '%s' AND user_id = '%s'
	)`, word, userId)

	result := database.DB.Raw(sql).Scan(&synonyms)

	return result.Error
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
