package exception

import (
	"ai-service/internal/util/logger"
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
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

func TranslateDatabaseError(ctx context.Context, err error) error {
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("empty result")
		}

		// Handle PostgreSQL errors
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				return errors.New("duplicate key value")
			case "23503": // foreign_key_violation
				return errors.New("foreign key violation")
			case "23502": // not_null_violation
				return errors.New("required field is missing")
			case "42P01": // undefined_table
				return errors.New("table does not exist")
			case "42703": // undefined_column
				return errors.New("column does not exist")
			case "23514": // check_violation
				return errors.New("check constraint violation")
			case "23513": // exclusion_violation
				return errors.New("exclusion constraint violation")
			default:
				logger.Error(ctx, err)
				return err
			}
		}

		// Handle other database errors
		logger.Error(ctx, err)
		return err
	}
	return nil
}

// TranslateMysqlError is kept for backward compatibility but now delegates to TranslateDatabaseError
func TranslateMysqlError(ctx context.Context, err error) error {
	return TranslateDatabaseError(ctx, err)
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
