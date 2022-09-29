package serviceb

import "fmt"

type ServiceB struct {
}

func NewServiceB() *ServiceB {
	return &ServiceB{}
}

func (s *ServiceB) DoSomething() {
	fmt.Println("Hello World from service-b")
}
