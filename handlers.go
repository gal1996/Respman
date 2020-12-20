package main

import (
	"bufio"
	"fmt"
	"github.com/labstack/echo"
	"github.com/thoas/go-funk"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	FileName = "response.json"
)

func Default(c echo.Context) error {
	return c.Render(http.StatusOK, "res", "hello Respman")
}

// サーバーサイドがレスポンスのjsonが正しい形か確認するよう
func CheckResponse(c echo.Context) error {
	f, err := os.Open(FileName)
	if err != nil {
		fmt.Printf("[Fail] failed file open: %s\n", err.Error())
		return err
	}

	// ファイルから期待するjsonのkeyの抽出
	expectedKeys := extractJsonKeys(f)
	fmt.Printf("[Info] expectedKeys(len: %d): %#v\n", len(expectedKeys), expectedKeys)

	// リクエストからきたjsonのkeyの抽出
	actualKeys := extractJsonKeys(c.Request().Body)
	fmt.Printf("[Info] acctualKeys(len: %d): %#v\n", len(actualKeys), actualKeys)

	// そもそも二つの長さが違うのはアウト
	if len(expectedKeys) != len(actualKeys) {
		return c.String(http.StatusOK, "length is not equal")
	}

	invalidKey := make([]string, 0, len(actualKeys))
	// 同じkeyで複数データが存在する場合、長さが同じでも存在しないkeyがすり抜けてしまうので、actualとexpected両方でチェックする必要がある
	for _, actualKey := range actualKeys {
		if !funk.ContainsString(expectedKeys, actualKey) {
			actualKey = strings.TrimSpace(actualKey)
			data := fmt.Sprintf("%s(is invalid)", actualKey)
			invalidKey = append(invalidKey, data)
		}
	}

	for _, expectedKey := range expectedKeys {
		if !funk.ContainsString(actualKeys, expectedKey) {
			expectedKey = strings.TrimSpace(expectedKey)
			data := fmt.Sprintf("%s(is not in your response)", expectedKey)
			invalidKey = append(invalidKey, data)
		}
	}

	if len(invalidKey) != 0 {
		resStr := fmt.Sprintf("your response is not perfect invalid keys: %s", invalidKey)
		return c.String(http.StatusOK, resStr)
	}

	return c.String(http.StatusOK, "your response is perfect")
}

func extractJsonKeys(rowJson io.Reader) []string {
	jsonKeys := make([]string, 0, 30)
	scanner := bufio.NewScanner(rowJson)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, ":") {
			// 基本的にjsonは 'key: value' の形なので、':'の前がキーになる
			jsonKey := strings.TrimSpace(strings.Replace((strings.Split(text, ":"))[0], "\"", "", -1))
			jsonKeys = append(jsonKeys, jsonKey)
		}
	}

	return jsonKeys
}

// 指定したレスポンスを返すやつ
func ReturnResponse(c echo.Context) error {
	log := c.Logger()
	log.Infof("[Start] ReturnResponse\n")
	defer log.Infof("[End] ReturnResponse\n")

	f, err := os.Open(FileName)
	if err != nil {
		log.Fatalf("[Fail] failed load response.txt: %s\n", err.Error())
		return err
	}
	defer f.Close()

	buf := make([]byte, 2048)
	var response []byte
	for {
		n, err := f.Read(buf)
		// バイト数が0は読み取り終了
		if n == 0 {
			break
		}
		if err != nil {
			break
		}

		response = buf[:n]
	}

	if err := c.JSONBlob(http.StatusOK, response); err != nil {
		log.Fatalf("[Fail] failed jsonBlob: %s\n", err.Error())
		return err
	}

	return nil
}
