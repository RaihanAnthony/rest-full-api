package helper

import (
	"github.com/go-playground/validator/v10"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/locales/en"
	id_translations "github.com/go-playground/validator/v10/translations/id"
)

func TranslateToIndonesia() (ut.Translator, *validator.Validate) {
	en := en.New()
	id := ut.New(en, en)
	trans, _ := id.GetTranslator("id")

	validate := validator.New()
	id_translations.RegisterDefaultTranslations(validate, trans)

	return trans, validate
}