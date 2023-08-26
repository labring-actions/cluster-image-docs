/*
Copyright 2023 cuisongliu@qq.com.

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

import (
	"sort"
	"strings"
)

type ByLatest []string

func (vs ByLatest) Len() int      { return len(vs) }
func (vs ByLatest) Swap(i, j int) { vs[i], vs[j] = vs[j], vs[i] }
func (vs ByLatest) Less(i, j int) bool {
	containsLatestI := strings.Contains(vs[i], "latest")
	containsLatestJ := strings.Contains(vs[j], "latest")

	if containsLatestI && !containsLatestJ {
		return true
	}
	if !containsLatestI && containsLatestJ {
		return false
	}
	return i < j
}

// Sort sorts a list of semantic version strings using ByRootfs.
func Sort(list []string) {
	sort.Sort(ByLatest(list))
}

type ByImageInfo []ImageInfo

func (vs ByImageInfo) Len() int      { return len(vs) }
func (vs ByImageInfo) Swap(i, j int) { vs[i], vs[j] = vs[j], vs[i] }
func (vs ByImageInfo) Less(i, j int) bool {
	return vs[i].Name < vs[j].Name
}

func SortByImageInfo(list []ImageInfo) {
	sort.Sort(ByImageInfo(list))
}
