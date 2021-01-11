package main

import (
	"bufio"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"os"
	"strings"
)

const (
	PathConfigFile = "path.conf"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.MethodOverride())

	e.GET("/", Default)
	urlPaths, err := extractURLPath()
	fmt.Printf("paths: %v", urlPaths)
	if err != nil {
		fmt.Printf("error has occurred in extractURLPath: %s", err.Error())
		return
	}

	for _, path := range urlPaths {
		e.POST(path, ReturnResponse)
	}

	e.Logger.Fatal(e.Start(":3030"))
}

// path.confから一行ずつエンドポイントとして抜き出す
func extractURLPath() ([]string, error) {
	rowJson, err := os.Open(PathConfigFile)
	if err != nil {
		fmt.Printf("[Fail] failed file open: %s\n", err.Error())
		return nil, err
	}

	urlPaths := make([]string, 0, 30)
	scanner := bufio.NewScanner(rowJson)
	for scanner.Scan() {
		text := scanner.Text()
		word := strings.TrimSpace(strings.Replace(text, "\"", "", -1))
		urlPaths = append(urlPaths, word)
	}

	return urlPaths, nil
}
