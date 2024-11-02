# full-cycle-go-eda

docker compose up -d

docker compose exec wallet bash 
cd wallet
go run cmd/walletcore/main.go

docker compose exec balance bash 
cd balances
go run cmd/balances/main.go
