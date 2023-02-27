package seller

import (
	"KeySell/pkg/auth"
	"context"
)

type SellerStore interface {
	UpdateData(ctx context.Context, m map[string]interface{}, SellerID int) error
	SignUp(ctx context.Context, m map[string]interface{}) (int, error)
	SignIn(ctx context.Context, m map[string]interface{}) (int, error)
	GetInfo(ctx context.Context, UserID int) (map[string]interface{}, error)
}

type SellerService struct {
	SellerStore SellerStore
	tm          auth.TokenManager
}

func NewSellerService(sellerStore SellerStore, tm auth.TokenManager) *SellerService {
	return &SellerService{SellerStore: sellerStore, tm: tm}
}

func (s *SellerService) SignUp(ctx context.Context, m map[string]interface{}) (int, error) {
	return s.SellerStore.SignUp(ctx, m)
}

func (s *SellerService) SignIn(ctx context.Context, m map[string]interface{}) (string, error) {
	in, err := s.SellerStore.SignIn(ctx, m)
	if err != nil {
		return "", err
	}
	token, err := s.tm.GenerateToken(in)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *SellerService) UpdateData(ctx context.Context, m map[string]interface{}, UserID int) error {
	return s.SellerStore.UpdateData(ctx, m, UserID)
}

func (s *SellerService) GetInfo(ctx context.Context, UserID int) (map[string]interface{}, error) {
	return s.SellerStore.GetInfo(ctx, UserID)
}
