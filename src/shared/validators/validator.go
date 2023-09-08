package validators

import (
	"context"
	"encoding/json"
	"fmt"
	"go-clean/src/models/messages"
	"go-clean/src/shared/validators/lang"
	"net/http"
	"reflect"
	"strings"

	"gorm.io/gorm"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	idTranslations "github.com/go-playground/validator/v10/translations/id"
	"github.com/gofiber/fiber/v2"
)

type injection struct {
	payload interface{}
	db      *gorm.DB
	ctx     *fiber.Ctx
}

func ValidateStruct(ctx *fiber.Ctx, payload interface{}) (map[string][]string, error) {
	validate := validator.New()
	registerTagName(validate)
	trans := registerLanguage(validate, ctx.Locals("lang").(string))

	inject := injection{
		payload: payload,
		//db:      db,
		ctx: ctx,
	}

	validate.RegisterValidation("exists", inject.validateExists)
	validate.RegisterValidation("unique_email", inject.validateUniqueEmail)

	err := validate.Struct(payload)
	if err != nil {
		formErrors := make(map[string][]string)
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			field := e.Field()
			message := e.Translate(trans)
			formErrors[field] = append(formErrors[field], message)
		}

		if len(formErrors) > 0 {
			return formErrors, &messages.ErrorWrapper{
				Context:    ctx,
				ErrorCode:  "ERROR-400003",
				StatusCode: http.StatusBadRequest,
			}
		}
	}

	return nil, nil
}

func registerLanguage(validate *validator.Validate, language string) ut.Translator {
	customTranslations := make(map[string]map[string]string)
	customTranslations["en"] = lang.CustomEnValidatorTranslation

	translators := make(map[string]locales.Translator)
	translators["en"] = en.New()
	translators["id"] = id.New()

	uni := ut.New(translators[language], translators[language])
	trans, _ := uni.GetTranslator(language)

	registerTranslation(validate, customTranslations[language], trans)

	switch language {
	case "id":
		_ = idTranslations.RegisterDefaultTranslations(validate, trans)
	default:
		_ = enTranslations.RegisterDefaultTranslations(validate, trans)
	}
	return trans
}

func registerTagName(validate *validator.Validate) {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
			if name == "-" {
				return ""
			}
		}
		return name
	})
}

func registerTranslation(validate *validator.Validate, rules map[string]string, trans ut.Translator) {
	for key, value := range rules {
		err := validate.RegisterTranslation(key, trans, func(ut ut.Translator) error {
			return ut.Add(key, value, true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			name := strings.TrimSuffix(fe.Field(), "_id")
			name = strings.ReplaceAll(name, "_", " ")
			t, _ := ut.T(fe.Tag(), name)
			return t
		})
		if err != nil {
			return
		}
	}
}

func exists(ctx *fiber.Ctx, db *gorm.DB, slices []string, value any) bool {
	var count int64

	query := db.Table(slices[0]).Where(fmt.Sprintf("%s = ?", slices[1]), value)

	if len(slices) == 4 {
		if slices[3] == "NULL" {
			query = query.Where(fmt.Sprintf("%s IS NULL", slices[2]))
		} else {
			query = query.Where(fmt.Sprintf("%s = ?", slices[2]), slices[3])
		}
	}

	query.Count(&count)
	return count > 0
}

func (inject *injection) validateExists(fl validator.FieldLevel) bool {
	ctx := fl.Top().Interface().(*fiber.Ctx)
	value := fl.Field().Interface()
	slices := strings.Split(fl.Param(), ";")

	if fl.Field().Kind() == reflect.Slice {
		switch value.(type) {
		case []string:
			for _, i := range value.([]string) {
				if exists(ctx, inject.db, slices, i) == false {
					return false
				}
			}
		case []int:
			for _, i := range value.([]int) {
				if exists(ctx, inject.db, slices, i) == false {
					return false
				}
			}
		}

	} else {
		return value == "" || exists(ctx, inject.db, slices, value)
	}

	return true
}

func (inject *injection) validateUniqueEmail(fl validator.FieldLevel) bool {
	email := fl.Field().Interface().(string)
	slices := strings.Split(fl.Param(), ";")

	var jsonData map[string]interface{}

	marshalled, _ := json.Marshal(inject.payload)
	err := json.Unmarshal(marshalled, &jsonData)
	if err != nil {
		return false
	}

	except := ""

	if len(slices) > 0 {
		if slices[0] != "" && slices[0] != "-" {
			except = jsonData[slices[0]].(string)
		}
	}

	var count int64

	if email == "" {
		return false
	}

	ctx := context.Background()
	res := inject.db.WithContext(ctx).Table("users").Where("email = ?", email).Where("deleted_at is null")

	if except != "" {
		res = res.Where(fmt.Sprintf("%s != ?", slices[0]), except)
	}

	if err := res.Count(&count).Error; messages.HasError(err) {
		return false
	}

	return count == 0
}
