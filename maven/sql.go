package maven

import (
	"crypto/md5"
	"encoding/hex"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strconv"
	"time"
)

func Create(con DSN) Database {
	connec, err := sqlx.Open("mysql", con.Format())
	if err != nil {
		panic(err)
	}
	connec.SetConnMaxLifetime(time.Minute * 3)
	connec.SetMaxOpenConns(10)
	return Database{db: connec}
}

func (d DSN) Format() string {
	fm := ""
	if d.Username != "" {
		fm += d.Username
		if d.Password != "" {
			fm += ":"
			fm += d.Password
		}
		fm += "@"
	}
	if d.Host != "" {
		fm += "("
		fm += d.Host
		if d.Port != 0 {
			fm += ":"
			fm += strconv.Itoa(d.Port)
		}
		fm += ")"
	}
	fm += "/"
	fm += d.Database
	return fm
}

type Database struct {
	db *sqlx.DB
}

type User struct {
	Password string `DB:"password"`
}

func (d Database) HasAccess(path string, username string, password string) bool {
	row := d.db.QueryRowx("SELECT password FROM users WHERE username=?", username)
	user := User{}
	err := row.StructScan(&user)
	if err != nil {
		return false
	}
	return user.Password == Hash(password)
}

func Hash(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}
