package actor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nspcc-dev/neo-go/pkg/core/block"
	"github.com/nspcc-dev/neo-go/pkg/core/state"
	"github.com/nspcc-dev/neo-go/pkg/neorpc"
	"github.com/nspcc-dev/neo-go/pkg/neorpc/result"
	"github.com/nspcc-dev/neo-go/pkg/smartcontract/trigger"
	"github.com/nspcc-dev/neo-go/pkg/util"
)

// PollingWaiterRetryCount is a threshold for a number of subsequent failed
// attempts to get block count from the RPC server for PollingWaiter. If it fails
// to retrieve block count PollingWaiterRetryCount times in a raw then transaction
// awaiting attempt considered to be failed and an error is returned.
const PollingWaiterRetryCount = 3

var (
	// ErrTxNotAccepted is returned when transaction wasn't accepted to the chain
	// even after ValidUntilBlock block persist.
	ErrTxNotAccepted = errors.New("transaction was not accepted to chain")
	// ErrContextDone is returned when Waiter context has been done in the middle
	// of transaction awaiting process and no result was received yet.
	ErrContextDone = errors.New("waiter context done")
	// ErrAwaitingNotSupported is returned from Wait method if Waiter instance
	// doesn't support transaction awaiting.
	ErrAwaitingNotSupported = errors.New("awaiting not supported")
	// ErrMissedEvent is returned when RPCEventWaiter closes receiver channel
	// which happens if missed event was received from the RPC server.
	ErrMissedEvent = errors.New("some event was missed")
)

type (
	// Waiter is an interface providing transaction awaiting functionality to Actor.
	Waiter interface {
		// Wait allows to wait until transaction will be accepted to the chain. It can be
		// used as a wrapper for Send or SignAndSend and accepts transaction hash,
		// ValidUntilBlock value and an error. It returns transaction execution result
		// or an error if transaction wasn't accepted to the chain.
		Wait(h util.Uint256, vub uint32, err error) (*state.AppExecResult, error)
		// WaitAny waits until at least one of the specified transactions will be accepted
		// to the chain until vub (including). It returns execution result of this
		// transaction or an error if none of the transactions was accepted to the chain.
		// It uses underlying RPCPollingWaiter or RPCEventWaiter context to interrupt
		// awaiting process, but additional ctx can be passed as an argument for the same
		// purpose.
		WaitAny(ctx context.Context, vub uint32, hashes ...util.Uint256) (*state.AppExecResult, error)
	}
	// RPCPollingWaiter is an interface that enables transaction awaiting functionality
	// for Actor instance based on periodical BlockCount and ApplicationLog polls.
	RPCPollingWaiter interface {
		// Context should return the RPC client context to be able to gracefully
		// shut down all running processes (if so).
		Context() context.Context
		GetVersion() (*result.Version, error)
		GetBlockCount() (uint32, error)
		GetApplicationLog(hash util.Uint256, trig *trigger.Type) (*result.ApplicationLog, error)
	}
	// RPCEventWaiter is an interface that enables improved transaction awaiting functionality
	// for Actor instance based on web-socket Block and ApplicationLog notifications. RPCEventWaiter
	// contains RPCPollingWaiter under the hood and falls back to polling when subscription-based
	// awaiting fails.
	RPCEventWaiter interface {
		RPCPollingWaiter

		ReceiveBlocks(flt *neorpc.BlockFilter, rcvr chan<- *block.Block) (string, error)
		ReceiveExecutions(flt *neorpc.ExecutionFilter, rcvr chan<- *state.AppExecResult) (string, error)
		Unsubscribe(id string) error
	}
)

// NullWaiter is a Waiter stub that doesn't support transaction awaiting functionality.
type NullWaiter struct{}

// PollingWaiter is a polling-based Waiter.
type PollingWaiter struct {
	polling RPCPollingWaiter
	version *result.Version
}

// EventWaiter is a websocket-based Waiter.
type EventWaiter struct {
	ws      RPCEventWaiter
	polling Waiter
}

// newWaiter creates Waiter instance. It can be either websocket-based or
// polling-base, otherwise Waiter stub is returned.
func newWaiter(ra RPCActor, v *result.Version) Waiter {
	if eventW, ok := ra.(RPCEventWaiter); ok {
		return &EventWaiter{
			ws: eventW,
			polling: &PollingWaiter{
				polling: eventW,
				version: v,
			},
		}
	}
	if pollW, ok := ra.(RPCPollingWaiter); ok {
		return &PollingWaiter{
			polling: pollW,
			version: v,
		}
	}
	return NewNullWaiter()
}

// NewNullWaiter creates an instance of Waiter stub.
func NewNullWaiter() NullWaiter {
	return NullWaiter{}
}

// Wait implements Waiter interface.
func (NullWaiter) Wait(h util.Uint256, vub uint32, err error) (*state.AppExecResult, error) {
	return nil, ErrAwaitingNotSupported
}

// WaitAny implements Waiter interface.
func (NullWaiter) WaitAny(ctx context.Context, vub uint32, hashes ...util.Uint256) (*state.AppExecResult, error) {
	return nil, ErrAwaitingNotSupported
}

// NewPollingWaiter creates an instance of Waiter supporting poll-based transaction awaiting.
func NewPollingWaiter(waiter RPCPollingWaiter) (*PollingWaiter, error) {
	v, err := waiter.GetVersion()
	if err != nil {
		return nil, err
	}
	return &PollingWaiter{
		polling: waiter,
		version: v,
	}, nil
}

// Wait implements Waiter interface.
func (w *PollingWaiter) Wait(h util.Uint256, vub uint32, err error) (*state.AppExecResult, error) {
	if err != nil {
		return nil, err
	}
	return w.WaitAny(context.TODO(), vub, h)
}

// WaitAny implements Waiter interface.
func (w *PollingWaiter) WaitAny(ctx context.Context, vub uint32, hashes ...util.Uint256) (*state.AppExecResult, error) {
	var (
		currentHeight uint32
		failedAttempt int
		pollTime      = time.Millisecond * time.Duration(w.version.Protocol.MillisecondsPerBlock) / 2
	)
	if pollTime == 0 {
		pollTime = time.Second
	}
	timer := time.NewTicker(pollTime)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			blockCount, err := w.polling.GetBlockCount()
			if err != nil {
				failedAttempt++
				if failedAttempt > PollingWaiterRetryCount {
					return nil, fmt.Errorf("failed to retrieve block count: %w", err)
				}
				continue
			}
			failedAttempt = 0
			if blockCount-1 > currentHeight {
				currentHeight = blockCount - 1
			}
			t := trigger.Application
			for _, h := range hashes {
				res, err := w.polling.GetApplicationLog(h, &t)
				if err == nil {
					return &state.AppExecResult{
						Container: res.Container,
						Execution: res.Executions[0],
					}, nil
				}
			}
			if currentHeight >= vub {
				return nil, ErrTxNotAccepted
			}
		case <-w.polling.Context().Done():
			return nil, fmt.Errorf("%w: %v", ErrContextDone, w.polling.Context().Err())
		case <-ctx.Done():
			return nil, fmt.Errorf("%w: %v", ErrContextDone, ctx.Err())
		}
	}
}

// NewEventWaiter creates an instance of Waiter supporting websocket event-based transaction awaiting.
// EventWaiter contains PollingWaiter under the hood and falls back to polling when subscription-based
// awaiting fails.
func NewEventWaiter(waiter RPCEventWaiter) (*EventWaiter, error) {
	polling, err := NewPollingWaiter(waiter)
	if err != nil {
		return nil, err
	}
	return &EventWaiter{
		ws:      waiter,
		polling: polling,
	}, nil
}

// Wait implements Waiter interface.
func (w *EventWaiter) Wait(h util.Uint256, vub uint32, err error) (res *state.AppExecResult, waitErr error) {
	if err != nil {
		return nil, err
	}
	return w.WaitAny(context.TODO(), vub, h)
}

// WaitAny implements Waiter interface.
func (w *EventWaiter) WaitAny(ctx context.Context, vub uint32, hashes ...util.Uint256) (res *state.AppExecResult, waitErr error) {
	var wsWaitErr error
	defer func() {
		if wsWaitErr != nil {
			res, waitErr = w.polling.WaitAny(ctx, vub, hashes...)
			if waitErr != nil {
				waitErr = fmt.Errorf("WS waiter error: %w, simple waiter error: %v", wsWaitErr, waitErr)
			}
		}
	}()
	bRcvr := make(chan *block.Block)
	aerRcvr := make(chan *state.AppExecResult)
	defer func() {
	drainLoop:
		// Drain receivers to avoid other notification receivers blocking.
		for {
			select {
			case <-bRcvr:
			case <-aerRcvr:
			default:
				break drainLoop
			}
		}
		if wsWaitErr == nil || !errors.Is(wsWaitErr, ErrMissedEvent) {
			close(bRcvr)
			close(aerRcvr)
		}
	}()
	// Execution event precedes the block event, thus wait until the VUB-th block to be sure.
	since := vub
	blocksID, err := w.ws.ReceiveBlocks(&neorpc.BlockFilter{Since: &since}, bRcvr)
	if err != nil {
		wsWaitErr = fmt.Errorf("failed to subscribe for new blocks: %w", err)
		return
	}
	defer func() {
		err = w.ws.Unsubscribe(blocksID)
		if err != nil {
			errFmt := "failed to unsubscribe from blocks (id: %s): %v"
			errArgs := []interface{}{blocksID, err}
			if waitErr != nil {
				errFmt += "; wait error: %w"
				errArgs = append(errArgs, waitErr)
			}
			waitErr = fmt.Errorf(errFmt, errArgs...)
		}
	}()
	for _, h := range hashes {
		txsID, err := w.ws.ReceiveExecutions(&neorpc.ExecutionFilter{Container: &h}, aerRcvr)
		if err != nil {
			wsWaitErr = fmt.Errorf("failed to subscribe for execution results: %w", err)
			return
		}
		defer func() {
			err = w.ws.Unsubscribe(txsID)
			if err != nil {
				errFmt := "failed to unsubscribe from transactions (id: %s): %v"
				errArgs := []interface{}{txsID, err}
				if waitErr != nil {
					errFmt += "; wait error: %w"
					errArgs = append(errArgs, waitErr)
				}
				waitErr = fmt.Errorf(errFmt, errArgs...)
			}
		}()
	}

	select {
	case _, ok := <-bRcvr:
		if !ok {
			// We're toast, retry with non-ws client.
			wsWaitErr = ErrMissedEvent
			return
		}
		waitErr = ErrTxNotAccepted
	case aer, ok := <-aerRcvr:
		if !ok {
			// We're toast, retry with non-ws client.
			wsWaitErr = ErrMissedEvent
			return
		}
		res = aer
	case <-w.ws.Context().Done():
		waitErr = fmt.Errorf("%w: %v", ErrContextDone, w.ws.Context().Err())
	case <-ctx.Done():
		waitErr = fmt.Errorf("%w: %v", ErrContextDone, ctx.Err())
	}
	return
}
