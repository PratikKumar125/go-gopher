package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
    ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Name  string             `bson:"name" json:"name" validate:"required,lte=3"`
    Email string             `bson:"email" json:"email" validate:"required,email"`
}