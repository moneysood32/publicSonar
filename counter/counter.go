package counter

import (
	"encoding/json"
	"fmt"
)

// UserInfo contains the input data fields
type UserInfo struct {
	ID     string `json:"tenantID"`
	Tenant string `json:"tenant"`
}

// ParseAndSave parses the request body and if the data isn't malformed or corrupt, saves it.
func ParseAndSave(data []byte) (*UserInfo, error) {
	var user UserInfo
	fmt.Println(data)
	if err := json.Unmarshal(data, &user); err != nil {
		panic(err)
	}
	fmt.Println(user.ID)
	fmt.Println(user.Tenant)
	return &UserInfo{}, nil
}

// String implements stringer interface over UserInfo,
// required to print UserInfo in a user friendly way to the output stream.
func (u *UserInfo) String() string {
	return "TenantID : " + u.ID + " , " + "Tenant : " + u.Tenant
}
