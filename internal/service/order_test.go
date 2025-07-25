package service

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/mikhailshtv/stockLkBack/internal/model"
	"github.com/mikhailshtv/stockLkBack/internal/repository"
	mocks "github.com/mikhailshtv/stockLkBack/internal/service/mocks"

	"github.com/golang/mock/gomock"
)

func TestOrderService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	dbMock := mocks.NewMockOrder(ctrl)

	tests := []struct {
		name    string
		mock    func()
		args    model.OrderRequestBody
		want    *model.Order
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				dbMock.EXPECT().Create(
					model.OrderRequestBody{
						Products: []model.OrderProduct{
							{
								ProductID: 1,
								Quantity:  1,
								SellPrice: 74000,
							},
						},
					}, 1,
				).Return(
					&model.Order{
						ID:               1,
						Number:           1,
						TotalCost:        74000,
						CreatedDate:      time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
						LastModifiedDate: time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
						Status:           model.StatusActive,
						Products: []model.Product{
							{
								ID:            1,
								Code:          14823,
								Quantity:      215,
								Name:          "Cheese",
								PurchasePrice: 24000,
								SellPrice:     74000,
							},
						},
					}, nil,
				)
			},
			args: model.OrderRequestBody{
				Products: []model.OrderProduct{
					{
						ProductID: 1,
						Quantity:  1,
						SellPrice: 74000,
					},
				},
			},
			want: &model.Order{
				ID:               1,
				Number:           1,
				TotalCost:        74000,
				CreatedDate:      time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
				LastModifiedDate: time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
				Status:           model.StatusActive,
				Products: []model.Product{
					{
						ID:            1,
						Code:          14823,
						Quantity:      215,
						Name:          "Cheese",
						PurchasePrice: 24000,
						SellPrice:     74000,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			mock: func() {
				dbMock.EXPECT().Create(
					model.OrderRequestBody{
						Products: []model.OrderProduct{
							{
								ProductID: 1,
								Quantity:  1,
								SellPrice: 74000,
							},
						},
					}, 1,
				).Return(
					nil,
					errors.New("ошибка сохранения в файл"),
				)
			},
			args: model.OrderRequestBody{
				Products: []model.OrderProduct{
					{
						ProductID: 1,
						Quantity:  1,
						SellPrice: 74000,
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os := &Service{Order: dbMock}

			tt.mock()

			got, err := os.Order.Create(tt.args, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ошибка создания заказа error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ошибка создания заказа got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })
	dbMock := mocks.NewMockOrder(ctrl)

	type requestBody struct {
		id   int
		body model.OrderRequestBody
	}

	tests := []struct {
		name    string
		mock    func()
		args    requestBody
		want    *model.Order
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				dbMock.EXPECT().Update(
					1,
					model.OrderRequestBody{
						Products: []model.OrderProduct{
							{
								ProductID: 1,
								Quantity:  1,
								SellPrice: 74000,
							},
						},
					}, 1,
				).Return(
					&model.Order{
						ID:               1,
						Number:           1,
						TotalCost:        74000,
						CreatedDate:      time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
						LastModifiedDate: time.Date(2025, time.June, 15, 12, 0o0, 0o0, 0, time.UTC),
						Status:           model.StatusActive,
						Products: []model.Product{
							{
								ID:            1,
								Code:          14823,
								Quantity:      215,
								Name:          "Pizza",
								PurchasePrice: 24000,
								SellPrice:     74000,
							},
						},
					}, nil,
				)
			},
			args: requestBody{
				id: 1,
				body: model.OrderRequestBody{
					Products: []model.OrderProduct{
						{
							ProductID: 1,
							Quantity:  1,
							SellPrice: 74000,
						},
					},
				},
			},
			want: &model.Order{
				ID:               1,
				Number:           1,
				TotalCost:        74000,
				CreatedDate:      time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
				LastModifiedDate: time.Date(2025, time.June, 15, 12, 0o0, 0o0, 0, time.UTC),
				Status:           model.StatusActive,
				Products: []model.Product{
					{
						ID:            1,
						Code:          14823,
						Quantity:      215,
						Name:          "Pizza",
						PurchasePrice: 24000,
						SellPrice:     74000,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			mock: func() {
				dbMock.EXPECT().Update(
					1,
					model.OrderRequestBody{
						Products: []model.OrderProduct{
							{
								ProductID: 1,
								Quantity:  1,
								SellPrice: 74000,
							},
						},
					}, 1,
				).Return(
					nil,
					errors.New(repository.NotFoundErrorMessage),
				)
			},
			args: requestBody{
				id: 1,
				body: model.OrderRequestBody{
					Products: []model.OrderProduct{
						{
							ProductID: 1,
							Quantity:  1,
							SellPrice: 74000,
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os := &Service{Order: dbMock}
			tt.mock()
			got, err := os.Order.Update(tt.args.id, tt.args.body, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ошибка обноления заказа error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ошибка обноления заказа got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	dbMock := mocks.NewMockOrder(ctrl)

	tests := []struct {
		name    string
		mock    func()
		want    []model.Order
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				dbMock.EXPECT().GetAll(1, model.RoleEmployee).Return(
					[]model.Order{
						{
							ID:               1,
							Number:           1,
							TotalCost:        74000,
							CreatedDate:      time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
							LastModifiedDate: time.Date(2025, time.June, 15, 12, 0o0, 0o0, 0, time.UTC),
							Status:           model.StatusActive,
							Products: []model.Product{
								{
									ID:            1,
									Code:          14823,
									Quantity:      215,
									Name:          "Pizza",
									PurchasePrice: 24000,
									SellPrice:     74000,
								},
							},
						},
					}, nil,
				)
			},
			want: []model.Order{
				{
					ID:               1,
					Number:           1,
					TotalCost:        74000,
					CreatedDate:      time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
					LastModifiedDate: time.Date(2025, time.June, 15, 12, 0o0, 0o0, 0, time.UTC),
					Status:           model.StatusActive,
					Products: []model.Product{
						{
							ID:            1,
							Code:          14823,
							Quantity:      215,
							Name:          "Pizza",
							PurchasePrice: 24000,
							SellPrice:     74000,
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os := &Service{Order: dbMock}

			tt.mock()

			got, err := os.Order.GetAll(1, model.RoleEmployee)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ошибка получения списка заказов error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ошибка получения списка заказов got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderService_GetById(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	dbMock := mocks.NewMockOrder(ctrl)

	tests := []struct {
		name    string
		mock    func()
		args    int
		want    *model.Order
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				dbMock.EXPECT().GetByID(1, 1, model.RoleEmployee).Return(
					&model.Order{
						ID:               1,
						Number:           1,
						TotalCost:        74000,
						CreatedDate:      time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
						LastModifiedDate: time.Date(2025, time.June, 15, 12, 0o0, 0o0, 0, time.UTC),
						Status:           model.StatusActive,
						Products: []model.Product{
							{
								ID:            1,
								Code:          14823,
								Quantity:      215,
								Name:          "Pizza",
								PurchasePrice: 24000,
								SellPrice:     74000,
							},
						},
					}, nil,
				)
			},
			args: 1,
			want: &model.Order{
				ID:               1,
				Number:           1,
				TotalCost:        74000,
				CreatedDate:      time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
				LastModifiedDate: time.Date(2025, time.June, 15, 12, 0o0, 0o0, 0, time.UTC),
				Status:           model.StatusActive,
				Products: []model.Product{
					{
						ID:            1,
						Code:          14823,
						Quantity:      215,
						Name:          "Pizza",
						PurchasePrice: 24000,
						SellPrice:     74000,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			mock: func() {
				dbMock.EXPECT().GetByID(1, 1, model.RoleEmployee).Return(
					nil,
					errors.New(repository.NotFoundErrorMessage),
				)
			},
			args:    1,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os := &Service{Order: dbMock}

			tt.mock()

			got, err := os.Order.GetByID(tt.args, 1, model.RoleEmployee)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ошибка получения заказа по id error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ошибка получения заказа по id got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	dbMock := mocks.NewMockOrder(ctrl)

	tests := []struct {
		name    string
		mock    func()
		args    int
		want    error
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				dbMock.EXPECT().Delete(1, 1).Return(nil)
			},
			args:    1,
			want:    nil,
			wantErr: false,
		},
		{
			name: "error",
			mock: func() {
				dbMock.EXPECT().Delete(1, 1).Return(
					errors.New("ошибка сохранения в файл"),
				)
			},
			args:    1,
			want:    errors.New("ошибка сохранения в файл"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os := &Service{Order: dbMock}

			tt.mock()

			err := os.Order.Delete(tt.args, 1)
			if !errors.Is(err, tt.want) && err.Error() != tt.want.Error() {
				t.Errorf("Ошибка удаления заказа по id error = %v, want %v", err, tt.want)
				return
			}
		})
	}
}

func TestOrderService_UpdateStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })
	dbMock := mocks.NewMockOrder(ctrl)

	type requestBody struct {
		id   int
		body model.OrderStatusRequest
	}

	tests := []struct {
		name    string
		mock    func()
		args    requestBody
		want    *model.Order
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				dbMock.EXPECT().UpdateStatus(
					1,
					model.OrderStatusRequest{
						Status: model.OrderStatus{
							Key:         "executed",
							DisplayName: "Выполнен",
						},
					}, 1,
				).Return(
					&model.Order{
						ID:               1,
						Number:           1,
						TotalCost:        74000,
						CreatedDate:      time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
						LastModifiedDate: time.Date(2025, time.June, 15, 12, 0o0, 0o0, 0, time.UTC),
						Status:           model.StatusActive,
						Products: []model.Product{
							{
								ID:            1,
								Code:          14823,
								Quantity:      215,
								Name:          "Pizza",
								PurchasePrice: 24000,
								SellPrice:     74000,
							},
						},
					}, nil,
				)
			},
			args: requestBody{
				id: 1,
				body: model.OrderStatusRequest{
					Status: model.OrderStatus{
						Key:         "executed",
						DisplayName: "Выполнен",
					},
				},
			},
			want: &model.Order{
				ID:               1,
				Number:           1,
				TotalCost:        74000,
				CreatedDate:      time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
				LastModifiedDate: time.Date(2025, time.June, 15, 12, 0o0, 0o0, 0, time.UTC),
				Status:           model.StatusActive,
				Products: []model.Product{
					{
						ID:            1,
						Code:          14823,
						Quantity:      215,
						Name:          "Pizza",
						PurchasePrice: 24000,
						SellPrice:     74000,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			mock: func() {
				dbMock.EXPECT().UpdateStatus(
					1,
					model.OrderStatusRequest{
						Status: model.OrderStatus{
							Key:         "executed",
							DisplayName: "Выполнен",
						},
					}, 1,
				).Return(
					nil,
					errors.New(repository.NotFoundErrorMessage),
				)
			},
			args: requestBody{
				id: 1,
				body: model.OrderStatusRequest{
					Status: model.OrderStatus{
						Key:         "executed",
						DisplayName: "Выполнен",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os := &Service{Order: dbMock}
			tt.mock()
			got, err := os.Order.UpdateStatus(tt.args.id, tt.args.body, 1)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ошибка обноления заказа error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ошибка обноления заказа got = %v, want %v", got, tt.want)
			}
		})
	}
}
