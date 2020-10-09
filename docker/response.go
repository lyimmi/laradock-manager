package docker

import (
	b64 "encoding/base64"
	"encoding/json"

	"github.com/labstack/gommon/log"
)

// Response struct
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (t *Compose) emitError(err string) {
	uEnc := b64.URLEncoding.EncodeToString([]byte(err))
	t.runtime.Events.Emit("backendError", uEnc)
}

// returnResponse a response json
func returnResponse(success bool, message string) string {
	resp, err := json.Marshal(Response{Success: success, Message: message})
	if err != nil {
		log.Error(err)
	}
	return string(resp)
}
