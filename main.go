package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Bet represents a single bet in roulette
type Bet struct {
	Type   string
	Value  int
	Amount float64
}

// Strategy represents a roulette betting strategy
type Strategy struct {
	InitialBankroll float64
	Bets            []Bet
}

// RouletteWheel represents the roulette wheel
type RouletteWheel struct {
	Numbers []int
}

// NewRouletteWheel creates a new roulette wheel
func NewRouletteWheel() *RouletteWheel {
	numbers := make([]int, 38)
	for i := 0; i < 36; i++ {
		numbers[i] = i + 1
	}
	numbers[36] = 0  // Green 0
	numbers[37] = 00 // Green 00
	return &RouletteWheel{Numbers: numbers}
}

// Spin spins the roulette wheel and returns the winning number
func (rw *RouletteWheel) Spin() int {
	return rw.Numbers[rand.Intn(len(rw.Numbers))]
}

// ParseStrategy parses the DSL input and returns a Strategy
func ParseStrategy(input string) (*Strategy, error) {
	lines := strings.Split(input, "\n")
	strategy := &Strategy{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "bankroll:") {
			bankrollStr := strings.TrimPrefix(line, "bankroll:")
			bankroll, err := strconv.ParseFloat(strings.TrimSpace(bankrollStr), 64)
			if err != nil {
				return nil, fmt.Errorf("invalid bankroll: %v", err)
			}
			strategy.InitialBankroll = bankroll
		} else if strings.HasPrefix(line, "bet:") {
			betStr := strings.TrimPrefix(line, "bet:")
			parts := strings.Split(betStr, ",")
			if len(parts) != 3 {
				return nil, fmt.Errorf("invalid bet format: %s", line)
			}
			betType := strings.TrimSpace(parts[0])
			betValue, err := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid bet value: %v", err)
			}
			betAmount, err := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
			if err != nil {
				return nil, fmt.Errorf("invalid bet amount: %v", err)
			}
			strategy.Bets = append(strategy.Bets, Bet{Type: betType, Value: betValue, Amount: betAmount})
		}
	}

	return strategy, nil
}

// SimulateRoulette simulates roulette games using the given strategy
func SimulateRoulette(strategy *Strategy, numGames int) float64 {
	wheel := NewRouletteWheel()
	bankroll := strategy.InitialBankroll

	for i := 0; i < numGames; i++ {
		winningNumber := wheel.Spin()

		for _, bet := range strategy.Bets {
			if bankroll < bet.Amount {
				continue // Skip this bet if we don't have enough money
			}

			bankroll -= bet.Amount

			switch bet.Type {
			case "number":
				if bet.Value == winningNumber {
					bankroll += bet.Amount * 36
				}
			case "even":
				if winningNumber%2 == 0 && winningNumber != 0 {
					bankroll += bet.Amount * 2
				}
			case "odd":
				if winningNumber%2 != 0 && winningNumber != 0 {
					bankroll += bet.Amount * 2
				}
			case "red":
				redNumbers := []int{1, 3, 5, 7, 9, 12, 14, 16, 18, 19, 21, 23, 25, 27, 30, 32, 34, 36}
				if contains(redNumbers, winningNumber) {
					bankroll += bet.Amount * 2
				}
			case "black":
				blackNumbers := []int{2, 4, 6, 8, 10, 11, 13, 15, 17, 20, 22, 24, 26, 28, 29, 31, 33, 35}
				if contains(blackNumbers, winningNumber) {
					bankroll += bet.Amount * 2
				}
			}
		}
	}

	return bankroll
}

// contains checks if a slice contains a specific value
func contains(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Enter your roulette strategy (type 'done' on a new line when finished):")
	var input strings.Builder
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "done" {
			break
		}
		input.WriteString(line + "\n")
	}

	strategy, err := ParseStrategy(input.String())
	if err != nil {
		fmt.Printf("Error parsing strategy: %v\n", err)
		return
	}

	fmt.Print("Enter the number of games to simulate: ")
	scanner.Scan()
	numGames, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Printf("Invalid number of games: %v\n", err)
		return
	}

	finalBankroll := SimulateRoulette(strategy, numGames)
	fmt.Printf("Initial bankroll: $%.2f\n", strategy.InitialBankroll)
	fmt.Printf("Final bankroll after %d games: $%.2f\n", numGames, finalBankroll)
	fmt.Printf("Profit/Loss: $%.2f\n", finalBankroll-strategy.InitialBankroll)
}
