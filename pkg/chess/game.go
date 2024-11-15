package main

import (
	"image/color"
	"log"
	"path/filepath"
	"strings"

	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth  = 540
	ScreenHeight = 540
	BoardSize    = 8
	SquareSize   = ScreenWidth / BoardSize
)

var (
	orange      = color.RGBA{255, 165, 0, 255}
	lightSkin   = color.RGBA{255, 228, 196, 255}
	pieceImages = make(map[string]*ebiten.Image)
	FENPosition = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"
)

// Load image for a piece from assets
func loadPieceImage(filename string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(filepath.Join("assets", filename))
	if err != nil {
		log.Printf("Failed to load image %s: %v", filename, err)
		return nil
	}
	return img
}

// Initialize images for pieces
func initPieces() {
	pieceImages["K"] = loadPieceImage("white-king.png")
	pieceImages["Q"] = loadPieceImage("white-queen.png")
	pieceImages["R"] = loadPieceImage("white-rook.png")
	pieceImages["N"] = loadPieceImage("white-knight.png")
	pieceImages["B"] = loadPieceImage("white-bishop.png")
	pieceImages["P"] = loadPieceImage("white-pawn.png")

	pieceImages["k"] = loadPieceImage("black-king.png")
	pieceImages["q"] = loadPieceImage("black-queen.png")
	pieceImages["r"] = loadPieceImage("black-rook.png")
	pieceImages["n"] = loadPieceImage("black-knight.png")
	pieceImages["b"] = loadPieceImage("black-bishop.png")
	pieceImages["p"] = loadPieceImage("black-pawn.png")
}

// Parse FEN notation into a 2D board array
func parseFEN(fen string) [][]string {
	board := make([][]string, BoardSize)
	rows := strings.Split(fen, "/")

	for i, row := range rows {
		board[i] = make([]string, BoardSize)
		col := 0
		for _, char := range row {
			if char >= '1' && char <= '8' {
				col += int(char - '0')
			} else {
				board[i][col] = string(char)
				col++
			}
		}
	}
	return board
}

// Game structure for Ebiten with piece selection tracking
type Game struct {
	board         [][]string
	selectedPiece string
	selectedRow   int
	selectedCol   int
	isDragging    bool
	currentTurn   string
}

// Draw the chessboard and pieces
func (g *Game) Draw(screen *ebiten.Image) {
	for row := 0; row < BoardSize; row++ {
		for col := 0; col < BoardSize; col++ {
			var squareColor color.Color
			if (row+col)%2 == 0 {
				squareColor = orange
			} else {
				squareColor = lightSkin
			}
			x := float64(col * SquareSize)
			y := float64(row * SquareSize)
			ebitenutil.DrawRect(screen, x, y, SquareSize, SquareSize, squareColor)

			if piece := g.board[row][col]; piece != "" && !(g.isDragging && row == g.selectedRow && col == g.selectedCol) {
				drawPiece(screen, piece, row, col)
			}
		}
	}

	// Draw the selected piece at mouse position if dragging
	if g.isDragging && g.selectedPiece != "" {
		x, y := ebiten.CursorPosition()
		drawPieceAtPosition(screen, g.selectedPiece, float64(x), float64(y))
	}
}

// Draw a piece at a specific square
func drawPiece(screen *ebiten.Image, pieceName string, row, col int) {
	pieceImg := pieceImages[pieceName]
	if pieceImg == nil {
		log.Printf("Image for piece %s not found", pieceName)
		return
	}

	scale := float64(SquareSize) / float64(pieceImg.Bounds().Dx())
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(scale, scale)

	x := float64(col*SquareSize) + (SquareSize-float64(pieceImg.Bounds().Dx())*scale)/2
	y := float64(row*SquareSize) + (SquareSize-float64(pieceImg.Bounds().Dy())*scale)/2
	options.GeoM.Translate(x, y)

	screen.DrawImage(pieceImg, options)
}

// Helper function to draw a piece at an arbitrary position (for dragging)
func drawPieceAtPosition(screen *ebiten.Image, pieceName string, x, y float64) {
	pieceImg := pieceImages[pieceName]
	scale := float64(SquareSize) / float64(pieceImg.Bounds().Dx())
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(scale, scale)
	options.GeoM.Translate(x-float64(SquareSize)/2, y-float64(SquareSize)/2)
	screen.DrawImage(pieceImg, options)
}

func (g *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		col, row := x/SquareSize, y/SquareSize

		if !g.isDragging {
			// Select a piece only if it matches the current turn
			selectedPiece := g.board[row][col]
			if selectedPiece != "" && g.isPieceTurn(selectedPiece) {
				g.selectedPiece = selectedPiece
				g.selectedRow = row
				g.selectedCol = col
				g.isDragging = true
			}
		}
	} else if g.isDragging {
		x, y := ebiten.CursorPosition()
		col, row := x/SquareSize, y/SquareSize

		// Try to complete the move
		if g.isValidMove(g.selectedRow, g.selectedCol, row, col) {
			g.board[row][col] = g.selectedPiece
			g.board[g.selectedRow][g.selectedCol] = ""

			// Switch the turn after a valid move
			g.switchTurn()
		}
		// Reset selection state
		g.isDragging = false
		g.selectedPiece = ""
	}
	return nil
}

// Check if the piece belongs to the current turn
func (g *Game) isPieceTurn(piece string) bool {
	if g.currentTurn == "white" {
		return piece >= "A" && piece <= "Z" // Uppercase pieces for white
	}
	return piece >= "a" && piece <= "z" // Lowercase pieces for black
}

// Switch the turn to the other player
func (g *Game) switchTurn() {
	if g.currentTurn == "white" {
		g.currentTurn = "black"
	} else {
		g.currentTurn = "white"
	}
}

// Move validation with basic piece movement rules
func (g *Game) isValidMove(fromRow, fromCol, toRow, toCol int) bool {
	piece := g.board[fromRow][fromCol]

	rowDiff := math.Abs(float64(toRow - fromRow))
	colDiff := math.Abs(float64(toCol - fromCol))

	switch piece {
	case "P":
		if fromRow == 6 && toRow == 4 && colDiff == 0 { // Two-step pawn move
			return true
		}
		return toRow == fromRow-1 && colDiff == 0 // One-step move forward
	case "p":
		if fromRow == 1 && toRow == 3 && colDiff == 0 {
			return true
		}
		return toRow == fromRow+1 && colDiff == 0
	case "R", "r":
		return rowDiff == 0 || colDiff == 0 // Rook moves
	case "B", "b":
		return rowDiff == colDiff // Bishop moves diagonally
	case "Q", "q":
		return rowDiff == colDiff || rowDiff == 0 || colDiff == 0 // Queen moves
	case "N", "n":
		return (rowDiff == 2 && colDiff == 1) || (rowDiff == 1 && colDiff == 2) // Knight moves
	case "K", "k":
		return rowDiff <= 1 && colDiff <= 1 // King moves one square
	}
	return false
}

// Layout function to set screen dimensions
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	initPieces()
	game := &Game{
		board:       parseFEN(FENPosition),
		currentTurn: "white", // White moves first
	}

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Chessboard with Turn-Based Logic")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
