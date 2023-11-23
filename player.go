package main

// Player represents a player in the game.
type Player struct {
	Login    string
	Password string
	Name     string
	Elo      int
}

// newPlayer creates a new player with the given name, login, and password.
//
// Parameters:
// - name: the name of the player.
// - login: the login of the player.
// - password: the password of the player.
//
// Returns:
// - Player: the newly created player.
func newPlayer(name string, login string, password string) Player {
	var tmp Player
	tmp.Login = login
	tmp.Password = password
	tmp.Name = name
	tmp.Elo = 1500
	return tmp
}

// recalculateElo updates the Elo ratings of the winner and loser players.
//
// Parameters:
// - winner: a pointer to the Player who won the match.
// - loser: a pointer to the Player who lost the match.
func recalculateElo(winner, loser *Player) {
	winner.Elo += 10
	loser.Elo -= 10
}
