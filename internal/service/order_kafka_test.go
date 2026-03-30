package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mikhailshtv/stockLkBack/internal/model"
	repoMocks "github.com/mikhailshtv/stockLkBack/internal/repository/mocks"
	producerMocks "github.com/mikhailshtv/stockLkBack/pkg/kafka/producer/mocks"

	"github.com/golang/mock/gomock"
)

var (
	kafkaTestOrder = &model.Order{
		ID:        1,
		Number:    1,
		TotalCost: 74000,
		UserID:    1,
		Status:    model.StatusActive,
	}
	kafkaTestUser = &model.User{
		ID:        1,
		FirstName: "Ivan",
		LastName:  "Petrov",
		Email:     "ivan@test.com",
	}
	kafkaTestOrderBody = model.OrderRequestBody{
		Products: []model.OrderProduct{
			{ProductID: 1, Quantity: 1, SellPrice: 74000},
		},
	}
	kafkaTestStatusReq = model.OrderStatusRequest{
		Status: model.OrderStatus{Key: "executed", DisplayName: "Выполнен"},
	}
)

func waitPublished(t *testing.T, published <-chan struct{}) {
	t.Helper()
	select {
	case <-published:
	case <-time.After(time.Second):
		t.Error("publisher.Publish не был вызван в течение 1 секунды")
	}
}

func newOrdersServiceMocks(ctrl *gomock.Controller) (
	*OrdersService,
	*repoMocks.MockOrder,
	*repoMocks.MockUser,
	*producerMocks.MockEventPublisher,
) {
	orderRepo := repoMocks.NewMockOrder(ctrl)
	userRepo := repoMocks.NewMockUser(ctrl)
	publisher := producerMocks.NewMockEventPublisher(ctrl)
	svc := NewOrdersService(context.Background(), orderRepo, userRepo, publisher)
	return svc, orderRepo, userRepo, publisher
}

func expectPublishEvent(
	t *testing.T,
	userRepo *repoMocks.MockUser,
	publisher *producerMocks.MockEventPublisher,
	userID int,
) <-chan struct{} {
	t.Helper()
	userRepo.EXPECT().
		GetByID(gomock.Any(), userID).
		Return(kafkaTestUser, nil)

	published := make(chan struct{}, 1)
	publisher.EXPECT().
		Publish(gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, _ string, _ any) {
			published <- struct{}{}
		}).
		Return(nil)
	return published
}

func TestOrdersService_Create_PublishesEventOnSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	svc, orderRepo, userRepo, publisher := newOrdersServiceMocks(ctrl)

	orderRepo.EXPECT().
		Create(gomock.Any(), kafkaTestOrderBody, 1).
		Return(kafkaTestOrder, nil)
	orderRepo.EXPECT().
		WriteLog(gomock.Any(), "Create", logSuccessStatus, logOrdersTableName).
		Return(int64(1), nil)

	published := expectPublishEvent(t, userRepo, publisher, kafkaTestOrder.UserID)

	got, err := svc.Create(kafkaTestOrderBody, 1)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if got == nil {
		t.Fatal("ожидался заказ, получен nil")
	}

	waitPublished(t, published)
}

func TestOrdersService_Create_DoesNotPublishOnRepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	svc, orderRepo, _, _ := newOrdersServiceMocks(ctrl)

	orderRepo.EXPECT().
		Create(gomock.Any(), kafkaTestOrderBody, 1).
		Return(nil, errors.New("ошибка БД"))
	orderRepo.EXPECT().
		WriteLog(gomock.Any(), "Create", logErrorStatus, logOrdersTableName).
		Return(int64(1), nil)

	_, err := svc.Create(kafkaTestOrderBody, 1)
	if err == nil {
		t.Fatal("ожидалась ошибка")
	}
}

func TestOrdersService_Update_PublishesEventOnSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	svc, orderRepo, userRepo, publisher := newOrdersServiceMocks(ctrl)

	orderRepo.EXPECT().
		Update(gomock.Any(), 1, kafkaTestOrderBody, 1).
		Return(kafkaTestOrder, nil)
	orderRepo.EXPECT().
		WriteLog(gomock.Any(), "Update", logSuccessStatus, logOrdersTableName).
		Return(int64(1), nil)

	published := expectPublishEvent(t, userRepo, publisher, kafkaTestOrder.UserID)

	got, err := svc.Update(1, kafkaTestOrderBody, 1)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if got == nil {
		t.Fatal("ожидался заказ, получен nil")
	}

	waitPublished(t, published)
}

func TestOrdersService_Update_DoesNotPublishOnRepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	svc, orderRepo, _, _ := newOrdersServiceMocks(ctrl)

	orderRepo.EXPECT().
		Update(gomock.Any(), 1, kafkaTestOrderBody, 1).
		Return(nil, errors.New("заказ не найден"))
	orderRepo.EXPECT().
		WriteLog(gomock.Any(), "Update", logErrorStatus, logOrdersTableName).
		Return(int64(1), nil)

	_, err := svc.Update(1, kafkaTestOrderBody, 1)
	if err == nil {
		t.Fatal("ожидалась ошибка")
	}
}

func TestOrdersService_UpdateStatus_PublishesEventOnSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	svc, orderRepo, userRepo, publisher := newOrdersServiceMocks(ctrl)

	orderRepo.EXPECT().
		UpdateStatus(gomock.Any(), 1, kafkaTestStatusReq, 1).
		Return(kafkaTestOrder, nil)
	orderRepo.EXPECT().
		WriteLog(gomock.Any(), "UpdateStatus", logSuccessStatus, logOrdersTableName).
		Return(int64(1), nil)

	published := expectPublishEvent(t, userRepo, publisher, kafkaTestOrder.UserID)

	got, err := svc.UpdateStatus(1, kafkaTestStatusReq, 1)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if got == nil {
		t.Fatal("ожидался заказ, получен nil")
	}

	waitPublished(t, published)
}

func TestOrdersService_UpdateStatus_DoesNotPublishOnRepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	svc, orderRepo, _, _ := newOrdersServiceMocks(ctrl)

	orderRepo.EXPECT().
		UpdateStatus(gomock.Any(), 1, kafkaTestStatusReq, 1).
		Return(nil, errors.New("заказ не найден"))
	orderRepo.EXPECT().
		WriteLog(gomock.Any(), "UpdateStatus", logErrorStatus, logOrdersTableName).
		Return(int64(1), nil)

	_, err := svc.UpdateStatus(1, kafkaTestStatusReq, 1)
	if err == nil {
		t.Fatal("ожидалась ошибка")
	}
}

func TestOrdersService_Delete_NeverPublishesEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	svc, orderRepo, _, _ := newOrdersServiceMocks(ctrl)

	orderRepo.EXPECT().
		Delete(gomock.Any(), 1, 1).
		Return(kafkaTestOrder, nil)
	orderRepo.EXPECT().
		WriteLog(gomock.Any(), "Delete", logSuccessStatus, logOrdersTableName).
		Return(int64(1), nil)

	err := svc.Delete(1, 1)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
}

func TestOrdersService_PublishEvent_SkipsPublishOnUserRepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() { ctrl.Finish() })

	svc, orderRepo, userRepo, _ := newOrdersServiceMocks(ctrl)

	orderRepo.EXPECT().
		Create(gomock.Any(), kafkaTestOrderBody, 1).
		Return(kafkaTestOrder, nil)
	orderRepo.EXPECT().
		WriteLog(gomock.Any(), "Create", logSuccessStatus, logOrdersTableName).
		Return(int64(1), nil)

	userRepoCallDone := make(chan struct{}, 1)
	userRepo.EXPECT().
		GetByID(gomock.Any(), kafkaTestOrder.UserID).
		Do(func(_ context.Context, _ int) {
			userRepoCallDone <- struct{}{}
		}).
		Return(nil, errors.New("пользователь не найден"))

	_, err := svc.Create(kafkaTestOrderBody, 1)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	select {
	case <-userRepoCallDone:
	case <-time.After(time.Second):
		t.Error("userRepo.GetByID не был вызван в течение 1 секунды")
	}
}
