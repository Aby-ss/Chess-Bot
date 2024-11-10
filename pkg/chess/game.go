package main

import (
	"image/color"
	"log"
	"path/filepath"
	"strings"

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
	board       = parseFEN(FENPosition)

	selectedPiece string
	selectedRow   int
	selectedCol   int
	isDragging    bool
)

func loadPieceImage(filename string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(filepath.Join("assets", filename))
	if err != nil {
		log.Printf("Failed to load image %s: %v", filename, err)
		return nil
	}
	return img
}

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

type Game struct{}

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
		}
	}

	for row, rowValues := range board {
		for col, piece := range rowValues {
			if piece != "" && !(isDragging && row == selectedRow && col == selectedCol) {
				drawPiece(screen, piece, row, col)
			}
		}
	}

	if isDragging && selectedPiece != "" {
		x, y := ebiten.CursorPosition()
		drawPieceAt(screen, selectedPiece, x, y)
	}
}

func drawPiece(screen *ebiten.Image, pieceName string, row, col int) {
	pieceImg := pieceImages[pieceName]
	if pieceImg == nil {
		log.Printf("Image for piece %s not found ❌", pieceName)
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

func drawPieceAt(screen *ebiten.Image, pieceName string, x, y int) {
	pieceImg := pieceImages[pieceName]
	if pieceImg == nil {
		return
	}

	scale := float64(SquareSize) / float64(pieceImg.Bounds().Dx())
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(scale, scale)
	options.GeoM.Translate(float64(x)-float64(SquareSize)/2, float64(y)-float64(SquareSize)/2)
	screen.DrawImage(pieceImg, options)
}

func (g *Game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		col, row := x/SquareSize, y/SquareSize

		if !isDragging && board[row][col] != "" {
			selectedPiece = board[row][col]
			selectedRow, selectedCol = row, col
			board[row][col] = ""
			isDragging = true
		}
	} else if isDragging {
		x, y := ebiten.CursorPosition()
		col, row := x/SquareSize, y/SquareSize

		board[row][col] = selectedPiece
		isDragging = false
		selectedPiece = ""
	}
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	initPieces()

	game := &Game{}
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Chess Beta Version 0.1.2")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
