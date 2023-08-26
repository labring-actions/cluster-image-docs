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
	"context"
	"github.com/cuisongliu/logger"
	"github.com/labring-actions/cluster-image-docs/generator/markdown"
	"github.com/labring-actions/cluster-image-docs/generator/types"
	"golang.org/x/sync/errgroup"
	"os"
	"path/filepath"
)

func Do() {
	logger.Cfg(true, false)
	syncDir, _ := os.LookupEnv("SYNC_DIR")
	if syncDir == "" {
		logger.Fatal("SYNC_DIR is empty")
		return
	}
	logger.Info("using syncDir %s", syncDir)
	err := os.MkdirAll(syncDir, 0700|0055)
	if err != nil {
		logger.Fatal("mkdir failed: %v", err)
		return
	}

	syncHub, _ := os.LookupEnv("SYNC_HUB")
	if syncHub == "" {
		syncHub = "docker.io/labring"
	}

	err = autoRemoveGenerator(syncDir)
	if err != nil {
		logger.Fatal("autoRemoveGenerator sync config error %s", err.Error())
		return
	}
	got, err := fetchDockerHubAllRepo(syncHub)
	if err != nil {
		logger.Fatal("fetchDockerHubAllRepo error %s", err.Error())
		return
	}
	logger.Info("get docker hub all repo success")
	goRunData := make(map[markdown.Type][]types.ImageInfo)
	if got != nil {
		goRunData[markdown.Rootfs] = got.Rootfs
		goRunData[markdown.Sealos] = got.Sealos
		goRunData[markdown.Laf] = got.Laf
		goRunData[markdown.Apps] = got.Apps
	}

	if err = markdown.NewReadme().Generator(syncDir); err != nil {
		logger.Fatal("markdown.NewReadme().Generator error %s", err.Error())
		return
	}

	g, _ := errgroup.WithContext(context.Background())
	for k, v := range goRunData {
		// Capture the range variables.
		k, v := k, v
		g.Go(func() error {
			if len(v) == 0 {
				return nil
			}
			if err = markdown.New(k, syncHub, v).Generator(syncDir); err != nil {
				return err
			}
			return nil
		})
	}

	// Wait for all goroutines to finish and return the first error.
	if err = g.Wait(); err != nil {
		logger.Fatal(err.Error())
	}
}

func autoRemoveGenerator(dir string) error {
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			os.Remove(path)
		}
		return nil
	}); err != nil {
		return err
	}
	logger.Info("auto remove %s files success", dir)
	return nil
}
