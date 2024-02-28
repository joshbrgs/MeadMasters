package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/joshbrgs/meadmasters/products/reviews/protobuf/reviews"
)

var reviewCollection *mongo.Collection

type server struct {
	pb.UnimplementedReviewServiceServer
}

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("reviewService")
	reviewCollection = db.Collection("reviews")
}

func (s *server) GetReviewById(ctx context.Context, req *pb.ReviewByIdRequest) (*pb.Review, error) {
	reviewID := req.GetId()
	result := reviewCollection.FindOne(context.Background(), bson.M{"_id": reviewID})

	review := &pb.Review{}
	if err := result.Decode(review); err != nil {
		return nil, fmt.Errorf("review not found: %v", err)
	}

	return review, nil
}

func (s *server) CreateReview(
	ctx context.Context,
	req *pb.CreateReviewRequest,
) (*pb.ReviewIdResponse, error) {
	review := &pb.Review{
		Name:     req.GetName(),
		Location: req.GetLocation(),
	}

	insertResult, err := reviewCollection.InsertOne(context.Background(), review)
	if err != nil {
		return nil, fmt.Errorf("Failed to create review: %v", err)
	}

	return &pb.ReviewIdResponse{Id: insertResult.InsertedID.(string)}, nil
}

func (s *server) UpdateReview(
	ctx context.Context,
	req *pb.UpdateReviewRequest,
) (*pb.ReviewResponse, error) {
	updateResult, err := reviewCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": req.GetId()},
		bson.D{{"$set", bson.D{{"name", req.GetName()}, {"location", req.GetLocation()}}}},
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to update review: %v", err)
	}

	if updateResult.ModifiedCount == 0 {
		return nil, fmt.Errorf("review not found")
	}

	return &pb.ReviewResponse{Message: "review updated successfully"}, nil
}

func (s *server) DeleteReview(
	ctx context.Context,
	req *pb.ReviewByIdRequest,
) (*pb.ReviewResponse, error) {
	deleteResult, err := reviewCollection.DeleteOne(
		context.Background(),
		bson.M{"_id": req.GetId()},
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to delete review: %v", err)
	}

	if deleteResult.DeletedCount == 0 {
		return nil, fmt.Errorf("review not found")
	}

	return &pb.ReviewResponse{Message: "review deleted successfully"}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterReviewServiceServer(s, &server{})
	reflection.Register(s)

	log.Println("gRPC server is running on http://localhost:50051")
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
