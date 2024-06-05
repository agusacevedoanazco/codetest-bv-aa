# Mock-endpoint

Microservicio creado sobre Go encargado de responder para toda request tipo "GET" en la ruta "/" con la siguiente respuesta de tipo `application/json`

```json
{
  "insurance": {
    "name": "Nombre de Seguro",
    "description": "Descripcion del Seguro"
    "price": "Valor"
    "image": "url_de_imagen"
  }
}
```

> Este microservicio fue desarrollado únicamente para validación local, no deberia ser tomado en cuenta para la evaluación

### Estructura de la aplicación

```
mock-endpoint
├── .dockerignore
├── .gitignore
├── Dockerfile
├── README.md
├── cmd
│   └── web
│       ├── handlers.go
│       ├── helpers.go
│       └── main.go
└── go.mod
```

La aplicación contiene un único directorio principal `web/cmd` el cual contiene el servicio web y handler de la ruta "/" que retorna la respuesta definida al inicio del documento

### Inicio rápido de la aplicación

Teniendo el binario de go instalado, basta con ejecutar lo siguiente de modo de apuntar al paquete principal

```bash
go ./cmd/web
```

La aplicación levantará un servicio web en el puerto 8080 escuchando en **TODAS** las interfaces del ambiente de ejecución

### Buildeo de la aplicación

En el archivo `Dockerfile` se encuentra el manifesto de definicion del proceso de buildeo de la aplicacion para su despliegue, el proceso a grandes rasgos es el siguiente

- Define el entorno de buildeo sobre `golang:alpine3.20`
- Copia los archivos fuente al ambiente temporal
- Construye el target de buildeo sobre el paquete main contenido en `./cmd/web`
- Luego minimiza la imagen definiendo un stage del servicio
- Copia el binario ejecutable construido en los pasos anteriores
- Define el usuario y grupo a ejecutar el microservicio
- Expone el puerto definido en el microservicio
- Define el comando de inicializacion del servicio concluyendo el buildeo de la imagen

#### Buildeo local

Alternativamente para buildear la imagen de forma local, basta con tener instalados los binarios de desarrollo de golang en su version 1.22.2 y ejecutar

```
go build ./cmd/web
```

La ejecución tiene como salida el archivo binario ejecutable `mock-bv-58`, que hace referencia a los archivos estatico de forma relativa al directorio ./ui/
