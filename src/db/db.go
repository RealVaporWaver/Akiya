package db

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type Customer struct {
	bun.BaseModel `bun:"table:customers"`

	Token  string `bun:"token,notnull"`
	Unix   int    `bun:"unix_timestamp,notnull"`
	Amount int    `bun:"amount,notnull"`
}

func Open() *bun.DB {
	sqldb, err := sql.Open("postgres", "postgres://postgres:horsemanshoe@localhost:5432/akiya?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}

	db := bun.NewDB(sqldb, pgdialect.New())

	return db
}

func Insert(db *bun.DB) {

	ctx := context.Background()

	customer := &Customer{Token: "aaaaaaaaaaaaaaaa", Unix: 9000000, Amount: 34}
	_, err := db.NewInsert().Model(customer).Exec(ctx)

	if err != nil {
		log.Panic("err: ", err)
	}

	println("Successfully Added Customer")
}
