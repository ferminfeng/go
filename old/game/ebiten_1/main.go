package main

import (
	"fmt"
	"image/color"
	"io"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

const (
	screenWidth  = 480
	screenHeight = 480
	targetWidth  = 48
	targetHeight = 48
	sampleRate   = 44100
	incScore     = 1 // 每次增加分数
	maxScore     = 0 // 通关分数，0表示无通关分数
)

var (
	// 中文字体
	chineseFont font.Face
)

type Game struct {
	playerX     float64
	playerY     float64
	playerImage *ebiten.Image
	enemies     []Enemy
	items       []Item
	score       int
	gameOver    bool
	gameStart   bool
	gameWin     bool
	background  *ebiten.Image
	audioPlayer *audio.Player
	goldSound   *audio.Player
	winSound    *audio.Player
}

type Enemy struct {
	x, y  float64
	image *ebiten.Image
	speed float64
}

type Item struct {
	x, y  float64
	image *ebiten.Image
}

// loadImage 加载图片资源
func loadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatalf("Failed to load image %s: %v", path, err)
	}
	return img
}

// loadSound 加载音效资源
func loadSound(audioContext *audio.Context, path string) *audio.Player {
	file, err := ebitenutil.OpenFile(path)
	if err != nil {
		log.Fatalf("Failed to load sound %s: %v", path, err)
	}
	stream, err := mp3.Decode(audioContext, file)
	if err != nil {
		log.Fatalf("Failed to decode sound %s: %v", path, err)
	}
	player, err := audioContext.NewPlayer(stream)
	if err != nil {
		log.Fatalf("Failed to create sound player for %s: %v", path, err)
	}
	return player
}

// loadFont 加载中文字体
func loadFont() font.Face {
	fontFile, err := ebitenutil.OpenFile("assets/font.ttf") // 确保 assets/font.ttf 是中文字体文件
	if err != nil {
		log.Fatalf("Failed to load font: %v", err)
	}
	defer fontFile.Close()

	fontData, err := io.ReadAll(fontFile) // Read the file into a []byte
	if err != nil {
		log.Fatalf("Failed to read font data: %v", err)
	}

	tt, err := truetype.Parse(fontData)
	if err != nil {
		log.Fatalf("Failed to parse font: %v", err)
	}
	return truetype.NewFace(tt, &truetype.Options{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

// NewGame 初始化游戏
func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())
	audioContext := audio.NewContext(sampleRate)

	// 加载中文字体
	chineseFont = loadFont()

	return &Game{
		playerX:     screenWidth/2 - targetWidth/2,
		playerY:     screenHeight/2 - targetHeight/2,
		playerImage: loadImage("assets/player.png"),
		enemies:     makeEnemies(1),
		items:       makeItems(5),
		background:  loadImage("assets/background.png"),
		audioPlayer: loadSound(audioContext, "assets/background.mp3"),
		// goldSound:   loadSound(audioContext, "assets/gold.mp3"),
		// winSound: loadSound(audioContext, "assets/win.mp3"),
	}
}

// makeEnemies 初始化敌人
func makeEnemies(count int) []Enemy {
	enemies := make([]Enemy, count)
	for i := range enemies {
		enemies[i] = Enemy{
			x:     rand.Float64() * screenWidth,
			y:     rand.Float64() * screenHeight,
			image: loadImage("assets/enemy.png"),
			speed: 0.5 + rand.Float64()*0.5,
		}
	}
	return enemies
}

// makeItems 初始化道具
func makeItems(count int) []Item {
	items := make([]Item, count)
	for i := range items {
		items[i] = Item{
			x:     rand.Float64() * screenWidth,
			y:     rand.Float64() * screenHeight,
			image: loadImage("assets/gold.png"),
		}
	}
	return items
}

// checkCollision 检测两个矩形是否碰撞
func checkCollision(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 &&
		x1+w1 > x2 &&
		y1 < y2+h2 &&
		y1+h1 > y2
}

// handleInput 处理玩家输入
func (g *Game) handleInput() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.playerX -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.playerX += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.playerY -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.playerY += 2
	}
}

// updateEnemies 更新敌人位置和状态
func (g *Game) updateEnemies() {
	for i := range g.enemies {
		dx := g.playerX - g.enemies[i].x
		dy := g.playerY - g.enemies[i].y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist > 0 {
			g.enemies[i].x += dx / dist * g.enemies[i].speed
			g.enemies[i].y += dy / dist * g.enemies[i].speed
		}

		if checkCollision(g.playerX, g.playerY, targetWidth, targetHeight, g.enemies[i].x, g.enemies[i].y, targetWidth, targetHeight) {
			g.gameOver = true
		}
	}
}

// updateItems 更新道具位置和状态
func (g *Game) updateItems() {
	for i := range g.items {
		if checkCollision(g.playerX, g.playerY, targetWidth, targetHeight, g.items[i].x, g.items[i].y, targetWidth, targetHeight) {
			g.score += incScore
			g.items[i].x = rand.Float64() * screenWidth
			g.items[i].y = rand.Float64() * screenHeight
			// g.goldSound.Rewind()
			// g.goldSound.Play()
		}
	}
}

// Update 更新游戏逻辑
func (g *Game) Update() error {
	if !g.gameStart {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.gameStart = true
			g.audioPlayer.Play()
		}
		return nil
	}

	if g.gameOver || g.gameWin {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.reset()
		}
		return nil
	}

	g.handleInput()
	g.updateEnemies()
	g.updateItems()

	if maxScore > 0 && g.score >= maxScore {
		g.gameWin = true
		// g.winSound.Play()
	}

	return nil
}

// Draw 绘制游戏画面
func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.background, nil)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(targetWidth/float64(g.playerImage.Bounds().Dx()), targetHeight/float64(g.playerImage.Bounds().Dy()))
	op.GeoM.Translate(g.playerX, g.playerY)
	screen.DrawImage(g.playerImage, op)

	for _, enemy := range g.enemies {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(targetWidth/float64(enemy.image.Bounds().Dx()), targetHeight/float64(enemy.image.Bounds().Dy()))
		op.GeoM.Translate(enemy.x, enemy.y)
		screen.DrawImage(enemy.image, op)
	}

	for _, item := range g.items {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(targetWidth/float64(item.image.Bounds().Dx()), targetHeight/float64(item.image.Bounds().Dy()))
		op.GeoM.Translate(item.x, item.y)
		screen.DrawImage(item.image, op)
	}

	// 显示得分
	text.Draw(screen, fmt.Sprintf("得分: %d", g.score), chineseFont, 10, 30, color.White)

	// 显示游戏开始信息
	if !g.gameStart {
		text.Draw(screen, "按空格键开始", chineseFont, screenWidth/2-100, screenHeight/2, color.White)
	}

	// 显示游戏结束信息
	if g.gameOver {
		text.Draw(screen, fmt.Sprintf("游戏结束! 最终得分: %d. 按空格键重新开始", g.score), chineseFont, screenWidth/2-150, screenHeight/2, color.White)
	}

	// 显示游戏胜利信息
	if g.gameWin {
		text.Draw(screen, fmt.Sprintf("你赢了! 最终得分: %d. 按空格键重新开始", g.score), chineseFont, screenWidth/2-150, screenHeight/2, color.White)
	}
}

// Layout 设置游戏窗口大小
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// reset 重置游戏状态
func (g *Game) reset() {
	g.playerX = screenWidth/2 - targetWidth/2
	g.playerY = screenHeight/2 - targetHeight/2
	g.score = 0
	g.gameOver = false
	g.gameStart = false
	g.gameWin = false
	g.audioPlayer.Rewind()
}

// main 函数，游戏入口
func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("我的 Ebiten 游戏")

	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
