package port

import (
	"time"

	"github.com/OleksUMD/grpc_test_server/internal/adapter/database"
	"github.com/google/uuid"
)

type HelloServicePort interface {
	GenerateHello(name string) string
}

type BankServicePort interface {
	GetCurrentBalance(accountNumber string) (float64, error)
	GetCurrentExchangeRates(fromCurrency, toCurrency string, ts time.Time) (float64, error)
	GetEmptyTransactionSummary() *database.TransactionSummary
	CreateTransaction(accountNumber string, trans *database.BankTransaction) (uuid.UUID, error)
	SummarizeTransactions(curSum *database.TransactionSummary, trans *database.BankTransaction) error
	MultipleMoneyTransfer(fromAccntNum, toAccntNum, currency string, amount float64, ts time.Time) error
}

type ResiliencyServicePort interface {
	GenerateResiliency(minDelaySecond, maxDelaySecond int32, statusCodes []uint32) (string, uint32)
}
