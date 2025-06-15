package handler

import (
	"bytes"
	"errors"
	"golang/stockLkBack/internal/model"
	"golang/stockLkBack/internal/service"
	mock_service "golang/stockLkBack/internal/service/mocks"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestHandler_CreateOrder(t *testing.T) {
	type mockBehavior func(s *mock_service.MockOrder, orderReq model.OrderRequestBody)

	tests := []struct {
		name                 string
		inputBody            string
		inputOrder           model.OrderRequestBody
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"products":[{"id":1,"code":14823,"quantity":215,"name":"Cheese","purchasePrice":24000,"salePrice":74000}]}`,
			inputOrder: model.OrderRequestBody{
				Products: []model.Product{
					{
						Id:            1,
						Code:          14823,
						Quantity:      215,
						Name:          "Cheese",
						PurchasePrice: 24000,
						SalePrice:     74000,
					},
				},
			},
			mockBehavior: func(s *mock_service.MockOrder, orderReq model.OrderRequestBody) {
				s.EXPECT().Create(orderReq).Return(
					&model.Order{
						Id:               1,
						Number:           1,
						TotalCost:        74000,
						CreatedDate:      time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
						LastModifiedDate: time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
						Status:           1,
						Products: []model.Product{
							{
								Id:            1,
								Code:          14823,
								Quantity:      215,
								Name:          "Cheese",
								PurchasePrice: 24000,
								SalePrice:     74000,
							},
						},
					}, nil,
				)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"number":1,"totalCost":74000,"createdDate":"2025-05-25T12:17:16.550631Z","lastModifiedDate":"2025-05-25T12:17:16.550631Z","status":1,"products":[{"id":1,"code":14823,"quantity":215,"name":"Cheese","purchasePrice":24000,"salePrice":74000}]}`,
		},
		{
			name:                 "Некорректное тело запроса",
			inputBody:            `{"id":1,"code":14823,"quantity":215,"name":"Cheese","purchasePrice":24000,"salePrice":74000}`,
			mockBehavior:         func(s *mock_service.MockOrder, orderReq model.OrderRequestBody) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Некорректное тело запроса"}`,
		},
		{
			name:      "Ошибка сохранения в файл",
			inputBody: `{"products":[{"id":1,"code":14823,"quantity":215,"name":"Cheese","purchasePrice":24000,"salePrice":74000}]}`,
			inputOrder: model.OrderRequestBody{
				Products: []model.Product{
					{
						Id:            1,
						Code:          14823,
						Quantity:      215,
						Name:          "Cheese",
						PurchasePrice: 24000,
						SalePrice:     74000,
					},
				},
			},
			mockBehavior: func(s *mock_service.MockOrder, orderReq model.OrderRequestBody) {
				s.EXPECT().Create(orderReq).Return(nil, errors.New("ошибка сохранения в файл"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"ошибка сохранения в файл"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockOrder(c)
			test.mockBehavior(repo, test.inputOrder)

			services := &service.Service{Order: repo}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/orders", handler.CreateOrder)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/orders", bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
