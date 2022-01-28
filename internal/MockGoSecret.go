package internal

import "github.com/stretchr/testify/mock"

type MockGoSecret struct {
	mock.Mock
}

func (m *MockGoSecret) Bytes() []byte {
	retVal := m.Called().Get(0)
	if retVal == nil {
		return nil
	}
	return retVal.([]byte)
}

func (m *MockGoSecret) Keys() []string {
	retVal := m.Called().Get(0)
	if retVal == nil {
		return nil
	}
	return retVal.([]string)
}

func (m *MockGoSecret) Get(key string) (string, bool) {
	args := m.Called(key)
	return args.String(0), args.Bool(1)
}

func (m *MockGoSecret) Values(key string) ([]string, bool) {
	args := m.Called(key)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).([]string), args.Bool(1)
}

func (m *MockGoSecret) Set(key string, value interface{}) error {
	return m.Called(key, value).Error(0)
}

func (m *MockGoSecret) Add(key string, value interface{}) error {
	return m.Called(key, value).Error(0)
}

func (m *MockGoSecret) Del(key string) bool {
	return m.Called(key).Bool(0)
}

func (m *MockGoSecret) Body() string {
	return m.Called().String(0)
}

func (m *MockGoSecret) Password() string {
	return m.Called().String(0)
}

func (m *MockGoSecret) SetPassword(s string) {
	m.Called(s)
}
