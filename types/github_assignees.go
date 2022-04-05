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

import "sync"

// Assignees has GitHub users
type Assignees struct {
	once sync.Once
	// listMap contains GitHub users in Key
	listMap map[string]struct{}
}

// RemoveFromList removes users from the list
func (a *Assignees) RemoveFromList(names ...string) {
	for _, name := range names {
		delete(a.listMap, name)
	}
}

// Add adds users into the list
func (a *Assignees) Add(assignees ...string) {
	if len(assignees) == 0 {
		return
	}

	a.once.Do(func() {
		a.listMap = make(map[string]struct{})
	})

	for _, assignee := range assignees {
		a.listMap[assignee] = struct{}{}
	}
}

// List returns the list of users
func (a *Assignees) List() []string {
	size := len(a.listMap)
	if size == 0 {
		return nil
	}

	list := make([]string, 0, size)
	for key := range a.listMap {
		list = append(list, key)
	}
	return list
}
