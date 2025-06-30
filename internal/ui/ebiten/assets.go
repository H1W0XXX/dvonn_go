// File internal/ui/ebiten/assets.go
package ebiten

import (
	"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	_ "image/png"
)

//go:embed assets/red1.png
var redPNG []byte

//go:embed assets/white1.png
var whitePNG []byte

//go:embed assets/black1.png
var blackPNG []byte

var (
	imgRed, imgWhite, imgBlack *ebiten.Image
)

func init() {
	imgRed = decode(redPNG)
	imgWhite = decode(whitePNG)
	imgBlack = decode(blackPNG)
}
func decode(b []byte) *ebiten.Image {
	//fmt.Println("Decoding image...")
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		panic("Failed to decode image: " + err.Error())
	}
	return ebiten.NewImageFromImage(img)
}
