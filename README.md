# ipv6world

## Getting Started
### Requires
* go version 1.10.3+
* MongoDB version 4.0.1+ (no auth)

### Startup
`go run main.go`

### Testing
No tests as of yet :(

## Services
On service startup, data from the file csv/data/GeoLite2-City-Blocks-IPv6.csv is extracted and saved as `addresses` documents in the `ipv6` database and the REST API is started. For now, the client retrieves the entire data set and uses it to create the heat map. This behavior is likely to change in the future in favor of a less resource intensive method such as load-on-demand.

### Client
* Available at `localhost:8000`

### Server
#### Return all data
`localhost:8000/api/v1/addresses`

#### Return addresses within a geographical bounding box
`localhost:8000/api/v1/addresses?bbox=-79,35,-78,36`

### TODO
#### Client
* Add inputs for bounding box values
* Delay loading data until user enters bounding box values
* Fetch and display only that data which is needed for the bounded view?
* Add markers to display count information
* Adjust heatmap color gradient to be more accessible

#### Server
* Add tests and benchmarking
* Make configurable by environment variables (port, mongo host/port, CSV file location, etc)
* Watch CSV file for changes and update automatically
* Parse CSV file for latitude and longitude values rather than rely on column numbers
* Add support for auth-enabled MongoDB
