## 介绍

QrBridge 是一个二维码调度中转站，可以将多个二维码调度到不同的网站或服务器上。
通过加密算法，对二维码网址的请求参数进行加密，从而隐藏数据库ID。


## 运行

`CMD` 命令窗口 或 `Shell` 运行：

```
# for windows
qrbridge.exe

# for linux
./qrbridge
```

首次运行，会自动生成 `.env` 配置文件。手动修改配置文件后，重新运行即可。


## 配置

程序首次运行时，会自动生成 `.env` 配置文件。示例如下：

- `TO_BASE_URL`: 要跳转过去的目标服务器。此处填写URL前缀。如：`https://hello.yoursite.com`
- `ENCRYPT_MULTIPLE`：加密倍数
- `ENCRYPT_ADD`：加密增量

目前数据库只支持 `postgres`, 修改数据库配置，加密参数，和 `TO_BASE_URL` 即可使用。


## 参数加解密规则

假设二维码对应的网页，需要2个参数: `m` 和 `mid`，才可以进行正常的业务逻辑处理和数据库查询。
格式为：`m=module_name&mid=128`。其中，`mid` 的值，是需要隐藏的数据库ID字段，必须为正整数。

加密规则：

1. 在项目的 `.env` 文件定义 `ENCRYPT_MULTIPLE` 和 `ENCRYPT_ADD`。 数据类型都为正整数。
2. `mid` 的值，先乘以 `ENCRYPT_MULTIPLE`，再加上 `ENCRYPT_ADD`。如：`128 * 10 + 5 = 1285`。
3. 对刚处理过的整个参数，进行 `urlbase64` 转码: urlbase64_encode("m=module_name&mid=1285")。最后去掉末尾的 `=` 符号。

例：`m=module_name&mid=128` 加密结果为：`bT1tb2R1bGVfbmFtZSZtaWQ9MjAwNTA3MDg`

整个过程中，除 `mid` 之外的参数，不进行额外处理。


## 业务逻辑

1. 调试二维码参数。生成隐藏数据库ID后，需要展示给用户的二维码参数 `code`。
2. 访问二维码地址。通过前面生成的二维码参数 `code`。
3. 根据 `code` 解析结果，进行针对性的业务处理。例如，判断是否跳转到第三方网页。

假设 `ENCRYPT_ADD` 参数值为 `5006`。则 `调试二维码参数` 的路由地址为：`/codetest5006`
通过调试地址，生成二维码参数，可以验证加密结果是否正确。

```
http://127.0.0.1:8080/codetest5006?m=module_name&mid=128

# 返回结果：{"code":200,"msg":"code=bT1tb2R1bGVfbmFtZSZtaWQ9MjAwNTA3MDg","data":{}}
```

`访问二维码地址` 如下所示：

```
https://127.0.0.1:8080/qrcode?code=bT1tb2R1bGVfbmFtZSZtaWQ9MjAwNTA3MDg
```

## 访问结果记录

- `st_qrcode_list`: 记录每个二维码的基本信息，比如访问次数，code参数，解析是否成功等。
- `st_qrcode_query_log`: 记录每个二维码连接的访问历史。包括请求IP，请求头，访问时间等。
