#当前目录
CURRENT_DIR := $(shell pwd)

#######################参数#######################
#获取分支名 默认dev分支
ifeq ($(B),)
    A := $(info 缺少分支参数B. 不更新代码, 如果更新代码后编译请使用"make B=分支名字")
else
    A := $(info 本次编译分支为 $(B), 即将更新代码)
    A := $(shell cd .. && git checkout $(B) && git pull --rebase -f && cd $(CURRENT_DIR))
    A := $(info 代码更新完成)
endif

#######################版本#######################
#代码版本
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_TIME := $(shell date "+%F %T")
COMMIT_SHA1 := $(shell git rev-parse HEAD)
COMMIT_TIME := $(shell git log -1 --format="%ci")
COMMIT_LOG := $(shell git log -1 --format="%s")

#编译参数
BUILD_LDFLAG="-X 'laiya/share/version.BuildTime=${BUILD_TIME}' -X 'laiya/share/version.CommitID=${COMMIT_SHA1}' -X 'laiya/share/version.CommitTime=${COMMIT_TIME}' -X 'laiya/share/version.Branch=${BRANCH}' "
BUILD_FLAGS=-ldflags $(BUILD_LDFLAG)
BUILD_CMD=go build $(BUILD_FLAGS)

#proto可执行文件位置
export PATH := $(CURRENT_DIR)/tools/:$(PATH)

all: gate home web hall group wechat

version:
	@echo CURRENT_DIR:$(CURRENT_DIR)
	@echo BRANCH:$(BRANCH)
	@echo COMMIT_SHA1:$(COMMIT_SHA1)
	@echo COMMIT_TIME:$(COMMIT_TIME)
	@echo COMMIT_LOG:$(COMMIT_LOG)
	@echo BUILD_TIME:$(BUILD_TIME)
gate: version
	$(BUILD_CMD) -o ./bin/gate cmd/gate/main.go
home: version
	$(BUILD_CMD) -o ./bin/home cmd/home/main.go
web: version
	$(BUILD_CMD) -o ./bin/web cmd/web/main.go
hall: version
	$(BUILD_CMD) -o ./bin/hall cmd/hall/main.go
group: version
	$(BUILD_CMD) -o ./bin/group cmd/group/main.go
wechat: version
	$(BUILD_CMD) -o ./bin/wechat cmd/wechat/main.go

clean:
	cd bin && ./run.sh clean
auto:
ifdef T
	@echo $(shell echo "本次为单进程更新:$(T)")
	cd bin && ./run.sh stop $(T) && cd .. && make $(T) && cd bin && ./run.sh restart $(T)
else
	@echo $(shell echo "本次为全进程更新")
	cd bin && ./run.sh stop && ./run.sh clean && cd .. && make all && cd bin && ./run.sh restart
endif

proto: version
	protoc -I=./proto_file --go-new_out=paths=source_relative:./proto ./proto_file/outer/*.proto
	protoc -I=./proto_file --go-new_out=paths=source_relative:./proto ./proto_file/code/*.proto
	protoc -I=./proto_file --go-new_out=paths=source_relative:./proto ./proto_file/common/*.proto
	protoc -I=./proto_file --go-new_out=paths=source_relative:./proto ./proto_file/inner/*.proto
	protoc -I=./proto_file --go-new_out=paths=source_relative:./proto ./proto_file/gm/*.proto

dart: version
	@echo $(shell pwd)
	protoc -I=./proto_file --dart_out=./proto_file/gen-dart ./proto_file/outer/*.proto
	protoc -I=./proto_file --dart_out=./proto_file/gen-dart ./proto_file/code/*.proto
	protoc -I=./proto_file --dart_out=./proto_file/gen-dart ./proto_file/common/*.proto
	protoc -I=./proto_file --dart_out=./proto_file/gen-dart ./proto_file/google/protobuf/*.proto
	#cd ./proto_file/gen-dart && find . -type f -name "*.pb.dart" -exec sed -i 's|google/protobuf/any.pb.dart|./any.pb.dart|g' {} +

.PHONY: proto gate home group hall we web wechat