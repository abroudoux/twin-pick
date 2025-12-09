package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/charmbracelet/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	pb "github.com/abroudoux/twinpick/api/proto"
	"github.com/abroudoux/twinpick/internal/application"
	"github.com/abroudoux/twinpick/internal/domain"
)

type Server struct {
	pb.UnimplementedTwinPickServiceServer
	pickService *application.PickService
	spotService *application.SpotService
	grpcServer  *grpc.Server
	port        int
}

func NewServer(pickService *application.PickService, spotService *application.SpotService, port int) *Server {
	return &Server{
		pickService: pickService,
		spotService: spotService,
		port:        port,
	}
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s.grpcServer = grpc.NewServer()
	pb.RegisterTwinPickServiceServer(s.grpcServer, s)

	// Enable reflection for tools like grpcurl
	reflection.Register(s.grpcServer)

	log.Info("gRPC server started", "port", s.port)
	return s.grpcServer.Serve(lis)
}

func (s *Server) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}

// Pick retourne les films en commun entre plusieurs watchlists
func (s *Server) Pick(ctx context.Context, req *pb.PickRequest) (*pb.PickResponse, error) {
	if len(req.Usernames) == 0 {
		return nil, status.Error(codes.InvalidArgument, "usernames is required")
	}

	pickParams := convertPickRequest(req)

	films, err := s.pickService.Pick(pickParams)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to pick films: %v", err)
	}

	return &pb.PickResponse{
		Films:      convertFilmsToProto(films),
		TotalCount: int32(len(films)),
	}, nil
}

// Spot retourne des suggestions de films populaires
func (s *Server) Spot(ctx context.Context, req *pb.SpotRequest) (*pb.SpotResponse, error) {
	spotParams := convertSpotRequest(req)

	films, err := s.spotService.Spot(spotParams)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to spot films: %v", err)
	}

	return &pb.SpotResponse{
		Films:      convertFilmsToProto(films),
		TotalCount: int32(len(films)),
	}, nil
}

// PickStream retourne les films en streaming (un par un)
func (s *Server) PickStream(req *pb.PickRequest, stream grpc.ServerStreamingServer[pb.Film]) error {
	if len(req.Usernames) == 0 {
		return status.Error(codes.InvalidArgument, "usernames is required")
	}

	pickParams := convertPickRequest(req)

	films, err := s.pickService.Pick(pickParams)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to pick films: %v", err)
	}

	for _, film := range films {
		if err := stream.Send(convertFilmToProto(film)); err != nil {
			return err
		}
	}

	return nil
}

// SpotStream retourne les suggestions en streaming
func (s *Server) SpotStream(req *pb.SpotRequest, stream grpc.ServerStreamingServer[pb.Film]) error {
	spotParams := convertSpotRequest(req)

	films, err := s.spotService.Spot(spotParams)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to spot films: %v", err)
	}

	for _, film := range films {
		if err := stream.Send(convertFilmToProto(film)); err != nil {
			return err
		}
	}

	return nil
}

// Conversion helpers

func convertPickRequest(req *pb.PickRequest) *domain.PickParams {
	filters := convertFilters(req.Filters)
	scrapperFilters := convertScrapperFilters(req.ScrapperFilters)

	return domain.NewPickParams(
		req.Usernames,
		domain.NewParams(filters, scrapperFilters),
	)
}

func convertSpotRequest(req *pb.SpotRequest) *domain.SpotParams {
	filters := convertFilters(req.Filters)
	scrapperFilters := convertScrapperFilters(req.ScrapperFilters)

	return domain.NewSpotParams(
		domain.NewParams(filters, scrapperFilters),
	)
}

func convertFilters(f *pb.Filters) *domain.Filters {
	if f == nil {
		return domain.NewFilters(0, domain.Long)
	}

	duration := convertDuration(f.Duration)
	return domain.NewFilters(int(f.Limit), duration)
}

func convertDuration(d pb.Duration) domain.Duration {
	switch d {
	case pb.Duration_DURATION_SHORT:
		return domain.Short
	case pb.Duration_DURATION_MEDIUM:
		return domain.Medium
	default:
		return domain.Long
	}
}

func convertScrapperFilters(sf *pb.ScrapperFilters) *domain.ScrapperFilters {
	if sf == nil {
		return domain.NewScrapperFilters(nil, "", domain.OrderFilterPopular)
	}

	order := convertOrderFilter(sf.Order)
	return domain.NewScrapperFilters(sf.Genres, sf.Platform, order)
}

func convertOrderFilter(o pb.OrderFilter) domain.OrderFilter {
	switch o {
	case pb.OrderFilter_ORDER_HIGHEST_RATED:
		return domain.OrderFilterHighest
	case pb.OrderFilter_ORDER_NEWEST:
		return domain.OrderFilterNewest
	case pb.OrderFilter_ORDER_SHORTEST:
		return domain.OrderFilterShortest
	default:
		return domain.OrderFilterPopular
	}
}

func convertFilmsToProto(films []*domain.Film) []*pb.Film {
	protoFilms := make([]*pb.Film, len(films))
	for i, film := range films {
		protoFilms[i] = convertFilmToProto(film)
	}
	return protoFilms
}

func convertFilmToProto(film *domain.Film) *pb.Film {
	return &pb.Film{
		Title:     film.Title,
		Duration:  int32(film.Duration),
		Directors: film.Directors,
		Year:      int32(film.Year),
	}
}
