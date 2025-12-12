<?php
/**
 * SanticErpHelper - PO订单文件转换SDK
 * 用于调用PO订单文件转换HTTP API
 */
class SanticErpHelper
{
    private $baseUrl;
    private $timeout;
    
    /**
     * 构造函数
     * @param string $baseUrl API基础地址，默认 http://127.0.0.1:8999
     * @param int $timeout 请求超时时间(秒)，默认30
     */
    public function __construct(string $baseUrl = 'http://127.0.0.1:8999', int $timeout = 30)
    {
        $this->baseUrl = rtrim($baseUrl, '/');
        $this->timeout = $timeout;
    }
    
    /**
     * 转换PO订单文件为Excel
     * @param string $inputtpl 输入模板名称，如 "Rohnisch"
     * @param string $inputfile 输入文件路径
     * @param string $outputfile 输出文件路径
     * @return array API响应，包含 code, msg, data 字段
     * @throws Exception 当HTTP请求失败或响应解析失败时抛出
     */
    public function convertPoToExcel(string $inputtpl, string $inputfile, string $outputfile): array
    {
        $url = $this->baseUrl . '/api/poimport';
        
        $postData = [
            'inputtpl' => $inputtpl,
            'inputfile' => $inputfile,
            'outputfile' => $outputfile
        ];
        
        $ch = curl_init();
        
        curl_setopt_array($ch, [
            CURLOPT_URL => $url,
            CURLOPT_POST => true,
            CURLOPT_POSTFIELDS => json_encode($postData),
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_TIMEOUT => $this->timeout,
            CURLOPT_HTTPHEADER => [
                'Content-Type: application/json',
                'Accept: application/json'
            ]
        ]);
        
        $response = curl_exec($ch);
        $error = curl_error($ch);
        $errno = curl_errno($ch);
        
        curl_close($ch);
        
        if ($response === false) {
            throw new Exception("HTTP请求失败: [$errno] $error");
        }
        
        $result = json_decode($response, true);
        
        if (json_last_error() !== JSON_ERROR_NONE) {
            throw new Exception("响应JSON解析失败: " . json_last_error_msg());
        }
        
        if (!isset($result['code']) || !isset($result['msg'])) {
            throw new Exception("无效的API响应格式");
        }
        
        return $result;
    }
    
    /**
     * 获取基础URL（用于调试）
     * @return string
     */
    public function getBaseUrl(): string
    {
        return $this->baseUrl;
    }
}

// 测试代码
if (basename(__FILE__) === basename($_SERVER['PHP_SELF'])) {
    echo "=== SanticErpHelper 测试 ===\n";
    
    $helper = new SanticErpHelper('http://127.0.0.1:8999');
    
    echo "基础URL: " . $helper->getBaseUrl() . "\n";
    
    // 测试1: 正常调用
    echo "\n测试1: 正常调用\n";
    try {
        // 最好使用绝对路径
        $result = $helper->convertPoToExcel(
            'Rohnisch',
            'D:\\Users\\santic\\Downloads\\PO SS26 Main (1)(2).xlsx',
            'poss26main.xlsx'
        );
        
        if ($result['code'] === 200) {
            echo "✓ 转换成功: " . $result['msg'] . "\n";
        } else {
            echo "✗ 转换失败: [" . $result['code'] . "] " . $result['msg'] . "\n";
        }
    } catch (Exception $e) {
        echo "✗ 异常: " . $e->getMessage() . "\n";
    }
    echo "\n=== 测试完成 ===\n";
}
