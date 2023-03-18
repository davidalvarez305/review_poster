package actions

import (
	"fmt"

	"github.com/gosimple/slug"

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
	err := database.DB.Where("name = ?", groupName).First(&g).Error

	if err != nil {
		fmt.Printf("Creating new group...")
	}

	g.Group = &models.Group{
		Name: groupName,
		Slug: slug.Make(groupName),
	}

	return nil
}

func (s *SubCategory) GetOrCreateSubCategory(categoryName, subCategoryName, groupName string) error {

	err := database.DB.Where("name = ?", subCategoryName).Preload("SubCategory.Category.Group").First(&s).Error

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
		Name: subCategoryName,
		Slug: slug.Make(subCategoryName),
		Category: &models.Category{
			ID:   category.ID,
			Name: category.Name,
			Slug: category.Slug,
			Group: &models.Group{
				ID:   group.ID,
				Name: group.Name,
				Slug: group.Slug,
			},
		},
	}

	return nil
}

func (c *Category) GetOrCreateCategory(categoryName string, group *Group) error {
	err := database.DB.Where("name = ?", categoryName).Preload("Category.Group").First(&c).Error

	if err != nil {
		fmt.Printf("Creating new category...\n")
	}

	c.Category = &models.Category{
		Name: categoryName,
		Slug: slug.Make(categoryName),
		Group: &models.Group{
			ID:   group.ID,
			Name: group.Name,
			Slug: group.Slug,
		},
	}

	return nil
}
