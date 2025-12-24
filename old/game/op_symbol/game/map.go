package game

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type TerrainType int

const (
	Water TerrainType = iota
	Sand
	Plain
	Forest
	Mountain
)

// 地形颜色
var terrainColors = map[TerrainType]color.RGBA{
	Water:    {0, 100, 255, 255},   // 蓝色
	Sand:     {240, 230, 140, 255}, // 沙色
	Plain:    {120, 200, 80, 255},  // 绿色
	Forest:   {34, 139, 34, 255},   // 深绿色
	Mountain: {139, 137, 137, 255}, // 灰色
}

type GameMap struct {
	Width, Height int
	Terrain       [][]TerrainType
	NoiseMap      [][]float64
}

// 柏林噪声生成器
type PerlinNoise struct {
	permutation []int
	octaves     int
	persistence float64
}

func NewPerlinNoise(seed int64) *PerlinNoise {
	r := rand.New(rand.NewSource(seed))
	p := &PerlinNoise{
		permutation: make([]int, 512),
		octaves:     6,
		persistence: 0.5,
	}

	// 生成随机排列
	perm := make([]int, 256)
	for i := 0; i < 256; i++ {
		perm[i] = i
	}

	// 打乱排列
	for i := 255; i > 0; i-- {
		j := r.Intn(i + 1)
		perm[i], perm[j] = perm[j], perm[i]
	}

	// 复制到更大的数组中以简化边界处理
	for i := 0; i < 256; i++ {
		p.permutation[i] = perm[i]
		p.permutation[i+256] = perm[i]
	}

	return p
}

// 线性插值
func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

// 平滑插值
func smoothstep(t float64) float64 {
	return t * t * (3 - 2*t)
}

// 梯度函数
func (p *PerlinNoise) grad(hash int, x, y float64) float64 {
	h := hash & 15
	u := float64(0)
	v := float64(0)

	// 根据哈希值选择梯度向量
	switch h {
	case 0, 12:
		u, v = x, y
	case 1, 13:
		u, v = -x, y
	case 2, 14:
		u, v = x, -y
	case 3, 15:
		u, v = -x, -y
	case 4:
		u, v = x, 0
	case 5:
		u, v = -x, 0
	case 6:
		u, v = 0, y
	case 7:
		u, v = 0, -y
	case 8:
		u, v = x, x
	case 9:
		u, v = -x, y
	case 10:
		u, v = x, -y
	case 11:
		u, v = -x, -y
	}

	return u + v
}

// 噪声函数
func (p *PerlinNoise) noise(x, y float64) float64 {
	// 整数部分
	X := int(math.Floor(x)) & 255
	Y := int(math.Floor(y)) & 255

	// 小数部分
	x -= math.Floor(x)
	y -= math.Floor(y)

	// 计算淡入淡出曲线
	u := smoothstep(x)
	v := smoothstep(y)

	// 获取哈希值
	A := p.permutation[X] + Y
	B := p.permutation[X+1] + Y

	// 计算梯度并插值
	return lerp(
		lerp(p.grad(p.permutation[A], x, y),
			p.grad(p.permutation[B], x-1, y), u),
		lerp(p.grad(p.permutation[A+1], x, y-1),
			p.grad(p.permutation[B+1], x-1, y-1), u),
		v)
}

// 分形柏林噪声
func (p *PerlinNoise) fractalNoise(x, y float64) float64 {
	total := 0.0
	frequency := 1.0
	amplitude := 1.0
	maxValue := 0.0

	for i := 0; i < p.octaves; i++ {
		total += p.noise(x*frequency, y*frequency) * amplitude
		maxValue += amplitude
		amplitude *= p.persistence
		frequency *= 2
	}

	return total / maxValue
}

func NewGameMap(width, height int) *GameMap {
	m := &GameMap{
		Width:    width,
		Height:   height,
		Terrain:  make([][]TerrainType, height),
		NoiseMap: make([][]float64, height),
	}

	// 初始化柏林噪声生成器
	perlin := NewPerlinNoise(rand.Int63())

	// 生成噪声地图
	for y := 0; y < height; y++ {
		m.Terrain[y] = make([]TerrainType, width)
		m.NoiseMap[y] = make([]float64, width)

		for x := 0; x < width; x++ {
			// 生成噪声值，范围在 -1 到 1 之间
			nx := float64(x) / float64(width) * 4 // 调整缩放以获得更好的视觉效果
			ny := float64(y) / float64(height) * 4

			// 生成分形噪声
			noiseValue := perlin.fractalNoise(nx, ny)
			m.NoiseMap[y][x] = noiseValue

			// 根据噪声值确定地形类型
			if noiseValue < -0.25 {
				m.Terrain[y][x] = Water
			} else if noiseValue < 0.0 {
				m.Terrain[y][x] = Sand
			} else if noiseValue < 0.25 {
				m.Terrain[y][x] = Plain
			} else if noiseValue < 0.5 {
				m.Terrain[y][x] = Forest
			} else {
				m.Terrain[y][x] = Mountain
			}
		}
	}

	return m
}

func (m *GameMap) Draw(screen *ebiten.Image) {
	// 绘制地形
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			terrainType := m.Terrain[y][x]
			terrainColor := terrainColors[terrainType]

			// 绘制地形颜色
			ebitenutil.DrawRect(screen,
				float64(x*TileSize),
				float64(y*TileSize),
				TileSize, TileSize,
				terrainColor)
		}
	}
}
