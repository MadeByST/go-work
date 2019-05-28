package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var token = "7781b5a5-aa05-4899-9142-e89839204d9b"

//当天已经改变Oppo的推送Getui的标志位
var isChangToGetuiToday = false

const oppoStaticsLimit float64 = 9500

func main() {

	//查看oppo当天已经推送了多少条（定时任务）
	go func() {
		minClock := time.Tick(time.Second * 10)

		for {

		_ :<- minClock
			d := time.Now()
			leftTime := time.Date(d.Year(), d.Month(), d.Day(), 0, 0,10,0, d.Location())

			leftString := strconv.FormatInt(leftTime.UnixNano() / 1e6,10)
			rightString := strconv.FormatInt(d.UnixNano() / 1e6,10)

			queryString := "action:oppo_invoke_push AND timestamp:["+ leftString +" TO " + rightString +"] AND code:0"
			fmt.Println(queryString)

			total := getOppoHits(queryString)

			fmt.Fprintln(os.Stdout, d, ", oppo_invoke_push:", total)

			if !isChangToGetuiToday && total > oppoStaticsLimit{
				result := changeOppoToGetui()
				fmt.Println("changeOppoToGetui update success, result = ", result)
				isChangToGetuiToday = true
				fmt.Fprintln(os.Stdout,"changeOppoToGetui execed, current oppo_invoke_push:", total)
			}

		}

	}()


	go func() {

		//每天0点由Getui切换oppo
		hour, min, sec := time.Now().Clock()
		//() * time.Hour + (59 - min) * time.Minute + (59 - sec) * time.Second
		//zeroTime := time.After()
		leftTime := time.Duration(23-hour)*time.Hour + time.Duration(59-min)*time.Minute + time.Duration(59-sec)*time.Second
		leftTime += time.Minute * 10
		fmt.Fprintln(os.Stdout, "next update time at ", time.Now().Add(leftTime))

		time.AfterFunc(leftTime, func() {
			isChangToGetuiToday = false
			result := changeGetuiToOppo()
			fmt.Println(time.Now(), " changeGetuiToOppo update success, result = ", result)
		})

		for {
			time.Sleep(time.Hour * 24)
			isChangToGetuiToday = false
			result := changeGetuiToOppo()
			fmt.Println(time.Now(), " changeGetuiToOppo update success, result = ", result)
		}
	}()

	changeGetuiToOppoUrl := "/changeGetuiToOppo"
	changeOppoToGetuiUrl := "/changeOppoToGetui"

	fmt.Fprintln(os.Stdout, "changeGetuiToOppo: ", changeGetuiToOppoUrl + "?token=" + token)
	fmt.Fprintln(os.Stdout, "changeOppoToGetui: ", changeOppoToGetuiUrl + "?token=" + token)

	select {}

}

func changeGetuiToOppo() string {
	url := "https://api.msg.xescdn.com/push/v1/conf/push/?appId=xes10001"

	payload := strings.NewReader("{\"preferVendersForVersions\":[{\"versionList\":[\"default\"],\"preferVendersForBrands\":[{\"brand\":\"HUAWEI\",\"maxSdkNum\":2,\"notification\":[4,1],\"transmission\":[1,4],\"bindVenders\":[1,4]},{\"brand\":\"XIAOMI\",\"maxSdkNum\":2,\"notification\":[3,1],\"transmission\":[1,3],\"bindVenders\":[1,3]},{\"brand\":\"IPHONE\",\"maxSdkNum\":2,\"notification\":[1,2],\"transmission\":[1,2],\"bindVenders\":[1,2]},{\"brand\":\"VIVO\",\"maxSdkNum\":2,\"notification\":[1,2],\"transmission\":[1,2],\"bindVenders\":[1,2]},{\"brand\":\"OTHER\",\"maxSdkNum\":2,\"notification\":[1,2],\"transmission\":[1,2],\"bindVenders\":[1,2]}]},{\"versionList\":[\"v1.8\"],\"preferVendersForBrands\":[{\"brand\":\"HUAWEI\",\"maxSdkNum\":2,\"notification\":[4,1],\"transmission\":[1,4],\"bindVenders\":[1,4]},{\"brand\":\"XIAOMI\",\"maxSdkNum\":2,\"notification\":[3,1],\"transmission\":[1,3],\"bindVenders\":[1,3]},{\"brand\":\"IPHONE\",\"maxSdkNum\":2,\"notification\":[1,2],\"transmission\":[1,2],\"bindVenders\":[1,2]},{\"brand\":\"VIVO\",\"maxSdkNum\":2,\"notification\":[7,1,2],\"transmission\":[1,7,2],\"bindVenders\":[1,7,2]},{\"brand\":\"OPPO\",\"maxSdkNum\":2,\"notification\":[1,2,6],\"transmission\":[1,2,6],\"bindVenders\":[1,2,6]},{\"brand\":\"OTHER\",\"maxSdkNum\":2,\"notification\":[1,2],\"transmission\":[1,2],\"bindVenders\":[1,2]}]},{\"versionList\":[\"v1.9\"],\"preferVendersForBrands\":[{\"brand\":\"HUAWEI\",\"maxSdkNum\":2,\"notification\":[4,1],\"transmission\":[1,4],\"bindVenders\":[1,4]},{\"brand\":\"XIAOMI\",\"maxSdkNum\":2,\"notification\":[3,1],\"transmission\":[1,3],\"bindVenders\":[1,3]},{\"brand\":\"IPHONE\",\"maxSdkNum\":2,\"notification\":[1,2],\"transmission\":[1,2],\"bindVenders\":[1,2]},{\"brand\":\"VIVO\",\"maxSdkNum\":2,\"notification\":[7,1,2],\"transmission\":[1,7,2],\"bindVenders\":[1,7,2]},{\"brand\":\"OPPO\",\"maxSdkNum\":2,\"notification\":[6,1,2],\"transmission\":[1,6,2],\"bindVenders\":[1,6,2]},{\"brand\":\"OTHER\",\"maxSdkNum\":2,\"notification\":[1,2],\"transmission\":[1,2],\"bindVenders\":[1,2]}]}]}")

	req, _ := http.NewRequest("PUT", url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-PS-AppID", "xes10001")
	req.Header.Add("X-PS-Timestamp", "1560596400")
	req.Header.Add("X-PS-Version", "1")
	req.Header.Add("X-PS-Signature", "f51df7b081baaa0300c257ff82a1e3d6")
	req.Header.Add("User-Agent", "PostmanRuntime/7.11.0")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Postman-Token", "402ceeb5-4108-48f3-a33d-2b66372fa412,8bbdec79-e220-4481-9a6f-71c51dfa39b0")
	req.Header.Add("Host", "api.msg.xescdn.com")
	req.Header.Add("accept-encoding", "gzip, deflate")
	req.Header.Add("content-length", "1815")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return string(body)
}

func changeOppoToGetui() string {
	url := "https://api.msg.xescdn.com/push/v1/conf/push/?appId=xes10001"

	payload := strings.NewReader("{\"preferVendersForVersions\":[{\"versionList\":[\"default\"],\"preferVendersForBrands\":[{\"brand\":\"HUAWEI\",\"maxSdkNum\":2,\"notification\":[4,1],\"transmission\":[1,4],\"bindVenders\":[1,4]},{\"brand\":\"XIAOMI\",\"maxSdkNum\":2,\"notification\":[3,1],\"transmission\":[1,3],\"bindVenders\":[1,3]},{\"brand\":\"IPHONE\",\"maxSdkNum\":2,\"notification\":[1,2],\"transmission\":[1,2],\"bindVenders\":[1,2]},{\"brand\":\"VIVO\",\"maxSdkNum\":2,\"notification\":[1,2],\"transmission\":[1,2],\"bindVenders\":[1,2]},{\"brand\":\"OTHER\",\"maxSdkNum\":2,\"notification\":[1,2],\"transmission\":[1,2],\"bindVenders\":[1,2]}]},{\"versionList\":[\"v1.8\"],\"preferVendersForBrands\":[{\"brand\":\"HUAWEI\",\"maxSdkNum\":2,\"notification\":[4,1],\"transmission\":[1,4],\"bindVenders\":[1,4]},{\"brand\":\"XIAOMI\",\"maxSdkNum\":2,\"notification\":[3,1],\"transmission\":[1,3],\"bindVenders\":[1,3]},{\"brand\":\"IPHONE\",\"maxSdkNum\":2,\"notification\":[1,2],\"transmission\":[1,2],\"bindVenders\":[1,2]},{\"brand\":\"VIVO\",\"maxSdkNum\":2,\"notification\":[7,1,2],\"transmission\":[1,7,2],\"bindVenders\":[1,7,2]},{\"brand\":\"OPPO\",\"maxSdkNum\":2,\"notification\":[1,2,6],\"transmission\":[1,2,6],\"bindVenders\":[1,2,6]},{\"brand\":\"OTHER\",\"maxSdkNum\":2,\"notification\":[1,2],\"transmission\":[1,2],\"bindVenders\":[1,2]}]},{\"versionList\":[\"v1.9\"],\"preferVendersForBrands\":[{\"brand\":\"HUAWEI\",\"maxSdkNum\":2,\"notification\":[4,1],\"transmission\":[1,4],\"bindVenders\":[1,4]},{\"brand\":\"XIAOMI\",\"maxSdkNum\":2,\"notification\":[3,1],\"transmission\":[1,3],\"bindVenders\":[1,3]},{\"brand\":\"IPHONE\",\"maxSdkNum\":2,\"notification\":[1,2],\"transmission\":[1,2],\"bindVenders\":[1,2]},{\"brand\":\"VIVO\",\"maxSdkNum\":2,\"notification\":[7,1,2],\"transmission\":[1,7,2],\"bindVenders\":[1,7,2]},{\"brand\":\"OPPO\",\"maxSdkNum\":2,\"notification\":[1,6,2],\"transmission\":[1,6,2],\"bindVenders\":[1,6,2]},{\"brand\":\"OTHER\",\"maxSdkNum\":2,\"notification\":[1,2],\"transmission\":[1,2],\"bindVenders\":[1,2]}]}]}")

	req, _ := http.NewRequest("PUT", url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-PS-AppID", "xes10001")
	req.Header.Add("X-PS-Timestamp", "1560596400")
	req.Header.Add("X-PS-Version", "1")
	req.Header.Add("X-PS-Signature", "f51df7b081baaa0300c257ff82a1e3d6")
	req.Header.Add("User-Agent", "PostmanRuntime/7.11.0")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Postman-Token", "cb2e6bf5-f415-44ae-b46a-5cf9e471892e,cc3ef9e1-84ef-44f3-8b65-1d49044e1fea")
	req.Header.Add("Host", "api.msg.xescdn.com")
	req.Header.Add("accept-encoding", "gzip, deflate")
	req.Header.Add("content-length", "1815")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return string(body)
}

func getOppoHits(queryString string) float64{

	url := "http://kibana.platform.com/elasticsearch/_msearch"

	payload := strings.NewReader("{\"index\":[\"commonlog-*\"],\"ignore_unavailable\":true,\"preference\":1557999793754}\r\n{\"version\":true,\"size\":1000,\"sort\":[{\"_score\":{\"order\":\"desc\"}}],\"_source\":{\"excludes\":[]},\"stored_fields\":[\"*\"],\"script_fields\":{},\"docvalue_fields\":[\"@timestamp\"],\"query\":{\"bool\":{\"must\":[{\"query_string\":{\"query\":"+"\""+queryString+"\""+",\"analyze_wildcard\":true,\"default_field\":\"*\"}}],\"filter\":[],\"should\":[],\"must_not\":[]}},\"highlight\":{\"pre_tags\":[\"@kibana-highlighted-field@\"],\"post_tags\":[\"@/kibana-highlighted-field@\"],\"fields\":{\"*\":{}},\"fragment_size\":2147483647}}\r\n")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/x-ndjson")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.131 Safari/537.36")
	req.Header.Add("Referer", "http://kibana.platform.com/app/kibana")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Origin", "http://kibana.platform.com")
	req.Header.Add("Authorization", "Basic YWRtaW46eXVucGluZ3RhaTE=,Basic YWRtaW46eXVucGluZ3RhaTE=")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", "kibana.platform.com")
	req.Header.Add("kbn-version", "6.2.2")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Proxy-Connection", "keep-alive")
	req.Header.Add("Postman-Token", "936757c8-55e7-45f1-960c-9d472bbf1de3,5095ffe8-cac2-4503-9640-e60a34349c1b")
	req.Header.Add("content-length", "613")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return -1
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return -1
	}
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(body), &dat); err != nil {
		return -1
	}

	if responses, ok := dat["responses"]; ok {
		resps, ok := responses.([]interface{})
		if ok {
			if len(resps) > 0  {
				resp0 := resps[0]
				respMap0, ok := resp0.(map[string]interface{})
				if ok {
					hits, ok:= respMap0["hits"]
					if ok {
						hitsMap, ok := hits.(map[string]interface{})
						if ok {
							total,ok := hitsMap["total"]
							if ok {
								t, ok:= total.(float64)
								if ok {
									return t
								}
							}
						}
					}

				}
			}
		}
	}
	return -1
}