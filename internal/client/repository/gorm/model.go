package gorm

import "github.com/defryheryanto/mini-wallet/internal/client"

type Client struct {
	Xid   string `gorm:"primaryKey;column:xid"`
	Token string `gorm:"column:token"`
}

func (Client) TableName() string {
	return "clients"
}

func (Client) FromServiceModel(data *client.Client) *Client {
	if data == nil {
		return nil
	}

	return &Client{
		Xid:   data.Xid,
		Token: data.Token,
	}
}

func (c *Client) ToServiceModel() *client.Client {
	return &client.Client{
		Xid:   c.Xid,
		Token: c.Token,
	}
}
