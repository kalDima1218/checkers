package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
)

// поставить проверку от инъекций
// config.json {"User": "", "Passwd": "", "Addr": "", "DBName": ""}
func getDB() *sql.DB {
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
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func insertUser(login, password, username string) error {
	db := getDB()
	defer db.Close()

	_, err := db.Query(fmt.Sprintf("INSERT INTO Users (login, password, username, elo) VALUES ('%v','%v','%v',1500)", login, password, username))
	return err
}

func getUsername(login string) string {
	db := getDB()
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("SELECT username FROM Users WHERE login = '%v';", login))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var username string
	rows.Next()
	if rows.Scan(&username) != nil {
		log.Fatal(err)
	}
	return username
}

func getPassword(login string) string {
	db := getDB()
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("SELECT password FROM Users WHERE login = '%v';", login))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var password string
	rows.Next()
	if rows.Scan(&password) != nil {
		log.Fatal(err)
	}
	return password
}

func isFreeLogin(login string) bool {
	db := getDB()
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("SELECT COUNT(*) FROM Users WHERE login = '%v';", login))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var cnt int
	rows.Next()
	if rows.Scan(&cnt) != nil {
		log.Fatal(err)
	}
	return cnt == 0
}

func insertGame(id string, game *Game) error {
	db := getDB()
	defer db.Close()

	gameJsonByte, _ := json.Marshal(game)
	gameJson := string(gameJsonByte)

	_, err := db.Query(fmt.Sprintf("INSERT INTO Games (id, game) VALUES ('%v', '%v')", id, gameJson))
	return err
}

func setGame(id string, game *Game) {
	db := getDB()
	defer db.Close()

	gameJsonByte, _ := json.Marshal(game)
	gameJson := string(gameJsonByte)

	_, err := db.Query(fmt.Sprintf("UPDATE Games SET game = '%v' WHERE id = '%v';", gameJson, id))
	if err != nil {
		log.Fatal(err)
	}
}

func getGame(id string) (*Game, bool) {
	if !isGameExists(id) {
		return nil, false
	}

	db := getDB()
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("SELECT game FROM Games WHERE id = '%v';", id))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var gameJson string
	rows.Next()
	if rows.Scan(&gameJson) != nil {
		log.Fatal(err)
	}

	gameJsonByte := []byte(gameJson)
	var game Game
	json.Unmarshal(gameJsonByte, &game)
	return &game, true
}

func isGameExists(id string) bool {
	db := getDB()
	defer db.Close()

	rows, err := db.Query(fmt.Sprintf("SELECT COUNT(*) FROM Games WHERE id = '%v';", id))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var cnt int
	rows.Next()
	if rows.Scan(&cnt) != nil {
		log.Fatal(err)
	}
	return cnt != 0
}
