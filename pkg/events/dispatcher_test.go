package events

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestEventHandler struct {
	ID int
}

func (e *TestEventHandler) Handle(event Event, wg *sync.WaitGroup) {}

type DispatcherTestSuite struct {
	suite.Suite
	event      Event
	event2     Event
	handler    TestEventHandler
	handler2   TestEventHandler
	handler3   TestEventHandler
	dispatcher DispatcherInterface
}

func (suite *DispatcherTestSuite) SetupTest() {
	suite.dispatcher = NewDispatcher()
	suite.handler = TestEventHandler{1}
	suite.handler2 = TestEventHandler{2}
	suite.handler3 = TestEventHandler{3}
	suite.event = Event{Name: "test", Payload: "test"}
	suite.event2 = Event{Name: "test2", Payload: "test2"}
}

func (suite *DispatcherTestSuite) TestDispatcher_Register() {
	err := suite.dispatcher.Register(suite.event.Name, &suite.handler)
	suite.Nil(err)
	suite.Equal(1, suite.dispatcher.Total(suite.event.Name))

	err = suite.dispatcher.Register(suite.event.Name, &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, suite.dispatcher.Total(suite.event.Name))

	assert.Equal(suite.T(), &suite.handler, suite.dispatcher.GetByIndex(suite.event.Name, 0))
	assert.Equal(suite.T(), &suite.handler2, suite.dispatcher.GetByIndex(suite.event.Name, 1))
}

func (suite *DispatcherTestSuite) TestDispatcher_Register_WithSameHandler() {
	err := suite.dispatcher.Register(suite.event.Name, &suite.handler)
	suite.Nil(err)
	suite.Equal(1, suite.dispatcher.Total(suite.event.Name))

	err = suite.dispatcher.Register(suite.event.Name, &suite.handler)
	suite.NotNil(err)
	suite.Equal(1, suite.dispatcher.Total(suite.event.Name))
}

func (suite *DispatcherTestSuite) TestDispatcher_Clear() {
	// Event 01
	err := suite.dispatcher.Register(suite.event.Name, &suite.handler)
	suite.Nil(err)
	suite.Equal(1, suite.dispatcher.Total(suite.event.Name))

	err = suite.dispatcher.Register(suite.event.Name, &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, suite.dispatcher.Total(suite.event.Name))

	// Event 02
	err = suite.dispatcher.Register(suite.event2.Name, &suite.handler3)
	suite.Nil(err)
	suite.Equal(1, suite.dispatcher.Total(suite.event2.Name))

	// Clear
	suite.dispatcher.Clear()
	suite.Equal(0, suite.dispatcher.Total(""))
}

func (suite *DispatcherTestSuite) TestDispatcher_Has() {
	err := suite.dispatcher.Register(suite.event.Name, &suite.handler)
	suite.Nil(err)
	suite.Equal(1, suite.dispatcher.Total(suite.event.Name))

	err = suite.dispatcher.Register(suite.event.Name, &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, suite.dispatcher.Total(suite.event.Name))

	assert.True(suite.T(), suite.dispatcher.Has(suite.event.Name, &suite.handler))
	assert.True(suite.T(), suite.dispatcher.Has(suite.event.Name, &suite.handler2))
	assert.False(suite.T(), suite.dispatcher.Has(suite.event.Name, &suite.handler3))
}

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Handle(event Event, wg *sync.WaitGroup) {
	defer wg.Done()
	m.Called(event)
}

func (suite *DispatcherTestSuite) TestDispatcher_Dispatch() {
	handler := &MockHandler{}
	handler.On("Handle", suite.event)
	suite.dispatcher.Register(suite.event.Name, handler)
	suite.dispatcher.Dispatch(suite.event)
	handler.AssertExpectations(suite.T())
	handler.AssertNumberOfCalls(suite.T(), "Handle", 1)
}

func (suite *DispatcherTestSuite) TestDispatcher_Remove() {
	err := suite.dispatcher.Register(suite.event.Name, &suite.handler)
	suite.Nil(err)
	suite.Equal(1, suite.dispatcher.Total(suite.event.Name))

	err = suite.dispatcher.Register(suite.event.Name, &suite.handler2)
	suite.Nil(err)
	suite.Equal(2, suite.dispatcher.Total(suite.event.Name))

	err = suite.dispatcher.Register(suite.event2.Name, &suite.handler3)
	suite.Nil(err)
	suite.Equal(1, suite.dispatcher.Total(suite.event2.Name))

	err = suite.dispatcher.Remove(suite.event2.Name, &suite.handler3)
	suite.Nil(err)
	suite.Equal(0, suite.dispatcher.Total(suite.event2.Name))

	err = suite.dispatcher.Remove(suite.event.Name, &suite.handler2)
	suite.Nil(err)
	suite.Equal(1, suite.dispatcher.Total(suite.event.Name))

	err = suite.dispatcher.Remove(suite.event.Name, &suite.handler)
	suite.Nil(err)
	suite.Equal(0, suite.dispatcher.Total(suite.event.Name))
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(DispatcherTestSuite))
}
