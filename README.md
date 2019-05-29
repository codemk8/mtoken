# mtoken

## Prerequisite

Run the other microservice [muser](https://github.com/codemk8/muser). Assume it is run in the same `localhost` and the authentication endpoint is `http://localhost:8000/v1/user/auth`.

## Build and run mtoken service

```bash
make
./bin/mtoken -user_service http://localhost:8000/v1/user/auth
```

## Test

```
$ curl -X POST --user test_user:secret http://127.0.0.1:8001/v1/token/issue
```