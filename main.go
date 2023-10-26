package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/pechorka/htmx-snake/internal/snake"
	"github.com/pechorka/htmx-snake/pkg/enums/direction"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}

}

func run() error {
	templates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		return fmt.Errorf("error parsing templates: %w", err)
	}

	mux := chi.NewRouter()
	borders := [2]int{10, 50}
	initialSnakeHead := [2]int{5, 20}
	initialSnakeLength := 5
	snake := snake.NewSnake(initialSnakeHead, initialSnakeLength, direction.Left, borders)
	state := newGameState(borders, snake)
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if err := templates.ExecuteTemplate(w, "index.html", state); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.Post("/tick", func(w http.ResponseWriter, r *http.Request) {
		newDirection := r.FormValue("lastKey")
		slog.Info("request info", "lastKey", newDirection)
		err := state.MoveSnake(direction.Direction(newDirection))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := templates.ExecuteTemplate(w, "game-container", state); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	return http.ListenAndServe(":8080", mux)
}

type gameState struct {
	mu    *sync.Mutex
	Score int
	Board [][]cellType
	snake *snake.Snake
}

func newGameState(borders [2]int, snake *snake.Snake) *gameState {
	board := newBoard(borders[0], borders[1])
	state := &gameState{
		mu:    &sync.Mutex{},
		Score: 0,
		Board: board,
		snake: snake,
	}
	state.drawSnake()
	return state
}

type cellType string

const (
	emptyCell cellType = "X"
	snakeHead cellType = "🐔"
	snakeBody cellType = "🧠"
)

func newBoard(row, column int) [][]cellType {
	board := make([][]cellType, row)
	for i := range board {
		board[i] = make([]cellType, column)
		for j := range board[i] {
			board[i][j] = emptyCell
		}
	}
	return board
}

func (g *gameState) MoveSnake(to direction.Direction) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !to.IsValid() {
		to = g.snake.Direction()
	}

	if g.snake.CantMove(to) {
		return fmt.Errorf("can't move in to opposite direction")
	}

	g.eraseSnake()
	g.snake.Move(to)
	g.drawSnake()

	return nil
}

func (g *gameState) eraseSnake() {
	g.snake.Iterate(func(loc [2]int, bodyPart snake.Part) {
		g.Board[loc[0]][loc[1]] = emptyCell
	})
}

func (g *gameState) drawSnake() {
	g.snake.Iterate(func(loc [2]int, bodyPart snake.Part) {
		switch bodyPart {
		case snake.PartHead:
			g.Board[loc[0]][loc[1]] = snakeHead
		case snake.PartBody:
			g.Board[loc[0]][loc[1]] = snakeBody
		}
	})
}
