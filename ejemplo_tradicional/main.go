package main

import (
	"fmt"

	"github.com/nmarsollier/go_di_ioc/ejemplo_tradicional/dao"
	"github.com/nmarsollier/go_di_ioc/ejemplo_tradicional/service"
)

func main() {

	srv := service.NewService(dao.NewDao())

	fmt.Println(srv.SayHello())
}
