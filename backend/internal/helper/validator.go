package helper

import (
	"reflect"
	"strings"

	localeEN "github.com/go-playground/locales/en"
	translator "github.com/go-playground/universal-translator"
	validators "github.com/go-playground/validator/v10"
	vTranslator "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	vld *validators.Validate
	trn *translator.Translator
}

type validationErrorResponse struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

func NewValidator() *Validator {
	v := validators.New()
	t := newTranslator()

	registerCustom(v, *t)

	return &Validator{
		vld: v,
		trn: t,
	}
}

func (v *Validator) Validate(data interface{}) []validationErrorResponse {
	if err := v.vld.Struct(data); err != nil {
		vErrs := v.translateValidateError(err)

		return vErrs
	}

	return nil
}

func (v *Validator) translateValidateError(err error) (errs []validationErrorResponse) {
	if err == nil {
		return errs
	}

	validationErrs, ok := err.(validators.ValidationErrors)
	if !ok {
		return errs
	}

	for _, e := range validationErrs {
		errMsg := e.Translate(*v.trn)
		fieldName := e.Field()

		errs = append(errs, validationErrorResponse{
			Field:   fieldName,
			Message: errMsg,
		})
	}

	return errs
}

func newTranslator() *translator.Translator {
	en := localeEN.New()
	trans := translator.New(en, en)
	transEN, _ := trans.GetTranslator("en")

	return &transEN
}

func registerCustom(v *validators.Validate, t translator.Translator) {
	_ = vTranslator.RegisterDefaultTranslations(v, t)

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "" {
			name = strings.SplitN(fld.Tag.Get("query"), ",", 2)[0]
		}
		if name == "" {
			name = strings.SplitN(fld.Tag.Get("param"), ",", 2)[0]
		}
		if name == "-" {
			return fld.Name
		}

		return name
	})
}
