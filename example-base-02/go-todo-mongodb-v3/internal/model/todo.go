package model

import (
 "time"
 "go.mongodb.org/mongo-driver/v2/bson"
)

type Todo struct {
 ID bson.ObjectID `bson:"_id,omitempty" json:"id"`
 Title string `bson:"title" json:"title"`
 Description string `bson:"description" json:"description"`
 Completed bool `bson:"completed" json:"completed"`
 CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
 UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}
type Input struct { Title string `json:"title"`; Description string `json:"description"`; Completed bool `json:"completed"` }
type Patch struct { Title *string `json:"title"`; Description *string `json:"description"`; Completed *bool `json:"completed"` }
