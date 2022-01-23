package main

import (
	"log"
	"os"
)

type Die interface {
	// [1,100]
	Roll() uint8
	NumRolls() uint
}

type deterministicDie struct {
	lastRoll uint8
	rolls    uint
}

func NewDeterministicDie() Die {
	return new(deterministicDie)
}

func (d *deterministicDie) Roll() uint8 {
	d.lastRoll = (d.lastRoll % 100) + 1
	d.rolls++
	return d.lastRoll
}
func (d *deterministicDie) NumRolls() uint {
	return d.rolls
}

type Pawn interface {
	//Move i positions relative to current position
	Move(i uint8)
	// MoveTo particular position on board
	MoveTo(i uint8)
	// current position
	Position() uint8
}
type pawn struct {
	pos uint8
}

func NewPawn(startingPosition uint8) Pawn {
	return &pawn{startingPosition}
}

func (p *pawn) Move(i uint8) {
	p.MoveTo(p.pos + i)
}
func (p *pawn) MoveTo(i uint8) {
	p.pos = i
}
func (p pawn) Position() uint8 {
	return p.pos
}

type Board interface {
	Move(p Pawn, i uint8)
}
type board struct {
	ps   []Pawn
	size uint8
}

func NewBoard(pawns []Pawn) Board {
	return &board{pawns, 10}
}

func (b board) Move(p Pawn, i uint8) {
	for _, pown := range b.ps {
		if pown == p {
			// we need to think of this as [1,10]
			// not [0,9]
			for j := uint8(0); j < i; j++ {
				p.Move(1)
				if p.Position() >= b.size+1 {
					p.MoveTo(p.Position() - b.size)
				}
			}
			return
		}
	}
	panic("panik")
}

type Player interface {
	Turn(d Die, b Board)
	Score() uint16
}
type player struct {
	p     Pawn
	score uint16
}

func NewPlayer(p Pawn) Player {
	return &player{p: p}
}

func (p *player) Turn(d Die, b Board) {
	for i := 0; i < 3; i++ {
		roll := d.Roll()
		b.Move(p.p, roll)
	}
	p.score += uint16(p.p.Position())
}
func (p player) Score() uint16 {
	return p.score
}

type Game interface {
	Loser() *Player
	PlayTurn() bool
}
type game struct {
	scoreToWin          uint16
	players             []Player
	board               Board
	die                 Die
	nextTurnPlayerIndex uint8
}

func NewGame(players []Player, board Board, d Die) Game {
	return &game{1000, players, board, d, 0}
}

func (g game) Loser() *Player {
	var winner *Player
	for _, p := range g.players {
		if p.Score() >= g.scoreToWin {
			winner = &p
			break
		}
	}

	if winner == nil {
		return nil
	}

	for _, p := range g.players {
		if p != *winner {
			return &p
		}
	}
	return nil
}
func (g *game) PlayTurn() bool {
	g.players[g.nextTurnPlayerIndex].Turn(g.die, g.board)
	g.nextTurnPlayerIndex = (g.nextTurnPlayerIndex + 1) % uint8(len(g.players))
	return g.Loser() == nil
}

func main() {
	var f *os.File
	if len(os.Args) == 2 {
		var err error
		f, err = os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
	} else {
		f = os.Stdin
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	/*
		scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				_ = scanner.Text()
			}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	*/

	// pawns := []Pawn{NewPawn(4), NewPawn(8)}
	pawns := []Pawn{NewPawn(10), NewPawn(1)}
	d := NewDeterministicDie()
	g := NewGame([]Player{NewPlayer(pawns[0]), NewPlayer(pawns[1])}, NewBoard(pawns), d)
	for g.PlayTurn() {
	}
	log.Printf("Part 1: %d", uint((*g.Loser()).Score())*d.NumRolls())
}
