package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type bingoCard [25]uint8

func (b bingoCard) win(called map[uint8]bool) bool {
	for i := 0; i < 5; i++ {
		// did we win by row?
		var loser bool
		for j := 0; j < 5; j++ {
			if _, ok := called[b[i*5+j]]; !ok {
				loser = true
				break
			}
		}
		if !loser {
			return true
		}
		// did we win by column?
		loser = false
		for j := 0; j < 5; j++ {
			if _, ok := called[b[i+j*5]]; !ok {
				loser = true
				break
			}
		}
		if !loser {
			return true
		}
	}

	return false
}

type bingoCards struct {
	cards      []bingoCard
	called     map[uint8]bool
	lastCalled uint8
	winners    []int64
}

func newBingoCard(numCards int64) bingoCards {
	return bingoCards{cards: make([]bingoCard, 0, numCards), called: map[uint8]bool{}}
}

func (b *bingoCards) win() {
	var winner *int
	for i, card := range b.cards {
		if card.win(b.called) {
			winner = &i
			break
		}
	}

	if winner != nil {
		b.winners = append(b.winners, b.score(b.cards[*winner]))

		b.cards = append(b.cards[:*winner], b.cards[*winner+1:]...)
		// allow for simultaneous winners
		b.win()
	}
}

func (b *bingoCards) call(i uint8) {
	b.called[i] = true
	b.lastCalled = i
	b.win()
}

func (b *bingoCards) add(card bingoCard) {
	b.cards = append(b.cards, card)
}

func (b bingoCards) empty() bool {
	return len(b.cards) == 0
}

func (b bingoCards) score(card bingoCard) int64 {
	var sum int64
	for _, i := range card {
		if !b.called[i] {
			sum += int64(i)
		}
	}
	return sum * int64(b.lastCalled)
}

func (b bingoCards) lastWinner() *int64 {
	if len(b.winners) == 0 {
		return nil
	}
	return &b.winners[len(b.winners)-1]
}

func (b bingoCards) firstWinner() *int64 {
	if len(b.winners) == 0 {
		return nil
	}
	return &b.winners[0]
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

	scanner := bufio.NewScanner(f)

	if !scanner.Scan() {
		log.Fatal("Not enough input to draw numbers!")
	}
	numbersToDraw := strings.Split(scanner.Text(), ",")

	var bingoCardNumbers []uint8
	for scanner.Scan() {
		if scanner.Text() == "" {
			continue
		}

		var b, i, n, g, o uint8
		if n, err := fmt.Sscan(scanner.Text(), &b, &i, &n, &g, &o); n != 5 || err != nil {
			log.Fatal(err)
		}

		bingoCardNumbers = append(bingoCardNumbers, b, i, n, g, o)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if len(bingoCardNumbers)%25 != 0 {
		log.Fatal("Unexpected bingo card, not a full card: \n%q", bingoCardNumbers[len(bingoCardNumbers)-25])
	}

	bingoCards := newBingoCard(int64(len(bingoCardNumbers) / 25))
	for i := 0; i < len(bingoCardNumbers); i += 25 {
		var card [25]uint8
		if copy(card[:], bingoCardNumbers[i:i+25]) != 25 {
			log.Fatal("Did not copy slice->arr correctly")
		}
		bingoCards.add(card)
	}

	for i := 0; i < len(numbersToDraw) && !bingoCards.empty(); i++ {
		n := numbersToDraw[i]
		var i uint8
		if n, err := fmt.Sscan(n, &i); n != 1 || err != nil {
			log.Fatal(err)
		}

		bingoCards.call(i)
	}
	firstWinner := bingoCards.firstWinner()
	if firstWinner == nil {
		log.Fatal("Expected a winner!")
	}
	log.Printf("Part 1: %d", *firstWinner)

	// reasonably sure this won't panic, because of line 170!
	lastWinner := bingoCards.lastWinner()
	log.Printf("Part 2: %d", *lastWinner)
}
