cmake_minimum_required(VERSION 3.17)
project(vision_demo
    VERSION 0.0.1
    LANGUAGES C CXX)

set(CMAKE_CXX_EXTENSIONS OFF)
set_property(GLOBAL PROPERTY USE_FOLDERS ON)

find_package(OpenCV REQUIRED COMPONENTS core imgproc highgui)
find_package(gRPC CONFIG REQUIRED)
find_package(cameraserver CONFIG REQUIRED PATH_SUFFIXES "/wpilib")
find_package(cscore CONFIG REQUIRED PATH_SUFFIXES "/wpilib")
find_package(ntcore CONFIG REQUIRED PATH_SUFFIXES "/wpilib")
find_package(wpiutil CONFIG REQUIRED PATH_SUFFIXES "/wpilib")

add_executable(vision src/main.cpp "src/pipeline_controller.cpp" "src/pipeline.cpp" "src/grpc-common/pipeline_controller.pb.cc" "src/grpc-common/pipeline_controller.grpc.pb.cc" "src/pipeline.hpp")

target_include_directories(vision PRIVATE gRPC::grpc++)

target_link_libraries(vision PRIVATE ${OpenCV_LIBS})
target_link_libraries(vision PRIVATE gRPC::grpc++)
target_link_libraries(vision PRIVATE cameraserver)
target_link_libraries(vision PRIVATE cscore)
target_link_libraries(vision PRIVATE ntcore)
target_link_libraries(vision PRIVATE wpiutil)

target_compile_features(vision PUBLIC cxx_std_17)
set_target_properties(vision PROPERTIES CXX_EXTENSIONS OFF)
