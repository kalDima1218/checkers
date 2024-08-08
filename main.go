package main

func _init() {
	for i := 1; i <= 7; i++ {
		POSSIBLE_TURNS[i-1] = [2]int{i, i}
		POSSIBLE_TURNS[7+i-1] = [2]int{i, -i}
		POSSIBLE_TURNS[14+i-1] = [2]int{-i, i}
		POSSIBLE_TURNS[21+i-1] = [2]int{-i, -i}
	}

}

func main() {
	_init()
	//_init_load()
	//go _autosave()
	//insertGame("{}")
	startSite()
	//print(isGameExists(0))
	//print(insertUser("b", "b", "b"))
	//game := newGame("", "")
	//modifyGame(0, &game)
	//print(game.Board.Board[0][0])
	//var game *Game
	//game, _ = getGame("")
	//print(game)
	//game = getGame(0)
	//_bot_vs_bot()
}
