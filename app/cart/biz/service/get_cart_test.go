package service

import (
	"context"
	cart "github.com/jazwu/tiktok-ecommerce/app/cart/app/cart/kitex_gen/cart"
	"testing"
)

func TestGetCart_Run(t *testing.T) {
	ctx := context.Background()
	s := NewGetCartService(ctx)
	// init req and assert value

	req := &cart.GetCartReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
