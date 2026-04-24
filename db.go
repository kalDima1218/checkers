package main

import (
	"database/sql"
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var DB *sql.DB

// TODO поставить проверку от инъекций
// config.json {"User": "", "Passwd": "", "Addr": "", "DBName": ""}
func loadDB() {
	config := make(map[string]string)
	configByte, _ := os.ReadFile("config.json")
	json.Unmarshal(configByte, &config)
	cfg := mysql.Config{
		User:   config["User"],
		Passwd: config["Passwd"],
		Net:    "tcp",
		Addr:   config["Addr"],
		DBName: config["DBName"],
	}
	var err error
	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
}

func insertUser(login, password, username string) error {
	_, err := DB.Exec("INSERT INTO Users (login, password, username, elo) VALUES (?, ?, ?, 1500)", login, password, username)
	return err
}

func getUsername(login string) string {
	var username string
	err := DB.QueryRow("SELECT username FROM Users WHERE login = ?;", login).Scan(&username)
	if err != nil {
		log.Println(err)
		return ""
	}
	return username
}

func getPassword(login string) string {
	var password string
	err := DB.QueryRow("SELECT password FROM Users WHERE login = ?;", login).Scan(&password)
	if err != nil {
		log.Println(err)
		return ""
	}
	return password
}

func getLastSeen(login string) int64 {
	var lastSeen int64
	err := DB.QueryRow("SELECT last_seen FROM Users WHERE login = ?;", login).Scan(&lastSeen)
	if err != nil {
		return 0
	}
	return lastSeen
}

func setLastSeen(login string, lastSeen int64) {
	_, err := DB.Exec("UPDATE Users SET last_seen = ? WHERE login = ?;", lastSeen, login)
	if err != nil {
		log.Fatal(err)
	}
}

func isFreeLogin(login string) bool {
	var cnt int
	if err := DB.QueryRow("SELECT COUNT(*) FROM Users WHERE login = ?;", login).Scan(&cnt); err != nil {
		log.Fatal(err)
	}
	return cnt == 0
}

func insertGame(id string, game *Game) error {
	gameJsonByte, _ := json.Marshal(game)
	gameJson := string(gameJsonByte)

	_, err := DB.Exec("INSERT INTO Games (id, game) VALUES (?, ?)", id, gameJson)
	return err
}

func setGame(id string, game *Game) {
	gameJsonByte, _ := json.Marshal(game)
	gameJson := string(gameJsonByte)

	_, err := DB.Exec("UPDATE Games SET game = ? WHERE id = ?;", gameJson, id)
	if err != nil {
		log.Fatal(err)
	}
}

func getGame(id string) (*Game, bool) {
	if !isGameExists(id) {
		return nil, false
	}

	var gameJson string
	if err := DB.QueryRow("SELECT game FROM Games WHERE id = ?;", id).Scan(&gameJson); err != nil {
		log.Println(err)
		return nil, false
	}

	gameJsonByte := []byte(gameJson)
	var game Game
	json.Unmarshal(gameJsonByte, &game)
	return &game, true
}

func isGameExists(id string) bool {
	var cnt int
	if err := DB.QueryRow("SELECT COUNT(*) FROM Games WHERE id = ?;", id).Scan(&cnt); err != nil {
		log.Fatal(err)
	}
	return cnt != 0
}
