package port

import (
	"time"

	db "github.com/OleksUMD/grpc_test_server/internal/adapter/database"
	"github.com/google/uuid"
)

type BankDatabasePort interface {
	SaveAccount(data *db.BankAccount) (uuid.UUID, error)
	GetAccountByUuid(id uuid.UUID) (*db.BankAccount, error)
	GetAccountByNumber(accountNumber string) (*db.BankAccount, error)
	GetCurrentExchangeRates(fromCurrency, toCurrency string, ts time.Time) (*db.BankExchangeRate, error)
	CreateTransaction(accnt *db.BankAccount, t *db.BankTransaction) (uuid.UUID, error)
	CalculateTransactionSummary(tcur *db.TransactionSummary, trans *db.BankTransaction) error
	CreateTransfer(fromAccntNum, toAccntNum, currency string, amount float64, ts time.Time) error
}
