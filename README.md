# mtoken

## Prerequisite

Run the other microservice [muser](https://github.com/codemk8/muser). Assume it is run in the same `localhost`.

## Build and run mtoken service

```bash
make
./bin/mtoken
```

## Test

```
$ curl -X POST --user test_user:secret http://127.0.0.1:8001/v1/token/login
```