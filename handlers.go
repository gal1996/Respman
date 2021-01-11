package main

import (
	"github.com/labstack/echo"
	"net/http"
	"os"
	"strings"
)

func Default(c echo.Context) error {
	return c.Render(http.StatusOK, "res", "hello Respman")
}

// 指定したレスポンスを返すやつ
func ReturnResponse(c echo.Context) error {
	log := c.Logger()
	log.Infof("[Start] ReturnResponse\n")
	defer log.Infof("[End] ReturnResponse\n")

	// どのエンドポイントに来たリクエストかを調べる
	path := strings.Trim(c.Path(), "/")

	fileName := path + ".json"

	f, err := os.Open(fileName)
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
