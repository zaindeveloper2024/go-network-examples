package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func testEndpoint(client *http.Client, method, url string) {
	fmt.Printf("\n=== Testing %s %s ===\n", method, url)

	// リクエストの作成
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Printf("Error creating request: %v\n", err)
		return
	}

	// カスタムヘッダーの追加
	req.Header.Add("User-Agent", "ProxyTestClient/1.0")
	req.Header.Add("X-Custom-Header", "test-value")

	// リクエストの実行と時間計測
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	duration := time.Since(start)

	// レスポンスの詳細表示
	fmt.Printf("Time taken: %v\n", duration)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Headers:\n")
	for k, v := range resp.Header {
		fmt.Printf("  %s: %v\n", k, v)
	}

	// ボディの読み取り
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading body: %v\n", err)
		return
	}

	fmt.Printf("Body preview (first 200 bytes):\n%s\n", string(body[:min(len(body), 200)]))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// プロキシのURL設定
	proxyURL, err := url.Parse("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	// クライアントの設定
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 30 * time.Second,
	}

	// 様々なエンドポイントでテスト
	testURLs := []string{
		"http://example.com",
		"http://httpbin.org/get",
		"http://httpbin.org/headers",
		"http://httpbin.org/ip",
	}

	// 各URLに対してGETリクエストを実行
	for _, testURL := range testURLs {
		testEndpoint(client, "GET", testURL)
	}
}
