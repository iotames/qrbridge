CC=go
APP_VERSION=v1.6.4
# 中文乱码，在CFLAGS添加-fexec-charset=UTF-8选项

# # For Windows:
# # 生成图标和版本信息
# go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest

# 根据操作系统设置目标文件名和链接库
ifeq ($(OS),Windows_NT)
	BUILD_FILE_NAME=POTool.exe
	BUILD_TIME=%date:~0,4%-%date:~5,2%-%date:~8,2%_%time:~0,2%_%time:~3,2%
	RELEASE_FILE="release\$(BUILD_FILE_NAME)"
# 	copy resource\amis\rest.js release\resource\amis
	GO_BUILD= goversioninfo versioninfo.json && go build -v
	GO_BUILD_ARGS=-trimpath -ldflags "-X 'main.BuildTime=$(BUILD_TIME)' -X 'main.Version=$(APP_VERSION)' -X 'main.DbFlag=false' "
	COPY=copy
	DIRSEP=\\
	RM=del /Q
	MKDIR=mkdir
else
	BUILD_FILE_NAME=qrbridge
	RELEASE_FILE="release/$(BUILD_FILE_NAME)"
	GO_BUILD=CGO_ENABLED=0 go build -v
	BUILD_TIME=$(shell date +%Y-%m-%d_%H_%M)
	GO_BUILD_ARGS=-trimpath -ldflags "-X 'main.BuildTime=$(BUILD_TIME)' -X 'main.Version=$(APP_VERSION)' -X 'main.DbFlag=true' "
	COPY=cp -rf
	DIRSEP=/
	RM=rm -rf
	MKDIR=mkdir -p
endif

SRC_AMIS_DIR=resource$(DIRSEP)amis
RELEASE_AMIS_DIR=release$(DIRSEP)resource$(DIRSEP)amis
ZIP_NAME = POTool_$(APP_VERSION).zip

# 执行make编译生成可执行文件
build: $(BUILD_FILE_NAME)

# 新增 initdir BUILD_FILE_NAME 来设置代码页
initdir:
ifeq ($(OS),Windows_NT)
	@echo "Begin Make For Windows app version $(APP_VERSION)"
# 	chcp 65001 >nul
else
	echo "Begin Make For Linux app version $(APP_VERSION)"
endif
	-$(MKDIR) release$(DIRSEP)tpl
	-$(MKDIR) release$(DIRSEP)resource$(DIRSEP)amis

# 执行make release时，编译生成最终的可执行文件，并复制到release目录下
release: initdir $(BUILD_FILE_NAME)
	$(COPY) $(BUILD_FILE_NAME) $(RELEASE_FILE)
	$(COPY) tpl release$(DIRSEP)tpl
	$(COPY) $(SRC_AMIS_DIR)$(DIRSEP)helper.css $(RELEASE_AMIS_DIR)$(DIRSEP)helper.css
	$(COPY) $(SRC_AMIS_DIR)$(DIRSEP)iconfont.css $(RELEASE_AMIS_DIR)$(DIRSEP)iconfont.css
	$(COPY) $(SRC_AMIS_DIR)$(DIRSEP)rest.js $(RELEASE_AMIS_DIR)$(DIRSEP)rest.js
	$(COPY) $(SRC_AMIS_DIR)$(DIRSEP)sdk.css $(RELEASE_AMIS_DIR)$(DIRSEP)sdk.css
	$(COPY) $(SRC_AMIS_DIR)$(DIRSEP)sdk.js $(RELEASE_AMIS_DIR)$(DIRSEP)sdk.js

zip: release
ifeq ($(OS),Windows_NT)
	powershell -Command "Compress-Archive -Path 'release' -DestinationPath '$(ZIP_NAME)' -Force"
else
	zip -r $(ZIP_NAME) release
endif
	@echo "ZIP Generate Done: $(ZIP_NAME)"

debug:
	@echo $(SRC_AMIS_DIR)
	@echo $(RELEASE_AMIS_DIR)
	@echo $(BUILD_TIME)
# 	$(COPY) $(SRC_AMIS_DIR)$(DIRSEP)helper.css $(RELEASE_AMIS_DIR)$(DIRSEP)helper.css

# 编译生成可执行文件
$(BUILD_FILE_NAME):
	$(GO_BUILD) -o $(BUILD_FILE_NAME) $(GO_BUILD_ARGS) .

# 
# $(RELEASE_FILE): $(BUILD_FILE_NAME)
# 	$(COPY) $(BUILD_FILE_NAME) $(RELEASE_FILE)
# 	$(COPY) tpl release$(DIRSEP)tpl

# 用于清理编译生成的中间文件和可执行文件。
clean:
	-$(RM) $(BUILD_FILE_NAME)
	-$(RM) $(RELEASE_FILE)
	-$(RM) release$(DIRSEP)tpl
	-$(RM) $(ZIP_NAME)

# 运行编译好的程序。
run: $(BUILD_FILE_NAME)
ifeq ($(OS),Windows_NT)
	$(BUILD_FILE_NAME)
else
	./$(BUILD_FILE_NAME)
endif

deploy: $(BUILD_FILE_NAME)
	scp $(BUILD_FILE_NAME) santicerp:/home/localapi/$(BUILD_FILE_NAME).new
	ssh santicerp "cd /home/localapi/; \
		./$(BUILD_FILE_NAME) --version; \
	    systemctl status pgprice; \
		mv /home/localapi/$(BUILD_FILE_NAME) /home/localapi/$(BUILD_FILE_NAME).$(BUILD_TIME); \
		mv /home/localapi/$(BUILD_FILE_NAME).new /home/localapi/$(BUILD_FILE_NAME); \
		systemctl restart pgprice; \
		systemctl status pgprice; \
		./$(BUILD_FILE_NAME) --version;"

# .PHONY 声明后面的目标（如 release、clean、run）不是实际的文件名，而是“伪目标”。保证release、clean、run 始终被执行
.PHONY: initdir build release clean run zip deploy debug
