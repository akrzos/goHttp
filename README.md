# goHttp

Simple golang http server configured for kubernetes probes

## Configuration

| Env Var                       | Default | Description                                                                                |
| ----------------------------- | ------- | ------------------------------------------------------------------------------------------ |
| `PORT`                        | `8000`  | Port to listen on                                                                          |
| `LISTEN_DELAY_SECONDS`        | `10`    | Delay before application will listen for requests in seconds                               |
| `LIVENESS_DELAY_SECONDS`      | `2`     | Delay before application will report HTTP 200 on /livez in seconds                         |
| `READINESS_DELAY_SECONDS`     | `10`    | Delay before application will report HTTP 200 on /readyz in seconds                        |
| `RESPONSE_DELAY_MILLISECONDS` | `0`     | Delay for responsiveness of all endpoints in milliseconds                                  |
| `LIVENESS_SUCCESS_MAX`        | `0`     | Maximum number of /livez replies before and http 503 is returned, 0 means infinite replies |


## Run

```console
[akrzos@fedora goHttp]$ source env.sh
[akrzos@fedora goHttp]$ go run main.go
2021/08/16 16:55:08 Starting the server...
2021/08/16 16:55:08 Using port 8000
2021/08/16 16:55:08 Using listen delay 3s
2021/08/16 16:55:08 Using live delay 3s
2021/08/16 16:55:08 Using readiness delay 3s
2021/08/16 16:55:08 Using response delay 500ms
2021/08/16 16:55:08 Using liveness success max 3
2021/08/16 16:55:08 Starting listen delay...
2021/08/16 16:55:11 Completed listen delay
2021/08/16 16:55:11 The service is listening on port 8000
2021/08/16 16:55:11 Starting readyz delay...
2021/08/16 16:55:11 Starting livez delay...
2021/08/16 16:55:12 Completed livez delay
2021/08/16 16:55:14 Completed readyz delay
2021/08/16 16:55:32 /readyz request when ready
...
```

## Endpoints

| Endpoint | Purpose                        |
| -------- | ------------------------------ |
| /home    | Testing                        |
| /readyz  | Readiness check for kubernetes |
| /livez   | Liveness check for kubernetes  |
| /crash   | Crashes server                 |


```console
[akrzos@fedora goHttp]$ time curl -i http://127.0.0.1:8000/home
HTTP/1.1 200 OK
Date: Mon, 16 Aug 2021 19:05:21 GMT
Content-Length: 23
Content-Type: text/plain; charset=utf-8

/home request processed
real	0m0.511s
user	0m0.005s
sys	0m0.005s
[akrzos@fedora goHttp]$ curl -i http://127.0.0.1:8000/readyz
HTTP/1.1 200 OK
...

/readyz request processed
[akrzos@fedora goHttp]$ curl -i http://127.0.0.1:8000/livez
HTTP/1.1 200 OK
...

/livez request processed
[akrzos@fedora goHttp]$ curl -i http://127.0.0.1:8000/livez
HTTP/1.1 200 OK
...

/livez request processed
[akrzos@fedora goHttp]$ curl -i http://127.0.0.1:8000/livez
HTTP/1.1 200 OK
...

/livez request processed
[akrzos@fedora goHttp]$ curl -i http://127.0.0.1:8000/livez
HTTP/1.1 503 Service Unavailable
Date: Mon, 16 Aug 2021 19:05:42 GMT
Content-Length: 0
```

Note `LIVENESS_SUCCESS_MAX` was set to `3` thus after 3 requests the /livez endpoint returns a 503.
