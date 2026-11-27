package grpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/OleksUMD/grpc_test_server/internal/adapter/database"
	"github.com/OleksUMD/proto-example/protogen/go/bank"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/genproto/googleapis/type/datetime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *GrpcAdapter) GetCurrentBalance(ctx context.Context, req *bank.CurrentBalanceRequest) (*bank.CurrentBalanceResponse, error) {
	balance, err := a.bankService.GetCurrentBalance(req.AccountNumber)
	if err != nil {
		s := status.New(codes.FailedPrecondition, fmt.Sprintf("account %v not found", req.AccountNumber))
		s, _ = s.WithDetails(&errdetails.ErrorInfo{})
		s, _ = s.WithDetails(&errdetails.DebugInfo{})
		return nil, s.Err()
	}
	return &bank.CurrentBalanceResponse{
		Amount: balance,
	}, nil
}

func (a *GrpcAdapter) FetchExchangeRates(req *bank.ExchangeRateRequest, stream grpc.ServerStreamingServer[bank.ExchangeRateResponse]) error {
	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			log.Println("Client cancelled stream")
			return nil
		default:
		}
		now := time.Now()
		rate, err := a.bankService.GetCurrentExchangeRates(req.FromCurrency, req.ToCurrency, time.Now())
		if err != nil {
			s := status.New(codes.InvalidArgument,
				"Currency not valid. Please use valid currency for both from and to",
			)
			s, _ = s.WithDetails(&errdetails.ErrorInfo{
				Domain: "go-grpc-server.edu.com",
				Reason: "INVALID_CURRENCY",
				Metadata: map[string]string{
					"from_currency": req.FromCurrency,
					"to_currency":   req.ToCurrency,
				},
			})
			return s.Err()
		}
		resp := &bank.ExchangeRateResponse{
			FromCurrency: req.FromCurrency,
			ToCurrency:   req.ToCurrency,
			Rate:         rate,
			Timestamp:    now.Format(time.RFC3339),
		}
		fmt.Println("Sending exchange rate to client")
		stream.Send(resp)
		time.Sleep(3 * time.Second)
	}
}
func (a *GrpcAdapter) SummarizeTransactions(stream grpc.ClientStreamingServer[bank.Transaction, bank.TransactionSummary]) error {
	tsum := a.bankService.GetEmptyTransactionSummary()
	acct := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			res := bank.TransactionSummary{
				AccountNumber: acct,
				SumAmountIn:   tsum.SumIn,
				SumAmountOut:  tsum.SumOut,
				SumTotal:      tsum.SumTotal,
				TransactionDate: &date.Date{
					Year:  int32(tsum.SummaryOnDate.Year()),
					Month: int32(tsum.SummaryOnDate.Month()),
					Day:   int32(tsum.SummaryOnDate.Day()),
				},
			}
			return stream.SendAndClose(&res)
		}
		if err != nil {
			return err
		}
		acct = req.AccountNumber
		ts, err := toTime(req.Timestamp)
		if err != nil {
			return err
		}
		ttype := database.TransactionTypeUnknown
		if req.Type == bank.TransactionType_TRANSACTION_TYPE_IN {
			ttype = database.TransactionTypeIn
		} else if req.Type == bank.TransactionType_TRANSACTION_TYPE_OUT {
			ttype = database.TransactionTypeOut
		}
		tcur := database.BankTransaction{
			TransactionUUID:      uuid.New(),
			Amount:               req.Amount,
			TransactionType:      ttype,
			TransactionTimestamp: ts,
			Notes:                req.Notes,
		}
		_, err = a.bankService.CreateTransaction(acct, &tcur)
		err = a.bankService.SummarizeTransactions(tsum, &tcur)
		if err != nil {
			return err
		}
	}
}

func (a *GrpcAdapter) TransferMultiple(stream grpc.BidiStreamingServer[bank.TransferRequest, bank.TransferResponse]) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		transferStatus := bank.TransferStatus_TRANSFER_STATUS_SUCCESS
		now := time.Now()
		err = a.bankService.MultipleMoneyTransfer(req.FromAccountNumber, req.ToAccountNumber, req.Curency, req.Amount, now)
		if err != nil {
			log.Println(err)
			transferStatus = bank.TransferStatus_TRANSFER_STATUS_FAILED
		}
		resp := &bank.TransferResponse{
			FromAccountNumber: req.FromAccountNumber,
			ToAccountNumber:   req.ToAccountNumber,
			Curency:           req.Curency,
			Amount:            req.Amount,
			Status:            transferStatus,
			Timestamp:         toDateTime(now),
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}

func toTime(dt *datetime.DateTime) (time.Time, error) {
	if dt == nil {
		now := time.Now()
		dt = &datetime.DateTime{
			Year:    int32(now.Year()),
			Month:   int32(now.Month()),
			Day:     int32(now.Day()),
			Hours:   int32(now.Hour()),
			Minutes: int32(now.Minute()),
			Seconds: int32(now.Second()),
			Nanos:   int32(now.Nanosecond()),
		}
	}
	res := time.Date(
		int(dt.Year),
		time.Month(dt.Month),
		int(dt.Day),
		int(dt.Hours),
		int(dt.Minutes),
		int(dt.Seconds),
		int(dt.Nanos),
		time.UTC,
	)
	return res, nil
}

func toDateTime(t time.Time) *datetime.DateTime {
	return &datetime.DateTime{
		Year:    int32(t.Year()),
		Month:   int32(t.Month()),
		Day:     int32(t.Day()),
		Hours:   int32(t.Hour()),
		Minutes: int32(t.Minute()),
		Seconds: int32(t.Second()),
		Nanos:   int32(t.Nanosecond()),
	}
}
