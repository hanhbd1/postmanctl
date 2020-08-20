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
	"time"
)

// MonitorListResponse represents the top-level monitors response from the
// Postman API.
type MonitorListResponse struct {
	Monitors MonitorListItems `json:"monitors,omitempty"`
}

// MonitorListItems is a slice of MonitorListItem.
type MonitorListItems []MonitorListItem

// Format returns column headers and values for the resource.
func (r MonitorListItems) Format() ([]string, []interface{}) {
	s := make([]interface{}, len(r))
	for i, v := range r {
		s[i] = v
	}

	return []string{"UID", "Name"}, s
}

// MonitorListItem represents a single item in an MonitorListResponse.
type MonitorListItem struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	UID   string `json:"uid,omitempty"`
	Owner string `json:"owner,omitempty"`
}

// UnmarshalJSON sets the receiver to a copy of data.
func (m *MonitorListItem) UnmarshalJSON(data []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if val, ok := v["id"]; ok {
		m.ID = val.(string)
	}

	if val, ok := v["name"]; ok {
		m.Name = val.(string)
	}

	if val, ok := v["uid"]; ok {
		m.UID = val.(string)
	}

	if val, ok := v["owner"]; ok {
		m.Owner = strconv.Itoa(int(val.(float64)))
	}

	return nil
}

// MonitorResponse is the top-level monitor response from the
// Postman API.
type MonitorResponse struct {
	Monitor Monitor `json:"monitor,omitempty"`
}

// Monitor represents the single monitor response from the
// Postman API
type Monitor struct {
	ID             string         `json:"id,omitempty"`
	Name           string         `json:"name,omitempty"`
	UID            string         `json:"uid,omitempty"`
	Owner          string         `json:"owner,omitempty"`
	CollectionUID  string         `json:"collectionUid,omitempty"`
	EnvironmentUID string         `json:"environmentUid,omitempty"`
	Options        MonitorOptions `json:"options,omitempty"`
	Notifications  Notifications  `json:"notifications,omitempty"`
	Distribution   []interface{}  `json:"distribution,omitempty"`
	Schedule       Schedule       `json:"schedule,omitempty"`
}

type monitor struct {
	ID             string         `json:"id,omitempty"`
	Name           string         `json:"name,omitempty"`
	UID            string         `json:"uid,omitempty"`
	Owner          int            `json:"owner,omitempty"`
	CollectionUID  string         `json:"collectionUid,omitempty"`
	EnvironmentUID string         `json:"environmentUid,omitempty"`
	Options        MonitorOptions `json:"options,omitempty"`
	Notifications  Notifications  `json:"notifications,omitempty"`
	Distribution   []interface{}  `json:"distribution,omitempty"`
	Schedule       Schedule       `json:"schedule,omitempty"`
}

// Format returns column headers and values for the resource.
func (r Monitor) Format() ([]string, []interface{}) {
	s := make([]interface{}, 1)
	s[0] = r

	return []string{"UID", "Name"}, s
}

// UnmarshalJSON sets the receiver to a copy of data.
func (r *Monitor) UnmarshalJSON(data []byte) error {
	var m monitor
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	r.ID = m.ID
	r.Name = m.Name
	r.UID = m.UID
	r.Owner = strconv.Itoa(m.Owner)
	r.CollectionUID = m.CollectionUID
	r.EnvironmentUID = m.EnvironmentUID
	r.Options = m.Options
	r.Notifications = m.Notifications
	r.Distribution = m.Distribution
	r.Schedule = m.Schedule

	return nil
}

// MonitorSlice is a slice of Monitor.
type MonitorSlice []*Monitor

// Format returns column headers and values for the resource.
func (r MonitorSlice) Format() ([]string, []interface{}) {
	s := make([]interface{}, len(r))
	for i, v := range r {
		s[i] = v
	}

	return []string{"UID", "Name"}, s
}

// MonitorOptions list options for a monitor.
type MonitorOptions struct {
	StrictSSL       bool `json:"strictSSL,omitempty"`
	FollowRedirects bool `json:"followRedirects,omitempty"`
	RequestTimeout  *int `json:"requestTimeout,omitempty"`
	RequestDelay    int  `json:"requestDelay,omitempty"`
}

// OnError represents a communication mechanism for errors.
type OnError struct {
	Email string `json:"email,omitempty"`
}

// OnFailure represents a communication mechanism for failures.
type OnFailure struct {
	Email string `json:"email,omitempty"`
}

// Notifications represents a communication structure for notifications.
type Notifications struct {
	OnError   []OnError   `json:"onError,omitempty"`
	OnFailure []OnFailure `json:"onFailure,omitempty"`
}

// Schedule represents when the monitor is scheduled to run.
type Schedule struct {
	Cron     string    `json:"cron,omitempty"`
	Timezone string    `json:"timezone,omitempty"`
	NextRun  time.Time `json:"nextRun,omitempty"`
}
