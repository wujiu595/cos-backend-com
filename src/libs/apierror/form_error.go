package apierror

import (
	"reflect"
	"strings"

	"cos-backend-com/src/common/locales"

	"github.com/wujiu2020/strip/utils/apires"
	validator "gopkg.in/go-playground/validator.v9"
)

type FormError struct {
	Errors      *FieldError   `json:"error,omitempty"`
	ErrorFields []*FieldError `json:"errorFields,omitempty"`
}

type FieldError struct {
	Name      string      `json:"name"`
	NameSpace string      `json:"nameSpace"`
	Type      string      `json:"type"`
	Value     interface{} `json:"value"`
	locales.Content
}

func (p *FormError) Error() string {
	errs := make([]string, 0, len(p.ErrorFields)+1)
	if p.Errors != nil {
		errs = append(errs, p.Errors.Error())
	}
	for _, er := range p.ErrorFields {
		errs = append(errs, er.Error())
	}
	return strings.Join(errs, "\n")
}

func (p *FieldError) Error() string {
	return p.Message
}

func NewFormFieldError(namespace string, value interface{}, locale locales.Content) *FieldError {
	var name string
	parts := strings.SplitAfterN(namespace, ".", 1)
	if len(parts) > 0 {
		name = parts[len(parts)-1]
	}

	typ := reflect.TypeOf(value)

	return &FieldError{
		Name:      name,
		NameSpace: namespace,
		Type:      typ.Kind().String(),
		Value:     value,
		Content:   locale,
	}
}

func ResponseFormFieldError(errorFields ...*FieldError) *apires.ResError {
	ret := &FormError{}
	ret.ErrorFields = errorFields
	return ErrBadRequest.WithData(ret, ret.Error())
}

func ResponseFormError(namespace string, value interface{}, locale locales.Content) *apires.ResError {
	ret := &FormError{}
	ret.Errors = NewFormFieldError(namespace, value, locale)
	return ErrBadRequest.WithData(ret, ret.Error())
}

func validationErrorConvert(err interface{}) (data interface{}) {
	switch er := err.(type) {
	case validator.ValidationErrors:
		fields := make([]*FieldError, 0, len(er))
		for _, e := range er {
			fields = append(fields, validationFieldErrorConvert(e))
		}
		data = &FormError{ErrorFields: fields}

	case validator.FieldError:
		field := validationFieldErrorConvert(er)
		data = &FormError{ErrorFields: []*FieldError{field}}
	}
	return
}

func validationFieldErrorConvert(f validator.FieldError) *FieldError {
	t := f.Translate(nil)
	messageId := "validator." + f.Tag()
	ret := &FieldError{
		Name:      f.Field(),
		NameSpace: f.Namespace(),
		Type:      f.Kind().String(),
		Value:     f.Value(),
		Content: locales.Content{
			Message:       t,
			MessageId:     messageId,
			MessageParams: []string{f.Field()},
		},
	}
	if ret.MessageParams == nil {
		ret.MessageParams = []string{}
	}
	return ret
}
