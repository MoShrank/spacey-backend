package models

type Deck struct {
	ID        string `json:"id,omitempty" bson:"id"`
	Name      string `json:"name"         bson:"name"       validate:"required"`
	UserID    string `json:"-"            bson:"user_id"`
	CreatedAt string `json:"-"            bson:"created_at"`
	UpdatedAt string `json:"-"            bson:"updated_at"`
	DeletedAt string `json:"-"            bson:"deleted_at"`
}
