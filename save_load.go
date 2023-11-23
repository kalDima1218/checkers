package main

import (
	"encoding/json"
	"os"
	"time"
)

// _load loads data from a file and unmarshals it into the given struct pointer.
//
// Parameters:
// - data: a pointer to the struct that the data will be unmarshaled into.
// - file: the name of the file to load the data from.
func _load(data *any, file string) {
	data_byte, _ := os.ReadFile("json/" + file)
	json.Unmarshal(data_byte, data)
}

// _save writes data to a file in JSON format.
//
// Parameters:
//
//	data - the data to be saved.
//	file - the name of the file to save the data to.
//
// Return type: none.
func _save(data any, file string) {
	data_byte, _ := json.Marshal(data)
	os.WriteFile("json/"+file, data_byte, 0666)
}

// _autosave is a Go function that performs automatic saving of data at regular intervals.
//
// This function does not take any parameters.
// It does not return any values.
func _autosave() {
	for true {
		_save(players, "players.json")
		_save(logins, "logins.json")
		_save(names, "names.json")
		_save(games, "games.json")
		<-time.After(time.Second)
	}
}

// _init_load initializes the data by reading from JSON files.
//
// It reads the "games.json" file and unmarshals the data into the "games" variable.
// It reads the "players.json" file and unmarshals the data into the "players" variable.
// It reads the "logins.json" file and unmarshals the data into the "logins" variable.
// It reads the "names.json" file and unmarshals the data into the "names" variable.
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
