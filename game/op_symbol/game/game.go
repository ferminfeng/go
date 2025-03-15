package game

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"sort"
	"time"

	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
)

const (
	TileSize    = 40
	MapWidth    = 30
	MapHeight   = 25
	CameraSpeed = 5 // 相机移动速度
)

type Game struct {
	gameMap       *GameMap
	units         []*Unit
	selectedUnits []*Unit // 选中的单位列表
	gameFont      font.Face
	showGrid      bool    // 是否显示网格
	cameraX       float64 // 相机X位置
	cameraY       float64 // 相机Y位置

	// 圈选相关
	isSelecting     bool // 是否正在圈选
	selectionStartX int  // 圈选起始X坐标
	selectionStartY int  // 圈选起始Y坐标
	selectionEndX   int  // 圈选结束X坐标
	selectionEndY   int  // 圈选结束Y坐标

	// 战争迷雾相关
	fogOfWar        [][]bool // 战争迷雾状态，true表示有迷雾
	visibleToPlayer [][]bool // 玩家可见区域，true表示可见

	// 攻击目标相关
	targetUnit      *Unit     // 当前攻击目标单位
	targetFlashTime time.Time // 攻击目标闪烁时间

	// 单位信息面板相关
	uiButtons     []UIButton // 操作按钮列表
	showUnitPanel bool       // 是否显示单位信息面板
}

// UIButton 表示界面上的按钮
type UIButton struct {
	X, Y, Width, Height int
	Text                string
	Action              string // 按钮对应的操作
	Enabled             bool   // 按钮是否可用
	Hovered             bool   // 鼠标是否悬停在按钮上
}

func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())

	// 初始化单位类型配置
	configPath := "./unit_types.json"
	err := InitUnitTypes(configPath)
	if err != nil {
		fmt.Printf("警告：无法加载单位类型配置，将使用默认值: %v\n", err)
		// 创建默认配置
		err = SaveDefaultUnitTypesConfig(configPath)
		if err != nil {
			fmt.Printf("警告：无法创建默认单位类型配置: %v\n", err)
		} else {
			// 重新加载配置
			err = LoadUnitTypesFromJSON(configPath)
			if err != nil {
				fmt.Printf("警告：无法重新加载单位类型配置: %v\n", err)
			}
		}
	}

	// 加载支持中文的字体
	var gameFont font.Face

	// 尝试加载unifont字体（支持中文）
	fontData, err := os.ReadFile("../assets/font/unifont-16.0.02.otf")
	if err != nil {
		fmt.Printf("警告：无法加载中文字体: %v\n", err)
		// 如果找不到中文字体，尝试加载原来的字体
		fontData, err = os.ReadFile("game/assets/fonts/font.ttf")
		if err != nil {
			// 如果找不到任何字体文件，使用默认字体
			fontData = nil
		}
	}

	if fontData != nil {
		tt, err := opentype.Parse(fontData)
		if err == nil {
			// 创建字体，设置合适的大小
			gameFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
				Size:    18, // 中文字体通常需要稍小一些的尺寸
				DPI:     72,
				Hinting: font.HintingFull,
			})
			if err != nil {
				fmt.Printf("警告：无法创建字体: %v\n", err)
				gameFont = nil
			}
		} else {
			fmt.Printf("警告：无法解析字体: %v\n", err)
		}
	}

	// 确保 gameFont 不为 nil
	if gameFont == nil {
		// 使用基本的默认字体
		gameFont = basicfont.Face7x13
		fmt.Println("警告：使用默认字体，可能无法正确显示中文")
	}

	// 初始化战争迷雾
	fogOfWar := make([][]bool, MapHeight)
	visibleToPlayer := make([][]bool, MapHeight)
	for y := 0; y < MapHeight; y++ {
		fogOfWar[y] = make([]bool, MapWidth)
		visibleToPlayer[y] = make([]bool, MapWidth)
		for x := 0; x < MapWidth; x++ {
			fogOfWar[y][x] = true         // 初始时所有区域都有迷雾
			visibleToPlayer[y][x] = false // 初始时所有区域都不可见
		}
	}

	// 初始化UI按钮
	uiButtons := []UIButton{
		{Text: "攻击", Action: "attack", Enabled: false},
		{Text: "移动", Action: "move", Enabled: false},
		{Text: "停止", Action: "stop", Enabled: false},
		{Text: "巡逻", Action: "patrol", Enabled: false},
		{Text: "装载", Action: "load", Enabled: false},
		{Text: "卸载", Action: "unload", Enabled: false},
	}
	showUnitPanel := true

	g := &Game{
		gameMap:         NewGameMap(MapWidth, MapHeight),
		units:           make([]*Unit, 0),
		selectedUnits:   make([]*Unit, 0),
		gameFont:        gameFont,
		showGrid:        false, // 默认不显示网格
		cameraX:         0,     // 初始相机位置
		cameraY:         0,
		isSelecting:     false,
		fogOfWar:        fogOfWar,
		visibleToPlayer: visibleToPlayer,
		targetUnit:      nil,
		targetFlashTime: time.Time{},
		uiButtons:       uiButtons,
		showUnitPanel:   showUnitPanel,
	}

	// 创建一些初始单位
	// 玩家单位 (左侧)
	g.units = append(g.units, NewUnit(3, 5, Infantry, true))
	g.units = append(g.units, NewUnit(3, 8, Infantry, true))
	g.units = append(g.units, NewUnit(3, 11, Infantry, true))
	g.units = append(g.units, NewUnit(2, 8, Armor, true))
	g.units = append(g.units, NewUnit(2, 5, Artillery, true))
	g.units = append(g.units, NewUnit(2, 11, Recon, true))

	// 敌方单位 (右侧)
	g.units = append(g.units, NewUnit(26, 5, Infantry, false))
	g.units = append(g.units, NewUnit(26, 8, Infantry, false))
	g.units = append(g.units, NewUnit(26, 11, Infantry, false))
	g.units = append(g.units, NewUnit(27, 8, Armor, false))
	g.units = append(g.units, NewUnit(27, 5, Artillery, false))
	g.units = append(g.units, NewUnit(27, 11, AntiAir, false))

	// 初始化玩家单位的视野
	g.updateFogOfWar()

	return g
}

func (g *Game) Update() error {
	// 清除过期的攻击目标
	if g.targetUnit != nil && time.Now().After(g.targetFlashTime) {
		g.targetUnit = nil
	}

	// 更新单位的持续攻击
	g.updateUnitAttacks()

	// 处理WASD键盘输入移动相机
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.cameraY -= CameraSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.cameraY += CameraSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.cameraX -= CameraSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.cameraX += CameraSpeed
	}

	// 获取屏幕尺寸
	screenWidth, screenHeight := ebiten.WindowSize()

	// 计算地图的实际像素尺寸
	mapWidthPixels := MapWidth * TileSize
	mapHeightPixels := MapHeight * TileSize

	// 限制相机范围，防止移出地图
	// 相机最大位置是地图尺寸减去屏幕尺寸
	g.cameraX = math.Max(0, math.Min(g.cameraX, float64(mapWidthPixels-screenWidth)))
	g.cameraY = math.Max(0, math.Min(g.cameraY, float64(mapHeightPixels-screenHeight)))

	// 处理鼠标输入
	// 获取鼠标位置（相对于屏幕）
	mouseX, mouseY := ebiten.CursorPosition()

	// 将鼠标位置转换为相对于地图的位置
	mapMouseX := mouseX + int(g.cameraX)
	mapMouseY := mouseY + int(g.cameraY)

	// 将鼠标位置转换为网格坐标
	gridX := mapMouseX / TileSize
	gridY := mapMouseY / TileSize

	// 处理鼠标左键点击
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// 开始圈选
		g.isSelecting = true
		g.selectionStartX = mapMouseX
		g.selectionStartY = mapMouseY
		g.selectionEndX = mapMouseX
		g.selectionEndY = mapMouseY
	}

	// 如果正在圈选，更新结束位置
	if g.isSelecting {
		g.selectionEndX = mapMouseX
		g.selectionEndY = mapMouseY
	}

	// 处理鼠标左键释放
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		// 结束圈选
		g.isSelecting = false

		// 计算选择框的网格范围
		startGridX := g.selectionStartX / TileSize
		startGridY := g.selectionStartY / TileSize
		endGridX := g.selectionEndX / TileSize
		endGridY := g.selectionEndY / TileSize

		// 确保范围有效
		minGridX := min(startGridX, endGridX)
		maxGridX := max(startGridX, endGridX)
		minGridY := min(startGridY, endGridY)
		maxGridY := max(startGridY, endGridY)

		// 如果选择范围很小，视为点击而非圈选
		if maxGridX-minGridX <= 1 && maxGridY-minGridY <= 1 {
			// 点击操作
			// 检查是否点击了单位
			clickedUnit := g.getUnitAt(gridX, gridY)

			if clickedUnit != nil {
				// 如果点击的是玩家单位
				if clickedUnit.IsPlayerUnit {
					// 如果按住Shift键，添加到选择列表
					if ebiten.IsKeyPressed(ebiten.KeyShift) {
						// 检查单位是否已经被选中
						alreadySelected := false
						for _, unit := range g.selectedUnits {
							if unit == clickedUnit {
								alreadySelected = true
								break
							}
						}

						// 如果未被选中，添加到选择列表
						if !alreadySelected {
							g.selectedUnits = append(g.selectedUnits, clickedUnit)
							clickedUnit.Selected = true
							fmt.Printf("选择了玩家单位：%s\n", clickedUnit.Name)
							// 对选中的单位按照entityID排序
							g.sortSelectedUnitsByEntityID()
						}
					} else {
						// 如果没有按Shift，清除当前选择并选择新单位
						for _, unit := range g.selectedUnits {
							unit.Selected = false
						}
						g.selectedUnits = []*Unit{clickedUnit}
						clickedUnit.Selected = true
						fmt.Printf("选择了玩家单位：%s\n", clickedUnit.Name)
						// 单个单位不需要排序
					}
				} else {
					// 如果点击的是敌方单位，并且有选中的玩家单位
					if len(g.selectedUnits) > 0 {
						// 检查敌方单位是否在可见区域内
						if g.visibleToPlayer[clickedUnit.Y][clickedUnit.X] {
							// 设置攻击目标和闪烁时间
							g.targetUnit = clickedUnit
							g.targetFlashTime = time.Now().Add(1 * time.Second) // 闪烁1秒

							// 命令选中的单位攻击敌方单位
							for _, unit := range g.selectedUnits {
								// 设置目标单位，启用持续攻击
								unit.TargetUnit = clickedUnit

								// 计算与目标单位的距离
								dist := abs(unit.X-clickedUnit.X) + abs(unit.Y-clickedUnit.Y)

								// 检查是否在攻击范围内
								if dist <= unit.AttackRange && unit.CanAttack() {
									// 计算命中率
									hitRate := unit.CalculateHitRate(dist)

									// 随机决定是否命中
									if rand.Float64() <= hitRate {
										// 检查攻击力是否大于防御力
										if unit.AttackPower > clickedUnit.Defense {
											// 计算实际伤害
											damage := unit.AttackPower - clickedUnit.Defense
											clickedUnit.Health -= damage
											fmt.Printf("%s攻击%s成功！造成%d点伤害，命中率：%.2f\n", unit.Name, clickedUnit.Name, damage, hitRate)

											// 更新上次攻击时间
											unit.UpdateLastAttackTime()

											// 如果敌方单位被消灭
											if clickedUnit.Health <= 0 {
												fmt.Printf("%s被消灭！\n", clickedUnit.Name)
												g.removeUnit(clickedUnit)
											}
										} else {
											fmt.Printf("%s攻击%s失败！攻击力(%d)不足以突破防御(%d)，命中率：%.2f\n", unit.Name, clickedUnit.Name, unit.AttackPower, clickedUnit.Defense, hitRate)

											// 更新上次攻击时间
											unit.UpdateLastAttackTime()
										}
									} else {
										fmt.Printf("%s攻击%s未命中！距离：%d，命中率：%.2f\n", unit.Name, clickedUnit.Name, dist, hitRate)
									}
								} else if dist <= unit.AttackRange {
									fmt.Printf("%s无法攻击，攻击冷却中！\n", unit.Name)
								} else {
									// 如果不在攻击范围内，设置为目标单位，并寻找路径
									unit.TargetUnit = clickedUnit
									// 寻找靠近目标单位的空位置
									targetPos := g.findNearbyEmptyPositionForTarget(clickedUnit.X, clickedUnit.Y, unit.Type)
									unit.TargetX = targetPos[0]
									unit.TargetY = targetPos[1]
									unit.HasTarget = true
									// 使用A*寻路
									unit.Path = FindPathWithUnits(g.gameMap, g, unit.X, unit.Y, unit.TargetX, unit.TargetY, unit.Type)
								}
							}
						}
					} else {
						// 如果没有选中的玩家单位，选择敌方单位以查看信息
						// 只有当敌方单位在可见区域内才能选择
						if g.visibleToPlayer[clickedUnit.Y][clickedUnit.X] {
							for _, unit := range g.selectedUnits {
								unit.Selected = false
							}
							g.selectedUnits = []*Unit{clickedUnit}
							clickedUnit.Selected = true
							// 单个单位不需要排序
						} else {
							fmt.Println("无法选择：目标单位在迷雾中！")
						}
					}
				}
			} else {
				// 如果点击了空地，并且有选中的单位
				if len(g.selectedUnits) > 0 {
					// 检查是否所有选中的单位都是玩家单位
					allPlayerUnits := true
					for _, unit := range g.selectedUnits {
						if !unit.IsPlayerUnit {
							allPlayerUnits = false
							break
						}
					}

					// 只有当所有选中的单位都是玩家单位时，才执行移动命令
					if allPlayerUnits {
						// 左键点击空地不再控制单位移动，只取消选择
						fmt.Println("左键点击空地：取消选择所有单位")
						for _, unit := range g.selectedUnits {
							unit.Selected = false
						}
						g.selectedUnits = []*Unit{}
					} else {
						// 如果选中的单位包含敌方单位，取消选择
						for _, unit := range g.selectedUnits {
							unit.Selected = false
						}
						g.selectedUnits = []*Unit{}
					}
				}
			}
		} else {
			// 圈选操作
			// 清除当前选择
			for _, unit := range g.selectedUnits {
				unit.Selected = false
			}
			g.selectedUnits = []*Unit{}

			// 选择范围内的所有玩家单位
			for _, unit := range g.units {
				if unit.IsPlayerUnit && !unit.IsPassenger {
					// 检查单位是否在选择范围内
					if unit.X >= minGridX && unit.X <= maxGridX && unit.Y >= minGridY && unit.Y <= maxGridY {
						g.selectedUnits = append(g.selectedUnits, unit)
						unit.Selected = true
					}
				}
			}

			// 对选中的单位按照entityID排序
			g.sortSelectedUnitsByEntityID()
		}
	}

	// 处理鼠标右键点击（命令攻击或移动）
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		// 如果有选中的玩家单位
		if len(g.selectedUnits) > 0 {
			// 检查是否所有选中的单位都是玩家单位
			allPlayerUnits := true
			for _, unit := range g.selectedUnits {
				if !unit.IsPlayerUnit {
					allPlayerUnits = false
					break
				}
			}

			// 只有当所有选中的单位都是玩家单位时，才执行命令
			if allPlayerUnits {
				// 检查点击位置是否有敌方单位
				clickedUnit := g.getUnitAt(gridX, gridY)

				if clickedUnit != nil && !clickedUnit.IsPlayerUnit {
					// 检查敌方单位是否在可见区域内
					if g.visibleToPlayer[clickedUnit.Y][clickedUnit.X] {
						fmt.Printf("右键命令：攻击敌方单位 %s\n", clickedUnit.Name)
						// 设置攻击目标和闪烁时间
						g.targetUnit = clickedUnit
						g.targetFlashTime = time.Now().Add(1 * time.Second) // 闪烁1秒

						// 命令选中的单位攻击敌方单位
						for _, unit := range g.selectedUnits {
							// 设置目标单位，启用持续攻击
							unit.TargetUnit = clickedUnit

							// 计算与目标单位的距离
							dist := abs(unit.X-clickedUnit.X) + abs(unit.Y-clickedUnit.Y)

							// 检查是否在攻击范围内
							if dist <= unit.AttackRange && unit.CanAttack() {
								// 计算命中率
								hitRate := unit.CalculateHitRate(dist)

								// 随机决定是否命中
								if rand.Float64() <= hitRate {
									// 检查攻击力是否大于防御力
									if unit.AttackPower > clickedUnit.Defense {
										// 计算实际伤害
										damage := unit.AttackPower - clickedUnit.Defense
										clickedUnit.Health -= damage
										fmt.Printf("%s攻击%s成功！造成%d点伤害，命中率：%.2f\n", unit.Name, clickedUnit.Name, damage, hitRate)

										// 更新上次攻击时间
										unit.UpdateLastAttackTime()

										// 如果敌方单位被消灭
										if clickedUnit.Health <= 0 {
											fmt.Printf("%s被消灭！\n", clickedUnit.Name)
											g.removeUnit(clickedUnit)
										}
									} else {
										fmt.Printf("%s攻击%s失败！攻击力(%d)不足以突破防御(%d)，命中率：%.2f\n", unit.Name, clickedUnit.Name, unit.AttackPower, clickedUnit.Defense, hitRate)

										// 更新上次攻击时间
										unit.UpdateLastAttackTime()
									}
								} else {
									fmt.Printf("%s攻击%s未命中！距离：%d，命中率：%.2f\n", unit.Name, clickedUnit.Name, dist, hitRate)
								}
							} else if dist <= unit.AttackRange {
								fmt.Printf("%s无法攻击，攻击冷却中！\n", unit.Name)
							} else {
								// 如果不在攻击范围内，设置为目标单位，并寻找路径
								unit.TargetUnit = clickedUnit
								// 寻找靠近目标单位的空位置
								targetPos := g.findNearbyEmptyPositionForTarget(clickedUnit.X, clickedUnit.Y, unit.Type)
								if targetPos[0] != -1 && targetPos[1] != -1 {
									unit.TargetX = targetPos[0]
									unit.TargetY = targetPos[1]
									unit.HasTarget = true
									// 使用A*寻路
									unit.Path = FindPathWithUnits(g.gameMap, g, unit.X, unit.Y, unit.TargetX, unit.TargetY, unit.Type)
									fmt.Printf("右键命令单位 %s 移动到敌方单位 %s 附近\n", unit.Name, clickedUnit.Name)
								} else {
									fmt.Printf("无法找到靠近敌方单位 %s 的空位置\n", clickedUnit.Name)
								}
							}
						}
					} else {
						fmt.Println("无法攻击：目标单位在迷雾中！")
					}
				} else {
					// 如果点击了空地，命令选中的单位移动到点击位置
					fmt.Println("右键命令：移动到指定位置")
					for _, unit := range g.selectedUnits {
						// 清除目标单位
						unit.TargetUnit = nil
						// 设置目标位置
						unit.TargetX = gridX
						unit.TargetY = gridY
						unit.HasTarget = true
						// 使用A*寻路
						unit.Path = FindPathWithUnits(g.gameMap, g, unit.X, unit.Y, unit.TargetX, unit.TargetY, unit.Type)
						fmt.Printf("右键移动单位 %s 到位置 (%d, %d)\n", unit.Name, gridX, gridY)
					}
				}
			}
		}
	}

	// 更新所有单位
	for _, unit := range g.units {
		// 如果是玩家单位，或者在可见区域内的敌方单位，正常更新
		if unit.IsPlayerUnit || g.isVisibleToPlayer(unit.X, unit.Y) {
			// 保存当前位置
			oldX, oldY := unit.X, unit.Y
			oldCurrentX, oldCurrentY := unit.CurrentX, unit.CurrentY

			// 更新单位位置
			unit.Update()

			// 检查碰撞
			hasCollision := false
			var collidedUnit *Unit
			for _, otherUnit := range g.units {
				if unit != otherUnit && !otherUnit.IsPassenger && !unit.IsPassenger && g.unitsCollide(unit, otherUnit) {
					hasCollision = true
					collidedUnit = otherUnit
					break
				}
			}

			// 如果发生碰撞，恢复原位置
			if hasCollision {
				unit.X, unit.Y = oldX, oldY
				unit.CurrentX, unit.CurrentY = oldCurrentX, oldCurrentY

				// 如果是玩家单位，提供碰撞反馈
				if unit.IsPlayerUnit {
					fmt.Printf("单位 %s 与 %s 发生碰撞，无法移动！\n", unit.Name, collidedUnit.Name)
				}

				// 如果发生碰撞，尝试寻找新路径
				if unit.HasTarget && len(unit.Path) > 0 {
					// 寻找新的路径
					unit.Path = FindPathWithUnits(g.gameMap, g, unit.X, unit.Y, unit.TargetX, unit.TargetY, unit.Type)
				}
			}
		} else {
			// 如果是不可见区域内的敌方单位，只更新移动，不更新攻击
			// 这样敌方单位仍然会移动，但不会攻击玩家单位，除非被发现
			if len(unit.Path) > 0 {
				// 获取下一个路径点
				nextPoint := unit.Path[0]
				targetX := float64(nextPoint[0]*TileSize) + float64(TileSize)/2
				targetY := float64(nextPoint[1]*TileSize) + float64(TileSize)/2

				// 计算方向向量
				dx := targetX - unit.CurrentX
				dy := targetY - unit.CurrentY
				distance := math.Sqrt(dx*dx + dy*dy)

				// 如果已经非常接近目标点，则认为已到达
				if distance < unit.MoveSpeed {
					unit.X = nextPoint[0]
					unit.Y = nextPoint[1]
					unit.CurrentX = targetX
					unit.CurrentY = targetY
					// 移除已到达的路径点
					unit.Path = unit.Path[1:]
					// 如果路径为空且已到达目标，清除目标
					if len(unit.Path) == 0 && unit.X == unit.TargetX && unit.Y == unit.TargetY {
						unit.HasTarget = false
					}
				} else {
					// 移动单位（使用平滑的移动）
					moveSpeed := unit.MoveSpeed

					// 计算单位向量
					normalizedDx := dx / distance
					normalizedDy := dy / distance

					// 应用移动
					unit.CurrentX += normalizedDx * moveSpeed
					unit.CurrentY += normalizedDy * moveSpeed

					// 更新网格位置
					unit.X = int(unit.CurrentX) / TileSize
					unit.Y = int(unit.CurrentY) / TileSize
				}
			}

			// 清除目标单位，如果目标单位是玩家单位
			if unit.TargetUnit != nil && unit.TargetUnit.IsPlayerUnit {
				unit.TargetUnit = nil
			}
		}
	}

	// 更新AI
	g.updateAI()

	// 更新战争迷雾
	g.updateFogOfWar()

	// 更新UI按钮状态
	g.updateUIButtons()

	// 处理UI按钮点击
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.handleUIButtonClick(x, y)
	}

	return nil
}

// 检查位置是否被占用
func (g *Game) isPositionOccupied(x, y int) bool {
	// 将网格坐标转换为像素坐标（中心点）
	posX := float64(x*TileSize) + float64(TileSize)/2
	posY := float64(y*TileSize) + float64(TileSize)/2

	for _, unit := range g.units {
		if unit.IsPassenger {
			continue
		}

		// 计算与单位中心的距离
		dx := posX - unit.CurrentX
		dy := posY - unit.CurrentY
		distance := math.Sqrt(dx*dx + dy*dy)

		// 获取单位的碰撞半径
		unitRadius := getCollisionRadius(unit)

		// 如果距离小于单位的碰撞半径加上半个格子大小，则认为位置被占用
		if distance < (unitRadius + float64(TileSize)*0.4) {
			return true
		}
	}
	return false
}

// 检查位置是否对玩家可见
func (g *Game) isVisibleToPlayer(x, y int) bool {
	// 检查坐标是否在地图范围内
	if x < 0 || x >= MapWidth || y < 0 || y >= MapHeight {
		return false
	}
	return g.visibleToPlayer[y][x]
}

// 检查两个单位是否碰撞
func (g *Game) unitsCollide(unit1, unit2 *Unit) bool {
	// 如果任一单位是乘客，则不会碰撞
	if unit1.IsPassenger || unit2.IsPassenger {
		return false
	}

	// 计算单位中心点之间的距离
	unit1CenterX := unit1.CurrentX
	unit1CenterY := unit1.CurrentY
	unit2CenterX := unit2.CurrentX
	unit2CenterY := unit2.CurrentY

	// 计算距离
	dx := unit1CenterX - unit2CenterX
	dy := unit1CenterY - unit2CenterY
	distance := math.Sqrt(dx*dx + dy*dy)

	// 根据单位类型和大小计算碰撞半径
	unit1Radius := getCollisionRadius(unit1)
	unit2Radius := getCollisionRadius(unit2)

	// 如果距离小于两个单位的碰撞半径之和，则发生碰撞
	return distance < (unit1Radius + unit2Radius)
}

// 根据单位类型获取碰撞半径
func getCollisionRadius(unit *Unit) float64 {
	// 基础碰撞系数
	collisionFactor := 0.5

	// 根据单位类型调整碰撞系数
	switch unit.Type {
	case Infantry:
		collisionFactor = 0.45 // 步兵稍小一些
	case Armor, HeavyTank:
		collisionFactor = 0.55 // 装甲车和重型坦克更大一些
	case Artillery, RocketLauncher:
		collisionFactor = 0.52 // 火炮和火箭炮
	case Helicopter, FighterJet, Bomber:
		collisionFactor = 0.5 // 空中单位
	default:
		collisionFactor = 0.5 // 默认值
	}

	// 计算并返回碰撞半径
	return float64(TileSize) * unit.Size * collisionFactor
}

// 寻找附近的空位置
func (g *Game) findNearbyEmptyPosition(unit *Unit) {
	// 检查周围8个方向
	directions := [][2]int{
		{0, -1}, {1, -1}, {1, 0}, {1, 1},
		{0, 1}, {-1, 1}, {-1, 0}, {-1, -1},
	}

	for _, dir := range directions {
		newX := unit.X + dir[0]
		newY := unit.Y + dir[1]

		// 检查是否在地图范围内
		if newX >= 0 && newX < MapWidth && newY >= 0 && newY < MapHeight {
			// 检查是否被占用
			if !g.isPositionOccupied(newX, newY) {
				// 检查地形是否可通行
				if isWalkable(g.gameMap, newX, newY, unit.Type) {
					unit.X = newX
					unit.Y = newY
					unit.CurrentX = float64(newX*TileSize) + float64(TileSize)/2
					unit.CurrentY = float64(newY*TileSize) + float64(TileSize)/2
					return
				}
			}
		}
	}
}

// 为目标点寻找附近的空位置
func (g *Game) findNearbyEmptyPositionForTarget(x, y int, unitType UnitType) [2]int {
	// 检查周围8个方向
	directions := [][2]int{
		{0, -1}, {1, -1}, {1, 0}, {1, 1},
		{0, 1}, {-1, 1}, {-1, 0}, {-1, -1},
	}

	for _, dir := range directions {
		newX := x + dir[0]
		newY := y + dir[1]

		// 检查是否在地图范围内
		if newX >= 0 && newX < MapWidth && newY >= 0 && newY < MapHeight {
			// 检查是否被占用
			if !g.isPositionOccupied(newX, newY) {
				// 检查地形是否可通行
				if isWalkable(g.gameMap, newX, newY, unitType) {
					return [2]int{newX, newY}
				}
			}
		}
	}

	// 如果没有找到合适的位置，返回无效坐标
	return [2]int{-1, -1}
}

// 简单的AI逻辑
func (g *Game) updateAI() {
	// 每隔一段时间让AI单位移动
	if rand.Intn(120) == 0 { // 大约每2秒
		for _, unit := range g.units {
			if !unit.IsPlayerUnit && !unit.HasTarget && !unit.IsPassenger {
				// 随机选择一个目标位置
				targetX := rand.Intn(MapWidth)
				targetY := rand.Intn(MapHeight)

				// 设置目标并计算路径
				unit.TargetX = targetX
				unit.TargetY = targetY
				unit.HasTarget = true
				unit.Path = FindPathWithUnits(g.gameMap, g, unit.X, unit.Y, targetX, targetY, unit.Type)
			}
		}
	}

	// 注意：攻击逻辑已移至updateUnitAttacks函数中，实现了自动攻击功能
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 创建一个临时画布，用于应用相机偏移
	canvas := ebiten.NewImage(MapWidth*TileSize, MapHeight*TileSize)

	// 获取屏幕尺寸
	screenWidth, screenHeight := ebiten.WindowSize()

	// 绘制地图
	g.gameMap.Draw(canvas)

	// 如果需要显示网格，绘制网格线
	if g.showGrid {
		for y := 0; y <= g.gameMap.Height; y++ {
			ebitenutil.DrawLine(canvas,
				0,
				float64(y*TileSize),
				float64(g.gameMap.Width*TileSize),
				float64(y*TileSize),
				color.RGBA{50, 50, 50, 100})
		}

		for x := 0; x <= g.gameMap.Width; x++ {
			ebitenutil.DrawLine(canvas,
				float64(x*TileSize),
				0,
				float64(x*TileSize),
				float64(g.gameMap.Height*TileSize),
				color.RGBA{50, 50, 50, 100})
		}
	}

	// 绘制单位（只绘制可见区域的单位）
	for _, unit := range g.units {
		// 只绘制玩家单位或者在可见区域内的敌方单位
		if unit.IsPlayerUnit || g.visibleToPlayer[unit.Y][unit.X] {
			unit.Draw(canvas, g.gameFont)
		}
	}

	// 绘制攻击目标指示器
	if g.targetUnit != nil && time.Now().Before(g.targetFlashTime) {
		// 只有当目标单位在可见区域内才绘制
		if g.visibleToPlayer[g.targetUnit.Y][g.targetUnit.X] {
			// 计算单位中心位置
			centerX := g.targetUnit.CurrentX
			centerY := g.targetUnit.CurrentY

			// 绘制攻击目标指示器（红色十字准星和闪烁圆圈）
			// 十字准星
			crossSize := float32(TileSize * 0.6)
			vector.StrokeLine(canvas,
				float32(centerX)-crossSize, float32(centerY),
				float32(centerX)+crossSize, float32(centerY),
				3, color.RGBA{255, 0, 0, 255}, false)
			vector.StrokeLine(canvas,
				float32(centerX), float32(centerY)-crossSize,
				float32(centerX), float32(centerY)+crossSize,
				3, color.RGBA{255, 0, 0, 255}, false)

			// 闪烁圆圈
			// 根据时间计算闪烁效果
			alpha := uint8(255 * (g.targetFlashTime.Sub(time.Now()).Seconds() / 1.0))
			vector.StrokeCircle(canvas, float32(centerX), float32(centerY), float32(TileSize*0.8), 3, color.RGBA{255, 0, 0, alpha}, true)
		}
	}

	// 绘制迷雾战
	for y := 0; y < MapHeight; y++ {
		for x := 0; x < MapWidth; x++ {
			if g.fogOfWar[y][x] {
				// 绘制完全迷雾（黑色）
				ebitenutil.DrawRect(canvas,
					float64(x*TileSize), float64(y*TileSize),
					float64(TileSize), float64(TileSize),
					color.RGBA{0, 0, 0, 200})
			} else if !g.visibleToPlayer[y][x] {
				// 绘制已探索但当前不可见的区域（灰色）
				ebitenutil.DrawRect(canvas,
					float64(x*TileSize), float64(y*TileSize),
					float64(TileSize), float64(TileSize),
					color.RGBA{100, 100, 100, 150})
			}
		}
	}

	// 绘制圈选框
	if g.isSelecting {
		minX := float64(min(g.selectionStartX, g.selectionEndX))
		maxX := float64(max(g.selectionStartX, g.selectionEndX))
		minY := float64(min(g.selectionStartY, g.selectionEndY))
		maxY := float64(max(g.selectionStartY, g.selectionEndY))

		// 绘制半透明的选择框
		vector.StrokeRect(canvas, float32(minX), float32(minY), float32(maxX-minX), float32(maxY-minY), 2, color.RGBA{0, 255, 0, 255}, false)
		ebitenutil.DrawRect(canvas, minX, minY, maxX-minX, maxY-minY, color.RGBA{0, 255, 0, 50})
	}

	// 将画布绘制到屏幕上，应用相机偏移
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-g.cameraX, -g.cameraY)
	screen.DrawImage(canvas, op)

	// 获取屏幕尺寸
	screenWidth, screenHeight = ebiten.WindowSize()

	// 如果有选中的单位，绘制单位信息面板
	if g.showUnitPanel && len(g.selectedUnits) > 0 {
		g.drawUnitInfoPanel(screen, screenWidth, screenHeight)
	}
}

// 绘制游戏基本信息 - 已移除，不再显示左侧提示
func (g *Game) drawGameInfo(screen *ebiten.Image) {
	// 此函数保留但不再绘制任何内容
}

// 绘制单位信息面板
func (g *Game) drawUnitInfoPanel(screen *ebiten.Image, screenWidth, screenHeight int) {
	// 计算信息面板的位置和大小
	panelWidth := screenWidth - 40            // 左右各留20像素边距
	panelHeight := 250                        // 增加面板高度，确保所有内容都在边框内
	panelX := 20                              // 左边距
	panelY := screenHeight - panelHeight - 20 // 底部位置，留20像素边距

	// 绘制信息面板背景
	ebitenutil.DrawRect(screen,
		float64(panelX), float64(panelY),
		float64(panelWidth), float64(panelHeight),
		color.RGBA{0, 0, 0, 200}) // 半透明黑色背景

	// 绘制信息面板边框
	vector.StrokeRect(screen,
		float32(panelX), float32(panelY),
		float32(panelWidth), float32(panelHeight),
		2, color.RGBA{255, 165, 0, 255}, false) // 橙色边框

	// 显示选中单位数量
	selectedText := fmt.Sprintf("已选中: %d 个单位", len(g.selectedUnits))
	text.Draw(screen, selectedText, g.gameFont, panelX+10, panelY+30, color.RGBA{255, 255, 255, 255})

	// 如果只选中了一个单位，显示详细信息
	if len(g.selectedUnits) == 1 {
		unit := g.selectedUnits[0]

		// 获取单位类型数据
		unitTypeData, err := GetUnitTypeData(unit.Type)
		unitTypeName := "未知"
		if err == nil {
			unitTypeName = unitTypeData.Name
		} else {
			// 使用备用方式获取单位类型名称
			switch unit.Type {
			case Infantry:
				unitTypeName = "步兵"
			case Armor:
				unitTypeName = "装甲车"
			case Artillery:
				unitTypeName = "火炮"
			case Recon:
				unitTypeName = "侦察车"
			case AntiAir:
				unitTypeName = "防空炮"
			case HeavyTank:
				unitTypeName = "重型坦克"
			case Helicopter:
				unitTypeName = "直升机"
			case FighterJet:
				unitTypeName = "战斗机"
			case Bomber:
				unitTypeName = "轰炸机"
			case RocketLauncher:
				unitTypeName = "火箭炮"
			case Engineer:
				unitTypeName = "工程兵"
			case MedicUnit:
				unitTypeName = "医疗单位"
			}
		}

		// 绘制单位图标（简化版，使用彩色矩形代替）
		iconSize := 64
		iconX := panelX + 20
		iconY := panelY + 50

		// 根据单位类型选择不同的颜色
		iconColor := color.RGBA{100, 100, 255, 255} // 默认蓝色
		if unit.IsAirUnit() {
			iconColor = color.RGBA{100, 200, 255, 255} // 空中单位使用浅蓝色
		}
		if !unit.IsPlayerUnit {
			iconColor = color.RGBA{255, 100, 100, 255} // 敌方单位使用红色
		}

		// 绘制单位图标背景
		ebitenutil.DrawRect(screen,
			float64(iconX), float64(iconY),
			float64(iconSize), float64(iconSize),
			iconColor)

		// 绘制单位图标边框
		vector.StrokeRect(screen,
			float32(iconX), float32(iconY),
			float32(iconSize), float32(iconSize),
			2, color.RGBA{255, 255, 255, 255}, false) // 将color.White改为明确的白色RGBA

		// 绘制单位类型标识（简化版，使用字母代替）
		var typeChar string
		switch unit.Type {
		case Infantry:
			typeChar = "步"
		case Armor:
			typeChar = "装"
		case Artillery:
			typeChar = "炮"
		case Recon:
			typeChar = "侦"
		case AntiAir:
			typeChar = "防"
		case HeavyTank:
			typeChar = "坦"
		case Helicopter:
			typeChar = "直"
		case FighterJet:
			typeChar = "战"
		case Bomber:
			typeChar = "轰"
		case RocketLauncher:
			typeChar = "火"
		case Engineer:
			typeChar = "工"
		case MedicUnit:
			typeChar = "医"
		}
		text.Draw(screen, typeChar, g.gameFont, iconX+iconSize/2-10, iconY+iconSize/2+10, color.Black)

		// 绘制单位详细信息
		infoX := iconX + iconSize + 30
		infoY := iconY - 10

		// 单位名称
		text.Draw(screen, unitTypeName, g.gameFont, infoX, infoY+20, color.RGBA{255, 255, 255, 255})

		// 单位阵营
		teamText := "阵营: "
		if unit.IsPlayerUnit {
			teamText += "友方"
			text.Draw(screen, teamText, g.gameFont, infoX, infoY+45, color.RGBA{100, 100, 255, 255})
		} else {
			teamText += "敌方"
			text.Draw(screen, teamText, g.gameFont, infoX, infoY+45, color.RGBA{255, 100, 100, 255})
		}

		// 生命值
		healthBarWidth := 150
		healthBarHeight := 15
		healthPercentage := float64(unit.Health) / float64(unitTypeData.Health)
		if healthPercentage > 1.0 {
			healthPercentage = 1.0
		}

		// 绘制生命值背景
		ebitenutil.DrawRect(screen,
			float64(infoX), float64(infoY+55),
			float64(healthBarWidth), float64(healthBarHeight),
			color.RGBA{50, 50, 50, 255})

		// 绘制生命值条
		healthColor := color.RGBA{0, 255, 0, 255} // 绿色
		if healthPercentage < 0.3 {
			healthColor = color.RGBA{255, 0, 0, 255} // 红色
		} else if healthPercentage < 0.6 {
			healthColor = color.RGBA{255, 255, 0, 255} // 黄色
		}

		ebitenutil.DrawRect(screen,
			float64(infoX), float64(infoY+55),
			float64(healthBarWidth)*healthPercentage, float64(healthBarHeight),
			healthColor)

		// 显示生命值文本
		healthText := fmt.Sprintf("生命值: %d/%d", unit.Health, unitTypeData.Health)
		text.Draw(screen, healthText, g.gameFont, infoX+healthBarWidth+10, infoY+68, color.RGBA{255, 255, 255, 255})

		// 显示攻击力
		attackText := fmt.Sprintf("攻击力: %d", unit.AttackPower)
		text.Draw(screen, attackText, g.gameFont, infoX, infoY+90, color.RGBA{255, 150, 150, 255})

		// 显示防御力
		defenseText := fmt.Sprintf("防御力: %d", unit.Defense)
		text.Draw(screen, defenseText, g.gameFont, infoX, infoY+115, color.RGBA{150, 150, 255, 255})

		// 显示攻击范围
		rangeText := fmt.Sprintf("攻击范围: %d", unit.AttackRange)
		text.Draw(screen, rangeText, g.gameFont, infoX, infoY+140, color.RGBA{200, 200, 200, 255})

		// 显示弹药信息（如果有）
		if unit.MaxAmmo > 0 {
			ammoText := fmt.Sprintf("弹药: %d/%d", unit.CurrentAmmo, unit.MaxAmmo)
			text.Draw(screen, ammoText, g.gameFont, infoX, infoY+165, color.RGBA{255, 255, 150, 255})
		}

		// 如果是载具，显示乘客信息
		if unit.MaxPassengers > 0 {
			passengerText := fmt.Sprintf("乘客: %d/%d", len(unit.Passengers), unit.MaxPassengers)
			text.Draw(screen, passengerText, g.gameFont, infoX+200, infoY+90, color.RGBA{150, 255, 150, 255})
		}

		// 在右侧显示更多信息
		rightInfoX := infoX + 400

		// 显示单位位置
		posText := fmt.Sprintf("位置: (%d, %d)", unit.X, unit.Y)
		text.Draw(screen, posText, g.gameFont, rightInfoX, infoY+45, color.RGBA{200, 200, 200, 255})

		// 显示单位速度
		speedText := fmt.Sprintf("速度: %.1f", unit.MoveSpeed)
		text.Draw(screen, speedText, g.gameFont, rightInfoX, infoY+70, color.RGBA{200, 200, 200, 255})

		// 显示单位视野
		visionText := fmt.Sprintf("视野: %d", unit.VisionRange)
		text.Draw(screen, visionText, g.gameFont, rightInfoX, infoY+95, color.RGBA{200, 200, 200, 255})
	} else if len(g.selectedUnits) > 1 {
		// 如果选中了多个单位，显示单位类型统计
		unitTypeCounts := make(map[UnitType]int)
		for _, unit := range g.selectedUnits {
			unitTypeCounts[unit.Type]++
		}

		// 创建一个有序的UnitType切片，确保每次显示顺序一致
		unitTypes := []UnitType{
			Infantry, Armor, Artillery, Recon, AntiAir, HeavyTank,
			Helicopter, FighterJet, Bomber, RocketLauncher, Engineer, MedicUnit,
		}

		// 计算每列显示的单位类型数量
		typesPerColumn := 6
		columnWidth := 200

		// 按照固定顺序显示单位类型
		typeCount := 0
		for _, unitType := range unitTypes {
			// 只显示存在的单位类型
			if count, exists := unitTypeCounts[unitType]; exists {
				var typeName string
				switch unitType {
				case Infantry:
					typeName = "步兵"
				case Armor:
					typeName = "装甲车"
				case Artillery:
					typeName = "火炮"
				case Recon:
					typeName = "侦察车"
				case AntiAir:
					typeName = "防空炮"
				case HeavyTank:
					typeName = "重型坦克"
				case Helicopter:
					typeName = "直升机"
				case FighterJet:
					typeName = "战斗机"
				case Bomber:
					typeName = "轰炸机"
				case RocketLauncher:
					typeName = "火箭炮"
				case Engineer:
					typeName = "工程兵"
				case MedicUnit:
					typeName = "医疗单位"
				}

				// 计算当前单位类型应该显示在哪一列
				column := typeCount / typesPerColumn
				row := typeCount % typesPerColumn

				countText := fmt.Sprintf("%s: %d个", typeName, count)
				text.Draw(screen, countText, g.gameFont,
					panelX+40+column*columnWidth,
					panelY+60+row*25,
					color.RGBA{255, 255, 255, 255})

				typeCount++
			}
		}
	}

	// 绘制操作按钮
	buttonWidth := 80
	buttonHeight := 30
	buttonSpacing := 10

	// 计算按钮区域的起始位置（面板右侧）
	buttonsAreaWidth := len(g.uiButtons) * (buttonWidth + buttonSpacing)
	buttonsStartX := panelX + panelWidth - buttonsAreaWidth - 20 // 右侧留20像素边距
	buttonY := panelY + panelHeight - buttonHeight - 15          // 底部留15像素边距

	// 更新按钮位置
	for i := range g.uiButtons {
		g.uiButtons[i].X = buttonsStartX + (buttonWidth+buttonSpacing)*i
		g.uiButtons[i].Y = buttonY
		g.uiButtons[i].Width = buttonWidth
		g.uiButtons[i].Height = buttonHeight

		// 绘制按钮背景
		buttonColor := color.RGBA{50, 50, 50, 200}
		if !g.uiButtons[i].Enabled {
			buttonColor = color.RGBA{30, 30, 30, 200}
		} else if g.uiButtons[i].Hovered {
			buttonColor = color.RGBA{80, 80, 80, 200}
		}

		ebitenutil.DrawRect(screen,
			float64(g.uiButtons[i].X), float64(g.uiButtons[i].Y),
			float64(g.uiButtons[i].Width), float64(g.uiButtons[i].Height),
			buttonColor)

		// 绘制按钮边框
		borderColor := color.RGBA{100, 100, 100, 255}
		if g.uiButtons[i].Enabled {
			borderColor = color.RGBA{150, 150, 150, 255}
		}
		if g.uiButtons[i].Hovered {
			borderColor = color.RGBA{200, 200, 200, 255}
		}

		vector.StrokeRect(screen,
			float32(g.uiButtons[i].X), float32(g.uiButtons[i].Y),
			float32(g.uiButtons[i].Width), float32(g.uiButtons[i].Height),
			1, borderColor, false)

		// 绘制按钮文本
		textColor := color.RGBA{150, 150, 150, 255}
		if g.uiButtons[i].Enabled {
			textColor = color.RGBA{255, 255, 255, 255}
		}

		textWidth := len(g.uiButtons[i].Text) * 7 // 估算文本宽度
		textX := g.uiButtons[i].X + (g.uiButtons[i].Width-textWidth)/2
		textY := g.uiButtons[i].Y + g.uiButtons[i].Height/2 + 5

		text.Draw(screen, g.uiButtons[i].Text, g.gameFont, textX, textY, textColor)
	}
}

// 更新UI按钮状态
func (g *Game) updateUIButtons() {
	// 根据选中的单位类型启用/禁用相应的按钮
	hasSelectedUnits := len(g.selectedUnits) > 0

	for i := range g.uiButtons {
		// 默认禁用所有按钮
		g.uiButtons[i].Enabled = false

		// 如果有选中的单位，启用相应的按钮
		if hasSelectedUnits {
			switch g.uiButtons[i].Action {
			case "attack":
				// 攻击按钮：只有当选中的单位有攻击能力时才启用
				for _, unit := range g.selectedUnits {
					if unit.CanAttack() {
						g.uiButtons[i].Enabled = true
						break
					}
				}
			case "move":
				// 移动按钮：总是启用
				g.uiButtons[i].Enabled = true
			case "stop":
				// 停止按钮：总是启用
				g.uiButtons[i].Enabled = true
			case "patrol":
				// 巡逻按钮：总是启用
				g.uiButtons[i].Enabled = true
			case "load":
				// 装载按钮：只有当选中的单位可以装载其他单位时才启用
				for _, unit := range g.selectedUnits {
					if unit.MaxPassengers > 0 && len(unit.Passengers) < unit.MaxPassengers {
						g.uiButtons[i].Enabled = true
						break
					}
				}
			case "unload":
				// 卸载按钮：只有当选中的单位有乘客时才启用
				for _, unit := range g.selectedUnits {
					if len(unit.Passengers) > 0 {
						g.uiButtons[i].Enabled = true
						break
					}
				}
			}
		}

		// 更新按钮的悬停状态
		mx, my := ebiten.CursorPosition()
		g.uiButtons[i].Hovered = mx >= g.uiButtons[i].X && mx < g.uiButtons[i].X+g.uiButtons[i].Width &&
			my >= g.uiButtons[i].Y && my < g.uiButtons[i].Y+g.uiButtons[i].Height
	}
}

// 处理UI按钮点击
func (g *Game) handleUIButtonClick(x, y int) {
	for _, button := range g.uiButtons {
		if button.Enabled && x >= button.X && x < button.X+button.Width && y >= button.Y && y < button.Y+button.Height {
			// 执行按钮对应的操作
			switch button.Action {
			case "attack":
				// 进入攻击模式，等待玩家选择攻击目标
				fmt.Println("进入攻击模式")
				// 这里可以设置一个标志，表示当前处于攻击模式
			case "move":
				// 进入移动模式，等待玩家选择移动目标
				fmt.Println("进入移动模式")
			case "stop":
				// 停止所有选中的单位
				for _, unit := range g.selectedUnits {
					unit.HasTarget = false
					unit.TargetUnit = nil
					unit.Path = nil
				}
				fmt.Println("单位已停止")
			case "patrol":
				// 进入巡逻模式，等待玩家选择巡逻路径
				fmt.Println("进入巡逻模式")
			case "load":
				// 进入装载模式，等待玩家选择要装载的单位
				fmt.Println("进入装载模式")
			case "unload":
				// 卸载所有选中单位的乘客
				for _, unit := range g.selectedUnits {
					for len(unit.Passengers) > 0 {
						passenger := unit.Passengers[0]
						unit.RemovePassenger(passenger)
						// 将乘客放置在载具附近的空位置
						passenger.X = unit.X
						passenger.Y = unit.Y
						g.findNearbyEmptyPosition(passenger)
						passenger.CurrentX = float64(passenger.X * TileSize)
						passenger.CurrentY = float64(passenger.Y * TileSize)
					}
				}
				fmt.Println("已卸载所有乘客")
			}
			break
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	// 返回实际游戏窗口大小，而不是整个地图大小
	return outsideWidth, outsideHeight
}

func (g *Game) getUnitAt(x, y int) *Unit {
	// 将网格坐标转换为像素坐标（中心点）
	clickX := float64(x*TileSize) + float64(TileSize)/2
	clickY := float64(y*TileSize) + float64(TileSize)/2

	// 按照距离排序，优先选择最近的单位
	var closestUnit *Unit
	minDistance := math.MaxFloat64

	for _, unit := range g.units {
		if unit.IsPassenger {
			continue
		}

		// 计算单位中心点到点击位置的距离
		dx := unit.CurrentX - clickX
		dy := unit.CurrentY - clickY
		distance := math.Sqrt(dx*dx + dy*dy)

		// 获取单位的碰撞半径，并增加一点额外的点击容差
		unitRadius := getCollisionRadius(unit) * 1.2 // 增加20%的点击容差

		// 如果距离小于单位的碰撞半径，且是最近的单位，则记录该单位
		if distance < unitRadius && distance < minDistance {
			closestUnit = unit
			minDistance = distance
		}
	}

	return closestUnit
}

func (g *Game) removeUnit(unit *Unit) {
	// 如果是乘客，先从载具中移除
	if unit.IsPassenger && unit.ParentUnit != nil {
		unit.ParentUnit.RemovePassenger(unit)
	}

	// 如果是载具，释放所有乘客
	if unit.Type == Armor && len(unit.Passengers) > 0 {
		// 复制乘客列表，因为在循环中会修改原列表
		passengers := make([]*Unit, len(unit.Passengers))
		copy(passengers, unit.Passengers)

		for _, passenger := range passengers {
			unit.RemovePassenger(passenger)
			// 为乘客找一个附近的空位置
			g.findNearbyEmptyPosition(passenger)
		}
	}

	// 清除所有指向该单位的目标引用
	for _, otherUnit := range g.units {
		if otherUnit.TargetUnit == unit {
			otherUnit.TargetUnit = nil
			otherUnit.HasTarget = false
			otherUnit.Path = make([][2]int, 0)
		}
	}

	// 从选中单位列表中移除
	for i, u := range g.selectedUnits {
		if u == unit {
			g.selectedUnits = append(g.selectedUnits[:i], g.selectedUnits[i+1:]...)
			break
		}
	}

	// 从单位列表中移除
	for i, u := range g.units {
		if u == unit {
			g.units = append(g.units[:i], g.units[i+1:]...)
			break
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// updateFogOfWar 更新战争迷雾状态
func (g *Game) updateFogOfWar() {
	// 重置所有区域为有迷雾和不可见
	for y := 0; y < MapHeight; y++ {
		for x := 0; x < MapWidth; x++ {
			g.fogOfWar[y][x] = true
			g.visibleToPlayer[y][x] = false
		}
	}

	// 遍历所有玩家单位，更新可见区域
	for _, unit := range g.units {
		if unit.IsPlayerUnit && !unit.IsPassenger {
			// 使用单位的视野范围
			visionRange := unit.VisionRange

			// 更新单位周围的可见区域
			for y := unit.Y - visionRange; y <= unit.Y+visionRange; y++ {
				for x := unit.X - visionRange; x <= unit.X+visionRange; x++ {
					// 检查是否在地图范围内
					if x >= 0 && x < MapWidth && y >= 0 && y < MapHeight {
						// 计算曼哈顿距离
						dist := abs(x-unit.X) + abs(y-unit.Y)

						// 如果在视野范围内，移除迷雾并标记为可见
						if dist <= visionRange {
							g.fogOfWar[y][x] = false
							g.visibleToPlayer[y][x] = true
						}
					}
				}
			}
		}
	}
}

// updateUnitAttacks 处理单位持续攻击敌人的逻辑
func (g *Game) updateUnitAttacks() {
	// 遍历所有单位（包括玩家单位和敌方单位）
	for _, unit := range g.units {
		// 跳过乘客单位
		if unit.IsPassenger {
			continue
		}

		// 如果单位有目标单位，执行持续攻击
		if unit.TargetUnit != nil {
			// 检查目标单位是否还存在且在可见区域内
			if unit.TargetUnit.Health > 0 &&
				((unit.IsPlayerUnit && g.visibleToPlayer[unit.TargetUnit.Y][unit.TargetUnit.X]) ||
					(!unit.IsPlayerUnit && g.visibleToPlayer[unit.Y][unit.X])) {
				// 检查目标是否为空中单位，如果是空中单位，则需要检查当前单位是否能攻击空中单位
				if unit.TargetUnit.IsAirUnit() && !unit.CanAttackAir {
					fmt.Printf("%s无法攻击空中单位%s\n", unit.Name, unit.TargetUnit.Name)
					continue
				}

				// 检查弹药量
				if unit.CurrentAmmo <= 0 {
					fmt.Printf("%s弹药耗尽，无法攻击！\n", unit.Name)
					continue
				}

				// 计算与目标单位的距离
				dist := abs(unit.X-unit.TargetUnit.X) + abs(unit.Y-unit.TargetUnit.Y)

				// 如果在攻击范围内且可以攻击
				if dist <= unit.AttackRange && unit.CanAttack() {
					// 计算命中率
					hitRate := unit.CalculateHitRate(dist)

					// 随机决定是否命中
					if rand.Float64() <= hitRate {
						// 检查攻击力是否大于防御力
						if unit.AttackPower > unit.TargetUnit.Defense {
							// 计算实际伤害
							damage := unit.AttackPower - unit.TargetUnit.Defense

							// 如果是防空单位攻击空中单位，增加额外伤害
							if unit.Type == AntiAir && unit.TargetUnit.IsAirUnit() {
								damage += 2 // 防空单位对空中单位额外伤害
								fmt.Printf("防空单位对空中单位造成额外伤害！\n")
							}

							// 如果是战斗机攻击空中单位，增加额外伤害
							if unit.Type == FighterJet && unit.TargetUnit.IsAirUnit() {
								damage += 1 // 战斗机对空中单位额外伤害
								fmt.Printf("战斗机对空中单位造成额外伤害！\n")
							}

							// 消耗弹药
							unit.CurrentAmmo--

							unit.TargetUnit.Health -= damage
							fmt.Printf("%s持续攻击%s成功！造成%d点伤害，命中率：%.2f，剩余弹药：%d\n",
								unit.Name, unit.TargetUnit.Name, damage, hitRate, unit.CurrentAmmo)

							// 更新上次攻击时间
							unit.UpdateLastAttackTime()

							// 设置攻击目标闪烁效果
							g.targetUnit = unit.TargetUnit
							g.targetFlashTime = time.Now().Add(500 * time.Millisecond) // 闪烁0.5秒

							// 如果敌方单位被消灭
							if unit.TargetUnit.Health <= 0 {
								fmt.Printf("%s被消灭！\n", unit.TargetUnit.Name)
								g.removeUnit(unit.TargetUnit)
								unit.TargetUnit = nil
							}
						} else {
							// 即使攻击未造成伤害，也消耗弹药
							unit.CurrentAmmo--

							fmt.Printf("%s持续攻击%s失败！攻击力(%d)不足以突破防御(%d)，命中率：%.2f，剩余弹药：%d\n",
								unit.Name, unit.TargetUnit.Name, unit.AttackPower, unit.TargetUnit.Defense, hitRate, unit.CurrentAmmo)

							// 更新上次攻击时间
							unit.UpdateLastAttackTime()
						}
					} else {
						// 即使未命中，也消耗弹药
						unit.CurrentAmmo--

						fmt.Printf("%s持续攻击%s未命中！命中率：%.2f，剩余弹药：%d\n",
							unit.Name, unit.TargetUnit.Name, hitRate, unit.CurrentAmmo)

						// 更新上次攻击时间
						unit.UpdateLastAttackTime()
					}
				} else if dist > unit.AttackRange {
					// 如果超出攻击范围，更新目标位置，追踪敌人
					targetPos := g.findNearbyEmptyPositionForTarget(unit.TargetUnit.X, unit.TargetUnit.Y, unit.Type)
					if targetPos[0] != -1 && targetPos[1] != -1 {
						unit.TargetX = targetPos[0]
						unit.TargetY = targetPos[1]
						unit.HasTarget = true
						// 使用A*寻路
						unit.Path = FindPathWithUnits(g.gameMap, g, unit.X, unit.Y, unit.TargetX, unit.TargetY, unit.Type)
					}
				}
			} else {
				// 目标单位不再可见或已被消灭，清除目标
				unit.TargetUnit = nil
			}
		} else {
			// 自动攻击：如果单位没有目标，寻找范围内的敌方单位进行攻击
			// 检查弹药量
			if unit.CurrentAmmo <= 0 {
				continue
			}

			// 检查是否可以攻击
			if !unit.CanAttack() {
				continue
			}

			// 寻找范围内的敌方单位
			var nearestEnemy *Unit = nil
			minDist := 9999

			for _, potentialTarget := range g.units {
				// 跳过同阵营单位和乘客单位
				if potentialTarget.IsPlayerUnit == unit.IsPlayerUnit || potentialTarget.IsPassenger {
					continue
				}

				// 检查目标是否在可见区域内
				if unit.IsPlayerUnit && !g.visibleToPlayer[potentialTarget.Y][potentialTarget.X] {
					continue
				}

				// 如果是敌方单位，检查玩家单位是否在其可见区域内
				if !unit.IsPlayerUnit && !g.visibleToPlayer[unit.Y][unit.X] {
					continue
				}

				// 计算距离
				dist := abs(unit.X-potentialTarget.X) + abs(unit.Y-potentialTarget.Y)

				// 检查是否在攻击范围内
				if dist <= unit.AttackRange {
					// 如果是空中单位，检查当前单位是否能攻击空中单位
					if potentialTarget.IsAirUnit() && !unit.CanAttackAir {
						continue
					}

					// 找到最近的敌人
					if dist < minDist {
						minDist = dist
						nearestEnemy = potentialTarget
					}
				}
			}

			// 如果找到了敌方单位，设置为攻击目标
			if nearestEnemy != nil {
				unit.TargetUnit = nearestEnemy
				if unit.IsPlayerUnit {
					fmt.Printf("%s自动锁定敌方目标%s！\n", unit.Name, nearestEnemy.Name)
				} else {
					fmt.Printf("敌方%s自动锁定玩家单位%s！\n", unit.Name, nearestEnemy.Name)
				}
			}
		}
	}
}

// sortSelectedUnitsByEntityID 按照实体ID对选中的单位进行排序
func (g *Game) sortSelectedUnitsByEntityID() {
	if len(g.selectedUnits) <= 1 {
		return // 如果只有0或1个单位，不需要排序
	}

	// 使用sort包对选中的单位进行排序
	sort.Slice(g.selectedUnits, func(i, j int) bool {
		return g.selectedUnits[i].EntityID < g.selectedUnits[j].EntityID
	})
}
