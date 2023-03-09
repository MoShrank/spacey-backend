package entity

import "time"

type Tag struct {
	ID        string     `bson:"_id,omitempty"`
	Name      string     `bson:"name"`
	SubTags   []Tag      `bson:"sub_tags"`
	CreatedAt *time.Time `bson:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at"`
}

type TagRes struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	SubTags []TagRes `json:"sub_tags"`
}

type TagUseCaseInterface interface {
}

type TagStoreInterface interface {
	CreateTag(tag *Tag) (string, error)
	FindAll() ([]Tag, error)
	FindByID(tagID string) (*Tag, error)
	Delete(tagID string) error
}
