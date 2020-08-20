/*
Copyright © 2020 Kevin Swiber <kswiber@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package resources

import (
	"encoding/json"
	"strconv"
)

// UserResponse represents the top-level struct of a user response in the
// Postman API.
type UserResponse struct {
	User User `json:"user,omitempty"`
}

// User represents the user info associated with a user request in the
// Postman API.
type User struct {
	ID string `json:"-,omitempty"`
}

// UnmarshalJSON sets the receiver to a copy of data.
func (r *User) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	r.ID = strconv.Itoa(int(v["id"].(float64)))

	return nil
}

// Format returns column headers and values for the resource.
func (r User) Format() ([]string, []interface{}) {
	s := make([]interface{}, 1)
	s[0] = r

	return []string{"ID"}, s
}
