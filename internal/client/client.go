package client

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/defryheryanto/mini-wallet/internal/wallet"
)

type Client struct {
	Xid   string `json:"xid"`
	Token string `json:"token"`
}

type ClientRepository interface {
	Insert(ctx context.Context, data *Client) error
	FindByXid(ctx context.Context, xid string) (*Client, error)
	FindByToken(ctx context.Context, token string) (*Client, error)
}

type ClientIService interface {
	Create(ctx context.Context, xid string) (*Client, error)
	GetByToken(ctx context.Context, token string) (*Client, error)
}

type ClientService struct {
	repository    ClientRepository
	walletService wallet.WalletIService
}

func NewClientService(repository ClientRepository, walletService wallet.WalletIService) ClientIService {
	return &ClientService{repository, walletService}
}

func (s *ClientService) Create(ctx context.Context, xid string) (*Client, error) {
	existingClient, err := s.repository.FindByXid(ctx, xid)
	if err != nil {
		return nil, err
	}
	if existingClient != nil {
		return nil, ErrXidAlreadyTaken
	}

	var token string
	for {
		token = s.generateToken()
		clientByToken, err := s.repository.FindByToken(ctx, token)
		if err != nil {
			return nil, err
		}
		if clientByToken == nil {
			break
		}
	}

	// TODO: Wrap this with DB transaction
	err = s.repository.Insert(ctx, &Client{
		Xid:   xid,
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	err = s.walletService.Create(ctx, &wallet.CreateWalletParams{
		OwnedBy: xid,
	})
	if err != nil {
		return nil, err
	}

	existingClient, err = s.repository.FindByXid(ctx, xid)
	if err != nil {
		return nil, err
	}

	return existingClient, nil
}

func (s *ClientService) GetByToken(ctx context.Context, token string) (*Client, error) {
	currentClient, err := s.repository.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return currentClient, nil
}

func (s *ClientService) generateToken() string {
	b := make([]byte, 10)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
