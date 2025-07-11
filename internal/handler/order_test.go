package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mikhailshtv/stockLkBack/internal/model"
	"github.com/mikhailshtv/stockLkBack/internal/service"
	mock_service "github.com/mikhailshtv/stockLkBack/internal/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockBehaviorCreateOrder func(s *mock_service.MockOrder, orderReq model.OrderRequestBody)

func TestHandler_CreateOrder(t *testing.T) {
	tests := []struct {
		name                 string
		inputBody            string
		inputOrder           model.OrderRequestBody
		mockBehavior         mockBehaviorCreateOrder
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			inputBody: `{
				"products":[
					{
						"productId":1,
						"quantity":1,
						"sellPrice":74000
					}
				]
			}`,
			inputOrder: model.OrderRequestBody{
				Products: []model.OrderProduct{
					{
						ProductID: 1,
						Quantity:  1,
						SellPrice: 74000,
					},
				},
			},
			mockBehavior: func(s *mock_service.MockOrder, orderReq model.OrderRequestBody) {
				s.EXPECT().Create(orderReq, 1).Return(
					&model.Order{
						ID:               1,
						Number:           1,
						TotalCost:        74000,
						CreatedDate:      time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
						LastModifiedDate: time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
						Status:           model.StatusActive,
						Products: []model.Product{
							{
								ID:        1,
								Code:      137207,
								Quantity:  1,
								Name:      "Bacon",
								SellPrice: 74000,
							},
						},
						UserID: 1,
					}, nil,
				)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `{
				"id":1,
				"number":1,
				"totalCost":74000,
				"createdDate":"2025-05-25T12:17:16.550631Z",
				"lastModifiedDate":"2025-05-25T12:17:16.550631Z",
				"status":{"key":"active","displayName":"Активный"},
				"products":[
					{
						"id":1,
						"code":137207,
						"quantity":1,
						"name":"Bacon",
						"sellPrice":74000
					}
				],
				"userId":1
			}`,
		},
		{
			name: "Некорректное тело запроса",
			inputBody: `
				{
					"id":1,
					"code":14823,
					"quantity":215,
					"name":"Cheese",
					"purchasePrice":24000,
					"sellPrice":74000
				}`,
			mockBehavior:         func(_ *mock_service.MockOrder, _ model.OrderRequestBody) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"Некорректное тело запроса"}`,
		},
		{
			name: "Ошибка 500",
			inputBody: `
				{
					"products":[]
				}`,
			inputOrder: model.OrderRequestBody{
				Products: []model.OrderProduct{},
			},
			mockBehavior: func(s *mock_service.MockOrder, orderReq model.OrderRequestBody) {
				s.EXPECT().Create(orderReq, 1).Return(nil, errors.New("список товаров не может быть пустым"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"список товаров не может быть пустым"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := repo(t)
			test.mockBehavior(repo, test.inputOrder)
			services := &service.Service{Order: repo}
			handler := NewHandler(services)
			r := gin.New()
			r.POST("/orders", func(ctx *gin.Context) {
				ctx.Set("userId", 1)
				handler.CreateOrder(ctx)
			})
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/orders", bytes.NewBufferString(test.inputBody))
			r.ServeHTTP(w, req)
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.JSONEq(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_ListOrders(t *testing.T) {
	type mockBehavior func(s *mock_service.MockOrder)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			mockBehavior: func(s *mock_service.MockOrder) {
				s.EXPECT().GetAll(1, model.RoleEmployee).Return(
					[]model.Order{
						{
							ID:               1,
							Number:           1,
							TotalCost:        74000,
							CreatedDate:      time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
							LastModifiedDate: time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
							Status:           model.StatusActive,
							Products: []model.Product{
								{
									ID:        1,
									Code:      14823,
									Quantity:  1,
									Name:      "Cheese",
									SellPrice: 74000,
								},
							},
							UserID: 1,
						},
					}, nil,
				)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `[
				{
					"id":1,
					"number":1,
					"totalCost":74000,
					"createdDate":"2025-05-25T12:17:16.550631Z",
					"lastModifiedDate":"2025-05-25T12:17:16.550631Z",
					"status":{"key":"active","displayName":"Активный"},
					"products":[
						{
							"id":1,
							"code":14823,
							"quantity":1,
							"name":"Cheese",
							"sellPrice":74000
						}
					],
					"userId":1
				}
			]`,
		},
		{
			name: "Ошибка коннекта к базе",
			mockBehavior: func(s *mock_service.MockOrder) {
				s.EXPECT().GetAll(1, model.RoleEmployee).Return(
					nil, errors.New("Ошибка коннекта к базе"),
				)
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"Ошибка коннекта к базе"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := repo(t)
			test.mockBehavior(repo)
			services := &service.Service{Order: repo}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/orders", func(ctx *gin.Context) {
				ctx.Set("userId", 1)
				ctx.Set("role", model.RoleEmployee)
				handler.ListOrders(ctx)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/orders", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.JSONEq(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_EditOrders(t *testing.T) {
	type mockBehavior func(s *mock_service.MockOrder, requestBody model.OrderRequestBody)

	tests := []struct {
		name                 string
		inputBody            string
		inputOrder           model.OrderRequestBody
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			inputBody: `
				{
					"products":[
						{
							"productId":1,
							"quantity":1,
							"sellPrice":74000
						}
					]
				}`,
			inputOrder: model.OrderRequestBody{
				Products: []model.OrderProduct{
					{
						ProductID: 1,
						Quantity:  1,
						SellPrice: 74000,
					},
				},
			},
			mockBehavior: func(s *mock_service.MockOrder, requestBody model.OrderRequestBody) {
				s.EXPECT().Update(1, requestBody, 1).Return(
					&model.Order{
						ID:               1,
						Number:           1,
						TotalCost:        74000,
						CreatedDate:      time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
						LastModifiedDate: time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
						Status:           model.StatusActive,
						Products: []model.Product{
							{
								ID:        1,
								Code:      14823,
								Quantity:  1,
								Name:      "Cheese",
								SellPrice: 74000,
							},
						},
						UserID: 1,
					}, nil,
				)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `
				{
					"id":1,
					"number":1,
					"totalCost":74000,
					"createdDate":"2025-05-25T12:17:16.550631Z",
					"lastModifiedDate":"2025-05-25T12:17:16.550631Z",
					"status":{"key":"active","displayName":"Активный"},
					"products":[
						{
							"id":1,
							"code":14823,
							"quantity":1,
							"name":"Cheese",
							"sellPrice":74000
						}
					],
					"userId":1
				}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := repo(t)
			test.mockBehavior(repo, test.inputOrder)
			services := &service.Service{Order: repo}
			handler := NewHandler(services)
			r := gin.New()
			r.PUT("/orders/:id", func(ctx *gin.Context) {
				ctx.Set("userId", 1)
				handler.EditOrder(ctx)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/orders/1", bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.JSONEq(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_GetOrderById(t *testing.T) {
	type mockBehavior func(s *mock_service.MockOrder)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			mockBehavior: func(s *mock_service.MockOrder) {
				s.EXPECT().GetByID(1, 1, model.RoleEmployee).Return(
					&model.Order{
						ID:               1,
						Number:           1,
						TotalCost:        74000,
						CreatedDate:      time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
						LastModifiedDate: time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
						Status:           model.StatusActive,
						Products: []model.Product{
							{
								ID:        1,
								Code:      14823,
								Quantity:  1,
								Name:      "Cheese",
								SellPrice: 74000,
							},
						},
						UserID: 1,
					}, nil,
				)
			},
			expectedStatusCode: 200,
			expectedResponseBody: `
				{
					"id":1,
					"number":1,
					"totalCost":74000,
					"createdDate":"2025-05-25T12:17:16.550631Z",
					"lastModifiedDate":"2025-05-25T12:17:16.550631Z",
					"status":{"key":"active","displayName":"Активный"},
					"products":[
						{
							"id":1,
							"code":14823,
							"quantity":1,
							"name":"Cheese",
							"sellPrice":74000
						}
					],
					"userId":1
				}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := repo(t)
			test.mockBehavior(repo)
			services := &service.Service{Order: repo}
			handler := NewHandler(services)
			r := gin.New()
			r.GET("/orders/:id", func(ctx *gin.Context) {
				ctx.Set("userId", 1)
				ctx.Set("role", model.RoleEmployee)
				handler.GetOrderByID(ctx)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/orders/1", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.JSONEq(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_DeleteOrder(t *testing.T) {
	type mockBehavior func(s *mock_service.MockOrder)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			mockBehavior: func(s *mock_service.MockOrder) {
				s.EXPECT().Delete(1, 1).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"Success","message":"Объект успешно удален"}`,
		},
		{
			name: "Ошибка удаления",
			mockBehavior: func(s *mock_service.MockOrder) {
				s.EXPECT().Delete(1, 1).Return(errors.New("ошибка сохранения в файл"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"ошибка сохранения в файл"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := repo(t)
			test.mockBehavior(repo)

			services := &service.Service{Order: repo}
			handler := NewHandler(services)

			r := gin.New()
			r.DELETE("/orders/:id", func(ctx *gin.Context) {
				ctx.Set("userId", 1)
				handler.DeleteOrder(ctx)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/orders/1", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func repo(t *testing.T) *mock_service.MockOrder {
	t.Helper()
	c := gomock.NewController(t)
	t.Cleanup(func() { c.Finish() })
	return mock_service.NewMockOrder(c)
}
