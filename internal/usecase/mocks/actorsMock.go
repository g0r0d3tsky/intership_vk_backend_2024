// Code generated by MockGen. DO NOT EDIT.
// Source: actors.go
//
// Generated by this command:
//
//	mockgen -source=actors.go -destination=mocks/actorsMock.go
//

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	domain "cinema_service/internal/domain"
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockActorsRepo is a mock of ActorsRepo interface.
type MockActorsRepo struct {
	ctrl     *gomock.Controller
	recorder *MockActorsRepoMockRecorder
}

// MockActorsRepoMockRecorder is the mock recorder for MockActorsRepo.
type MockActorsRepoMockRecorder struct {
	mock *MockActorsRepo
}

// NewMockActorsRepo creates a new mock instance.
func NewMockActorsRepo(ctrl *gomock.Controller) *MockActorsRepo {
	mock := &MockActorsRepo{ctrl: ctrl}
	mock.recorder = &MockActorsRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockActorsRepo) EXPECT() *MockActorsRepoMockRecorder {
	return m.recorder
}

// CreateActor mocks base method.
func (m *MockActorsRepo) CreateActor(ctx context.Context, act *domain.Actor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateActor", ctx, act)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateActor indicates an expected call of CreateActor.
func (mr *MockActorsRepoMockRecorder) CreateActor(ctx, act any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateActor", reflect.TypeOf((*MockActorsRepo)(nil).CreateActor), ctx, act)
}

// DeleteActor mocks base method.
func (m *MockActorsRepo) DeleteActor(ctx context.Context, actorID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteActor", ctx, actorID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteActor indicates an expected call of DeleteActor.
func (mr *MockActorsRepoMockRecorder) DeleteActor(ctx, actorID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteActor", reflect.TypeOf((*MockActorsRepo)(nil).DeleteActor), ctx, actorID)
}

// GetActors mocks base method.
func (m *MockActorsRepo) GetActors(ctx context.Context) (map[*domain.Actor][]*domain.Movie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActors", ctx)
	ret0, _ := ret[0].(map[*domain.Actor][]*domain.Movie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActors indicates an expected call of GetActors.
func (mr *MockActorsRepoMockRecorder) GetActors(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActors", reflect.TypeOf((*MockActorsRepo)(nil).GetActors), ctx)
}

// UpdateActor mocks base method.
func (m *MockActorsRepo) UpdateActor(ctx context.Context, act *domain.Actor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateActor", ctx, act)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateActor indicates an expected call of UpdateActor.
func (mr *MockActorsRepoMockRecorder) UpdateActor(ctx, act any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateActor", reflect.TypeOf((*MockActorsRepo)(nil).UpdateActor), ctx, act)
}
