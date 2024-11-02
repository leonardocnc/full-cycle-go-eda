package main

import (
	"balances/internal/database"
	"balances/internal/usecase/balance/create_balance"
	"balances/internal/usecase/balance/get_balance"
	"balances/internal/web"
	"balances/internal/web/webserver"
	"encoding/json"

	// "balances/pkg/kafka"
	"balances/pkg/uow"
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
)

type BalanceEvent struct {
	Name    string      `json:"Name"`
	Payload BalanceData `json:"Payload"`
}

type BalanceData struct {
	AccountIDFrom      string  `json:"account_id_from"`
	AccountIDTo        string  `json:"account_id_to"`
	BalanceAccountFrom float64 `json:"balance_account_id_from"`
	BalanceAccountTo   float64 `json:"balance_account_id_to"`
}

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "balance-mysql", "3306", "balance"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Exec("CREATE DATABASE IF NOT EXISTS balance;")
	db.Exec("DROP TABLE IF EXISTS balances;")
	db.Exec("Create table balances (id varchar(255), account_id varchar(255), amount int, created_at timestamp);")
	db.Exec("INSERT INTO balances(id, account_id, amount, created_at) VALUES('254ec385-e4ca-4dc8-ba7f-478baacad022', '5c4e66ba-2042-4248-9df5-3cba4c72b3c3', 100, NOW());")
	db.Exec("INSERT INTO balances(id, account_id, amount, created_at) VALUES('254ec385-e4ca-4dc8-ba7f-478baacad021', '5c4e66ba-2042-4248-9df5-3cba4c72b3c5', 100, NOW());")

	balanceDb := database.NewBalanceDB(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("BalanceDB", func(tx *sql.Tx) interface{} {
		return database.NewBalanceDB(db)
	})

	createBalanceUseCase := create_balance.NewCreateBalanceUseCase(uow)
	getBalanceUseCase := get_balance.NewGetBalanceUseCase(balanceDb)

	go func() {
		webserver := webserver.NewWebServer(":3003")

		balanceHandler := web.NewWebBalanceHandler(
			*createBalanceUseCase,
			*getBalanceUseCase,
		)

		webserver.AddHandler("/balances", balanceHandler.CreateBalance)

		webserver.AddGetHandler("/balances/{account_id}/account", balanceHandler.GetBalance)
		webserver.AddGetHandler("/health", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Health check")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok"}`))
		})

		fmt.Println("Server is running")
		webserver.Start()
	}()

	configMap := &kafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"client.id":         "balances",
		"group.id":          "balances",
		"auto.offset.reset": "earliest",
	}

	c, err := kafka.NewConsumer(configMap)
	if c != nil {
		topics := []string{"balances"}
		c.SubscribeTopics(topics, nil)

		for {
			msg, err := c.ReadMessage(-1)
			if msg != nil || err == nil {
				fmt.Println(string(msg.Value), msg.TopicPartition)

				var event BalanceEvent
				if err := json.Unmarshal(msg.Value, &event); err != nil {
					fmt.Println("Error JSON:", err)
					continue
				}

				input := create_balance.CreateBalanceInputDTO{
					AccountID: event.Payload.AccountIDFrom,
					Amount:    int(event.Payload.BalanceAccountFrom),
				}
				createBalanceUseCase.Execute(ctx, input)

				input = create_balance.CreateBalanceInputDTO{
					AccountID: event.Payload.AccountIDTo,
					Amount:    int(event.Payload.BalanceAccountTo),
				}
				createBalanceUseCase.Execute(ctx, input)
				c.CommitMessage(msg)
			}
		}
	}
}
