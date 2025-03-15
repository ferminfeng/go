package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// UnitTypeData 定义单位类型的所有属性
type UnitTypeData struct {
	TypeID          UnitType      `json:"type_id"`          // 单位类型ID
	Name            string        `json:"name"`             // 单位名称
	Health          int           `json:"health"`           // 生命值
	MoveSpeed       float64       `json:"move_speed"`       // 移动速度
	AllowedTerrains []TerrainType `json:"allowed_terrains"` // 允许行动的地形类型
	AttackPower     int           `json:"attack_power"`     // 攻击力
	AttackRange     int           `json:"attack_range"`     // 攻击范围
	AttackFrequency float64       `json:"attack_frequency"` // 攻击频率（每秒攻击次数）
	Defense         int           `json:"defense"`          // 防御力
	MaxPassengers   int           `json:"max_passengers"`   // 最大搭载数量
	Size            float64       `json:"size"`             // 单位大小（相对于TileSize的比例）
	MaxHitRate      float64       `json:"max_hit_rate"`     // 最大命中率（0.0-1.0）
	VisionRange     int           `json:"vision_range"`     // 视野范围（格子数）
	CanAttackAir    bool          `json:"can_attack_air"`   // 是否能够攻击空中单位
	Ammo            int           `json:"ammo"`             // 最大弹药量
	Description     string        `json:"description"`      // 单位描述
}

// UnitTypesConfig 存储所有单位类型的配置
type UnitTypesConfig struct {
	UnitTypes map[UnitType]UnitTypeData `json:"unit_types"`
}

// 全局单位类型配置
var GlobalUnitTypesConfig UnitTypesConfig

// LoadUnitTypesFromJSON 从JSON文件加载单位类型配置
func LoadUnitTypesFromJSON(filePath string) error {
	// 读取JSON文件
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("无法读取单位类型配置文件: %v", err)
	}

	// 解析JSON数据
	err = json.Unmarshal(jsonData, &GlobalUnitTypesConfig)
	if err != nil {
		return fmt.Errorf("解析单位类型配置失败: %v", err)
	}

	return nil
}

// GetUnitTypeData 获取指定单位类型的配置数据
func GetUnitTypeData(unitType UnitType) (UnitTypeData, error) {
	data, exists := GlobalUnitTypesConfig.UnitTypes[unitType]
	if !exists {
		return UnitTypeData{}, fmt.Errorf("未找到单位类型 %d 的配置数据", unitType)
	}
	return data, nil
}

// SaveDefaultUnitTypesConfig 保存默认的单位类型配置到JSON文件
func SaveDefaultUnitTypesConfig(filePath string) error {
	// 创建默认配置
	config := UnitTypesConfig{
		UnitTypes: map[UnitType]UnitTypeData{
			Infantry: {
				TypeID:          Infantry,
				Name:            "步兵",
				Health:          3,
				MoveSpeed:       1.0,
				AllowedTerrains: []TerrainType{Plain, Sand, Forest}, // 步兵可以在平原、沙地和森林上行动
				AttackPower:     1,
				AttackRange:     2,
				AttackFrequency: 1.2, // 步兵攻击频率较高
				Defense:         0,
				MaxPassengers:   0,
				Size:            0.8,
				MaxHitRate:      0.8,
				VisionRange:     3,     // 步兵视野范围
				CanAttackAir:    false, // 步兵不能攻击空中单位
				Ammo:            10,    // 步兵弹药量
				Description:     "基础步兵单位，可以被装甲单位搭载",
			},
			Armor: {
				TypeID:          Armor,
				Name:            "装甲车",
				Health:          5,
				MoveSpeed:       0.8,
				AllowedTerrains: []TerrainType{Plain, Sand}, // 装甲车只能在平原和沙地上行动
				AttackPower:     2,
				AttackRange:     2,
				AttackFrequency: 0.8, // 装甲车攻击频率较低
				Defense:         2,
				MaxPassengers:   2,
				Size:            1.2,
				MaxHitRate:      0.8,
				VisionRange:     4,     // 装甲车视野范围
				CanAttackAir:    false, // 装甲车不能攻击空中单位
				Ammo:            15,    // 装甲车弹药量
				Description:     "装甲单位，可以搭载步兵，具有较高防御力",
			},
			Artillery: {
				TypeID:          Artillery,
				Name:            "火炮",
				Health:          3,
				MoveSpeed:       0.5,
				AllowedTerrains: []TerrainType{Plain, Sand}, // 火炮只能在平原和沙地上行动
				AttackPower:     3,
				AttackRange:     5,
				AttackFrequency: 0.5, // 火炮攻击频率很低
				Defense:         0,
				MaxPassengers:   0,
				Size:            1.0,
				MaxHitRate:      0.5,
				VisionRange:     4,     // 火炮视野范围
				CanAttackAir:    false, // 火炮不能攻击空中单位
				Ammo:            8,     // 火炮弹药量
				Description:     "远程攻击单位，攻击范围大但移动慢",
			},
			Recon: {
				TypeID:          Recon,
				Name:            "侦察车",
				Health:          2,
				MoveSpeed:       1.5,
				AllowedTerrains: []TerrainType{Plain, Sand, Forest}, // 侦察车可以在平原、沙地和森林上行动
				AttackPower:     1,
				AttackRange:     2,
				AttackFrequency: 1.5, // 侦察车攻击频率较高
				Defense:         0,
				MaxPassengers:   0,
				Size:            0.9,
				MaxHitRate:      0.8,
				VisionRange:     6,     // 侦察车视野范围更大
				CanAttackAir:    false, // 侦察车不能攻击空中单位
				Ammo:            12,    // 侦察车弹药量
				Description:     "侦察单位，移动速度快",
			},
			AntiAir: {
				TypeID:          AntiAir,
				Name:            "防空炮",
				Health:          3,
				MoveSpeed:       1.2,
				AllowedTerrains: []TerrainType{Plain, Sand}, // 防空炮只能在平原和沙地上行动
				AttackPower:     2,
				AttackRange:     3,
				AttackFrequency: 1.2, // 防空炮攻击频率中等
				Defense:         1,
				MaxPassengers:   0,
				Size:            1.0,
				MaxHitRate:      0.8,
				VisionRange:     5,    // 防空炮视野范围
				CanAttackAir:    true, // 防空炮可以攻击空中单位
				Ammo:            20,   // 防空炮弹药量
				Description:     "防空单位，对空中单位有特殊攻击加成",
			},
			HeavyTank: {
				TypeID:          HeavyTank,
				Name:            "重型坦克",
				Health:          8,
				MoveSpeed:       0.6,
				AllowedTerrains: []TerrainType{Plain, Sand}, // 重型坦克只能在平原和沙地上行动
				AttackPower:     4,
				AttackRange:     2,
				AttackFrequency: 0.6, // 重型坦克攻击频率低
				Defense:         3,
				MaxPassengers:   0,
				Size:            1.4,
				MaxHitRate:      0.6,
				VisionRange:     3,     // 重型坦克视野范围
				CanAttackAir:    false, // 重型坦克不能攻击空中单位
				Ammo:            10,    // 重型坦克弹药量
				Description:     "重型坦克，具有极高的攻击力和防御力，但移动缓慢",
			},
			Helicopter: {
				TypeID:          Helicopter,
				Name:            "直升机",
				Health:          4,
				MoveSpeed:       2.0,
				AllowedTerrains: []TerrainType{}, // 直升机是空中单位，不受地形限制
				AttackPower:     2,
				AttackRange:     3,
				AttackFrequency: 1.8, // 直升机攻击频率高
				Defense:         1,
				MaxPassengers:   4,
				Size:            1.1,
				MaxHitRate:      0.8,
				VisionRange:     7,    // 直升机视野范围更大
				CanAttackAir:    true, // 直升机可以攻击空中单位
				Ammo:            16,   // 直升机弹药量
				Description:     "空中单位，可以飞越障碍物，能搭载步兵快速部署",
			},
			FighterJet: {
				TypeID:          FighterJet,
				Name:            "战斗机",
				Health:          3,
				MoveSpeed:       3.0,
				AllowedTerrains: []TerrainType{}, // 战斗机是空中单位，不受地形限制
				AttackPower:     3,
				AttackRange:     3,
				AttackFrequency: 2.0, // 战斗机攻击频率最高
				Defense:         1,
				MaxPassengers:   0,
				Size:            1.0,
				MaxHitRate:      0.8,
				VisionRange:     8,    // 战斗机视野范围最大
				CanAttackAir:    true, // 战斗机可以攻击空中单位
				Ammo:            6,    // 战斗机弹药量
				Description:     "高速空中单位，对其他空中单位有攻击加成",
			},
			Bomber: {
				TypeID:          Bomber,
				Name:            "轰炸机",
				Health:          5,
				MoveSpeed:       1.8,
				AllowedTerrains: []TerrainType{}, // 轰炸机是空中单位，不受地形限制
				AttackPower:     5,
				AttackRange:     3,
				AttackFrequency: 0.7, // 轰炸机攻击频率低
				Defense:         0,
				MaxPassengers:   0,
				Size:            1.3,
				MaxHitRate:      0.5,
				VisionRange:     6,     // 轰炸机视野范围
				CanAttackAir:    false, // 轰炸机不能攻击空中单位
				Ammo:            5,     // 轰炸机弹药量
				Description:     "空中轰炸单位，对地面单位造成大范围伤害",
			},
			RocketLauncher: {
				TypeID:          RocketLauncher,
				Name:            "火箭炮",
				Health:          3,
				MoveSpeed:       0.7,
				AllowedTerrains: []TerrainType{Plain, Sand}, // 火箭炮只能在平原和沙地上行动
				AttackPower:     4,
				AttackRange:     6,
				AttackFrequency: 0.4, // 火箭炮攻击频率最低
				Defense:         0,
				MaxPassengers:   0,
				Size:            1.1,
				MaxHitRate:      0.4,
				VisionRange:     5,    // 火箭炮视野范围
				CanAttackAir:    true, // 火箭炮可以攻击空中单位
				Ammo:            4,    // 火箭炮弹药量
				Description:     "远程火箭发射单位，攻击范围极大",
			},
			Engineer: {
				TypeID:          Engineer,
				Name:            "工程兵",
				Health:          2,
				MoveSpeed:       0.9,
				AllowedTerrains: []TerrainType{Plain, Sand, Forest}, // 工程兵可以在平原、沙地和森林上行动
				AttackPower:     1,
				AttackRange:     2,
				AttackFrequency: 1.0, // 工程兵攻击频率中等
				Defense:         0,
				MaxPassengers:   0,
				Size:            0.8,
				MaxHitRate:      0.8,
				VisionRange:     3,     // 工程兵视野范围
				CanAttackAir:    false, // 工程兵不能攻击空中单位
				Ammo:            5,     // 工程兵弹药量
				Description:     "工程兵单位，可以修复友方单位和建造防御工事",
			},
			MedicUnit: {
				TypeID:          MedicUnit,
				Name:            "医疗单位",
				Health:          3,
				MoveSpeed:       1.1,
				AllowedTerrains: []TerrainType{Plain, Sand, Forest}, // 医疗单位可以在平原、沙地和森林上行动
				AttackPower:     0,
				AttackRange:     0,
				AttackFrequency: 0.5, // 医疗单位攻击频率低（实际上是治疗频率）
				Defense:         0,
				MaxPassengers:   0,
				Size:            0.9,
				MaxHitRate:      0.5,
				VisionRange:     4,     // 医疗单位视野范围
				CanAttackAir:    false, // 医疗单位不能攻击空中单位
				Ammo:            0,     // 医疗单位没有弹药
				Description:     "医疗单位，可以治疗附近的友方单位",
			},
		},
	}

	// 将配置转换为JSON
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化单位类型配置失败: %v", err)
	}

	// 写入文件
	err = ioutil.WriteFile(filePath, jsonData, os.ModePerm)
	if err != nil {
		return fmt.Errorf("写入单位类型配置文件失败: %v", err)
	}

	return nil
}

// CanMoveOnTerrain 检查指定单位类型是否可以在特定地形上行动
func CanMoveOnTerrain(unitType UnitType, terrain TerrainType) (bool, error) {
	// 获取单位类型数据
	unitData, err := GetUnitTypeData(unitType)
	if err != nil {
		return false, err
	}

	// 空中单位不受地形限制
	if unitType == Helicopter || unitType == FighterJet || unitType == Bomber {
		return true, nil
	}

	// 检查地形是否在允许列表中
	for _, allowedTerrain := range unitData.AllowedTerrains {
		if allowedTerrain == terrain {
			return true, nil
		}
	}

	return false, nil
}

// InitUnitTypes 初始化单位类型配置
// 如果配置文件不存在，则创建默认配置
func InitUnitTypes(configPath string) error {
	// 检查配置文件是否存在
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		// 配置文件不存在，创建默认配置
		err = SaveDefaultUnitTypesConfig(configPath)
		if err != nil {
			return err
		}
	}

	// 加载配置
	return LoadUnitTypesFromJSON(configPath)
}
