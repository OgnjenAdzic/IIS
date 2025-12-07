package handler

import (
	"context"

	"AUTH/internal/models"
	"AUTH/internal/service"
	pb "COMMON/auth/proto"

	"github.com/golang-jwt/jwt/v5"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	userService *service.UserService
	jwtService  *service.JWTService
}

func NewAuthHandler(us *service.UserService, js *service.JWTService) *AuthHandler {
	return &AuthHandler{userService: us, jwtService: js}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	// 1. Validation
	if req.Role == pb.UserRole_ROLE_UNKNOWN {
		return nil, status.Error(codes.InvalidArgument, "Role is required")
	}

	// 2. Map Proto -> Model
	domainRole := mapProtoRoleToModel(req.Role)

	// 3. Call Service
	user, err := h.userService.Register(req.Username, req.Password, string(domainRole))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Registration failed: %v", err)
	}

	return &pb.RegisterResponse{
		Id:       user.Id.String(),
		Username: user.Username,
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	// 1. Check if user exists and password is correct
	user, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
	}

	// 2. Generate Token
	// Note: string(user.Role) converts the Enum/String model to string for JWT
	token, err := h.jwtService.GenerateToken(user.Id.String(), string(user.Role), user.Username)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to generate token")
	}

	// 3. Return Response
	return &pb.LoginResponse{
		Token: token,
	}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	token, err := h.jwtService.ValidateToken(req.Token)

	// Check if token is nil or invalid
	if err != nil || !token.Valid {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	// 2. FIX: Now 'jwt' is defined so this works
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	roleStr, _ := claims["role"].(string)
	userId, _ := claims["sub"].(string)

	// Map Model -> Proto
	protoRole := mapModelRoleToProto(models.Role(roleStr))

	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: userId,
		Role:   protoRole,
	}, nil
}

// --- HELPER FUNCTIONS ---

// 3. FIX: We now use 'models.RoleX' constants.
// This makes the code safer and satisfies the "unused import" error.

func mapProtoRoleToModel(pRole pb.UserRole) models.Role {
	switch pRole {
	case pb.UserRole_ROLE_ADMIN:
		return models.RoleAdmin
	case pb.UserRole_ROLE_CUSTOMER:
		return models.RoleCustomer
	case pb.UserRole_ROLE_DELIVERY_PERSON:
		return models.RoleDeliveryPerson
	case pb.UserRole_ROLE_RESTAURANT_WORKER:
		return models.RoleRestaurantWorker
	default:
		return models.RoleCustomer
	}
}

func mapModelRoleToProto(mRole models.Role) pb.UserRole {
	switch mRole {
	case models.RoleAdmin:
		return pb.UserRole_ROLE_ADMIN
	case models.RoleCustomer:
		return pb.UserRole_ROLE_CUSTOMER
	case models.RoleDeliveryPerson:
		return pb.UserRole_ROLE_DELIVERY_PERSON
	case models.RoleRestaurantWorker:
		return pb.UserRole_ROLE_RESTAURANT_WORKER
	default:
		return pb.UserRole_ROLE_UNKNOWN
	}
}
