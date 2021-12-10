# users-micro
go micro service using PostgreSQL

## How to build?

from project folder executes commands below, it will generate binary files in the `./bin/` folder

* just build with current operating system
```sh
make build
```

* build for a linux distro operating system
```sh
make build-linux
```

## How to run a test environment quickly?

1. make sure you have docker-compose installed.
2. give initdb.sh execution permissions.
```sh
chmod +x sql/ddl/initdb.sh
```
3. run the docker compose.
```sh
docker-compose up --build
```

or run this shortcut

```sh
make run-local
```

3. once you finished to use the environment, follow these steps

    * ctrl + c
    * make clean-local

## How to test?

from project folder run the following command

```sh
go test ./...
```


docker run -d -p 80:80 -v /var/run/docker.sock:/tmp/docker.sock:ro jwilder/nginx-proxy

## Useful Queries

* look for a user who is happy

```sql
SELECT * FROM public.jobseeker WHERE skills @> '["happy"]';
```
