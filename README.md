A daemon providing system informations

sysd collects the system informations, organizes them into well-structured data, and provides them in the RESTful HTTP API.

## Features

- Thin and fast
- HTTP RESTful API
- JSON format
- Less runtime dependencies
- Save time from parsing variant commands's output

## Development

Compile with the commands below:

```
make
```

and then launch the daemon with following command:

```
sudo DEBUG=1 ./sysd/sysd
```

## Usage

To launch the daemon, just run:

```
sudo ./sysd/sysd
```

To get network interfaces list:

```
curl -sL http://127.0.0.1:8/ifconfig | jq
```

To get registered API functions:

```
curl -sL http://127.0.0.1:8/apilist | jq
```
