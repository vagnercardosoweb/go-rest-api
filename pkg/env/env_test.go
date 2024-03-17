package env

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAsString(t *testing.T) {
	const expectedValue = "any"
	const expectedDefault = "any_default"

	t.Setenv("ANY", expectedValue)
	assert.Equal(t, expectedValue, GetAsString("ANY", expectedDefault))

	assert.Equal(t, expectedDefault, GetAsString("ANY_NOT_EXIST", expectedDefault))
}

func TestGetAsInt(t *testing.T) {
	const expectedValue = 100
	const expectedDefault = "200"

	t.Setenv("ANY", "100")
	assert.Equal(t, expectedValue, GetAsInt("ANY", expectedDefault))

	expectedDefaultToInt, _ := strconv.Atoi(expectedDefault)
	assert.Equal(t, expectedDefaultToInt, GetAsInt("ANY_NOT_EXIST", expectedDefault))
}

func TestGetAsBool(t *testing.T) {
	const expectedValue = true
	const expectedDefault = "false"

	t.Setenv("ANY", "true")
	assert.Equal(t, expectedValue, GetAsBool("ANY", expectedDefault))

	expectedDefaultToBool, _ := strconv.ParseBool(expectedDefault)
	assert.Equal(t, expectedDefaultToBool, GetAsBool("ANY_NOT_EXIST", expectedDefault))
}

func TestGetAsFloat64(t *testing.T) {
	const expectedValue = 100.0
	const expectedDefault = "200.0"

	t.Setenv("ANY", "100.0")
	assert.Equal(t, expectedValue, GetAsFloat64("ANY", expectedDefault))

	expectedDefaultToFloat64, _ := strconv.ParseFloat(expectedDefault, 64)
	assert.Equal(t, expectedDefaultToFloat64, GetAsFloat64("ANY_NOT_EXIST", expectedDefault))
}

func TestGetRequired(t *testing.T) {
	t.Setenv("ANY", "any")
	assert.Equal(t, "any", Required("ANY"))
}

func TestGetRequiredPanicError(t *testing.T) {
	assert.Panicsf(t, func() {
		Required("ANY_NOT_EXIST")
	}, `Environment "ANY_NOT_EXIST" is required`)
}

func TestGetAsIntPanicError(t *testing.T) {
	assert.Panicsf(t, func() {
		GetAsInt("ANY_NOT_EXIST")
	}, `Environment "ANY_NOT_EXIST" is not a integer`)
}

func TestGetAsBoolPanicError(t *testing.T) {
	assert.Panicsf(t, func() {
		GetAsBool("ANY_NOT_EXIST")
	}, `Environment "ANY_NOT_EXIST" is not a boolean`)
}

func TestGetAsFloat64PanicError(t *testing.T) {
	assert.Panicsf(t, func() {
		GetAsFloat64("ANY_NOT_EXIST")
	}, `Environment "ANY_NOT_EXIST" is not a float64`)
}
