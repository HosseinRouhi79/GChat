package validation

import (
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	NUMBERS = "0123456789"
	LETTERS  = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	SIGNS  = "!@#$%^&*"
)

type Specification interface {
	IsSatisfied(pass string) bool
}

func PasswordValidator(fld validator.FieldLevel) bool {
	var finalValue bool
	value, ok := fld.Field().Interface().(string)
	if !ok {
		return false
	}
	var passNum Number
	var passLetter Letter
	var passSign Sign
	resNum := passNum.IsSatisfied(value)
	resLetter := passLetter.IsSatisfied(value)
	resSign := passSign.IsSatisfied(value)
	if resNum && resLetter && resSign {
		finalValue = true
	}
	return finalValue
}

type Number struct {
	Content string
	Satisfy bool
}

type Letter struct {
	Content string
	Satisfy bool
}

type Sign struct {
	Content string
	Satisfy bool
}

func (n Number) IsSatisfied(pass string) bool {
	n.Content = pass
	elements := strings.Split(n.Content, "")
	numberSlice := strings.Split(NUMBERS, "")
	for _, e := range elements {
		if  slices.Contains(numberSlice, e){
			n.Satisfy = true 
			break
		} else {
			n.Satisfy = false
		}
	}
	return n.Satisfy
}

func (l Letter) IsSatisfied(pass string) bool {
	l.Content = pass
	elements := strings.Split(l.Content, "")
	numberSlice := strings.Split(LETTERS, "")
	for _, e := range elements {
		if  slices.Contains(numberSlice, e){
			l.Satisfy = true
			break
		} else {
			l.Satisfy = false
		}
	}
	return l.Satisfy
}

func (s Sign) IsSatisfied(pass string) bool {
	s.Content = pass
	elements := strings.Split(s.Content, "")
	numberSlice := strings.Split(SIGNS, "")
	for _, e := range elements {
		if  slices.Contains(numberSlice, e){
			s.Satisfy = true
			break
		} else {
			s.Satisfy = false
		}
	}
	return s.Satisfy
}
