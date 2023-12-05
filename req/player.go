package req

type AddItemReq struct {
	// 道具ID
	ItemId int `json:"item_id"`
}

type ShareReq struct {
	// 游戏ID
	GameId int `json:"game_id"`
}

type TeamGetReq struct {
	GameId uint64 `json:"game_id,string" form:"game_id,string"` // 游戏id
}

type ResultGetReq struct {
	GameId uint64 `json:"game_id,string" form:"game_id,string"` // 游戏id
}

type CoinGetReq struct {
	GameId uint64 `json:"game_id,string" form:"game_id,string"` // 游戏id
}

type CoinSetReq struct {
	GameId uint64 `json:"game_id,string" form:"game_id,string"` // 游戏id
	Coin   uint64 `form:"coin"`
}

type PlayerTeamReq struct {
	GameId      uint64         `json:"game_id,string"` // 游戏id
	Team        []*TeamInfoReq `json:"team"`
	LastBalance int            `json:"last_balance"` // 上一把余额，默认为0
}

type TeamInfoReq struct {
	Position int `json:"position"` // 位序 最多5个
	LocalID  int `json:"localid"`  // 角色code
	LevelNum int `json:"levelNum"` // 角色等级
	BaseAtk  int `json:"baseAtk"`  // 攻击力
	TempAtk  int `json:"tempAtk"`  // 临时攻击力
	PerAtk   int `json:"perAtk"`   // 永久攻击力
	BaseHp   int `json:"baseHp"`   // 生命力
	TempHp   int `json:"tempHp"`   // 临时生命力
	PerHp    int `json:"perHp"`    // 永久生命力
	ExSkill  int `json:"exSkill"`  // 全局技能
}

type PlayerResultReq struct {
	GameId uint64 `json:"game_id,string" form:"game_id,string"` // 游戏id
	Result int    `json:"result" form:"result"`                 // 胜负平
}

type AvatarRoleListReq struct {
	*PageReq
}

type LaAdReq struct {
	ItemID   uint   `json:"item_id" url:"item_id"` // 道具ID
	AdId     string `json:"ad_id" url:"ad_id"`     // 广告id
	ItemType string `json:"item_type"`             // 道具类型
}

type PostAvatarReq struct {
	CharacterID int    `json:"character_id" url:"character_id"` // 角色编号
	AdId        string `json:"ad_id"`                           // 广告id
}

type PutNickNameReq struct {
	NickName string `json:"nickname" url:"nickname"`
}
