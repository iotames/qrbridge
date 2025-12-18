## Prompt For AI

我做了一个基于HTTP的API接口微服务，用于把PO订单文件转换成Excel。

HTTP请求示例:

```json
// 正常返回
curl -X POST http://127.0.0.1:8080/api/poimport -H "Content-Type: application/json" -d '{"inputtpl":"A89SP","inputfile":"D:\\Users\\santic\\Downloads\\PO SS26 Main (1)(2).xlsx","outputfile":"poss26main.xlsx"}'
{"code":200,"msg":"success","data":{"inputfile":"D:\\Users\\santic\\Downloads\\PO SS26 Main (1)(2).xlsx","outputfile":"poss26main.xlsx"}}

// curl -X POST http://127.0.0.1:8080/api/poimport -H "Content-Type: application/json" -d '{"inputtpl":"A89SP","inputfile":"https://s1.sendike.com/salesContract/1766020611PO SS26 Main (1)(2).xlsx","outputfile":"poss26main.xlsx"}'


// 异常返回
curl -X POST http://localhost:8080/api/poimport   -H "Content-Type: application/json"   -d '{"inputfile":"D:\\Users\\santic\\Downloads\\PO SS26 Main (1)(2).xlsx","outputfile":"poss26main.xlsx","inputtpl":"Rohnisch9"}'
{"code":400,"msg":"QueryArgsError.请求参数错误:inputtpl仅支持:A89SP","data":{}}
```


SDK调用示例:

```php
    $helper = new SanticErpHelper('http://127.0.0.1:8999');
    try {
        // 最好使用绝对路径
        $result = $helper->convertPoToExcel(
            '',
            'customer-po-1.xlsx',
            'export-1.xlsx'
        );
        if ($result['code'] === 200) {
            echo "✓ 转换成功: " . $result['msg'] . "\n";
        } else {
            echo "✗ 转换失败: [" . $result['code'] . "] " . $result['msg'] . "\n";
        }
    } catch (Exception $e) {
        echo "✗ 异常: " . $e->getMessage() . "\n";
    }
```
