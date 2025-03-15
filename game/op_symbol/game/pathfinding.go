package game

import (
	"container/heap"
	"math"
)

// 定义节点结构体
type Node struct {
	X, Y   int     // 坐标
	G      float64 // 从起点到当前节点的成本
	H      float64 // 从当前节点到终点的估计成本
	F      float64 // G + H
	Parent *Node   // 父节点
	index  int     // 在优先队列中的索引
}

// 优先队列实现
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].F < pq[j].F
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	node := x.(*Node)
	node.index = n
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.index = -1
	*pq = old[0 : n-1]
	return node
}

// 计算两点之间的曼哈顿距离
func manhattanDistance(x1, y1, x2, y2 int) float64 {
	return float64(abs(x1-x2) + abs(y1-y2))
}

// A*寻路算法
func FindPath(gameMap *GameMap, startX, startY, endX, endY int, unitType UnitType) [][2]int {
	// 如果起点和终点相同，返回空路径
	if startX == endX && startY == endY {
		return [][2]int{}
	}

	// 检查终点是否可达
	if !isWalkable(gameMap, endX, endY, unitType) {
		// 如果终点不可达，尝试找到附近可达的点
		nearestX, nearestY := findNearestWalkable(gameMap, endX, endY, unitType)
		if nearestX == -1 && nearestY == -1 {
			return [][2]int{} // 没有找到可达点
		}
		endX, endY = nearestX, nearestY
	}

	// 定义方向数组（上、右、下、左、左上、右上、左下、右下）
	directions := [][2]int{
		{0, -1}, {1, 0}, {0, 1}, {-1, 0},
		{-1, -1}, {1, -1}, {-1, 1}, {1, 1},
	}

	// 创建开放列表和关闭列表
	openList := make(PriorityQueue, 0)
	closedSet := make(map[string]bool)
	nodeMap := make(map[string]*Node)

	// 创建起点节点
	startNode := &Node{
		X:      startX,
		Y:      startY,
		G:      0,
		H:      manhattanDistance(startX, startY, endX, endY),
		Parent: nil,
	}
	startNode.F = startNode.G + startNode.H

	// 将起点加入开放列表
	heap.Push(&openList, startNode)
	nodeMap[nodeKey(startX, startY)] = startNode

	// 开始寻路
	for openList.Len() > 0 {
		// 取出F值最小的节点
		current := heap.Pop(&openList).(*Node)
		key := nodeKey(current.X, current.Y)

		// 如果到达终点，构建路径并返回
		if current.X == endX && current.Y == endY {
			return buildPath(current)
		}

		// 将当前节点加入关闭列表
		closedSet[key] = true

		// 检查相邻节点
		for _, dir := range directions {
			newX := current.X + dir[0]
			newY := current.Y + dir[1]
			newKey := nodeKey(newX, newY)

			// 检查是否已在关闭列表中
			if closedSet[newKey] {
				continue
			}

			// 检查是否可行走
			if !isWalkable(gameMap, newX, newY, unitType) {
				continue
			}

			// 计算新的G值（对角线移动成本为1.414，直线移动成本为1）
			moveCost := 1.0
			if abs(dir[0]) == 1 && abs(dir[1]) == 1 {
				moveCost = 1.414
			}

			// 地形影响移动成本
			if gameMap.Terrain[newY][newX] == Mountain {
				moveCost *= 2.0 // 山地移动成本加倍
			} else if gameMap.Terrain[newY][newX] == Forest {
				moveCost *= 1.5 // 森林移动成本增加50%
			} else if gameMap.Terrain[newY][newX] == Water && unitType != Armor {
				continue // 非装甲单位不能进入水域
			}

			newG := current.G + moveCost

			// 检查节点是否已在开放列表中
			neighbor, exists := nodeMap[newKey]
			if !exists {
				// 创建新节点
				neighbor = &Node{
					X:      newX,
					Y:      newY,
					G:      newG,
					H:      manhattanDistance(newX, newY, endX, endY),
					Parent: current,
				}
				neighbor.F = neighbor.G + neighbor.H
				nodeMap[newKey] = neighbor
				heap.Push(&openList, neighbor)
			} else if newG < neighbor.G {
				// 更新节点
				neighbor.G = newG
				neighbor.F = neighbor.G + neighbor.H
				neighbor.Parent = current
				heap.Fix(&openList, neighbor.index)
			}
		}
	}

	// 没有找到路径
	return [][2]int{}
}

// A*寻路算法（考虑其他单位）
func FindPathWithUnits(gameMap *GameMap, game *Game, startX, startY, endX, endY int, unitType UnitType) [][2]int {
	// 如果起点和终点相同，返回空路径
	if startX == endX && startY == endY {
		return [][2]int{}
	}

	// 检查终点是否可达
	if !isWalkableWithUnits(gameMap, game, endX, endY, unitType) {
		// 如果终点不可达，尝试找到附近可达的点
		nearestX, nearestY := findNearestWalkableWithUnits(gameMap, game, endX, endY, unitType)
		if nearestX == -1 && nearestY == -1 {
			return [][2]int{} // 没有找到可达点
		}
		endX, endY = nearestX, nearestY
	}

	// 定义方向数组（上、右、下、左、左上、右上、左下、右下）
	directions := [][2]int{
		{0, -1}, {1, 0}, {0, 1}, {-1, 0},
		{-1, -1}, {1, -1}, {-1, 1}, {1, 1},
	}

	// 创建开放列表和关闭列表
	openList := make(PriorityQueue, 0)
	closedSet := make(map[string]bool)
	nodeMap := make(map[string]*Node)

	// 创建起点节点
	startNode := &Node{
		X:      startX,
		Y:      startY,
		G:      0,
		H:      manhattanDistance(startX, startY, endX, endY),
		Parent: nil,
	}
	startNode.F = startNode.G + startNode.H

	// 将起点加入开放列表
	heap.Push(&openList, startNode)
	nodeMap[nodeKey(startX, startY)] = startNode

	// 开始寻路
	for openList.Len() > 0 {
		// 取出F值最小的节点
		current := heap.Pop(&openList).(*Node)
		key := nodeKey(current.X, current.Y)

		// 如果到达终点，构建路径并返回
		if current.X == endX && current.Y == endY {
			return buildPath(current)
		}

		// 将当前节点加入关闭列表
		closedSet[key] = true

		// 检查相邻节点
		for _, dir := range directions {
			newX := current.X + dir[0]
			newY := current.Y + dir[1]
			newKey := nodeKey(newX, newY)

			// 检查是否已在关闭列表中
			if closedSet[newKey] {
				continue
			}

			// 检查是否可行走（考虑其他单位）
			if !isWalkableWithUnits(gameMap, game, newX, newY, unitType) {
				continue
			}

			// 计算新的G值（对角线移动成本为1.414，直线移动成本为1）
			moveCost := 1.0
			if abs(dir[0]) == 1 && abs(dir[1]) == 1 {
				moveCost = 1.414
			}

			// 地形影响移动成本
			if gameMap.Terrain[newY][newX] == Mountain {
				moveCost *= 2.0 // 山地移动成本加倍
			} else if gameMap.Terrain[newY][newX] == Forest {
				moveCost *= 1.5 // 森林移动成本增加50%
			} else if gameMap.Terrain[newY][newX] == Water && unitType != Armor {
				continue // 非装甲单位不能进入水域
			}

			newG := current.G + moveCost

			// 检查节点是否已在开放列表中
			neighbor, exists := nodeMap[newKey]
			if !exists {
				// 创建新节点
				neighbor = &Node{
					X:      newX,
					Y:      newY,
					G:      newG,
					H:      manhattanDistance(newX, newY, endX, endY),
					Parent: current,
				}
				neighbor.F = neighbor.G + neighbor.H
				nodeMap[newKey] = neighbor
				heap.Push(&openList, neighbor)
			} else if newG < neighbor.G {
				// 更新节点
				neighbor.G = newG
				neighbor.F = neighbor.G + neighbor.H
				neighbor.Parent = current
				heap.Fix(&openList, neighbor.index)
			}
		}
	}

	// 没有找到路径
	return [][2]int{}
}

// 检查位置是否可行走
func isWalkable(gameMap *GameMap, x, y int, unitType UnitType) bool {
	// 检查是否在地图范围内
	if x < 0 || x >= gameMap.Width || y < 0 || y >= gameMap.Height {
		return false
	}

	// 检查地形是否可行走
	terrainType := gameMap.Terrain[y][x]

	// 使用CanMoveOnTerrain函数检查单位是否可以在该地形上行动
	canMove, err := CanMoveOnTerrain(unitType, terrainType)
	if err != nil || !canMove {
		return false
	}

	return true
}

// 检查位置是否可行走（考虑其他单位）
func isWalkableWithUnits(gameMap *GameMap, game *Game, x, y int, unitType UnitType) bool {
	// 首先检查基本的地形可行走性
	if !isWalkable(gameMap, x, y, unitType) {
		return false
	}

	// 检查该位置是否有其他单位
	for _, unit := range game.units {
		if unit.X == x && unit.Y == y && !unit.IsPassenger {
			// 判断单位类型是否会造成障碍
			// 空中单位只考虑空中障碍物
			if unitType == Helicopter || unitType == FighterJet || unitType == Bomber {
				// 如果当前单位是空中单位，只有其他空中单位会阻挡它
				if unit.IsAirUnit() {
					return false
				}
			} else if gameMap.Terrain[y][x] == Water {
				// 水中单位只考虑水中障碍物
				// 假设只有装甲车可以在水中移动
				if unitType == Armor {
					// 如果当前单位是水中单位，只有其他水中单位会阻挡它
					if unit.Type == Armor {
						return false
					}
				}
			} else {
				// 地面单位只考虑地面障碍物
				if !unit.IsAirUnit() {
					return false
				}
			}
		}
	}

	return true
}

// 寻找最近的可行走位置
func findNearestWalkable(gameMap *GameMap, x, y int, unitType UnitType) (int, int) {
	// 搜索范围
	searchRadius := 5
	bestDistance := math.MaxFloat64
	bestX, bestY := -1, -1

	for dy := -searchRadius; dy <= searchRadius; dy++ {
		for dx := -searchRadius; dx <= searchRadius; dx++ {
			newX, newY := x+dx, y+dy
			if isWalkable(gameMap, newX, newY, unitType) {
				distance := math.Sqrt(float64(dx*dx + dy*dy))
				if distance < bestDistance {
					bestDistance = distance
					bestX, bestY = newX, newY
				}
			}
		}
	}

	return bestX, bestY
}

// 寻找最近的可行走位置（考虑其他单位）
func findNearestWalkableWithUnits(gameMap *GameMap, game *Game, x, y int, unitType UnitType) (int, int) {
	// 搜索范围
	searchRadius := 5
	bestDistance := math.MaxFloat64
	bestX, bestY := -1, -1

	for dy := -searchRadius; dy <= searchRadius; dy++ {
		for dx := -searchRadius; dx <= searchRadius; dx++ {
			newX, newY := x+dx, y+dy
			if isWalkableWithUnits(gameMap, game, newX, newY, unitType) {
				distance := math.Sqrt(float64(dx*dx + dy*dy))
				if distance < bestDistance {
					bestDistance = distance
					bestX, bestY = newX, newY
				}
			}
		}
	}

	return bestX, bestY
}

// 构建路径
func buildPath(endNode *Node) [][2]int {
	path := make([][2]int, 0)
	current := endNode

	// 从终点回溯到起点
	for current != nil {
		path = append([][2]int{{current.X, current.Y}}, path...)
		current = current.Parent
	}

	// 移除起点（因为单位已经在起点）
	if len(path) > 0 {
		path = path[1:]
	}

	return path
}

// 生成节点的唯一键
func nodeKey(x, y int) string {
	return string(rune(x)) + "," + string(rune(y))
}
