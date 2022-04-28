package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type Customer struct {
	bun.BaseModel `bun:"table:customers"`

	Token  string `bun:"token,notnull"`
	Unix   int    `bun:"unix_timestamp,notnull"`
	Amount int    `bun:"amount,notnull"`
	Userid string `bun:"customer_id,pk"`
}

type Table struct {
	Db *bun.DB
}

var DbClient *bun.DB

func Open() {
	sqldb, err := sql.Open("postgres", "postgres://postgres:horsemanshoe@localhost:5432/akiya?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}

	db := bun.NewDB(sqldb, pgdialect.New())

	DbClient = db
}

func CreateCustomer(userid string, usertoken string) {

	ctx := context.Background()

	customer := &Customer{Token: usertoken, Userid: userid}
	_, err := DbClient.NewInsert().Model(customer).Exec(ctx)

	if err != nil {
		log.Panic("err: ", err)
	}
}

func CreateOrder(quantity int, userid string) {
	ctx := context.Background()

	order := &Customer{Userid: userid, Amount: quantity}
	orderTime := &Customer{Unix: int(time.Now().Unix())}
	_, err := DbClient.NewUpdate().Model(order).Column("amount").WherePK().Exec(ctx)
	DbClient.NewUpdate().Model(orderTime).Column("unix_timestamp").WherePK().Exec(ctx)

	if err != nil {
		log.Panic("err: ", err)
	}

	println("successfully updated")
}
