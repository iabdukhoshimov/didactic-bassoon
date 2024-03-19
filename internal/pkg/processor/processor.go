package processor

import (
	"context"
	"encoding/json"

	"gitlab.com/tour/internal/pkg/logger"

	"go.uber.org/zap"
)

type ExecutorWithResp[K any, Z any] func(ctx context.Context, input K) (Z, error)
type Executor[K any] func(ctx context.Context, input K) error

func ExecuteWithResp[T any, K any, U any](ctx context.Context, input T, f ExecutorWithResp[K, U]) (U, error) {
	var (
		resp U
	)

	logger.Log.Debug("Execute with Resp: ", zap.Any("input", input))

	rawBytes, err := json.Marshal(input)
	if err != nil {
		logger.Log.Debug("Execute with Resp: ", zap.Any("input", input))
		return resp, err
	}

	var dbPayload K
	err = json.Unmarshal(rawBytes, &dbPayload)
	if err != nil {
		logger.Log.Error("failed to convert json to struct", zap.Any("error", err))
		return resp, err
	}

	resp, err = f(ctx, dbPayload)
	if err != nil {
		logger.Log.Error("failed to execute database operation", zap.Any("error", err))
		return resp, err
	}

	logger.Log.Debug("Request succeeded")

	return resp, nil
}

func Execute[T any, K any](ctx context.Context, input T, f Executor[K]) error {
	logger.Log.Debug("Execute with Resp: ", zap.Any("input", input))

	rawBytes, err := json.Marshal(input)
	if err != nil {
		logger.Log.Error("failed to convert input to json", zap.Any("error", err))
		return err
	}

	var dbPayload K
	err = json.Unmarshal(rawBytes, &dbPayload)
	if err != nil {
		logger.Log.Error("failed to convert json to struct", zap.Any("error", err))
		return err
	}

	err = f(ctx, dbPayload)
	if err != nil {
		logger.Log.Error("failed to execute database operation", zap.Any("error", err))
		return err
	}

	logger.Log.Debug("Request succeeded")

	return nil
}
