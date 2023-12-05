package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"singo/conf"
	"singo/logger"
	"time"
)

type TokenResponse struct {
	Code    string     `json:"code"`
	Data    *TokenData `json:"data"`
	Message string     `json:"message"`
	Success bool       `json:"success"`
}

type TokenData struct {
	AccessToken  string `json:"accessToken"`
	ExpireTime   int    `json:"expireTime"`
	RefreshToken string `json:"refreshToken"`
}

func getToken() (tokenData *TokenData) {
	url := "https://test-apimall.ubanquan.cn/dapp/token"

	requestBody, err := json.Marshal(map[string]string{
		"appId":     conf.GetConfig().Ubq.AppId,
		"appSecret": conf.GetConfig().Ubq.Secret,
	})
	if err != nil {
		fmt.Println("Error in request body:", err)
		return
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		logger.Error("Error creating request:", err)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error("Error reading response body:", err)
		return
	}

	var tokenResponse TokenResponse
	if err = json.Unmarshal(body, &tokenResponse); err != nil {
		logger.Error("Error decoding JSON:", err)
		return
	}

	if tokenResponse.Success {
		return tokenResponse.Data
	}
	return
}

type FlushResponse struct {
	Code    string     `json:"code"`
	Data    *FlushData `json:"data"`
	Message string     `json:"message"`
	Success bool       `json:"success"`
}

type FlushData struct {
	AccessToken  string `json:"accessToken"`
	ExpireTime   int    `json:"expireTime"`
	RefreshToken string `json:"refreshToken"`
}

func flushToken(refreshToken string) (flushData *FlushData) {
	url := "https://test-apimall.ubanquan.cn/dapp/flush?refreshToken=" + refreshToken

	response, err := http.Get(url)
	if err != nil {
		logger.Error("Error creating request:", err)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error("Error reading response body:", err)
		return
	}

	var flushResponse FlushResponse
	err = json.Unmarshal(body, &flushResponse)
	if err != nil {
		logger.Error("Error decoding JSON:", err)
		return
	}
	if flushResponse.Success {
		return flushResponse.Data
	}
	return
}
func InitUbq() {

	var accessToken, refreshToken string
	var expireTime int
	tokenData := new(TokenData)
	for tokenData.ExpireTime == 0 {
		tokenData = getToken()
	}
	accessToken, refreshToken, expireTime = tokenData.AccessToken, tokenData.RefreshToken, tokenData.ExpireTime

	for {
		logger.Info("获取token : ", accessToken)
		logger.Info("刷新token : ", refreshToken)
		logger.Info("过期时间 : ", expireTime)
		if err := redis().SetUbqToken(accessToken, 2*time.Hour); err != nil {
			logger.Error("保存ubq token报错", err)
			continue
		}
		time.Sleep(time.Hour)
		flushData := new(FlushData)
		for flushData.ExpireTime == 0 {
			flushData = flushToken(tokenData.RefreshToken)
		}
		accessToken, refreshToken, expireTime = flushData.AccessToken, flushData.RefreshToken, flushData.ExpireTime

	}
}

type AuthenticationResponse struct {
	Code    string    `json:"code"`
	Data    *AuthData `json:"data"`
	Message string    `json:"message"`
	Success bool      `json:"success"`
}

type AuthData struct {
	HeadImg  string `json:"headImg"`
	NickName string `json:"nickName"`
	OpenID   string `json:"openId"`
}

func AuthenticationCode(code string) (authResponse *AuthenticationResponse) {
	url := "https://test-apimall.ubanquan.cn/dapp/authentication?code=" + code

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("Error creating request:", err)
		return
	}
	accessToken := redis().GetUbqToken()
	req.Header.Set("access-token", accessToken)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		logger.Error("Error in GET request:", err)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error("Error reading response body:", err)
		return
	}

	err = json.Unmarshal(body, &authResponse)
	if err != nil {
		logger.Error("Error decoding JSON:", err)
		return
	}

	return
}

type CardResponse struct {
	Code    string     `json:"code"`
	Data    []CardData `json:"data"`
	Message string     `json:"message"`
	Success bool       `json:"success"`
}

type CardData struct {
	MetaProductImg  string        `json:"metaProductImg"`
	MetaProductName string        `json:"metaProductName"`
	MetaProductNo   string        `json:"metaProductNo"`
	NfrInfoList     []NfrInfoData `json:"nfrInfoList"`
}

type NfrInfoData struct {
	AuctionNo       string `json:"auctionNo"`
	Cd              int    `json:"cd"`
	CoverImg        string `json:"coverImg"`
	CreationNo      string `json:"creationNo"`
	LockTag         int    `json:"lockTag"`
	MetaProductName string `json:"metaProductName"`
	MetaProductNo   string `json:"metaProductNo"`
	Name            string `json:"name"`
	ProductNo       string `json:"productNo"`
	SerialNum       string `json:"serialNum"`
	ThemeKey        string `json:"themeKey"`
	ThemeName       string `json:"themeName"`
}

func GetCardInfo(openID string) (cardResponse *CardResponse) {
	url := "https://test-apimall.ubanquan.cn/dapp/card?openId=" + openID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("Error creating request:", err)
		return
	}

	accessToken := redis().GetUbqToken()
	req.Header.Set("access-token", accessToken)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		logger.Error("Error in GET request::", err)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error("Error reading response body:", err)
		return
	}

	err = json.Unmarshal(body, &cardResponse)
	if err != nil {
		logger.Error("Error decoding JSON:", err)
		return
	}
	return
}

type UserInfoResponse struct {
	Code    string   `json:"code"`
	Data    UserData `json:"data"`
	Message string   `json:"message"`
	Success bool     `json:"success"`
}

type UserData struct {
	FusionEnergy  int    `json:"fusionEnergy"`
	HasRealName   bool   `json:"hasRealName"`
	HeadImg       string `json:"headImg"`
	NaturalEnergy int    `json:"naturalEnergy"`
	NickName      string `json:"nickName"`
	OpenID        string `json:"openId"`
}

func GetUserInfo(openID string) (userInfoResponse *UserInfoResponse) {
	url := "https://test-apimall.ubanquan.cn/dapp/user/" + openID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("Error creating request:", err)
		return
	}

	accessToken := redis().GetUbqToken()
	req.Header.Set("access-token", accessToken)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		logger.Error("Error in GET request:", err)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Error("Error reading response body:", err)
		return
	}

	err = json.Unmarshal(body, &userInfoResponse)
	if err != nil {
		logger.Error("Error decoding JSON:", err)
		return
	}

	return
}

type DeductionResponse struct {
	Code    string `json:"code"`
	Data    string `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func DeductEnergy(openID string, energy int) (deductionResponse *DeductionResponse) {
	url := fmt.Sprintf("https://test-apimall.ubanquan.cn/dapp/energy/deduction?openId=%s&energy=%d", openID, energy)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		logger.Error("Error creating request:", err)
		return
	}

	accessToken := redis().GetUbqToken()
	req.Header.Set("access-token", accessToken)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		logger.Error("Error in POST request:", err)
		return
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&deductionResponse)
	if err != nil {
		logger.Error("Error decoding JSON:", err)
		return
	}
	return deductionResponse
}

type FuiouWalletResponse struct {
	Code    string `json:"code"`
	Data    bool   `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func CheckFuiouWallet(openID string) (fuiouWalletResponse *FuiouWalletResponse) {
	url := "https://test-apimall.ubanquan.cn/dapp/pool/isFuiouWallet?openId=" + openID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("Error creating request:", err)
		return
	}
	accessToken := redis().GetUbqToken()
	req.Header.Set("access-token", accessToken)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		logger.Error("Error in GET request:", err)
		return
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&fuiouWalletResponse)
	if err != nil {
		logger.Error("Error decoding JSON:", err)
		return
	}
	return fuiouWalletResponse
}
