package errors

import (
	"github.com/stretchr/testify/assert"

	"net/http"
	"testing"
)

func TestCreateEmptyError(t *testing.T) {
	err := New(Input{})
	assert.NotNil(t, err)
	assert.NotNil(t, err.ErrorId)
	assert.Equal(t, err.Code, "InternalServerError")
	assert.Equal(t, err.Message, "InternalServerError")
	assert.Equal(t, err.Error(), "InternalServerError")
	assert.Equal(t, err.StatusCode, http.StatusInternalServerError)
	assert.NotEmpty(t, err.Metadata)
	assert.Equal(t, err.Logging, false)
	assert.Equal(t, err.SendToSlack, false)
	assert.Nil(t, err.Arguments)
}
