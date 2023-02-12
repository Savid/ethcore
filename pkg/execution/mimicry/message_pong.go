// RLPx pong https://github.com/ethereum/devp2p/blob/master/rlpx.md#pong-0x03
package mimicry

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

const (
	PongCode = 0x03
)

type Pong struct{}

func (h *Pong) Code() int { return PongCode }

func (h *Pong) ReqID() uint64 { return 0 }

func (m *Mimicry) sendPong(ctx context.Context) error {
	m.log.WithFields(logrus.Fields{
		"code": PongCode,
	}).Debug("sending Pong")

	if _, err := m.rlpxConn.Write(PongCode, []byte{}); err != nil {
		return fmt.Errorf("error sending pong: %w", err)
	}

	return nil
}
