package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}

}

type gameState struct {
	mu    *sync.Mutex
	Score int
	Board [][]gameCell // don't modify directly, use Modify
}

func (g *gameState) Modify(mod func(*gameState)) *gameState {
	g.mu.Lock()
	mod(g)
	gCopy := *g
	g.mu.Unlock()
	return &gCopy
}

type gameCell struct {
	Symbol string
}

func run() error {
	templates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		return fmt.Errorf("error parsing templates: %w", err)
	}

	mux := chi.NewRouter()
	state := &gameState{
		mu:    &sync.Mutex{},
		Score: 10,
		Board: newBoard(50, 10),
	}
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if err := templates.ExecuteTemplate(w, "index.html", state); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.Post("/tick", func(w http.ResponseWriter, r *http.Request) {
		state = state.Modify(func(s *gameState) {
			s.Score++
		})
		if err := templates.ExecuteTemplate(w, "game-container", state); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	return http.ListenAndServe(":8080", mux)
}

func newBoard(x, y int) [][]gameCell {
	board := make([][]gameCell, 0, y)
	for i := 0; i < y; i++ {
		row := make([]gameCell, 0, x)
		for j := 0; j < x; j++ {
			row = append(row, newCell())
		}
		board = append(board, row)
	}
	return board
}

func newCell() gameCell {
	return gameCell{
		Symbol: "X",
	}
}
