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
	"sync"
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

	// WaitGroupを宣言する
	// WaitGroupを使用することで、すべてのgoroutineが終了した時点で次の処理に移ることができる
	wg := new(sync.WaitGroup)

	// channelの生成
	// channel生成時にバッファの数を指定する
	chn := make(chan HackerNews, len(ids))

	count := 0

	for _, v := range ids {
		if count > n-1 {
			break
		}

		// 終了待ちするgoroutineの数を設定する
		wg.Add(1)

		url := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%s.json?print=pretty", strconv.Itoa(v))

		var hn HackerNews

		// goroutineとしてリクエストを行う処理を呼び出す
		go func(url string) {
			res, _ := http.Get(url)
			body, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(body, &hn)

			if hn.Url != "" {
				description := getDescription(hn.Url)

				if description != "" {
					hn.Description = getDescription(hn.Url)
				}
			}

			// チャネルに書き込みを行う
			chn <- hn

			// WaitGroupに設定された数だけDoneが実行されるまで待機
			defer wg.Done()
		}(url)

		hns = append(hns, <-chn)

		count += 1
	}

	// 【重要】データを送信し終わったらチャネルをcloseする
	close(chn)

	// 全ての処理が終わるまで待機
	wg.Wait()
	fmt.Println(hns)

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
