package service

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"singo/conf"
	"singo/data"
	"singo/logger"
	"singo/middleware"
	"singo/repo"
	"singo/req"
	"singo/util"
	"time"
)

func PlayerLogin(code string) *data.Response {

	authResponse := util.AuthenticationCode(code)
	if !authResponse.Success {
		return data.NewErrorResponse(30002, authResponse.Message)
	}

	if authResponse.Data.OpenID == "" {
		return data.NewErrorResponse(30002, "openId为空")
	}

	// 用户是否存在
	_, err := rep().GetPlayerInfoByOpenId(authResponse.Data.OpenID)
	if err != nil && !repo.IsNotFound(err) {
		logger.Error("查询库内用户信息失败", err)
		return data.NewErrorResponse(30009, "查询用户失败")
	}
	if repo.IsNotFound(err) {
		// 创建新用户
		if err = makeNewPlayer(authResponse.Data.OpenID); err != nil {
			logger.Error("创建新用户出错", err)
			return data.NewErrorResponse(30010, "创建用户失败")
		}
	}

	// 校验钱包
	walletResponse := util.CheckFuiouWallet(authResponse.Data.OpenID)
	if !walletResponse.Success || !walletResponse.Data {
		return data.NewErrorResponse(30011, "未开通钱包")
	}

	// 颁布token
	expirationTime := time.Now().Add(12 * time.Hour)
	claims := &middleware.Claims{
		Username: authResponse.Data.OpenID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(conf.GetConfig().Server.Secret))
	if err != nil {
		logger.Error("颁发Token错误", err)
		return data.NewErrorResponse(30003, "登录失败")
	}

	return data.NewDataResponse(data.PlayerLoginResp{
		Token: tokenString,
		Purse: true,
	})
}

func makeNewPlayer(openId string) (err error) {

	// 获取当前时间戳
	currentTime := time.Now().Unix()

	// 拼接编号字符串
	nickName := fmt.Sprintf("%s%d", "无名鼠辈_", currentTime)

	logger.Info("创建新用户->", openId, "--", nickName)

	if err = rep().Create(&repo.PlayerInfo{
		OpenID:       openId,
		Nickname:     nickName,
		AvatarId:     0,
		LieId:        10000,
		BackgroundId: 22000,
		StandId:      20000,
		ExpressionId: 21000,
		AppearanceId: 0,
		Trophies:     0,
		Energy:       conf.GetConfig().Game.MaxEnergy,
	}); err != nil {
		logger.Error("创建新用户出错", err)
		return err
	}
	// 创建默认的四个道具
	if err = rep().Create(&repo.Backpack{
		OpenID:   openId,
		ItemID:   21000,
		ItemType: "Expression",
	}); err != nil {
		logger.Error("背包保存报错:", err)
		return err
	}
	if err = rep().Create(&repo.Backpack{
		OpenID:   openId,
		ItemID:   20000,
		ItemType: "Stand",
	}); err != nil {
		logger.Error("背包保存报错:", err)
		return err
	}
	if err = rep().Create(&repo.Backpack{
		OpenID:   openId,
		ItemID:   22000,
		ItemType: "Background",
	}); err != nil {
		logger.Error("背包保存报错:", err)
		return err
	}
	if err = rep().Create(&repo.Backpack{
		OpenID:   openId,
		ItemID:   0,
		ItemType: "Appearance",
	}); err != nil {
		logger.Error("背包保存报错:", err)
		return err
	}
	return
}

func StartGame(openId string) *data.Response {

	resp := new(data.StartGameResp)
	resp.Permit = false
	cardResponse := util.GetCardInfo(openId)

	if !cardResponse.Success {
		return data.NewErrorResponse(30004, cardResponse.Message)
	}
	for _, cardData := range cardResponse.Data {
		if cardData.MetaProductNo == conf.GetConfig().Ubq.PermitMetaProductNo {
			// 是否有卡
			if len(cardData.NfrInfoList) > 0 {
				resp.Permit = true
			}
			break
		}
	}

	if resp.Permit {

		// 体力是否充足
		myInfo, err := rep().GetPlayerInfoByOpenId(openId)
		if err != nil && !repo.IsNotFound(err) {
			logger.Error("查询个人信息失败", err)
			return data.NewErrorResponse(30051, "查询个人信息失败")
		}
		if myInfo.Energy < 1 {
			logger.Info("能量不足", myInfo.ID)
			return data.NewErrorResponse(30052, "能量不足")
		}
		if myInfo.Energy == conf.GetConfig().Game.MaxEnergy {
			if err = rep().EnergyDecreaseWithEnergyTime(openId, 1, time.Now().Unix()); err != nil {
				logger.Info("能量扣除报错", myInfo.ID)
				return data.NewErrorResponse(30053, "能量扣除失败")
			}
		} else {
			if err = rep().EnergyDecrease(openId, 1); err != nil {
				logger.Info("能量扣除报错", myInfo.ID)
				return data.NewErrorResponse(30054, "能量扣除失败")
			}
		}

		gameRecord := &repo.GameRecord{
			ID:       util.NextID(),
			OpenID:   openId,
			Health:   5,
			Trophies: 0,
			Result:   false,
		}
		if err = rep().Create(gameRecord); err != nil {
			logger.Error("新增游戏错误", err)
			return data.NewErrorResponse(30055, "开始游戏失败")
		}

		// 给这家伙场次
		total, err := rep().TotalGameRecord(openId)
		if err != nil {
			logger.Error("游戏错误:", err)
			return data.NewErrorResponse(30056, "游戏场次设置失败")
		}
		if err = rep().GamesPlayedSet(openId, total); err != nil {
			logger.Error("游戏计次错误:", err)
			return data.NewErrorResponse(30057, "游戏场次设置失败")
		}

		if err = rep().SaveCoinRecord(openId, gameRecord.ID, 10); err != nil {
			logger.Error("保存金币错误", err)
			return data.NewErrorResponse(30058, "保存金币失败")
		}
		resp.GameRecord = gameRecord
	}

	return data.NewDataResponse(resp)
}

func MyEnergy(openId string) *data.Response {

	resp := new(data.EnergyResp)
	userInfoResponse := util.GetUserInfo(openId)

	if !userInfoResponse.Success {
		return data.NewErrorResponse(30005, userInfoResponse.Message)
	}

	resp.NaturalEnergy = userInfoResponse.Data.NaturalEnergy
	resp.FusionEnergy = userInfoResponse.Data.FusionEnergy

	return data.NewDataResponse(resp)
}

func AddHealth(openId string) *data.Response {

	energy := redis().GetUbqHealth()
	deductionResponse := util.DeductEnergy(openId, energy)

	if !deductionResponse.Success {
		return data.NewErrorResponse(30006, deductionResponse.Message)
	}

	//  体力恢复
	myInfo, err := rep().GetPlayerInfoByOpenId(openId)
	if err != nil {
		logger.Error("查询个人信息报错:", err)
		return data.NewErrorResponse(30012, "查询个人信息失败")
	}
	if myInfo.Energy >= conf.GetConfig().Game.MaxEnergy {
		return data.NewErrorResponse(30013, "体力已达上限")
	}
	if err = rep().EnergyIncrease(openId, 1); err != nil {
		logger.Error("增加能量报错:", err)
		return data.NewErrorResponse(30014, "体力增加失败")
	}

	return data.NewSuccessResponse()
}

func AddItem(openId string, param *req.AddItemReq) *data.Response {

	// todo 根据道具id查询能量
	energy := 500

	deductionResponse := util.DeductEnergy(openId, energy)

	if !deductionResponse.Success {
		return data.NewErrorResponse(30007, deductionResponse.Message)
	}

	return data.NewSuccessResponse()
}

func Share(openId string, param *req.ShareReq) *data.Response {

	return data.NewSuccessResponse()
}

func GetTeam(openId string, param *req.TeamGetReq) *data.Response {
	resp := &data.GetTeamResp{
		Data: make([]*data.GetTeamInfoResp, 0),
	}
	res, err := rep().GetTeamInfoRecordByGameId(param.GameId, openId)
	if err != nil && !repo.IsNotFound(err) {
		logger.Error("查询上把对局信息报错", err)
		return data.NewErrorResponse(30015, "查询对局信息失败")
	}
	// if repo.IsNotFound(err) {
	// 	return resp, nil
	// }
	resp.LastResult = res.Result

	array, err := rep().GetTeamInfoByGameId(param.GameId, openId)
	if err != nil {
		logger.Error("查询队伍信息报错", err)
		return data.NewErrorResponse(30016, "查询队伍信息失败")
	}
	for _, info := range array {
		if info.Position == 999 {
			break
		}
		resp.Data = append(resp.Data, &data.GetTeamInfoResp{
			RecordId:    info.RecordId,
			Trophies:    info.Trophies,
			Round:       info.Round,
			Position:    info.Position,
			CharacterID: info.CharacterID,
			LevelNum:    info.LevelNum,
			BaseAtk:     info.BaseAtk,
			TempAtk:     info.TempAtk,
			PerAtk:      info.PerAtk,
			BaseHp:      info.BaseHp,
			TempHp:      info.TempHp,
			PerHp:       info.PerHp,
			ExSkill:     info.ExSkill,
		})
	}
	return data.NewDataResponse(resp)
}

func GetResult(openId string, param *req.ResultGetReq) *data.Response {

	// 对局校验
	record, err := rep().GetGameRecord(param.GameId)
	if err != nil || record.OpenID != openId {
		logger.Error("查询对局报错:", err)
		return data.NewErrorResponse(30017, "查询对局失败")
	}

	res, err := rep().GetTeamInfoRecordByGameId(param.GameId, openId)
	if err != nil && !repo.IsNotFound(err) {
		logger.Error("查询上把对局信息报错", err)
		return data.NewErrorResponse(30018, "查询对局失败")
	}
	// 无保存result结果不认可金币和队伍
	if res.ID == 0 || res.Round != record.LastBalanceRound {
		record.LastBalance = 0
	}
	resp := &data.GameResultResp{
		Health:      record.Health,
		Trophies:    record.Trophies,
		LastBalance: record.LastBalance,
		Round:       record.Round,
	}
	return data.NewDataResponse(resp)
}

func GetCoin(openId string, param *req.CoinGetReq) *data.Response {
	resp, err := rep().GetCoinRecordByGameID(param.GameId)
	if err != nil {
		logger.Error("查询金币错误", err)
		return data.NewErrorResponse(30019, "查询金币失败")
	}
	return data.NewDataResponse(resp.Coin)
}

func SaveCoin(openId string, param *req.CoinSetReq) *data.Response {
	if err := rep().SaveCoinRecord(openId, param.GameId, param.Coin); err != nil {
		logger.Error("保存金币错误", err)
		return data.NewErrorResponse(30020, "保存金币失败")
	}
	return data.NewSuccessResponse()
}

func SaveTeam(param *req.PlayerTeamReq) *data.Response {

	// 对局校验
	record, err := rep().GetGameRecord(param.GameId)
	if err != nil {
		logger.Error("查询对局报错:", err)
		return data.NewErrorResponse(30021, "查询对局失败")
	}
	if checkRecord(record) {
		logger.Info(record.ID, "对局已结束")
		return data.NewErrorResponse(30022, "对局已结束")
	}
	//  本玩家的OpenId
	openId := record.OpenID
	logger.Info(record.ID, record.Round, " 保存队伍 ", param.Team)

	if err = rep().Transaction(func(r *repo.Repository) (e error) {
		// 删除原有team 虽然不会有 以防万一
		if e = r.DeleteTeamInfoByOpenIdAndRecordId(openId, record.ID, record.Round); e != nil {
			logger.Error("删除team缓存报错:", e)
			return e
		}
		if len(param.Team) == 0 {
			// 空阵容保存999
			if e = r.Create(&repo.TeamInfo{
				OpenID:   openId,
				Position: 999,
				RecordId: record.ID,
				Trophies: record.Trophies,
				Round:    record.Round,
			}); e != nil {
				logger.Error("插入team信息报错:", e)
				return e
			}
		}
		for _, teamInfo := range param.Team {
			if e = r.Create(&repo.TeamInfo{
				OpenID:      openId,
				Position:    teamInfo.Position,
				CharacterID: teamInfo.LocalID,
				LevelNum:    teamInfo.LevelNum,
				BaseAtk:     teamInfo.BaseAtk,
				BaseHp:      teamInfo.BaseHp,
				RecordId:    record.ID,
				Trophies:    record.Trophies,
				Round:       record.Round,
				TempAtk:     teamInfo.TempAtk,
				PerAtk:      teamInfo.PerAtk,
				TempHp:      teamInfo.TempHp,
				PerHp:       teamInfo.PerHp,
				ExSkill:     teamInfo.ExSkill,
			}); e != nil {
				logger.Error("插入team信息报错:", e)
				return e
			}
		}

		return
	}); err != nil {
		return data.NewErrorResponse(30023, "队伍保存失败")
	}
	// 对满阶角色进行更新
	go maxLevelRole(openId, param.Team)
	// 随机一个同回合对手 那肯定不会匹配自己啊
	trophies := record.Trophies
	round := record.Round
	array, warId, err := rep().GetTeamInfosByRound(round, openId)
	if err != nil {
		logger.Error("查询队伍报错", err)
		return data.NewErrorResponse(30024, "查询队伍失败")
	}
	// 组建返回
	chs := make([]*data.TeamInfoResp, 0)
	for _, teamInfo := range array {
		chs = append(chs, &data.TeamInfoResp{
			Position: teamInfo.Position,
			LocalID:  teamInfo.CharacterID,
			LevelNum: teamInfo.LevelNum,
			BaseAtk:  teamInfo.BaseAtk,
			BaseHp:   teamInfo.BaseHp,
			TempAtk:  teamInfo.TempAtk,
			PerAtk:   teamInfo.PerAtk,
			TempHp:   teamInfo.TempHp,
			PerHp:    teamInfo.PerHp,
			ExSkill:  teamInfo.ExSkill,
		})
	}
	warPlayer, err := rep().GetPlayerInfoByOpenId(warId)
	if err != nil && !repo.IsNotFound(err) {
		logger.Error("匹配对手失败", err)
		return data.NewErrorResponse(30025, "匹配对手失败")
	}
	if repo.IsNotFound(err) {
		warPlayer.Nickname = "小U"
		warPlayer.StandId = 20000
		warPlayer.AvatarId = 21000
		warPlayer.AppearanceId = 0
		warPlayer.BackgroundId = 22000
	}
	resp := &data.Player{
		Nickname:     warPlayer.Nickname,
		BackgroundId: warPlayer.BackgroundId,
		StandId:      warPlayer.StandId,
		ExpressionId: warPlayer.ExpressionId,
		AppearanceId: warPlayer.AppearanceId,
		Trophies:     trophies,
		Health:       record.Health,
		Round:        round,
		Team:         chs,
	}
	// 如果需要对上一把金币进行保存
	record.LastBalance = param.LastBalance
	record.LastBalanceRound = record.Round
	if err = rep().UpdateGameRecordBalance(record); err != nil {
		logger.Error("留存金币报错", err)
	}
	return data.NewDataResponse(resp)
}

func checkRecord(record *repo.GameRecord) bool {
	return record.Trophies > 8 || record.Health == 0
}

func maxLevelRole(openId string, team []*req.TeamInfoReq) {
	// 对满阶角色进行更新
	for _, teamInfo := range team {
		if teamInfo.LevelNum == 5 {
			// if err := rep().SaveRoleHistory(&repo.RoleHistory{
			// 	OpenID:      openId,
			// 	CharacterID: teamInfo.LocalID,
			// 	LieId:       teamInfo.LocalID + 10000,
			// 	AvatarId:    teamInfo.LocalID,
			// }); err != nil {
			// 	logger.Error("插入满阶角色信息报错", err)
			// 	return
			// }

			// 是否符合满级解锁角色
			roleId := teamInfo.LocalID
			role, err := rep().GetRoleByCharacterID(roleId)
			if err != nil {
				logger.Error("查询角色信息报错:", err)
				return
			}
			// 同时解锁立绘和头像
			if role.SkinPriceType == 2 {
				if err = rep().Create(&repo.PlayerRoleInfo{
					OpenID:     openId,
					UnlockType: "avatar",
					UnlockId:   roleId,
				}); err != nil {
					return
				}
			}
			if role.LiePriceType == 2 {
				if err = rep().Create(&repo.PlayerRoleInfo{
					OpenID:     openId,
					UnlockType: "lie",
					UnlockId:   roleId,
				}); err != nil {
					return
				}
			}
		}
	}
}

func GameResult(openId string, param *req.PlayerResultReq) *data.Response {
	resp := &data.ResultResp{}

	myInfo, err := rep().GetPlayerInfoByOpenId(openId)
	if err != nil && !repo.IsNotFound(err) {
		logger.Error("查询个人信息失败", err)
		return data.NewErrorResponse(30026, "查询个人信息失败")
	}
	// 对局校验
	record, err := rep().GetGameRecord(param.GameId)
	if err != nil || record.OpenID != openId {
		logger.Error("查询对局报错:", err)
		return data.NewErrorResponse(30027, "查询对局失败")
	}
	if checkRecord(record) {
		logger.Info(record.ID, "对局已结束")
		return data.NewErrorResponse(30028, "对局已结束")
	}
	if param.Result == 1 {
		// 赢了
		// if record.Trophies == 0 {
		// 	record.Trophies = 7
		// }
		record.Trophies++
		// go func() {
		// 胜场增加
		if err = rep().WinsIncrease(openId); err != nil {
			logger.Error("胜场增加错误:", err)
		}
		// 上次对战记录
		result, err := rep().LastGameRecord(openId, record.ID)
		if err != nil {
			logger.Error("查询上次记录错误:", err)
		}
		// 连胜了这货
		if result.Result == 1 {
			winStreak := myInfo.WinStreak
			if winStreak == 0 {
				winStreak = 2
			} else {
				winStreak++
			}
			if err = rep().WinStreakUpdate(openId, winStreak); err != nil {
				logger.Error("连胜增加错误:", err)
			}
			// 如果连胜超过最高连胜次数就更新最高连胜次数
			if winStreak > myInfo.HighestWinStreak {
				if err = rep().HighestWinStreakUpdate(openId, winStreak); err != nil {
					logger.Error("最高连胜更新错误:", err)
				}
			}

		} else {
			if err = rep().WinStreakUpdate(openId, 0); err != nil {
				logger.Error("连胜增加错误:", err)
			}
		}
		// }()

		// 对战记录增加
		if err = rep().Create(&repo.GameResultRecord{
			OpenID: openId,
			GameId: record.ID,
			Result: 1,
			Round:  record.Round,
		}); err != nil {
			logger.Error("记录增加错误:", err)
		}
	} else if param.Result == 0 {
		// 输了
		// if record.Health == 5 {
		// 	record.Health = 2
		// }
		record.Health--
		// 对战记录增加
		if err = rep().Create(&repo.GameResultRecord{
			OpenID: openId,
			GameId: record.ID,
			Result: 0,
			Round:  record.Round,
		}); err != nil {
			logger.Error("记录增加错误:", err)
		}
	} else {
		// 对战记录增加
		if err = rep().Create(&repo.GameResultRecord{
			OpenID: openId,
			GameId: record.ID,
			Result: -1,
			Round:  record.Round,
		}); err != nil {
			logger.Error("记录增加错误:", err)
		}
	}
	// 查询复活历史
	reviveRecord, err := rep().LastReviveRecord(openId, record.ID)
	if err != nil && !repo.IsNotFound(err) {
		logger.Error("查询复活记录报错:", err)
	}
	// 打满了9个
	if record.Trophies == 9 {
		// 复活过了 直接按逻辑发放奖励
		resp.Result = true
		record.Result = true
		myInfo.Trophies += record.Trophies
		if err = rep().TrophiesIncrease(openId, record.Trophies); err != nil {
			logger.Error("奖杯增加报错", err)
			return data.NewErrorResponse(30029, "奖杯增加失败")
		}
		if record.Trophies > 8 {
			// 奖杯属于哪个档次的奖励
			rewardConfig, err := rep().GetRewardConfigBy(record.Trophies)
			if err != nil {
				logger.Error("查询奖励配置信息报错", err)
				return data.NewErrorResponse(30030, "查询奖励失败")
			}
			resp.RewardID = rewardConfig.RewardID
			logger.Info(openId, "宝箱ID:", rewardConfig.RewardID)
			// 根据档次的权重发奖，你啥档次啊
			rewardArray, err := rep().GetRewardByRewardID(rewardConfig.RewardID)
			if err != nil {
				logger.Error("查询奖励配置报错", err)
				return data.NewErrorResponse(30031, "查询奖励失败")
			}
			reward := getRandomReward(rewardArray)
			logger.Info(openId, "奖励池:", reward.RelatedID)
			// 抽出的权重从奖池抽卡，你多冒昧啊
			itemReward, err := rep().GetRandomItemReward(reward.RelatedID)
			if err != nil {
				logger.Error("查询道具奖池报错", err)
				return data.NewErrorResponse(30032, "查询道具失败")
			}
			item, err := rep().GetItemById(itemReward.ItemID)
			if err != nil {
				logger.Error("查询道具报错", err)
				return data.NewErrorResponse(30033, "查询道具失败")
			}

			resp.Reward = &data.WinRewardResp{
				ItemID:      item.ItemID,
				ItemType:    item.ItemType,
				Name:        item.Name,
				Description: item.Description,
				ResourceURL: item.ResourceURL,
			}
			// 判断道具是否重复发布
			if repeatItemResolve(openId, item.ItemID) {
				logger.Info(openId, "重复道具", item.ItemID)

				balance, err := redis().GetRepeatBalance()
				if err != nil {
					logger.Error("查询配置报错报错:", err)
					return data.NewErrorResponse(30034, "查询配置失败")
				}
				if err = rep().BalanceIncrease(openId, uint(balance)); err != nil {
					logger.Error("增加金币报错:", err)
					return data.NewErrorResponse(30035, "增加金币失败")
				}
				resp.Trophies = record.Trophies
				resp.Health = record.Health
				resp.Round = record.Round
				resp.Repeat = true
				resp.BalanceData = uint(balance)
				if err = rep().UpdateGameRecord(record); err != nil {
					logger.Error("更新对局报错", err)
					return data.NewErrorResponse(30036, "更新对局失败")
				}
				return data.NewDataResponse(resp)
			}
			// 道具入库
			if err = rep().Create(&repo.Backpack{
				OpenID:   openId,
				ItemID:   item.ItemID,
				ItemType: item.ItemType,
			}); err != nil {
				logger.Error("背包保存报错:", err)
				return data.NewErrorResponse(30037, "背包保存失败")
			}
		}
		if err = rep().UpdateGameRecord(record); err != nil {
			logger.Error("更新对局报错", err)
			return data.NewErrorResponse(30038, "更新对局失败")
		}
	} else if record.Trophies < 9 && record.Health == 0 {
		// 生命值归0时，若蟠桃数量不足9 提供复活途径
		// 没有复活过 给予复活机会
		if reviveRecord.ID == 0 {
			resp.Revive = true
		} else {
			// 复活过了 直接按逻辑发放奖励
			resp.Result = true
			record.Result = true
			myInfo.Trophies += record.Trophies
			if err = rep().TrophiesIncrease(openId, record.Trophies); err != nil {
				logger.Error("奖杯增加报错", err)
				return data.NewErrorResponse(30039, "增加奖杯失败")
			}
			if record.Trophies > 8 {
				// 奖杯属于哪个档次的奖励
				rewardConfig, err := rep().GetRewardConfigBy(record.Trophies)
				if err != nil {
					logger.Error("查询奖励配置信息报错", err)
					return data.NewErrorResponse(30040, "查询奖励失败")
				}
				resp.RewardID = rewardConfig.RewardID
				logger.Info(openId, "宝箱ID:", rewardConfig.RewardID)
				// 根据档次的权重发奖，你啥档次啊
				rewardArray, err := rep().GetRewardByRewardID(rewardConfig.RewardID)
				if err != nil {
					logger.Error("查询奖励配置报错", err)
					return data.NewErrorResponse(30041, "查询奖励失败")
				}
				reward := getRandomReward(rewardArray)
				logger.Info(openId, "奖励池:", reward.RelatedID)
				// 抽出的权重从奖池抽卡，你多冒昧啊
				itemReward, err := rep().GetRandomItemReward(reward.RelatedID)
				if err != nil {
					logger.Error("查询道具奖池报错", err)
					return data.NewErrorResponse(30042, "查询道具失败")
				}
				item, err := rep().GetItemById(itemReward.ItemID)
				if err != nil {
					logger.Error("查询道具报错", err)
					return data.NewErrorResponse(30043, "查询道具失败")
				}

				resp.Reward = &data.WinRewardResp{
					ItemID:      item.ItemID,
					ItemType:    item.ItemType,
					Name:        item.Name,
					Description: item.Description,
					ResourceURL: item.ResourceURL,
				}
				// 判断道具是否重复发布
				if repeatItemResolve(openId, item.ItemID) {
					logger.Info(openId, "重复道具", item.ItemID)

					balance, err := redis().GetRepeatBalance()
					if err != nil {
						logger.Error("查询配置报错报错:", err)
						return data.NewErrorResponse(30044, "查询配置失败")
					}
					if err = rep().BalanceIncrease(openId, uint(balance)); err != nil {
						logger.Error("增加金币报错:", err)
						return data.NewErrorResponse(30045, "增加金币失败")
					}
					resp.Trophies = record.Trophies
					resp.Health = record.Health
					resp.Round = record.Round
					resp.Repeat = true
					resp.BalanceData = uint(balance)
					if err = rep().UpdateGameRecord(record); err != nil {
						logger.Error("更新对局报错", err)
						return data.NewErrorResponse(30046, "更新对局失败")
					}
					return data.NewDataResponse(resp)
				}
				// 道具入库
				if err = rep().Create(&repo.Backpack{
					OpenID:   openId,
					ItemID:   item.ItemID,
					ItemType: item.ItemType,
				}); err != nil {
					logger.Error("背包保存报错:", err)
					return data.NewErrorResponse(30047, "背包保存失败")
				}
			}
			if err = rep().UpdateGameRecord(record); err != nil {
				logger.Error("更新对局报错", err)
				return data.NewErrorResponse(30048, "更新对局失败")
			}

		}
	}
	resp.Trophies = record.Trophies
	resp.Health = record.Health
	resp.Round = record.Round
	record.Round++
	if err = rep().UpdateGameRecord(record); err != nil {
		logger.Error("更新对局报错", err)
		return data.NewErrorResponse(30049, "更新对局失败")
	}
	if err = rep().SaveCoinRecord(openId, param.GameId, 10); err != nil {
		logger.Error("保存金币错误", err)
		return data.NewErrorResponse(30050, "保存金币失败")
	}
	return data.NewDataResponse(resp)
}

func getRandomReward(rewards []*repo.Reward) *repo.Reward {
	totalWeight := 0
	for _, reward := range rewards {
		totalWeight += reward.Weight
	}
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(totalWeight)
	currentWeight := 0
	for _, reward := range rewards {
		currentWeight += reward.Weight
		if currentWeight > randomNumber {
			return reward
		}
	}
	// 在极端情况下，如果没有奖励满足条件，则返回第一个奖励
	return rewards[len(rewards)-1]
}

// 是否重复获得道具
func repeatItemResolve(openId string, itemId uint) bool {
	// 是否重复获得
	backpack, err := rep().GetBackpackByOpenIdAndItemId(openId, itemId)
	if err != nil && !repo.IsNotFound(err) {
		logger.Error("查询背包出错:", err)
		return true
	}
	if backpack.ID != 0 {
		return true
	}
	return false
}
