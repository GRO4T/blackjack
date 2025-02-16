package blackjack

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strconv"

	"github.com/GRO4T/bjack-api/constant"
	"github.com/GRO4T/bjack-api/deck"
)

var (
	ErrNotFound           = errors.New("not found")
	ErrGameIsFull         = errors.New("game is full")
	ErrGameAlreadyStarted = errors.New("game already started")
	ErrCardsAlreadyDealt  = errors.New("cards already dealt")
	ErrGameNotInProgress  = errors.New("game not in progress")
	ErrOtherPlayerTurn    = errors.New("other player's turn")
)

const (
	initialChips = 100
)

type State int

const (
	WaitingForPlayers State = iota
	CardsDealt
	Finished
)

type Outcome int

const (
	Undecided Outcome = iota
	Win
	Lose
	Push
)

type Action int

const (
	Hit Action = iota
	Stand
)

type Player struct {
	Id      string  `json:"-"`
	Name    string  `json:"name"`
	IsReady bool    `json:"isReady"`
	Chips   int     `json:"chips"`
	Bet     int     `json:"bet"`
	Outcome Outcome `json:"outcome"`
}

type Blackjack struct {
	Deck           []deck.Card   `json:"-"`
	Hands          [][]deck.Card `json:"hands"`
	Players        []*Player     `json:"players"`
	State          State         `json:"state"`
	CurrentPlayer  int           `json:"currentPlayer"`
	onStateChanged func()        `json:"-"`
}

func NewPlayer(id string, name string) Player {
	return Player{
		Id:      id,
		Name:    name,
		IsReady: false,
		Chips:   initialChips,
		Bet:     0,
		Outcome: Undecided,
	}
}

func New(onStateChanged func()) Blackjack {
	dealerHand := []deck.Card{}
	return Blackjack{
		Deck:           deck.New(deck.WithShuffle()),
		Hands:          [][]deck.Card{dealerHand},
		Players:        []*Player{},
		State:          WaitingForPlayers,
		CurrentPlayer:  1,
		onStateChanged: onStateChanged,
	}
}

func (b *Blackjack) AddPlayer(name string) (*Player, error) {
	if b.State != WaitingForPlayers {
		return nil, ErrGameAlreadyStarted
	}
	if len(b.Players) >= constant.MaxPlayers {
		return nil, ErrGameIsFull
	}
	for _, player := range b.Players {
		if player.Name == name {
			return nil, errors.New("Player with name " + name + " already exists")
		}
	}
	newPlayer := NewPlayer(getRandomId(), name)
	b.Players = append(b.Players, &newPlayer)
	b.Hands = append(b.Hands, []deck.Card{})
	if b.onStateChanged != nil {
		b.onStateChanged()
	}
	return &newPlayer, nil
}

func (b *Blackjack) TogglePlayerReady(playerId string) (*Player, error) {
	if b.State != WaitingForPlayers {
		return nil, ErrGameAlreadyStarted
	}

	var targetPlayer *Player = nil
	allPlayersReady := true
	for _, player := range b.Players {
		if player.Id == playerId {
			player.IsReady = !player.IsReady
			targetPlayer = player
		}
		if !player.IsReady {
			allPlayersReady = false
		}
	}

	if targetPlayer == nil {
		return nil, ErrNotFound
	}

	if allPlayersReady {
		err := b.Deal()
		if err != nil {
			return nil, err
		}
		b.State = CardsDealt
	}

	if b.onStateChanged != nil {
		b.onStateChanged()
	}

	return targetPlayer, nil
}

func (b *Blackjack) Deal() error {
	if b.State == CardsDealt {
		return ErrCardsAlreadyDealt
	}
	for range 2 {
		for player := 0; player < len(b.Hands); player++ {
			b.Hands[player] = append(b.Hands[player], b.Deck[0])
			b.Deck = b.Deck[1:]
		}
	}
	return nil
}

func (b *Blackjack) GetPlayerCount() int {
	return len(b.Hands) - 1
}

func (b *Blackjack) GetPlayerHand(playerNumber int) []deck.Card {
	return b.Hands[playerNumber+1]
}

func (b *Blackjack) GetDealerHand() []deck.Card {
	return b.Hands[0]
}

func (b *Blackjack) PlayerAction(playerId string, action Action) error {
	if b.State != CardsDealt {
		return ErrGameNotInProgress
	}

	playerIndex := -1
	for i, p := range b.Players {
		if p.Id == playerId {
			playerIndex = i
			break
		}
	}
	if playerIndex == -1 {
		return ErrNotFound
	}

	if playerIndex+1 != b.CurrentPlayer {
		return ErrOtherPlayerTurn
	}

	switch action {
	case Hit:
		b.Hands[playerIndex+1] = append(b.Hands[playerIndex+1], b.Deck[0])
		b.Deck = b.Deck[1:]
		b.CurrentPlayer++
	case Stand:
		b.CurrentPlayer++
	}

	// Dealer's turn
	if b.CurrentPlayer == len(b.Hands) { // TODO(fix): CurrentPlayer should start from 0
		// TODO(feat): Implement dealer AI
		b.State = Finished
		b.DetermineOutcomes()
	}

	return nil
}

func (b *Blackjack) DetermineOutcomes() {
	if b.State != Finished {
		return
	}
	for i := range b.GetPlayerCount() {
		b.Players[i].Outcome = determineOutcome(b.GetDealerHand(), b.GetPlayerHand(i))
	}
}

// nolint: mnd
func determineOutcome(dealerHand []deck.Card, playerHand []deck.Card) Outcome {
	dealerScore := getScore(dealerHand)
	playerScore := getScore(playerHand)
	dealerHasBlackjack := isBlackjack(dealerHand)
	playerHasBlackjack := isBlackjack(playerHand)

	if playerHasBlackjack && dealerHasBlackjack {
		return Push
	}
	if playerHasBlackjack {
		return Win
	}
	if dealerHasBlackjack {
		return Lose
	}
	if playerScore > 21 {
		return Lose
	}
	if dealerScore > 21 {
		return Win
	}
	if playerScore > dealerScore {
		return Win
	}
	if playerScore < dealerScore {
		return Lose
	}
	return Push
}

func getScore(hand []deck.Card) int {
	score := 0
	aceCount := 0
	for _, card := range hand {
		if card.Rank == deck.Ace {
			aceCount++
		}
		score += int(card.Rank)
	}
	for aceCount > 0 {
		if score > 21 { //nolint: mnd
			score -= 10
			aceCount--
		} else {
			break
		}
	}
	return score
}

func isBlackjack(hand []deck.Card) bool {
	if len(hand) != 2 { //nolint: mnd
		return false
	}
	isAce := func(i int) bool {
		return hand[i].Rank == deck.Ace
	}
	isTenOrQKJ := func(i int) bool {
		return hand[i].Rank == deck.Ten ||
			hand[i].Rank == deck.Jack ||
			hand[i].Rank == deck.Queen ||
			hand[i].Rank == deck.King
	}
	return (isAce(0) && isTenOrQKJ(1)) || (isAce(1) && isTenOrQKJ(0))
}

func getRandomId() string {
	id, err := rand.Int(rand.Reader, big.NewInt(constant.MaxId))
	if err != nil {
		panic(err)
	}
	return strconv.Itoa(int(id.Int64()))
}
