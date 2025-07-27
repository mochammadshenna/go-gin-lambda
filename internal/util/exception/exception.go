package exception

import (
	"ai-service/internal/util/logger"
	"context"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

func PanicOnError(err error) {
	if err != nil {
		logger.Logger.Error(err)
		panic(err)
	}
}

func PanicOnErrorContext(ctx context.Context, err error) {
	if err != nil {
		logger.Error(ctx, err)
		panic(err)
	}
}

func TranslateMysqlError(ctx context.Context, err error) error {
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = errors.New("empty result")
		}
		var myErr *mysql.MySQLError
		if errors.As(err, &myErr) {
			if myErr.Number == 1062 {
				err = errors.New("duplicate key value")
			} else if myErr.Number == 1452 {
				err = errors.New("foreign key violation")
			} else {
				logger.Error(ctx, err)
			}
		}
	}
	return err
}

func CancelBackground(ctx context.Context, cancel context.CancelFunc, errorMessage string, successMessage string) {
	select {
	case <-ctx.Done():
		if len(errorMessage) > 0 {
			logger.Errorf(ctx, "%s: %v", errorMessage, ctx.Err())
		}
		cancel()
		return
	default:
		if len(successMessage) > 0 {
			logger.Info(ctx, successMessage)
		}
		return
	}
}
