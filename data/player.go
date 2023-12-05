package data

import "singo/repo"

type PlayerLoginResp struct {
	// 登录token
	Token string `json:"token"`
	// 是否开通钱包
	Purse bool `json:"purse"`
}

type StartGameResp struct {
	// 是否准许进入游戏
	Permit     bool             `json:"permit"`
	GameRecord *repo.GameRecord `json:"game_record,omitempty"`
}

type EnergyResp struct {
	// 融合能量
	FusionEnergy int `json:"rh_energy"`
	// 自然能量
	NaturalEnergy int `json:"zr_energy"`
}

type ShareResp struct {
	// 是否有效分享
	Effective bool `json:"effective"`
}

type GetTeamResp struct {
	LastResult int                `gorm:"-" json:"last_result"` // 上一把对局结果
	Data       []*GetTeamInfoResp `json:"data"`
}

type GetTeamInfoResp struct {
	RecordId    uint64 // 对局ID
	Trophies    int    // 此时的奖杯数 用来给同奖杯的人匹配
	Round       int    // 回合数
	Position    int    // 位序
	CharacterID int    // 角色编号
	LevelNum    int    // 角色等级
	BaseAtk     int    // 攻击力，默认为0
	TempAtk     int    // 临时攻击力，默认为0
	PerAtk      int    // 永久攻击力，默认为0
	BaseHp      int    // 生命力，默认为0
	TempHp      int    // 临时生命力，默认为0
	PerHp       int    // 永久生命力，默认为0
	ExSkill     int    `json:"exSkill"` // 全局技能
}

type GameResultResp struct {
	Health      uint `json:"health"`       // 体力
	Trophies    int  ` json:"trophies"`    // 奖杯数
	LastBalance int  `json:"last_balance"` // 金币留存
	Round       int  `json:"round"`        // 回合数
}

type TeamInfoResp struct {
	Position int `json:"position"` // 位序 最多5个
	LocalID  int `json:"localid"`  // 角色code
	LevelNum int `json:"levelNum"` // 角色等级
	BaseAtk  int `json:"baseAtk"`  // 攻击力
	BaseHp   int `json:"baseHp"`   // 生命力
	TempAtk  int `json:"tempAtk"`  // 临时攻击力，默认为0
	PerAtk   int `json:"perAtk"`   // 永久攻击力，默认为0
	TempHp   int `json:"tempHp"`   // 临时生命力，默认为0
	PerHp    int `json:"perHp"`    // 永久生命力，默认为0
	ExSkill  int `json:"exSkill"`  // 全局技能
}

type Player struct {
	Nickname     string          `json:"nickname"`      // 对战的昵称
	Trophies     int             `json:"trophies"`      // 奖杯数
	Health       uint            `json:"health"`        // 体力数
	Round        int             `json:"round"`         // 回合数
	BackgroundId int             `json:"background_id"` // 背景id
	StandId      int             `json:"stand_id"`      // 站台id
	ExpressionId int             `json:"expression_id"` // 表情id
	AppearanceId int             `json:"appearance_id"` // 登场id
	Team         []*TeamInfoResp `json:"team"`
}

type ResultResp struct {
	Trophies    int            `json:"trophies"`               // 奖杯数
	Health      uint           `json:"health"`                 // 体力数
	Round       int            `json:"round"`                  // 回合数
	Result      bool           `json:"result"`                 // 对局是否结束
	Reward      *WinRewardResp `json:"reward,omitempty"`       // 奖励道具
	Repeat      bool           `json:"repeat"`                 // 是否重复
	BalanceData uint           `json:"balance_data,omitempty"` // 金额信息
	RewardID    uint           `json:"reward_id"`              // 奖励的宝箱ID
	Revive      bool           `json:"revive"`                 // 是否复活
}

type WinRewardResp struct {
	ItemID      uint   `json:"item_id"`      // 道具ID
	ItemType    string `json:"item_type"`    // 道具类型
	Name        string `json:"name"`         // 道具名称，非空字段
	Description string `json:"description"`  // 道具描述
	ResourceURL string `json:"resource_url"` // 资源路径
}

type GameInfoResp struct {
	OpenID       string `json:"open_id"`   // 玩家的OpenID，唯一
	Nickname     string `json:"nickname"`  // 玩家昵称
	AvatarId     int    `json:"avatar_id"` // 玩家头像编号
	LieId        int    `json:"lie_id"`    // 立绘编号
	BackgroundId int    `json:"background_id"`
	StandId      int    `json:"stand_id"`
	ExpressionId int    `json:"expression_id"`
	AppearanceId int    `json:"appearance_id"`
	// Health           int    `json:"health"`             // 体力数，默认为0
	Energy             int    `json:"energy"`               // 能量，默认为0
	Trophies           int    `json:"trophies"`             // 奖杯数，默认为0
	Balance            int    `json:"balance"`              // 余额，默认为0
	GamesPlayed        int    `json:"games_played"`         // 对局数，默认为0
	Wins               int    `json:"wins"`                 // 胜利场次，默认为0
	ChestsOpened       int    `json:"chests_opened"`        // 金宝箱场次，默认为0
	HighestWinStreak   int    `json:"highest_win_streak"`   // 最高连胜次数，默认为0
	Ranking            int64  `json:"ranking"`              // 排名
	LastGame           uint64 `json:"last_game,string"`     // 上一把游戏id
	LastEnergyRecovery int64  `json:"last_energy_recovery"` // 上次能量恢复时间戳，默认为0

}
