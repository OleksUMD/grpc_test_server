package database

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (a *DatabaseAdapter) SaveAccount(data *BankAccount) (uuid.UUID, error) {
	if err := a.db.Save(data).Error; err != nil {
		return uuid.Nil, err
	}
	return data.AccountUUID, nil
}

func (a *DatabaseAdapter) GetAccountByUuid(id uuid.UUID) (*BankAccount, error) {
	var res BankAccount
	if err := a.db.First(&res, id).Error; err != nil {
		return nil, err
	}
	return &res, nil
}

func (a *DatabaseAdapter) GetAccountByNumber(accountNumber string) (*BankAccount, error) {
	var res BankAccount
	if err := a.db.Where("account_number = ?", accountNumber).First(&res).Error; err != nil {
		return nil, err
	}
	return &res, nil
}

func (a *DatabaseAdapter) GetCurrentExchangeRates(fromCurrency, toCurrency string, ts time.Time) (*BankExchangeRate, error) {
	var res BankExchangeRate
	if err := a.db.First(
		&res,
		"from_currency = ? AND to_currency = ? AND (? BETWEEN valid_from_timestamp AND valid_to_timestamp)",
		fromCurrency, toCurrency, ts).Error; err != nil {
		return nil, err
	}
	return &res, nil
}

func (a *DatabaseAdapter) CreateTransaction(account *BankAccount, t *BankTransaction) (uuid.UUID, error) {
	tx := a.db.Begin()
	defer func() {
		if err := tx.Rollback().Error; err != nil {
			return
		}
	}()
	if err := tx.Create(&t).Error; err != nil {
		return uuid.Nil, err
	}
	newAmount := t.Amount
	if t.TransactionType == TransactionTypeOut {
		newAmount = -1 * t.Amount
	}
	account.CurrentBalance += newAmount
	account.UpdatedAt = time.Now()
	_, err := a.SaveAccount(account)
	if err != nil {
		return uuid.Nil, err
	}
	tx.Commit()
	return t.TransactionUUID, nil
}

func (a *DatabaseAdapter) CalculateTransactionSummary(tcur *TransactionSummary, trans *BankTransaction) error {
	switch trans.TransactionType {
	case TransactionTypeIn:
		tcur.SumIn += trans.Amount
	case TransactionTypeOut:
		tcur.SumOut += trans.Amount
	default:
		return fmt.Errorf("unknown transaction type : %s", trans.TransactionType)
	}
	tcur.SumTotal = tcur.SumIn - tcur.SumOut
	return nil
}

func (a *DatabaseAdapter) CreateTransfer(fromAccntNum, toAccntNum, currency string, amount float64, ts time.Time) error {
	tx := a.db.Begin()
	defer func() {
		if err := tx.Rollback().Error; err != nil {
			return
		}
	}()
	var toAccnt, fromAccnt BankAccount
	if err := a.db.Where("account_number = ?", fromAccntNum).First(&fromAccnt).Error; err != nil {
		return err
	}
	if err := a.db.Where("account_number = ?", toAccntNum).First(&toAccnt).Error; err != nil {
		return err
	}
	t := BankTransfer{
		TransferUUID:      uuid.New(),
		FromAccountUUID:   fromAccnt.AccountUUID,
		ToAccountUUID:     toAccnt.AccountUUID,
		Amount:            amount,
		Currency:          currency,
		TransferTimestamp: ts,
	}
	if err := tx.Create(&t).Error; err != nil {
		return err
	}
	fromTrans := BankTransaction{
		TransactionUUID:      uuid.New(),
		AccountUUID:          fromAccnt.AccountUUID,
		TransactionTimestamp: ts,
		TransactionType:      TransactionTypeOut,
		Amount:               amount,
	}
	if err := tx.Create(&fromTrans).Error; err != nil {
		return err
	}
	toTrans := BankTransaction{
		TransactionUUID:      uuid.New(),
		AccountUUID:          toAccnt.AccountUUID,
		TransactionTimestamp: ts,
		TransactionType:      TransactionTypeIn,
		Amount:               amount,
	}
	if err := tx.Create(&toTrans).Error; err != nil {
		return err
	}
	fromAccnt.CurrentBalance -= amount
	if err := tx.Save(&fromAccnt).Error; err != nil {
		return err
	}
	toAccnt.CurrentBalance += amount
	if err := tx.Save(&toAccnt).Error; err != nil {
		return err
	}
	t.TransferSuccess = true
	if err := tx.Save(&t).Error; err != nil {
		return err
	}
	tx.Commit()
	return nil
}
