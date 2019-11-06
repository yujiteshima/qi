package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "qiitasearch"
	app.Usage = "search qiita articles"
	app.Version = "0.1.0"

	app.Commands = []cli.Command{
		{
			Name:  "mine",
			Usage: "qiita + mine : you get yours articles",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "example, e",
					Value: "exmple",
					Usage: "This help text is example",
				},
			},
			Action: func(c *cli.Context) error {
				//fmt.Println("Hello friend!")
				qiitaToken := os.Getenv("QIITA_TOKEN")
				// Debug用 fmt.Println(qiitaToken)
				//id_ary := getId()
				datas := fetchQiitaData(qiitaToken)
				outputQiitaData(datas)
				return nil
			},
		},
	}
	app.Run(os.Args)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

// jsonをパースする為の構造体を定義する

type Data struct {
	ID         string `json:id`
	Url        string `json:"url"`
	Title      string `json:"title"`
	LikesCount int    `json:"likes_count"`
	//ReactionsCount int    `json:"reactions_count"`
	PageViewsCount int `json:"page_views_count"`
}

func fetchQiitaData(accessToken string) []Data {
	baseUrl := "https://qiita.com/api/v2/items?query=user:yujiteshima"
	// 様々な検索条件をかけるときはbaseUrlをv2/までにして他を変数で定義してurl.Parseで合体させる
	endpointURL, err := url.Parse(baseUrl)
	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(Data{})
	if err != nil {
		panic(err)
	}

	var resp = &http.Response{}
	// qiitaのアクセストークンがない場合はAuthorizationを付与しない
	// 2パターン作っておく。
	// accessトークンは環境変数に入れておく。自分の場合は.bash_profileにexport文を書いている。

	if len(accessToken) > 0 {
		resp, err = http.DefaultClient.Do(&http.Request{
			URL:    endpointURL,
			Method: "GET",
			Header: http.Header{
				"Content-Type":  {"application/json"},
				"Authorization": {"Bearer " + accessToken},
			},
		})
	} else {
		fmt.Println("***** Access Token 無しでQiiitaAPIを叩いています アクセス制限に注意して下さい*****")

		resp, err = http.DefaultClient.Do(&http.Request{
			URL:    endpointURL,
			Method: "GET",
			Header: http.Header{
				"Content-Type": {"application/json"},
			},
		})
	}
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var datas []Data

	if err := json.Unmarshal(b, &datas); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return nil
	}

	/*********************************************************************************/
	for i, val := range datas {
		//fmt.Println("id:", val.ID)
		article_id := val.ID
		baseUrl := "https://qiita.com/api/v2/items/"
		endpointURL2, err := url.Parse(baseUrl + article_id)
		if err != nil {
			panic(err)
		}

		b, err := json.Marshal(Data{})
		if err != nil {
			panic(err)
		}

		resp, err = http.DefaultClient.Do(&http.Request{
			URL:    endpointURL2,
			Method: "GET",
			Header: http.Header{
				"Content-Type":  {"application/json"},
				"Authorization": {"Bearer " + accessToken},
			},
		})

		if err != nil {
			panic(err)
		}

		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var m map[string]interface{}

		if err := json.Unmarshal(b, &m); err != nil {
			fmt.Println("JSON Unmarshal error:", err)
			return nil
		}

		//fmt.Println(m["page_views_count"])
		datas[i].PageViewsCount = int(m["page_views_count"].(float64))
	}
	return datas
}

// データの出力
func outputQiitaData(datas []Data) {
	fmt.Println("************************自分のQiita投稿一覧******************************")
	for _, val := range datas {
		fmt.Printf("%-15v%v%v\n", "ID", ": ", val.ID)
		fmt.Printf("%-15v%v%v\n", "Title", ": ", val.Title)
		fmt.Printf("%-12v%v%v\n", "いいね", ": ", val.LikesCount)
		fmt.Printf("%-9v%v%v\n", "ページビュー", ": ", val.PageViewsCount)
		fmt.Printf("%-15v%v%v\n", "URL", ": ", val.Url)
		fmt.Println("-------------------------------------------------------------------------")
	}
}
