package model

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	ID        uuid.UUID  `json:"id"`         // type:uuid
	Name      string     `json:"name"`       // type:varchar
	Email     string     `json:"email"`      // type:varchar
	Phone     string     `json:"phone"`      // type:varchar
	Password  string     `json:"-"`          // type:varchar
	Gamertag  string     `json:"gamertag"`   // type:varchar
	Deleted   bool       `json:"is_deleted"` // type:bool
	DeletedAt *time.Time `json:"deleted_at"` // type:timestamp
	UpdatedAt time.Time  `json:"updated_at"` // type:timestamp
	CreatedAt time.Time  `json:"created_at"` // type:timestamp
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

func (u *UserModel) ComparePasswords(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}

	return nil
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

func (u *UserModel) GetDeletedAt() *time.Time {
	return u.DeletedAt
}

func (u *UserModel) GetCreatedAt() time.Time {
	return u.CreatedAt
}

func (u *UserModel) GetUpdatedAt() time.Time {
	return u.UpdatedAt
}
