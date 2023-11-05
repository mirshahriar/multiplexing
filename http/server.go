package http

import (
	echo "github.com/labstack/echo/v4"
	mgrpc "github.com/mirshahriar/multiplexing/grpc"
	"github.com/mirshahriar/multiplexing/grpc/proto"
	"net/http"
)

type httpServer struct {
	grpcHandler proto.EchoServiceServer
}

func NewHTTPServer() *http.Server {
	echoRouter := echo.New()

	hServer := &httpServer{
		grpcHandler: mgrpc.NewGRPCHandler(),
	}

	echoRouter.GET("/echo", hServer.echoMessage)
	return echoRouter.Server
}

func (h *httpServer) echoMessage(c echo.Context) error {
	var req proto.EchoRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	resp, err := h.grpcHandler.EchoMessage(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"response": resp,
		"from":     "http",
	})
}
