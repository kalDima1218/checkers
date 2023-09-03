package main

type Player struct {
	Login    string
	Password string
	Name     string
	Elo      int
}

func newPlayer(name string, login string, password string) Player {
	var tmp Player
	tmp.Login = login
	tmp.Password = password
	tmp.Name = name
	return tmp
}
