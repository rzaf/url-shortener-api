package models

import (
	"github.com/rzaf/url-shortener-api/database"
	"github.com/rzaf/url-shortener-api/helpers"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Url struct {
	Id         int64      `json:"id"`
	Short      string     `json:"shortened"`
	Shortened  string     `json:"shortened_url"`
	Url        string     `json:"url"`
	User_id    int64      `json:"user_id"`
	Created_at *time.Time `json:"created_at"`
	Updated_at *time.Time `json:"updated_at,omitempty"`
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func intPow(n, m int) int {
	if m == 0 {
		return 1
	}
	result := n
	for i := 0; i < m-1; i++ {
		result *= n
	}
	return result
}

func IdToShort(id int) string {
	var reminders []int
	for id > 0 {
		reminders = append(reminders, id%len(letters))
		id = id / len(letters)
	}
	result := ""
	for i := len(reminders) - 1; i >= 0; i-- {
		result += string(letters[reminders[i]])
	}
	return result
}

func ShortToId(shortened string) int {
	var id int
	index := 0
	for i := len(shortened) - 1; i >= 0; i-- {
		m := strings.IndexByte(letters, shortened[index])
		id += m * intPow(len(letters), i)
		index++
	}
	return id
}

func GetUrl(short string) *Url {
	sql := "SELECT id,url,user_id,created_at,updated_at FROM `urls` WHERE id=? ;"
	rows, err := database.Db.Query(sql, ShortToId(short))
	if err != nil {
		panic(err.Error())
	}
	if !rows.Next() {
		panic(helpers.NewServerError("url not found", 404))
	}
	var newUrl Url
	rows.Scan(&newUrl.Id, &newUrl.Url, &newUrl.User_id, &newUrl.Created_at, &newUrl.Updated_at)
	newUrl.Short = short
	newUrl.Shortened = os.Getenv("URL") + "/urls/" + short
	return &newUrl
}

func DeleteUrl(short string) {
	sql := "DELETE FROM `urls` WHERE id=?;"
	res, err := database.Db.Exec(sql, ShortToId(short))
	if err != nil {
		panic(err.Error())
	}

	if _, err = res.RowsAffected(); err != nil {
		panic(helpers.NewServerError("url `"+short+"` not found", 404))
	}
}

func GetUserUrls(user_id int) []Url {
	sql := "SELECT id,url,user_id,created_at,updated_at FROM `urls` WHERE user_id=? ORDER BY id;"
	rows, err := database.Db.Query(sql, user_id)
	if err != nil {
		panic(err.Error())
	}
	var urls []Url
	for rows.Next() {
		var u Url
		err := rows.Scan(&u.Id, &u.Url, &u.User_id, &u.Created_at, &u.Updated_at)
		if err != nil {
			panic(err.Error())
		}
		u.Short = IdToShort(int(u.Id))
		u.Shortened = os.Getenv("URL") + "/urls/" + u.Short
		urls = append(urls, u)
	}
	// fmt.Println("users:", users)
	return urls
}

func CreateUrl(url string, userId int64) int64 {
	sql := "INSERT INTO `urls` (url,user_id,created_at) VALUES (?,?,?)"
	res, err := database.Db.Exec(sql, url, userId, time.Now())
	if err != nil {
		if mySqlErr, ok := err.(*mysql.MySQLError); ok {
			if mySqlErr.Number == 1062 {
				panic(helpers.NewServerError("duplicate url "+url, 400))
			}
		}
		panic(err.Error())
	}
	id, _ := res.LastInsertId()
	return id
}

func EditUrl(short string, newUrl string) {
	sql := "UPDATE `urls` SET url=?,updated_at=? WHERE id=? ;"
	res, err := database.Db.Exec(sql, newUrl, time.Now(), ShortToId(short))
	if err != nil {
		if mySqlErr, ok := err.(*mysql.MySQLError); ok {
			if mySqlErr.Number == 1062 {
				panic(helpers.NewServerError("duplicate url "+newUrl, 400))
			}
		}
		panic(err.Error())
	}
	if _, err = res.RowsAffected(); err != nil {
		// panic("id not found")
		panic(helpers.NewServerError("url `"+short+"` not found", 404))
	}
}
