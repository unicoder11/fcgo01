package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	http.HandleFunc("/usdbrl", handler)
	http.ListenAndServe(":8080", nil)
}
func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*300)
	defer cancel()

	result, err := getUsdbrl(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%q", result)
}

func insertUsdbrl(usdbrl string) error {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*10)
	defer cancel()

	db, err := sql.Open("sqlite3", "./usdbrl.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	select {
	case <-ctx.Done():
		fmt.Println("timeout db")
	default:
		datetime := time.Now().Format("2006-01-02 15:04:05")
		stmt, err := db.Prepare("INSERT INTO usdbrl VALUES(?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(datetime, usdbrl)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("inserted into database")
	}
	return nil
}

func getUsdbrl(ctx context.Context) (string, error) {
	var data USDBRL

	select {
	case <-ctx.Done():
		fmt.Println("timeout apicall")
	default:
		resp, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Fatal(err)
		}
		fmt.Println(data.Usdbrl.Bid)
		insertUsdbrl(data.Usdbrl.Bid)
		return data.Usdbrl.Bid, nil
	}
	return data.Usdbrl.Bid, nil
}

type USDBRL struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}
