# Taller Docker/contenedores

1. Deberá crear un contenedor de Docker para ejecutar el backend escrito en GO,
   tenga presente que, para agregar las dependencias del proyecto se debe
   ejecutar: `go mod init main` y `go mod tidy` además para compilar el
   proyecto debe: `go build -o server .` esto producirá un binario `server` en
   el directorio donde se compile el código; el backend deberá exponer el
   puerto 8080 

2. Deberá crear un contenedor de *postgresql* que exponga el puerto 5432, la base
   de datos debe tener las siguientes variables de entorno: 
- `POSTGRES_PASSWORD`: password
- `POSTGRES_USER`: postgres
- `APP_DB_USER`: postgres
- `APP_DB_PASS`: password
- `APP_DB_NAME`: postgres

Además de esto el contenedor deberá tener un volumen para garantizar
persistencia de datos, para inicializar la misma deberá ejecutar el script de
bash dado 

1. Deberá crear un contenedor con *nginx* instalado, este deberá tener la
   configuración facilitada, y deberá ejecutar, antes de ejecutar *nginx*, el
   script de *bash* dado 

2. Deberá probar, que cada uno de los contenedores funcione correctamente y que
   la aplicación funcione

## Parte 2: Docker compose

Para la segunda parte del ejercicio deberá crear un `docker compose` del proyecto

```yaml
services: # servicios que vamos a crear dentro de nuestro proyectos 
    service_name: # Nombre del servicio especifico puede ser:
                  # Backend, Fronted un Proxy, una base de datos

        build: # como vamos a contruirla, puede ser con una imagen 
               # de docker o desde un docker file 
        networks: # definimos la red a la que vamos a cenectar nuestro 
                  # servicio 
        ports: # definimos puertos internos y externos del servicio 

        depends_on: # definimos si nuestro servicio depende de otro 

        deploy: # definimos caracteristicas para el desplieque de 
                # nuestro servicio, como recursos 
                # y replicas 

        volumes: # definimos un espacio de disco donde vamos a almacenar 
                 # nuestra información 
        enviroment: # variables de entorno del servicio
```
