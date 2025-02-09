package nft

import (
	"crypto/sha256"
	"fmt"

	"github.com/xssnick/tonutils-go/tvm/cell"
)

type ContentAny interface {
	ContentCell() (*cell.Cell, error)
}

type ContentOffchain struct {
	URI string
}

type ContentOnchain struct {
	Name        string
	Description string
	Image       string
	ImageData   []byte
	attributes  *cell.Dictionary
}

type ContentSemichain struct {
	ContentOffchain
	ContentOnchain
}

func ContentFromCell(c *cell.Cell) (ContentAny, error) {
	s := c.BeginParse()
	if s.BitsLeft() < 8 {
		return nil, nil
	}

	typ, err := s.LoadUInt(8)
	if err != nil {
		return nil, fmt.Errorf("failed to load type: %w", err)
	}
	t := uint8(typ)

	switch t {
	case 0x00:
		dict, err := s.LoadDict(256)
		if err != nil {
			return nil, fmt.Errorf("failed to load dict onchain data: %w", err)
		}

		return &ContentOnchain{
			Name:        string(getOnchainVal(dict, "name")),
			Description: string(getOnchainVal(dict, "description")),
			Image:       string(getOnchainVal(dict, "image")),
			ImageData:   getOnchainVal(dict, "image_data"),
			attributes:  dict,
		}, err
	case 0x01:
		str, err := s.LoadStringSnake()
		if err != nil {
			return nil, fmt.Errorf("failed to load snake offchain data: %w", err)
		}

		return &ContentOffchain{
			URI: str,
		}, err
	default:
		str, err := s.LoadStringSnake()
		if err != nil {
			return nil, fmt.Errorf("failed to load snake offchain data: %w", err)
		}

		return &ContentOffchain{
			URI: string(t) + str,
		}, err
	}
}

func getOnchainVal(dict *cell.Dictionary, key string) []byte {
	h := sha256.New()
	h.Write([]byte(key))

	val := dict.Get(cell.BeginCell().MustStoreSlice(h.Sum(nil), 256).EndCell())
	if val != nil {
		v := val.BeginParse()
		typ, err := v.LoadUInt(8)
		if err != nil {
			return nil
		}

		switch typ {
		case 0x01:
			data, _ := v.LoadBinarySnake()
			return data
		default:
			// TODO: add support for chunked
			return nil
		}
	}

	return nil
}

func setOnchainVal(dict *cell.Dictionary, key string, val []byte) error {
	h := sha256.New()
	h.Write([]byte(key))

	v := cell.BeginCell().MustStoreUInt(0x00, 8)
	if err := v.StoreBinarySnake(val); err != nil {
		return err
	}

	err := dict.Set(cell.BeginCell().MustStoreSlice(h.Sum(nil), 256).EndCell(), v.EndCell())
	if err != nil {
		return err
	}

	return nil
}

func (c *ContentOffchain) ContentCell() (*cell.Cell, error) {
	// https://github.com/ton-blockchain/TIPs/issues/64
	// Standard says that prefix should be 0x01, but looks like it was misunderstanding in other implementations and 0x01 was dropped
	// so we make compatibility
	return cell.BeginCell().MustStoreStringSnake(c.URI).EndCell(), nil
}

func (c *ContentSemichain) ContentCell() (*cell.Cell, error) {
	if c.attributes == nil {
		c.attributes = cell.NewDict(256)
	}

	if c.URI != "" && getOnchainVal(c.attributes, "uri") == nil {
		ci := cell.BeginCell()

		err := ci.StoreStringSnake(c.URI)
		if err != nil {
			return nil, err
		}

		err = setOnchainVal(c.attributes, "uri", []byte(c.URI))
		if err != nil {
			return nil, err
		}
	}

	return c.ContentOnchain.ContentCell()
}

func (c *ContentOnchain) ContentCell() (*cell.Cell, error) {
	if c.attributes == nil {
		c.attributes = cell.NewDict(256)
	}

	if len(c.Image) > 0 {
		err := setOnchainVal(c.attributes, "image", []byte(c.Image))
		if err != nil {
			return nil, fmt.Errorf("failed to store image: %w", err)
		}
	}
	if len(c.ImageData) > 0 {
		err := setOnchainVal(c.attributes, "image_data", c.ImageData)
		if err != nil {
			return nil, fmt.Errorf("failed to store image_data: %w", err)
		}
	}
	if len(c.Name) > 0 {
		err := setOnchainVal(c.attributes, "name", []byte(c.Name))
		if err != nil {
			return nil, fmt.Errorf("failed to store name: %w", err)
		}
	}
	if len(c.Description) > 0 {
		err := setOnchainVal(c.attributes, "description", []byte(c.Description))
		if err != nil {
			return nil, fmt.Errorf("failed to store description: %w", err)
		}
	}

	return cell.BeginCell().MustStoreUInt(0x00, 8).MustStoreDict(c.attributes).EndCell(), nil
}
