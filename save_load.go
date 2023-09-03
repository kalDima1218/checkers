package main

import (
	"encoding/json"
	"os"
	"time"
)

func _load(data *any, file string) {
	data_byte, _ := os.ReadFile("json/" + file)
	json.Unmarshal(data_byte, data)
}

func _save(data any, file string) {
	data_byte, _ := json.Marshal(data)
	os.WriteFile("json/"+file, data_byte, 0666)
}

func _autosave() {
	for true {
		_save(players, "players.json")
		_save(logins, "logins.json")
		_save(names, "names.json")
		_save(games, "games.json")
		<-time.After(time.Second)
	}
}

func _init_load() {
	games_byte, _ := os.ReadFile("json/games.json")
	json.Unmarshal(games_byte, &games)
	players_byte, _ := os.ReadFile("json/players.json")
	json.Unmarshal(players_byte, &players)
	logins_byte, _ := os.ReadFile("json/logins.json")
	json.Unmarshal(logins_byte, &logins)
	names_byte, _ := os.ReadFile("json/names.json")
	json.Unmarshal(names_byte, &names)
}
