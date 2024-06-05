# Infrastructure HCL

IaC para el despliegue de la aplicacion `microservice` en Google Cloud Platform usando Google Cloud Run v2

### Estructura del directorio

Se definen 2 estados (ambientes) de despliegue `staging` y `prod` cada uno con sus recursos correspondientes, siguiendo las recomendaciones de organizacion de ambientes
conforme a la documentación oficial de hashicorp [(link)](https://developer.hashicorp.com/terraform/tutorials/modules/organize-configuration#separate-states), en el potencial
caso de que los ambientes requieran el despliegue de recursos distintos.

```
infrastructure
├── .gitignore
├── README.md
├── prod
│   ├── .terraform.lock.hcl
│   ├── main.tf
│   ├── outputs.tf
│   ├── providers.tf
│   ├── terraform.tfvars
│   └── variables.tf
└── staging
    ├── .terraform.lock.hcl
    ├── main.tf
    ├── outputs.tf
    ├── providers.tf
    └── variables.tf
```

Para cada ambiente existen los siguientes archivos de definición:

- `main.tf`: Recursos que son desplegados, correspondientes a un Google Cloud Run Service v2 y un IAM Binding para el servicio de modo de aceptar el trafico red desde todos los servicios
- `outputs.tf`: Valores de salida de los recursos definidos, en este caso unicamente la url generada para el Cloud Run Service
- `providers.tf`: Archivo de definición y configuracion de proveedores, se utiliza Google Cloud Storage para el mantenimiento de los archivos de estado y se utiliza el proveedor oficial de Google Cloud Platform en
  su última versión estable para la creación de los recursos
- `variables.tf`: Variables utilizadas para la creación de los recursos, se definen las variables siguientes
  - gcp-project: Proyecto de Google Cloud Platform a utilizar
  - gcp-region: Región de Google Cloud Platform donde se despliegan los recursos (con valor por defecto: "us-central1")
  - bv-ms-img: Imagen de contenedor del microservicio a utilizar
  - endpoint-url: URL del endpoint al que debe apuntar el microservicio, utilizado como variable de entorno por éste

### Consideraciones para su uso

Conforme a la definición de recursos en los archivos correspondientes, existen los siguientes requerimientos para la aplicación de los estados definidos

- Recursos de GCP
  - Se requiere una cuenta de servicio o rol equivalente para la ejecución de los planes definidos
  - Debe existir el bucket `tf-state-codetest-bv-aa` y un proyecto creado para la ejecución del plan
  - La cuenta de servicio o rol debe tener los permisos suficientes y necesarios para la creación de los recursos definidos.
- La imagen de contenedor definida por la variable `bv-ms-img` debe encontrarse disponible para la descarga y utilización
