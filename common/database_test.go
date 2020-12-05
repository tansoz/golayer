package common_test

import (
	"encoding/json"
	"fmt"
	"github.com/tansoz/golayer/common"
	"testing"
)

type user struct {
	UserId     int    `json:"userid"`
	Username   string `json:"username" default:"1566"`
	Password   string `json:"password" default:"15"`
	Age        int    `json:"age"`
	CreateTime int    `json:"createDate"`
}

func TestDatabaseImpl_Connection(t *testing.T) {
	db := common.NewDatabase(
		"127.0.0.1",
		3309,
		"root",
		"root",
		"player_test",
		"utf8mb4",
	)
	db.Connection()

	result := db.Query("select * from `user`", nil)
	if result != nil {
		ok := db.FetchAll(result, user{})
		fmt.Println(ok)
		result.Close()

		jsonStr, _ := json.Marshal(ok)
		fmt.Println(string(jsonStr))

	}

}
