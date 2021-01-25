// https://shiro-16.hatenablog.com/entry/2020/05/29/130508
package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

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

func CreateImage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

	card, err := drawFrame()
	if err != nil {
		isErr = true
	}

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

	// reply := ReplyInfo{
	// 	Hashtag: "#Hashtag",
	// 	Age:     20,
	// }
	// w.Header().Del("Content-Type")
	// w.WriteHeader(200)

	// if err := json.NewEncoder(w).Encode(reply); err != nil {
	// 	panic(err)
	// }

	var img image.Image = card
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, img, nil); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}

func drawFrame() (image.Image, error) {
	// f, err := os.Open("white.png")
	// if err != nil {
	// 	fmt.Println("open:", err)
	// 	return
	// }
	// defer f.Close()

	// img, _, err := image.Decode(f)
	// if err != nil {
	// 	fmt.Println("decode:", err)
	// 	return
	// }

	var img image.Image

	fso, err := os.Create("out.png")
	if err != nil {
		fmt.Println("create:", err)
		return img, errors.New("something is nil")
	}
	defer fso.Close()

	m := image.NewRGBA(image.Rect(0, 0, 200, 200)) // 200x200 の画像に test.jpg をのせる
	c := color.RGBA{50, 200, 255, 255}             // RGBA で色を指定(B が 255 なので青?)

	draw.Draw(m, m.Bounds(), &image.Uniform{c}, image.ZP, draw.Src) // 青い画像を描画
	rct := image.Rectangle{image.Point{25, 25}, m.Bounds().Size()}  // test.jpg をのせる位置を指定する(中央に配置する為に横:25 縦:25 の位置を指定)
	draw.Draw(m, rct, img, image.Point{0, 0}, draw.Src)             // 合成する画像を描画
	// jpeg.Encode(fso, m, &jpeg.Options{Quality: 100})

	return m, nil
}
