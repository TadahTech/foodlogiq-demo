package model

type Event struct {
	CreatedAt string `json:"created_at" bson:"created_at"`
	IsDeleted bool   `json:"is_deleted" bson:"is_deleted"`
	CreatedBy string `json:"created_by" bson:"created_by"`
	Contents  []struct {
		Lot            string `json:"lot" bson:"lot"`
		Gtin           string `json:"gtin" bson:"gtin"`
		BestByDate     string `json:"best_by_date" bson:"best_by_date"`
		ExpirationDate string `json:"expiration_date" bson:"expiration_date"`
	} `json:"contents" bson:"contents"`
	ID   string `json:"id" bson:"id"`
	Type string `json:"type" bson:"type"`
}
