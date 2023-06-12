package game

import (
	"math/rand"
	"minesweeper/types"
)

type ClientId = types.ClientId

type Counter struct {
	Score     map[ClientId]uint `json:"score"`
	BombsLeft uint              `json:"BombsLeft"`
}

type Turn struct {
	Count uint
	Curr  ClientId
	Next  ClientId
}

type Game struct {
	Counter Counter
	Turn    Turn
	Board   [][]*Block
	Winner  ClientId
}

const DEFAULT_SIZE = 26
const DEFAULT_BOMB_COUNT = 100

func CreateGame() *Game {
	return &Game{
		Counter: Counter{},
		Turn: Turn{
			Curr: "",
			Next: "",
		},
		Board: GetBoard(DEFAULT_SIZE, DEFAULT_BOMB_COUNT),
	}
}

func (g *Game) InitCounter() {
	g.Counter.BombsLeft = DEFAULT_BOMB_COUNT
	g.Counter.Score = make(map[ClientId]uint)
	g.Counter.Score[g.Turn.Curr] = 0
	g.Counter.Score[g.Turn.Next] = 0
}

func (g *Game) GetCounter() Counter {
	return g.Counter
}

func (g *Game) GetTurn() Turn {
	return g.Turn
}
func (g *Game) AssignTurn(cId ClientId) (ClientId, ClientId) {
	isEmpty := g.Turn.Curr == "" && g.Turn.Next == ""

	if isEmpty {
		if rand.Intn(2) == 0 {
			g.Turn.Curr = cId
		} else {
			g.Turn.Next = cId
		}
	} else if g.Turn.Curr == "" {
		g.Turn.Curr = cId
	} else if g.Turn.Next == "" { // Not 'else'. Could be more than 3 clients in a room
		g.Turn.Next = cId
	}
	return g.Turn.Curr, g.Turn.Next
}

func (g *Game) UnassignTurn(cId ClientId) (ClientId, ClientId) {
	if g.Turn.Curr == cId {
		g.Turn.Curr = ""
	} else if g.Turn.Next == cId {
		g.Turn.Next = ""
	}
	return g.Turn.Curr, g.Turn.Next
}

func (g *Game) AdvanceTurn() Turn {
	g.Turn.Count++
	g.Turn.Curr, g.Turn.Next = g.Turn.Next, g.Turn.Curr
	return g.Turn
}

func (g *Game) ScoreCurrPlayer() (Counter, bool) {
	g.Counter.BombsLeft--
	g.Counter.Score[g.Turn.Curr]++

	isWon := g.Counter.Score[g.Turn.Curr] > DEFAULT_BOMB_COUNT/2
	if isWon {
		g.Winner = g.Turn.Curr
	}
	return g.Counter, isWon
}

func (g *Game) GetWinner() ClientId {
	return g.Winner
}

func (g *Game) getVisibleBlocks() []BlockInfo {
	s := []BlockInfo{}
	for i := range g.Board {
		for j := range g.Board {
			if g.Board[i][j].Visited {
				s = append(s, BlockInfo{
					X:     j,
					Y:     i,
					Block: *g.Board[i][j],
				})
			}
		}
	}
	return s
}
