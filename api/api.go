package api

import (
	"encoding/json"
	"fmt"
	"github.com/dyatlov/go-opengraph/opengraph"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Marshalする際に構造体タグを利用することでMarshalの結果を一部制御できる
type HackerNews struct {
	By          string `json:"by"`
	Score       int    `json:"score"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	Url         string `json:"url"`
	Description string
}

func GetHackerNews(n int) []HackerNews {
	resp, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json?print=pretty")
	if err != nil {
		log.Fatal(err)
	}

	// ioutil・・・データの読み書きに必要な機能をまとめたパッケージ
	// ReadAll・・・内容を全て読み込んでバイトスライスとして返却
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// この時点では、bodyはバイト列になっている
	// log.Print(body)

	var ids []int

	// jsonを構造体に変換
	json.Unmarshal(body, &ids)

	var hns []HackerNews
	var hn HackerNews
	count := 0

	for _, v := range ids {
		if count > n-1 {
			break
		}

		url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%s.json?print=pretty", strconv.Itoa(v))
		res, _ := http.Get(url)
		body, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(body, &hn)

		if hn.Url != "" {
			description := getDescription(hn.Url)

			if description != "" {
				hn.Description = getDescription(hn.Url)
			}
		}

		hns = append(hns, hn)

		count += 1
	}

	// hacker newsの情報を含む構造体を返却する
	return hns
}

// og:descriptionを取得
// MEMO: og:descriptionしか取得できない
func getDescription(url string) string {
	res, _ := http.Get(url)
	body, _ := ioutil.ReadAll(res.Body)

	og := opengraph.NewOpenGraph()
	err := og.ProcessHTML(strings.NewReader(string(body)))
	if err != nil {
		log.Fatal(err)
	}

	return og.Description
}
