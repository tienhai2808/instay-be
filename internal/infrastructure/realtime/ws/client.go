package ws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/domain/constants"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	httpRes "github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/response"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/validator"
	"github.com/InstaySystem/is_v2-be/pkg/utils"
	validatorPkg "github.com/InstaySystem/is_v2-be/pkg/validator"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 2048

	eventMarkRead   = "mark_read"
	eventNewMessage = "new_message"
	eventTyping     = "typing"
	eventReaction   = "reaction"
	eventError      = "error"
)

var (
	newline = []byte{'\n'}
	// space   = []byte{' '}
)

type WSRequest struct {
	Event string          `json:"event" validate:"required,oneof=new_message mark_read reaction typing"`
	Data  json.RawMessage `json:"data" validate:"required"`
}

type WSResponse struct {
	Event   string `json:"event"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type WSClient struct {
	schema       string
	hub          *HubWS
	conn         *websocket.Conn
	send         chan []byte
	id           string
	clientID     int64
	departmentID *int64
	clientType   string
	ActiveChats  map[int64]bool
}

func NewWSClient(
	schema string,
	hub *HubWS,
	conn *websocket.Conn,
	clientID int64,
	clientType string,
	departmentID *int64,
) *WSClient {
	return &WSClient{
		schema,
		hub,
		conn,
		make(chan []byte, 256),
		utils.GenerateUUID(),
		clientID,
		departmentID,
		clientType,
		make(map[int64]bool),
	}
}

func (c *WSClient) ReadPump() {
	defer func() {
		c.hub.Unregister <- c
		if err := c.conn.Close(); err != nil {
			return
		}
	}()

	c.conn.SetReadLimit(maxMessageSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		return
	}
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var req WSRequest
		if err := json.Unmarshal(data, &req); err != nil {
			c.sendError(exceptions.ErrBadRequest.Code, exceptions.ErrBadRequest.Message, nil)
			break
		}

		if err := validator.Validate.Struct(req); err != nil {
			field, tag, param := validatorPkg.HandleRequestError(err)
			c.sendError(exceptions.ErrBadRequest.Code, exceptions.ErrBadRequest.Message, map[string]any{
				"field": field,
				"tag":   tag,
				"param": param,
			})
			break
		}

		switch req.Event {
		case eventNewMessage:
			c.handleSendMessage(req.Data)
		case eventMarkRead:
			c.handleMarkRead(req.Data)
		case eventReaction:
			c.handleReaction(req.Data)
		case eventTyping:
			c.handleTyping(req.Data)
		default:
			c.sendError(exceptions.ErrBadRequest.Code, exceptions.ErrBadRequest.Message, nil)
			return
		}
	}
}

func (c *WSClient) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := c.conn.Close(); err != nil {
			return
		}
	}()

	for {
		select {
		case message, ok := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}
			if !ok {
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					return
				}
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if _, err := w.Write(message); err != nil {
				return
			}

			n := len(c.send)
			for range n {
				if _, err := w.Write(newline); err != nil {
					return
				}
				if _, err := w.Write(<-c.send); err != nil {
					return
				}
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *WSClient) handleSendMessage(content []byte) {
	var req request.CreateMessageRequest
	if err := json.Unmarshal(content, &req); err != nil {
		c.sendError(exceptions.ErrBadRequest.Code, exceptions.ErrBadRequest.Message, nil)
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		field, tag, param := validatorPkg.HandleRequestError(err)
		c.sendError(exceptions.ErrBadRequest.Code, exceptions.ErrBadRequest.Message, map[string]any{
			"field": field,
			"tag":   tag,
			"param": param,
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	messageRes, err := c.hub.messageUC.Create(ctx, c.schema, req)
	if err != nil {
		code, msg := handleErrorCode(err)
		c.sendError(code, msg, nil)
		return
	}

	data := map[string]any{
		"message": messageRes,
	}

	c.pushToHub(eventNewMessage, data, req.DepartmentID, req.BookingRoomID)
}

func (c *WSClient) handleMarkRead(content []byte) {
	var req request.UpdateMessageReadRequest
	if err := json.Unmarshal(content, &req); err != nil {
		c.sendError(exceptions.ErrBadRequest.Code, exceptions.ErrBadRequest.Message, nil)
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		field, tag, param := validatorPkg.HandleRequestError(err)
		c.sendError(exceptions.ErrBadRequest.Code, exceptions.ErrBadRequest.Message, map[string]any{
			"field": field,
			"tag":   tag,
			"param": param,
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.hub.messageUC.UpdateRead(ctx, c.schema, req); err != nil {
		code, msg := handleErrorCode(err)
		c.sendError(code, msg, nil)
		return
	}

	updatedRes := &response.UpdateMessageReadResponse{
		ChatID:     req.ChatID,
		ReaderID:   req.ReaderID,
		ReaderName: req.ReaderName,
		ReaderType: req.ReaderType,
	}

	data := map[string]any{
		"read": updatedRes,
	}

	c.pushToHub(eventMarkRead, data, req.DepartmentID, req.BookingRoomID)
}

func (c *WSClient) handleReaction(content []byte) {
	var req request.UpdateMessageReactionRequest
	if err := json.Unmarshal(content, &req); err != nil {
		c.sendError(exceptions.ErrBadRequest.Code, exceptions.ErrBadRequest.Message, nil)
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		field, tag, param := validatorPkg.HandleRequestError(err)
		c.sendError(exceptions.ErrBadRequest.Code, exceptions.ErrBadRequest.Message, map[string]any{
			"field": field,
			"tag":   tag,
			"param": param,
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updatedRes, err := c.hub.messageUC.UpdateReaction(ctx, c.schema, req)
	if err != nil {
		code, msg := handleErrorCode(err)
		c.sendError(code, msg, nil)
		return
	}

	data := map[string]any{
		"reaction": updatedRes,
	}

	c.pushToHub(eventReaction, data, req.DepartmentID, req.BookingRoomID)
}

func (c *WSClient) handleTyping(content []byte) {
	var req request.TypingMessageRequest
	if err := json.Unmarshal(content, &req); err != nil {
		c.sendError(exceptions.ErrBadRequest.Code, exceptions.ErrBadRequest.Message, nil)
		return
	}
	if err := validator.Validate.Struct(req); err != nil {
		field, tag, param := validatorPkg.HandleRequestError(err)
		c.sendError(exceptions.ErrBadRequest.Code, exceptions.ErrBadRequest.Message, map[string]any{
			"field": field,
			"tag":   tag,
			"param": param,
		})
		return
	}

	typingRes := &response.TypingMessageResponse{
		ChatID:     req.ChatID,
		SenderType: req.SenderType,
	}

	data := map[string]any{
		"typing": typingRes,
	}

	c.pushToHub(eventTyping, data, req.DepartmentID, req.BookingRoomID)
}

func (c *WSClient) sendError(code int, msg string, data any) {
	resp := WSResponse{
		Event:   eventError,
		Code:    code,
		Message: msg,
		Data:    data,
	}
	b, _ := json.Marshal(resp)
	c.send <- b
}

func (c *WSClient) getKey() string {
	switch {
	case c.clientType == constants.ClientTypeStaff && c.departmentID != nil:
		return fmt.Sprintf("%s-%s_%d", c.schema, c.clientType, *c.departmentID)
	case c.clientType == constants.ClientTypeStaff && c.departmentID == nil:
		return fmt.Sprintf("%s-admin", c.schema)
	default:
		return fmt.Sprintf("%s-%s_%d", c.schema, c.clientType, c.clientID)
	}
}

func (c *WSClient) pushToHub(event string, data any, departmentID, bookingRoomID int64) {
	targets := []string{
		fmt.Sprintf("%s-admin", c.schema),
		fmt.Sprintf("%s-%s_%d", c.schema, constants.ClientTypeStaff, departmentID),
		fmt.Sprintf("%s-%s_%d", c.schema, constants.ClientTypeGuest, bookingRoomID),
	}

	res := WSResponse{
		Event:   event,
		Code:    constants.CodeSuccess,
		Message: httpRes.ActionSuccess,
		Data:    data,
	}

	resBytes, _ := json.Marshal(res)
	c.hub.Broadcast(targets, resBytes)
}

func handleErrorCode(err error) (int, string) {
	var ex *exceptions.Exception
	if errors.As(err, &ex) {
		return ex.Code, ex.Message
	}

	return exceptions.ErrInternalServerError.Code, exceptions.ErrInternalServerError.Message
}
