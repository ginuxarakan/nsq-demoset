package service

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"nsq-demoset/app/app-services/model"
	marketpb "nsq-demoset/app/app-services/proto/market/v1"
	"sync"
)

type MarketService struct {
	mu       sync.RWMutex
	client   marketpb.MarketClient
	CoinChan chan *model.CoinData
}

func NewMarketService(addr string) *MarketService {
	svc := &MarketService{}
	svc.CoinChan = make(chan *model.CoinData)

	go svc.getData(addr)

	return svc
}

func (s *MarketService) getData(addr string) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := marketpb.NewMarketClient(conn)
	s.client = client

	stream, err := client.Subscribe(context.Background(), &marketpb.MarketRequest{
		Coin: marketpb.Coin_ALL,
	})
	if err != nil {
		panic(err)
	}

	for {
		value, err := stream.Recv()
		if err != nil {
			panic(err)
		}
		if value.GetSymbol() != "" {
			coin := &model.CoinData{
				Symbol:  value.GetSymbol(),
				Price:   value.GetLastPrice(),
				Percent: value.GetPriceChangePercent(),
			}
			s.CoinChan <- coin
		}
	}
}

func (s *MarketService) GetClient() marketpb.MarketClient {
	return s.client
}
