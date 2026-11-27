package database

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	TransactionTypeUnknown string = "UNKNOWN"
	TransactionTypeIn      string = "IN"
	TransactionTypeOut     string = "OUT"
)

var ErrTransferSourceAccountNotFound = errors.New("source account not found")
var ErrTransferDestinationAccountNotFound = errors.New("destination account not found")

type BankAccount struct {
	AccountUUID    uuid.UUID `gorm:"primaryKey"`
	AccountNumber  string
	AccountName    string
	Currency       string
	CurrentBalance float64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Transactions   []BankTransaction `gorm:"foreignKey:AccountUUID"`
}

type BankTransaction struct {
	TransactionUUID      uuid.UUID `gorm:"primaryKey"`
	AccountUUID          uuid.UUID
	TransactionTimestamp time.Time
	Amount               float64
	TransactionType      string
	Notes                string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type BankExchangeRate struct {
	ExchangeRateUUID   uuid.UUID `gorm:"primaryKey"`
	FromCurrency       string
	ToCurrency         string
	ValidFromTimestamp time.Time
	ValidToTimestamp   time.Time
	Rate               float64
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type TransactionSummary struct {
	SummaryOnDate time.Time
	SumIn         float64
	SumOut        float64
	SumTotal      float64
}

type BankTransfer struct {
	TransferUUID      uuid.UUID `gorm:"primaryKey"`
	FromAccountUUID   uuid.UUID
	ToAccountUUID     uuid.UUID
	Currency          string
	Amount            float64
	TransferTimestamp time.Time
	TransferSuccess   bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
