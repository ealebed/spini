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

// StageExecutionWindowJitter represents random jitter to add to execution window
type StageExecutionWindowJitter struct {
	Enabled    bool `json:"enabled"`
	MaxDelay   int  `json:"maxDelay"`
	MinDelay   int  `json:"minDelay"`
	SkipManual bool `json:"skipManual"`
}

// StageExecutionWindowWhitelist represents hours allowed to perform deploy
type StageExecutionWindowWhitelist struct {
	EndHour   int `json:"endHour"`
	EndMin    int `json:"endMin"`
	StartHour int `json:"startHour"`
	StartMin  int `json:"startMin"`
}

// StageExecutionWindow represents time slot, when allowed to execute pipeline stage
type StageExecutionWindow struct {
	Days      []int                             `json:"days"`
	Jitter    *StageExecutionWindowJitter       `json:"jitter,omitempty"`
	Whitelist *[]*StageExecutionWindowWhitelist `json:"whitelist"`
}

// defaultExecutionWindowWhitelist return StageExecutionWindowWhitelist object with default values for execution window whitelist
func defaultExecutionWindowWhitelist() *StageExecutionWindowWhitelist {
	return &StageExecutionWindowWhitelist{
		EndHour:   11,
		EndMin:    0,
		StartHour: 7,
		StartMin:  0,
	}
}

// defaultStageExecutionWindow return StageExecutionWindow object with default values
func defaultStageExecutionWindow() *StageExecutionWindow {
	return &StageExecutionWindow{
		Days:      []int{2, 3, 4, 5, 6},
		Whitelist: &[]*StageExecutionWindowWhitelist{defaultExecutionWindowWhitelist()},
	}
}
