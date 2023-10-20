package drivers

import (
	"bufio"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

// get SSE messages on connection and on update
// on update, send message to client, where the client will reload the page
// see ./hot-reload.templ
func (c *Controller) getSSE(fc *fiber.Ctx) error {
	fc.Set("Content-Type", "text/event-stream")
	fc.Set("Cache-Control", "no-cache")
	fc.Set("Connection", "keep-alive")

	fc.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		c.log.Debug("SSE connection established")
		fmt.Fprintf(w, "data: connected\n\n")
		err := w.Flush()

		if err != nil {
			c.log.Error("Error while flushing", "error", err)
		}

		<-c.onStart
		fmt.Fprintf(w, "data: updated\n\n")
		c.log.Debug("SSE update message sent")

		err = w.Flush()
		if err != nil {
			c.log.Error("Error while flushing", "error", err)
		}
		return
	}))

	return nil
}
