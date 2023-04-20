package actions

import (
	"fmt"
	"strings"

	"github.com/gosimple/slug"
	"gorm.io/gorm/clause"

	"github.com/davidalvarez305/review_poster/crawler/server/database"
	"github.com/davidalvarez305/review_poster/crawler/server/models"
)

type Group struct {
	*models.Group
}

func (g *Group) newGroup(groupName string) error {
	group := models.Group{
		Name: strings.ToLower(groupName),
		Slug: slug.Make(groupName),
	}

	err := database.DB.Clauses(clause.OnConflict{DoNothing: true}).FirstOrCreate(&group).Error

	if err != nil {
		return err
	}

	return database.DB.Where("name = ?", group.Name).First(&g).Error
}

func newSubCategory(categoryName, subCategoryName, groupName string) (*models.SubCategory, error) {
	var s models.SubCategory

	err := database.DB.Where("name = ?", subCategoryName).Preload("Category.Group").First(&s).Error

	// If there is no error, it means that the subcategory was found, so we can return early.
	if err == nil {
		return &s, nil
	} else {
		fmt.Printf("CREATING NEW SUB_CATEGORY...%+v\n", err)
	}

	group := Group{}
	err = group.newGroup(groupName)

	if err != nil {
		return &s, err
	}

	category, err := newCategory(categoryName, &group)

	if err != nil {
		return &s, err
	}

	s.Name = strings.ToLower(subCategoryName)
	s.Slug = slug.Make(subCategoryName)
	s.CategoryID = category.ID

	err = database.DB.Save(&s).Error

	if err != nil {
		return &s, err
	}

	err = database.DB.Preload("Category.Group").First(&s).Error

	return &s, err
}

func newCategory(categoryName string, group *Group) (*models.Category, error) {
	var c models.Category

	c.Name = strings.ToLower(categoryName)
	c.Slug = slug.Make(categoryName)
	c.GroupID = group.ID

	err := database.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&c).Error

	if err != nil {
		return &c, err
	}

	err = database.DB.Where("slug = ?", c.Slug).Preload("Group").First(&c).Error

	return &c, err
}
