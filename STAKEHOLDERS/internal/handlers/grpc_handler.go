package handler

import (
	pb "COMMON/stakeholders/proto"
	"context"
	"fmt"
	"stakeholders/internal/models"
	"stakeholders/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StakeholdersHandler struct {
	pb.UnimplementedStakeholdersServiceServer
	service *service.StakeholdersService
}

func NewStakeholdersHandler(service *service.StakeholdersService) *StakeholdersHandler {
	return &StakeholdersHandler{service: service}
}

// --- CUSTOMER ---
func (h *StakeholdersHandler) CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest) (*pb.CustomerResponse, error) {
	fmt.Printf("Received Create Request: %+v\n", req)
	c, err := h.service.CreateCustomer(req.UserId, req.FirstName, req.LastName, req.Address, req.Latitude, req.Longitude)
	if err != nil {
		fmt.Printf("ERROR creating customer: %v\n", err)
		return nil, status.Errorf(codes.Internal, "Detailed Error: %v", err)
	}
	return &pb.CustomerResponse{
		UserId:    c.UserId.String(),
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Address:   c.Address,
		Latitude:  c.Latitude,
		Longitude: c.Longitude,
	}, nil
}

func (h *StakeholdersHandler) GetCustomer(ctx context.Context, req *pb.GetRequest) (*pb.CustomerResponse, error) {
	c, err := h.service.GetCustomer(req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Customer not found")
	}
	return &pb.CustomerResponse{
		UserId:    c.UserId.String(),
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Address:   c.Address,
		Latitude:  c.Latitude,
		Longitude: c.Longitude,
	}, nil
}

// --- DELIVERY PERSON ---
func (h *StakeholdersHandler) CreateDeliveryPerson(ctx context.Context, req *pb.CreateDeliveryPersonRequest) (*pb.DeliveryPersonResponse, error) {
	// Map Enum to String
	vehicleStr := "CAR" // Default
	switch req.Vehicle {
	case pb.VehicleType_BIKE:
		vehicleStr = "BIKE"
	case pb.VehicleType_SCOOTER:
		vehicleStr = "SCOOTER"
	}

	dp, err := h.service.CreateDeliveryPerson(req.UserId, req.FirstName, req.LastName, vehicleStr)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create delivery person")
	}
	return mapDeliveryToProto(dp), nil
}

func (h *StakeholdersHandler) GetDeliveryPerson(ctx context.Context, req *pb.GetRequest) (*pb.DeliveryPersonResponse, error) {
	dp, err := h.service.GetDeliveryPerson(req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Delivery person not found")
	}
	return mapDeliveryToProto(dp), nil
}

func (h *StakeholdersHandler) UpdateWorkingStatus(ctx context.Context, req *pb.UpdateWorkStatusRequest) (*pb.DeliveryPersonResponse, error) {
	dp, err := h.service.UpdateWorkingStatus(req.UserId, req.IsWorking)
	if err != nil {
		return nil, err
	}
	return mapDeliveryToProto(dp), nil
}

func (h *StakeholdersHandler) CreateRestaurantWorker(ctx context.Context, req *pb.CreateWorkerRequest) (*pb.WorkerResponse, error) {
	w, err := h.service.CreateWorker(req.UserId, req.FirstName, req.LastName)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed")
	}

	return &pb.WorkerResponse{
		UserId:    w.UserId.String(),
		FirstName: w.FirstName,
		LastName:  w.LastName,
	}, nil
}

func (h *StakeholdersHandler) GetAllRestaurantWorkers(ctx context.Context, req *pb.GetAllWorkersRequest) (*pb.GetAllWorkersResponse, error) {
	workers, err := h.service.GetAllWorkers()
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed")
	}

	var protoWorkers []*pb.WorkerResponse
	for _, w := range workers {
		protoWorkers = append(protoWorkers, &pb.WorkerResponse{
			UserId:    w.UserId.String(),
			FirstName: w.FirstName,
			LastName:  w.LastName,
		})
	}
	return &pb.GetAllWorkersResponse{Workers: protoWorkers}, nil
}

func (h *StakeholdersHandler) GetRestaurantWorker(ctx context.Context, req *pb.GetRequest) (*pb.WorkerResponse, error) {
	w, err := h.service.GetWorker(req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "Worker not found")
	}
	return &pb.WorkerResponse{
		UserId:    w.UserId.String(),
		FirstName: w.FirstName,
		LastName:  w.LastName,
	}, nil
}

// Helper
func mapDeliveryToProto(dp *models.DeliveryPerson) *pb.DeliveryPersonResponse {
	vType := pb.VehicleType_CAR
	switch dp.Vehicle {
	case models.VehicleBike:
		vType = pb.VehicleType_BIKE
	case models.VehicleScooter:
		vType = pb.VehicleType_SCOOTER
	}

	return &pb.DeliveryPersonResponse{
		UserId:        dp.UserId.String(),
		FirstName:     dp.FirstName,
		LastName:      dp.LastName,
		Vehicle:       vType,
		IsWorking:     dp.IsWorking,
		DeliveryCount: int32(dp.DeliveryCount),
	}
}
