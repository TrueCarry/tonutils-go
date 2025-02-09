package nft

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/tvm/cell"
)

type ItemEditPayload struct {
	_       tlb.Magic  `tlb:"#1a0b9d51"`
	QueryID uint64     `tlb:"## 64"`
	Content *cell.Cell `tlb:"^"`
}

type ItemEditableClient struct {
	*ItemClient
}

func NewItemEditableClient(api *ton.APIClient, nftAddr *address.Address) *ItemEditableClient {
	return &ItemEditableClient{
		ItemClient: NewItemClient(api, nftAddr),
	}
}

func (c *ItemEditableClient) GetEditor(ctx context.Context) (*address.Address, error) {
	b, err := c.api.CurrentMasterchainInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get masterchain info: %w", err)
	}

	res, err := c.api.RunGetMethod(ctx, b, c.addr, "get_editor")
	if err != nil {
		return nil, fmt.Errorf("failed to run get_editor method: %w", err)
	}

	x, ok := res[0].(*cell.Slice)
	if !ok {
		return nil, fmt.Errorf("result is not slice")
	}

	addr, err := x.LoadAddr()
	if err != nil {
		return nil, fmt.Errorf("failed to load address from result slice: %w", err)
	}

	return addr, nil
}

func (c *ItemEditableClient) BuildEditPayload(content ContentAny) (*cell.Cell, error) {
	con, err := content.ContentCell()
	if err != nil {
		return nil, err
	}

	body, err := tlb.ToCell(ItemEditPayload{
		QueryID: rand.Uint64(),
		Content: con,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to convert ItemEditPayload to cell: %w", err)
	}

	return body, nil
}
