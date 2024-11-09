package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "image/color"
    "log"
    "path/filepath"
)

const (
    ScreenWidth  = 540
    ScreenHeight = 540
    BoardSize    = 8
    SquareSize   = ScreenWidth / BoardSize
)

var (
    orange    = color.RGBA{255, 165, 0, 255}
    lightSkin = color.RGBA{255, 228, 196, 255}
    pieceImages = make(map[string]*ebiten.Image)
)

// Load a single image file from assets
func loadPieceImage(filename string) *ebiten.Image {
    img, _, err := ebitenutil.NewImageFromFile(filepath.Join("assets", filename))
    if err != nil {
        log.Printf("Failed to load image %s: %v", filename, err)
        return nil
    }
    log.Printf("Loaded image: %s", filename)
    return img
}

// Initialize all piece images
func initPieces() {
    pieceImages["white_king"] = loadPieceImage("white-king.png")
    pieceImages["white_queen"] = loadPieceImage("white-queen.png")
    pieceImages["white_rook"] = loadPieceImage("white-rook.png")
    pieceImages["white_knight"] = loadPieceImage("white-knight.png")
    pieceImages["white_bishop"] = loadPieceImage("white-bishop.png")
    pieceImages["white_pawn"] = loadPieceImage("white-pawn.png")

    pieceImages["black_king"] = loadPieceImage("black-king.png")
    pieceImages["black_queen"] = loadPieceImage("black-queen.png")
    pieceImages["black_rook"] = loadPieceImage("black-rook.png")
    pieceImages["black_knight"] = loadPieceImage("black-knight.png")
    pieceImages["black_bishop"] = loadPieceImage("black-bishop.png")
    pieceImages["black_pawn"] = loadPieceImage("black-pawn.png")

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

    // Example of drawing pieces at specific squares
    drawPiece(screen, "white_king", 0, 4)  // Place white king on row 0, col 4
    drawPiece(screen, "black_queen", 7, 3) // Place black queen on row 7, col 3
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

    log.Printf("Drawing %s at row %d, col %d", pieceName, row, col)
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
