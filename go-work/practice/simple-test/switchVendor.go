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

var isChangToGetuiToday = false

const oppoStaticsLimit float64 = 9500

func main() {
	sendDingMsg("switch oppo to getui script start")
	//查看oppo当天已经推送了多少条（定时任务）
	go func() {
		clockChan := time.Tick(time.Second * 10)

		i := 0
		for {
		_ :<-clockChan
			i++
			d := time.Now()
			leftTime := time.Date(d.Year(), d.Month(), d.Day(), 0, 0,10,0, d.Location())
			leftString := strconv.FormatInt(leftTime.UnixNano() / 1e6,10)
			rightString := strconv.FormatInt(d.UnixNano() / 1e6,10)

			queryString := "action:oppo_invoke_push AND timestamp:["+ leftString +" TO " + rightString +"] "
			fmt.Println(queryString)
			total := getOppoHits(queryString)
			fmt.Fprintln(os.Stdout, d, ", oppo_invoke_push:", total)

			if i >= 360 {
				i = 0
				if !isChangToGetuiToday  && total > oppoStaticsLimit - 3000{
					sendDingMsg("oppo invoke push num : " + strconv.Itoa(int(total)))
				}
			}

			if !isChangToGetuiToday && total > oppoStaticsLimit{
				result := changeOppoToGetui()
				fmt.Println("changeOppoToGetui update success, result = ", result)
				isChangToGetuiToday = true
				fmt.Fprintln(os.Stdout,"changeOppoToGetui execed, current oppo_invoke_push:", total)
				sendDingMsg("changeOppoToGetui update success")
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
			sendDingMsg("changeGetuiToOppo update success")
			fmt.Println(time.Now(), " changeGetuiToOppo update success, result = ", result)

			for {
				time.Sleep(time.Hour * 24)
				isChangToGetuiToday = false
				result := changeGetuiToOppo()
				sendDingMsg("changeGetuiToOppo update success")
				fmt.Println(time.Now(), " changeGetuiToOppo update success, result = ", result)
			}
		})

	}()

	select {}
}

func changeGetuiToOppo() string {
	//请求切换
}

func changeOppoToGetui() string {
	 //请求切换
}

func getOppoHits(queryString string) float64{

	url := "http://xx/elasticsearch/_msearch"

	payload := strings.NewReader("{\"index\":[\"commonlog-*\"],\"ignore_unavailable\":true,\"preference\":1557999793754}\r\n{\"version\":true,\"size\":1000,\"sort\":[{\"_score\":{\"order\":\"desc\"}}],\"_source\":{\"excludes\":[]},\"stored_fields\":[\"*\"],\"script_fields\":{},\"docvalue_fields\":[\"@timestamp\"],\"query\":{\"bool\":{\"must\":[{\"query_string\":{\"query\":"+"\""+queryString+"\""+",\"analyze_wildcard\":true,\"default_field\":\"*\"}}],\"filter\":[],\"should\":[],\"must_not\":[]}},\"highlight\":{\"pre_tags\":[\"@kibana-highlighted-field@\"],\"post_tags\":[\"@/kibana-highlighted-field@\"],\"fields\":{\"*\":{}},\"fragment_size\":2147483647}}\r\n")

	req, _ := http.NewRequest("POST", url, payload)
	//构建请求的Header

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

func sendDingMsg(msg string)  {
	//钉钉发送消息
}