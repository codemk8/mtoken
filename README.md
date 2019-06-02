# mtoken

## Prerequisite

Run the other microservice [muser](https://github.com/codemk8/muser). Assume it is run in the same `localhost` and the authentication endpoint is `http://localhost:8000/v1/user/auth`.

## Build and run mtoken service

```bash
make
./bin/mtoken -user_service http://localhost:8000/v1/user/auth
# On the first run, it'll generate pub and private key pairs. Keep them in a safe place. If restart, remember to specify them in the command line:
./bin/mtoken -user_service http://localhost:8000/v1/user/auth -priv_key private_key_file_path -pub_key public_key_file_path
```

## Test

```
$ curl -X POST --user test_user:secret http://127.0.0.1:8001/v1/token/issue
eyJhbGciOiJSUzI1NiIsImtpZCI6IiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NjIwNzk5MjMsImlzcyI6InRlc3RfdXNlciIsInN1YiI6InlvdXJhcHAifQ.VP1vL4Gtdsa4IHdCmbAIPswnmPgAcWpCNVlHz_hKfgF1oBZ_YC3ifWbO49vZXYxDNylWmubsHLn4Q196T7gjntLF2bnGFSvFZ3pLWGaCmVU-QcOXE4IHQSUu0gP0mrP-tfCB-aaaYGFDEoKFcl8ECPwPbYbEAx1WJFXf_2y5M2pDXQ6fybTY9NOqdllc1OI4Z6YsJjYNEIFT41FCKHErgx1oFpQ20CHxzuOCXkTa2rfIvs-QZMd9_8qYKI4ZaTAEfLd9ZEqyLMc6Jr2rylSmgOrGDFoldhb4khmLR5FcQHyjOFLFJqRxEROT-IDdCu94A7OeN914CaE8vNLytCrmhCY0xb03F_cKE22i6PWOygMISS24UrEwVFnHP-GLiniO5-vawdqiuBq_JKtwYbA0i5sCp3iXWEoiUdAMNI-KY_X_7xC502i8NA7mMVLELQpLoSbz1_J0NOBs1qpkvWw-1uGKaddHVpm1wXjBx4qOGg75h5mW-aaG2XJAb4gmQ6Kv5MXu4qszvUmeSKpONW9a2OAC5h-pDFUnwyySaN9JAtKXV4Q0PzxIjt58EaSV6xdtPd05yxZkeeibgU0RBvLxvUQaomaxELopFKF1YX7QP_-AfNOfZDZVBIHjM9ZaLqEHJ06C7uaKjGwRMsfC7AaAJ7nNwPvEiKxIx0L7D8

# visit jwt.io, copy the above token string to check claims are:
{
  "exp": 1562079923,
  "iss": "test_user",
  "sub": "yourapp"
}
```