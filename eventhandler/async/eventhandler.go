// Copyright (c) 2017 - Max Ekman <max@looplab.se>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package async

import (
	"context"
	"fmt"

	eh "github.com/looplab/eventhorizon"
)

// EventHandler is an async event handler middleware. It will run the event in
// a new go routine and report any errors to the error channel obtained by Errors().
type EventHandler struct {
	eh.EventHandler
	errCh chan Error
}

// NewEventHandler creates a new EventHandler.
func NewEventHandler(h eh.EventHandler) *EventHandler {
	return &EventHandler{
		EventHandler: h,
		errCh:        make(chan Error, 20),
	}
}

// HandleEvent implements the HandleEvent method of the EventHandler interface.
func (h *EventHandler) HandleEvent(ctx context.Context, event eh.Event) error {
	go func() {
		if err := h.EventHandler.HandleEvent(ctx, event); err != nil {
			// Always try to deliver errors.
			h.errCh <- Error{err, ctx, event}
		}
	}()
	return nil
}

// Errors returns an error channel where async handling errors are sent.
func (h *EventHandler) Errors() <-chan Error {
	return h.errCh
}

// Error is an async error containing the error and the event.
type Error struct {
	Err   error
	Ctx   context.Context
	Event eh.Event
}

// Error implements the Error method of the error interface.
func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Event.String(), e.Err.Error())
}
