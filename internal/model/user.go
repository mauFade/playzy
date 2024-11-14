package model

import (
	"time"

	"github.com/google/uuid"
)

type UserModel struct {
	ID        uuid.UUID  // type:uuid
	Name      string     // type:varchar
	Email     string     // type:varchar
	Phone     string     // type:varchar
	Password  string     // type:varchar
	Gamertag  string     // type:varchar
	Deleted   bool       // type:bool
	DeletedAt *time.Time // type:timestamp
	UpdatedAt time.Time  // type:timestamp
	CreatedAt time.Time  // type:timestamp
}

func NewUserModel(
	id uuid.UUID,
	name,
	email,
	phone,
	gamertag,
	password string,
	deleted bool,
	deletedAt *time.Time,
	updatedAt,
	createdAt time.Time,
) *UserModel {
	return &UserModel{
		ID:        id,
		Name:      name,
		Email:     email,
		Phone:     phone,
		Password:  password,
		Gamertag:  gamertag,
		Deleted:   deleted,
		DeletedAt: deletedAt,
		UpdatedAt: updatedAt,
		CreatedAt: createdAt,
	}
}

func (u *UserModel) GetID() uuid.UUID {
	return u.ID
}

func (u *UserModel) GetName() string {
	return u.Name
}

func (u *UserModel) SetName(name string) {
	u.Name = name
}

func (u *UserModel) GetEmail() string {
	return u.Email
}

func (u *UserModel) SetEmail(email string) {
	u.Email = email
}

func (u *UserModel) GetPhone() string {
	return u.Phone
}

func (u *UserModel) SetPhone(phone string) {
	u.Phone = phone
}

func (u *UserModel) GetGamertag() string {
	return u.Gamertag
}

func (u *UserModel) SetGamertag(gamertag string) {
	u.Gamertag = gamertag
}

func (u *UserModel) GetPassword() string {
	return u.Password
}

func (u *UserModel) SetPassword(password string) {
	u.Password = password
}

func (u *UserModel) IsDeleted() bool {
	return u.Deleted
}

func (u *UserModel) SetDeleted(deleted bool) {
	u.Deleted = deleted
}
