# steve


##TODO

- Use id from Meta as index
- Determine if plural should be kept or not
- Go plugin

## Learn


## Methods

**Alternative 1**
- Get
- Set
- Delete

1. Template protobuf files
  - Add methods
2. Generate gRPC
  - Run protoc
3. Template plugin
  - Add protobuf into to plugin
3. Build shared library
  - Compile protobuf code and plugin together

- add inotify
  - resources templates -> generated new protobuf files
  - profobuf files -> generate new stubs and new plugins
  - plugins -> try load

## Configuration

**steve config**

- run dir
  - generated protobuf
  - generated templates
- plugins
- plugin based config parsing
  - mongo parser

**resource config**


References
----------


https://github.com/grpc/grpc-go
