# Programación Funciona e Inversión de Control en GO

Este repositorio plantea alternativas de manejo de dependencias, a la programacion tradicional de un proyecto Go. 

## Inyección de Dependencias

_Spoiler alert: Es lo que debemos cambiar_

Es esa estrategia de IoC que nos permite insertar dependencias en una clase para que sean usadas internamente

En la carpeta [ejemplo_tradicional](./ejemplo_tradicional/) tenemos los ejemplos de codigo.

La mayoria de los autores recomiendan inyeccion de dependencias para separa capas logicas de codigo. 

Nuestro codigo luce como el siguiente: 

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

// SayHello es nuestro metodo de negocio
func (s HelloService) SayHello() string {
	return s.dao.Hello()
}
```

Se implementa generalmente pasando una instancia de las dependencias necesarias en constructor o en el método que la necesite.

Segun la bibliografia este tipo de patron nos :

- Permite desacoplar nuestro codigo, de forma que pueda ser configurable
- Reduce la complejidad de código
- Permite hacer codigo reusable, y testable
- Permite mockear tests
 
_Y esto es cierto, pero hasta cierto punto_

## El patrón Strategy

El patrón estrategia se fortaleció en la programación OO.
Permite establecer diferentes estrategias de resolución de un problema a través de una interfaz y múltiples implementaciones.

[Una busqueda en wikipedia nos da mas detalles](https://es.wikipedia.org/wiki/Strategy_(patr%C3%B3n_de_dise%C3%B1o)).

Lo cierto es que la existencia de Strategy, es lo que le da sentido a la inyeccion de dependencias.

No debemos usar DI cuando no tenemos strategy. 

Porque digo esto ? Porque es muy comun observar las siguientes conductas a la hora de programar : 

- Implementar interfaces si o si, para separar capas
- Implementar interfaces cuando solo existe una sola implementacion
- Utilizar interfaces para poder mockear tests, cuando en realidad existe una sola implementacion

### Lo que realmente deberiamos considerar es que :

- No debemos usar strategy cuando no tenemos varias implementaciones. (Esto quiere decir, no hacemos interfaces con una sola implementación)
- Una clase mock para testear no es excusa para implementar strategy.
- Solo debemos hacer DI cuando realmente tenemos una estrategia.
- Cuando *por las dudas* generalizamos y hacemos DI, estamos escribiendo código extra innecesario.
- Cuando queremos mockear para unit test, es preferible sobreescribir.

### Cuales son los problemas de la DI cuando se usa mal:

- Sobrecargamos los factories y/o metodos con instancias innecesariamente
- Generamos confusion al dejar abierta la puerta a multiples implementaciones, cuando en realidad no las hay.
- Acoplamos codigo. Por ejemplo un controller no deberia saber que instancia de un DAO utilizar un Servicio de negocio.
- Hacemos el codigo dificil de leer y por consiguiente de mantener

### Cuando SI deberiamos usar DI

- Cuando tenemos una estrategia, o sea varias implementaciones para resolver un problema.
- Cuando estamos programando un modulo y la implementacion del comportamiento se define fuera del modulo
- Cuando queremos programar callbacks que dependen de quien lo llame

## Uso de Factory Methods como IoC

Si partimos los patrones generales de asignacion de responsabilidades [GRASP](https://es.wikipedia.org/wiki/GRASP), una de las formas clasicas y adecuadas de uso de Inversion de control es el uso de Factory Methods.

Esta estrategia nos permite evitar inyectar las dependencias en los constructores y delegar la instanciacion a metodos factory.

Este ejemplo lo encontramos en [ioc_factory](./ioc_factory/)

Como vemos en el metodo main, la creacion del service no esta acoplada a la creacion del dao.

```
	srv := service.NewService()

	fmt.Println(srv.SayHello())
```

Sino mas bien el mismo service se encarga de crear el dao que corresponda segun el contexto

```
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

// mockedDao es un DAO que nos permite mockear los tests
var mockedDao IHelloDao = nil

func buildDao() IHelloDao {
	if mockedDao != nil {
		return mockedDao
	}

	return dao.NewDao()
}
```

Esta estrategia es muy util cuando solo tenemos una implementacion, sin embargo queremos tener control de la creacion de nuevas instancias.
En este caso el factory del DAO esta implementado en el Service, lo que nos permite tomar decisiones sobre la estrategia de construccion para el service en cuestion.

Para realizar mocks en los tests solo tenemos que definir un valor para mockedDao

```
	mockedDao = new(daoMock)

	s := NewService()
	assert.Equal(t, "Hello", s.SayHello())

	// Volvemos al original
	mockedDao = nil
```

Al ser privado mockedDao no deberia afectar al uso de la libreria. Suena hacky , pero en realidad si hay algo que puede ser hacky es el test, es preferible que el test quede hacky a que toda la aplicacion quede sobre estructurada solo para poder testear.

A su vez, siguiendo los lineamientos de no realizar estrategias donde no es necesario, el dao, no expone interfaces, es solo una estructura.

```
// HelloDao Es nuestra implementacion de Dao
type HelloDao struct {
}

// NewDao es el factory
func NewDao() *HelloDao {
	return new(HelloDao)
}

// Hello es nuestro metodo de negocio
func (d *HelloDao) Hello() string {
	return "Holiiis"
}
```

Ventajas:
- Permite encapsular el código de forma correcta, definiendo lo que se necesita en el lugar que se necesita 
- Permite reducir complejidad de constructores
- Podemos utilizar el patrón experto de forma más clara y concisa


## Un enfoque mas funcional

En entornos de microservicios, la mayoria de nuestro codigo es privado al proyecto, o mas bien propio del microservicio, nuestro codigo es generalmente interno y privado, por los que la implementacion interna del codigo, no tiene porque ser compleja ni escalable ni adaptable. 

Esa es la gracia de los microservicios, son soluciones adaptadas a un problema puntual, donde las interfaces del microservicio (REST, GRPC, etc) es todo lo que se usa desde afuera.

### Programamos Go como si fuera Java...

Cuando en Go definimos una estructura, asociamos codigo a una estructura y generlizamos el uso de esa estructura con una interfaz (como lo hice en el ejemplo anterior), basicamente estamos programando en forma orientada a objetos.

Podran decir que Go no es orientado a objeto, pero sin embargo en la definicion misma del lenguaje estos artifactos de Go, explicitamente hacen referencia a la programacion OO. 

Podemos citar muchas referencias a ésto desde el mismo sitio de [Go](https://golang.org/doc/effective_go#interfaces_and_types), y la mayoria de los libros que he leido expresan conceptos de la misma forma.

Ahora bien, si en lugar de tomar un enfoque OO, aprovechamos las capacidades de Go para programar en forma funcional, podemos encontrarnos con un codigo mas prolijo y directo.

### Go en forma Funciona ? se puede.

Siguiendo los lineamientos de single Responsability debería ser bastante común en nuestro código tener servicios con una sola función.

Analicemos que significa en el codigo anterior, la siguiente estructura :

```
// HelloService es el servicio de negocio
type HelloService struct {
	dao IHelloDao
}
```

Es basicamnete encapsular un puntero a una estructura que define una funcion. Ademas este tipo de estructuras es un antipatron de programacion OO.

No estoy en contra de las estrcuturas, pero este tipo de estrcuturas es solo un pasamanos a una llamada a una o varias funciones, no tiene razon de existir porque no tiene estado. 

Que tal si nos evitamos todo esto y vamos directo a una solucion donde no existan estructuras innecesarias ?

Nuestro archivo main, entonces no necesita crear ningun instancia, solo llamamos a una funcion :

```
func main() {
	fmt.Println(service.SayHello())
}
```

Nuestro DAO es muy simple tambien, es solo una funcion.

```
func Hello() string {
	return "Holiiiiis"
}
```

Nuestro Service es un poquito mas complejo, pero no llega a ser tan complejo como usar interfaces :

```
// Nos va a permitir mockear respuestas para los tests
var sayHelloMock func() string = nil

// SayHello es nuestro negocio
func SayHello() string {
	if sayHelloMock != nil {
		return sayHelloMock()
	}

	return dao.Hello()
}
```

Dado que nuestro servicio es algo que debemos poder testear usando mocks de daos, no nos quedan muchas opciones que permitirnos mockear esta funcion con un puntero a funcion.

El test en cuestion es el siguiente :

```
	// Cuando testeamos la reescribimos con el
	// mock que queramos
	sayHelloMock = func() string {
		return "Hello"
	}

	assert.Equal(t, "Hello", SayHello())
	sayHelloMock = nil
```

La estrategia de utilizar un puntero a una funcion, conceptualmente es la misma que utilizar una interfaz, en su forma mas simple un puntero a una funcion define una interfaz a respetar.

Es medio hacky ? si pero como dije antes, si algo puede ser hacky es el test.

### Ventajas de este enfoque 

Este concepto de programar en forma funcional simplifica mucho la programacion tradicional, y en la mayoria de las soluciones es el balance ideal porque es la forma mas simple de definir, entender y mantener codigo. 

- No tenemos que escribir una interfaz, ni una estructura, tan solo para mockear
- No hacemos DI
- Single Responsibility
- Código mas simple de leer y mantener
- Podemos visualizar mejor la programacion declarativa del paradigma funcional
- Llevamos el concepto de [Interface segregation](https://en.wikipedia.org/wiki/Interface_segregation_principle) a su minima expresion, una funcion, algo deseable en POO

Incluso podemos hacer strategy si planificamos bien los factories, sin necesidad del uso de interfaces.

## Opinion personal sobre POO

La programacion OO esta muy bien, solo que esta muy devaluada la forma en la que se debe realizar. Programar OO es muy complejo, y cuando los proyectos crecen, no se realiza el mantenimiento adecuado, por lo que en general terminamos teniendo codigo espagueti.

El libro Domain Driven Design de Erik Evans, expresa una forma conceptualmente adecuada de implementar POO, sin embargo es muy raro ver algo claro y bien implementado.

Muchos desarrolladores entienden que el concepto de Clean Architecture y DDD solo va en generar interfaces hacia todo lo que entra y sale del negocio, pero olvidan lo fundamental, volver a las raices, respetar los patrones basicos, que es donde se encuentra el balance de simplicidad necesario.

La programacion OO no es aplicar la misma regla y norma a todo, sino, requiere el uso del energia cerebral para que nuestros diseños tengan la complejidad justa y necesaria. Y eso es muy dificil de adquirir. Y a su vez, la POO requiere un mantenimiento adecuado, el refactor continuo debe ser una norma, y no siempre los equipos de desarrolladores lo entienden asi.

Por otro lado la programacion funcional es simple, y refactor es simple, lo que debemos adoptar es simplemente una separacion clara del negocio con las dependencias que usa y que usan al mismo. Teniendo esa separacion en capas bien lograda, el resultado es elegante, simple y sobre todo, muy eficiente.
