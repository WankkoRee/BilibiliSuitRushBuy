package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func init() {
	//proxyUrl, err := url.Parse("http://127.0.0.1:10801")
	//if err != nil {
	//	log.Fatal(err)
	//	return
	//}
	//http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
}

func orderCreate(itemId int, addMonth int, buyNum int64, csrf string, apiCookie string, userAgent string, accessKey string) string {
	req, err := http.NewRequest("POST", "https://api.bilibili.com/x/garb/trade/create", strings.NewReader(fmt.Sprintf("item_id=%d&platform=android&currency=bp&add_month=%d&buy_num=%d&coupon_token=&hasBiliapp=true&csrf=%s", itemId, addMonth, buyNum, csrf)))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("native_api_from", "h5")
	req.Header.Set("Referer", fmt.Sprintf("https://www.bilibili.com/h5/mall/suit/detail?id=%d&navhide=1", itemId))
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	req.Header.Set("Cookie", apiCookie)
	req.Header.Set("X-CSRF-TOKEN", csrf)
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Fatal(err)
	}
	respb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(respb))

	OrderCreateResultCode, err := jsonparser.GetInt(respb, "code")
	if err != nil {
		log.Fatal(err)
	}
	if OrderCreateResultCode == 26127 {
		return ""
	} else if OrderCreateResultCode != 0 {
		log.Fatal(string(respb))
	}

	payData0, err := jsonparser.GetString(respb, "data", "pay_data")
	if err != nil {
		log.Fatal(err)
	}
	payData1, err := jsonparser.Set([]byte(payData0), []byte(`"`+accessKey+`"`), "accessKey")
	if err != nil {
		log.Fatal(err)
	}
	payData2, err := jsonparser.Set(payData1, []byte(`"tv.danmaku.bili"`), "appName")
	if err != nil {
		log.Fatal(err)
	}
	payData1, err = jsonparser.Set(payData2, []byte("6560300"), "appVersion")
	if err != nil {
		log.Fatal(err)
	}
	payData2, err = jsonparser.Set(payData1, []byte(`"ANDROID"`), "device")
	if err != nil {
		log.Fatal(err)
	}
	payData1, err = jsonparser.Set(payData2, []byte(`"3"`), "deviceType")
	if err != nil {
		log.Fatal(err)
	}
	payData2, err = jsonparser.Set(payData1, []byte(`"WiFi"`), "network")
	if err != nil {
		log.Fatal(err)
	}
	payData1, err = jsonparser.Set(payData2, []byte(`"bp"`), "payChannel")
	if err != nil {
		log.Fatal(err)
	}
	payData2, err = jsonparser.Set(payData1, []byte("99"), "payChannelId")
	if err != nil {
		log.Fatal(err)
	}
	payData1, err = jsonparser.Set(payData2, []byte(`"bp"`), "realChannel")
	if err != nil {
		log.Fatal(err)
	}
	payData2, err = jsonparser.Set(payData1, []byte(`"1.4.9"`), "sdkVersion")
	if err != nil {
		log.Fatal(err)
	}
	return string(payData2)
}

func pay(payData string, userAgent string, Buvid string, DeviceID string, fpLocal string, fpRemote string, sessionId string, deviceFingerprint string) (string, error) {
	req, err := http.NewRequest("POST", "https://pay.bilibili.com/payplatform/pay/pay", strings.NewReader(payData))
	if err != nil {
		return "", err
	}
	req.Header.Set("cLocale", "zh_CN")
	req.Header.Set("sLocale", "zh_CN")
	req.Header.Set("Buvid", Buvid)
	req.Header.Set("Device-ID", DeviceID)
	req.Header.Set("fp_local", fpLocal)
	req.Header.Set("fp_remote", fpRemote)
	req.Header.Set("session_id", sessionId)
	req.Header.Set("deviceFingerprint", deviceFingerprint)
	req.Header.Set("buildId", "6560300")
	req.Header.Set("env", "prod")
	req.Header.Set("APP-KEY", "android64")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("APP-KEY", "android64")
	req.Header.Set("bili-bridge-engine", "cronet")
	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", err
	}
	respb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	log.Println(string(respb))
	payChannelParam, err := jsonparser.GetString(respb, "data", "payChannelParam")
	if err != nil {
		return "", err
	}
	return payChannelParam, nil
}

func payBp(payChannelParam string, userAgent string, Buvid string, DeviceID string, fpLocal string, fpRemote string, sessionId string, deviceFingerprint string) (string, error) {
	req, err := http.NewRequest("POST", "https://pay.bilibili.com/paywallet/pay/payBp", strings.NewReader(payChannelParam))
	if err != nil {
		return "", err
	}
	req.Header.Set("cLocale", "zh_CN")
	req.Header.Set("sLocale", "zh_CN")
	req.Header.Set("Buvid", Buvid)
	req.Header.Set("Device-ID", DeviceID)
	req.Header.Set("fp_local", fpLocal)
	req.Header.Set("fp_remote", fpRemote)
	req.Header.Set("session_id", sessionId)
	req.Header.Set("deviceFingerprint", deviceFingerprint)
	req.Header.Set("buildId", "6560300")
	req.Header.Set("env", "prod")
	req.Header.Set("APP-KEY", "android64")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("APP-KEY", "android64")
	req.Header.Set("bili-bridge-engine", "cronet")
	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", err
	}
	respb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respb), err
}

func buy(payData string, userAgent string, Buvid string, DeviceID string, fpLocal string, fpRemote string, sessionId string, deviceFingerprint string) (string, error) {
	payChannelParam, err := pay(payData, userAgent, Buvid, DeviceID, fpLocal, fpRemote, sessionId, deviceFingerprint)
	if err != nil {
		return "", err
	}

	payResult, err := payBp(payChannelParam, userAgent, Buvid, DeviceID, fpLocal, fpRemote, sessionId, deviceFingerprint)
	if err != nil {
		return "", err
	}
	return payResult, nil
}

func watch(itemId int, limitT int64) bool {
	log.Println("开始观望")
	for {
		resp, err := http.Get(fmt.Sprintf("https://api.bilibili.com/x/garb/mall/item/suit/v2?item_id=%d&part=suit", itemId))
		if err != nil {
			log.Fatal(err)
			return false
		}
		respb, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
			return false
		}
		SuitRecentResultCode, err := jsonparser.GetInt(respb, "code")
		if err != nil {
			log.Fatal(err)
			return false
		}
		if SuitRecentResultCode != 0 {
			log.Fatal(respb)
			return false
		}

		saleLeftTime, err := jsonparser.GetInt(respb, "data", "sale_left_time")
		if err != nil {
			log.Fatal(err)
			return false
		}
		log.Println(saleLeftTime)

		if saleLeftTime <= limitT {
			log.Println("开始抢购")
			return true
		} else if saleLeftTime <= limitT*10 { // 10倍观望时间内加速刷新频率
			time.Sleep(time.Duration(1) * time.Second)
		} else {
			time.Sleep(time.Duration(2) * time.Second)
		}
	}
}

func main() {
	// 请务必确保脚本运行时号里的B币充足，>=永久价格*buyNum

	var limitT int64 = 3 // 观望时无CD刷新的临界值，当前时间与目标装扮开售时间相差小于等于临界值(秒)时，脚本观望将不再sleep，开启暴风下单状态
	// 建议3s~5s，过小容易错过，过大容易ban IP，请自行合理预估

	itemId := 33998      // 2022拜年纪 33998
	addMonth := -1       // -1:永久, 1:一个月，应该不会有人买一个月吧，这个我功能性测试时候用的，就是为了测试能省点钱
	var buyNum int64 = 1 // 一次性买几个

	// 以下参数在余额不足时抓取下单请求即可获取，在https://api.bilibili.com/x/garb/trade/create数据包中可找到
	// 以下参数在余额充足时抓取下单请求亦可获取，在https://api.bilibili.com/x/garb/trade/create数据包中可找到
	// 余额充足时务必开启发送断点，切勿放行https://pay.bilibili.com/paywallet/pay/payBp，否则将会购买成功
	apiCsrf := ""
	apiCookie := ""
	apiUserAgent := ""
	// 以下参数在余额不足时抓取下单请求即可获取，在https://pay.bilibili.com/paywallet/recharge/getRechargePanel/v2数据包中可找到
	// 以下参数在余额充足时抓取下单请求亦可获取，在https://pay.bilibili.com/payplatform/pay/pay数据包中可找到
	// 余额充足时务必开启发送断点，切勿放行https://pay.bilibili.com/paywallet/pay/payBp，否则将会购买成功
	apiAccessKey := ""
	payUserAgent := ""
	payBuvid := ""
	payDeviceID := ""
	payFpLocal := ""
	payFpRemote := ""
	paySessionId := ""
	payDeviceFingerprint := ""

	if watch(itemId, limitT) { // 观望，不使用Cookie，防止被临时ban号
		payData := ""
		for payData == "" {
			payData = orderCreate(itemId, addMonth, buyNum, apiCsrf, apiCookie, apiUserAgent, apiAccessKey) // 下单
		}
		payResult, err := buy(payData, payUserAgent, payBuvid, payDeviceID, payFpLocal, payFpRemote, paySessionId, payDeviceFingerprint) // 购买
		if err != nil {
			log.Fatal(err)
		}
		log.Println(payResult)
	}

	log.Println("End.")
}
