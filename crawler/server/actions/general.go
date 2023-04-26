package actions

import (
	"strings"

	"github.com/gosimple/slug"
	"gorm.io/gorm/clause"

	"github.com/davidalvarez305/review_poster/crawler/server/database"
	"github.com/davidalvarez305/review_poster/crawler/server/models"
)

func createSubCategories(subCategories []string, category models.Category) ([]models.SubCategory, error) {
	var createdSubcategories []models.SubCategory

	var sc []models.SubCategory
	for _, subcategory := range subCategories {
		sc = append(sc, models.SubCategory{
			Name:       strings.ToLower(subcategory),
			Slug:       slug.Make(subcategory),
			CategoryID: category.ID,
		})
	}

	err := database.DB.Save(&sc).Error

	if err != nil {
		return createdSubcategories, err
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

	err = database.DB.Clauses(clause.OnConflict{DoNothing: true}).Preload("Group").FirstOrCreate(&category).Error

	if err != nil {
		return category, err
	}

	return category, nil
}
