package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "net"
    "sync"

    "google.golang.org/grpc"
    pb "github.com/DeathStroke19891/transactions/transactions"
)

var (
    port = flag.Int("port", 50051, "The server port")
)

type server struct {
    pb.UnimplementedTransactionsServer
    mu           sync.Mutex
    transactions []*pb.Transaction
}

func (s *server) CommitTransaction(ctx context.Context, in *pb.Transaction) (*pb.Status, error) {
    log.Printf("Received Transaction: %v", in.GetId())

    s.mu.Lock()
    defer s.mu.Unlock()

    s.transactions = append(s.transactions, in)
	log.Printf("Size of minePool: %d", len(s.transactions))
    return &pb.Status{
        Res:     1,
        Message: "Transaction committed successfully",
    }, nil
}

func main() {
    flag.Parse()
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()
    pb.RegisterTransactionsServer(s, &server{})

    log.Printf("server listening at %v", lis.Addr())
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
