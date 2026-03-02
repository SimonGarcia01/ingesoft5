# Ejercicio con Docker Compose

Partes de un archivo `.yaml`: 

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

Estructura del proyecto

```bash
.
├── docker-compose.yml
│
├── db
│   └── init.sh
│
├── golang
│   ├── Dockerfile
│   └── main.go
│
└── proxy
    ├── Dockerfile
    ├── init.sh
    └── nginx.conf

```

## Definamos nuestras imágenes con los Dockerfile

Definamos la imagen del `proxy` 

```Dockerfile
FROM nginx:mainline-alpine
COPY init.sh /init.sh
RUN apk update && apk add --no-cache openssl
COPY ./nginx.conf /etc/nginx/nginx.conf
ENTRYPOINT ["/bin/sh", "init.sh"]
CMD ["nginx", "-g", "daemon off;"]
```

Definimos la imagen del `backend` 

```Dockerfile
FROM golang:1.24-alpine
WORKDIR /app
COPY . .
RUN apk --no-cache add ca-certificates
RUN go mod init main && go mod tidy && go build -o server .

EXPOSE 8080
CMD ["/app/server"]
```

NOTA: La imagen anterior puede ser *optimizada* si se aplica una configuración 
[*Multi-stage*:](https://docs.docker.com/build/building/multi-stage/) 

```Dockerfile
# Definimos un 'builder'
FROM golang:1.24-alpine AS builder-go
WORKDIR /app
COPY . .
RUN go mod init main && go mod tidy && go build -o server .

# Definimos la imagen del contenedor 'final'
FROM alpine:latest 
RUN apk --no-cache add ca-certificates
COPY --from=builder-go /app/server /app/server

EXPOSE 8080
CMD ["/app/server"]
```

## Definamos como administrar nuestra aplicación con Docker compose

Creación del proxy server: 

```yaml
  proxy:
    container_name: nginx_proxy 
    build: 
      context: proxy
      dockerfile: Dockerfile
    restart: always
    networks:
      - app
    ports:
      - 80:80
      - 443:443
    links:
      - backend
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 500M
```

Creación de la base de datos: 

```yaml
  db:
    container_name: postgres_db 
    image: postgres
    ports:
      - "5432:5432"
    restart: always
    environment:
      POSTGRES_PASSWORD: password
      POSTGRES_USER: postgres
      APP_DB_USER: postgres
      APP_DB_PASS: password
      APP_DB_NAME: postgres
    volumes:
      # https://hub.docker.com/_/postgres#initialization-scripts
      - ./db:/docker-entrypoint-initdb.d/
      - ./postgres_data:/var/lib/postgresql
    networks:
      - app
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 500M
```

Definición del `backed`: 

```yaml
  backend:
    build: 
      context: golang
      dockerfile: Dockerfile
    restart: always
    expose:
      - "8080"
    depends_on:
      - db
    networks:
      - app
    environment:
      RDS_HOSTNAME: db
    # https://docs.docker.com/reference/compose-file/deploy/
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 500M
      mode: replicated
      replicas: 4
      placement:
        max_replicas_per_node: 1
      endpoint_mode: dnsrr
      update_config:
        parallelism: 2
        delay: 10s
      restart_policy:
        condition: on-failure
```

Definición del volumen: 

```yaml
volumes:
  postgres_data:
```

Definición de la red: 

```bash
networks:
  app:
    ipam:
      driver: default
      config:
        - subnet: 192.168.200.0/24
```
