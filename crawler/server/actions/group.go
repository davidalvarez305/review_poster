package actions

import (
	"github.com/davidalvarez305/review_poster/crawler/server/database"
	"github.com/davidalvarez305/review_poster/crawler/server/models"
)

type Group struct {
	*models.Group
}

func (g *Group) GetGroup(name string) error {
	return database.DB.Where("name = ?", name).First(&g).Error
}
