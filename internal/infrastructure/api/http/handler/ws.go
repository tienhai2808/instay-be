package handler

import (
	"net/http"

	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/realtime/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WSHandler struct {
	wsHub      *ws.HubWS
	domainHTTP *middleware.DomainMiddleware
}

func NewWSHandler(wsHub *ws.HubWS, domainHTTP *middleware.DomainMiddleware) *WSHandler {
	return &WSHandler{
		wsHub:      wsHub,
		domainHTTP: domainHTTP,
	}
}

func (h *WSHandler) ServeWS(c *gin.Context) {
	wsUpgrader := upgrader
	wsUpgrader.CheckOrigin = func(r *http.Request) bool {
		if h.domainHTTP == nil {
			return false
		}

		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = r.Host
		}

		return h.domainHTTP.AllowOriginWithContext(c, origin)
	}

	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		_ = c.Error(err)
		return
	}

	schema := c.GetString(middleware.CtxSchema)
	clientID := c.GetInt64(middleware.CtxClientID)
	clientType := c.GetString(middleware.CtxClientType)
	departmentIDAny, _ := c.Get(middleware.CtxDepartmentID)
	departmentID, _ := departmentIDAny.(*int64)

	if clientID == 0 || clientType == "" {
		_ = c.Error(exceptions.ErrUnauth)
		return
	}

	client := ws.NewWSClient(schema, h.wsHub, conn, clientID, clientType, departmentID)
	h.wsHub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
