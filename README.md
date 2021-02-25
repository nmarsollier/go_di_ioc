# DI y IoC en GO

Este repositorio plantea alternativas de manejo de dependencias, a la programación tradicional de un proyecto Go. 
## Inyección de Dependencias

_Spoiler alert: Es lo que debemos cambiar_

Es esa estrategia de IoC que nos permite insertar dependencias en una clase para que sean usadas internamente

En la carpeta [ejemplo_tradicional](./ejemplo_tradicional/) tenemos los ejemplos de código.

La mayoría de los autores recomiendan inyección de dependencias para separa capas lógicas de código. 

Nuestro código luce como el siguiente: 

```
srv := service.NewService(dao.NewDao())
fmt.Println(srv.SayHello())
```

donde Service es algo como lo siguiente :

```
// IHelloDao interface DAO necesaria a inyectar en el service
type IHelloDao interface {
	Hello() string
}

// HelloService es el servicio de negocio
type HelloService struct {
	dao IHelloDao
}

// NewService Es un factory del servicio de Negocio , depende de IHelloDao
func NewService(dao IHelloDao) *HelloService {
	return &HelloService{dao}
}

// SayHello es nuestro método de negocio
func (s HelloService) SayHello() string {
	return s.dao.Hello()
}
```

Se implementa generalmente pasando una instancia de las dependencias necesarias en constructor o en el método que la necesite.

Según la bibliografía este tipo de patrón nos :

- Permite desacoplar nuestro código, de forma que pueda ser configurable
- Reduce la complejidad de código
- Permite hacer código reusable, y testable
- Permite mockear tests
 
_Y esto es cierto, pero hasta cierto punto_

Porque no desacoplamos realmente, todo lo contrario, terminamos acoplando mucho mas, nuestro código debe definir métodos bootstraps en lugares donde no deberían estar, acoplando todo el negocio en un archivo main.go por ejemplo. 

## Uso de Factory Methods como IoC

Veamos como podemos mejorar la situación anterior.

Si partimos los patrones generales de asignación de responsabilidades [GRASP](https://es.wikipedia.org/wiki/GRASP), una de las formas clásicas y adecuadas de uso de IoC es el uso de Factory Methods.

Esta estrategia nos permite evitar inyectar las dependencias en los constructores y delegar la instanciación a funciones factory.

Este ejemplo lo encontramos en [ioc_factory](./ioc_factory/)

Como vemos en la función main: la creación del service no esta acoplada a la creación del dao.

```
	srv := service.NewService()

	fmt.Println(srv.SayHello())
```

Sino mas bien el mismo service se encarga de crear el dao que corresponda según el contexto. 

Esto esta muy en linea con el patrón experto.

```
// IHelloDao interface DAO necesaria a inyectar en el service
type IHelloDao interface {
	Hello() string
}

// HelloService es el servicio de negocio
type HelloService struct {
	dao IHelloDao
}

// NewService es una función que puede mockearse
func NewService() *HelloService {
	return &HelloService{
		buildDao(),
	}
}

// SayHello es nuestro método de negocio
func (s *HelloService) SayHello() string {
	return s.dao.Hello()
}

// mockedDao es un DAO que nos permite mockear los tests
var mockedDao IHelloDao = nil

func buildDao() IHelloDao {
	if mockedDao != nil {
		return mockedDao
	}

	return dao.NewDao()
}
```

Para realizar mocks en los tests solo tenemos que definir un valor para mockedDao

```
mockedDao = new(daoMock)

s := NewService()
assert.Equal(t, "Hello", s.SayHello())

// Volvemos al original
mockedDao = nil
```

Al ser privado mockedDao no afecta al uso de la librería. Suena hacky, pero en realidad si hay algo que puede ser hacky es el test, es preferible que el test quede hacky a que toda la aplicación quede sobre estructurada solo para poder testear.

A su vez, siguiendo los lineamientos de no realizar estrategias donde no es necesario, el dao, no expone interfaces, es solo una estructura.

---
Este mock es muy rudimentario, en próximas notas voy a mostrar como mejorar el tema de funciones mock.

---


```
// HelloDao Es nuestra implementación de Dao
type HelloDao struct {
}

// NewDao es el factory
func NewDao() *HelloDao {
	return new(HelloDao)
}

// Hello es nuestro método de negocio
func (d *HelloDao) Hello() string {
	return "Holiiis"
}
```

Ventajas:
- Permite encapsular el código de forma correcta, definiendo lo que se necesita en el lugar que se necesita 
- Permite reducir complejidad de constructores
- Podemos utilizar el patrón experto de forma más clara y concisa

## Ahora veamos los fundamentos 

### El patrón Strategy

El patrón estrategia se fortaleció en la programación OO.
Permite establecer diferentes estrategias de resolución de un problema a través de una interfaz y múltiples implementaciones.

[Una búsqueda en wikipedia nos da mas detalles](https://es.wikipedia.org/wiki/Strategy_(patr%C3%B3n_de_dise%C3%B1o)).

Lo cierto es que la existencia de Strategy, es lo que le da sentido a la inyección de dependencias.

No debemos usar DI cuando no tenemos strategy. 

Porque digo esto ? Porque es muy común observar las siguientes conductas a la hora de programar : 

- Implementar interfaces si o si, para separar capas
- Implementar interfaces cuando solo existe una sola implementación
- Utilizar interfaces para poder mockear tests, cuando en realidad existe una sola implementación
- O simplemente porque es la forma que todos dicen

### Lo que realmente deberíamos considerar es que :

- No debemos usar strategy cuando no tenemos varias implementaciones. (Esto quiere decir, no hacemos interfaces si no hay polimorfismo)
- Una clase mock para testear no es excusa para implementar strategy.
- Solo debemos hacer DI cuando realmente tenemos una estrategia.
- Cuando *por las dudas* generalizamos y hacemos DI, estamos escribiendo código extra innecesario.
- Cuando queremos mockear para unit test, es preferible sobrescribir.

### Cuales son los problemas de la DI cuando se usa mal:

Aclarando que la inyección de dependencias es una buena practica, y recomendable, los vicios de implementarla en cualquier lado cuando no es necesario serian:

- Sobrecargamos los factories y/o métodos con instancias innecesariamente.
- Generamos confusión al dejar abierta las puertas al polimorfismo , cuando en realidad no lo hay.
- Acoplamos código. Por ejemplo un controller no debería saber que instancia de un DAO utilizar un Servicio de negocio.
- Hacemos el código difícil de leer y por consiguiente de mantener.

### Cuando SI deberíamos usar DI

- Cuando tenemos una estrategia, o sea polimorfismo para resolver un problema.
- Cuando estamos programando un modulo y la implementación del comportamiento se define fuera del modulo.
- Cuando queremos programar callbacks que dependen de quien lo llame.
- Cuando programamos una librería y queremos ser user friendly para terceros.

## Nota

Esta es una serie de tutoriales sobre patrones simples de programación en GO.

[Tabla de Contenidos](https://github.com/nmarsollier/go_index)
