# rev

**Rev** is a simple (***rev***)erse proxy built with Go for all kinds of REST servers.

Rev does simple prevention of DDOS attacks by establishing rate limits cooldowns for the incoming clients.

Moreover, Rev also provides a way to view telemetry data coming from the requests to the actual server. Persistent data collection is available through redirection on `stdout`.

### Simple Distributed Denial of Service Protection

The following are 2 examples of trying to spam an API endpoint directly to the server versus to the reverse proxy.

Here, the actual server is hosted on `localhost:7000`, and the reverse proxy server is hosted on `localhost:7001`.

[Demo without Reverse Proxy Server](https://i.imgur.com/NzSWk7z.mp4)

We ran a script that spam calls an API endpoint 20 times, and our actual server processed all 20 calls.

The following example shows the rate-limiting block from the proxy server.

[Demo with Reverse Proxy Server](https://i.imgur.com/sP6PXyN.mp4)

After 10 consecutive calls, we detect the rate of request is too stronger and puts the IP into a temporary cool down list and prevents it from making any interactions with the server further for the next 5 seconds.


### Viewing Additional Telemetry Data

![sQZRPz](https://github.com/user-attachments/assets/7fb45cd6-0b62-4e57-ab7e-cc1ea5be83fd)

By setting `SHOW_REQUEST_DETAIL=true`, the proxy server will display additional telemetry data including `Method`, `Path`, `URL`, `ContentLength`, and `Header`. For details on these fields, visit https://pkg.go.dev/net/http#Request for reference.

## Set Up

##### Requirements

- Go version 1.22 or above

#### Configurations
The `.example.env` has an example of available environment variables:
- `PROXY_PORT` the (localhost) port that the proxy server will run on
- `BACKEND_SERVER_PORT` the port of the actual REST server this proxy will forward requests to
- `SHOW_REQUEST_DETAIL` enables additional telemetry data for incoming requests
- `RATE_LIMIT` specifies the maximum frequency in milliseconds between 10 calls to not be considered as spam (trigger rate limit prevention)

See [Viewing Additional Telemetry Data](#viewing-additional-telemetry-data) for more information on `SHOW_REQUEST_DETAIL`

To run the proxy server:
`go run main.go`

We also have a test REST API server that can be run wih `go run test_backend/server.go`. See [server.go](test_backend/server.go)
