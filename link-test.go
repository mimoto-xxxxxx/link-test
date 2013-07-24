package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"time"
)

type TestResultSet struct {
	Title    string
	Url      string
	Selector string
	Results  []*TestResult
	Err      error
	StartAt  time.Time
	EndAt    time.Time
}

type TestResult struct {
	Url        string
	Caption    string
	StatusCode int
	Err        error
	StartAt    time.Time
	EndAt      time.Time
}

const DATETIME_FORMAT = "2006-01-02 15:04:05"

func findUrl(e *goquery.Selection) (string, bool) {
	u, ok := e.Attr("href")
	if ok {
		return u, ok
	}

	u, ok = e.Attr("src")
	if ok {
		return u, ok
	}

	return "", false
}

func TestLink(title string, url string, selector string) (r TestResultSet) {
	r = TestResultSet{
		Title:    title,
		Url:      url,
		Selector: selector,
		StartAt:  time.Now(),
	}
	defer func() {
		r.EndAt = time.Now()
	}()

	doc, err := goquery.NewDocument(url)
	if err != nil {
		r.Err = err
		return
	}

	if title == "" {
		r.Title = doc.Find("title").Text()
	}

	doc.Find(selector).Each(func(index int, e *goquery.Selection) {
		u, ok := findUrl(e)
		if !ok {
			//href や src に URL が存在しなかった
			return
		}

		tr := TestResult{
			Url:     u,
			Caption: e.Text(),
			StartAt: time.Now(),
		}
		defer func() {
			tr.EndAt = time.Now()
		}()

		resp, err := http.Get(u)
		if err != nil {
			tr.Err = err
			r.Results = append(r.Results, &tr)
			return
		}
		defer resp.Body.Close()

		tr.StatusCode = resp.StatusCode
		r.Results = append(r.Results, &tr)
	})
	return
}

func PrintTestResultSet(r TestResultSet) {
	fmt.Println(r.Title)
	fmt.Println("========")
	fmt.Println()
	fmt.Printf("調査日時: %s\n", r.StartAt.Format(DATETIME_FORMAT))
	fmt.Printf("全所要時間: %0.3f秒\n", r.EndAt.Sub(r.StartAt).Seconds())
	fmt.Println()
	fmt.Printf("ページ `%s` 内で  \n", r.Url)
	fmt.Printf("セレクタ `%s` に一致するリンクを生存チェックしました。\n", r.Selector)
	fmt.Println()
	if r.Err != nil {
		fmt.Printf("`%s` へのアクセス時にエラーが発生しました\n", r.Url)
		fmt.Printf("エラー: %v\n", r.Err)
		return
	}

	if len(r.Results) == 0 {
		fmt.Println("エラー: 生存チェック対象の URL がひとつも見つかりませんでした。")
		return
	}

	fmt.Printf("全 %d 件\n", len(r.Results))
	fmt.Println()

	for i, result := range r.Results {
		fmt.Printf("%d. %s\n", i+1, result.Caption)
		fmt.Println("--------")
		fmt.Printf(" * URL: %s\n", result.Url)
		if result.Err != nil {
			fmt.Printf(" * エラー: %v\n", result.Err)
		} else {
			fmt.Printf(" * ステータス: %d\n", result.StatusCode)
		}
		fmt.Printf(" * 所要時間: %0.3f秒\n", result.EndAt.Sub(result.StartAt).Seconds())
		fmt.Println()
	}
}

var (
	title    = flag.String("t", "サンプルページ", "レポートの見出し")
	url      = flag.String("u", "", "リンク掲載ページの URL")
	selector = flag.String("s", "", "リンク検出用のCSSセレクタ")
)

func main() {
	flag.Parse()
	if *url == "" || *selector == "" {
		flag.Usage()
		return
	}

	ret := TestLink(*title, *url, *selector)
	PrintTestResultSet(ret)
}
