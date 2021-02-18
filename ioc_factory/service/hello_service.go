package service

import "github.com/nmarsollier/go_di_ioc/ioc_factory/dao"

// IHelloDao interface DAO necesaria a inyectar en el service
type IHelloDao interface {
	Hello() string
}

// HelloService es el servicio de negocio
type HelloService struct {
	dao IHelloDao
}

// NewService es una funcion que puede mockearse
func NewService() *HelloService {
	return &HelloService{
		buildDao(),
	}
}

// SayHello es nuestro metodo de negocio
func (s *HelloService) SayHello() string {
	return s.dao.Hello()
}

// MockedDao es un DAO que nos permite mockear los tests
var MockedDao IHelloDao = nil

func buildDao() IHelloDao {
	if MockedDao != nil {
		return MockedDao
	}

	return dao.NewDao()
}
