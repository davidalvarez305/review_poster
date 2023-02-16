package models

type CategoryGroup struct {
	ID          int
	Name        string
	ParentGroup ParentGroup `gorm:"column:parent_group_id"`
}
