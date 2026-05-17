package pkg

import "go.uber.org/zap"

func NewProductionLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

func NewDevelopmentLogger() (*zap.Logger, error) {
	return zap.NewDevelopment()
}
