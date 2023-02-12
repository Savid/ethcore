// eth protocol block headers https://github.com/ethereum/devp2p/blob/master/caps/eth.md#blockheaders-0x04
package mimicry

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/sirupsen/logrus"
)

const (
	BlockHeadersCode = 0x14
)

type BlockHeaders eth.BlockHeadersPacket66

func (msg *BlockHeaders) Code() int { return BlockHeadersCode }

func (msg *BlockHeaders) ReqID() uint64 { return msg.RequestId }

func (m *Mimicry) receiveBlockHeaders(ctx context.Context, data []byte) (*BlockHeaders, error) {
	s := new(BlockHeaders)
	if err := rlp.DecodeBytes(data, &s); err != nil {
		return nil, fmt.Errorf("error decoding block headers: %w", err)
	}

	return s, nil
}

func (m *Mimicry) handleBlockHeaders(ctx context.Context, code uint64, data []byte) error {
	m.log.WithField("code", code).Debug("received BlockHeaders")

	blockHeaders, err := m.receiveBlockHeaders(ctx, data)
	if err != nil {
		return err
	}

	err = m.sendBlockHeaders(ctx, blockHeaders)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mimicry) sendBlockHeaders(ctx context.Context, bh *BlockHeaders) error {
	m.log.WithFields(logrus.Fields{
		"code":          BlockHeadersCode,
		"request_id":    bh.RequestId,
		"headers_count": len(bh.BlockHeadersPacket),
	}).Debug("sending BlockHeaders")

	encodedData, err := rlp.EncodeToBytes(bh)
	if err != nil {
		return fmt.Errorf("error encoding block headers: %w", err)
	}

	if _, err := m.rlpxConn.Write(BlockHeadersCode, encodedData); err != nil {
		return fmt.Errorf("error sending block headers: %w", err)
	}

	return nil
}
