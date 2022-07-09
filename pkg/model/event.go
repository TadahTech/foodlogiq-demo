package model

type Event struct {
	CreatedAt string     `json:"created_at,omitempty"`
	IsDeleted bool       `json:"is_deleted,omitempty"`
	CreatedBy int        `json:"created_by,omitempty"`
	Contents  []*Content `json:"contents"`
	ID        string     `json:"id,omitempty"`
	Type      string     `json:"type"`
}

type Content struct {
	Lot            string `json:"lot"`
	Gtin           string `json:"gtin"`
	BestByDate     string `json:"best_by_date,omitempty"`
	ExpirationDate string `json:"expiration_date,omitempty"`
}
