package common_test

import (
	"fmt"
	"testing"
)
import "github.com/tansoz/golayer/common"

func TestNewTemplateDataType(t *testing.T) {

	tdt := common.NewTemplateDataType()
	tdt.AddDataType("INT", "[-+]?\\d")
	tdt.AddDataType("VARCHAR", ".")
	tdt.AddDataType("DATETIME", "\\d{4}-\\d{1,2}-\\d{1,2} \\d{1,2}:\\d{1,2}:\\d{1,2}")
	tpl := tdt.NewTemplate("SELECT * FROM `goods` WHERE `goods_id` = '{goods_id:INT(+)}'")

	fmt.Println(tpl.ToString(map[string]interface{}{
		"goods_id": "4",
		"man1":     "0",
	}, "\\'"))
}
