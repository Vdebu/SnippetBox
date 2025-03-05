package models

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// EmailRX 使用正则表达式匹配用户输入的邮箱
var EmailRX = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Validator 验证器:包含验证信息是否准确的方法与存储错误信息的字段
type Validator struct {
	// 用于存储与存储结构或字段无关的错误 这里由于没有字段所以不是map 只需加入信息
	NonFieldErrors []string
	// 创建用于存储错误信息的map
	FieldErrors map[string]string
}

// Valid 如果没有返回任何错误就返回true,用于判断是否发生错误执行处理错误的逻辑
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// AddFieldError 向FieldErrors添加错误信息
func (v *Validator) AddFieldError(key, message string) {
	// 如果当前的map还没有初始化就初始化再加入进行操作防止panic
	if v.FieldErrors == nil {
		v.FieldErrors = map[string]string{}
	}
	// 如果当前错误信息没在map中就创建
	if _, exists := v.FieldErrors[key]; !exists {
		// 如果存在就什么都不干
		v.FieldErrors[key] = message
	}
}

// AddNonFieldError 向NonFieldError中添加错误信息
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// CheckField 通过与其他返回bool值的函数组合用于检测是否发生错误
func (v *Validator) CheckField(ok bool, key, msg string) {
	// 发生错误就将信息存入FieldErrors
	if !ok {
		v.AddFieldError(key, msg)
	}
}

// NotBlank 判断字符串是否是空的
func (v *Validator) NotBlank(value string) bool {
	// 更易理解的条件判断语句
	return strings.TrimSpace(value) != ""
}

// MaxChars 判断字符串的长度是超出最大限制
func (v *Validator) MaxChars(value string, max int) bool {
	return utf8.RuneCountInString(value) <= max
}

// PermittedValue 使用泛型进行优化 泛型不能绑定在结构体上作为方法
// 检查用户填写的值是否是合法取值
func PermittedValue[T comparable](val T, permittedValues ...T) bool {
	for _, v := range permittedValues {
		if val == v {
			return true
		}
	}
	return false
}

// MinChars 检查是否满足最小字符要求
func (v *Validator) MinChars(value string, min int) bool {
	return utf8.RuneCountInString(value) >= min
}

// Matches 检查是否匹配字符串格式(通过正则表达式)
func (v *Validator) Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func Confirms[T comparable](s1, s2 T) bool {
	return s1 != s2
}
