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

// Update function to handle dragging and dropping of pieces
func (g *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		col, row := x/SquareSize, y/SquareSize

		if !g.isDragging {
			g.selectedPiece = g.board[row][col]
			if g.selectedPiece != "" {
				g.selectedRow, g.selectedCol = row, col
				g.isDragging = true
			}
		}
	} else if g.isDragging {
		x, y := ebiten.CursorPosition()
		col, row := x/SquareSize, y/SquareSize

		// Complete the move if mouse is released on a valid square
		if g.isValidMove(g.selectedRow, g.selectedCol, row, col) {
			g.board[row][col] = g.selectedPiece
			g.board[g.selectedRow][g.selectedCol] = ""
		}
		g.isDragging = false
		g.selectedPiece = ""
	}
	return nil
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
	game := &Game{board: parseFEN(FENPosition)}

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Chessboard with Drag-and-Drop Pieces")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
