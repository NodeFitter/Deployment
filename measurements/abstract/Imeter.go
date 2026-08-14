package abstract

type Imeter interface {
	Init() error
	Increment() error
	Get() (uint, error)
}
