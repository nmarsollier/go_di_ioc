package service

import "github.com/nmarsollier/go_di_ioc/go_funcional/dao"

// Nos va a permitir mockear respuestas para los tests
var sayHelloMock func() string = nil

// SayHello es nuestro negocio
func SayHello() string {
	if sayHelloMock != nil {
		return sayHelloMock()
	}

	return dao.Hello()
}
