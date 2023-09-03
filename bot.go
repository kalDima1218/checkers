package main

import (
	"fmt"
	"math"
	"math/rand"
)

var MAX_DEPTH = 6

var POSSIBLE_TURNS [28][2]int

var BOT = newBot()

type Bot struct {
	max_depth                      int
	cell_cost, king_cost, win_cost float64
	cost_matrix                    [8][8]float64
	cost_vertical                  [8]float64
	moves_table                    MapGame
}

func newBot() *Bot {
	var tmp Bot
	tmp.max_depth = MAX_DEPTH
	tmp.cell_cost = 10
	tmp.king_cost = 20
	tmp.win_cost = 1000
	tmp.cost_matrix = [8][8]float64{
		{1.0, 0, 1.2, 0, 1.2, 0, 1.15, 0},
		{0, 1.15, 0, 1.2, 0, 1.15, 0, 1.0},
		{1.0, 0, 1.2, 0, 1.2, 0, 1.15, 0},
		{0, 1.15, 0, 1.2, 0, 1.15, 0, 1.0},
		{1.0, 0, 1.2, 0, 1.2, 0, 1.15, 0},
		{0, 1.15, 0, 1.2, 0, 1.15, 0, 1.0},
		{1.0, 0, 1.2, 0, 1.2, 0, 1.15, 0},
		{0, 1.15, 0, 1.2, 0, 1.15, 0, 1.0},
	}
	tmp.moves_table = *newMapGame()
	return &tmp
}

func (bot *Bot) evaluate(game *Game) float64 {
	if game.isGameEnded() {
		if game.isWin(0) {
			return bot.win_cost
		} else {
			return -bot.win_cost
		}
	}
	var d_cell = 0.0
	var d_king = 0.0
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if game.Board[i][j] == 1 {
				//math.Exp(-math.Abs(3.5-float64(j))/7+0.25)
				d_cell += bot.cost_matrix[i][j] * math.Pow(1.1, float64(i))
			} else if game.Board[i][j] == 3 {
				d_king += bot.cost_matrix[i][j]
			} else if game.Board[i][j] == 2 {
				d_cell -= bot.cost_matrix[i][j] * math.Pow(1.1, float64(7-i))
			} else if game.Board[i][j] == 4 {
				d_king -= bot.cost_matrix[i][j]
			}
		}
	}
	return bot.cell_cost*d_cell + bot.king_cost*d_king
}

func (bot *Bot) _dfsStreak(game Game, me int, enemy int) Game {
	var max_game = game
	for _, k := range POSSIBLE_TURNS {
		if game.canMove(game.Last_piece, _add(game.Last_piece, k)) {
			var _game = game
			_game.makeMove(game.Last_piece, _add(game.Last_piece, k))
			if game.Turn == 0 {
				_game := bot._dfsStreak(_game, me, enemy)
				if bot.evaluate(&_game) > bot.evaluate(&max_game) {
					max_game = _game
				}
			} else {
				_game := bot._dfsStreak(_game, me, enemy)
				if bot.evaluate(&_game) < bot.evaluate(&max_game) {
					max_game = _game
				}
			}
		}
	}
	return max_game
}

func (bot *Bot) _findBestMove(game Game, depth int, me int, enemy int, prev_score float64) float64 {
	//val, ok := bot.moves_table.get(newItemGame(game))
	//if ok {
	//	return val
	//}
	val, ok := bot.moves_table.get(game)
	if ok {
		return val
	}
	if depth == bot.max_depth || game.isGameEnded() {
		return bot.evaluate(&game)
	}
	var possible_turns [28][2]int
	for i := 1; i <= 7; i++ {
		possible_turns[i-1] = [2]int{i, i}
		possible_turns[7+i-1] = [2]int{i, -i}
		possible_turns[14+i-1] = [2]int{-i, i}
		possible_turns[21+i-1] = [2]int{-i, -i}
	}
	var max_score float64
	if game.Turn == 0 {
		max_score = -1e9
	} else {
		max_score = 1e9
	}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if game.Board[i][j] != game.Turn+1 && game.Board[i][j] != 2+game.Turn+1 {
				continue
			}
			for _, k := range possible_turns {
				if game.canMove([2]int{i, j}, _add([2]int{i, j}, k)) {
					var _game = game
					_game.makeMove([2]int{i, j}, _add([2]int{i, j}, k))
					_game = bot._dfsStreak(_game, me, enemy)
					_game.endMove()
					if game.Turn == 0 {
						max_score = math.Max(max_score, bot._findBestMove(_game, depth+1, me, enemy, max_score))
						if max_score < prev_score {
							return max_score
						}
					} else {
						max_score = math.Min(max_score, bot._findBestMove(_game, depth+1, me, enemy, max_score))
						if max_score > prev_score {
							return max_score
						}
					}
				}
			}
		}
	}
	bot.moves_table.insert(game, max_score)
	return max_score
}

func (bot *Bot) gameTemp(game *Game) float64 {
	if game.Turn == 0 {
		return bot._findBestMove(*game, 0, 0, 1, -1e9)
	} else {
		return bot._findBestMove(*game, 0, 0, 1, 1e9)
	}
}

func _findBestMove_goroutine(bot *Bot, game Game, me int, enemy int, move_chanel chan Move) {
	if game.Turn == 0 {
		move_chanel <- newMove(bot._findBestMove(game, 1, me, enemy, -1e9), game)
	} else {
		move_chanel <- newMove(bot._findBestMove(game, 1, me, enemy, 1e9), game)
	}
}

func (bot *Bot) findBestMove(game Game, me int, enemy int) Game {
	var possible_turns [28][2]int
	for i := 1; i <= 7; i++ {
		possible_turns[i-1] = [2]int{i, i}
		possible_turns[7+i-1] = [2]int{i, -i}
		possible_turns[14+i-1] = [2]int{-i, i}
		possible_turns[21+i-1] = [2]int{-i, -i}
	}
	var move_chanel = make(chan Move)
	var max_score float64
	if game.Turn == 0 {
		max_score = -1e9
	} else {
		max_score = 1e9
	}
	var turns []Move
	var cnt = 0
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if game.Board[i][j] != game.Turn+1 && game.Board[i][j] != 2+game.Turn+1 {
				continue
			}
			for _, k := range possible_turns {
				if game.canMove([2]int{i, j}, _add([2]int{i, j}, k)) {
					var _game = game
					_game.makeMove([2]int{i, j}, _add([2]int{i, j}, k))
					_game = bot._dfsStreak(_game, me, enemy)
					_game.endMove()
					cnt++
					go _findBestMove_goroutine(bot, _game, me, enemy, move_chanel)
				}
			}
		}
	}
	for cnt > 0 {
		move := <-move_chanel
		if (game.Turn == 0 && move.score > max_score) || (game.Turn == 1 && move.score < max_score) {
			turns = make([]Move, 1)
			turns[0] = move
			max_score = move.score
		} else if move.score == max_score {
			turns = append(turns, move)
		}
		cnt--
	}
	bot.moves_table.clear()
	return turns[rand.Int()%len(turns)].game
}

func _print_board(game *Game) {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			fmt.Print(game.Board[i][j])
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}

func _bot_vs_bot() {
	var player1 = newPlayer("", "", "")
	var player2 = newPlayer("", "", "")
	var game = newGame(player1, player2)
	var cnt = 0
	for !game.isGameEnded() {
		game = BOT.findBestMove(game, cnt, (cnt+1)%2)
		fmt.Println(math.Round(BOT.gameTemp(&game)))
		_print_board(&game)
		cnt++
		cnt %= 2
	}
	fmt.Println(game.whoWin())
}
