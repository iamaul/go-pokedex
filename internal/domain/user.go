package domain

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	Monsters  []primitive.ObjectID `json:"monsters" bson:"monsters"`
	Username  string               `json:"username" bson:"username" validate:"required"`
	Password  string               `json:"password,omitempty" bson:"password" validate:"required,gte=6"`
	Role      *string              `json:"role" bson:"role" validate:"required"`
	CreatedAt time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time            `json:"updated_at" bson:"updated_at"`
}

type UserUpdate struct {
	ID       primitive.ObjectID `json:"_id,omitempty"`
	Username string             `json:"username"`
	Role     *string            `json:"role"`
}

type UserMonsterBody struct {
	MonsterID primitive.ObjectID `json:"monster_id"`
}

type UserLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,gte=6"`
}

type UserList struct {
	TotalCount int     `json:"total_count"`
	TotalPages int     `json:"total_pages"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
	HasMore    bool    `json:"has_more"`
	Users      []*User `json:"users"`
}

type UserWithToken struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ComparePasswords(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

func (u *User) SanitizePassword() {
	u.Password = ""
}

func (u *User) PrepareCreate() error {
	u.Password = strings.TrimSpace(u.Password)

	if err := u.HashPassword(); err != nil {
		return err
	}

	if u.Role != nil {
		*u.Role = strings.ToLower(strings.TrimSpace(*u.Role))
	}
	return nil
}

func (u *UserUpdate) PrepareUpdate() error {
	if u.Role != nil {
		*u.Role = strings.ToLower(strings.TrimSpace(*u.Role))
	}
	return nil
}
