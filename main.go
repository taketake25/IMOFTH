// https://lnly.hatenablog.com/entry/2020/02/26/225722
// https://shiro-16.hatenablog.com/entry/2020/05/29/130508
// curl -d "{\"Hashtag\":\"aaaaaaa\", \"Age\":40}" -H "Content-type: application/json" -X POST localhost:8080/createImage

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"

	"github.com/golang/freetype/truetype"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// "database/sql"
// https://lnly.hatenablog.com/entry/2020/02/26/225722

type CardInfo struct {
	Hashtag     string `json:"hashtag,omitempty"`
	Age         int    `json:"age,omitempty"`
	Position    string `json:"position,omitempty"`
	Sex         int    `json:"sex,omitempty"`
	Work        string `json:"work,omitempty"`
	Background1 string `json:"background1,omitempty"`
	Background2 string `json:"background2,omitempty"`
	Background3 string `json:"background3,omitempty"`
	TwitterId   string `json:"twitterId,omitempty"`
}

type ReplyInfo struct {
	Message string `json:"message,omitempty"`
	Cause   string `json:"cause,omitempty"`
}

type Page struct {
	Title     string
	CardImage string
}

var ImageTemplate string = `<!DOCTYPE html>
<html lang="en"><head></head>
<body><img src="data:image/jpg;base64,{{.Image}}"></body>`

func sendJsonResponse(w http.ResponseWriter, headerVal int, repMessage string, causeMessage string) {
	reply := ReplyInfo{
		Message: "#" + repMessage,
		Cause:   causeMessage,
	}
	w.Header().Del("Content-Type")
	w.WriteHeader(headerVal)

	json.NewEncoder(w).Encode(reply)
	return
}

func checkInputFormat(input CardInfo) (isErr bool, causeText string) {
	// 実装すべき点；
	// 1.必須情報が入力されているか
	// 1.5. 正しい型であるか
	// 2.文字列の長さが長すぎないか
	// 3.変な記号が含まれていないか
	// 4.Hashtagが既に存在しているか
	// 5.詳細すぎる情報がかかれていないか
	// 6.多重投稿者でないか

	// requireな内容が足りない
	log.Println(input.Hashtag)
	if input.Hashtag == "" {
		isErr = true
		causeText = "require hashtag"
	}

	if len(input.Hashtag) > 15 {
		isErr = true
		causeText = "length of hashtag is not correct"
	}

	// err := validation.Validate(uid, is.Alphanumeric)
	var err interface{} = nil // 直接nilは代入できなかった
	if err != nil {
		isErr = true
		causeText = "hashtag is not correct format"
	}
	return
}

func CreateTemplate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func CreateImage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Println("createImage")
	var input CardInfo

	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Bodyから受信内容を読み取る
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &input)
	log.Println("request:", input)

	// ハッシュタグが空白だった時
	if input.Hashtag == "" {
		sendJsonResponse(w, 400, "Please input hashtag name", "value empty")
		return
	}

	// requestの入力情報が正しいフォーマットかをチェック
	isErr, causeText := checkInputFormat(input)

	card, status := drawFrame(input)
	if status == false {
		isErr = true
	}

	if isErr {
		sendJsonResponse(w, 400, "create image error", ""+causeText)
		return
	}

	fso, err := os.Create("out.png")
	// fso, err := os.Create("out.jpg")
	defer fso.Close()
	if err != nil {
		log.Println("create error:", err)
	}
	defer fso.Close()

	// 画像を送信用バッファに代入する
	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, card); err != nil {
		// if err := jpeg.Encode(fso, m, nil); err != nil {
		log.Println("error:png\n", err)
	}
	if err := png.Encode(fso, card); err != nil {
		log.Println("error:png\n", err)
	}

	// log.Println(buffer.Bytes())

	// http形式で送信する場合
	// str := base64.StdEncoding.EncodeToString(buffer.Bytes())
	// if tmpl, err := template.New("image").Parse(ImageTemplate); err != nil {
	// 	log.Println("unable to parse image template.")
	// } else {
	// 	data := map[string]interface{}{"Image": str}
	// 	if err = tmpl.Execute(w, data); err != nil {
	// 		log.Println("unable to execute template.")
	// 	}
	// }

	// w.Header().Set("Content-Type", "text")
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}

func drawTextOnImage(dr *font.Drawer, drawedText string, x, y int) {
	// dr.Dot.X = (fixed.I(x) - dr.MeasureString(drawedText)) / 2
	dr.Dot.X = fixed.I(x)
	dr.Dot.Y = fixed.I(y)
	dr.DrawString(drawedText)
}

func drawFrame(input CardInfo) (image.Image, bool) {
	log.Println("drawFrame")

	imageWidth := 1200
	imageHeight := 675

	m := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight)) // 16:9 のpng画像を生成
	c := color.RGBA{50, 200, 255, 255}                            // RGBA で色を指定(B が 255 なので青?)
	c2 := color.RGBA{255, 255, 255, 255}                          // RGBA で色を指定(B が 255 なので青?)

	draw.Draw(m, m.Bounds(), &image.Uniform{c}, image.ZP, draw.Src)               // 青い画像を描画
	rct := image.Rectangle{image.Point{25, 25}, image.Point{1200 - 25, 675 - 25}} // test.jpg をのせる位置を指定する(中央に配置する為に横:25 縦:25 の位置を指定)
	draw.Draw(m, rct, &image.Uniform{c2}, image.Point{0, 0}, draw.Src)            // 合成する画像を描画

	//************************************************
	// https://qiita.com/n-noguchi/items/566e83c5cc0d3b80852c
	// フォントファイルを読み込み
	ftBinary, err := ioutil.ReadFile("ipaexm.ttf")
	// ftBinary, err := ioutil.ReadFile("SourceHanSerifJP-Medium.otf")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ft, err := truetype.Parse(ftBinary)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	opt := truetype.Options{
		Size:              80,
		DPI:               0,
		Hinting:           0,
		GlyphCacheEntries: 0,
		SubPixelsX:        0,
		SubPixelsY:        0,
	}

	face := truetype.NewFace(ft, &opt)

	dr := &font.Drawer{
		Dst:  m,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.Point26_6{},
	}

	textTopMargin := 110
	text2Margin := textTopMargin + (int(opt.Size)+20)*1
	text3Margin := textTopMargin + (int(opt.Size)+20)*2
	text4Margin := textTopMargin + (int(opt.Size)+20)*3
	text5Margin := textTopMargin + (int(opt.Size)+20)*4
	text6Margin := textTopMargin + (int(opt.Size)+20)*5
	textLeftMargin := 40
	textRightMargin := 1150

	// 座標情報を事前に保持してそれをロードしてくる形で実装したい。
	drawTextOnImage(dr, input.Hashtag, textLeftMargin, textTopMargin)
	drawTextOnImage(dr, strconv.Itoa(input.Age), textRightMargin-dr.MeasureString(strconv.Itoa(input.Age)).Ceil(), textTopMargin)
	drawTextOnImage(dr, input.Position, textLeftMargin, text2Margin)
	drawTextOnImage(dr, input.Work, textRightMargin-(dr.MeasureString(strconv.Itoa(input.Sex)).Ceil()+dr.MeasureString(input.Work).Ceil()), text2Margin)
	drawTextOnImage(dr, strconv.Itoa(input.Sex), textRightMargin-dr.MeasureString(strconv.Itoa(input.Sex)).Ceil(), text2Margin)
	drawTextOnImage(dr, input.Background1, textLeftMargin, text3Margin)
	drawTextOnImage(dr, input.Background2, textLeftMargin, text4Margin)
	drawTextOnImage(dr, input.Background3, textLeftMargin, text5Margin)
	drawTextOnImage(dr, input.TwitterId, (imageWidth-dr.MeasureString(input.TwitterId).Ceil())/2, text6Margin)

	//************************************************

	// gocv.PutText(&atom, timeStr, image.Pt(20, atom.Rows()-40), gocv.FontHersheyComplex, 1, black, 1)
	// jpeg.Encode(fso, m, &jpeg.Options{Quality: 100})

	return m, true
}

func ViewHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// func ViewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Del("Content-Type")
	w.WriteHeader(200)

	page := Page{"Hello", ""}
	tmpl, err := template.ParseFiles("html/index.html") // ParseFilesを使う
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, page)
	if err != nil {
		panic(err)
	}
}

func Build() *httprouter.Router {
	router := httprouter.New()

	router.GET("/", ViewHandler)
	router.POST("/createTemplate", CreateTemplate)
	router.POST("/createImage", CreateImage)

	router.NotFound = http.FileServer(http.Dir("html/index.html"))
	router.MethodNotAllowed = http.FileServer(http.Dir("html/index.html"))
	// router.NotFound = http.HandlerFunc(handler.ApiNotFound)
	// router.MethodNotAllowed = http.HandlerFunc(handler.ApiNotFound)

	return router
}

func main() {
	r := Build()
	log.Fatal(http.ListenAndServe(":8080", r))
}
