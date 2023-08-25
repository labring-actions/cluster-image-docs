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
	"os"
)

func Do() {
	logger.Cfg(true, false)
	syncDir := os.Getenv("SYNC_DIR")
	if syncDir == "" {
		logger.Fatal("SYNC_DIR is empty")
		return
	}
	logger.Info("using syncDir %s", syncDir)
	err := autoRemoveGenerator(syncDir)
	if err != nil {
		logger.Fatal("autoRemoveGenerator sync config error %s", err.Error())
		return
	}
	_, err = fetchDockerHubAllRepo()
	if err != nil {
		logger.Fatal("fetchDockerHubAllRepo error %s", err.Error())
		return
	}
	logger.Info("get docker hub all repo success")
	//g, _ := errgroup.WithContext(context.Background())

	//for k, v := range got {
	//	// Capture the range variables.
	//	k, v := k, v
	//	g.Go(func() error {
	//		if len(v.Repos) == 0 {
	//			return nil
	//		}
	//		if err = generatorSyncFile(syncDir, k, v); err != nil {
	//			return fmt.Errorf("generatorSyncFile %s error: %w", k, err)
	//		}
	//		return nil
	//	})
	//}
	//
	//// Wait for all goroutines to finish and return the first error.
	//if err = g.Wait(); err != nil {
	//	logger.Fatal(err.Error())
	//}
}
