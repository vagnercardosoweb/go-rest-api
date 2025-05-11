package events

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestEventHandler struct{ ID int }

func (e *TestEventHandler) Handle(event *Event) {}

type DispatcherSuite struct {
	suite.Suite
	event      *Event
	event2     *Event
	handler    *TestEventHandler
	handler2   *TestEventHandler
	handler3   *TestEventHandler
	dispatcher DispatcherInterface
}

func (suite *DispatcherSuite) SetupTest() {
	suite.dispatcher = NewDispatcher()
	suite.handler = &TestEventHandler{1}
	suite.handler2 = &TestEventHandler{2}
	suite.handler3 = &TestEventHandler{3}
	suite.event = &Event{Name: "test", Input: "test"}
	suite.event2 = &Event{Name: "test2", Input: "test2"}
}

func (t *DispatcherSuite) TestRegisterEventHandlers() {
	err := t.dispatcher.Register(t.event.Name, t.handler)
	t.Nil(err)
	t.Equal(1, t.dispatcher.Total(t.event.Name))

	err = t.dispatcher.Register(t.event.Name, t.handler2)
	t.Nil(err)
	t.Equal(2, t.dispatcher.Total(t.event.Name))

	t.Equal(t.handler, t.dispatcher.GetByIndex(t.event.Name, 0))
	t.Equal(t.handler2, t.dispatcher.GetByIndex(t.event.Name, 1))
}

func (t *DispatcherSuite) TestRegisterWithSameHandler() {
	err := t.dispatcher.Register(t.event.Name, t.handler)
	t.Nil(err)
	t.Equal(1, t.dispatcher.Total(t.event.Name))

	err = t.dispatcher.Register(t.event.Name, t.handler)
	t.NotNil(err)
	t.Equal(1, t.dispatcher.Total(t.event.Name))
}

func (t *DispatcherSuite) TestClearEventHandlers() {
	// Event 01
	err := t.dispatcher.Register(t.event.Name, t.handler)
	t.Nil(err)
	t.Equal(1, t.dispatcher.Total(t.event.Name))

	err = t.dispatcher.Register(t.event.Name, t.handler2)
	t.Nil(err)
	t.Equal(2, t.dispatcher.Total(t.event.Name))

	// Event 02
	err = t.dispatcher.Register(t.event2.Name, t.handler3)
	t.Nil(err)
	t.Equal(1, t.dispatcher.Total(t.event2.Name))

	// Clear
	t.dispatcher.Clear()
	t.Equal(0, t.dispatcher.Total(""))
}

func (t *DispatcherSuite) TestHasEventHandlers() {
	err := t.dispatcher.Register(t.event.Name, t.handler)
	t.Nil(err)
	t.Equal(1, t.dispatcher.Total(t.event.Name))

	err = t.dispatcher.Register(t.event.Name, t.handler2)
	t.Nil(err)
	t.Equal(2, t.dispatcher.Total(t.event.Name))

	t.True(t.dispatcher.Has(t.event.Name, t.handler))
	t.True(t.dispatcher.Has(t.event.Name, t.handler2))
	t.False(t.dispatcher.Has(t.event.Name, t.handler3))
}

type MockHandler struct{ mock.Mock }

func (m *MockHandler) Handle(event *Event) {
	m.Called(event)
}

func (t *DispatcherSuite) TestDispatchEventRegistration() {
	handler := &MockHandler{}
	handler.On("Handle", t.event)

	_ = t.dispatcher.Register(t.event.Name, handler)
	_ = t.dispatcher.Dispatch(t.event)

	handler.AssertExpectations(t.T())
	handler.AssertNumberOfCalls(t.T(), "Handle", 1)
}

func (t *DispatcherSuite) TestRemoveEventHandlers() {
	err := t.dispatcher.Register(t.event.Name, t.handler)
	t.Nil(err)
	t.Equal(1, t.dispatcher.Total(t.event.Name))

	err = t.dispatcher.Register(t.event.Name, t.handler2)
	t.Nil(err)
	t.Equal(2, t.dispatcher.Total(t.event.Name))

	err = t.dispatcher.Register(t.event2.Name, t.handler3)
	t.Nil(err)
	t.Equal(1, t.dispatcher.Total(t.event2.Name))

	err = t.dispatcher.Remove(t.event2.Name, t.handler3)
	t.Nil(err)
	t.Equal(0, t.dispatcher.Total(t.event2.Name))

	err = t.dispatcher.Remove(t.event.Name, t.handler2)
	t.Nil(err)
	t.Equal(1, t.dispatcher.Total(t.event.Name))

	err = t.dispatcher.Remove(t.event.Name, t.handler)
	t.Nil(err)
	t.Equal(0, t.dispatcher.Total(t.event.Name))
}

func TestDispatcherSuite(t *testing.T) {
	suite.Run(t, new(DispatcherSuite))
}
