package target

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/target/request"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

type callerRequest interface {
	OnMessage(ctx context.Context, msg *types.MessageBody) error
	ResponseChan() <-chan commoncap.CapabilityResponse
	Expired() bool
	Cancel(reason string)
}

// caller/Receiver are shims translating between capability API calls and network messages
type caller struct {
	lggr                 logger.Logger
	remoteCapabilityInfo commoncap.CapabilityInfo
	localDONInfo         capabilities.DON
	dispatcher           types.Dispatcher
	requestTimeout       time.Duration

	messageIDToExecuteRequest map[string]callerRequest
	mutex                     sync.Mutex
}

var _ commoncap.TargetCapability = &caller{}
var _ types.Receiver = &caller{}

func NewCaller(ctx context.Context, lggr logger.Logger, remoteCapabilityInfo commoncap.CapabilityInfo, localDonInfo capabilities.DON, dispatcher types.Dispatcher,
	requestTimeout time.Duration) *caller {

	c := &caller{
		lggr:                      lggr,
		remoteCapabilityInfo:      remoteCapabilityInfo,
		localDONInfo:              localDonInfo,
		dispatcher:                dispatcher,
		requestTimeout:            requestTimeout,
		messageIDToExecuteRequest: make(map[string]callerRequest),
	}

	go func() {
		ticker := time.NewTicker(requestTimeout)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.ExpireRequests()
			}
		}
	}()

	return c
}

func (c *caller) ExpireRequests() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for messageID, req := range c.messageIDToExecuteRequest {
		if req.Expired() {
			req.Cancel("request expired")
			delete(c.messageIDToExecuteRequest, messageID)
		}
	}
}

func (c *caller) Info(ctx context.Context) (commoncap.CapabilityInfo, error) {
	return c.remoteCapabilityInfo, nil
}

func (c *caller) RegisterToWorkflow(ctx context.Context, request commoncap.RegisterToWorkflowRequest) error {
	return errors.New("not implemented")
}

func (c *caller) UnregisterFromWorkflow(ctx context.Context, request commoncap.UnregisterFromWorkflowRequest) error {
	return errors.New("not implemented")
}

func (c *caller) Execute(ctx context.Context, req commoncap.CapabilityRequest) (<-chan commoncap.CapabilityResponse, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	messageID, err := GetMessageIDForRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get message ID for request: %w", err)
	}

	if _, ok := c.messageIDToExecuteRequest[messageID]; ok {
		return nil, fmt.Errorf("request for message ID %s already exists", messageID)
	}

	callerReq, err := request.NewCallerRequest(ctx, c.lggr, req, messageID, c.remoteCapabilityInfo, c.localDONInfo, c.dispatcher,
		c.requestTimeout)

	c.messageIDToExecuteRequest[messageID] = callerReq

	return callerReq.ResponseChan(), nil
}

func (c *caller) Receive(msg *types.MessageBody) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// TODO should the dispatcher be passing in a context?
	ctx := context.Background()

	messageID := GetMessageID(msg)

	req := c.messageIDToExecuteRequest[messageID]
	if req == nil {
		c.lggr.Warnw("received response for unknown message ID ", "messageID", messageID)
		return
	}

	go func() {
		if err := req.OnMessage(ctx, msg); err != nil {
			c.lggr.Errorw("failed to add response to request", "messageID", messageID, "err", err)
		}
	}()

}

func GetMessageIDForRequest(req commoncap.CapabilityRequest) (string, error) {
	if req.Metadata.WorkflowID == "" || req.Metadata.WorkflowExecutionID == "" {
		return "", errors.New("workflow ID and workflow execution ID must be set in request metadata")
	}

	return req.Metadata.WorkflowID + req.Metadata.WorkflowExecutionID, nil
}
