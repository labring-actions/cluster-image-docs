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

package dockerhub

import (
	"github.com/cuisongliu/logger"
	"testing"
)

func TestFetchDockerHubAllVersion(t *testing.T) {
	logger.Cfg(true, false)
	err := autoRemoveGenerator("../docs")
	if err != nil {
		t.Error(err)
		return
	}
	//got, err := fetchDockerHubAllRepo()
	//if err != nil {
	//	t.Error(err)
	//	return
	//}
	//t.Log("get docker hub all repo success")
	//for k, v := range got {
	//	err = generatorSyncFile("../docs", k, v)
	//	if err != nil {
	//		t.Errorf("generatorSyncFile %s error %s", k, err.Error())
	//		continue
	//	}
	//}
}
