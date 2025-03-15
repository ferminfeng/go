package game

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

// 创建一个全局的空白纹理，用于绘制多边形
var emptyImage *ebiten.Image

// 全局实体ID计数器，用于生成唯一的实体ID
var nextEntityID int = 1

func init() {
	// 初始化一个 1x1 的空白纹理
	emptyImage = ebiten.NewImage(1, 1)
	emptyImage.Fill(color.White)

	// 初始化实体ID计数器
	nextEntityID = 1
}

type UnitType int

const (
	Infantry       UnitType = iota // 步兵
	Armor                          // 装甲车
	Artillery                      // 火炮
	Recon                          // 侦察车
	AntiAir                        // 防空炮
	HeavyTank                      // 重型坦克
	Helicopter                     // 直升机
	FighterJet                     // 战斗机
	Bomber                         // 轰炸机
	RocketLauncher                 // 火箭炮
	Engineer                       // 工程兵
	MedicUnit                      // 医疗单位
)

type Unit struct {
	X, Y         int
	Type         UnitType
	Health       int
	IsPlayerUnit bool
	EntityID     int // 唯一实体ID，用于排序

	// 添加实时战略所需的属性
	TargetX, TargetY   int      // 目标位置
	HasTarget          bool     // 是否有目标
	Path               [][2]int // 寻路路径
	MoveSpeed          float64  // 移动速度
	CurrentX, CurrentY float64  // 实际位置（浮点数，用于平滑移动）
	Selected           bool     // 是否被选中
	TargetUnit         *Unit    // 目标单位，用于持续追踪

	// 装甲单位搭载步兵相关属性
	Passengers    []*Unit // 搭载的步兵单位
	MaxPassengers int     // 最大搭载数量
	IsPassenger   bool    // 是否是乘客
	ParentUnit    *Unit   // 所属的载具单位

	// 新增属性
	AttackPower     int       // 攻击力
	AttackRange     int       // 攻击范围
	AttackFrequency float64   // 攻击频率（每秒攻击次数）
	LastAttackTime  time.Time // 上次攻击时间
	Defense         int       // 防御力
	Size            float64   // 单位大小（相对于TileSize的比例）
	MaxHitRate      float64   // 最大命中率（0.0-1.0）
	VisionRange     int       // 视野范围（格子数）
	CanAttackAir    bool      // 是否能够攻击空中单位
	MaxAmmo         int       // 最大弹药量
	CurrentAmmo     int       // 当前弹药量
	Name            string    // 单位名称

	// 攻击动画相关
	IsAttacking        bool      // 是否正在攻击
	AttackAnimTime     time.Time // 攻击动画开始时间
	AttackAnimDuration float64   // 攻击动画持续时间（秒）
}

func NewUnit(x, y int, unitType UnitType, isPlayerUnit bool) *Unit {
	// 获取唯一的实体ID并递增计数器
	entityID := nextEntityID
	nextEntityID++

	// 从配置中获取单位类型数据
	typeData, err := GetUnitTypeData(unitType)
	if err != nil {
		// 如果配置不存在，使用默认值
		health := 3
		moveSpeed := 1.0
		maxPassengers := 0
		attackPower := 1
		attackRange := 2 // 默认攻击范围增加到2
		attackFrequency := 1.0
		defense := 0
		size := 1.0
		maxHitRate := 0.8
		visionRange := 4      // 默认视野范围
		canAttackAir := false // 默认不能攻击空中单位
		maxAmmo := 10         // 默认弹药量
		name := "未知单位"

		if unitType == Armor {
			health = 5
			moveSpeed = 0.8
			maxPassengers = 2
			attackPower = 2
			attackRange = 2 // 装甲车攻击范围
			attackFrequency = 0.8
			defense = 2
			size = 1.2
			maxHitRate = 0.8
			visionRange = 4      // 装甲车视野范围
			canAttackAir = false // 装甲车不能攻击空中单位
			maxAmmo = 15         // 装甲车弹药量
			name = "装甲车"
		} else if unitType == Infantry {
			attackFrequency = 1.2
			attackRange = 2 // 步兵攻击范围
			maxHitRate = 0.8
			visionRange = 3      // 步兵视野范围
			canAttackAir = false // 步兵不能攻击空中单位
			maxAmmo = 10         // 步兵弹药量
			name = "步兵"
			size = 0.8
		} else if unitType == Artillery {
			moveSpeed = 0.5
			attackPower = 3
			attackRange = 5 // 火炮攻击范围
			attackFrequency = 0.5
			maxHitRate = 0.5
			visionRange = 4      // 火炮视野范围
			canAttackAir = false // 火炮不能攻击空中单位
			maxAmmo = 8          // 火炮弹药量
			name = "火炮"
		} else if unitType == Recon {
			health = 2
			moveSpeed = 1.5
			attackRange = 2 // 侦察车攻击范围
			attackFrequency = 1.5
			size = 0.9
			maxHitRate = 0.8
			visionRange = 6      // 侦察车视野范围更大
			canAttackAir = false // 侦察车不能攻击空中单位
			maxAmmo = 12         // 侦察车弹药量
			name = "侦察车"
		} else if unitType == AntiAir {
			moveSpeed = 1.2
			attackPower = 2
			attackRange = 3 // 防空炮攻击范围
			attackFrequency = 1.2
			defense = 1
			maxHitRate = 0.8
			visionRange = 5     // 防空炮视野范围
			canAttackAir = true // 防空炮可以攻击空中单位
			maxAmmo = 20        // 防空炮弹药量
			name = "防空炮"
		} else if unitType == HeavyTank {
			health = 8
			moveSpeed = 0.6
			attackPower = 4
			attackRange = 2 // 重型坦克攻击范围
			attackFrequency = 0.6
			defense = 3
			size = 1.4
			maxHitRate = 0.6
			visionRange = 3      // 重型坦克视野范围
			canAttackAir = false // 重型坦克不能攻击空中单位
			maxAmmo = 10         // 重型坦克弹药量
			name = "重型坦克"
		} else if unitType == Helicopter {
			health = 4
			moveSpeed = 2.0
			attackPower = 2
			attackRange = 3 // 直升机攻击范围
			attackFrequency = 1.8
			defense = 1
			maxPassengers = 4
			size = 1.1
			maxHitRate = 0.8
			visionRange = 7     // 直升机视野范围更大
			canAttackAir = true // 直升机可以攻击空中单位
			maxAmmo = 16        // 直升机弹药量
			name = "直升机"
		} else if unitType == FighterJet {
			health = 3
			moveSpeed = 3.0
			attackPower = 3
			attackRange = 3 // 战斗机攻击范围
			attackFrequency = 2.0
			defense = 1
			size = 1.0
			maxHitRate = 0.8
			visionRange = 8     // 战斗机视野范围最大
			canAttackAir = true // 战斗机可以攻击空中单位
			maxAmmo = 6         // 战斗机弹药量
			name = "战斗机"
		} else if unitType == Bomber {
			health = 5
			moveSpeed = 1.8
			attackPower = 5
			attackRange = 3 // 轰炸机攻击范围
			attackFrequency = 0.7
			size = 1.3
			maxHitRate = 0.5
			visionRange = 6      // 轰炸机视野范围
			canAttackAir = false // 轰炸机不能攻击空中单位
			maxAmmo = 5          // 轰炸机弹药量
			name = "轰炸机"
		} else if unitType == RocketLauncher {
			health = 3
			moveSpeed = 0.7
			attackPower = 4
			attackRange = 6 // 火箭炮攻击范围
			attackFrequency = 0.4
			size = 1.1
			maxHitRate = 0.4
			visionRange = 5     // 火箭炮视野范围
			canAttackAir = true // 火箭炮可以攻击空中单位
			maxAmmo = 4         // 火箭炮弹药量
			name = "火箭炮"
		} else if unitType == Engineer {
			health = 2
			moveSpeed = 0.9
			attackPower = 1
			attackRange = 2 // 工程兵攻击范围
			attackFrequency = 1.0
			size = 0.8
			maxHitRate = 0.8
			visionRange = 3      // 工程兵视野范围
			canAttackAir = false // 工程兵不能攻击空中单位
			maxAmmo = 5          // 工程兵弹药量
			name = "工程兵"
		} else if unitType == MedicUnit {
			health = 3
			moveSpeed = 1.1
			attackPower = 0
			attackRange = 0 // 医疗单位没有攻击范围
			attackFrequency = 0.5
			size = 0.9
			maxHitRate = 0.5
			visionRange = 4      // 医疗单位视野范围
			canAttackAir = false // 医疗单位不能攻击空中单位
			maxAmmo = 0          // 医疗单位没有弹药
			name = "医疗单位"
		}

		unit := &Unit{
			X:               x,
			Y:               y,
			Type:            unitType,
			Health:          health,
			IsPlayerUnit:    isPlayerUnit,
			EntityID:        entityID, // 设置唯一实体ID
			HasTarget:       false,
			Path:            make([][2]int, 0),
			MoveSpeed:       moveSpeed,
			CurrentX:        float64(x*TileSize) + float64(TileSize)/2,
			CurrentY:        float64(y*TileSize) + float64(TileSize)/2,
			Selected:        false,
			Passengers:      make([]*Unit, 0),
			MaxPassengers:   maxPassengers,
			IsPassenger:     false,
			ParentUnit:      nil,
			AttackPower:     attackPower,
			AttackRange:     attackRange,
			AttackFrequency: attackFrequency,
			LastAttackTime:  time.Now().Add(-10 * time.Second), // 设置一个过去的时间，使单位可以立即攻击
			Defense:         defense,
			Size:            size,
			MaxHitRate:      maxHitRate,
			VisionRange:     visionRange,
			CanAttackAir:    canAttackAir,
			MaxAmmo:         maxAmmo,
			CurrentAmmo:     maxAmmo, // 初始弹药为最大值
			Name:            name,
			// 初始化攻击动画相关字段
			IsAttacking:        false,
			AttackAnimTime:     time.Time{},
			AttackAnimDuration: 0.2,
		}
		return unit
	}

	// 使用配置数据创建单位
	unit := &Unit{
		X:               x,
		Y:               y,
		Type:            unitType,
		Health:          typeData.Health,
		IsPlayerUnit:    isPlayerUnit,
		EntityID:        entityID, // 设置唯一实体ID
		HasTarget:       false,
		Path:            make([][2]int, 0),
		MoveSpeed:       typeData.MoveSpeed,
		CurrentX:        float64(x*TileSize) + float64(TileSize)/2,
		CurrentY:        float64(y*TileSize) + float64(TileSize)/2,
		Selected:        false,
		Passengers:      make([]*Unit, 0),
		MaxPassengers:   typeData.MaxPassengers,
		IsPassenger:     false,
		ParentUnit:      nil,
		AttackPower:     typeData.AttackPower,
		AttackRange:     typeData.AttackRange,
		AttackFrequency: typeData.AttackFrequency,
		LastAttackTime:  time.Now().Add(-10 * time.Second), // 设置一个过去的时间，使单位可以立即攻击
		Defense:         typeData.Defense,
		Size:            typeData.Size,
		MaxHitRate:      typeData.MaxHitRate,
		VisionRange:     typeData.VisionRange,
		CanAttackAir:    typeData.CanAttackAir,
		MaxAmmo:         typeData.Ammo,
		CurrentAmmo:     typeData.Ammo, // 初始弹药为最大值
		Name:            typeData.Name,
		// 初始化攻击动画相关字段
		IsAttacking:        false,
		AttackAnimTime:     time.Time{},
		AttackAnimDuration: 0.2,
	}
	return unit
}

// 添加步兵上车方法
func (u *Unit) AddPassenger(infantry *Unit) bool {
	// 检查是否是装甲单位
	if u.Type != Armor {
		return false
	}

	// 检查是否有空位
	if len(u.Passengers) >= u.MaxPassengers {
		return false
	}

	// 检查是否是步兵
	if infantry.Type != Infantry {
		return false
	}

	// 设置步兵为乘客状态
	infantry.IsPassenger = true
	infantry.ParentUnit = u

	// 添加到乘客列表
	u.Passengers = append(u.Passengers, infantry)
	return true
}

// 添加步兵下车方法
func (u *Unit) RemovePassenger(infantry *Unit) bool {
	// 检查是否是装甲单位
	if u.Type != Armor {
		return false
	}

	// 查找并移除乘客
	for i, passenger := range u.Passengers {
		if passenger == infantry {
			// 重置步兵状态
			infantry.IsPassenger = false
			infantry.ParentUnit = nil

			// 设置步兵位置为装甲单位附近的位置
			infantry.X = u.X
			infantry.Y = u.Y + 1 // 默认放在下方，后续会检查是否可行
			infantry.CurrentX = float64(infantry.X*TileSize) + float64(TileSize)/2
			infantry.CurrentY = float64(infantry.Y*TileSize) + float64(TileSize)/2

			// 从乘客列表中移除
			u.Passengers = append(u.Passengers[:i], u.Passengers[i+1:]...)
			return true
		}
	}

	return false
}

// 更新单位位置
func (u *Unit) Update() {
	// 如果是乘客，不需要独立更新位置
	if u.IsPassenger {
		return
	}

	// 更新攻击动画状态
	if u.IsAttacking {
		// 检查攻击动画是否结束
		if time.Since(u.AttackAnimTime).Seconds() >= u.AttackAnimDuration {
			u.IsAttacking = false
		}
	}

	// 如果有目标单位，更新目标位置
	if u.TargetUnit != nil && !u.TargetUnit.IsPassenger {
		// 检查目标单位是否还存在（未被消灭）
		if u.TargetUnit.Health > 0 {
			// 计算与目标单位的距离
			dist := abs(u.X-u.TargetUnit.X) + abs(u.Y-u.TargetUnit.Y)

			if dist <= u.AttackRange {
				// 如果在攻击范围内，停止移动并攻击
				u.HasTarget = false
				u.Path = make([][2]int, 0)
			} else {
				// 如果目标单位移动了，更新目标位置
				if u.TargetX != u.TargetUnit.X || u.TargetY != u.TargetUnit.Y {
					// 寻找靠近目标单位的空位置
					// 注意：这里需要游戏实例来调用findNearbyEmptyPositionForTarget
					// 我们将在Game的Update方法中处理这个逻辑
				}
			}
		} else {
			// 如果目标单位已被消灭，清除目标
			u.TargetUnit = nil
			u.HasTarget = false
			u.Path = make([][2]int, 0)
		}
	}

	// 如果有路径，沿着路径移动
	if len(u.Path) > 0 {
		// 获取下一个路径点
		nextPoint := u.Path[0]
		targetX := float64(nextPoint[0]*TileSize) + float64(TileSize)/2
		targetY := float64(nextPoint[1]*TileSize) + float64(TileSize)/2

		// 计算方向向量
		dx := targetX - u.CurrentX
		dy := targetY - u.CurrentY
		distance := math.Sqrt(dx*dx + dy*dy)

		// 如果已经非常接近目标点，则认为已到达
		if distance < u.MoveSpeed {
			u.X = nextPoint[0]
			u.Y = nextPoint[1]
			u.CurrentX = targetX
			u.CurrentY = targetY
			// 移除已到达的路径点
			u.Path = u.Path[1:]
			// 如果路径为空且已到达目标，清除目标（除非有目标单位）
			if len(u.Path) == 0 && u.X == u.TargetX && u.Y == u.TargetY && u.TargetUnit == nil {
				u.HasTarget = false
			}
		} else {
			// 移动单位（使用平滑的移动）
			moveSpeed := u.MoveSpeed

			// 计算单位向量
			normalizedDx := dx / distance
			normalizedDy := dy / distance

			// 应用移动
			u.CurrentX += normalizedDx * moveSpeed
			u.CurrentY += normalizedDy * moveSpeed

			// 更新网格位置
			u.X = int(u.CurrentX) / TileSize
			u.Y = int(u.CurrentY) / TileSize
		}
	}
}

func (u *Unit) Draw(screen *ebiten.Image, font font.Face) {
	// 如果是乘客，不绘制
	if u.IsPassenger {
		return
	}

	// 计算单位的实际大小
	unitSize := float64(TileSize) * u.Size

	// 如果正在攻击，应用收缩动画效果
	if u.IsAttacking {
		// 计算动画进度（0.0-1.0）
		animProgress := time.Since(u.AttackAnimTime).Seconds() / u.AttackAnimDuration

		// 创建一个收缩-恢复的动画效果
		// 在动画前半段收缩，后半段恢复
		var scaleFactor float64
		if animProgress < 0.5 {
			// 前半段：从1.0缩小到0.8
			scaleFactor = 1.0 - (0.2 * (animProgress / 0.5))
		} else {
			// 后半段：从0.8恢复到1.0
			scaleFactor = 0.8 + (0.2 * ((animProgress - 0.5) / 0.5))
		}

		// 应用缩放因子
		unitSize *= scaleFactor
	}

	// 计算单位的中心位置
	centerX := u.CurrentX
	centerY := u.CurrentY

	// 计算单位的左上角位置
	x := centerX - unitSize/2
	y := centerY - unitSize/2

	// 绘制单位
	var unitColor color.RGBA
	if u.IsPlayerUnit {
		unitColor = color.RGBA{0, 0, 255, 255} // 蓝色表示玩家单位
	} else {
		unitColor = color.RGBA{255, 0, 0, 255} // 红色表示敌方单位
	}

	// 绘制单位形状
	switch u.Type {
	case Infantry:
		// 步兵绘制为圆形
		vector.DrawFilledCircle(screen, float32(centerX), float32(centerY), float32(unitSize/2.5), unitColor, true)
	case Armor:
		// 装甲单位绘制为方形
		vector.DrawFilledRect(screen, float32(x), float32(y), float32(unitSize), float32(unitSize), unitColor, true)
	case Artillery:
		// 炮兵绘制为三角形
		vector.StrokeLine(screen, float32(centerX), float32(y), float32(x+unitSize), float32(y+unitSize), 3, unitColor, false)
		vector.StrokeLine(screen, float32(x+unitSize), float32(y+unitSize), float32(x), float32(y+unitSize), 3, unitColor, false)
		vector.StrokeLine(screen, float32(x), float32(y+unitSize), float32(centerX), float32(y), 3, unitColor, false)
		// 使用三个顶点绘制三角形
		vertices := []ebiten.Vertex{
			{
				DstX:   float32(centerX),
				DstY:   float32(y),
				SrcX:   0,
				SrcY:   0,
				ColorR: float32(unitColor.R) / 255,
				ColorG: float32(unitColor.G) / 255,
				ColorB: float32(unitColor.B) / 255,
				ColorA: float32(unitColor.A) / 255,
			},
			{
				DstX:   float32(x + unitSize),
				DstY:   float32(y + unitSize),
				SrcX:   0,
				SrcY:   0,
				ColorR: float32(unitColor.R) / 255,
				ColorG: float32(unitColor.G) / 255,
				ColorB: float32(unitColor.B) / 255,
				ColorA: float32(unitColor.A) / 255,
			},
			{
				DstX:   float32(x),
				DstY:   float32(y + unitSize),
				SrcX:   0,
				SrcY:   0,
				ColorR: float32(unitColor.R) / 255,
				ColorG: float32(unitColor.G) / 255,
				ColorB: float32(unitColor.B) / 255,
				ColorA: float32(unitColor.A) / 255,
			},
		}
		indices := []uint16{0, 1, 2}
		screen.DrawTriangles(vertices, indices, emptyImage, &ebiten.DrawTrianglesOptions{})
	case Recon:
		// 侦察单位绘制为菱形
		vector.DrawFilledCircle(screen, float32(centerX), float32(centerY), float32(unitSize/3), unitColor, true)
		vector.StrokeLine(screen, float32(centerX), float32(y), float32(x+unitSize), float32(centerY), 2, unitColor, false)
		vector.StrokeLine(screen, float32(x+unitSize), float32(centerY), float32(centerX), float32(y+unitSize), 2, unitColor, false)
		vector.StrokeLine(screen, float32(centerX), float32(y+unitSize), float32(x), float32(centerY), 2, unitColor, false)
		vector.StrokeLine(screen, float32(x), float32(centerY), float32(centerX), float32(y), 2, unitColor, false)
	case AntiAir:
		// 防空单位绘制为X形
		vector.StrokeLine(screen, float32(x), float32(y), float32(x+unitSize), float32(y+unitSize), 3, unitColor, false)
		vector.StrokeLine(screen, float32(x), float32(y+unitSize), float32(x+unitSize), float32(y), 3, unitColor, false)
	case HeavyTank:
		// 重型坦克绘制为大方形加炮塔
		vector.DrawFilledRect(screen, float32(x), float32(y), float32(unitSize), float32(unitSize), unitColor, true)
		// 绘制炮塔
		vector.DrawFilledCircle(screen, float32(centerX), float32(centerY), float32(unitSize/3), color.RGBA{
			R: unitColor.R - 40,
			G: unitColor.G - 40,
			B: unitColor.B - 40,
			A: unitColor.A,
		}, true)
		// 绘制炮管
		vector.StrokeLine(screen, float32(centerX), float32(centerY), float32(centerX+unitSize/2), float32(centerY), 4, color.RGBA{
			R: unitColor.R - 40,
			G: unitColor.G - 40,
			B: unitColor.B - 40,
			A: unitColor.A,
		}, false)
	case Helicopter:
		// 直升机绘制为圆形加旋翼
		vector.DrawFilledCircle(screen, float32(centerX), float32(centerY), float32(unitSize/3), unitColor, true)
		// 绘制旋翼
		vector.StrokeLine(screen, float32(x), float32(centerY), float32(x+unitSize), float32(centerY), 2, unitColor, false)
		vector.StrokeLine(screen, float32(centerX), float32(y), float32(centerX), float32(y+unitSize), 2, unitColor, false)
		// 绘制尾部
		vector.StrokeLine(screen, float32(centerX), float32(centerY), float32(centerX-unitSize/3), float32(centerY+unitSize/3), 2, unitColor, false)
	case FighterJet:
		// 战斗机绘制为三角形
		// 机身
		vector.DrawFilledCircle(screen, float32(centerX), float32(centerY), float32(unitSize/4), unitColor, true)
		// 机翼
		vector.StrokeLine(screen, float32(centerX-unitSize/2), float32(centerY), float32(centerX+unitSize/2), float32(centerY), 3, unitColor, false)
		// 机头
		vector.StrokeLine(screen, float32(centerX), float32(centerY), float32(centerX+unitSize/2), float32(centerY-unitSize/4), 2, unitColor, false)
		vector.StrokeLine(screen, float32(centerX+unitSize/2), float32(centerY-unitSize/4), float32(centerX+unitSize/2), float32(centerY), 2, unitColor, false)
	case Bomber:
		// 轰炸机绘制为大型机身加机翼
		// 机身
		vector.DrawFilledCircle(screen, float32(centerX), float32(centerY), float32(unitSize/3), unitColor, true)
		// 机翼
		vector.StrokeLine(screen, float32(centerX-unitSize/1.5), float32(centerY), float32(centerX+unitSize/1.5), float32(centerY), 4, unitColor, false)
		// 尾翼
		vector.StrokeLine(screen, float32(centerX-unitSize/3), float32(centerY), float32(centerX-unitSize/3), float32(centerY-unitSize/3), 2, unitColor, false)
	case RocketLauncher:
		// 火箭炮绘制为方形底座加发射管
		// 底座
		vector.DrawFilledRect(screen, float32(x+unitSize/4), float32(y+unitSize/2), float32(unitSize/2), float32(unitSize/2), unitColor, true)
		// 发射管
		vector.DrawFilledRect(screen, float32(x+unitSize/4), float32(y), float32(unitSize/2), float32(unitSize/2), color.RGBA{
			R: unitColor.R - 40,
			G: unitColor.G - 40,
			B: unitColor.B - 40,
			A: unitColor.A,
		}, true)
	case Engineer:
		// 工程兵绘制为圆形加工具标志
		vector.DrawFilledCircle(screen, float32(centerX), float32(centerY), float32(unitSize/2.5), unitColor, true)
		// 绘制工具标志（扳手形状）
		vector.StrokeLine(screen, float32(centerX-unitSize/4), float32(centerY-unitSize/4), float32(centerX+unitSize/4), float32(centerY+unitSize/4), 2, color.White, false)
		vector.StrokeLine(screen, float32(centerX-unitSize/4), float32(centerY+unitSize/4), float32(centerX+unitSize/4), float32(centerY-unitSize/4), 2, color.White, false)
	case MedicUnit:
		// 医疗单位绘制为圆形加十字标志
		vector.DrawFilledCircle(screen, float32(centerX), float32(centerY), float32(unitSize/2.5), unitColor, true)
		// 绘制医疗十字标志
		vector.StrokeLine(screen, float32(centerX), float32(centerY-unitSize/4), float32(centerX), float32(centerY+unitSize/4), 2, color.White, false)
		vector.StrokeLine(screen, float32(centerX-unitSize/4), float32(centerY), float32(centerX+unitSize/4), float32(centerY), 2, color.White, false)
	}

	// 如果单位被选中，绘制选择指示器
	if u.Selected {
		if u.IsPlayerUnit {
			// 友方单位选中指示器 - 绿色圆圈
			vector.StrokeCircle(screen, float32(centerX), float32(centerY), float32(unitSize/2+2), 2, color.RGBA{0, 255, 0, 255}, true)
		} else {
			// 敌方单位选中指示器 - 黄色十字准星
			crossSize := float32(unitSize / 2)
			// 水平线
			vector.StrokeLine(screen,
				float32(centerX)-crossSize, float32(centerY),
				float32(centerX)+crossSize, float32(centerY),
				2, color.RGBA{255, 255, 0, 255}, false)
			// 垂直线
			vector.StrokeLine(screen,
				float32(centerX), float32(centerY)-crossSize,
				float32(centerX), float32(centerY)+crossSize,
				2, color.RGBA{255, 255, 0, 255}, false)
			// 外圈
			vector.StrokeCircle(screen, float32(centerX), float32(centerY), float32(unitSize/2+4), 2, color.RGBA{255, 255, 0, 255}, true)
		}
	}

	// 显示单位生命值
	healthStr := fmt.Sprintf("%d", u.Health)
	text.Draw(screen, healthStr, font, int(centerX-float64(len(healthStr)*3)), int(centerY+3), color.White)

	// 显示单位防御力（如果有防御力）
	if u.Defense > 0 {
		defenseStr := fmt.Sprintf("防:%d", u.Defense)
		text.Draw(screen, defenseStr, font, int(centerX-float64(len(defenseStr)*3)), int(centerY+15), color.RGBA{220, 220, 100, 255})
	}

	// 显示单位弹药信息（如果有攻击能力）
	if u.AttackPower > 0 && u.MaxAmmo > 0 {
		ammoStr := fmt.Sprintf("弹:%d", u.CurrentAmmo)
		// 如果弹药不足，显示红色
		ammoColor := color.RGBA{100, 220, 100, 255} // 绿色
		if u.CurrentAmmo < u.MaxAmmo/4 {
			ammoColor = color.RGBA{220, 100, 100, 255} // 红色
		} else if u.CurrentAmmo < u.MaxAmmo/2 {
			ammoColor = color.RGBA{220, 220, 100, 255} // 黄色
		}
		text.Draw(screen, ammoStr, font, int(centerX-float64(len(ammoStr)*3)), int(centerY+27), ammoColor)
	}

	// 如果是装甲单位或直升机，显示乘客信息
	if u.Type == Armor || u.Type == Helicopter {
		passengerInfo := fmt.Sprintf("%d/%d", len(u.Passengers), u.MaxPassengers)
		text.Draw(screen, passengerInfo, font, int(centerX-float64(len(passengerInfo)*3)), int(centerY+39), color.White)
	}

	// 如果单位正在移动，绘制路径
	if u.HasTarget && len(u.Path) > 0 {
		// 绘制到下一个路径点的线
		ebitenutil.DrawLine(screen,
			float64(centerX),
			float64(centerY),
			float64(u.Path[0][0]*TileSize)+float64(TileSize)/2,
			float64(u.Path[0][1]*TileSize)+float64(TileSize)/2,
			color.RGBA{0, 255, 0, 128})

		// 绘制剩余路径
		for i := 0; i < len(u.Path)-1; i++ {
			ebitenutil.DrawLine(screen,
				float64(u.Path[i][0]*TileSize)+float64(TileSize)/2,
				float64(u.Path[i][1]*TileSize)+float64(TileSize)/2,
				float64(u.Path[i+1][0]*TileSize)+float64(TileSize)/2,
				float64(u.Path[i+1][1]*TileSize)+float64(TileSize)/2,
				color.RGBA{0, 255, 0, 128})
		}

		// 绘制目标点
		ebitenutil.DrawCircle(screen,
			float64(u.TargetX*TileSize)+float64(TileSize)/2,
			float64(u.TargetY*TileSize)+float64(TileSize)/2,
			5, color.RGBA{255, 255, 0, 255})
	}

	// 如果单位有攻击目标，绘制攻击线
	if u.TargetUnit != nil && u.IsPlayerUnit {
		// 绘制从单位到目标单位的攻击线
		ebitenutil.DrawLine(screen,
			float64(centerX),
			float64(centerY),
			float64(u.TargetUnit.CurrentX),
			float64(u.TargetUnit.CurrentY),
			color.RGBA{255, 0, 0, 200})

		// 在单位上方绘制攻击图标
		vector.StrokeLine(screen,
			float32(centerX-TileSize*0.2), float32(centerY-TileSize*0.4),
			float32(centerX+TileSize*0.2), float32(centerY-TileSize*0.4),
			2, color.RGBA{255, 0, 0, 255}, false)
		vector.StrokeLine(screen,
			float32(centerX), float32(centerY-TileSize*0.6),
			float32(centerX), float32(centerY-TileSize*0.2),
			2, color.RGBA{255, 0, 0, 255}, false)
	}
}

// IsAirUnit 判断单位是否为空中单位
func (u *Unit) IsAirUnit() bool {
	// 根据单位类型判断是否为空中单位
	return u.Type == Helicopter || u.Type == FighterJet || u.Type == Bomber
}

// CanAttack 检查单位是否可以攻击（基于攻击频率和弹药量）
func (u *Unit) CanAttack() bool {
	// 如果单位没有攻击能力，直接返回false
	if u.AttackPower <= 0 || u.AttackRange <= 0 {
		return false
	}

	// 检查弹药量
	if u.CurrentAmmo <= 0 {
		return false
	}

	// 计算攻击间隔（秒）
	attackInterval := 1.0 / u.AttackFrequency

	// 检查是否已经过了足够的时间
	timeSinceLastAttack := time.Since(u.LastAttackTime).Seconds()
	return timeSinceLastAttack >= attackInterval
}

// UpdateLastAttackTime 更新上次攻击时间
func (u *Unit) UpdateLastAttackTime() {
	u.LastAttackTime = time.Now()

	// 设置攻击动画状态
	u.IsAttacking = true
	u.AttackAnimTime = time.Now()
	u.AttackAnimDuration = 0.2 // 攻击动画持续0.2秒
}

// CalculateHitRate 根据距离计算命中率
// 距离越近，命中率越高；在最大攻击范围时，命中率等于MaxHitRate
// 在距离为1时，命中率为100%
func (u *Unit) CalculateHitRate(dist int) float64 {
	// 如果距离为1，则100%命中
	if dist <= 1 {
		return 1.0
	}

	// 如果超出攻击范围，则命中率为0
	if dist > u.AttackRange {
		return 0.0
	}

	// 线性插值计算命中率
	// 距离为1时，命中率为1.0
	// 距离为AttackRange时，命中率为MaxHitRate
	hitRateRange := 1.0 - u.MaxHitRate
	distanceRatio := float64(dist-1) / float64(u.AttackRange-1)

	// 命中率随距离线性减少
	return 1.0 - (hitRateRange * distanceRatio)
}

// RefillAmmo 补充弹药
func (u *Unit) RefillAmmo(amount int) {
	// 如果单位没有攻击能力或最大弹药量为0，则不需要补充
	if u.AttackPower <= 0 || u.MaxAmmo <= 0 {
		return
	}

	// 补充指定数量的弹药，但不超过最大弹药量
	u.CurrentAmmo += amount
	if u.CurrentAmmo > u.MaxAmmo {
		u.CurrentAmmo = u.MaxAmmo
	}
}

// RefillAmmoToMax 将弹药补充到最大值
func (u *Unit) RefillAmmoToMax() {
	// 如果单位没有攻击能力或最大弹药量为0，则不需要补充
	if u.AttackPower <= 0 || u.MaxAmmo <= 0 {
		return
	}

	// 将弹药补充到最大值
	u.CurrentAmmo = u.MaxAmmo
}

// CanMoveOnTerrain 检查单位是否可以在特定地形上行动
func (u *Unit) CanMoveOnTerrain(terrain TerrainType) bool {
	// 如果单位是空中单位，可以在任何地形上行动
	if u.IsAirUnit() {
		return true
	}

	// 获取单位类型数据
	unitData, err := GetUnitTypeData(u.Type)
	if err != nil {
		return false
	}

	// 检查地形是否在允许列表中
	for _, allowedTerrain := range unitData.AllowedTerrains {
		if allowedTerrain == terrain {
			return true
		}
	}

	return false
}
