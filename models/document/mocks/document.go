package mocks

import (
	"github.com/bapjiws/local_times_dashboard_backend/models/document"
)

// Mock documents
type MockDocument struct {
	SelfDescription string
}

var (
	BadDocument = MockDocument{
		SelfDescription: "I'm bad! I'm, like, really bad!",
	}

	GoodDocument = MockDocument{
		SelfDescription: "S'all good, man.",
	}

	GoodDocuments = []document.Document{
		MockDocument{
			SelfDescription: "S'all good, man.",
		},
		MockDocument{
			SelfDescription: "I will send you to Belize!",
		},
	}

	EmptyDocumentList = []document.Document{}
)

// Mock IDs
var (
	BadId  = "123"
	GoodId = "4b06dcbe-aba8-4314-a7c4-03cdde197a38"
)

// Mock suggest params
var (
	SuggesterName = "mockDoc_suggest"
	Field         = "suggest"
	PayloadKey    = "mockDoc_id"

	BadText                     = "\033c"
	GoodTextWithNoSuggestions   = "no need in suggestions"
	GoodTextWithSomeSuggestions = "gimme some"
)
