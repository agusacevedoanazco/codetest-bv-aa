# Microservice

Microservicio creado sobre Go que consume un endpoint con una respuesta predefinida bajo la siguiente estructura:

```json
{
  "insurance": {
    "name": "string-value",
    "description": "string-value"
    "price": "string-value"
    "image": "string-value"
  }
}
```

El microservicio consume los datos y los presenta en un template html en la ruta "/".

### Estructura de la aplicación

```
microservice
├── .dockerignore
├── Dockerfile
├── README.md
├── cmd
│   └── web
│       ├── handlers.go
│       ├── helpers.go
│       ├── main.go
│       └── main_test.go
├── go.mod
├── internal
│   └── client
│       ├── client.go
│       └── client_test.go
└── ui
    ├── html
    │   └── homepage.tmpl
    └── static
        ├── css
        │   ├── fonts.css
        │   ├── index.html
        │   └── main.css
        ├── fonts
        │   ├── MesloLGSNerdFont-Bold.ttf
        │   ├── MesloLGSNerdFont-Italic.ttf
        │   ├── MesloLGSNerdFont-Regular.ttf
        │   ├── MesloLGSNerdFontBold-Italic.ttf
        │   └── index.html
        └── index.html
```

La aplicacion está dividida en 3 directorios principales:

- web/cmd: que contiene el paquete `main` con el servicio web y el handler para servir el contenido desde la ruta "/" y sus correspondientes test
- internal/client: que contiene el paquete `client` encargado de conectarse con el endpoint y transformar los datos en una estructura definida
- ui: que mantiene los achivos estaticos que son utilizados por los templates renderizados por el handler del microservicio

### Inicio rápido de la aplicación

Teniendo el binario de go instalado, basta con definir las variables de entorno del microservicio `ENDPOINT_URL` en el formato `{ http || https }://{ urn }` y ejecutar el paquete main

```bash
go ./cmd/web
```

La aplicación levantará un servicio web en el puerto 32080 escuchando en **TODAS** las interfaces del ambiente de ejecución

### Buildeo de la aplicación

En el archivo `Dockerfile` se encuentra el manifesto de definicion del proceso de buildeo de la aplicacion para su despliegue, el proceso a grandes rasgos es el siguiente

- Define el entorno de buildeo sobre `golang:alpine3.20`
- Copia los archivos fuente al ambiente temporal
- Construye el target de buildeo sobre el paquete main contenido en `./cmd/web`
- Luego minimiza la imagen definiendo un stage del servicio
- Copia los archivos estaticos y el ejecutable construido en los pasos anteriores
- Define el usuario y grupo a ejecutar el microservicio
- Expone el puerto definido en el microservicio
- Define el comando de inicializacion del servicio concluyendo el buildeo de la imagen

#### Buildeo local

Alternativamente para buildear la imagen de forma local, basta con tener instalados los binarios de desarrollo de golang en su version 1.22.2 y ejecutar

```
go build ./cmd/web
```

La ejecución tiene como salida el archivo binario ejecutable `microservice`, que hace referencia a los archivos estatico de forma relativa al directorio ./ui/
