package counter

//Provider represents customized metric provider
type Provider interface {
	//Map maps value int slice index
	Map(interface{}) int
	//Keys  returns mapped keys
	Keys() []string
}

type CustomCounter interface {
	Aggregate(value interface{})
}

type CustomProvider interface {
	NewCounter() CustomCounter
}
