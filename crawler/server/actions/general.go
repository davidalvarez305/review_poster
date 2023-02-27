package actions

import (
	"github.com/davidalvarez305/review_poster/crawler/server/database"
	"github.com/davidalvarez305/review_poster/crawler/server/models"
)

type Group struct {
	*models.Group
}

type SubCategory struct {
	*models.SubCategory
}

func (g *Group) GetOrCreateGroup(name string) error {
	err := database.DB.Where("name = ?", name).First(&g).Error

	if err != nil {
		fmt.Printf("Creating new group...")
	}

	g = &models.Group{
		Name: groupName,
		Slug: groupName,
	}

	return nil
}

func (s *SubCategory) GetOrCreateSubCategory(categoryName, subCategoryName, groupName string) error {
	
	err := database.DB.Where("name = ?", subCategoryName).Preload("SubCategory.Category.Group").Error

	// If there is no error, it means that the subcategory was found, so we can return early.
	if err == nil {
		return err
	}

	group := Group{}
	category := Category{}

	err = group.GetOrCreateGroup(groupName)

	if err != nil {
		return err
	}

	err = category.GetOrCreateCategory(categoryName)

	if err != nil {
		return err
	}

	// create GetOrCreateCategory function
	// create slugify and toTitle functions

	s = &models.SubCategory{
		Name: subCategoryName,
		Slug: subCategorySlug,
		Category: &models.Category{
			ID:    category.ID,
			Name:  category.Name,
			Slug:  category.Slug,
			Group: group.Group,
		},
	}

	return nil
}
