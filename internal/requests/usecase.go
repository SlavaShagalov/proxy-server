package requests

type Usecase interface {
	Create(*Request) error
	Get(id string) (*Request, error)
	List() ([]Request, error)
}
