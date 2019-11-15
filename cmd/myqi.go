package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const baseURL = "https://qiita.com/api/v2/items"

// jsonをパースする為の構造体を定義する

type Data struct {
	ID             string `json:"id"`
	Url            string `json:"url"`
	Title          string `json:"title"`
	LikesCount     int    `json:"likes_count"`
	PageViewsCount int    `json:"page_views_count"`
}

func FetchMyQiitaData(accessToken string, qiitaUser string) ([]Data, error) {
	b, err := json.Marshal(Data{})
	if err != nil {
		return nil, err
	}

	var resp *http.Response
	// qiitaのアクセストークンがない場合はAuthorizationを付与しない
	// 2パターン作っておく。
	// accessトークンは環境変数に入れておく。自分の場合は.bash_profileにexport文を書いている。

	req, err := http.NewRequest(http.MethodGet, baseURL+"?query=user:"+qiitaUser, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", "application/json")

	if len(accessToken) > 0 {
		fmt.Println("***** Access Token 無しでQiitaAPIを叩いています アクセス制限に注意して下さい*****")
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	// QiitaAPIにリクエストを送ってレスポンスをrespに格納する。
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data []Data

	if err := json.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("JSON Unmarshal error: %w", err)
	}

	/*********一覧取得では、ページビューがnilになるので個別で取りに行ってデータを得る*****************/
	for i, val := range data {
		itemsURL := "https://qiita.com/api/v2/items/" + val.ID

		req, err := http.NewRequest(http.MethodGet, itemsURL, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("content-type", "application/json")
		if len(accessToken) > 0 {
			fmt.Println("***** Access Token 無しでQiitaAPIを叩いています アクセス制限に注意して下さい*****")
			req.Header.Set("Authorization", "Bearer "+accessToken)
		}

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		b, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		var m map[string]interface{}

		if err := json.Unmarshal(b, &m); err != nil {
			return nil, fmt.Errorf("JSON Unmarshal error: %w", err)
		}

		if v, ok := m["page_views_count"]; ok {
			if count, ok := v.(float64); ok {
				data[i].PageViewsCount = int(count)
			}
		}
	}
	return data, nil
}

// データの出力
func OutputQiitaData(data []Data) {
	fmt.Println("************************自分のQiita投稿一覧******************************")
	for _, val := range data {
		fmt.Printf("%-15v%v%v\n", "ID", ": ", val.ID)
		fmt.Printf("%-15v%v%v\n", "Title", ": ", val.Title)
		fmt.Printf("%-12v%v%v\n", "いいね", ": ", val.LikesCount)
		fmt.Printf("%-9v%v%v\n", "ページビュー", ": ", val.PageViewsCount)
		fmt.Printf("%-15v%v%v\n", "URL", ": ", val.Url)
		fmt.Println("-------------------------------------------------------------------------")
	}
}
