package main

import (
	"fmt"

	"github.com/nmarsollier/go_di_ioc/ioc_factory/service"
)

func main() {

	srv := service.NewService()

	fmt.Println(srv.SayHello())
}
