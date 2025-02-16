package api

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func BindJSONStrict(c *gin.Context, obj interface{}) error {
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(obj); err != nil {
		return err
	}

	return nil
}
