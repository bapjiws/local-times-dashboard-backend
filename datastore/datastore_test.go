package datastore

import (
	"errors"
	"testing"
	"timezones_mc/datastore/mocks"

	"github.com/stretchr/testify/assert"
)

func TestAddDocument(t *testing.T) {
	assert := assert.New(t)

	testDataStore := new(mocks.Datastore)

	testCases := []struct {
		document mocks.MockDocument
		error    error
	}{
		{
			document: mocks.BadDocument,
			error:    errors.New("Bad document!"),
		},
		{
			document: mocks.GoodDocument,
			error:    nil,
		},
	}

	for _, testCase := range testCases {
		testDataStore.On("AddDocument", testCase.document).Return(testCase.error)
		err := testDataStore.AddDocument(testCase.document)
		assert.Equal(testCase.error, err)
	}

	// Assert that everything specified with On and Return was in fact called as expected.
	testDataStore.AssertExpectations(t)
}
