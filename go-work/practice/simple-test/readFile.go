package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/baidubce/bce-sdk-go/util/log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
)

type UserIdToPsId struct {
	UserId string `json:"userId"`
	PsId   string `json:"psId"`
}

func read(file string) (userSlice []UserIdToPsId, err error) {

	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	read := bufio.NewReader(fd)
	for {
		line, err := read.ReadString('\n') //以'\n'为结束符读入一行
		if err != nil || io.EOF == err {
			break
		}
		split := strings.Split(line, ",")

		user := UserIdToPsId{}
		user.UserId = strings.TrimSpace(split[0])
		user.PsId = strings.TrimSpace(split[1])

		userSlice = append(userSlice, user)
	}

	return
}

type RestPushMsg struct {
	AppId       string   `json:"appId"`
	BusinessId  string   `json:"businessId"`
	SkipType    int      `json:"skipType"`
	MessageType int      `json:"messageType"`
	PushType    int      `json:"pushType"`
	Uids        []string `json:"uids"`
	Tags        []string `json:"tags"`
	Save        bool     `json:"save"`
	PushTime    int64    `json:"pushTime"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Payload     string   `json:"payload"`
}

type RestPushResp struct {
	State   int           `json:"state"`
	Message string        `json:"message"`
	Content *TaskIdHolder `json:"content"`
}

type TaskIdHolder struct {
	TaskId string `json:"taskId"`
}

func invokePush(pushMsg *RestPushMsg) (resp *RestPushResp, err error) {

	url := PS_PUSH_URL

	reqBody, err := json.Marshal(pushMsg)

	if err != nil {
		return
	}

	payload := strings.NewReader(string(reqBody))

	req, err := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-PS-AppID", "xes20001")
	req.Header.Add("X-PS-Timestamp", "1560596400")
	req.Header.Add("X-PS-Version", "1")
	req.Header.Add("X-PS-Signature", "aadasdasdasd")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return
	}
	resp = new(RestPushResp)
	err = json.Unmarshal(body, resp)

	return

}

type TaskStat struct {
	UserIdToPsId *UserIdToPsId
	PushResp     *RestPushResp
}

const (
	MAX_TRY_NUM      = 3
	GO_GOROUTINE_NUM = 4
	PS_PUSH_URL      = "https://xxx/push/v1/msg/push/?appId=xes20001"
)

var wg sync.WaitGroup

func pushWorker(wg *sync.WaitGroup, psIdChan <-chan UserIdToPsId, failRes chan<- string, successRes chan<- string) {
	for userIdToPsId := range psIdChan {
		pushMsg := &RestPushMsg{
			BusinessId:  "10014",
			AppId:       "xes20001",
			SkipType:    0,
			MessageType: 0,
			PushType:    0,
			Uids:        []string{userIdToPsId.PsId},
			Tags:        []string{},
			Save:        false,
			PushTime:    0,
			Title:       "您报名的课程讲次已更新！",
			Description: "诗词大会冠军陈更【理科思维记古诗】第2讲、第3讲已更新完结，学完即可领取结课证书！ >>",
			Payload:     "",
		}

		for i := 0; i < MAX_TRY_NUM; i++ {
			pushResp, err := invokePush(pushMsg)
			taskStat := TaskStat{&userIdToPsId, pushResp}
			res, _ := json.Marshal(taskStat)
			if err != nil {
				failRes <- string(res)
				continue
			}
			if pushResp.State != 0 {
				failRes <- string(res)
				continue
			}

			successRes <- string(res)
			break
		}
	}
	wg.Done()
}

func main() {

	psIds := make(chan UserIdToPsId)
	failRes := make(chan string)
	successRes := make(chan string)

	userIdToPsIds, err := read("user.txt")

	fmt.Println(len(userIdToPsIds))

	if err != nil {
		fmt.Println(err)
	}

	out, err := os.OpenFile("out.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

	if err != nil {
		fmt.Println(err)
	}

	errout, err := os.OpenFile("error.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)

	if err != nil {
		log.Fatal(err)
	}

	wg.Add(GO_GOROUTINE_NUM)
	for i := 0; i < GO_GOROUTINE_NUM; i++ {
		go pushWorker(&wg, psIds, failRes, successRes)
	}

	go func() {
		for index := range userIdToPsIds {
			userIdToPsId := userIdToPsIds[index]
			psIds <- userIdToPsId
		}
		close(psIds)
	}()

	var endChan = make(chan struct{})

	//错误输出
	go func() {
		for res := range failRes {
			_, _ = fmt.Fprintf(errout, "%s\n", res)
		}
		endChan <- struct{}{}
	}()

	//成功输出
	go func() {
		for res := range successRes {
			_, _ = fmt.Fprintf(out, "%s\n", res)
		}
		endChan <- struct{}{}
	}()

	//所有协程推出后关闭 输出chan
	go func() {
		wg.Wait()
		close(failRes)
		close(successRes)
	}()

	//等待文件io完成
	_ = <-endChan
	_ = <-endChan
}
