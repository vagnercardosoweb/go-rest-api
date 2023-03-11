package errors_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"net/http"
	"testing"
)

func TestCreateEmptyError(t *testing.T) {
	err := errors.New(errors.Input{})
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
