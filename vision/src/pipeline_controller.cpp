#include "pipeline_controller.hpp"

ServerImpl::ServerImpl (Pipeline* pipeline) : pipeline (pipeline) { }

grpc::Status ServerImpl::GetThresholds (grpc::ServerContext* ctx, const google::protobuf::Empty* _, Vision::Thresholds* thresholds) {
    thresholds->set_high (pipeline->high ());
    thresholds->set_low (pipeline->low ());

    return grpc::Status::OK;
}

grpc::Status ServerImpl::SetThresholds (grpc::ServerContext* ctx, const Vision::Thresholds* thresholds, google::protobuf::Empty* _) {
    pipeline->updateThresholds (static_cast<int>(thresholds->low ()), static_cast<int>(thresholds->high ()));

    return grpc::Status::OK;
}

void ServerImpl::Run (Pipeline* pipeline) {
    std::string server_address ("0.0.0.0:50051");

    ServerImpl instance = ServerImpl (pipeline);

    grpc::ServerBuilder builder;
    builder.AddListeningPort (server_address, grpc::InsecureServerCredentials ());
    builder.RegisterService (&instance);

    std::unique_ptr<grpc::Server> server (builder.BuildAndStart ());
    std::cout << "Server listening on " << server_address << std::endl;

    server->Wait ();
}