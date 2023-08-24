# sync-aliyun
auto sync images to aliyun

set vars for action 

- A_REGISTRY_USERNAME: login user cuisongliu@1186388767322904
- A_REGISTRY_NAME: registry addr  registry.cn-hongkong.aliyuncs.com
- A_REGISTRY_REPOSITORY: registry repo cuisongliu-labring

set secret for action:

- A_REGISTRY_TOKEN: registry password
- GH_PAT: auto commit pr to github

经过测试sync两个600M镜像大概50s

修改 generator的`reviewers: cuisongliu` 改成对应的reviewer
