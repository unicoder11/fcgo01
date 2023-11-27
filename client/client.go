package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {

	file, err := os.OpenFile("cotacao.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*30000)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/usdbrl", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	usdbrl, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	usdstring := string(usdbrl)
	_, err2 := file.WriteString("Dolar :" + usdstring + "\n")
	if err2 != nil {
		panic(err2)
	}
	fmt.Println("wrote to file")

}
