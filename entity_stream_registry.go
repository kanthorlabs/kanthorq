package kanthorq

type StreamRegistry struct {
	Id        string `json:"id" validate:"required"`
	Name      string `json:"name" validate:"required,is_collection_name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
