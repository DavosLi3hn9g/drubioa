package fun

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf16"
)

//strip_tags
func StripTags(str string) string {
	//将HTML标签全转换成小写
	//re, _ := regexp.Compile(`\<[\S\s]+?\>`)
	//str = re.ReplaceAllStringFunc(str, strings.ToLower)

	//去除STYLE
	re, _ := regexp.Compile(`\<style[\S\s]+?\</style\>`)
	str = re.ReplaceAllString(str, "")

	//去除SCRIPT
	re, _ = regexp.Compile(`\<script[\S\s]+?\</script\>`)
	str = re.ReplaceAllString(str, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile(`\<[\S\s]+?\>`)
	str = re.ReplaceAllString(str, "")

	return str
}

//substr 截取字符串，支持中文
func SubStr(str string, start int, end int) string {
	rs := []rune(str)
	length := int(len(rs))

	if start < 0 {
		start = length + start //负数 - 在从字符串结尾开始的指定位置开始
	}
	if end < 0 {
		end = length + end //负数 - 从字符串末端返回的长度
	} else if end == 0 {
		end = length
	} else if start > 0 {
		end = start + end //正数 - 从 start 参数所在的位置返回的长度
	}

	if int(math.Abs(float64(end))) > length || int(math.Abs(float64(start))) > length {
		return ""
	}

	return string(rs[start:end])
}

// addslashes() 函数返回在预定义字符之前添加反斜杠的字符串。
// 预定义字符是：
// 单引号（'）
// 双引号（"）
// 反斜杠（\）
func Addslashes(str string) string {
	tmpRune := []rune{}
	strRune := []rune(str)
	for _, ch := range strRune {
		switch ch {
		case []rune{'\\'}[0], []rune{'"'}[0], []rune{'\''}[0]:
			tmpRune = append(tmpRune, []rune{'\\'}[0])
			tmpRune = append(tmpRune, ch)
		default:
			tmpRune = append(tmpRune, ch)
		}
	}
	return string(tmpRune)
}

// stripslashes() 函数删除由 addslashes() 函数添加的反斜杠。
func Stripslashes(str string) string {
	dstRune := []rune{}
	strRune := []rune(str)
	strLenth := len(strRune)
	for i := 0; i < strLenth; i++ {
		if strRune[i] == []rune{'\\'}[0] {
			i++
		}
		dstRune = append(dstRune, strRune[i])
	}
	return string(dstRune)
}

//stripos , 未找到 -1， 查找字符串在另一字符串中第一次出现的位置（不区分大小写）, 未找到 -1
func Stripos(str string, index string) int {
	return strings.Index(strings.ToLower(str), strings.ToLower(index))
}

//strpos , 未找到 -1，查找字符串在另一字符串中第一次出现的位置（区分大小写）
func Strpos(str string, index string) int {
	return strings.Index(str, index)
}

//strripos , 未找到 -1， 查找字符串在另一字符串中最后一次出现的位置（不区分大小写）
func Strripos(str string, index string) int {
	return strings.LastIndex(strings.ToLower(str), strings.ToLower(index))
}

//strrpos , 未找到 -1， 查找字符串在另一字符串中最后一次出现的位置（区分大小写）
func Strrpos(str string, index string) int {
	return strings.LastIndex(str, index)
}

func PregMatchAll(pattern string, subject string, matches *[][]string, flags string, offset int) bool {
	data := regexp.MustCompile(pattern).FindAllStringSubmatch(subject, -1)
	switch flags {
	case "PREG_PATTERN_ORDER":
		*matches = data
	case "PREG_OFFSET_CAPTURE":
		//todo
	default: //"PREG_SET_ORDER" //有bug
		matchAll := make(map[int][]string, 2)
		for _, va := range data {
			for kb, vb := range va {
				matchAll[kb] = append(matchAll[kb], vb)
			}
		}
		for _, a := range matchAll {
			*matches = append(*matches, a)
		}
	}
	if *matches == nil {
		return false
	} else {
		return true
	}
}

func PregReplace(str string, repl string, src string) string {
	return regexp.MustCompile(str).ReplaceAllString(src, repl)
}

func Unicode2String(form string) (to string, err error) {
	bs, err := hex.DecodeString(strings.Replace(form, `\u`, ``, -1))
	if err != nil {
		return
	}
	for i, bl, br, r := 0, len(bs), bytes.NewReader(bs), uint16(0); i < bl; i += 2 {
		binary.Read(br, binary.BigEndian, &r)
		to += string(r)
	}
	return
}

func HexUTF16FromString(s string) string {
	hex := fmt.Sprintf("%04x", utf16.Encode([]rune(s)))
	return strings.Replace(hex[1:len(hex)-1], " ", "", -1)
}

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
