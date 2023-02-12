// eth protocol get get block headers https://github.com/ethereum/devp2p/blob/master/caps/eth.md#getblockreceipts-0x05
package mimicry

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/sirupsen/logrus"
)

const (
	GetReceiptsCode = 0x1f
)

type GetReceipts eth.GetReceiptsPacket66

func (msg *GetReceipts) Code() int { return GetReceiptsCode }

func (msg *GetReceipts) ReqID() uint64 { return msg.RequestId }

func (c *Client) receiveGetReceipts(ctx context.Context, data []byte) (*GetReceipts, error) {
	s := new(GetReceipts)
	if err := rlp.DecodeBytes(data, &s); err != nil {
		return nil, fmt.Errorf("error decoding get block receipts: %w", err)
	}

	return s, nil
}

func (c *Client) handleGetReceipts(ctx context.Context, code uint64, data []byte) error {
	c.log.WithField("code", code).Debug("received GetReceipts")

	blockBodies, err := c.receiveGetReceipts(ctx, data)
	if err != nil {
		return err
	}

	err = c.sendReceipts(ctx, &Receipts{
		RequestId:      blockBodies.RequestId,
		ReceiptsPacket: [][]*types.Receipt{},
	})
	if err != nil {
		c.handleSessionError(ctx, err)
		return err
	}

	return nil
}

func (c *Client) sendGetReceipts(ctx context.Context, bh *GetReceipts) error {
	c.log.WithFields(logrus.Fields{
		"code":       GetReceiptsCode,
		"request_id": bh.RequestId,
		"receipts":   bh.GetReceiptsPacket,
	}).Debug("sending GetReceipts")

	encodedData, err := rlp.EncodeToBytes(bh)
	if err != nil {
		return fmt.Errorf("error encoding get block receipts: %w", err)
	}

	if _, err := c.rlpxConn.Write(GetReceiptsCode, encodedData); err != nil {
		return fmt.Errorf("error sending get block receipts: %w", err)
	}

	return nil
}
