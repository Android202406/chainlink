// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	functions "github.com/smartcontractkit/chainlink/v2/core/services/functions"
	mock "github.com/stretchr/testify/mock"
)

// OffchainTransmitter is an autogenerated mock type for the OffchainTransmitter type
type OffchainTransmitter struct {
	mock.Mock
}

// ReportChannel provides a mock function with given fields:
func (_m *OffchainTransmitter) ReportChannel() chan *functions.OffchainResponse {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ReportChannel")
	}

	var r0 chan *functions.OffchainResponse
	if rf, ok := ret.Get(0).(func() chan *functions.OffchainResponse); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan *functions.OffchainResponse)
		}
	}

	return r0
}

// TransmitReport provides a mock function with given fields: ctx, report
func (_m *OffchainTransmitter) TransmitReport(ctx context.Context, report *functions.OffchainResponse) error {
	ret := _m.Called(ctx, report)

	if len(ret) == 0 {
		panic("no return value specified for TransmitReport")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *functions.OffchainResponse) error); ok {
		r0 = rf(ctx, report)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewOffchainTransmitter creates a new instance of OffchainTransmitter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOffchainTransmitter(t interface {
	mock.TestingT
	Cleanup(func())
}) *OffchainTransmitter {
	mock := &OffchainTransmitter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
