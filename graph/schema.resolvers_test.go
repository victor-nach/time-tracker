package graph

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/victor-nach/time-tracker/lib/rerrors"
	"github.com/victor-nach/time-tracker/mocks"
	"github.com/victor-nach/time-tracker/models"
	"testing"
)

func TestQueryResolver_CalculatePrice(t *testing.T) {
	const (
		success = iota
		priceSrvErr
	)

	var tests = []struct {
		testType     int
		name         string
		margin       float64
		tradeType    string
		exchangeRate float64
	}{
		{
			name:         "Test calculate price successfully",
			testType:     success,
			margin:       0.2,
			tradeType:    string(models.TradeTypeBuy),
			exchangeRate: 498,
		},
		{
			name:         "Test invalid request",
			testType:     priceSrvErr,
			margin:       0.2,
			tradeType:    string(models.TradeTypeBuy),
			exchangeRate: 498,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			priceServiceClient := new(mocks.PriceService)

			resolver := NewResolver(priceServiceClient).Query()

			switch testCase.testType {
			case success:
				priceServiceClient.On("GetCurrentPrice").Return(&models.CurrentPrice{
					Bpi: models.Bpi{
						USD: models.USD{
							RateFloat: 57673.4539,
						},
					},
				}, nil)

				response, err := resolver.CalculatePrice(
					context.Background(),
					models.TradeType(testCase.tradeType),
					testCase.margin,
					testCase.exchangeRate,
				)
				assert.NoError(t, err)
				assert.NotNil(t, response)
			case priceSrvErr:
				priceServiceClient.On("GetCurrentPrice").Return(nil, rerrors.Format(rerrors.InternalErr, nil))

				_, err := resolver.CalculatePrice(
					context.Background(),
					models.TradeType(testCase.tradeType),
					testCase.margin,
					testCase.exchangeRate,
				)
				assert.Error(t, err)
			}

			priceServiceClient.AssertExpectations(t)
		})
	}
}
