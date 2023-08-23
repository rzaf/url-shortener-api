package models

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/rzaf/url-shortener-api/database"
	"log"
	"math"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id              int64      `json:"id"`
	Email           string     `json:"email"`
	hashed_password string     `json:"-"`
	Api_key         string     `json:"api_key"`
	IsAdmin         bool       `json:"-"`
	Created_at      *time.Time `json:"created_at"`
	Updated_at      *time.Time `json:"updated_at,omitempty"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func GetUsers() []User {
	sql := "SELECT id,email,api_key,is_admin,created_at,updated_at FROM `users` ORDER BY id;"
	rows, err := database.Db.Query(sql)
	if err != nil {
		log.Panicln(err.Error())
	}
	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.Id, &u.Email, &u.Api_key, &u.IsAdmin, &u.Created_at, &u.Updated_at)
		if err != nil {
			log.Panicln(err.Error())
		}
		users = append(users, u)
	}
	// fmt.Println("users:", users)
	return users
}

func GetUserByApiKey(apiKey string) (*User, error) {
	sql := "SELECT id,email,hashed_password,api_key,is_admin,created_at,updated_at FROM `users` WHERE api_key=?"
	rows, err := database.Db.Query(sql, apiKey)
	if err != nil {
		log.Panicln(err.Error())
	}
	if !rows.Next() {
		return nil, nil
	}
	var user User
	rows.Scan(&user.Id, &user.Email, &user.hashed_password, &user.Api_key, &user.IsAdmin, &user.Created_at, &user.Updated_at)
	// fmt.Println("user:", user)
	return &user, nil
}

func GetUserById(id int) (*User, error) {
	sql := "SELECT id,email,hashed_password,api_key,is_admin,created_at,updated_at FROM `users` WHERE id=?"
	rows, err := database.Db.Query(sql, id)
	if err != nil {
		log.Panicln(err.Error())
	}
	if !rows.Next() {
		return nil, nil
	}
	var user User
	rows.Scan(&user.Id, &user.Email, &user.hashed_password, &user.Api_key, &user.IsAdmin, &user.Created_at, &user.Updated_at)
	// fmt.Println("user:", user)
	return &user, nil
}

func CreateUser(email string, pass string) (*User, error) {
	db := database.Db
	query := "INSERT INTO `users` (email,hashed_password,api_key) VALUES (?,?,?)"
	hashedPassword, err := HashPassword(pass)
	if err != nil {
		log.Panicln(err.Error())
	}
	apiKey := generateSecureToken(128)
	res, err := db.Exec(query, email, hashedPassword, apiKey)
	if err != nil {
		if mySqlErr, ok := err.(*mysql.MySQLError); ok {
			if mySqlErr.Number == 1062 {
				return nil, errors.New(mySqlErr.Message)
			}
		}
		log.Panicln(err.Error())
	}
	id, _ := res.LastInsertId()
	t := time.Now()
	return &User{
		Id:              id,
		hashed_password: hashedPassword,
		Email:           email,
		Api_key:         apiKey,
		Created_at:      &t,
	}, nil
}

func EditUserEmailOrPassword(oldUser *User, newEmail string, newPassword string) error {
	db := database.Db
	query := "UPDATE `users` SET email=?,hashed_password=?,updated_at=? WHERE id=?;"

	var err error
	if newPassword != "" {
		oldUser.hashed_password, err = HashPassword(newPassword)
		if err != nil {
			log.Panicln(err.Error())
		}
	}
	if newEmail != "" {
		oldUser.Email = newEmail
	}

	_, err = db.Exec(query, oldUser.Email, oldUser.hashed_password, time.Now(), oldUser.Id)
	if err != nil {
		if mySqlErr, ok := err.(*mysql.MySQLError); ok {
			if mySqlErr.Number == 1062 {
				return errors.New(mySqlErr.Message)
			}
		}
		log.Panicln(err.Error())
	}
	return nil
}

func EditUserApiKey(user *User, oldPass string) error {
	db := database.Db
	query := "UPDATE `users` SET api_key=?,updated_at=? WHERE id=?;"

	if h := bcrypt.CompareHashAndPassword([]byte(user.hashed_password), []byte(oldPass)); h != nil {
		return errors.New("wrong password")
	}
	apiKey := generateSecureToken(128)
	user.Api_key = apiKey
	_, err := db.Exec(query, apiKey, time.Now(), user.Id)
	if err != nil {
		log.Panicln(err.Error())
	}
	return nil
}

func DeleteUser(user *User) {
	db := database.Db
	query := "DELETE FROM `users` WHERE id=?;"
	_, err := db.Exec(query, user.Id)
	if err != nil {
		log.Panicln(err.Error())
	}
}

func generateSecureToken(length int) string {
	buff := make([]byte, int(math.Ceil(float64(length)/2)))
	if _, err := rand.Read(buff); err != nil {
		log.Panicln(err.Error())
	}
	str := hex.EncodeToString(buff)
	return str[:length]
}
