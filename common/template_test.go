package common_test

import (
	"fmt"
	"testing"
)
import "github.com/tansoz/golayer/common"

func TestNewTemplateDataType(t *testing.T) {

	tdt := common.NewTemplateDataType()
	tdt.AddDataType("uint", ".+")
	tpl := tdt.NewTemplate("123{man:uint}45{man1:uint5}6")

	fmt.Println(tpl.Match("123445{man1:uint5}6"))
	fmt.Println(tpl.Test("123456"))
	fmt.Println(tpl.Test("123445{man1:uint5}6"))

	fmt.Println(tpl.ToString(map[string]interface{}{
		"man":  "'\\4",
		"man1": "0",
	}, "\\'"))
}
