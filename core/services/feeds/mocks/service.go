// Code generated by mockery v2.42.2. DO NOT EDIT.

package mocks

import (
	context "context"

	feeds "github.com/smartcontractkit/chainlink/v2/core/services/feeds"
	mock "github.com/stretchr/testify/mock"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// ApproveSpec provides a mock function with given fields: ctx, id, force
func (_m *Service) ApproveSpec(ctx context.Context, id int64, force bool) error {
	ret := _m.Called(ctx, id, force)

	if len(ret) == 0 {
		panic("no return value specified for ApproveSpec")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, bool) error); ok {
		r0 = rf(ctx, id, force)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CancelSpec provides a mock function with given fields: ctx, id
func (_m *Service) CancelSpec(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for CancelSpec")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Close provides a mock function with given fields:
func (_m *Service) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CountJobProposalsByStatus provides a mock function with given fields:
func (_m *Service) CountJobProposalsByStatus() (*feeds.JobProposalCounts, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for CountJobProposalsByStatus")
	}

	var r0 *feeds.JobProposalCounts
	var r1 error
	if rf, ok := ret.Get(0).(func() (*feeds.JobProposalCounts, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() *feeds.JobProposalCounts); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*feeds.JobProposalCounts)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CountManagers provides a mock function with given fields:
func (_m *Service) CountManagers() (int64, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for CountManagers")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func() (int64, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateChainConfig provides a mock function with given fields: ctx, cfg
func (_m *Service) CreateChainConfig(ctx context.Context, cfg feeds.ChainConfig) (int64, error) {
	ret := _m.Called(ctx, cfg)

	if len(ret) == 0 {
		panic("no return value specified for CreateChainConfig")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, feeds.ChainConfig) (int64, error)); ok {
		return rf(ctx, cfg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, feeds.ChainConfig) int64); ok {
		r0 = rf(ctx, cfg)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, feeds.ChainConfig) error); ok {
		r1 = rf(ctx, cfg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteChainConfig provides a mock function with given fields: ctx, id
func (_m *Service) DeleteChainConfig(ctx context.Context, id int64) (int64, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteChainConfig")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (int64, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) int64); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteJob provides a mock function with given fields: ctx, args
func (_m *Service) DeleteJob(ctx context.Context, args *feeds.DeleteJobArgs) (int64, error) {
	ret := _m.Called(ctx, args)

	if len(ret) == 0 {
		panic("no return value specified for DeleteJob")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *feeds.DeleteJobArgs) (int64, error)); ok {
		return rf(ctx, args)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *feeds.DeleteJobArgs) int64); ok {
		r0 = rf(ctx, args)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *feeds.DeleteJobArgs) error); ok {
		r1 = rf(ctx, args)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetChainConfig provides a mock function with given fields: id
func (_m *Service) GetChainConfig(id int64) (*feeds.ChainConfig, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetChainConfig")
	}

	var r0 *feeds.ChainConfig
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) (*feeds.ChainConfig, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int64) *feeds.ChainConfig); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*feeds.ChainConfig)
		}
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetJobProposal provides a mock function with given fields: id
func (_m *Service) GetJobProposal(id int64) (*feeds.JobProposal, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetJobProposal")
	}

	var r0 *feeds.JobProposal
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) (*feeds.JobProposal, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int64) *feeds.JobProposal); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*feeds.JobProposal)
		}
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetManager provides a mock function with given fields: id
func (_m *Service) GetManager(id int64) (*feeds.FeedsManager, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetManager")
	}

	var r0 *feeds.FeedsManager
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) (*feeds.FeedsManager, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int64) *feeds.FeedsManager); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*feeds.FeedsManager)
		}
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSpec provides a mock function with given fields: id
func (_m *Service) GetSpec(id int64) (*feeds.JobProposalSpec, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetSpec")
	}

	var r0 *feeds.JobProposalSpec
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) (*feeds.JobProposalSpec, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int64) *feeds.JobProposalSpec); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*feeds.JobProposalSpec)
		}
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsJobManaged provides a mock function with given fields: ctx, jobID
func (_m *Service) IsJobManaged(ctx context.Context, jobID int64) (bool, error) {
	ret := _m.Called(ctx, jobID)

	if len(ret) == 0 {
		panic("no return value specified for IsJobManaged")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (bool, error)); ok {
		return rf(ctx, jobID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) bool); ok {
		r0 = rf(ctx, jobID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, jobID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListChainConfigsByManagerIDs provides a mock function with given fields: mgrIDs
func (_m *Service) ListChainConfigsByManagerIDs(mgrIDs []int64) ([]feeds.ChainConfig, error) {
	ret := _m.Called(mgrIDs)

	if len(ret) == 0 {
		panic("no return value specified for ListChainConfigsByManagerIDs")
	}

	var r0 []feeds.ChainConfig
	var r1 error
	if rf, ok := ret.Get(0).(func([]int64) ([]feeds.ChainConfig, error)); ok {
		return rf(mgrIDs)
	}
	if rf, ok := ret.Get(0).(func([]int64) []feeds.ChainConfig); ok {
		r0 = rf(mgrIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]feeds.ChainConfig)
		}
	}

	if rf, ok := ret.Get(1).(func([]int64) error); ok {
		r1 = rf(mgrIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListJobProposals provides a mock function with given fields:
func (_m *Service) ListJobProposals() ([]feeds.JobProposal, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ListJobProposals")
	}

	var r0 []feeds.JobProposal
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]feeds.JobProposal, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []feeds.JobProposal); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]feeds.JobProposal)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListJobProposalsByManagersIDs provides a mock function with given fields: ids
func (_m *Service) ListJobProposalsByManagersIDs(ids []int64) ([]feeds.JobProposal, error) {
	ret := _m.Called(ids)

	if len(ret) == 0 {
		panic("no return value specified for ListJobProposalsByManagersIDs")
	}

	var r0 []feeds.JobProposal
	var r1 error
	if rf, ok := ret.Get(0).(func([]int64) ([]feeds.JobProposal, error)); ok {
		return rf(ids)
	}
	if rf, ok := ret.Get(0).(func([]int64) []feeds.JobProposal); ok {
		r0 = rf(ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]feeds.JobProposal)
		}
	}

	if rf, ok := ret.Get(1).(func([]int64) error); ok {
		r1 = rf(ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListManagers provides a mock function with given fields:
func (_m *Service) ListManagers() ([]feeds.FeedsManager, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ListManagers")
	}

	var r0 []feeds.FeedsManager
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]feeds.FeedsManager, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []feeds.FeedsManager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]feeds.FeedsManager)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListManagersByIDs provides a mock function with given fields: ids
func (_m *Service) ListManagersByIDs(ids []int64) ([]feeds.FeedsManager, error) {
	ret := _m.Called(ids)

	if len(ret) == 0 {
		panic("no return value specified for ListManagersByIDs")
	}

	var r0 []feeds.FeedsManager
	var r1 error
	if rf, ok := ret.Get(0).(func([]int64) ([]feeds.FeedsManager, error)); ok {
		return rf(ids)
	}
	if rf, ok := ret.Get(0).(func([]int64) []feeds.FeedsManager); ok {
		r0 = rf(ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]feeds.FeedsManager)
		}
	}

	if rf, ok := ret.Get(1).(func([]int64) error); ok {
		r1 = rf(ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListSpecsByJobProposalIDs provides a mock function with given fields: ids
func (_m *Service) ListSpecsByJobProposalIDs(ids []int64) ([]feeds.JobProposalSpec, error) {
	ret := _m.Called(ids)

	if len(ret) == 0 {
		panic("no return value specified for ListSpecsByJobProposalIDs")
	}

	var r0 []feeds.JobProposalSpec
	var r1 error
	if rf, ok := ret.Get(0).(func([]int64) ([]feeds.JobProposalSpec, error)); ok {
		return rf(ids)
	}
	if rf, ok := ret.Get(0).(func([]int64) []feeds.JobProposalSpec); ok {
		r0 = rf(ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]feeds.JobProposalSpec)
		}
	}

	if rf, ok := ret.Get(1).(func([]int64) error); ok {
		r1 = rf(ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProposeJob provides a mock function with given fields: ctx, args
func (_m *Service) ProposeJob(ctx context.Context, args *feeds.ProposeJobArgs) (int64, error) {
	ret := _m.Called(ctx, args)

	if len(ret) == 0 {
		panic("no return value specified for ProposeJob")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *feeds.ProposeJobArgs) (int64, error)); ok {
		return rf(ctx, args)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *feeds.ProposeJobArgs) int64); ok {
		r0 = rf(ctx, args)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *feeds.ProposeJobArgs) error); ok {
		r1 = rf(ctx, args)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterManager provides a mock function with given fields: ctx, params
func (_m *Service) RegisterManager(ctx context.Context, params feeds.RegisterManagerParams) (int64, error) {
	ret := _m.Called(ctx, params)

	if len(ret) == 0 {
		panic("no return value specified for RegisterManager")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, feeds.RegisterManagerParams) (int64, error)); ok {
		return rf(ctx, params)
	}
	if rf, ok := ret.Get(0).(func(context.Context, feeds.RegisterManagerParams) int64); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, feeds.RegisterManagerParams) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RejectSpec provides a mock function with given fields: ctx, id
func (_m *Service) RejectSpec(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for RejectSpec")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RevokeJob provides a mock function with given fields: ctx, args
func (_m *Service) RevokeJob(ctx context.Context, args *feeds.RevokeJobArgs) (int64, error) {
	ret := _m.Called(ctx, args)

	if len(ret) == 0 {
		panic("no return value specified for RevokeJob")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *feeds.RevokeJobArgs) (int64, error)); ok {
		return rf(ctx, args)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *feeds.RevokeJobArgs) int64); ok {
		r0 = rf(ctx, args)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *feeds.RevokeJobArgs) error); ok {
		r1 = rf(ctx, args)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Start provides a mock function with given fields: ctx
func (_m *Service) Start(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Start")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SyncNodeInfo provides a mock function with given fields: ctx, id
func (_m *Service) SyncNodeInfo(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for SyncNodeInfo")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Unsafe_SetConnectionsManager provides a mock function with given fields: _a0
func (_m *Service) Unsafe_SetConnectionsManager(_a0 feeds.ConnectionsManager) {
	_m.Called(_a0)
}

// UpdateChainConfig provides a mock function with given fields: ctx, cfg
func (_m *Service) UpdateChainConfig(ctx context.Context, cfg feeds.ChainConfig) (int64, error) {
	ret := _m.Called(ctx, cfg)

	if len(ret) == 0 {
		panic("no return value specified for UpdateChainConfig")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, feeds.ChainConfig) (int64, error)); ok {
		return rf(ctx, cfg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, feeds.ChainConfig) int64); ok {
		r0 = rf(ctx, cfg)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, feeds.ChainConfig) error); ok {
		r1 = rf(ctx, cfg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateManager provides a mock function with given fields: ctx, mgr
func (_m *Service) UpdateManager(ctx context.Context, mgr feeds.FeedsManager) error {
	ret := _m.Called(ctx, mgr)

	if len(ret) == 0 {
		panic("no return value specified for UpdateManager")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, feeds.FeedsManager) error); ok {
		r0 = rf(ctx, mgr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateSpecDefinition provides a mock function with given fields: ctx, id, spec
func (_m *Service) UpdateSpecDefinition(ctx context.Context, id int64, spec string) error {
	ret := _m.Called(ctx, id, spec)

	if len(ret) == 0 {
		panic("no return value specified for UpdateSpecDefinition")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, string) error); ok {
		r0 = rf(ctx, id, spec)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewService(t interface {
	mock.TestingT
	Cleanup(func())
}) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
