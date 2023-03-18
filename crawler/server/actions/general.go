package actions

import (
	"github.com/gosimple/slug"
	"gorm.io/gorm/clause"

	"github.com/davidalvarez305/review_poster/crawler/server/database"
	"github.com/davidalvarez305/review_poster/crawler/server/models"
)

type Group struct {
	*models.Group
}

type SubCategory struct {
	*models.SubCategory
}

type Category struct {
	*models.Category
}

func (g *Group) GetOrCreateGroup(groupName string) error {
	return database.DB.Where("name = ?", groupName).First(&g).Error
}

func (s *SubCategory) GetOrCreateSubCategory(categoryName, subCategoryName, groupName string) error {

	err := database.DB.Where("name = ?", subCategoryName).Preload("Category.Group").First(&s).Error

	// If there is no error, it means that the subcategory was found, so we can return early.
	if err == nil {
		return nil
	}

	group := Group{}
	category := Category{}

	err = group.GetOrCreateGroup(groupName)

	if err != nil {
		return err
	}

	err = category.GetOrCreateCategory(categoryName, &group)

	if err != nil {
		return err
	}

	s.SubCategory = &models.SubCategory{
		Name:     subCategoryName,
		Slug:     slug.Make(subCategoryName),
		Category: category.Category,
	}

	return nil
}

func (c *Category) GetOrCreateCategory(categoryName string, group *Group) error {

	c.Category = &models.Category{
		Name:    categoryName,
		Slug:    slug.Make(categoryName),
		GroupID: group.ID,
	}

	err := database.DB.Clauses(clause.OnConflict{DoNothing: true}).Save(&c).Error

	if err != nil {
		return err
	}

	return database.DB.Preload("Group").First(&c).Error
}
