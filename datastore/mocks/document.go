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
