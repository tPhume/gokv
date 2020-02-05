# GoKv
GoKv is an in memory key-value pair database using Btree as the default data structure.
It supports insert,update,search and remove operations on the data structure.
The operations are exposed through a simple REST api and gRPC. 

## Getting started
There are two simple ways to start the example main service.
### `Docker`
The easiest way to start the example main service is to build the Dockerfile.
and the run it. The service REST (details about the REST api is further down) server is exposed through port `:8888` 
and the gRPC server is exposed through port `:9999`.

### `Go compiler`
If you have Go 1.13 installed, then you can simply build and run the main service as you would a normal Go program.
The just execute the binary file given.
 
### `REST`
The REST api is accessed through `/store/v1/:key`.
* **POST** - must include json body, which will be used as the value for the key-value pair. Does not return value.
* **PATCH** - must include json body, will replace existing value with given json. Does not return value.
* **GET** - no body needed, will search for given key and return the value in json format.
* **DELETE** - no body needed, will delete given key from the store. Does not return value.

## Directories
### `examples`
The examples directory contains example on running the REST server and the gRPC server (and the client).
The `main` application runs both the REST server and the gRPC server, and is used as the entrypoint
for our docker image.

### `btree`
The btree directory contains the btree implementation.
The package exposes Btree which is an encapsulation of the
node struct that does most of the heavy lifting.

### `store`
The store directory contains the interface Store that needs to be implemented by any
data structure that wants to allow itself as an alternative the btree data structure.

### `kv`
The kv directory contains the the REST server which depends on Gin framework,
and also the gRPC server alongside its protobuf definition. The default of both the REST and gRPC server uses
btree with a minimum of 3 degree (easy to visualize and check).