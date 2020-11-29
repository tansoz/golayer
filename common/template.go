package common

import (
	"fmt"
	"regexp"
)

type Template interface {
	Match(string) map[interface{}]string
	Test(string) bool                               // test string is match with template
	ToString(map[string]interface{}, string) string // template output string
}

type TemplateImpl struct {
	Template
	RawString      string // Unprocessed text
	DataTypeList   TemplateDataType
	Args           []string
	TemplateParts  []string
	TemplateRegexp *regexp.Regexp
}

func (this *TemplateImpl) Match(matchStr string) map[interface{}]string {

	result := make(map[interface{}]string)

	v := this.TemplateRegexp.FindStringSubmatch(matchStr)
	if len(v) > 1 {
		for k, i := range v[1:] {
			result[this.Args[k]] = i
		}
	}

	return result
}

func (this *TemplateImpl) Test(testStr string) bool {
	return this.TemplateRegexp.MatchString(testStr)
}

func (this *TemplateImpl) escape(str string, charlist string) string {

	return regexp.MustCompile("(["+regexp.MustCompile("\\\\").ReplaceAllString(charlist, "\\\\")+"])").ReplaceAllString(str, "\\$1")

}

func (this *TemplateImpl) ToString(data map[string]interface{}, charlist string) string {

	str := ""
	index := 0

	for _, i := range this.TemplateParts {

		if i != "{{%TEMPLATE%}}" {
			str += i
		} else {
			if data[this.Args[index]] == nil {
				return ""
			} else {
				str += this.escape(fmt.Sprint(data[this.Args[index]]), charlist)
			}
			index++
		}

	}
	if !this.TemplateRegexp.MatchString(str) {

		return ""
	}

	return str
}

func (this *TemplateImpl) parse() {

	matchRegexp := regexp.MustCompile("\\{ *([^ :]+) *: *([^ })]+)(\\([^\\]:]+\\))? *\\}")

	index := matchRegexp.FindAllStringSubmatchIndex(this.RawString, -1)

	begin := 0
	tempRegexpStr := ""

	for _, i := range index {

		this.TemplateParts = append(this.TemplateParts, this.RawString[begin:i[0]])
		tempRegexpStr += this.RawString[begin:i[0]]
		begin = i[1]
		if typeValue := this.DataTypeList.GetDataType(this.RawString[i[4]:i[5]]); typeValue != "" {
			this.TemplateParts = append(this.TemplateParts, "{{%TEMPLATE%}}")
			this.Args = append(this.Args, this.RawString[i[2]:i[3]])
			if i[6] != -1 && i[7] != -1 {
				tempRegexpStr += "(" + typeValue + "{" + this.RawString[i[6]:i[7]] + "}" + ")"
			} else {
				tempRegexpStr += "(" + typeValue + ")"
			}
		} else {
			this.TemplateParts = append(this.TemplateParts, this.RawString[i[0]:i[1]])
			tempRegexpStr += this.RawString[i[0]:i[1]]
		}

	}

	this.TemplateParts = append(this.TemplateParts, this.RawString[begin:])
	tempRegexpStr += this.RawString[begin:]

	this.TemplateRegexp = regexp.MustCompile("^" + tempRegexpStr + "$")
}
