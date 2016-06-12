// An example implementation the qlearning interfaces. Can be run
// with go run hangman.go.
//
// Word list provided by https://github.com/first20hours/google-10000-english
package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"

	"github.com/ecooper/qlearning"
)

const (
	startingLives = 6

	Lost   = -1
	Active = 0
	Won    = 1
)

var (
	Alphabet string   = "abcdefghijklmnopqrstuvwxyz"
	WordList []string = make([]string, 0)

	wordListPath string = "./wordlist.txt"
	debug        bool   = false
	progressAt   int    = 1000
	wordCount    int    = 10000
	playFor      int    = 5000000
)

func loadWords() error {
	f, err := os.Open(wordListPath)
	if err != nil {
		return err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if len(WordList) >= wordCount {
			break
		}
		WordList = append(WordList, scanner.Text())
	}

	return nil
}

// Game represents the state of any given game of Hangman. It implements
// qlearning.Agent, qlearning.Rewarder, and qlearning.State.
type Game struct {
	Word          string
	Characters    int
	StartingLives int
	Lives         int
	Attempted     map[string]bool
	Correct       []string

	debug bool
}

// NewGame creates a new Hangman game for the given word. If debug
// is true, Game.Log messages will print to stdout.
func NewGame(word string, debug bool) *Game {
	game := &Game{debug: debug}
	game.New(word)

	return game
}

// NewWord returns a random word from Wordlist.
func NewWord() string {
	return WordList[rand.Intn(len(WordList))]
}

// New resets the current game to a new game for the given word.
func (game *Game) New(word string) {
	game.Word = word
	game.Characters = len(word)
	game.StartingLives = startingLives
	game.Lives = startingLives
	game.Attempted = make(map[string]bool, len(Alphabet))
	game.Correct = make([]string, len(word), len(word))
}

// Returns Lost, Active, or Won based on the game's current state.
func (game *Game) IsComplete() int {
	if game.Lives < 1 {
		return Lost
	}

	if game.Characters > 0 {
		return Active
	}

	return Won
}

// Choose applies a character attempt in the current game, returning
// true if char is present in Game.Word.
//
// Choose updates the game's state.
func (game *Game) Choose(char string) bool {
	game.Attempted[char] = true

	hit := false

	for i, check := range game.Word {
		if string(check) == char {
			game.Correct[i] = char
			game.Characters -= 1
			hit = true
		}
	}

	if !hit {
		game.Lives -= 1
		return false
	}

	return true
}

// Reward returns a score for a given qlearning.StateAction. Reward is a
// member of the qlearning.Rewarder interface. If the choice is found in
// the game's word, a positive score is returned. Otherwise, a static
// -1000 is returned.
func (game *Game) Reward(action *qlearning.StateAction) float32 {
	choice := action.Action.String()
	for _, char := range game.Word {
		if string(char) == choice {
			return 24.0 / float32(len(game.Attempted))
		}
	}

	return -1000
}

// Next creates a new slice of qlearning.Action instances. A possible
// action is created for each character that has not been attempted in
// in the game.
func (game *Game) Next() []qlearning.Action {
	actions := make([]qlearning.Action, 0, len(Alphabet))

	for _, char := range Alphabet {
		attempted := game.Attempted[string(char)]
		if !attempted {
			actions = append(actions, &Choice{Character: string(char)})
		}
	}

	return actions
}

// Log is a wrapper of fmt.Printf. If Game.debug is true, Log will print
// to stdout.
func (game *Game) Log(msg string, args ...interface{}) {
	if game.debug {
		logMsg := fmt.Sprintf("[GAME %s] (%d moves, %d lives) %s\n", game.Word, len(game.Attempted), game.Lives, msg)
		fmt.Printf(logMsg, args...)
	}
}

// String returns a consistent hash for the current game state to be
// used in a qlearning.Agent.
func (game *Game) String() string {
	return fmt.Sprintf("%s", game.Correct)
}

// Choice implements qlearning.Action for a character choice in a game
// of Hangman.
type Choice struct {
	Character string
}

// String returns the character for the current action.
func (choice *Choice) String() string {
	return choice.Character
}

// Apply updates the state of the game for a given character choice.
func (choice *Choice) Apply(state qlearning.State) qlearning.State {
	game := state.(*Game)
	game.Choose(choice.Character)

	return game
}

func init() {
	flag.StringVar(&wordListPath, "wordlist", wordListPath, "Path to a wordlist")
	flag.BoolVar(&debug, "debug", debug, "Set debug")
	flag.IntVar(&progressAt, "progress", progressAt, "Print progress messages every N games")
	flag.IntVar(&wordCount, "words", wordCount, "Use N words from wordlist")
	flag.IntVar(&playFor, "games", playFor, "Play N games")

	flag.Parse()

	loadWords()
	fmt.Printf("%d words loaded\n", len(WordList))
}

func main() {
	var (
		wins     = 0
		lastWins = 0
		count    = 0

		// Our agent has a learning rate of 0.7 and discount of 1.0.
		agent = qlearning.NewSimpleAgent(0.7, 1.0)
	)

	progress := func() {
		// Print our progress every 1000 rows.
		if count > 0 && count%progressAt == 0 {
			accuracy := float32(wins-lastWins) / float32(progressAt) * 100.0
			lastWins = wins
			fmt.Printf("%d games played: %d WINS %d LOSSES %.0f ACCURACY\n", count, wins, count-wins, accuracy)
		}
	}

	// Let's play 5 million games
	for count = 0; count < playFor; count++ {
		// Get a new word and game for each iteration...
		word := NewWord()
		game := NewGame(word, debug)

		game.Log("Game created")

		// While the game is still active, we'll continue to update
		// our agent and learn from its choices.
		for game.IsComplete() == 0 {
			// Pick the next move, which is going to be a letter choice.
			action := qlearning.Next(agent, game)

			// Whatever that choice is, let's update our model for its
			// impact. If the character chosen is in the game's word,
			// then this action will be positive. Otherwise, it will be
			// negative.
			agent.Learn(action, game)

			// Reward doesn't change state so we can check what the
			// reward would be for this action, and report how the
			// game changed.
			if game.Reward(action) > 0.0 {
				game.Log("%s was correct", action.Action.String())
			} else {
				game.Log("%s was incorrect", action.Action.String())
			}
		}

		// If we won the game, record it as a victory.
		if game.IsComplete() == Won {
			game.Log("Victory!")
			wins += 1
		} else {
			game.Log("Defeat!")
		}

		progress()
	}

	progress()

	fmt.Printf("\nAgent performance: %d games played, %d WINS %d LOSSES %.0f ACCURACY\n", count, wins, count-wins, float32(wins)/float32(count)*100.0)
}
