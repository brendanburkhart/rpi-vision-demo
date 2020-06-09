# RPi Vision Demo

This is a rough demo of potential stack for FRC vision coprocessor software. The software would be split into three sections - the vision processing module, a web server, and the configuration dashboard. Communication between the OpenCV module and the server is accomplished via gRPC.

The vision module is OpenCV, written in C++, and would run pipelines configured via a standard set of options. It publishes a camera stream via WPILib's CameraServer that can be easily accessed via url or on the driver station dashboard, and could be easily made to publish additional streams. NetworkTables can be utilized to provide pipeline output data to consumers such as software on the Roborio, or the driver station. Finally, the module provides a gRPC server, which could be used to configure the vision pipelines.

The Go web server functions as both a static file server to serve the configuration dashboard, and provides a RESTful API to configure the vision pipelines, rather than having the dashboard communicate with the vision module directly via gRPC.

The dashboard consumes the CameraServer stream, and would be used to configure and calibrate the vision pipelines via the web server's API.

## Building

Note that installing the required dependencies for the C++ vision module will likely necessitate building several libraries from source and could take upwards of several hours, particularly on less powerful hardware such as a Raspberry Pi.

You will need Go and CMake installed to build this demo. After you've cloned or downloaded this repository, you will need to install several dependencies for the C++ vision module. While you can install these anyway you like, it is recommended you install them via [vcpkg](https://github.com/Microsoft/vcpkg) and use the vcpkg toolchain for CMake. If you install via vcpkg and want to build with Visual Studio, you need to run `vcpkg integrate install` so Visual Studio can find the packages. You will need to install [OpenCV](https://github.com/opencv/opencv), [WPILib](https://github.com/wpilibsuite/allwpilib), [gRPC](https://github.com/grpc/grpc), and [protobuf](https://github.com/protocolbuffers/protobuf). Using vcpkg, this can be accomplished by

```
vcpkg install opencv wpilib[cameraserver] grpc
```

(grpc will install protobuf as a dependency).

At this point, Visual Studio should be able to open the vision subfolder and generate build scripts from CMake and build the OpenCV vision module. You will need to add a CMake configuration for the vcpkg triplet you used - by default, vcpkg will install x86 packages for Windows, while Visual Studio's default config is x64. You can also add configurations targeting local WSL or remote Linux machines.

Building the webserver can be accomplished by navigating to `webserver/cmd/server` and executing `go build`.

Running both the vision module and webserver simultaneously, and navigating to `http://localhost:8080` should display your default camera feed run through the Canny edge detection algorithm. You can adjust the gradient thresholds used via either the REST or gRPC API.

### Modifying the gRPC service

If you plan on making changes to the gRPC API, you will need version 3 of the protobuf compiler, `protoc`. See https://grpc.io/docs/protoc-installation/ for installation instructions on Linux or MacOS, as well as binary downloads for all platforms. On WIndows, protoc can be easily installed with the `protobuf` package from [Scoop](https://github.com/lukesampson/scoop).

After modifying the gRPC service or message definitions, you will need to recompile the service definition for both C++ and Go. The service definition `pipeline_controller.proto`is duplicated between `/vision/src/grpc-common` and `/webserver/grpc-common` (just for this demo, it would normally be a Git submodule). In `grpc-common`, execute `protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative pipeline_controller.proto` to compile for Go, and `protoc -I . --grpc_out=. --plugin=protoc-gen-grpc=$grpc_cpp_plugin ./pipeline_controller.proto`, followed by `protoc -I . --cpp_out=. ./pipeline_controller.proto` to compile for C++. This assumes `$grpc_cpp_plugin` contains the path to the grpc C++ plugin, likely found in `~/vcpkg/packages/grpc_x86-windows/tools/grpc/grpc_cpp_plugin.exe` or similar.
