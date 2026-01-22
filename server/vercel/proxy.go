package vercel

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	apiv1 "github.com/usememos/memos/server/router/api/v1"
)

type VercelProxy struct {
	echoServer *echo.Echo
}

func NewVercelProxy(apiV1Service *apiv1.APIV1Service) *VercelProxy {
	echoServer := echo.New()
	echoServer.HideBanner = true
	echoServer.HidePort = true
	echoServer.Use(middleware.Recover())
	echoServer.Use(middleware.CORS())
	echoServer.Use(middleware.Logger())

	ctx := context.Background()

	if err := apiV1Service.RegisterGateway(ctx, echoServer); err != nil {
		panic(err)
	}

	return &VercelProxy{
		echoServer: echoServer,
	}
}

func (v *VercelProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	recorder := httptest.NewRecorder()

	v.echoServer.ServeHTTP(recorder, r)

	for key, values := range recorder.Header() {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(recorder.Code)
	w.Write(recorder.Body.Bytes())
}

func (v *VercelProxy) HandleConnectRPC(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/memos.api.v1.") {
		v.ServeHTTP(w, r)
	}
}

func (v *VercelProxy) HandleGatewayAPI(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/v1/") {
		v.ServeHTTP(w, r)
	}
}
