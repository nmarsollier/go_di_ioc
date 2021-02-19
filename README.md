# Capítulo 1: Programación Funcional e Inversión de Control en GO

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

## El patrón Strategy

El patrón estrategia se fortaleció en la programación OO.
Permite establecer diferentes estrategias de resolución de un problema a través de una interfaz y múltiples implementaciones.

[Una búsqueda en wikipedia nos da mas detalles](https://es.wikipedia.org/wiki/Strategy_(patr%C3%B3n_de_dise%C3%B1o)).

Lo cierto es que la existencia de Strategy, es lo que le da sentido a la inyección de dependencias.

No debemos usar DI cuando no tenemos strategy. 

Porque digo esto ? Porque es muy común observar las siguientes conductas a la hora de programar : 

- Implementar interfaces si o si, para separar capas
- Implementar interfaces cuando solo existe una sola implementación
- Utilizar interfaces para poder mockear tests, cuando en realidad existe una sola implementación

### Lo que realmente deberíamos considerar es que :

- No debemos usar strategy cuando no tenemos varias implementaciones. (Esto quiere decir, no hacemos interfaces con una sola implementación)
- Una clase mock para testear no es excusa para implementar strategy.
- Solo debemos hacer DI cuando realmente tenemos una estrategia.
- Cuando *por las dudas* generalizamos y hacemos DI, estamos escribiendo código extra innecesario.
- Cuando queremos mockear para unit test, es preferible sobrescribir.

### Cuales son los problemas de la DI cuando se usa mal:

Aclarando que la inyección de dependencias es una buena practica, y recomendable, los vicios de implementarla en cualquier lado cuando no es necesario serian:

- Sobrecargamos los factories y/o métodos con instancias innecesariamente
- Generamos confusión al dejar abierta la puerta a múltiples implementaciones, cuando en realidad no las hay.
- Acoplamos código. Por ejemplo un controller no debería saber que instancia de un DAO utilizar un Servicio de negocio.
- Hacemos el código difícil de leer y por consiguiente de mantener

### Cuando SI deberíamos usar DI

- Cuando tenemos una estrategia, o sea varias implementaciones para resolver un problema.
- Cuando estamos programando un modulo y la implementación del comportamiento se define fuera del modulo
- Cuando queremos programar callbacks que dependen de quien lo llame

## Uso de Factory Methods como IoC

Si partimos los patrones generales de asignación de responsabilidades [GRASP](https://es.wikipedia.org/wiki/GRASP), una de las formas clásicas y adecuadas de uso de IoC es el uso de Factory Methods.

Esta estrategia nos permite evitar inyectar las dependencias en los constructores y delegar la instanciación a funciones factory.

Este ejemplo lo encontramos en [ioc_factory](./ioc_factory/)

Como vemos en la función main: la creación del service no esta acoplada a la creación del dao.

```
	srv := service.NewService()

	fmt.Println(srv.SayHello())
```

Sino mas bien el mismo service se encarga de crear el dao que corresponda según el contexto

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

Esta estrategia es muy util cuando solo tenemos una implementación, sin embargo queremos tener control de la creación de nuevas instancias.
En este caso el factory del DAO esta implementado en el Service, lo que nos permite tomar decisiones sobre la estrategia de construcción para el service en cuestión.

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


## Un enfoque mas funcional

En entornos de microservicios, la mayoría de nuestro código es privado al proyecto, o mas bien propio del microservicio, nuestro código es generalmente interno y privado, por los que la implementación interna del código, no tiene porque ser compleja ni escalable ni adaptable. 

Esa es la gracia de los microservicios, son soluciones adaptadas a un problema puntual, donde las interfaces del microservicio (REST, GRPC, etc) es todo lo que se usa desde afuera.

### Programamos Go como si fuera Java...

Cuando en Go definimos una estructura, asociamos código a una estructura y generalizamos el uso de esa estructura con una interfaz (como lo hice en el ejemplo anterior), básicamente estamos programando en forma orientada a objetos.

Se suele decir que Go no es orientado a objetos, pero sin embargo en la definición misma del lenguaje estos artefactos de Go, explícitamente hacen referencia a la programación OO. 

Podemos citar muchas referencias a ésto desde el mismo sitio de [Go](https://golang.org/doc/effective_go#interfaces_and_types), y la mayoría de los libros que he leído expresan conceptos de la misma forma.

Ahora bien, si en lugar de tomar un enfoque OO, aprovechamos las capacidades de Go para programar en forma funcional, podemos encontrarnos con un código mas prolijo y directo.

### Go en forma Funcional ? se puede.

Siguiendo los lineamientos de single Responsability debería ser bastante común en nuestro código tener servicios con una sola función.

Analicemos que significa en el código anterior, la siguiente estructura :

```
// HelloService es el servicio de negocio
type HelloService struct {
	dao IHelloDao
}
```

Es básicamente encapsular un puntero a una estructura que define una función. Ademas este tipo de estructuras es un antipatrón de programación OO, muy popularizado en Java con EJB, cuando no se sabia como separar capas.

No estoy en contra de las estructuras, pero este tipo de estructuras es solo un pasamanos a una llamada a una o varias funciones, no tiene razón de existir porque no tiene estado real en nuestro sistema. 

Que tal si nos evitamos todo esto y vamos directo a una solución donde no existan estructuras innecesarias ?

Nuestro archivo main, entonces no necesita crear ningún instancia, solo llamamos a una función :

```
func main() {
	fmt.Println(service.SayHello())
}
```

Nuestro DAO es muy simple también, es solo una función.

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

Dado que nuestro servicio es algo que debemos poder testear usando mocks de DAO, no nos quedan muchas opciones que permitirnos mockear esta función con un puntero a función. 
En caso que el if nos cause ruido, podemos apuntar directamente el puntero sayHelloMock a la función original, y nos evitamos ese if de mas.

El test en cuestión es el siguiente :

```
	// Cuando testeamos la reescribimos con el
	// mock que queramos
	sayHelloMock = func() string {
		return "Hello"
	}

	assert.Equal(t, "Hello", SayHello())
	sayHelloMock = nil
```

La estrategia de utilizar un puntero a una función, conceptualmente es la misma que utilizar una interfaz, en su forma mas simple un puntero a una función define una interfaz a respetar.

Es medio hacky ? si pero como dije antes, si algo puede ser hacky es el test.

### Ventajas de este enfoque 

Este concepto de programar en forma funcional simplifica mucho la programación tradicional, y en la mayoría de las soluciones es el balance ideal porque es la forma mas simple de definir, entender y mantener código. 

- No tenemos que escribir una interfaz, ni una estructura, tan solo para mockear
- No hacemos DI
- Single Responsibility
- Código mas simple de leer y mantener
- Podemos visualizar mejor la programación declarativa del paradigma funcional
- Llevamos el concepto de [Interface segregation](https://en.wikipedia.org/wiki/Interface_segregation_principle) a su mínima expresión, una función, algo deseable en POO

Incluso podemos hacer strategy si planificamos bien los factories, sin necesidad del uso de interfaces.

## Opinión personal sobre POO

La programación OO esta muy bien, solo que esta muy subestimada la forma en la que se debe realizar. Programar OO es muy complejo, y cuando los proyectos crecen, no se realiza el mantenimiento adecuado, por lo que en general terminamos teniendo código espagueti.

El libro Domain Driven Design de Erik Evans, expresa una forma conceptualmente adecuada de implementar POO, sin embargo es muy raro ver algo claro y bien implementado.

Muchos desarrolladores entienden que el concepto de Clean Architecture y DDD solo va en generar interfaces hacia todo lo que entra y sale del negocio, pero olvidan lo fundamental, respetar los patrones básicos para que el código sea sustentable, que es donde se encuentra el balance de simplicidad necesario.

La programación OO no es aplicar la misma regla y norma a todo, sino, requiere el uso del energía cerebral para que nuestros diseños tengan la complejidad justa y necesaria. Y eso es muy difícil de adquirir. 

Y a su vez, la POO requiere un mantenimiento adecuado, el refactor continuo debe ser una norma, y no siempre los equipos de desarrolladores lo entienden asi.

Por otro lado un enfoque funcional es mucho mas simple, el refactor es simple, lo que debemos adoptar es simplemente una separación clara del negocio con las dependencias que usa y que usan al mismo. Teniendo esa separación en capas bien lograda, el resultado es elegante, simple y sobre todo, muy eficiente.

En las empresas en general mas de la mitad de los desarrolladores tendrán poca experiencia, muchos de ellos estarán dando sus primeros pasos, por lo que ésta simplicidad es bienvenida.

## Nota

Esta es una serie de tutoriales sobre patrones simples de programación en GO.

[Capítulo 2: REST Controllers en go](https://github.com/nmarsollier/go_rest_controller)
