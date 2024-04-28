# Project Sprint Cat Social App

CatSocial is an application where cat owners can match cats with each other

# How to run local (dev purposes)

- Create file `cats-social.env`

```go
DB_NAME=
DB_HOST=
DB_USER=
DB_PASSWORD=
DB_PORT=5432

BCRYPT_SALT=8
JWT_SECRET=
```

- run `make build-dev`
- run `make run-dev`
