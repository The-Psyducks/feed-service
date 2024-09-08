# go-microservice-template



## Como usar este template

Para iniciar el proyecto de go, correr el siguiente comando:

```bash
cd server
go mod init ${PROJECTNAME}
go mod tidy
```
para reflejar los ultimos cambios:
```bash
docker compose build
```

para correr el servicio:

```bash

docker compose up service
```

para correr los tests:

```bash
docker compose up tests
```


Para cambiar el entorno en el que se encuentra la aplicacion, se debe cambiar la variable de entorno ENVIRONMENT en el archivo `server/.env`


Ojo! que sacar los .env del repo