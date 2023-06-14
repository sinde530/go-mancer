package model

type Group struct {
	ID                string   `json:"group_id" bson:"_id,omitempty"`
	CreatedByUID      string   `json:"created_by_uid"`
	CreatedByUsername string   `json:"created_by_username"`
	Thumbnail         string   `json:"thumbnail"`
	Title             string   `json:"title"`
	Members           []string `json:"members"`
	CreatedAt         string   `json:"created_at"`
	UpdatedAt         string   `json:"updated_at"`
}
