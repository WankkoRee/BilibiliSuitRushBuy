package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

func watch(itemId int, targetId int64, limitV int64, buyNum *int64) bool {
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

		saleQuantity, err := jsonparser.GetString(respb, "data", "item", "properties", "sale_quantity")
		if err != nil {
			log.Fatal(err)
			return false
		}
		saleQuantityI, err := strconv.ParseInt(saleQuantity, 10, 64)
		if err != nil {
			log.Fatal(err)
			return false
		}
		saleSurplus, err := jsonparser.GetInt(respb, "data", "sale_surplus")
		if err != nil {
			log.Fatal(err)
			return false
		}
		nowId := saleQuantityI - saleSurplus
		log.Println(nowId)

		if nowId < targetId && nowId+*buyNum >= targetId {
			log.Println("开始抢购")
			*buyNum = targetId - nowId
			return true
		} else if nowId >= targetId {
			log.Println("已经错过")
			return false
		}

		if nowId < targetId-limitV {
			time.Sleep(time.Duration(2) * time.Second)
		}
	}
}

func main() {
	// 请务必确保脚本运行时号里的B币充足，>=永久价格*buyNum

	var limitV int64 = 20 // 观望时无CD刷新的临界值，当前id与目标id相差小于等于临界值时，脚本观望将不再sleep，开启暴风刷新状态
	// 个人推荐：装扮热门时，id冷门10, id中等70, id热门140；过小容易错过，过大容易ban IP，请自行合理预估

	itemId := 33998            // 2022拜年纪 33998
	addMonth := -1             // -1:永久, 1:一个月，应该不会有人买一个月吧，这个我功能性测试时候用的，就是为了测试能省点钱
	var targetId int64 = 18168 // 目标粉丝编号，豹子号和炸弹号基本上都有人x10、x10地抢，所以如果有这类需求的话就把下面的参数改大点，大于等于其他人即可，就是比较烧钱
	var buyNum int64 = 1       // 一次性买几个，当需要抢热门号时建议多买几个，因为号贩子就是这样的，由于是一次性购入，脚本无法比其快，所以只能以暴制暴；如果是冷门id的话那就无所谓了

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

	if watch(itemId, targetId, limitV, &buyNum) { // 观望，并且传递最多需要买几个，不使用Cookie，防止被临时ban号
		payData := orderCreate(itemId, addMonth, buyNum, apiCsrf, apiCookie, apiUserAgent, apiAccessKey)                                 // 下单
		payResult, err := buy(payData, payUserAgent, payBuvid, payDeviceID, payFpLocal, payFpRemote, paySessionId, payDeviceFingerprint) // 购买
		if err != nil {
			log.Fatal(err)
		}
		log.Println(payResult)
	}

	log.Println("End.")
}
