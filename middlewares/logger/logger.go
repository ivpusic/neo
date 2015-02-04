// Middleware for logging all requests and respones to Neo server.
// You will see output in your console and in all other appenders configured for
// ``github.com/ivpusic/neo`` logger.
package logger

import (
	"github.com/ivpusic/golog"
	"github.com/ivpusic/neo"
	"strconv"
	"time"
)

func Log(this *neo.Ctx, next neo.Next) {
	start := time.Now()
	logger := golog.GetLogger("github.com/ivpusic/neo")
	method := this.Req.Method
	path := this.Req.URL.Path

	logger.Info("--> [Req] " + method + " to " + path)

	next()

	status := strconv.Itoa(this.Res.Status)
	elapsed := int(time.Now().Sub(start) / time.Millisecond)

	logger.Info("<-- [Res] (" + status + ") " + method + " to " + path + " Took " + strconv.Itoa(elapsed))
}
