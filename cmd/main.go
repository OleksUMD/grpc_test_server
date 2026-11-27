package main

import (
	"database/sql"
	"log"

	dbMigrate "github.com/OleksUMD/grpc_test_server/db"
	mydb "github.com/OleksUMD/grpc_test_server/internal/adapter/database"
	mygrpc "github.com/OleksUMD/grpc_test_server/internal/adapter/grpc"
	app "github.com/OleksUMD/grpc_test_server/internal/application"
	bank_app "github.com/OleksUMD/grpc_test_server/internal/application/bank"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(logWriter{})
	sqlDB, err := sql.Open("pgx", "postgres://bank_user:Bank_user123!@localhost:5443/bank_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	dbMigrate.Migrate(sqlDB)
	databaseAdapter, err := mydb.NewDatabaseAdapter(sqlDB)
	hs := &app.HelloService{}
	rs := &app.ResiliencyService{}
	bs := bank_app.NewBankService(databaseAdapter)
	grpcAdapter := mygrpc.NewGrpcAdapter(hs, bs, rs, 8787)
	grpcAdapter.Run()
}
