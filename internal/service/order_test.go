package service

import (
	"errors"
	"golang/stockLkBack/internal/model"
	mocks "golang/stockLkBack/internal/service/mocks"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestOrderService_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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
				).Return(
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
			args: model.OrderRequestBody{
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
			want: &model.Order{
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
			},
			wantErr: false,
		},
		{
			name: "error",
			mock: func() {
				dbMock.EXPECT().Create(
					model.OrderRequestBody{
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
				).Return(
					nil,
					errors.New("ошибка сохранения в файл"),
				)
			},
			args: model.OrderRequestBody{
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
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			os := &Service{Order: dbMock}

			tt.mock()

			got, err := os.Order.Create(tt.args)
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
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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
						Products: []model.Product{
							{
								Id:            1,
								Code:          14823,
								Quantity:      215,
								Name:          "Pizza",
								PurchasePrice: 24000,
								SalePrice:     74000,
							},
						},
					},
				).Return(
					&model.Order{
						Id:               1,
						Number:           1,
						TotalCost:        74000,
						CreatedDate:      time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
						LastModifiedDate: time.Date(2025, time.June, 15, 12, 00, 00, 0, time.UTC),
						Status:           1,
						Products: []model.Product{
							{
								Id:            1,
								Code:          14823,
								Quantity:      215,
								Name:          "Pizza",
								PurchasePrice: 24000,
								SalePrice:     74000,
							},
						},
					}, nil,
				)
			},
			args: requestBody{
				id: 1,
				body: model.OrderRequestBody{
					Products: []model.Product{
						{
							Id:            1,
							Code:          14823,
							Quantity:      215,
							Name:          "Pizza",
							PurchasePrice: 24000,
							SalePrice:     74000,
						},
					},
				},
			},
			want: &model.Order{
				Id:               1,
				Number:           1,
				TotalCost:        74000,
				CreatedDate:      time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
				LastModifiedDate: time.Date(2025, time.June, 15, 12, 00, 00, 0, time.UTC),
				Status:           1,
				Products: []model.Product{
					{
						Id:            1,
						Code:          14823,
						Quantity:      215,
						Name:          "Pizza",
						PurchasePrice: 24000,
						SalePrice:     74000,
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
						Products: []model.Product{
							{
								Id:            1,
								Code:          14823,
								Quantity:      215,
								Name:          "Pizza",
								PurchasePrice: 24000,
								SalePrice:     74000,
							},
						},
					},
				).Return(
					nil,
					errors.New("элемент не найден"),
				)
			},
			args: requestBody{
				id: 1,
				body: model.OrderRequestBody{
					Products: []model.Product{
						{
							Id:            1,
							Code:          14823,
							Quantity:      215,
							Name:          "Pizza",
							PurchasePrice: 24000,
							SalePrice:     74000,
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
			t.Parallel()

			os := &Service{Order: dbMock}

			tt.mock()

			got, err := os.Order.Update(tt.args.id, tt.args.body)
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
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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
				dbMock.EXPECT().GetAll().Return(
					[]model.Order{
						{
							Id:               1,
							Number:           1,
							TotalCost:        74000,
							CreatedDate:      time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
							LastModifiedDate: time.Date(2025, time.June, 15, 12, 00, 00, 0, time.UTC),
							Status:           1,
							Products: []model.Product{
								{
									Id:            1,
									Code:          14823,
									Quantity:      215,
									Name:          "Pizza",
									PurchasePrice: 24000,
									SalePrice:     74000,
								},
							},
						},
					}, nil,
				)
			},
			want: []model.Order{
				{
					Id:               1,
					Number:           1,
					TotalCost:        74000,
					CreatedDate:      time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
					LastModifiedDate: time.Date(2025, time.June, 15, 12, 00, 00, 0, time.UTC),
					Status:           1,
					Products: []model.Product{
						{
							Id:            1,
							Code:          14823,
							Quantity:      215,
							Name:          "Pizza",
							PurchasePrice: 24000,
							SalePrice:     74000,
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			os := &Service{Order: dbMock}

			tt.mock()

			got, err := os.Order.GetAll()
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
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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
				dbMock.EXPECT().GetById(1).Return(
					&model.Order{
						Id:               1,
						Number:           1,
						TotalCost:        74000,
						CreatedDate:      time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
						LastModifiedDate: time.Date(2025, time.June, 15, 12, 00, 00, 0, time.UTC),
						Status:           1,
						Products: []model.Product{
							{
								Id:            1,
								Code:          14823,
								Quantity:      215,
								Name:          "Pizza",
								PurchasePrice: 24000,
								SalePrice:     74000,
							},
						},
					}, nil,
				)
			},
			args: 1,
			want: &model.Order{
				Id:               1,
				Number:           1,
				TotalCost:        74000,
				CreatedDate:      time.Date(2025, time.May, 25, 12, 17, 16, 550631000, time.UTC),
				LastModifiedDate: time.Date(2025, time.June, 15, 12, 00, 00, 0, time.UTC),
				Status:           1,
				Products: []model.Product{
					{
						Id:            1,
						Code:          14823,
						Quantity:      215,
						Name:          "Pizza",
						PurchasePrice: 24000,
						SalePrice:     74000,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			mock: func() {
				dbMock.EXPECT().GetById(1).Return(
					nil,
					errors.New("элемент не найден"),
				)
			},
			args:    1,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			os := &Service{Order: dbMock}

			tt.mock()

			got, err := os.Order.GetById(tt.args)
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
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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
				dbMock.EXPECT().Delete(1).Return(nil)
			},
			args:    1,
			want:    nil,
			wantErr: false,
		},
		{
			name: "error",
			mock: func() {
				dbMock.EXPECT().Delete(1).Return(
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
			t.Parallel()

			os := &Service{Order: dbMock}

			tt.mock()

			err := os.Order.Delete(tt.args)
			if err != tt.want && err.Error() != tt.want.Error() {
				t.Errorf("Ошибка удаления заказа по id error = %v, want %v", err, tt.want)
				return
			}
		})
	}
}
