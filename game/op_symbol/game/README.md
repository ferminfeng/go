# 单位类型配置系统

本文档说明如何使用JSON配置文件来定义游戏中不同类型单位的属性。

## 配置文件

单位类型配置文件位于 `op_symbol/game/unit_types.json`。这是一个JSON格式的文件，包含所有单位类型的属性定义。

## 配置文件结构

配置文件的基本结构如下：

```json
{
  "unit_types": {
    "0": {
      "type_id": 0,
      "name": "步兵",
      "health": 3,
      "move_speed": 1.0,
      "attack_power": 1,
      "attack_range": 1,
      "attack_frequency": 1.2,
      "defense": 0,
      "max_passengers": 0,
      "size": 0.8,
      "description": "基础步兵单位，可以被装甲单位搭载"
    },
    // 其他单位类型...
  }
}
```

## 单位类型属性

每个单位类型包含以下属性：

| 属性 | 类型 | 描述 |
|------|------|------|
| type_id | 整数 | 单位类型ID，对应代码中的UnitType枚举值 |
| name | 字符串 | 单位名称 |
| health | 整数 | 生命值 |
| move_speed | 浮点数 | 移动速度 |
| attack_power | 整数 | 攻击力 |
| attack_range | 整数 | 攻击范围（格子数） |
| attack_frequency | 浮点数 | 攻击频率（每秒攻击次数） |
| defense | 整数 | 防御力，攻击力必须大于防御力才能造成伤害 |
| max_passengers | 整数 | 最大搭载数量（仅对装甲单位和直升机有效） |
| size | 浮点数 | 单位大小（相对于TileSize的比例） |
| description | 字符串 | 单位描述 |

## 单位类型ID

目前支持的单位类型ID如下：

- 0: Infantry (步兵) - 基础步兵单位，可以被装甲单位和直升机搭载
- 1: Armor (装甲车) - 地面装甲单位，可以搭载步兵，具有较高防御力
- 2: Artillery (火炮) - 远程攻击单位，攻击范围大但移动慢
- 3: Recon (侦察车) - 侦察单位，移动速度快
- 4: AntiAir (防空炮) - 防空单位，对空中单位有特殊攻击加成
- 5: HeavyTank (重型坦克) - 重型坦克，具有极高的攻击力和防御力，但移动缓慢
- 6: Helicopter (直升机) - 空中单位，可以飞越障碍物，能搭载步兵快速部署
- 7: FighterJet (战斗机) - 高速空中单位，对其他空中单位有攻击加成
- 8: Bomber (轰炸机) - 空中轰炸单位，对地面单位造成大范围伤害
- 9: RocketLauncher (火箭炮) - 远程火箭发射单位，攻击范围极大
- 10: Engineer (工程兵) - 工程兵单位，可以修复友方单位和建造防御工事
- 11: MedicUnit (医疗单位) - 医疗单位，可以治疗附近的友方单位
- 12: Submarine (潜艇) - 水下单位，可以隐形并对水面单位发动突袭攻击

## 单位类型分类

单位可以按照以下几种类型进行分类：

### 步兵类单位
- Infantry (步兵)
- Engineer (工程兵)
- MedicUnit (医疗单位)

### 地面载具
- Armor (装甲车)
- Recon (侦察车)
- HeavyTank (重型坦克)

### 火炮类单位
- Artillery (火炮)
- AntiAir (防空炮)
- RocketLauncher (火箭炮)

### 飞行载具
- Helicopter (直升机)
- FighterJet (战斗机)
- Bomber (轰炸机)

## 如何修改配置

1. 直接编辑 `unit_types.json` 文件，修改现有单位类型的属性。
2. 如果配置文件不存在，游戏会自动创建一个包含默认配置的文件。
3. 修改配置文件后，重启游戏以应用新的配置。

## 防御系统说明

游戏实现了防御系统，具体规则如下：

1. 每个单位都有攻击力（attack_power）和防御力（defense）属性。
2. 当一个单位攻击另一个单位时，只有当攻击力大于防御力时才能造成伤害。
3. 实际造成的伤害为：攻击力 - 防御力。
4. 如果攻击力小于或等于防御力，则攻击无效，不会造成任何伤害。

这个系统使得高防御单位（如装甲单位）能够有效抵抗低攻击力单位的攻击，需要使用高攻击力单位（如炮兵）来对付它们。

## 攻击频率系统

游戏实现了攻击频率系统，具体规则如下：

1. 每个单位都有攻击频率（attack_frequency）属性，表示每秒可以攻击的次数。
2. 攻击频率越高，单位可以越频繁地发动攻击。
3. 攻击间隔（秒）= 1 / 攻击频率。例如，攻击频率为2.0的单位每0.5秒可以攻击一次。
4. 不同单位类型的攻击频率不同，反映了它们的战斗特性：
   - 步兵：攻击频率较高（1.2），但攻击力较低
   - 重型坦克：攻击频率较低（0.6），但攻击力很高
   - 火炮：攻击频率很低（0.4），但攻击范围大
   - 战斗机：攻击频率高（1.5），适合快速打击

这个系统使得游戏中的战斗更加动态和策略性，玩家需要考虑不同单位的攻击频率来制定战术。

## 命中率系统

游戏实现了基于距离的命中率系统，具体规则如下：

1. 每个单位都有最大命中率（max_hit_rate）属性，表示在最大攻击范围时的命中率。
2. 命中率随距离变化：
   - 距离为1时，命中率为100%，保证近距离攻击必定命中
   - 距离增加时，命中率线性降低
   - 在最大攻击范围时，命中率等于单位的最大命中率属性值
3. 不同单位类型的最大命中率不同：
   - 步兵、装甲车等近战单位：最大命中率为80%
   - 火炮、轰炸机等中程单位：最大命中率为50%
   - 火箭炮等远程单位：最大命中率为40%
4. 攻击时，系统会根据当前距离计算实际命中率，然后随机判断是否命中。

这个系统使得战斗更加真实和策略性，玩家需要权衡攻击距离和命中率之间的关系。近距离攻击虽然命中率高，但可能使单位暴露在敌方的攻击范围内；远距离攻击虽然安全，但命中率较低。

## 单位特性

### 搭载系统
- 装甲车可以搭载最多2名步兵
- 直升机可以搭载最多4名步兵，并能快速部署到地图各处

### 攻击特性
- 防空炮和战斗机对空中单位有额外伤害
- 轰炸机对地面单位有大范围伤害
- 火箭炮拥有最大的攻击范围
- 重型坦克拥有最高的单点攻击力

### 防御特性
- 重型坦克拥有最高的防御力(3)，可以完全抵消普通步兵的攻击
- 装甲车拥有较高的防御力(2)
- 大多数空中单位拥有基础防御力(1)

### 移动特性
- 战斗机拥有最快的移动速度(3.0)
- 直升机和侦察车也拥有较快的移动速度
- 重型坦克和火炮移动速度较慢

## 添加新的单位类型

要添加新的单位类型，需要：

1. 在 `unit.go` 文件中的 `UnitType` 枚举中添加新的单位类型。
2. 在 `unit_types.json` 文件中添加新单位类型的配置。
3. 在 `unit.go` 文件的 `Draw` 方法中添加新单位类型的绘制代码。
4. 在 `game.go` 文件的 `Draw` 方法中更新单位类型显示。

## 示例：添加新单位类型

假设我们要添加一个新的单位类型"潜艇"：

1. 在 `unit.go` 中添加新的枚举值：

```go
const (
    Infantry UnitType = iota
    Armor
    Artillery
    Recon
    AntiAir
    HeavyTank
    Helicopter
    FighterJet
    Bomber
    RocketLauncher
    Engineer
    MedicUnit
    Submarine // 新增的潜艇类型
)
```

2. 在 `unit_types.json` 中添加新的配置：

```json
"12": {
  "type_id": 12,
  "name": "潜艇",
  "health": 4,
  "move_speed": 1.0,
  "attack_power": 3,
  "attack_range": 2,
  "attack_frequency": 0.4,
  "defense": 2,
  "max_passengers": 0,
  "size": 1.1,
  "description": "水下单位，可以隐形并对水面单位发动突袭攻击"
}
```

3. 在 `unit.go` 的 `Draw` 方法中添加新单位类型的绘制代码。

## 注意事项

- 修改配置文件时，请确保JSON格式正确，否则游戏可能无法正确加载配置。
- 如果配置文件损坏，可以删除它，游戏会自动创建一个新的默认配置文件。
- 单位类型ID必须与代码中的枚举值一致，否则会导致错误。
- 设置防御力时要考虑游戏平衡性，防御力过高可能导致某些单位无法被攻击。
- 添加新单位类型时，需要考虑游戏平衡性，避免出现过于强大的单位。

## 迷雾战系统

游戏实现了迷雾战系统，具体规则如下：

1. 地图上的区域分为三种状态：
   - 完全迷雾（黑色）：从未探索过的区域
   - 已探索但当前不可见（灰色）：曾经探索过但当前不在视野范围内的区域
   - 可见区域（正常显示）：当前在友方单位视野范围内的区域

2. 每个单位都有视野范围（vision_range）属性，表示单位可以看到的格子数：
   - 侦察车：视野范围为6格，专门用于侦察
   - 战斗机：视野范围为8格，空中单位视野最广
   - 直升机：视野范围为7格，空中单位视野较广
   - 火炮和火箭炮：视野范围为5格，需要较好的观察位置
   - 大多数地面单位：视野范围为3-4格

3. 迷雾战对游戏的影响：
   - 玩家只能看到视野范围内的敌方单位
   - 玩家无法攻击迷雾中的敌方单位
   - 敌方单位只有在可见区域内才能被选中和攻击
   - 敌方AI只有在可见区域内才能攻击玩家单位

4. 战略意义：
   - 侦察变得非常重要，需要派出侦察单位探索地图
   - 高地形可能提供更好的视野（未来功能）
   - 需要保护自己的侦察单位，防止被敌方消灭
   - 可以利用迷雾隐藏自己的部队，进行战术性的伏击

这个系统大大增加了游戏的战略深度和不确定性，玩家需要谨慎行动，合理利用侦察单位，并随时准备应对从迷雾中出现的敌方单位。
