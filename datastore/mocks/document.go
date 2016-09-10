package mocks

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
)

var (
	BadId  = "123"
	GoodId = "4b06dcbe-aba8-4314-a7c4-03cdde197a38"
)


