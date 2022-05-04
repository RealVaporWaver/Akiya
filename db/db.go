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

var DbClient *bun.DB
var Ctx context.Context

func Open() {
	sqldb, err := sql.Open("postgres", "postgres://postgres:horsemanshoe@localhost:5432/akiya?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}

	db := bun.NewDB(sqldb, pgdialect.New())

	DbClient = db
	Ctx = context.Background()
}

func CreateCustomer(userid string, usertoken string) {

	customer := &Customer{Token: usertoken, Userid: userid, Unix: int(time.Now().UTC().Unix())}
	_, err := DbClient.NewInsert().Model(customer).Exec(Ctx)

	if err != nil {
		log.Panic("err: ", err)
	}
}

func CreateOrder(quantity int, userid string) {

	order := &Customer{Userid: userid, Amount: quantity}
	orderTime := &Customer{Unix: int(time.Now().Unix())}
	_, err := DbClient.NewUpdate().Model(order).Column("amount").WherePK().Exec(Ctx)
	DbClient.NewUpdate().Model(orderTime).Column("unix_timestamp").WherePK().Exec(Ctx)

	if err != nil {
		log.Panic("err: ", err)
	}
}

func GetUser(userid string) Customer {
	var customer Customer
	DbClient.NewSelect().Model(&customer).Where("customer_id = ?", userid).Scan(Ctx)

	return customer
}

func UpdateCustomer(userid string, token string) {
	test := &Customer{Userid: userid, Token: token}
	DbClient.NewUpdate().Model(test).Column("token").Exec(Ctx)
}
