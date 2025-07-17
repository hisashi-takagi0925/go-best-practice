package main

import (
	"fmt"
	"math/rand"
)

type Hand int

const (
	Rock Hand = iota
	Paper
	Scissors
)

func (h Hand) String() string {
	switch h {
	case Rock:
		return "👊"
	case Paper:
		return "✋"
	case Scissors:
		return "✌️"
	default:
		return "?"
	}
}

func (h Hand) Beats(other Hand) bool {
	return (h == Rock && other == Scissors) ||
		(h == Paper && other == Rock) ||
		(h == Scissors && other == Paper)
}

type Player struct {
	name string
}

func NewPlayer(name string) *Player {
	return &Player{name: name}
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) ChooseHand() Hand {
	fmt.Printf("%sの手を選択してください\n", p.name)
	fmt.Println("1: 👊 (グー)")
	fmt.Println("2: ✋ (パー)")
	fmt.Println("3: ✌️ (チョキ)")

	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		return Rock
	case 2:
		return Paper
	case 3:
		return Scissors
	default:
		fmt.Println("無効な選択です。もう一度選んでください")
		return p.ChooseHand()
	}
}

func (p *Player) ChooseRandomHand() Hand {
	return Hand(rand.Intn(3))
}

type GameResult int

const (
	Win GameResult = iota
	Lose
	Draw
)

func (gr GameResult) String() string {
	switch gr {
	case Win:
		return "あなたの勝ちです"
	case Lose:
		return "あなたの負けです"
	case Draw:
		return "引き分け"
	default:
		return "不明"
	}
}

func Judge(playerHand, computerHand Hand) GameResult {
	if playerHand == computerHand {
		return Draw
	}
	if playerHand.Beats(computerHand) {
		return Win
	}
	return Lose
}

type Record struct {
	totalGames int
	wins       int
	losses     int
	draws      int
}

func NewRecord() *Record {
	return &Record{}
}

func (r *Record) AddResult(result GameResult) {
	r.totalGames++
	switch result {
	case Win:
		r.wins++
	case Lose:
		r.losses++
	case Draw:
		r.draws++
	}
}

func (r *Record) TotalGames() int {
	return r.totalGames
}

func (r *Record) WinRate() float64 {
	if r.totalGames == 0 {
		return 0.0
	}
	return float64(r.wins) / float64(r.totalGames) * 100
}

func (r *Record) Summary() string {
	return fmt.Sprintf("試合数: %d, 勝ち: %d, 負け: %d, 引き分け: %d, 勝率: %.2f%%",
		r.totalGames, r.wins, r.losses, r.draws, r.WinRate())
}

func main() {

	player := NewPlayer("あなた")
	computer := NewPlayer("コンピューター")
	record := NewRecord()

	fmt.Println("じゃんけんゲームを開始します！")

	for {
		playerHand := player.ChooseHand()
		computerHand := computer.ChooseRandomHand()

		fmt.Printf("%s: %s\n", player.Name(), playerHand)
		fmt.Printf("%s: %s\n", computer.Name(), computerHand)

		result := Judge(playerHand, computerHand)
		fmt.Println("結果:", result)

		record.AddResult(result)

		fmt.Println("続けますか？ (y/n)")
		var cont string
		fmt.Scanln(&cont)
		if cont != "y" {
			break
		}
	}

	fmt.Println("\n=== 試合結果 ===")
	fmt.Println(record.Summary())
}