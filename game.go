package main

import (
	"encoding/json"
)

// Game represents a game.
type Game struct {
	Id           int
	Board        [8][8]int
	Whose_turn   int
	Players      [2]string
	Last_piece   [2]int
	Current_turn int
	Turns        [][8][8]int
}

// newGame initializes a new game with two players.
//
// It takes two parameters: player1 of type Player and player2 of type Player.
// It returns a Game.
func newGame(player1 Player, player2 Player) Game {
	var tmp Game
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
	tmp.Players[0] = player1.Login
	tmp.Players[1] = player2.Login
	tmp.Whose_turn = 0
	tmp.Last_piece = [2]int{-1, -1}
	tmp.Current_turn = 0
	tmp.Turns = append(tmp.Turns, tmp.Board)
	return tmp
}

// getJsonBoard returns the JSON representation of the game board.
//
// No parameters.
// Returns a string.
func (game *Game) getJsonBoard() string {
	board_json, _ := json.Marshal(game.Board)
	return string(board_json)
}

// checkKings updates the game board to replace the kings with their respective upgraded pieces.
//
// No parameters.
// No return values.
func (game *Game) checkKings() {
	for i := 0; i < 8; i++ {
		if game.Board[0][i] == 2 {
			game.Board[0][i] = 4
		}
		if game.Board[7][i] == 1 {
			game.Board[7][i] = 3
		}
	}
}

// whoWin determines the winner of the game.
//
// This function checks if either player 0 or player 1 can make a move on the board.
// If player 0 cannot make a move, it returns 1. If player 1 cannot make a move, it returns 0.
// If both players can make a move, it returns -1.
//
// Return:
// - int: The winner of the game. 1 for player 0, 0 for player 1, -1 for a tie.
func (game *Game) whoWin() int {
	var _turn = game.Whose_turn
	var possible_turns [28][2]int
	for i := 1; i <= 7; i++ {
		possible_turns[i-1] = [2]int{i, i}
		possible_turns[7+i-1] = [2]int{i, -i}
		possible_turns[14+i-1] = [2]int{-i, i}
		possible_turns[21+i-1] = [2]int{-i, -i}
	}
	var can_move = [2]bool{false, false}
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
	if !can_move[0] {
		return 1
	} else if !can_move[1] {
		return 0
	} else {
		return -1
	}
}

// isWin checks if the specified player has won the game.
//
// player - The player to check for a win.
// bool - Returns true if the player has won, false otherwise.
func (game *Game) isWin(player int) bool {
	return game.whoWin() == player
}

// isGameEnded checks if the game has ended.
//
// Returns true if the game has ended, false otherwise.
func (game *Game) isGameEnded() bool {
	if game.whoWin() != -1 {
		return true
	} else {
		return false
	}
}

// isEating checks if a piece is eating another piece in the game.
//
// Parameters:
//
//	from [2]int: the starting position of the piece
//	to [2]int: the destination position of the piece
//
// Returns:
//
//	bool: true if the piece is eating another piece, false otherwise
func (game *Game) isEating(from [2]int, to [2]int) bool {
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

// canMove checks if a move is valid in the game.
//
// Parameters:
// - from: the starting position of the piece.
// - to: the destination position of the piece.
//
// Returns:
// - bool: true if the move is valid, false otherwise.
func (game *Game) canMove(from [2]int, to [2]int) bool {
	if from[0] < 0 || from[0] > 7 || from[1] < 0 || from[1] > 7 || to[0] < 0 || to[0] > 7 || to[1] < 0 || to[1] > 7 {
		return false
	}
	if (game.Board[from[0]][from[1]] != game.Whose_turn+1 && game.Board[from[0]][from[1]] != 2+(game.Whose_turn+1)) || game.Board[to[0]][to[1]] != 0 {
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

// makeMove is a function that makes a move in the game.
//
// It takes two parameters: from - the starting position of the piece to be moved,
// and to - the destination position where the piece will be moved.
//
// It returns a boolean value indicating whether the move was successful or not.
func (game *Game) makeMove(from [2]int, to [2]int) bool {
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
		//game.Turns[len(game.Turns)-1] = append(game.Turns[len(game.Turns)-1], [3]int{to[0], to[1], game.Board[to[0]][to[1]]})
		return true
	} else {
		return false
	}
}

// endMove updates the game state after a player's turn is completed.
//
// The function does not take any parameters.
// It does not return any values.
func (game *Game) endMove() {
	if game.Last_piece == [2]int{-1, -1} {
		return
	}
	game.Whose_turn = (game.Whose_turn + 1) % 2
	game.Last_piece = [2]int{-1, -1}
	game.Current_turn = game.Current_turn + 1
	game.Turns = append(game.Turns, game.Board)
}
