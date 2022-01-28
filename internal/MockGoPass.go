package internal

import (
	"context"
	"github.com/gopasspw/gopass/pkg/gopass"
	"github.com/stretchr/testify/mock"
)

type MockGoPass struct {
	mock.Mock
}

func (m *MockGoPass) String() string {
	return m.Called().String(0)
}

func (m *MockGoPass) List(ctx context.Context) ([]string, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockGoPass) Get(ctx context.Context, name, revision string) (gopass.Secret, error) {
	args := m.Called(ctx, name, revision)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(gopass.Secret), args.Error(1)

}

func (m *MockGoPass) Set(ctx context.Context, name string, sec gopass.Byter) error {
	return m.Called(ctx, name, sec).Error(0)
}

func (m *MockGoPass) Revisions(ctx context.Context, name string) ([]string, error) {
	args := m.Called(ctx, name)
	if args.Get(0) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockGoPass) Remove(ctx context.Context, name string) error {
	return m.Called(ctx, name).Error(0)
}

func (m *MockGoPass) RemoveAll(ctx context.Context, prefix string) error {
	return m.Called(ctx, prefix).Error(0)
}

func (m *MockGoPass) Rename(ctx context.Context, src, dest string) error {
	return m.Called(ctx, src, dest).Error(0)
}

func (m *MockGoPass) Sync(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func (m *MockGoPass) Close(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}
