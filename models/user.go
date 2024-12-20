package models

import (
	"errors"

	"example.com/rest-api/db"
	"example.com/rest-api/utils"
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *User) Save() error {
	query := `
		INSERT INTO users (email, password)
		VALUES (?, ?)
	`

	hashed, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	result, err := db.DB.Exec(query, u.Email, hashed)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = id
	return nil
}

func GetUsers() ([]User, error) {
	query := `
		SELECT id, email
		FROM users
	`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *User) ValidadeCredentials() error {
	query := `
		SELECT id, password
		FROM users
		WHERE email = ?
	`
	var hashedPassword string
	err := db.DB.QueryRow(query, u.Email).Scan(&u.ID, &hashedPassword)
	if err != nil {
		return err
	}

	if !utils.CheckPasswordHash(u.Password, hashedPassword) {
		return errors.New("invalid credentials")
	}

	return nil
}
