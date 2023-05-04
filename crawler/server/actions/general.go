package actions

import (
	"fmt"
	"strings"

	"github.com/gosimple/slug"
	"gorm.io/gorm/clause"

	"github.com/davidalvarez305/review_poster/crawler/server/database"
	"github.com/davidalvarez305/review_poster/crawler/server/models"
)

func createSubCategories(subCategories []string, category models.Category) ([]models.SubCategory, error) {
	var createdSubcategories []models.SubCategory
	var existingSubcategories []models.SubCategory

	// Ensure that that subcategories that will be inserted don't already exist -- otherwise the entire category will not be inserted
	err := database.DB.Find(&existingSubcategories).Error

	if err != nil {
		fmt.Printf("ERROR FING EXISTING SUBCATEGORIES: %+v", err)
	}

	var subCategoriesToCreate []models.SubCategory
	for _, subcategory := range subCategories {
		var lowerCaseSubCategory = strings.ToLower(subcategory)
		var slugSubCategory = slug.Make(subcategory)

		exists := false

		for _, existingSubCategory := range existingSubcategories {
			if existingSubCategory.Name == lowerCaseSubCategory {
				exists = true
				break
			}
		}

		// subcategory will be omitted if it exists already, avoiding that this function errors out for existing entities
		if exists {
			break
		}

		subCategoriesToCreate = append(subCategoriesToCreate, models.SubCategory{
			Name:       lowerCaseSubCategory,
			Slug:       slugSubCategory,
			CategoryID: category.ID,
		})
	}

	if len(subCategoriesToCreate) > 0 {
		err = database.DB.Clauses(clause.OnConflict{UpdateAll: true}).Save(&subCategoriesToCreate).Error

		if err != nil {
			return createdSubcategories, err
		}
	}

	err = database.DB.Preload("Category.Group").Find(&createdSubcategories).Error

	return createdSubcategories, err
}

func createOrFindCategory(categoryName, groupName string) (models.Category, error) {
	var group models.Group
	var category models.Category

	err := database.DB.Where("name = ?", groupName).First(&group).Error

	if err != nil {
		return category, err
	}

	category.Name = strings.ToLower(categoryName)
	category.Slug = slug.Make(categoryName)
	category.GroupID = group.ID

	err = database.DB.Where("name = ?", category.Name).Preload("Group").First(&category).Error

	if err != nil {
		err = database.DB.Save(&category).Preload("Group").First(&category).Error

		if err != nil {
			return category, err
		}
	}

	return category, nil
}
