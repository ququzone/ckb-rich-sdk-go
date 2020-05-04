package indexer

import (
	"context"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/rpc"
)

type Client interface {
	// GetCells returns the live cells collection by the lock or type script.
	GetCells(ctx context.Context, searchKey *SearchKey, order SearchOrder, limit uint64, afterCursor string) (*LiveCells, error)

	// GetTransactions returns the transactions collection by the lock or type script.
	GetTransactions(ctx context.Context, searchKey *SearchKey, order SearchOrder, limit uint64, afterCursor string) (*Transactions, error)

	// Close close client
	Close()
}

type client struct {
	c *rpc.Client
}

func Dial(url string) (Client, error) {
	return DialContext(context.Background(), url)
}

func DialContext(ctx context.Context, url string) (Client, error) {
	c, err := rpc.DialContext(ctx, url)
	if err != nil {
		return nil, err
	}
	return NewClient(c), nil
}

func NewClient(c *rpc.Client) Client {
	return &client{c}
}

func (cli *client) Close() {
	cli.c.Close()
}

func (cli *client) GetCells(ctx context.Context, searchKey *SearchKey, order SearchOrder, limit uint64, afterCursor string) (*LiveCells, error) {
	var result liveCells
	var err error
	if afterCursor == "" {
		err = cli.c.CallContext(ctx, &result, "get_cells", fromSearchKey(searchKey), order, hexutil.Uint64(limit))
	} else {
		err = cli.c.CallContext(ctx, &result, "get_cells", fromSearchKey(searchKey), order, hexutil.Uint64(limit), afterCursor)
	}
	if err != nil {
		return nil, err
	}
	return toLiveCells(result), err
}

func (cli *client) GetTransactions(ctx context.Context, searchKey *SearchKey, order SearchOrder, limit uint64, afterCursor string) (*Transactions, error) {
	var result transactions
	var err error
	if afterCursor == "" {
		err = cli.c.CallContext(ctx, &result, "get_transactions", fromSearchKey(searchKey), order, hexutil.Uint64(limit))
	} else {
		err = cli.c.CallContext(ctx, &result, "get_transactions", fromSearchKey(searchKey), order, hexutil.Uint64(limit), afterCursor)
	}
	if err != nil {
		return nil, err
	}
	return toTransactions(result), err
}
