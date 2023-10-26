package snake

import (
	"github.com/pechorka/htmx-snake/pkg/enums/direction"
)

type Snake struct {
	head      *snakePart
	tail      *snakePart
	direction direction.Direction
	borders   [2]int
}

func NewSnake(head [2]int, length int, direction direction.Direction, borders [2]int) *Snake {
	headPart := &snakePart{
		part:     PartHead,
		location: head,
	}

	tailPart := headPart
	for i := 1; i < length; i++ {
		tailPart.next = &snakePart{
			part:     PartBody,
			location: [2]int{head[0], head[1] + i},
			prev:     tailPart,
		}
		tailPart = tailPart.next
	}

	return &Snake{
		head:      headPart,
		tail:      tailPart,
		direction: direction,
		borders:   borders,
	}
}

func (s *Snake) Move(to direction.Direction) {
	// calculate new head location
	newHead := s.head.location
	switch to {
	case direction.Up:
		newHead[0]--
		if newHead[0] < 0 {
			newHead[0] = s.borders[0] - 1
		}
	case direction.Down:
		newHead[0]++
		if newHead[0] >= s.borders[0] {
			newHead[0] = 0
		}
	case direction.Left:
		newHead[1]--
		if newHead[1] < 0 {
			newHead[1] = s.borders[1] - 1
		}
	case direction.Right:
		newHead[1]++
		if newHead[1] >= s.borders[1] {
			newHead[1] = 0
		}
	}

	// make new head
	curHead := s.head
	s.head = &snakePart{
		part:     PartHead,
		location: newHead,
		next:     s.head,
	}
	// make previous head a body
	curHead.part = PartBody
	curHead.prev = s.head
	s.direction = to

	// remove tail
	s.tail = s.tail.prev
	s.tail.next = nil
}

func (s *Snake) CantMove(newDirection direction.Direction) bool {
	switch s.direction {
	case direction.Up:
		return newDirection == direction.Down
	case direction.Down:
		return newDirection == direction.Up
	case direction.Left:
		return newDirection == direction.Right
	case direction.Right:
		return newDirection == direction.Left
	default:
		return false
	}
}

func (s *Snake) Iterate(fn func(loc [2]int, bodyPart Part)) {
	for part := s.head; part != nil; part = part.next {
		fn(part.location, part.part)
	}
}

func (s *Snake) Direction() direction.Direction {
	return s.direction
}

type Part int

const (
	PartHead Part = 0
	PartBody Part = 1
)

type snakePart struct {
	part     Part
	location [2]int
	next     *snakePart
	prev     *snakePart
}
