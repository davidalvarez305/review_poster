package actions

import (
	"fmt"
	"strings"

	"github.com/gosimple/slug"
	"gorm.io/gorm/clause"

	"github.com/davidalvarez305/review_poster/crawler/server/database"
	"github.com/davidalvarez305/review_poster/crawler/server/models"
)

func newGroup(groupName string) (models.Group, error) {
	group := models.Group{
		Name: strings.ToLower(groupName),
		Slug: slug.Make(groupName),
	}

	err := database.DB.Clauses(clause.OnConflict{DoNothing: true}).FirstOrCreate(&group).Error

	if err != nil {
		return group, err
	}

	return group, database.DB.Where("name = ?", group.Name).First(&group).Error
}

func newSubCategory(categoryName, subCategoryName, groupName string) (models.SubCategory, error) {
	var subCategory models.SubCategory

	err := database.DB.Where("name = ?", subCategoryName).Preload("Category.Group").First(&subCategory).Error

	// If there is no error, it means that the subcategory was found, so we can return early.
	if err == nil {
		return subCategory, nil
	} else {
		fmt.Printf("CREATING NEW SUB_CATEGORY...%+v\n", err)
	}

	group, err := newGroup(groupName)

	if err != nil {
		return subCategory, err
	}

	category, err := newCategory(categoryName, group)

	if err != nil {
		return subCategory, err
	}

	subCategory.Name = strings.ToLower(subCategoryName)
	subCategory.Slug = slug.Make(subCategoryName)
	subCategory.CategoryID = category.ID

	err = database.DB.Save(&subCategory).Error

	if err != nil {
		return subCategory, err
	}

	err = database.DB.Preload("Category.Group").First(&subCategory).Error

	return subCategory, err
}

func newCategory(categoryName string, group models.Group) (models.Category, error) {
	var category models.Category

	category.Name = strings.ToLower(categoryName)
	category.Slug = slug.Make(categoryName)
	category.GroupID = group.ID

	err := database.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&category).Error

	if err != nil {
		return category, err
	}

	err = database.DB.Where("slug = ?", category.Slug).Preload("Group").First(&category).Error

	return category, err
}
