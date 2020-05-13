package validate_test

import (
	"fmt"
	"regexp"
	"testing"

	"cos-backend-com/src/common/validate"
)

type errFuncTypeTest struct {
	C int `validate:"func=parent.WrongRetValue"`
}

type errFuncNotExitTest struct {
	D int `validate:"func=self.NoDefinedFunc"`
}

type errValidateFuncParamTest struct {
	E int `validate:"func=NoDefinedFunc"`
}

type funcTest struct {
	A int `validate:"func=parent.Ge"`
	B int `validate:"func=parent.SumAGe0"`
}

type funcTestParent struct {
	test funcTest
}

func (t errFuncTypeTest) WrongRetValue() error {
	return nil
}

func (t funcTest) Ge() bool {
	if t.A > t.B {
		return true
	}
	return false
}

func (t funcTest) SumAGe0() bool {
	if (t.B + t.A) > 0 {
		return true
	}
	return false
}

var (
	noDefinedFuncErrRegex = regexp.MustCompile(`valid func:(.*?) of (.*?) does not exit`)
	funcTypeErrRegex      = regexp.MustCompile(`valid func:(.*?) of (.*?) must be "func\(\)bool"`)
	funcParamErrRegex     = regexp.MustCompile(`valid func:value must be styled "self.func" or "parent.func"`)
)

func TestRfe(t *testing.T) {
	type Foo struct {
		A int `validate:"required"`
		B int `validate:"required_with_eq=A:1"`
	}

	foo := Foo{
		A: 1,
		B: 0,
	}

	err := validate.Default.Struct(&foo)
	if err == nil {
		t.Fatal("should return error")
	}

	foo = Foo{
		A: 1,
		B: 1,
	}

	err = validate.Default.Struct(&foo)
	if err != nil {
		t.Fatal("should not return error")
	}

	foo = Foo{
		A: 2,
		B: 0,
	}

	err = validate.Default.Struct(&foo)
	if err != nil {
		t.Fatal("should not return error")
	}

	ft := funcTest{
		A: 3,
		B: 2,
	}
	err = validate.Default.Struct(&ft)
	if err != nil {
		t.Fatal("should not return error")
	}
	ft = funcTest{
		A: 2,
		B: -1,
	}
	err = validate.Default.Struct(&ft)
	if err != nil {
		t.Fatal("should not return error")
	}
}

func TestValidFuncWithWrongFuncType(t *testing.T) {
	defer func() {
		if re := recover(); re != nil {
			if !funcTypeErrRegex.MatchString(re.(string)) {
				t.Fatal("unexpected err panic")
			}
		}
	}()

	eftt := errFuncTypeTest{
		C: 1,
	}

	err := validate.Default.Struct(&eftt)
	if err != nil {
		t.Fatal("should not return error")
	}
}

func TestValidFuncNotExit(t *testing.T) {
	defer func() {
		if re := recover(); re != nil {
			if !noDefinedFuncErrRegex.MatchString(re.(string)) {
				t.Fatal("unexpected err panic")
			}
		}
	}()

	efnet := errFuncNotExitTest{
		D: 100,
	}

	err := validate.Default.Struct(&efnet)
	if err != nil {
		t.Fatal("should not return error")
	}
}

type A struct {
	X string `validate:"func=parent.Validate"`
}

func (a A) Validate() bool {
	fmt.Println("a:parent.Validate")
	return true
}

type B struct {
	X A
}

type TestA string

func (t TestA) Validate() bool {
	fmt.Println("c:self.validate")
	return true
}

type C struct {
	X B
	Y TestA `validate:"func=self.Validate,func=parent.Validate"`
	Z TestA `validate:"func=self.Validate|func=parent.Validate"`
}

func (c C) Validate() bool {
	fmt.Println("c:parent.Validate")
	return true
}

func TestValidFuncParent(t *testing.T) {
	efnet := C{X: B{A{"111"}}}

	err := validate.Default.Struct(&efnet)
	if err != nil {
		t.Fatal("should not return error")
	}
}

func TestValidFuncParamErr(t *testing.T) {
	defer func() {
		if re := recover(); re != nil {
			if !funcParamErrRegex.MatchString(re.(string)) {
				fmt.Println(re.(string))
				t.Fatal("unexpected err panic")
			}
		}
	}()

	efnet := errValidateFuncParamTest{
		E: 100,
	}

	err := validate.Default.Struct(&efnet)
	if err != nil {
		t.Fatal("should not return error")
	}
}
