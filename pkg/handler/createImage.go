package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ReplyInfo struct {
	Hashtag string `json:"hashtag,omitempty"`
	Age     int    `json:"age,omitempty"`
}

func checkIdAndPasswordFormat(uid string) (isErr bool, causeText string) {
	// requireな内容が足りない
	if uid == "" {
		isErr = true
		causeText = "require user_id and password"
	}

	if (len(uid) < 6) || (len(uid) > 20) {
		isErr = true
		causeText = "length of user_id is not correct"
	}

	// err := validation.Validate(uid, is.Alphanumeric)
	var err interface{} = nil // 直接nilは代入できなかった
	if err != nil {
		isErr = true
		causeText = "user_id is not correct format"
	}
	return
}

func createImage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var rep ReplyInfo

	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Bodyから受信内容を読み取る
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &rep)
	fmt.Println(rep)

	// user_idやpasswordが変更されそうだったとき
	if rep.Hashtag != "" {
		reply := ReplyInfo{
			Hashtag: "#nothing",
			Age:     10,
		}
		w.Header().Del("Content-Type")
		w.WriteHeader(400)

		json.NewEncoder(w).Encode(reply)
		return
	}
	isErr, causeText := checkIdAndPasswordFormat(rep.Hashtag)

	if isErr {
		reply := ReplyInfo{
			Hashtag: "#" + causeText,
			Age:     10,
		}
		w.Header().Del("Content-Type")
		w.WriteHeader(400)

		json.NewEncoder(w).Encode(reply)
		return
	}

	reply := ReplyInfo{
		Hashtag: "#Hashtag",
		Age:     20,
	}
	w.Header().Del("Content-Type")
	w.WriteHeader(200)

	if err := json.NewEncoder(w).Encode(reply); err != nil {
		panic(err)
	}
}
