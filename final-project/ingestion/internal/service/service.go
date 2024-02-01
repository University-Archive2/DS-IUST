package service

import (
	"context"
	"encoding/json"
	"ingestion/pkg/broker"
)

type IngestionService struct {
	producer broker.Producer
}

func NewIngestionService(producer broker.Producer) *IngestionService {
	return &IngestionService{
		producer: producer,
	}
}

func (s *IngestionService) ProduceFinancialData(dataType string, data any) error {
	bytesData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return s.producer.Produce(context.Background(), dataType, bytesData)
}
