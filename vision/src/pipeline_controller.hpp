#ifndef PIPELINE_CONTROLLER_H
#define PIPELINE_CONTROLLER_H

#include <iostream>

#include <grpcpp/grpcpp.h>

#include "grpc-common/pipeline_controller.grpc.pb.h"
#include "pipeline.hpp"

class ServerImpl final : public Vision::PipelineController::Service {
public:
    ServerImpl (Pipeline* pipeline);

    grpc::Status GetThresholds (grpc::ServerContext* ctx, const google::protobuf::Empty* _, Vision::Thresholds* thresholds) override;
    grpc::Status SetThresholds (grpc::ServerContext* ctx, const Vision::Thresholds* thresholds, google::protobuf::Empty* _) override;

    static void Run (Pipeline* pipeline);

private:
    Pipeline* pipeline;
};

#endif
