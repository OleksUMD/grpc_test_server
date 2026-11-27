package bank

import (
	"time"

	"github.com/OleksUMD/grpc_test_server/internal/adapter/database"
	"github.com/OleksUMD/grpc_test_server/internal/port"
	"github.com/google/uuid"
)

type BankService struct {
	dbAdapter port.BankDatabasePort
}

func NewBankService(db port.BankDatabasePort) *BankService {
	return &BankService{dbAdapter: db}
}

func (a *BankService) GetCurrentBalance(accountNumber string) (float64, error) {
	account, err := a.dbAdapter.GetAccountByNumber(accountNumber)
	if err != nil {
		return 0, err
	}
	return account.CurrentBalance, nil
}

func (a *BankService) GetCurrentExchangeRates(fromCurrency, toCurrency string, ts time.Time) (float64, error) {
	excRateInst, err := a.dbAdapter.GetCurrentExchangeRates(fromCurrency, toCurrency, ts)
	if err != nil {
		return 0, err
	}
	return excRateInst.Rate, nil
}

func (a *BankService) GetEmptyTransactionSummary() *database.TransactionSummary {
	return &database.TransactionSummary{
		SummaryOnDate: time.Now(),
		SumIn:         0,
		SumOut:        0,
		SumTotal:      0,
	}
}

func (a *BankService) CreateTransaction(accountNumber string, trans *database.BankTransaction) (uuid.UUID, error) {
	accnt, err := a.dbAdapter.GetAccountByNumber(accountNumber)
	if err != nil {
		return uuid.Nil, nil
	}
	trans.AccountUUID = accnt.AccountUUID
	return a.dbAdapter.CreateTransaction(accnt, trans)
}

func (a *BankService) SummarizeTransactions(curSum *database.TransactionSummary, trans *database.BankTransaction) error {
	return a.dbAdapter.CalculateTransactionSummary(curSum, trans)
}

func (a *BankService) MultipleMoneyTransfer(fromAccntNum, toAccntNum, currency string, amount float64, ts time.Time) error {
	return a.dbAdapter.CreateTransfer(fromAccntNum, toAccntNum, currency, amount, ts)
}
