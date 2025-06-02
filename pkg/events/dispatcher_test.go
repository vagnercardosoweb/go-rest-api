package events

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestHandler struct{ ID int }

func (e *TestHandler) Handle(event *Event) error {
	return nil
}

type DispatcherSuite struct {
	suite.Suite
	event      *Event
	event2     *Event
	handler    *TestHandler
	handler2   *TestHandler
	handler3   *TestHandler
	dispatcher DispatcherInterface
}

func (suite *DispatcherSuite) SetupTest() {
	suite.dispatcher = NewDispatcher()

	suite.handler = &TestHandler{1}
	suite.handler2 = &TestHandler{2}
	suite.handler3 = &TestHandler{3}

	suite.event = &Event{Name: "test", Input: "test"}
	suite.event2 = &Event{Name: "test2", Input: "test2"}
}

func (t *DispatcherSuite) TestRegisterHandlers() {
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

func (t *DispatcherSuite) TestClearHandlers() {
	err := t.dispatcher.Register(t.event.Name, t.handler)
	t.Nil(err)
	t.Equal(1, t.dispatcher.Total(t.event.Name))

	err = t.dispatcher.Register(t.event.Name, t.handler2)
	t.Nil(err)
	t.Equal(2, t.dispatcher.Total(t.event.Name))

	err = t.dispatcher.Register(t.event2.Name, t.handler3)
	t.Nil(err)
	t.Equal(1, t.dispatcher.Total(t.event2.Name))

	t.dispatcher.Clear()
	t.Equal(0, t.dispatcher.Total(""))
}

func (t *DispatcherSuite) TestHasHandlers() {
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

func (m *MockHandler) Handle(event *Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func (t *DispatcherSuite) TestDispatchHandlers() {
	h := new(MockHandler)
	h.On("Handle", t.event).
		Return(nil).
		Once()

	err := t.dispatcher.Register(t.event.Name, h)
	t.Nil(err)

	err = t.dispatcher.Dispatch(t.event)
	t.Nil(err)

	h.AssertExpectations(t.T())
}

func (t *DispatcherSuite) TestDispatchError() {
	h := new(MockHandler)
	h.On("Handle", t.event).
		Return(errors.New("test")).
		Once()

	err := t.dispatcher.Register(t.event.Name, h)
	t.Nil(err)

	err = t.dispatcher.Dispatch(t.event)
	t.NotNil(err)

	h.AssertExpectations(t.T())
}

func (t *DispatcherSuite) TestRemoveHandlers() {
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
