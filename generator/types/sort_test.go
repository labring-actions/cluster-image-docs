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

import "testing"

func TestSort(t *testing.T) {
	type args struct {
		list []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				list: []string{"v1.1.1", "latest"},
			},
		},
		{
			name: "test1",
			args: args{
				list: []string{"v1.1.1", "latest", "v1.2.3"},
			},
		},
		{
			name: "test1",
			args: args{
				list: []string{"latest", "v1.1.1", "v1.2.3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Sort(tt.args.list)
			t.Log(tt.args.list)
		})
	}
}
