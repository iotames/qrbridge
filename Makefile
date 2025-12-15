CC=go
APP_VERSION=v1.1.2
# 中文乱码，在CFLAGS添加-fexec-charset=UTF-8选项

# # For Windows:
# # 生成图标和版本信息
# go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest

# 根据操作系统设置目标文件名和链接库
ifeq ($(OS),Windows_NT)
	BUILD_FILE_NAME=POTool.exe
	RELEASE_FILE="release\$(BUILD_FILE_NAME)"
	GO_BUILD= goversioninfo versioninfo.json && go build -v
	GO_BUILD_ARGS=-trimpath -ldflags "-X 'main.BuildTime=%date:~0,4%-%date:~5,2%-%date:~8,2%_%time:~0,2%:%time:~3,2%' -X 'main.Version=$(APP_VERSION)' -X 'main.DbFlag=false' "
	COPY=copy
	RM=del /Q
	DIRSEP="\"
else
	BUILD_FILE_NAME=qrbridge
	RELEASE_FILE="release/$(BUILD_FILE_NAME)"
	GO_BUILD=CGO_ENABLED=0 go build -v
	BUILD_TIME=$(shell date +%Y-%m-%d_%H:%M)
	GO_BUILD_ARGS=-trimpath -ldflags "-X 'main.BuildTime=$(BUILD_TIME)' -X 'main.Version=$(APP_VERSION)' -X 'main.DbFlag=true' "
	COPY=cp
	DIRSEP="/"
	RM=rm -f
endif

# 新增 init BUILD_FILE_NAME 来设置代码页
init:
ifeq ($(OS),Windows_NT)
	echo "Begin Make For Windows app version $(APP_VERSION)"
# 	chcp 65001 >nul
else
	echo "Begin Make For Linux app version $(APP_VERSION)"
endif

# 执行make编译生成可执行文件
build: init $(BUILD_FILE_NAME)

# 执行make all时，编译生成最终的可执行文件
all: init $(RELEASE_FILE)

# 编译生成可执行文件
$(BUILD_FILE_NAME):
	$(GO_BUILD) -o $(BUILD_FILE_NAME) $(GO_BUILD_ARGS) .

# 拷贝可执行文件到指定目录
$(RELEASE_FILE): $(BUILD_FILE_NAME)
	$(COPY) $(BUILD_FILE_NAME) $(RELEASE_FILE)

# 用于清理编译生成的中间文件和可执行文件。
clean: init
	-$(RM) $(RELEASE_FILE)
	-$(RM) $(BUILD_FILE_NAME)

# 运行编译好的程序。
run: init $(BUILD_FILE_NAME)
ifeq ($(OS),Windows_NT)
	$(BUILD_FILE_NAME)
else
	./$(BUILD_FILE_NAME)
endif

# .PHONY 声明后面的目标（如 all、clean、run）不是实际的文件名，而是“伪目标”。保证all、clean、run 始终被执行
.PHONY: init build all clean run
