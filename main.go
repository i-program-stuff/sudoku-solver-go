package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const BOARD_WIDTH = 9
const BOARD_HEIGHT = 9

const NUMBERS_BITMAP_MASK = 0b111111111 // 2^9 - 1

type SudokuBoard struct {
	Board [BOARD_HEIGHT][BOARD_WIDTH]byte
}

func NewSudokuBoard() *SudokuBoard {
	return &SudokuBoard{}
}

// func (b *SudokuBoard) Copy() *SudokuBoard {
// 	newBoard := NewSudokuBoard()

// 	for y := 0; y < BOARD_HEIGHT; y++ {
// 		for x := 0; x < BOARD_WIDTH; x++ {
// 			newBoard.Board[y][x] = b.Board[y][x]
// 		}
// 	}

// 	return newBoard
// }

func (b *SudokuBoard) GetPossibleBitmapInRow(column int) uint16 {
	not_possible := uint16(0)

	for x := 0; x < BOARD_WIDTH; x++ {
		if b.Board[column][x] != 0 {
			not_possible |= 1 << (b.Board[column][x] - 1)
		}
	}

	// Convert not_possible to possible numbers
	return ^not_possible & NUMBERS_BITMAP_MASK
}

func (b *SudokuBoard) GetPossibleBitmapInColumn(row int) uint16 {
	not_possible := uint16(0)

	for y := 0; y < BOARD_HEIGHT; y++ {
		if b.Board[y][row] != 0 {
			not_possible |= 1 << (b.Board[y][row] - 1)
		}
	}

	return ^not_possible & NUMBERS_BITMAP_MASK
}

func (b *SudokuBoard) GetPossibleBitmapInBlock(row int, column int) uint16 {
	not_possible := uint16(0)

	// blockRow and blockColumn are the top left position of the block
	blockRow := row - row%3
	blockColumn := column - column%3

	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			if b.Board[y][x] != 0 {
				not_possible |= 1 << (b.Board[blockColumn+y][blockRow+x] - 1)
			}
		}
	}

	return ^not_possible & NUMBERS_BITMAP_MASK
}

func (b *SudokuBoard) GetPossibleBitmapForCell(row int, column int) uint16 {
	row_possible := b.GetPossibleBitmapInRow(column)
	col_possible := b.GetPossibleBitmapInColumn(row)
	block_possible := b.GetPossibleBitmapInBlock(row, column)

	return row_possible & col_possible & block_possible
}

// Creates a random solvable board.
func (b *SudokuBoard) PlaceNumbers(amountOfGivens int) {
	tries := 0

	start_filling:
	for y := 0; y < BOARD_HEIGHT; y++ {
		for x := 0; x < BOARD_WIDTH; x++ {
			possible := b.GetPossibleBitmapForCell(x, y)

			// If no possible numbers can be placed, restart filling the row
			if possible == 0 {
				tries++
				if tries > 10 {
					// Refill the entire board
					b.Board = NewSudokuBoard().Board
					goto start_filling
				}

				b.Board[y] = [BOARD_WIDTH]byte{}
				x = -1

				continue
			}

			b.Board[y][x] = pickRandomNumberFromBitmap(possible)
		}

		tries = 0
	}

	// Remove numbers until only amountOfGivens is left
	const numberOfCells = BOARD_WIDTH * BOARD_HEIGHT
	amountToRemove := numberOfCells - amountOfGivens

	for amountToRemove > 0 {
		x := rand.Intn(BOARD_WIDTH)
		y := rand.Intn(BOARD_HEIGHT)

		if b.Board[y][x] == 0 {
			continue
		}

		b.Board[y][x] = 0

		amountToRemove--
	}
}

func (b *SudokuBoard) PrintAllSolutions() {
	for y := 0; y < BOARD_HEIGHT; y++ {
		for x := 0; x < BOARD_WIDTH; x++ {

			if b.Board[y][x] != 0 {
				continue
			}

			possible := b.GetPossibleBitmapForCell(x, y)

			if possible == 0 {
				return
			}

			possibleNumbers := NumbersFromBitmap(possible)

			for _, n := range possibleNumbers {
				b.Board[y][x] = n
				b.PrintAllSolutions()
			}

			b.Board[y][x] = 0
			return
		}
	}

	b.Print()
}

func (b *SudokuBoard) Print() {
	fmt.Println()

	for y := 0; y < BOARD_HEIGHT; y++ {
		if y != 0 && y % 3 == 0 {
			fmt.Println(" " + strings.Repeat("-", BOARD_WIDTH*3 + BOARD_WIDTH%3 + 1))
		}

		for x := 0; x < BOARD_WIDTH; x++ {
			if x % 3 == 0 {
				fmt.Print("|")
			}

			fmt.Printf(" %d ", b.Board[y][x])
		}
		fmt.Println("|")
	}

	fmt.Println()
}

// func randomInRange(min int, max int) int {
// 	return rand.Intn(max - min + 1) + min
// }

// Convert bitmap numbers to slice
func NumbersFromBitmap(bitmap uint16) []byte {
	bitmapNumbers := make([]byte, 0, 9)

	for i := 1; i <= 9; i++ {
		if bitmap & (1 << (i-1)) > 0 {
			bitmapNumbers = append(bitmapNumbers, byte(i))
		}
	}

	return bitmapNumbers
}

func pickRandomNumberFromBitmap(bitmap uint16) byte {
	bitmapNumbers := NumbersFromBitmap(bitmap)
	number := bitmapNumbers[rand.Intn(len(bitmapNumbers))]
	return number
}

func main() {
	rand.Seed(time.Now().UnixNano())

	board := NewSudokuBoard()

	board.PlaceNumbers(45)
	board.Print()

	board.PrintAllSolutions()
}
