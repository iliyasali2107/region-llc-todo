package models

import "time"

type Todo struct {
	Id       string    `bson:"_id,omitempty"`
	Title    string    `bson:"title"`
	ActiveAt time.Time `bson:"active_at"`
}
