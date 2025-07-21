package grpc

import (
	"context"

	pb "github.com/bignyap/go-admin/internal/gatekeeper/proto"
	gatekeeping "github.com/bignyap/go-admin/internal/gatekeeper/service/GateKeeping"
)

type GatekeeperGRPCHandler struct {
	pb.UnimplementedGatekeeperServiceServer
	Service *gatekeeping.GateKeepingService
}

func NewGatekeeperGRPCHandler(service *gatekeeping.GateKeepingService) *GatekeeperGRPCHandler {
	return &GatekeeperGRPCHandler{Service: service}
}

func (g *GatekeeperGRPCHandler) RecordUsage(ctx context.Context, req *pb.RecordUsageRequest) (*pb.RecordUsageResponse, error) {
	input := &gatekeeping.RecordUsageInput{
		Method:           req.Method,
		Path:             req.Path,
		OrganizationName: req.OrganizationName,
	}

	cost, err := g.Service.RecordUsage(ctx, input)
	if err != nil {
		return nil, err
	}

	return &pb.RecordUsageResponse{Cost: cost}, nil
}

func (g *GatekeeperGRPCHandler) ValidateRequest(ctx context.Context, req *pb.ValidateRequestRequest) (*pb.ValidateRequestResponse, error) {
	input := &gatekeeping.ValidateRequestInput{
		OrganizationName: req.OrganizationName,
		Method:           req.Method,
		Path:             req.Path,
	}

	_, err := g.Service.ValidateRequest(ctx, input)
	if err != nil {
		return &pb.ValidateRequestResponse{Valid: false, Message: err.Error()}, nil
	}

	return &pb.ValidateRequestResponse{Valid: true, Message: "Valid"}, nil
}
