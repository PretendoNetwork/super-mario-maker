# Super Mario Maker replacement server
Includes both the authentication and secure servers. Works on both WiiU and 3DS

## DataStore (S3)
This server requires an S3 compatible server to store user generated content. The WiiU and 3DS only support TLS versions 1.0 and 1.1 with RSA SSL ciphers, so an S3 server supporting these is required. The server must also support presigned `POST` URLs. This does not leave many options, as nearly all S3 cloud providers no longer support one of these 2 things. Even AWS, the creators of S3 and the provider Nintendo uses, is [dropping support for TLS versions below 1.2 in December 2023](https://aws.amazon.com/blogs/security/tls-1-2-required-for-aws-endpoints/)

For this reason, the recommended setup is using [MinIO](https://min.io/) to self host your S3 server and using [Cloudflare Tunnels](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/) to act as your reverse proxy. Tunnels support TLS versions 1.0 and 1.1 with the required ciphers, essentially using it as a TLS proxy. This does come at some cost, and you now must manage data storage and security yourself through MinIO, but there are no other options at this time outside of self hosting

## Compiling

### Setup
Install [Go](https://go.dev/doc/install) and [git](https://git-scm.com/downloads), then clone and enter the repository

```bash
$ git clone https://github.com/PretendoNetwork/super-mario-maker
$ cd super-mario-maker
```

### Compiling using `go`
To compile using Go, `go get` the required modules and then `go build` to your desired location. You may also want to tidy the go modules, though this is optional

```bash
$ go get -u
$ go mod tidy
$ go build -o build/super-mario-maker
```

The server is now built to `build/super-mario-maker`

When compiling with only Go, the authentication servers build string is not automatically set. This should not cause any issues with gameplay, but it means that the server build will not be visible in any packet dumps or logs a title may produce

To compile the servers with the authentication server build string, add `-ldflags "-X 'main.serverBuildString=BUILD_STRING_HERE'"` to the build command, or use `make` to compile the server

### Compiling using `make`
Compiling using `make` will read the local `.git` directory to create a dynamic authentication server build string, based on your repositories remote origin and current commit

Install `make` either through your systems package manager or the [official download](https://www.gnu.org/software/make/). We provide a `default` rule which compiles [using `go`](#compiling-using-go)

To build using `go`

```bash
$ make
```

The server is now built to `build/super-mario-maker`

## Configuration
All configuration options are handled via environment variables

`.env` files are supported

| Name                                | Description                                                           | Required                                      |
|-------------------------------------|-----------------------------------------------------------------------|-----------------------------------------------|
| `PN_SMM_POSTGRES_URI`               | Fully qualified URI to your Postgres server                           | Yes                                           |
| `PN_SMM_AUTHENTICATION_SERVER_PORT` | Port for the authentication server                                    | Yes                                           |
| `PN_SMM_SECURE_SERVER_HOST`         | Host name for the secure server                                       | Yes                                           |
| `PN_SMM_SECURE_SERVER_PORT`         | Port for the secure server                                            | Yes                                           |
| `PN_SMM_CONFIG_S3_ENDPOINT`         | S3 server endpoint                                                    | Yes                                           |
| `PN_SMM_CONFIG_S3_ACCESS_KEY`       | S3 access key ID                                                      | Yes                                           |
| `PN_SMM_CONFIG_S3_ACCESS_SECRET`    | S3 secret                                                             | Yes                                           |
| `PN_SMM_CONFIG_S3_BUCKET`           | S3 bucket                                                             | Yes                                           |
| `PN_SMM_ACCOUNT_GRPC_HOST`          | Host name for your account server gRPC service                        | Yes                                           |
| `PN_SMM_ACCOUNT_GRPC_PORT`          | Port for your account server gRPC service                             | Yes                                           |
| `PN_SMM_ACCOUNT_GRPC_API_KEY`       | API key for your account server gRPC service                          | No (Assumed to be an open gRPC API)           |