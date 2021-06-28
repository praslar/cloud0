package ginext

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type userValidation struct {
	FirstName string `validate:"required"`
	LastName  string `validate:"required" json:"last_name"`
	Age       uint8  `validate:"gte=0,lte=150" json:"-"`
	Email     string `validate:"email"`
}

func TestNewValidator(t *testing.T) {
	v := NewValidator()
	_, ok := v.(*validatorImpl)
	assert.True(t, ok)
}

func TestValidateNoError(t *testing.T) {
	testData := userValidation{
		FirstName: "Eric",
		LastName:  "Huynh",
		Age:       1,
		Email:     "eric.huynh@gmail.com",
	}
	err := NewValidator().ValidateStruct(testData)
	assert.NoError(t, err)
}

func TestValidateWithError(t *testing.T) {
	testData := userValidation{
		FirstName: "",
		LastName:  "Huynh",
		Age:       1,
		Email:     "eric.huynh@gmail.com",
	}
	err := NewValidator().ValidateStruct(testData)
	assert.NotNil(t, err)

	assert.Contains(t, err.Error(), "failed validation on 1 field(s)")

	vErrors, ok := err.(ValidatorErrors)
	assert.True(t, ok)
	assert.NotNil(t, vErrors)
	assert.Equal(t, 1, len(vErrors.GetErrors()))
	assert.Contains(t, vErrors.GetErrors()[0].Error(), "FirstName: failed validation on tag required")

	// test get map
	m := vErrors.GetErrorsMap()
	assert.Len(t, m, 1)
	assert.Equal(t, "failed validation on tag required", m["FirstName"])
}

func TestValidateFailedWithJSONField(t *testing.T) {
	testData := userValidation{
		FirstName: "Eric",
		LastName:  "",
		Age:       1,
		Email:     "eric.huynh@gmail.com",
	}
	err := NewValidator().ValidateStruct(testData)
	vErrors, _ := err.(ValidatorErrors)
	assert.Contains(t, vErrors.GetErrors()[0].Error(), "last_name: failed validation on tag required")
}

func TestValidateFailedWithParamDetailIfSet(t *testing.T) {
	testData := userValidation{
		FirstName: "Eric",
		LastName:  "Huynh",
		Age:       200,
		Email:     "eric.huynh@gmail.com",
	}
	err := NewValidator().ValidateStruct(testData)
	vErrors, _ := err.(ValidatorErrors)
	assert.Contains(t, vErrors.GetErrors()[0].Error(), "Age: failed validation on tag lte (param: 150, value: 200)")
}

func TestReturnErrorOnInvalidStruct(t *testing.T) {
	err := NewValidator().ValidateStruct(nil)
	assert.Error(t, err)
	_, ok := err.(*validator.InvalidValidationError)
	assert.True(t, ok)
}
