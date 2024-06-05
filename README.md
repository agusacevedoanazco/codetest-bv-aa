# CodeTest

Implementación de un pipeline de despliegue de una aplicación sobre Google Cloud Platform

### Tabla de contenidos (TOC)

- //TODO

### Contenidos del repositorio

Bajo el repositorio se encuentran 4 directorios de los cuales 3 son parte de los requerimientos de implementación

- `.github/workflows`: Directorio que contiene el pipeline de flujo CI/CD con los distintos jobs para la prueba, escaneo, buildeo y despliegue en ambientes staging y prod del microservicio/aplicación.
- `infrastructure`: Directorio con el IaC utilizado para la creación y mantención de los estados de los recursos a desplegar en Google Cloud Platform.
- `microservice`: Directorio que contiene el microservicio (aplicación) a desplegar en los ambientes ejecutados por el pipeline.
- `mock-endpoint`: _(adicional)_ Directorio que contiene un microservicio que emula el comportamiento del api endpoint provisto.

```
.
├── .github
│   └── workflows
├── infrastructure
│   ├── prod
│   └── staging
├── microservice
│   ├── cmd
│   │   └── web
│   ├── internal
│   │   └── client
│   └── ui
│       ├── html
│       └── static
└── mock-endpoint
    └── cmd
        └── web
```

### Aplicación (microservicio)

Bajo el directorio `microservice` se encuentra una aplicación que consume un endpoint provisto por la variable de entorno `ENDPOINT_URL`, ésta aplicación se encarga de obtener los datos del endpoint y servir
una vista en el puerto `32080` bajo la ruta `/` con los contenidos del endpoint. El servicio está desarrollado unicamente para leer los datos del endpoint en el formato provisto y no realiza ningún tipo de
digestión de la ruta del endpoint.

La aplicación provee interfaces para la prueba de las funcionalidades principales, de modo de automatizar su proceso de integración y despliegue, estas se encuentran bajo la ruta de cada paquete con el nombre
`*_test.go`

Adicionalmente, se define un archivo de buildeo de imagen de contenedor para la utilización en su proceso de despliegue.

### Despliegue en Cloud Pública

Como se menciona en el punto anterior, se define una imagen de contenedor del microservicio para el despliegue en cualquier servicio/plataforma Cloud Pública que sea compatible con Docker containers, la
alternativa elegida en este caso es _Google Cloud Run_, no obstante es compatible con EKS, ECS, GKE, GCRun, AKS, ACI, entre otros.

### Pipeline CI/CD

Como se menciona en los contenidos del repositorio se definen distintos workflows files conforme a los contenidos del directorio `workflows`

```
.github/workflows
├── build_scan_push.yaml
├── deploy_prod.yaml
├── deploy_staging.yaml
├── pipeline.yaml
├── source_scan.yaml
└── test.yaml
```

Cada workflow tiene un objetivo particular que corresponde a

- `build_scan_push.yaml` - Se encarga de realizar el buildeo de la imagen de contenedor taggeando la imagen con el `short-sha` del commit, el tag latest y un tag interno, realizar un escaner de seguridad del
  artefacto utilizando el tag interno (y fallar en caso de que la imagen tenga vulnerabilidades que superen el nivel de CVE `high`), hacer inicio de sesion y pushear la imagen en el container registry
  (Google Artifact Registry) con las credenciales almacenadas en el repositorio como secretos y variables (segun su nivel de confidencialidad) y finalmente actualizar el tag de la imagen en el repositorio
  para que pueda ser utilizada por el flujo de despliegue.
- `deploy_prod.yaml` - Realiza el proceso de despliegue del servicio en el ambiente prod usando las definiciones definidas en `/infrastructure/prod` y utilizando las variables de entorno definidas a nivel de repositorio
- `deploy_staging.yaml` - Realiza el proceso de despliegue del servicio en el ambiente staging usando las definiciones definidas en `/infrastructure/staging` y utilizando las variables de entorno definidas a nivel de repositorio
- `pipeline.yaml` - Define el flujo de dependencias y condiciones del flujo principal de integracion y despliegue
- `source_scan.yaml` - Realiza un proceso de validación de codigo fuente
- `test.yaml` - Ejecuta los unit tests definidos para cada uno de los paquetes del microservicio a desplegar

Adicionalmente, para que el pipeline se ejecute correctamente es necesario configurar los secretos y variables

- Variables:
  - CR_URL : Repositorio donde se almacena la imagen de contenedor (Google Artifact Registry)
  - CR_LOCATION : Localizacion de Google Artifact Registry (Google Artifact Registry)
  - CR_IMG : Imagen de contenedor, se debe inicializar con un valor "placeholder" ya que es actualizada por el pipeline
  - GCP_PROJECT : Proyecto a usar sobre Google Cloud Platform
  - ENDPOINT_URL : URL del endpoint que consume el microservicio
- Secretos:
  - GCP_SA : Google Cloud Platform service-account codificado en base64
  - GITHUB_TOKEN : Token de Github con permisos de escritura sobre las variables del repositorio Requerido para actualizar el valor de la variable CR_IMG

#### Workflows triggers

- El workflow definido `pipeline.yaml` se ejecuta unicamente cuando se completa el proceso definido en `source_scan.yaml` y solamente cuando éste último fue realizado por un evento de `push` a la rama `main`
- En el caso de los workflows definidos en `build_scan_push.yaml`, `deploy_staging.yaml` y `deploy_prod.yaml` son ejecutados de forma
  secuencial y en el mismo orden nombrados al ejecutar el workflow de `pipeline.yaml`
- Por otra parte el workflow definido en `source_scan.yaml` se ejecuta ante un evento de tipo `push` a cualquier rama
- Finalmente el workflow definido en `test.yaml` se ejecuta ante cualquier evento de tipo `push` y `pull_request`

#### Resumen

A continuación se presentan tablas resumen de los workflows antes mencionados

##### Workflows ejecutado por push events

```
Stage  Job ID          Job name        Workflow name                   Workflow file         Events
0      source-scan     source-scan     Source code vulnerability scan  source_scan.yaml      push
```

##### Workflow ejecutado por push o pull request events

```
Stage  Job ID          Job name        Workflow name                   Workflow file         Events
0      test            test            Run go tests                    test.yaml             push,pull_request
```

##### Workflow ejecutado por workflow_run.completed event de source-scan si el origen es un push event a main

```
Stage  Job ID          Job name        Workflow name                   Workflow file         Events
0      exec_artifact   exec_artifact   CI/CD Workflow pipeline         pipeline.yaml         workflow_run
1      deploy_staging  deploy_staging  CI/CD Workflow pipeline         pipeline.yaml         workflow_run
2      deploy_prod     deploy_prod     CI/CD Workflow pipeline         pipeline.yaml         workflow_run
```

> Considerar que job-id.exec_artifact apunta a build_scan_push.yaml, job-id.deploy_staging apunta a deploy_staging.yaml y job-id.deploy_prod apunta a deploy-prod.yaml

##### Workflows que son ejecutados por llamadas de otros workflows (workflow_call)

```
Stage  Job ID          Job name        Workflow name                   Workflow file         Events
0      artifact        artifact        Build, scan and push artifact   build_scan_push.yaml  workflow_call
0      deploy-prod     deploy-prod     Deploy to PROD                  deploy_prod.yaml      workflow_call
0      deploy-staging  deploy-staging  Deploy to STAGING               deploy_staging.yaml   workflow_call
```

### Flujo de Git soportado

Considerando los workflows definidos en la sección anterior, el flujo utilizado es GithubFlow [(referencia)](https://githubflow.github.io/), las implicancias
de su aplicación son las siguientes

1. Todo lo que se escriba sobre la rama `main` puede ser desplegado, para ello se tienen las siguientes condiciones

- La rama `main` es rama protegida (requiere configuracion sobre el repositorio) y solo puede ser escrita por medio de un `Pull request`
- Toda otra rama es considerada WIP (trabajo en curso), debe ser nombrada con el nombre descriptivo de para que se usa y es validada en cada push
  (aplican los flujos de `code_scan.yaml` y `test.yaml`).

2. Las ramas adicionales desde `main` tienen nombres descriptivos, en este caso se usan los nombres `feature/*` aunque el wokflow admite cualquier nombre
3. Una vez que el cambio pasa a `main` **debe** ser desplegado, en este caso, aplica el workflow definido en `pipeline.yaml` que asegura que al completarse
   el `push` a la rama `main` (desencadenado por el push event generado por el merge completed event) y el workflow de `code_scan.yaml` es exitoso, entonces
   pasa al proceso de buildeo y despliegue, lo que generaria los 2 entornos `staging` y `prod` con la nueva imagen

### IaC

Para el proceso de despliegue en el Pipeline de CI/CD se utiliza Terraform (directorio `infrastructure`), cada ambiente es separado en su correspondiente directorio.

Ambos ambientes son una copia exacta, a excepción de que el ambiente `prod` está configurado para levantar hasta 100 contenedores mediante el autoscaler del servicio. Por otro
lado el ambiente `staging` está limitado a 2 contenedores. De modo de que el despliegue sea exitoso, ambos recursos requieren las siguientes configuraciones y variables

- Se requiere de antemano la creación de un proyecto en Google Cloud Platform, el bucket `codetest-bv-aa` para almacenar el estado de los recursos (versionados \*.tfstate)
- Se debe configurar la variable de entorno `GOOGLE_CREDENTIALS` con una service account valida con los permisos suficientes y necesarios para los recursos definidos.
- Las APIs de IAM y Google Compute Engine deben ser activadas para ejecutar el despliegue para el proyecto a utilizar.
- Se deben definir las variables descritas en `<ambiente>/variables.tf` para desplegar el recurso, definidas en cada uno de los wokflows, estas corresponden a
  - gcp-project : el proyecto definido en GCP
  - gcp-region : (opcional, por defecto "us-central1") la region a usar
  - bv-ms-img : la imagen de contenedor a utilizar
  - endpoint-url : la url del endpoint que consume el microservicio
