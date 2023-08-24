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
	"bufio"
	"github.com/cuisongliu/logger"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func count(dir string) int64 {
	c := int64(0)
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 如果不是目录，并且文件名以指定的前缀开始，就删除它
		if !info.IsDir() && strings.HasPrefix(info.Name(), prefix) {
			//read file
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "#count") {
					// 提取数量
					parts := strings.Split(line, ":")
					if len(parts) == 2 {
						countStr := strings.TrimSpace(parts[1])
						newCount, err := strconv.ParseInt(countStr, 10, 64)
						if err != nil {
							logger.Error("auto count count.json files error", err)
							continue
						}
						c += newCount
					}
				}
			}
		}

		return nil
	}); err != nil {
		return 0
	}
	logger.Info("auto count count.json files success")
	return c
}
