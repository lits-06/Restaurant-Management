package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"restaurant-management/api-gateway/internal/grpcclient"
	authpb "restaurant-management/proto/auth"
	"restaurant-management/shared/pkg/logger"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		switch origin {
		case "http://localhost:5173", "http://localhost:5174", "http://localhost:5175":
			return true
		}
		return false
	},
}

type NotificationHandler struct {
	notifClient *grpcclient.NotificationClient
	authClient  *grpcclient.AuthClient
}

func NewNotificationHandler(notifClient *grpcclient.NotificationClient, authClient *grpcclient.AuthClient) *NotificationHandler {
	return &NotificationHandler{notifClient: notifClient, authClient: authClient}
}

// Subscribe handles GET /ws/notifications?token=<jwt>&role=<CHEF|WAITER>
//
// The browser passes the access token as a query param because the WebSocket
// API does not allow custom request headers.
func (h *NotificationHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimSpace(r.URL.Query().Get("token"))
	if token == "" {
		token = extractBearerToken(r.Header.Get("Authorization"))
	}
	if token == "" {
		http.Error(w, "token required", http.StatusUnauthorized)
		return
	}

	resp, err := h.authClient.VerifyToken(r.Context(), &authpb.VerifyTokenRequest{AccessToken: token})
	if err != nil || !resp.Valid {
		http.Error(w, "invalid or expired token", http.StatusUnauthorized)
		return
	}

	requestedRole := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("role")))
	if requestedRole == "" {
		http.Error(w, "role required", http.StatusBadRequest)
		return
	}
	if !canSubscribeRole(resp.Roles, requestedRole) {
		http.Error(w, "insufficient role", http.StatusForbidden)
		return
	}

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("websocket upgrade failed", zap.Error(err))
		return
	}
	defer conn.Close()

	logger.Info("ws client connected",
		zap.String("user_id", resp.UserId),
		zap.String("role", requestedRole),
	)

	// ctx is cancelled when the WebSocket client disconnects.
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// Goroutine: drain incoming WebSocket frames so gorilla can detect close events.
	// Cancels ctx (and thus the gRPC stream) when the connection closes.
	go func() {
		defer cancel()
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	stream, err := h.notifClient.Subscribe(ctx, requestedRole)
	if err != nil {
		logger.Error("failed to open notification stream", zap.Error(err))
		conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "subscription failed"))
		return
	}

	for {
		notif, err := stream.Recv()
		if err != nil {
			if err != io.EOF && ctx.Err() == nil {
				logger.Warn("notification stream recv error", zap.Error(err))
			}
			break
		}

		data, err := json.Marshal(notif)
		if err != nil {
			continue
		}
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			break
		}
	}

	logger.Info("ws client disconnected",
		zap.String("user_id", resp.UserId),
		zap.String("role", requestedRole),
	)
}

// canSubscribeRole checks the caller has the requested role or is ADMIN/MANAGER.
func canSubscribeRole(callerRoles []string, requestedRole string) bool {
	for _, r := range callerRoles {
		if r == "ADMIN" || r == "MANAGER" || r == requestedRole {
			return true
		}
	}
	return false
}
