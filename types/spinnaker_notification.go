/*
Copyright © 2022 Yevhen Lebid ealebed@gmail.com

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

package types

// NotificationMessage represents text of a stage state
type NotificationMessage struct {
	Text string `json:"text"`
}

// Notifications represents full notifications config
type Notification struct {
	ID      string                         `json:"id,omitempty"`
	Address string                         `json:"address"`
	Level   string                         `json:"level"`
	Message map[string]NotificationMessage `json:"message"`
	Type    string                         `json:"type"`
	When    []string                       `json:"when"`
}

// NewNotification new notification
func NewNotification() *Notification {
	return &Notification{}
}
