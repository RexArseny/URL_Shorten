package shorten

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"os"

	_ "github.com/lib/pq"

	__ "URL_Shorten/proto"
)

type urls struct {
	full_url  string
	short_url string
}

type GRPCServer struct{}

func (s *GRPCServer) Create(ctx context.Context, req *__.Request) (*__.Response, error) {
	password := os.Getenv("DB_PASS")
	connStr := "user=postgres password=" + password + " dbname=urlsdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	rows, err := db.Query("select * from urls")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	url := []urls{}

	for rows.Next() {
		p := urls{}
		err := rows.Scan(&p.full_url, &p.short_url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		url = append(url, p)
	}

	for _, p := range url {
		if p.full_url == req.FullUrl {
			return &__.Response{ShortUrl: p.short_url}, nil
		}
	}

	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_"
	var res_url string
	var equal bool
	for {
		temp := make([]byte, 10)
		for i := range temp {
			temp[i] = letters[rand.Intn(len(letters))]
		}
		res_url = string(temp)
		for _, p := range url {
			if p.short_url == res_url {
				equal = true
				break
			} else {
				equal = false
			}
		}
		if !equal {
			break
		}
	}

	_, err = db.Exec("insert into urls (full_url, short_url) values ($1, $2)", req.FullUrl, res_url)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	return &__.Response{ShortUrl: res_url}, nil
}

func (s *GRPCServer) Get(ctx context.Context, req *__.Response) (*__.Request, error) {
	password := os.Getenv("DB_PASS")
	connStr := "user=postgres password=" + password + " dbname=urlsdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	rows, err := db.Query("select * from urls")

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	url := []urls{}

	for rows.Next() {
		p := urls{}
		err := rows.Scan(&p.full_url, &p.short_url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		url = append(url, p)
	}

	for _, p := range url {
		if p.short_url == req.ShortUrl {
			return &__.Request{FullUrl: p.full_url}, nil
		}
	}

	return &__.Request{FullUrl: "none"}, nil
}
