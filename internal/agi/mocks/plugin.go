// Code generated by mockery v2.22.1. DO NOT EDIT.

package mocks

import (
	agi "github.com/selesy/asdf-go-install/internal/agi"
	mock "github.com/stretchr/testify/mock"
)

// Plugin is an autogenerated mock type for the Plugin type
type Plugin struct {
	mock.Mock
}

type Plugin_Expecter struct {
	mock *mock.Mock
}

func (_m *Plugin) EXPECT() *Plugin_Expecter {
	return &Plugin_Expecter{mock: &_m.Mock}
}

// Download provides a mock function with given fields:
func (_m *Plugin) Download() agi.ExitCode {
	ret := _m.Called()

	var r0 agi.ExitCode
	if rf, ok := ret.Get(0).(func() agi.ExitCode); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(agi.ExitCode)
	}

	return r0
}

// Plugin_Download_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Download'
type Plugin_Download_Call struct {
	*mock.Call
}

// Download is a helper method to define mock.On call
func (_e *Plugin_Expecter) Download() *Plugin_Download_Call {
	return &Plugin_Download_Call{Call: _e.mock.On("Download")}
}

func (_c *Plugin_Download_Call) Run(run func()) *Plugin_Download_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Plugin_Download_Call) Return(_a0 agi.ExitCode) *Plugin_Download_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Plugin_Download_Call) RunAndReturn(run func() agi.ExitCode) *Plugin_Download_Call {
	_c.Call.Return(run)
	return _c
}

// HelpOverview provides a mock function with given fields:
func (_m *Plugin) HelpOverview() agi.ExitCode {
	ret := _m.Called()

	var r0 agi.ExitCode
	if rf, ok := ret.Get(0).(func() agi.ExitCode); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(agi.ExitCode)
	}

	return r0
}

// Plugin_HelpOverview_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HelpOverview'
type Plugin_HelpOverview_Call struct {
	*mock.Call
}

// HelpOverview is a helper method to define mock.On call
func (_e *Plugin_Expecter) HelpOverview() *Plugin_HelpOverview_Call {
	return &Plugin_HelpOverview_Call{Call: _e.mock.On("HelpOverview")}
}

func (_c *Plugin_HelpOverview_Call) Run(run func()) *Plugin_HelpOverview_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Plugin_HelpOverview_Call) Return(_a0 agi.ExitCode) *Plugin_HelpOverview_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Plugin_HelpOverview_Call) RunAndReturn(run func() agi.ExitCode) *Plugin_HelpOverview_Call {
	_c.Call.Return(run)
	return _c
}

// Install provides a mock function with given fields:
func (_m *Plugin) Install() agi.ExitCode {
	ret := _m.Called()

	var r0 agi.ExitCode
	if rf, ok := ret.Get(0).(func() agi.ExitCode); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(agi.ExitCode)
	}

	return r0
}

// Plugin_Install_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Install'
type Plugin_Install_Call struct {
	*mock.Call
}

// Install is a helper method to define mock.On call
func (_e *Plugin_Expecter) Install() *Plugin_Install_Call {
	return &Plugin_Install_Call{Call: _e.mock.On("Install")}
}

func (_c *Plugin_Install_Call) Run(run func()) *Plugin_Install_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Plugin_Install_Call) Return(_a0 agi.ExitCode) *Plugin_Install_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Plugin_Install_Call) RunAndReturn(run func() agi.ExitCode) *Plugin_Install_Call {
	_c.Call.Return(run)
	return _c
}

// ListAll provides a mock function with given fields:
func (_m *Plugin) ListAll() agi.ExitCode {
	ret := _m.Called()

	var r0 agi.ExitCode
	if rf, ok := ret.Get(0).(func() agi.ExitCode); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(agi.ExitCode)
	}

	return r0
}

// Plugin_ListAll_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListAll'
type Plugin_ListAll_Call struct {
	*mock.Call
}

// ListAll is a helper method to define mock.On call
func (_e *Plugin_Expecter) ListAll() *Plugin_ListAll_Call {
	return &Plugin_ListAll_Call{Call: _e.mock.On("ListAll")}
}

func (_c *Plugin_ListAll_Call) Run(run func()) *Plugin_ListAll_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Plugin_ListAll_Call) Return(_a0 agi.ExitCode) *Plugin_ListAll_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Plugin_ListAll_Call) RunAndReturn(run func() agi.ExitCode) *Plugin_ListAll_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewPlugin interface {
	mock.TestingT
	Cleanup(func())
}

// NewPlugin creates a new instance of Plugin. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPlugin(t mockConstructorTestingTNewPlugin) *Plugin {
	mock := &Plugin{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
