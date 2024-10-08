// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/Azure/ARO-RP/pkg/metrics (interfaces: Emitter)
//
// Generated by this command:
//
//	mockgen -destination=../util/mocks/metrics/metrics.go github.com/Azure/ARO-RP/pkg/metrics Emitter
//

// Package mock_metrics is a generated GoMock package.
package mock_metrics

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockEmitter is a mock of Emitter interface.
type MockEmitter struct {
	ctrl     *gomock.Controller
	recorder *MockEmitterMockRecorder
}

// MockEmitterMockRecorder is the mock recorder for MockEmitter.
type MockEmitterMockRecorder struct {
	mock *MockEmitter
}

// NewMockEmitter creates a new mock instance.
func NewMockEmitter(ctrl *gomock.Controller) *MockEmitter {
	mock := &MockEmitter{ctrl: ctrl}
	mock.recorder = &MockEmitterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmitter) EXPECT() *MockEmitterMockRecorder {
	return m.recorder
}

// EmitFloat mocks base method.
func (m *MockEmitter) EmitFloat(arg0 string, arg1 float64, arg2 map[string]string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "EmitFloat", arg0, arg1, arg2)
}

// EmitFloat indicates an expected call of EmitFloat.
func (mr *MockEmitterMockRecorder) EmitFloat(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EmitFloat", reflect.TypeOf((*MockEmitter)(nil).EmitFloat), arg0, arg1, arg2)
}

// EmitGauge mocks base method.
func (m *MockEmitter) EmitGauge(arg0 string, arg1 int64, arg2 map[string]string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "EmitGauge", arg0, arg1, arg2)
}

// EmitGauge indicates an expected call of EmitGauge.
func (mr *MockEmitterMockRecorder) EmitGauge(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EmitGauge", reflect.TypeOf((*MockEmitter)(nil).EmitGauge), arg0, arg1, arg2)
}
