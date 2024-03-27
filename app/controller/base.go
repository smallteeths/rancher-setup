package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rancher-setup/app/request"
	"github.com/rancher-setup/internal/variable/consts"
)

var validator *request.Request

type base struct {
}

func init() {
	var err error
	validator, err = request.New()
	if err != nil {
		log.Fatal(consts.ErrorInitConfig)
	}
}

func (base) Validate(ctx *gin.Context, param any) map[string]string {
	return validator.Validator(ctx, param)
}
