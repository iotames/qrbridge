CC=go
# 中文乱码，在CFLAGS添加-fexec-charset=UTF-8选项
# SRCDIR=src
# OBJDIR=release

# # For Windows:
# # 生成图标和版本信息
# go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest

# 根据操作系统设置目标文件名和链接库
ifeq ($(OS),Windows_NT)
	BUILD_FILE_NAME=POTool.exe
	RELEASE_FILE_PATH="release\$(BUILD_FILE_NAME)"
	GO_BUILD= goversioninfo versioninfo.json && go build -v
	GO_BUILD_ARGS=-trimpath -ldflags "-X 'main.BuildTime=%date:~0,4%-%date:~5,2%-%date:~8,2%_%time:~0,2%:%time:~3,2%' -X 'main.Version=v1.1.1' -X 'main.DbFlag=false' "
	COPY=copy
	RM=del /Q
	DIRSEP="\"
	RMREDIR=2>nul
else
	BUILD_FILE_NAME=qrbridge
	RELEASE_FILE_PATH="release/$(BUILD_FILE_NAME)"
	GO_BUILD=CGO_ENABLED=0 go build -v
	BUILD_TIME=$(shell date +%Y-%m-%d_%H:%M)
	GO_BUILD_ARGS=-trimpath -ldflags "-X 'main.BuildTime=$(BUILD_TIME)' -X 'main.Version=v1.1.1' -X 'main.DbFlag=true' "
	COPY=cp
	DIRSEP="/"
	RM=rm -f
	RMREDIR=
endif

# 新增 init BUILD_FILE_NAME 来设置代码页
init:
ifeq ($(OS),Windows_NT)
	echo "Begin Make For Windows..............."
# 	chcp 65001 >nul
endif

build: init $(BUILD_FILE_NAME)

# 执行make或make all时，编译生成最终的可执行文件
all: init $(RELEASE_FILE_PATH)

$(BUILD_FILE_NAME):
	$(GO_BUILD) -o $(BUILD_FILE_NAME) $(GO_BUILD_ARGS) .

$(RELEASE_FILE_PATH): $(BUILD_FILE_NAME)
	$(COPY) $(BUILD_FILE_NAME) $(RELEASE_FILE_PATH)

# 用于清理编译生成的中间文件和可执行文件。执行 make clean 会删除 obj 目录下的 .o 文件和最终的可执行文件
clean: init
	-$(RM) $(RELEASE_FILE_PATH)
	-$(RM) $(BUILD_FILE_NAME)
# 	-$(RM) $(OBJDIR)$(DIRSEP)$(BUILD_FILE_NAME) $(RMREDIR)
# 	-$(RM) $(BUILD_FILE_NAME) $(RMREDIR)

# 用于运行编译好的程序。
# Windows命令窗口启动后，先执行chcp 65001，再运行程序。避免中文乱码问题
run: init $(BUILD_FILE_NAME)
ifeq ($(OS),Windows_NT)
# 	chcp 65001 >nul
	$(BUILD_FILE_NAME)
else
	./$(BUILD_FILE_NAME)
endif

# 声明后面的目标（如 all、clean、run）不是实际的文件名，而是“伪目标”。
# .PHONY保证all、clean、run 始终被执行
.PHONY: init build all clean run
