# rev

**Rev** is a simple (***rev***)erse proxy built with Go for all kinds of REST servers.

Rev does simple prevention of DDOS attacks by establishing rate limits cooldowns for the incoming clients.

Moreover, Rev also provides a way to view telemetry data coming from the requests to the actual server. Persistent data collection is available through redirection on `stdout`.

### Simple Distributed Denial of Service Protection

### Viewing Additional Telemetry Data

## Set Up

##### Requirements

- Go version 1.22 or above

#### Configurations
The `.example.env` has an example of available environment variables:
- `PROXY_PORT` the (localhost) port that the proxy server will run on
- `BACKEND_SERVER_PORT` the port of the actual REST server this proxy will forward requests to
- `SHOW_REQUEST_DETAIL` enables additional telemetry data for incoming requests

See [Viewing Additional Telemetry Data](#viewing-additional-telemetry-data) for more information on `SHOW_REQUEST_DETAIL`

To run the proxy server:
`go run main.go`

We also have a test REST API server that can be run wih `go run test_backend/server.go`. See [server.go](test_backend/server.go)
