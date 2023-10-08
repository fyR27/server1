package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

var per = `{
	"order_uid": "b563feb7b2b84b6test",
	"track_number": "WBILMTESTTRACK",
	"entry": "WBIL",
	"delivery": {
	  "name": "Test Testov",
	  "phone": "+9720000000",
	  "zip": "2639809",
	  "city": "Kiryat Mozkin",
	  "address": "Ploshad Mira 15",
	  "region": "Kraiot",
	  "email": "test@gmail.com"
	},
	"payment": {
	  "transaction": "b563feb7b2b84b6test",
	  "request_id": "",
	  "currency": "USD",
	  "provider": "wbpay",
	  "amount": 1817,
	  "payment_dt": 1637907727,
	  "bank": "alpha",
	  "delivery_cost": 1500,
	  "goods_total": 317,
	  "custom_fee": 0
	},
	"items": [
	  {
		"chrt_id": 9934930,
		"track_number": "WBILMTESTTRACK",
		"price": 453,
		"rid": "ab4219087a764ae0btest",
		"name": "Mascaras",
		"sale": 30,
		"size": "0",
		"total_price": 317,
		"nm_id": 2389212,
		"brand": "Vivienne Sabo",
		"status": 202
	  }
	],
	"locale": "en",
	"internal_signature": "",
	"customer_id": "test",
	"delivery_service": "meest",
	"shardkey": "9",
	"sm_id": 99,
	"date_created": "2021-11-26T06:22:19Z",
	"oof_shard": "1"
	}`

type person struct {
	OrderUID    string `json:"order_uid"`
	TrackNumber string `json:"track_number"`
	Entry       string `json:"entry"`
	Delivery    struct {
		Name    string `json:"name"`
		Phone   string `json:"phone"`
		Zip     string `json:"zip"`
		City    string `json:"city"`
		Address string `json:"address"`
		Region  string `json:"region"`
		Email   string `json:"email"`
	} `json:"delivery"`
	Payment struct {
		Transaction  string `json:"transaction"`
		RequestID    string `json:"request_id"`
		Currency     string `json:"currency"`
		Provider     string `json:"provider"`
		Amount       int    `json:"amount"`
		PaymentDt    int    `json:"payment_dt"`
		Bank         string `json:"bank"`
		DeliveryCost int    `json:"delivery_cost"`
		GoodsTotal   int    `json:"goods_total"`
		CustomFee    int    `json:"custom_fee"`
	} `json:"payment"`
	Items []struct {
		ChrtID      int    `json:"chrt_id"`
		TrackNumber string `json:"track_number"`
		Price       int    `json:"price"`
		Rid         string `json:"rid"`
		Name        string `json:"name"`
		Sale        int    `json:"sale"`
		Size        string `json:"size"`
		TotalPrice  int64  `json:"total_price"`
		NmID        int    `json:"nm_id"`
		Brand       string `json:"brand"`
		Status      int64  `json:"status"`
	} `json:"items"`
	Locale            string `json:"locale"`
	InternalSignature string `json:"internal_signature"`
	CustomerID        string `json:"customer_id"`
	DeliveryService   string `json:"delivery_service"`
	Shardkey          string `json:"shardkey"`
	SmID              int    `json:"sm_id"`
	DateCreated       string `json:"date_created"`
	OofShard          string `json:"oof_shard"`
}

const (
	host     = "localhost"
	port     = 1111
	user     = "postgres"
	password = "1111"
	dbname   = "test"
)

func sayhello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, per)
}

func main() {
	var order1 map[string]interface{}

	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)

	defer c.Close()

	c.Publish("test", "WB")

	if err := json.Unmarshal([]byte(per), &order1); err != nil {
		panic(err)
	}

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	insertStmt := `insert into "head"("order_uid", "track_number", "entry", "locale", "internal_signature", "customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard") 
	values('b563feb7b2b84b6test', 'WBILMTESTTRACK', 'WBIL', 'en', '', 'test', 'meest', '9', '99', '20012', '1')`

	_, e := db.Exec(insertStmt)
	if e != nil {
		panic(e)
	}

	insertStmt = `insert into "delivery"("order_uid", "name", "phone", "zip", "city", "address", "region", "email") 
	values('b563feb7b2b84b6test', 'Test Testov', '+9720000000', '2639809', 'Kiryat Mozkin', 'Ploshad Mira 15', 'Kraiot', 'test@gmail.com')`

	_, err = db.Exec(insertStmt)
	if e != nil {
		panic(err)
	}

	insertStmt = `insert into "payment"("order_uid", "transaction": "b563feb7b2b84b6test", "request_id", "currency", "provider", "amount", "payment_dt", "bank", "delivery_cost", "goods_total", "custom_fee") 
	values('b563feb7b2b84b6test', '111', '', 'USD', 'wbpay', '1817', '1637907727', 'alpha', '1500', '317', '0')`

	_, err = db.Exec(insertStmt)
	if e != nil {
		panic(err)
	}

	insertStmt = `insert into "items"("order_uid", "chrt_id", "track_number", "price", "rid", "name", "sale", "size", "total_price", "nm_id", "brand", "status") 
	values('b563feb7b2b84b6test', '9934930', 'WBILMTESTTRACK', '453', 'ab4219087a764ae0btest', 'Mascaras', '30', '0', '317', '2389212' , 'Vivienne Sabo' , '202')`

	_, err = db.Exec(insertStmt)
	if e != nil {
		panic(err)
	}

	http.HandleFunc("/", sayhello)
	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
