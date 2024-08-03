package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (u *User) MarshalBSON() ([]byte, error) {
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	u.UpdatedAt = time.Now()
	type userAlias User

	return bson.Marshal((*userAlias)(u))
}

// User model
type User struct {
	ID        string    `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName string    `json:"firstName" bson:"first_name"`
	LastName  string    `json:"lastName" bson:"last_name"`
	UserName  string    `json:"userName" bson:"user_name"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"password" bson:"password"`
	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}
