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

package generator

import (
	"github.com/cuisongliu/logger"
	"os"
	"testing"
)

func TestFetchDockerHubAllVersion(t *testing.T) {
	logger.Cfg(true, false)
	syncDir := "../docs/docker"
	err := os.MkdirAll(syncDir, 0700|0055)
	if err != nil {
		t.Error(err)
		return
	}
	err = autoRemoveGenerator(syncDir)
	if err != nil {
		t.Error(err)
		return
	}
	got, err := fetchDockerHubAllRepo("docker.io/labring")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("get docker hub all repo success")
	if got == nil {
		t.Error("got is nil")
		return
	}
	t.Log("rootfs infos")
	for _, v := range got.Rootfs {
		t.Logf("name: %s, tags: %+v", v.Name, v.Tags)
	}
	t.Log("sealos infos")
	for _, v := range got.Sealos {
		t.Logf("name: %s, tags: %+v", v.Name, v.Tags)
	}
	t.Log("laf infos")
	for _, v := range got.Laf {
		t.Logf("name: %s, tags: %+v", v.Name, v.Tags)
	}
	t.Log("app infos")
	for _, v := range got.Apps {
		t.Logf("name: %s, tags: %+v", v.Name, v.Tags)
	}

}

func TestFetchAliyunAllVersion(t *testing.T) {
	logger.Cfg(true, false)
	syncDir := "../docs/aliyun-hk"
	err := os.MkdirAll(syncDir, 0700|0055)
	if err != nil {
		t.Error(err)
		return
	}
	err = autoRemoveGenerator(syncDir)
	if err != nil {
		t.Error(err)
		return
	}
	got, err := fetchDockerHubAllRepo("registry.cn-hongkong.aliyuncs.com/labring")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("get docker hub all repo success")
	if got == nil {
		t.Error("got is nil")
		return
	}
	t.Log("rootfs infos")
	for _, v := range got.Rootfs {
		t.Logf("name: %s, tags: %+v", v.Name, v.Tags)
	}
	t.Log("sealos infos")
	for _, v := range got.Sealos {
		t.Logf("name: %s, tags: %+v", v.Name, v.Tags)
	}
	t.Log("laf infos")
	for _, v := range got.Laf {
		t.Logf("name: %s, tags: %+v", v.Name, v.Tags)
	}
	t.Log("app infos")
	for _, v := range got.Apps {
		t.Logf("name: %s, tags: %+v", v.Name, v.Tags)
	}

}
