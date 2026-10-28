package initialization

import (
	"time"

	aipb "github.com/InstaySystem/is_v2-be/proto/gen/ai"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type ClientAddresses struct {
	AIAddr string
}

type GRPCClients struct {
	AIClient aipb.AIServiceClient
	aiConn   *grpc.ClientConn
}

func InitClients(ca ClientAddresses) (*GRPCClients, error) {
	opts := dialOptions()
	aiConn, err := grpc.NewClient(ca.AIAddr, opts...)
	if err != nil {
		return nil, err
	}
	aiClient := aipb.NewAIServiceClient(aiConn)

	return &GRPCClients{
		aiClient,
		aiConn,
	}, nil
}

func (g *GRPCClients) Close() {
	_ = g.aiConn.Close()
}

func dialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{
			"methodConfig": [{
				"name": [{}],
				"retryPolicy": {
					"MaxAttempts": 4,
					"InitialBackoff": "0.1s",
					"MaxBackoff": "1s", 
					"BackoffMultiplier": 2.0,
					"RetryableStatusCodes": ["UNAVAILABLE", "DEADLINE_EXCEEDED"]
				}
			}]
		}`),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                5 * time.Minute,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	}
}
