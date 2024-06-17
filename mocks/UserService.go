// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	models "chat-system/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// UserService is an autogenerated mock type for the UserService type
type UserService struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: userInput
func (_m *UserService) CreateUser(userInput *models.RegisterInput) (*models.User, error) {
	ret := _m.Called(userInput)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(*models.RegisterInput) (*models.User, error)); ok {
		return rf(userInput)
	}
	if rf, ok := ret.Get(0).(func(*models.RegisterInput) *models.User); ok {
		r0 = rf(userInput)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(*models.RegisterInput) error); ok {
		r1 = rf(userInput)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByCreds provides a mock function with given fields: credentials
func (_m *UserService) GetUserByCreds(credentials models.LoginInput) (*models.User, error) {
	ret := _m.Called(credentials)

	if len(ret) == 0 {
		panic("no return value specified for GetUserByCreds")
	}

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(models.LoginInput) (*models.User, error)); ok {
		return rf(credentials)
	}
	if rf, ok := ret.Get(0).(func(models.LoginInput) *models.User); ok {
		r0 = rf(credentials)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(models.LoginInput) error); ok {
		r1 = rf(credentials)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserExists provides a mock function with given fields: username
func (_m *UserService) UserExists(username string) (bool, error) {
	ret := _m.Called(username)

	if len(ret) == 0 {
		panic("no return value specified for UserExists")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (bool, error)); ok {
		return rf(username)
	}
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(username)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserService creates a new instance of UserService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserService(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserService {
	mock := &UserService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
