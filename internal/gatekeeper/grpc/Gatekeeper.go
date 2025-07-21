package grpc

import (
	"context"

	"github.com/bignyap/go-admin/internal/common"
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

	output, err := g.Service.ValidateRequest(ctx, input)
	if err != nil {
		return nil, err
	}

	orgStruct, err := common.ConvertProtoStruct(output.Organization)
	if err != nil {
		return nil, err
	}

	endpointStruct, err := common.ConvertProtoStruct(output.Endpoint)
	if err != nil {
		return nil, err
	}

	subStruct, err := common.ConvertProtoStruct(output.Subscription)
	if err != nil {
		return nil, err
	}

	return &pb.ValidateRequestResponse{
		Organization: orgStruct,
		Endpoint:     endpointStruct,
		Subscription: subStruct,
		Remaining:    int64(output.Remaining),
	}, nil
}
