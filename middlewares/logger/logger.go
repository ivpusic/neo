// Middleware for logging all requests and respones to Neo server.
// You will see output in your console and in all other appenders configured for
// ``github.com/ivpusic/neo`` logger.
package logger

import (
	"strconv"
	"time"

	"github.com/ivpusic/golog"
	"github.com/ivpusic/neo"
)

func Log(ctx *neo.Ctx, next neo.Next) {
	start := time.Now()
	logger := golog.GetLogger("github.com/ivpusic/neo")
	method := ctx.Req.Method
	path := ctx.Req.URL.Path

	logger.Info("--> [Req] " + method + " to " + path)

	next()

	status := strconv.Itoa(ctx.Res.Status)
	elapsed := int(time.Now().Sub(start) / time.Millisecond)

	logger.Info("<-- [Res] (" + status + ") " + method + " to " + path + " Took " + strconv.Itoa(elapsed) + "ms")
}
