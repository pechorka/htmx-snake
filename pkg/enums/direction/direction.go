package direction

type Direction string

const (
	Up    Direction = "ArrowUp"
	Down  Direction = "ArrowDown"
	Left  Direction = "ArrowLeft"
	Right Direction = "ArrowRight"
)

func (d Direction) IsValid() bool {
	switch d {
	case Up, Down, Left, Right:
		return true
	}
	return false
}
