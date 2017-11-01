package checkers

import (
	"encoding/json"
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/go-kit/kit/log"
	"github.com/valyala/fasthttp"
	"os"
)

var (
	handlerChecker Checker
)

type Handler struct {
	Client fasthttp.Client
	addr   string
}

func NewHandler(addr string) Handler {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	logger.Log("Start Handler")
	return Handler{
		Client: fasthttp.Client{},
		addr:   addr,
	}
}

func (h *Handler) AddRoute(route string, checker Checker) {
	logger.Log("msg", "Start listen")
	router := fasthttprouter.New()
	router.GET(route, HealthCheck)
	handlerChecker = checker
	logger.Log("Fatal", fasthttp.ListenAndServe(h.addr, router.Handler))
}

func HealthCheck(ctx *fasthttp.RequestCtx) {
	if !ctx.IsGet() {
		return
	}
	health, err := handlerChecker.Check()

	if err != nil {
		logger.Log("Error in process check: %s", err.Error())
	}
	bytes, err := json.Marshal(health)
	if err != nil {
		logger.Log("Cannot marshal %+v", health)
	}
	ctx.Response.BodyWriter().Write(bytes)
	logger.Log("msg", "send:", "health", fmt.Sprintf("%+v", health))
}
