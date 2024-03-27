package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/rancher-setup/internal/request"
)

type RancherConfig struct {
	Host         string `binding:"required" form:"host" query:"host" json:"host"`
	Version      string `binding:"required" form:"version" query:"version" json:"version"`
	HarborUser   string `form:"harborUser" query:"harborUser" json:"harborUser"`
	HarborPasswd string `form:"harborPasswd" query:"harborPasswd" json:"harborPasswd"`
}

func (f RancherConfig) Message() validator.ValidationErrorsTranslations {
	return map[string]string{
		"Host.required":    "host 必填",
		"Version.required": "version 必填",
	}
}

var _ request.IValidator = (*RancherConfig)(nil)
