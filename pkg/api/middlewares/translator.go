package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translation "github.com/go-playground/validator/v10/translations/en"
	pt_br_translation "github.com/go-playground/validator/v10/translations/pt_BR"
	apicontext "github.com/vagnercardosoweb/go-rest-api/pkg/api/context"
	apirequest "github.com/vagnercardosoweb/go-rest-api/pkg/api/request"
	"golang.org/x/text/language"
)

var universalTranslator *ut.UniversalTranslator

func Translator(c *gin.Context) {
	acceptLanguage := apirequest.GetAcceptLanguage(c)

	languageTags, _, _ := language.ParseAcceptLanguage(acceptLanguage)
	acceptLanguage = strings.Replace(strings.ToLower(languageTags[0].String()), "-", "_", -1)

	// I18n Translator
	c.Set(apicontext.AcceptLanguageKey, acceptLanguage)
	// TODO: Implement I18n Translator

	// Validator Translator
	translator, _ := universalTranslator.GetTranslator(acceptLanguage)
	c.Set(apicontext.ValidatorTranslatorKey, &translator)

	c.Next()
}

func init() {
	enLocale := en.New()
	ptBRLocale := pt_BR.New()

	universalTranslator = ut.New(
		enLocale, // Fallback
		enLocale,
		ptBRLocale,
	)

	enTranslator, _ := universalTranslator.GetTranslator(enLocale.Locale())
	ptBRTranslator, _ := universalTranslator.GetTranslator(ptBRLocale.Locale())

	if val, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = en_translation.RegisterDefaultTranslations(val, enTranslator)
		_ = pt_br_translation.RegisterDefaultTranslations(val, ptBRTranslator)
	}
}
