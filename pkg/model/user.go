package model

type User struct {
	UserID int    `json:"user_id" bson:"user_id"`
	Name   string `json:"name" bson:"name"`
	Token  string `json:"token" bson:"token"`
}
