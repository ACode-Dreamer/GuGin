package repo

import (
	"gorm.io/gorm"
	"singo/req"
)

// @Description 单局游戏数据表 原先的逻辑-湛鑫加油
type GameRecord struct {
	ID               uint64 `gorm:"primary_key;auto_increment" json:"id,string"`  // 记录ID，主键自增
	OpenID           string `gorm:"index;not null" json:"open_id"`                // 玩家的OpenID，加索引，不能为空
	Health           uint   `gorm:"not null;default:2" json:"health"`             // 体力
	Trophies         int    `gorm:"not null;default:0" json:"trophies"`           // 奖杯数，默认为0
	LastBalance      int    `gorm:"not null;default:0" json:"last_balance"`       // 金币留存
	LastBalanceRound int    `gorm:"not null;default:0" json:"last_balance_round"` // 金币留存的回合
	Result           bool   `gorm:"not null;default:false" json:"-"`              // 对局结果
	Round            int    `gorm:"not null;default:0" json:"round"`              // 回合数
	Ad               bool   `gorm:"not null;default:0" json:"ad"`                 // 是否观看广告
	CreatedAt        int64  `json:"-"`                                            // 创建时间
	UpdatedAt        int64  `gorm:"autoUpdateTime" json:"-"`                      // 修改时间
}

func (rep *Repository) GetGameRecordByOpenId(openId string) (resp *GameRecord, err error) {
	err = rep.Where("open_id = ?", openId).Order("created_at DESC").First(&resp).Error
	return
}
func (rep *Repository) TotalChestsOpened(openId string) (total int64, err error) {
	err = rep.Model(&GameRecord{}).Where("open_id = ? and trophies = 9", openId).Count(&total).Error
	return
}
func (rep *Repository) UpdateGameRecord(record *GameRecord) (err error) {
	err = rep.Model(&GameRecord{}).Where("id = ?", record.ID).
		Updates(map[string]interface{}{
			"trophies": record.Trophies,
			"health":   record.Health,
			"round":    record.Round,
			"result":   record.Result,
		}).Error
	return err
}

func (rep *Repository) TotalGameRecord(openId string) (total int64, err error) {
	err = rep.Model(&GameRecord{}).Where("open_id = ?", openId).Count(&total).Error
	return
}

func (rep *Repository) GetGameRecord(gameId uint64) (resp *GameRecord, err error) {
	err = rep.Where("id = ?", gameId).First(&resp).Error
	return
}

func (rep *Repository) UpdateGameRecordBalance(record *GameRecord) (err error) {
	err = rep.Model(&GameRecord{}).Where("id = ?", record.ID).
		Updates(map[string]interface{}{
			"last_balance":       record.LastBalance,
			"last_balance_round": record.LastBalanceRound,
		}).Error
	return err
}

// PlayerInfo 玩家信息表
type PlayerInfo struct {
	ID               uint64 `gorm:"primary_key;auto_increment"`
	OpenID           string `gorm:"not null;unique"`          // 玩家的OpenID，唯一
	Nickname         string `gorm:"not null"`                 // 玩家昵称
	Avatar           string `gorm:"not null"`                 // 玩家头像
	AvatarId         int    `gorm:"not null"`                 // 玩家头像ID
	LieId            int    `gorm:"not null"`                 // 玩家立绘ID
	BackgroundId     int    `gorm:"not null"`                 // 背景ID
	StandId          int    `gorm:"not null"`                 // 站台ID
	ExpressionId     int    `gorm:"not null"`                 // 表情ID
	AppearanceId     int    `gorm:"not null"`                 // 登场ID
	Trophies         int    `gorm:"not null;default:0"`       // 奖杯数，默认为0
	Balance          int    `gorm:"not null;default:0"`       // 余额，默认为0
	Health           int    `gorm:"not null;default:0"`       // 体力数，默认为0
	GamesPlayed      int    `gorm:"not null;default:0"`       // 对局数，默认为0
	Wins             int    `gorm:"not null;default:0"`       // 胜利场次，默认为0
	ChestsOpened     int    `gorm:"not null;default:0"`       // 金宝箱场次，默认为0
	WinStreak        int    `gorm:"not null;default:0"`       // 连胜次数，默认为0
	Energy           int    `gorm:"not null;default:2;index"` // 能量，默认为0
	HighestWinStreak int    `gorm:"not null;default:0"`       // 最高连胜次数，默认为0

	LastEnergyRecovery int64 `gorm:"not null;default:0"`      // 上次能量恢复时间戳，默认为0
	CreatedAt          int64 `gorm:"autoCreateTime" json:"-"` // 创建时间
	UpdatedAt          int64 `gorm:"autoUpdateTime" json:"-"` // 修改时间
}

func (rep *Repository) UpdateNickname(openId, nickname string) (err error) {
	err = rep.Model(PlayerInfo{}).Where("open_id = ?", openId).Update("nickname", nickname).Error
	return
}
func (rep *Repository) BalanceDecrease(openId string, balance uint) (err error) {
	err = rep.Model(&PlayerInfo{}).Where("open_id = ?", openId).Update("balance", gorm.Expr("balance - ?", balance)).Error
	return
}
func (rep *Repository) GetRankingByOpenId(id uint64, trophies int) (ranking int64, err error) {
	var count int64
	if err = rep.Model(&PlayerInfo{}).
		Where("trophies > ?", trophies).
		Or("trophies = ? AND id > ?", trophies, id).Count(&count).Error; err != nil {
		return 0, err
	}
	return count + 1, nil
}

func (rep *Repository) ChestsOpenedSet(openId string, gamesPlayed int64) (err error) {
	err = rep.Model(&PlayerInfo{}).Where("open_id = ?", openId).Update("chests_opened", gamesPlayed).Error
	return
}
func (rep *Repository) BalanceIncrease(openId string, balance uint) (err error) {
	err = rep.Model(&PlayerInfo{}).Where("open_id = ?", openId).Update("balance", gorm.Expr("balance + ?", balance)).Error
	return
}
func (rep *Repository) WinsIncrease(openId string) (err error) {
	err = rep.Model(&PlayerInfo{}).Where("open_id = ?", openId).Update("wins", gorm.Expr("wins + 1")).Error
	return
}

func (rep *Repository) GetPlayerInfoByOpenId(openId string) (res *PlayerInfo, err error) {
	err = rep.Where("open_id = ?", openId).First(&res).Error
	return
}

func (rep *Repository) EnergyDecreaseWithEnergyTime(openId string, energy uint64, energyTime int64) (err error) {
	err = rep.Model(&PlayerInfo{}).Where("open_id = ? ", openId).Update("energy", gorm.Expr("energy - ?", energy)).Update("last_energy_recovery", energyTime).Error
	return
}

func (rep *Repository) EnergyDecrease(openId string, energy uint64) (err error) {
	err = rep.Model(&PlayerInfo{}).Where("open_id = ? ", openId).Update("energy", gorm.Expr("energy - ?", energy)).Error
	return
}

func (rep *Repository) GamesPlayedSet(openId string, gamesPlayed int64) (err error) {
	err = rep.Model(&PlayerInfo{}).Where("open_id = ?", openId).Update("games_played", gamesPlayed).Error
	return
}

func (rep *Repository) EnergyIncrease(openId string, energy uint64) (err error) {
	err = rep.Model(&PlayerInfo{}).Where("open_id = ? ", openId).Update("energy", gorm.Expr("energy + ?", energy)).Error
	return
}

func (rep *Repository) WinStreakUpdate(openId string, winStreak int) (err error) {
	err = rep.Model(&PlayerInfo{}).Where("open_id = ?", openId).Update("win_streak", winStreak).Error
	return
}

func (rep *Repository) HighestWinStreakUpdate(openId string, highestWinStreak int) (err error) {
	err = rep.Model(&PlayerInfo{}).Where("open_id = ?", openId).Update("highest_win_streak", highestWinStreak).Error
	return
}

func (rep *Repository) TrophiesIncrease(openId string, trophies int) (err error) {
	err = rep.Model(&PlayerInfo{}).Where("open_id = ?", openId).Update("trophies", gorm.Expr("trophies + ?", trophies)).Error
	return
}

// CoinRecord 金币数量
type CoinRecord struct {
	ID     uint64 `gorm:"primary_key;auto_increment"` // 记录ID，主键自增
	OpenID string `gorm:"index;not null"`             // 玩家的OpenID，加索引，不能为空
	GameId uint64 `gorm:"not null;"`                  // 游戏id
	Coin   uint64 `gorm:"not null;"`                  // 金币数量
}

func (rep *Repository) SaveCoinRecord(openID string, gameID, coin uint64) error {
	existingRecord, err := rep.GetCoinRecordByGameID(gameID)
	if err != nil && !IsNotFound(err) {
		return err
	}

	if existingRecord == nil {
		record := &CoinRecord{
			OpenID: openID,
			GameId: gameID,
			Coin:   coin,
		}
		err = rep.Create(record)
	} else {
		err = rep.Model(CoinRecord{}).Where("game_id = ?", gameID).Update("coin", coin).Error
	}
	return err
}

func (rep *Repository) GetCoinRecordByGameID(gameID uint64) (*CoinRecord, error) {
	var resp *CoinRecord
	err := rep.Where("game_id = ?", gameID).Last(&resp).Error
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Backpack 背包表
type Backpack struct {
	ID        uint64 `gorm:"primary_key;auto_increment"` // 背包ID，主键自增
	OpenID    string `gorm:"index;not null"`             // 玩家的OpenID，加索引，不能为空
	ItemID    uint   `gorm:"not null"`                   // 道具ID，不能为空
	ItemType  string `gorm:"not null"`                   // 道具类型，不能为空
	CreatedAt int64  // 创建时间
}

func (rep *Repository) GetBackpackByOpenIdAndItemId(openId string, itemId uint) (resp *Backpack, err error) {
	db := rep.Model(&Backpack{}).Where("open_id=? and item_id = ?", openId, itemId)
	err = db.First(&resp).Error
	return
}

// GameResultRecord 回合对局结果记录表
type GameResultRecord struct {
	ID        uint64 `gorm:"primary_key;auto_increment"`      // 记录ID，主键自增
	OpenID    string `gorm:"index;not null"`                  // 玩家的OpenID，加索引，不能为空
	GameId    uint64 `gorm:"not null;"`                       // 游戏id
	Result    int    `gorm:"not null;default:0"`              // 对局结果
	Round     int    `gorm:"not null;default:0" json:"round"` // 回合数
	CreatedAt int64  // 创建时间
}

func (rep *Repository) GetTeamInfoRecordByGameId(gameId uint64, openId string) (res *GameResultRecord, err error) {
	err = rep.Model(&GameResultRecord{}).
		Where("game_id = ? AND open_id = ?", gameId, openId).
		Order("round DESC").
		Limit(1).
		First(&res).Error
	return
}

func (rep *Repository) LastGameRecord(openId string, gameId uint64) (resp *GameResultRecord, err error) {
	err = rep.Where("open_id = ? and game_id = ?", openId, gameId).Order("created_at DESC").First(&resp).Error
	return
}

// TeamInfo 队伍配置表 - 每回合留存 供匹配机制
type TeamInfo struct {
	ID          uint64 `gorm:"primary_key;auto_increment"`
	OpenID      string `gorm:"not null;index"`                           // 玩家的OpenID
	RecordId    uint64 `gorm:"not null;index"`                           // 对局ID
	Trophies    int    `gorm:"not null;default:0;index" json:"trophies"` // 此时的奖杯数 用来给同奖杯的人匹配
	Round       int    `gorm:"not null;default:0" json:"round"`          // 回合数
	Position    int    `gorm:"not null"`                                 // 位序
	CharacterID int    `gorm:"not null"`                                 // 角色编号
	LevelNum    int    `gorm:"not null"`                                 // 角色等级
	BaseAtk     int    `gorm:"not null;default:0"`                       // 攻击力，默认为0
	TempAtk     int    `gorm:"not null;default:0"`                       // 临时攻击力，默认为0
	PerAtk      int    `gorm:"not null;default:0"`                       // 永久攻击力，默认为0
	BaseHp      int    `gorm:"not null;default:0"`                       // 生命力，默认为0
	TempHp      int    `gorm:"not null;default:0"`                       // 临时生命力，默认为0
	PerHp       int    `gorm:"not null;default:0"`                       // 永久生命力，默认为0
	ExSkill     int    `gorm:"not null;default:0" json:"exSkill"`        // 全局技能
	CreatedAt   int64  `gorm:"autoCreateTime"`                           // 创建时间
	UpdatedAt   int64  `gorm:"autoUpdateTime"`                           // 修改时间
}

func (rep *Repository) GetTeamInfoByGameId(gameId uint64, openId string) (array []*TeamInfo, err error) {
	var round int
	// // 根据上回合结果拿回合
	// rep.db.Model(&GameResultRecord{}).Select("round").
	// 	Where("game_id = ? AND open_id = ?", gameId, openId).
	// 	Group("round").
	// 	Order("round DESC").
	// 	Limit(1).
	// 	Pluck("round", &round1)

	// 根据保存的队伍回合拿回合
	rep.Model(&TeamInfo{}).Select("round").
		Where("record_id = ? AND open_id = ?", gameId, openId).
		Group("round").
		Order("round DESC").
		Limit(1).
		Pluck("round", &round)

	// if round1 > round2 {
	// 	round = round1
	// } else {
	// 	round = round2
	// }
	// 谁的回合大就取谁 可以同时解决
	// 1 空阵容打完回来会拉取到之前的非空阵容
	// 2 战备状态保存了信息退出去回来会保存到

	err = rep.Model(&TeamInfo{}).
		Where("round = ? AND open_id = ? and record_id = ?", round, openId, gameId).
		Find(&array).
		Error

	return
}

func (rep *Repository) DeleteTeamInfoByOpenIdAndRecordId(openId string, recordId uint64, round int) (err error) {
	err = rep.Where("open_id = ? and record_id = ? and round = ?", openId, recordId, round).Delete(&TeamInfo{}).Error
	return
}

func (rep *Repository) GetTeamInfosByRound(round int, openId string) (array []*TeamInfo, warId string, err error) {
	rep.Model(TeamInfo{}).Select("open_id").
		Where("round = ? AND open_id != ? AND position != 999", round, openId).
		Group("open_id").
		Having("COUNT(open_id) > ?", 1).
		Order("RAND()").
		Limit(1).
		Pluck("open_id", &warId)

	if warId != "" {
		recordId := ""
		// 获取随机选择的 record_id
		rep.Model(&TeamInfo{}).
			Select("record_id").
			Where("round = ? AND open_id = ?", round, warId).
			Group("record_id").
			Order("RAND()").
			Limit(1).
			Pluck("record_id", &recordId)
		if recordId != "" {
			err = rep.Model(&TeamInfo{}).
				Where("record_id =? and round = ? and open_id=?", recordId, round, warId).
				Find(&array).Error
		}
	}
	return
}

// Role 游戏角色表
type Role struct {
	ID            uint64 `gorm:"primary_key;auto_increment"`
	CharacterID   int    `gorm:"not null"`                                  // 角色编号
	Name          string `gorm:"not null"`                                  // 角色名字
	CreatedAt     int64  `gorm:"autoCreateTime" json:"-"`                   // 创建时间
	UpdatedAt     int64  `gorm:"autoUpdateTime" json:"-"`                   // 修改时间
	SkinPriceType uint   `gorm:"not null;default:0" json:"skin_price_type"` // 头像解锁方式
	SkinPrice     uint   `json:"skin_price"`                                // 头像价格
	LiePriceType  uint   `gorm:"not null;default:0" json:"lie_price_type"`  // 立绘解锁方式
	LiePrice      uint   `json:"lie_price"`                                 // 立绘价格
	SkinId        int    `gorm:"-" json:"skin_id,omitempty"`
	AvatarId      int    `gorm:"-" json:"avatar_id,omitempty"`
}

func (rep *Repository) GetRoleByCharacterID(characterId int) (resp *Role, err error) {
	err = rep.Where("character_id = ?", characterId).First(&resp).Error
	return
}

// PlayerRoleInfo 玩家解锁的头像、立绘记录
type PlayerRoleInfo struct {
	ID         uint64 `gorm:"primary_key;auto_increment"`
	OpenID     string `gorm:"not null"` // 玩家的OpenID
	UnlockType string `gorm:"not null"` // 解锁的类型
	UnlockId   int    `gorm:"not null"` // 解锁的编号
}

func (rep *Repository) GetPlayerRoleInfoByUnlockIdAndType(openId string, unlockId int, lockType string) (resp *PlayerRoleInfo, err error) {

	err = rep.Model(&PlayerRoleInfo{}).Where("open_id = ? and unlock_id = ? and unlock_type = ?", openId, unlockId, lockType).First(&resp).Error
	return
}

// ReviveRecord  复活记录表
type ReviveRecord struct {
	ID        uint64 `gorm:"primary_key;auto_increment"` // 记录ID，主键自增
	OpenID    string `gorm:"index;not null"`             // 玩家的OpenID，加索引，不能为空
	GameId    uint64 `gorm:"not null;"`                  // 游戏id
	CreatedAt int64  `json:"-"`                          // 创建时间
}

func (rep *Repository) LastReviveRecord(openId string, gameId uint64) (resp *ReviveRecord, err error) {
	err = rep.Where("open_id = ? and game_id = ?", openId, gameId).First(&resp).Error
	return
}

// RewardConfig 奖杯所属奖励配置表
type RewardConfig struct {
	TrophyCount uint // 奖杯数
	RewardID    uint // 奖励ID
}

func (rep *Repository) GetRewardConfigBy(trophies int) (resp *RewardConfig, err error) {
	err = rep.Where("trophy_count <= ?", trophies).Order("trophy_count DESC").First(&resp).Error
	return
}

// Reward 道具奖池及几率配置表
type Reward struct {
	RewardID  uint   // 奖励ID，作为主键
	Name      string `gorm:"not null"` // 奖励名称，非空字段
	Weight    int    // 权重
	RelatedID uint   // 关联ID
}

func (rep *Repository) GetRewardByRewardID(rewardID uint) (array []*Reward, err error) {
	err = rep.Where("reward_id = ?", rewardID).Order("weight DESC").Find(&array).Error
	return
}

// ItemReward 道具奖池中间表  多对多的关系
type ItemReward struct {
	ItemID    uint `json:"item_id"` // 道具ID，作为主键
	RelatedID uint `json:"-"`       // 关联ID
}

func (rep *Repository) GetRandomItemReward(relatedId uint) (res *ItemReward, err error) {
	err = rep.Where("related_id = ?", relatedId).Order("RAND()").First(&res).Error
	return
}

// Item 游戏道具表
type Item struct {
	ItemID      uint   `gorm:"primaryKey" json:"item_id"`            // 道具ID，作为主键
	ItemType    string `gorm:"not null" json:"item_type"`            // 道具类型，非空字段
	Name        string `gorm:"not null" json:"name"`                 // 道具名称，非空字段
	Description string `json:"description,omitempty"`                // 道具描述
	ResourceURL string `json:"resource_url,omitempty"`               // 资源路径
	PriceType   uint   `gorm:"not null;default:0" json:"price_type"` // 解锁方式
	Price       uint   `json:"price"`                                // 价格
}

func (rep *Repository) GetItemById(itemId uint) (resp *Item, err error) {
	err = rep.Where("item_id = ?", itemId).Find(&resp).Error
	return
}

type RoleFlag struct {
	*Role
	Flag       bool
	OwnerPrice uint `json:"owner_price"`
}

func (rep *Repository) GetAvatarPlayer(openId string, param *req.AvatarRoleListReq) (array []*RoleFlag, total int64, err error) {
	db := rep.Model(&Role{}).
		Select("role.*, CASE WHEN player_role_info.id IS NULL THEN 0 ELSE 1 END AS flag").
		Joins("LEFT JOIN player_role_info ON role.character_id = player_role_info.unlock_id  AND  player_role_info.unlock_type = 'avatar' AND player_role_info.open_id = ?", openId).
		Order(" flag  DESC,role.character_id ASC")

	if err = db.Count(&total).Error; err != nil {
		return
	}

	err = db.Limit(param.Limit).Offset(param.Offset()).Find(&array).Error
	return
}

// LAAdHistory 广告观看计次
type LAAdHistory struct {
	ID        uint64 `gorm:"primary_key;auto_increment"` // 背包ID，主键自增
	OpenID    string `gorm:"index;not null"`             // 玩家的OpenID，加索引，不能为空
	ItemID    uint   `gorm:"not null"`                   // 角色ID，不能为空
	ItemType  string `json:"item_type"`                  // 道具类型
	AdId      string `gorm:"not null"`                   // 广告id
	CreatedAt int64  // 创建时间
}

// LAAdCount 广告计次
type LAAdCount struct {
	ItemID     int  `gorm:"not null"` // 道具ID，不能为空
	OwnerPrice uint // 已观看次数
}

func (rep *Repository) GetLAAdCountHistory(openID, itemType string) ([]*LAAdCount, error) {
	var result []*LAAdCount
	err := rep.
		Model(&LAAdHistory{}).
		Select("item_id, count(0) as owner_price").
		Where("open_id = ? and item_type = ?", openID, itemType).
		Group("item_id").
		Find(&result).
		Error

	return result, err
}

// NickNameRecord 改名记录表
type NickNameRecord struct {
	ID        uint64 `gorm:"primary_key;auto_increment"` // 记录ID，主键自增
	OpenID    string `gorm:"index;not null"`             // 玩家的OpenID，加索引，不能为空
	Nickname  string `gorm:"not null"`                   // 玩家昵称
	CreatedAt int64  // 创建时间
}

func (rep *Repository) LastNickNameRecord(openId string) (resp *NickNameRecord, err error) {
	err = rep.Where("open_id = ?", openId).First(&resp).Error
	return
}
