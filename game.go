package main

type Board struct {
	Board      [8][8]int
	Whose_turn int
	Last_piece [2]int
}

func newBoard(whose_turn int, last_piece [2]int) Board {
	var tmp Board
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if (i*3+j)%2 == 0 {
				if i <= 2 {
					tmp.Board[i][j] = 1
				} else if i >= 5 {
					tmp.Board[i][j] = 2
				}
			}
		}
	}
	tmp.Whose_turn = whose_turn
	tmp.Last_piece = last_piece
	return tmp
}

func (game *Board) checkKings() {
	for i := 0; i < 8; i++ {
		if game.Board[0][i] == 2 {
			game.Board[0][i] = 4
		}
		if game.Board[7][i] == 1 {
			game.Board[7][i] = 3
		}
	}
}

func (game *Board) whoWin() int {
	var _turn = game.Whose_turn
	var _last_piece = game.Last_piece
	var possible_turns [28][2]int
	for i := 1; i <= 7; i++ {
		possible_turns[i-1] = [2]int{i, i}
		possible_turns[7+i-1] = [2]int{i, -i}
		possible_turns[14+i-1] = [2]int{-i, i}
		possible_turns[21+i-1] = [2]int{-i, -i}
	}
	var can_move = [2]bool{false, false}
	game.Last_piece = [2]int{-1, -1}
	game.Whose_turn = 0
	for i := 0; i < 8 && !can_move[0]; i++ {
		for j := 0; j < 8 && !can_move[0]; j++ {
			if game.Board[i][j] != 1 && game.Board[i][j] != 3 {
				continue
			}
			for _, k := range possible_turns {
				if game.canMove([2]int{i, j}, _add([2]int{i, j}, k)) {
					can_move[0] = true
					break
				}
			}
		}
	}
	game.Whose_turn = 1
	for i := 0; i < 8 && !can_move[1]; i++ {
		for j := 0; j < 8 && !can_move[1]; j++ {
			if game.Board[i][j] != 2 && game.Board[i][j] != 4 {
				continue
			}
			for _, k := range possible_turns {
				if game.canMove([2]int{i, j}, _add([2]int{i, j}, k)) {
					can_move[1] = true
					break
				}
			}
		}
	}
	game.Whose_turn = _turn
	game.Last_piece = _last_piece
	if !can_move[0] {
		return 1
	} else if !can_move[1] {
		return 0
	} else {
		return -1
	}
}

func (game *Board) isWin(player int) bool {
	return game.whoWin() == player
}

func (game *Board) isGameEnded() bool {
	if game.whoWin() != -1 {
		return true
	} else {
		return false
	}
}

func (game *Board) isEating(from [2]int, to [2]int) bool {
	var dist = _dist(from, to)
	var dir = _div(dist, _len(from, to))
	var cnt = 0
	for i := _add(from, dir); _isBetwen(_add(from, dir), i, _sub(to, dir)); i = _add(i, dir) {
		if game.Board[i[0]][i[1]] != 0 && game.Board[i[0]][i[1]] != game.Whose_turn+1 && game.Board[i[0]][i[1]] != 2+game.Whose_turn+1 {
			cnt++
		}
	}
	return cnt > 0
}

func (game *Board) canMove(from [2]int, to [2]int) bool {
	if from[0] < 0 || from[0] > 7 || from[1] < 0 || from[1] > 7 || to[0] < 0 || to[0] > 7 || to[1] < 0 || to[1] > 7 {
		return false
	}
	if (game.Board[from[0]][from[1]] != game.Whose_turn+1 && game.Board[from[0]][from[1]] != 2+(game.Whose_turn+1)) || game.Board[to[0]][to[1]] != 0 {
		return false
	}
	if !game.isEating(from, to) && game.Last_piece != [2]int{-1, -1} {
		return false
	}
	if game.Last_piece != from && game.Last_piece != [2]int{-1, -1} {
		return false
	}
	var dist = _dist(from, to)
	if _abs(dist[0]) != _abs(dist[1]) {
		return false
	}
	if game.Board[from[0]][from[1]] <= 2 {
		if game.isEating(from, to) {
			if _len(from, to) == 2 {
				return true
			} else {
				return false
			}
		} else {
			if (game.Board[from[0]][from[1]] == 1 && dist[0] == 1) || (game.Board[from[0]][from[1]] == 2 && dist[0] == -1) {
				return true
			} else {
				return false
			}
		}
	} else {
		return true
	}
}

func (game *Board) makeMove(from [2]int, to [2]int) bool {
	if game.canMove(from, to) {
		game.Board[to[0]][to[1]] = game.Board[from[0]][from[1]]
		var dir = _div(_dist(from, to), _len(from, to))
		if game.isEating(from, to) {
			game.Last_piece = to
		} else {
			game.Last_piece = [2]int{-2, -2}
		}
		for i := from; _isBetwen(from, i, _sub(to, dir)); i = _add(i, dir) {
			//game.Turns[len(game.Turns)-1] = append(game.Turns[len(game.Turns)-1], [3]int{i[0], i[1], 0})
			game.Board[i[0]][i[1]] = 0
		}
		game.checkKings()
		return true
	} else {
		return false
	}
}

func (game *Board) endMove() {
	if game.Last_piece == [2]int{-1, -1} {
		return
	}
	game.Whose_turn = (game.Whose_turn + 1) % 2
	game.Last_piece = [2]int{-1, -1}
}

type Game struct {
	Id      int
	Board   Board
	Players [2]string
	Turns   [][8][8]int
}

func newGame(player1 Player, player2 Player) Game {
	var tmp Game
	tmp.Board = newBoard(0, [2]int{-1, -1})
	tmp.Players[0] = player1.Login
	tmp.Players[1] = player2.Login
	tmp.Turns = append(tmp.Turns, tmp.Board.Board)
	return tmp
}

func (game *Game) checkKings() {
	game.Board.checkKings()
}

func (game *Game) whoWin() int {
	return game.Board.whoWin()
}

func (game *Game) isWin(player int) bool {
	return game.Board.isWin(player)
}

func (game *Game) isGameEnded() bool {
	return game.Board.isGameEnded()
}

func (game *Game) isEating(from [2]int, to [2]int) bool {
	return game.Board.isEating(from, to)
}

func (game *Game) canMove(from [2]int, to [2]int) bool {
	return game.Board.canMove(from, to)
}

func (game *Game) makeMove(from [2]int, to [2]int) bool {
	if game.Board.Last_piece == [2]int{-1, -1} {
		game.Turns = append(game.Turns, game.Board.Board)
	}
	response := game.Board.makeMove(from, to)
	game.Turns[len(game.Turns)-1] = game.Board.Board
	return response
}

func (game *Game) endMove() {
	game.Board.endMove()
}
