package validators

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	errors2 "go-clean/shared/errors"
	"go-clean/shared/validators/lang"
	"reflect"
	"strings"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	idTranslations "github.com/go-playground/validator/v10/translations/id"
)

type injection struct {
	payload interface{}
	db      *gorm.DB
	ctx     *gin.Context
}

func ValidateStruct(ctx *gin.Context, payload interface{}) error {
	validate := validator.New()
	registerTagName(validate)
	trans := registerLanguage(validate, ctx.Value("lang").(string))

	inject := injection{
		payload: payload,
		//db:      db,
		ctx: ctx,
	}

	validate.RegisterValidation("exists", inject.validateExists)
	validate.RegisterValidation("unique_email", inject.validateUniqueEmail)

	err := validate.Struct(payload)
	if err == nil {
		return nil
	}

	validationErrors := err.(validator.ValidationErrors)
	for _, e := range validationErrors {
		return errors.New(e.Translate(trans))
	}

	return nil
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

func exists(ctx *gin.Context, db *gorm.DB, slices []string, value any) bool {
	var count int64

	res := db.WithContext(ctx).Raw(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ?", slices[0], slices[1]), value)

	if len(slices) == 4 {
		if slices[3] == "NULL" {
			res = db.WithContext(ctx).Raw(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ? and %s is NULL", slices[0], slices[1], slices[2]), value)
		} else {
			if slices[3] == "NULL" {
				res = db.WithContext(ctx).Raw(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ? and %s = ?", slices[0], slices[1], slices[2]), value, slices[3])
			}
		}
	}

	res.Count(&count)
	return count > 0
}

func (inject *injection) validateExists(fl validator.FieldLevel) bool {
	ctx := context.Background().(*gin.Context)
	value := fl.Field().Interface()
	slices := strings.Split(fl.Param(), ";")
	validated := true

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

	return validated
}

func (inject *injection) validateUniqueEmail(fl validator.FieldLevel) bool {
	ctx := context.Background().(*gin.Context)
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

	var count1, count2 int64

	if email == "" {
		return false
	}

	res := inject.db.WithContext(ctx).Table("employers").Where("email = ?", email).Where("deleted_at is null")
	res2 := inject.db.WithContext(ctx).Table("talents").Where("email = ?", email).Where("deleted_at is null")

	if except != "" {
		res = res.Where(fmt.Sprintf("%s != ?", slices[0]), except)
		res2 = res2.Where(fmt.Sprintf("%s != ?", slices[0]), except)
	}

	if err := res.Count(&count1).Error; errors2.HasError(err) {
		return false
	}
	if err := res2.Count(&count2).Error; errors2.HasError(err) {
		return false
	}

	return (count1 + count2) == 0
}
