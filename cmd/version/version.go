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

package version

import (
	"fmt"

	"github.com/fatih/color"
)

// Version represents main version number being run right now.
var Version = "v0.1.0"

// ReleasePhase represents pre-release marker for the version. If this is an empty string,
// then the release is a final release. Otherwise this is a pre-release
// version e.g. "dev", "alpha", etc.
var ReleasePhase = ""

// String prints the version of the spini CLI.
func String() string {
	info := color.New(color.Bold, color.FgGreen).SprintFunc()

	if ReleasePhase != "" {
		return fmt.Sprint(info(fmt.Sprintf("%s-%s", Version, ReleasePhase)))
	}
	return fmt.Sprint(info(Version))
}
