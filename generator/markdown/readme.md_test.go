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

package markdown

import "testing"

func Test_readme_Generator(t *testing.T) {
	type fields struct {
		addr string
	}
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "default",
			fields: fields{
				addr: "docker.io",
			},
			args: args{
				dir: "../../tmp",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &readme{
				addr: tt.fields.addr,
			}
			if err := r.Generator(tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("Generator() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
