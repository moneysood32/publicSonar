package counter

import (
	"encoding/json"
	"fmt"
)

// UserInfo contains the input data fields
type UserInfo struct {
	Tenant string `json:"tenant"`
	ItemID string `json:"itemID"`
}

// UserInfoSet is a temporary in-memory store for the Items
var UserInfoSet = make(map[string][]UserInfo)

// ParseAndSave parses the request body and if the data isn't malformed or corrupt, saves it.
func ParseAndSave(data []byte) (*UserInfo, error) {
	var user UserInfo
	if err := json.Unmarshal(data, &user); err != nil {
		panic(err)
	}
	if _, ok := UserInfoSet[user.Tenant]; ok {
		return nil, fmt.Errorf("Item already exists")
	}
	UserInfoSet[user.Tenant] = append(UserInfoSet[user.Tenant], user)
	fmt.Println(UserInfoSet)
	if len(UserInfoSet) >= 2 {
		fmt.Println("max memory limit")
	}
	return &user, nil
}

// String implements stringer interface over UserInfo,
// required to print UserInfo in a user friendly way to the output stream.
func (u UserInfo) String() string {
	return "Tenant : " + u.Tenant + " , " + "Item : " + u.ItemID
}

// // WriteSetToDB writes input map to memory
// func (m) WriteSetToDB(userDataset map[string]string) err {
// }
