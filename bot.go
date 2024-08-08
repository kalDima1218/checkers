package main

import (
	"fmt"
	"math"
	"math/rand"
)

var MAX_DEPTH = 6
var POSSIBLE_TURNS = [28][2]int{{1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}, {6, 6}, {7, 7}, {1, -1}, {2, -2}, {3, -3}, {4, -4}, {5, -5}, {6, -6}, {7, -7}, {-1, 1}, {-2, 2}, {-3, 3}, {-4, 4}, {-5, 5}, {-6, 6}, {-7, 7}, {-1, -1}, {-2, -2}, {-3, -3}, {-4, -4}, {-5, -5}, {-6, -6}, {-7, -7}}
var BOT = newBot()

type Bot struct {
	maxDepth                    int
	cellCost, kingCost, winCost float64
	costMatrix                  [8][8]float64
	costVertical                [8]float64
	movesTable                  MapGame
}

func newBot() Bot {
	var tmp Bot
	tmp.maxDepth = MAX_DEPTH
	tmp.cellCost = 10
	tmp.kingCost = 20
	tmp.winCost = 1000
	tmp.costMatrix = [8][8]float64{
		{1.0, 0, 1.2, 0, 1.2, 0, 1.15, 0},
		{0, 1.15, 0, 1.2, 0, 1.15, 0, 1.0},
		{1.0, 0, 1.2, 0, 1.2, 0, 1.15, 0},
		{0, 1.15, 0, 1.2, 0, 1.15, 0, 1.0},
		{1.0, 0, 1.2, 0, 1.2, 0, 1.15, 0},
		{0, 1.15, 0, 1.2, 0, 1.15, 0, 1.0},
		{1.0, 0, 1.2, 0, 1.2, 0, 1.15, 0},
		{0, 1.15, 0, 1.2, 0, 1.15, 0, 1.0},
	}
	tmp.movesTable = *newMapGame()
	return tmp
}

func (bot *Bot) evaluate(game Board) float64 {
	if game.isGameEnded() {
		if game.isWin(0) {
			return bot.winCost
		} else {
			return -bot.winCost
		}
	}
	var dCell = 0.0
	var dKing = 0.0
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if game.Board[i][j] == 1 {
				//math.Exp(-math.Abs(3.5-float64(j))/7+0.25)
				dCell += bot.costMatrix[i][j] * math.Pow(1.1, float64(i))
			} else if game.Board[i][j] == 3 {
				dKing += bot.costMatrix[i][j]
			} else if game.Board[i][j] == 2 {
				dCell -= bot.costMatrix[i][j] * math.Pow(1.1, float64(7-i))
			} else if game.Board[i][j] == 4 {
				dKing -= bot.costMatrix[i][j]
			}
		}
	}
	return bot.cellCost*dCell + bot.kingCost*dKing
}

func (bot *Bot) _dfsStreak(game Board, me int, enemy int) Board {
	var maxGame = game
	for _, k := range POSSIBLE_TURNS {
		if game.canMove(game.Last_piece, _add(game.Last_piece, k)) {
			var _game = game
			_game.makeMove(_game.Last_piece, _add(_game.Last_piece, k))
			_game = bot._dfsStreak(_game, me, enemy)
			if game.Whose_turn == 0 {
				if bot.evaluate(_game) > bot.evaluate(maxGame) {
					maxGame = _game
				}
			} else {
				if bot.evaluate(_game) < bot.evaluate(maxGame) {
					maxGame = _game
				}
			}
		}
	}
	return maxGame
}

func (bot *Bot) _findBestMove(game Board, depth int, me int, enemy int, prev_score float64) float64 {
	//val, ok := bot.moves_table.get(newItemGame(game))
	//if ok {
	//	return val
	//}
	val, ok := bot.movesTable.get(game)
	if ok {
		return val
	}
	if depth == bot.maxDepth || game.isGameEnded() {
		return bot.evaluate(game)
	}
	var maxScore float64
	if game.Whose_turn == 0 {
		maxScore = -1e9
	} else {
		maxScore = 1e9
	}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if game.Board[i][j] != game.Whose_turn+1 && game.Board[i][j] != 2+game.Whose_turn+1 {
				continue
			}
			for _, k := range POSSIBLE_TURNS {
				if game.canMove([2]int{i, j}, _add([2]int{i, j}, k)) {
					var _game = game
					_game.makeMove([2]int{i, j}, _add([2]int{i, j}, k))
					_game = bot._dfsStreak(_game, me, enemy)
					_game.endMove()
					if game.Whose_turn == 0 {
						maxScore = math.Max(maxScore, bot._findBestMove(_game, depth+1, me, enemy, maxScore))
						if maxScore < prev_score {
							return maxScore
						}
					} else {
						maxScore = math.Min(maxScore, bot._findBestMove(_game, depth+1, me, enemy, maxScore))
						if maxScore > prev_score {
							return maxScore
						}
					}
				}
			}
		}
	}
	bot.movesTable.insert(game, maxScore)
	return maxScore
}

func (bot *Bot) gameTemp(game Board) float64 {
	if game.Whose_turn == 0 {
		return bot._findBestMove(game, 0, 0, 1, -1e9)
	} else {
		return bot._findBestMove(game, 0, 0, 1, 1e9)
	}
}

func _findBestMoveGoroutine(bot *Bot, game Board, me int, enemy int, move_chanel chan Move) {
	if game.Whose_turn == 0 {
		move_chanel <- newMove(bot._findBestMove(game, 1, me, enemy, -1e9), game)
	} else {
		move_chanel <- newMove(bot._findBestMove(game, 1, me, enemy, 1e9), game)
	}
}

func (bot *Bot) findBestMove(game Board, me int, enemy int) Board {
	var moveChanel = make(chan Move)
	var maxScore float64
	if game.Whose_turn == 0 {
		maxScore = -1e9
	} else {
		maxScore = 1e9
	}
	var cnt = 0
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if game.Board[i][j] != game.Whose_turn+1 && game.Board[i][j] != 2+game.Whose_turn+1 {
				continue
			}
			for _, k := range POSSIBLE_TURNS {
				if game.canMove([2]int{i, j}, _add([2]int{i, j}, k)) {
					_game := game
					_game.makeMove([2]int{i, j}, _add([2]int{i, j}, k))
					_game = bot._dfsStreak(_game, me, enemy)
					_game.endMove()
					cnt++
					go _findBestMoveGoroutine(bot, _game, me, enemy, moveChanel)
				}
			}
		}
	}
	var turns []Board
	for cnt > 0 {
		move := <-moveChanel
		if (game.Whose_turn == 0 && move.score > maxScore) || (game.Whose_turn == 1 && move.score < maxScore) {
			turns = make([]Board, 1)
			turns[0] = move.game
			maxScore = move.score
		} else if move.score == maxScore {
			turns = append(turns, move.game)
		}
		cnt--
	}
	bot.movesTable.clear()
	return turns[rand.Int()%len(turns)]
}

func (bot *Bot) makeMove(game *Game) {
	game.Board = bot.findBestMove(game.Board, game.Board.Whose_turn, (game.Board.Whose_turn+1)%2)
	game.Turns = append(game.Turns, game.Board.Board)
}

func _printBoard(game Board) {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			fmt.Print(game.Board[i][j])
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}

func botVsBot() {
	var game = newGame("", "")
	var cnt = 0
	for !game.isGameEnded() {
		game.Board = BOT.findBestMove(game.Board, cnt, (cnt+1)%2)
		fmt.Println(math.Round(BOT.gameTemp(game.Board)))
		_printBoard(game.Board)
		cnt++
		cnt %= 2
	}
	fmt.Println(game.whoWin())
}
