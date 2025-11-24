package api

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/google/uuid"
)

func RegisterValidation() (ut.Translator, error) {

	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := en_translations.RegisterDefaultTranslations(v, trans)
		if err != nil {
			return nil, err
		}

		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

	return trans, nil
}

type UUIDParam struct {
	UUID string `uri:"uuid" binding:"required,uuid" json:"uuid"`
}

func getUUIDParam(c *gin.Context) (uuid.UUID, bool) {
	var uuidParam UUIDParam
	err := c.ShouldBindUri(&uuidParam)
	if err != nil {
		_ = c.Error(err)
		return uuid.Max, false
	}

	return uuid.MustParse(uuidParam.UUID), true
}
