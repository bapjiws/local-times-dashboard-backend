package datastore

import (
	"errors"
	"testing"

	datastore_mocks "github.com/bapjiws/local_times_dashboard_backend/models/datastore/mocks"
	"github.com/bapjiws/local_times_dashboard_backend/models/document"
	document_mocks "github.com/bapjiws/local_times_dashboard_backend/models/document/mocks"
	"github.com/stretchr/testify/assert"
)

// AddDocument(doc models.Document) error
func TestAddDocument(t *testing.T) {
	assert := assert.New(t)

	testDataStore := new(datastore_mocks.Datastore)

	testCases := []struct {
		document document.Document
		error    error
	}{
		{
			document: document_mocks.BadDocument,
			error:    errors.New("Bad document!"),
		},
		{
			document: document_mocks.GoodDocument,
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

// FindDocumentById(id string) (models.Document, error)
func TestFindDocumentById(t *testing.T) {
	assert := assert.New(t)

	testDataStore := new(datastore_mocks.Datastore)

	testCases := []struct {
		id       string
		document document.Document
		error    error
	}{
		{
			id:       document_mocks.BadId,
			document: nil,
			error:    errors.New("Bad id!"),
		},
		{
			id:       document_mocks.GoodId,
			document: document_mocks.GoodDocument,
			error:    nil,
		},
	}

	for _, testCase := range testCases {
		testDataStore.On("FindDocumentById", testCase.id).Return(testCase.document, testCase.error)
		document, err := testDataStore.FindDocumentById(testCase.id)
		assert.Equal(testCase.document, document)
		assert.Equal(testCase.error, err)
	}

	// Assert that everything specified with On and Return was in fact called as expected.
	testDataStore.AssertExpectations(t)
}

// SuggestDocuments(suggesterName string, text string, field string, payloadKey string) ([]models.Document, error)
func TestSuggestDocuments(t *testing.T) {
	assert := assert.New(t)

	testDataStore := new(datastore_mocks.Datastore)

	suggesterName := document_mocks.SuggesterName
	field := document_mocks.Field
	payloadKey := document_mocks.PayloadKey

	testCases := []struct {
		text        string
		suggestions []document.Document
		error       error
	}{
		{
			text:        document_mocks.BadText,
			suggestions: nil,
			error:       errors.New("Bad suggestion text!"),
		},
		{
			text:        document_mocks.GoodTextWithNoSuggestions,
			suggestions: document_mocks.EmptyDocumentList,
			error:       nil,
		},
		{
			text:        document_mocks.GoodTextWithSomeSuggestions,
			suggestions: document_mocks.GoodDocuments,
			error:       nil,
		},
	}

	for _, testCase := range testCases {
		testDataStore.
			On("SuggestDocuments", suggesterName, testCase.text, field, payloadKey).
			Return(testCase.suggestions, testCase.error)
		suggestions, err := testDataStore.SuggestDocuments(suggesterName, testCase.text, field, payloadKey)
		assert.Equal(testCase.suggestions, suggestions)
		assert.Equal(testCase.error, err)
	}

	// Assert that everything specified with On and Return was in fact called as expected.
	testDataStore.AssertExpectations(t)
}
