package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "image/color"
    "log"
    "path/filepath"
    "strings"
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

// Load a single image file from assets
func loadPieceImage(filename string) *ebiten.Image {
    img, _, err := ebitenutil.NewImageFromFile(filepath.Join("assets", filename))
    if err != nil {
        log.Printf("Failed to load image %s: %v", filename, err)
        return nil
    }
    return img
}

// Initialize all piece images
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
                // Empty squares
                col += int(char - '0')
            } else {
                // Chess piece
                board[i][col] = string(char)
                col++
            }
        }
    }
    return board
}

// Game structure for Ebiten
type Game struct{}

// Draw the chessboard and pieces
func (g *Game) Draw(screen *ebiten.Image) {
    // Draw the chessboard squares
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

    // Draw pieces from the FEN position
    board := parseFEN(FENPosition)
    for row, rowValues := range board {
        for col, piece := range rowValues {
            if piece != "" {
                drawPiece(screen, piece, row, col)
            }
        }
    }
}

// Function to draw a piece at the center of a specific square
func drawPiece(screen *ebiten.Image, pieceName string, row, col int) {
    pieceImg := pieceImages[pieceName]
    if pieceImg == nil {
        log.Printf("Image for piece %s not found", pieceName)
        return
    }

    // Scale the image to fit within the square size
    scale := float64(SquareSize) / float64(pieceImg.Bounds().Dx())
    options := &ebiten.DrawImageOptions{}
    options.GeoM.Scale(scale, scale)

    // Calculate centered position within the square
    x := float64(col*SquareSize) + (SquareSize-float64(pieceImg.Bounds().Dx())*scale)/2
    y := float64(row*SquareSize) + (SquareSize-float64(pieceImg.Bounds().Dy())*scale)/2
    options.GeoM.Translate(x, y)

    screen.DrawImage(pieceImg, options)
}

// Update function for Ebiten (unused here)
func (g *Game) Update() error {
    return nil
}

// Layout function to set screen dimensions
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return ScreenWidth, ScreenHeight
}

func main() {
    initPieces() // Load all pieces

    game := &Game{}
    ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
    ebiten.SetWindowTitle("Chessboard with Pieces")

    // Start the game loop
    if err := ebiten.RunGame(game); err != nil {
        log.Fatal(err)
    }
}
