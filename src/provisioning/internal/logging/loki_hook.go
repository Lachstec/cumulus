package logging

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

type LokiHook struct {
	Endpoint string
}

func (h *LokiHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ts := strconv.FormatInt(time.Now().UnixNano(), 10)

	labels := map[string]string{
		"level": level.String(),
	}

	payload := map[string]interface{}{
		"streams": []map[string]interface{}{
			{
				"stream": labels,
				"values": [][]string{
					{ts, msg},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", h.Endpoint, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
}
