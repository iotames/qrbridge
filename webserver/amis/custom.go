package amis

import "fmt"

// Custom 组件，用于嵌入自定义 HTML/JavaScript
type Custom struct {
	BaseComponent
	Id      string `json:"id,omitempty"`
	Html    string `json:"html"`              // 静态 HTML 内容
	OnMount string `json:"onMount,omitempty"` // 组件渲染后执行的 JavaScript 代码
}

// NewCustom 创建 Custom 组件，必须提供 HTML 内容
func NewCustom(html string) *Custom {
	return &Custom{
		BaseComponent: BaseComponent{Type: "custom"},
		Html:          html,
	}
}

// SetID 设置组件 ID
func (c *Custom) SetID(id string) *Custom {
	c.Id = id
	return c
}

// SetOnMount 设置组件挂载后执行的脚本
func (c *Custom) SetOnMount(script string) *Custom {
	c.OnMount = script
	return c
}

// Map 转换为 AMIS 配置
func (c *Custom) Map() map[string]any {
	m := map[string]any{
		"type": c.Type,
		"html": c.Html,
	}
	if c.Id != "" {
		m["id"] = c.Id
	}
	if c.OnMount != "" {
		m["onMount"] = c.OnMount
	}
	return m
}

// BuildWebSocketCustom 是一个辅助函数，用于快速构建一个连接 WebSocket 并显示输出的 custom 组件
// wsURL: WebSocket 服务端地址
// outputElementID: 用于显示输出的 HTML 元素 ID
func BuildWebSocketCustom(wsURL, outputElementID string) *Custom {
	html := fmt.Sprintf(`<div id="%s" style="background:#1e1e1e; color:#d4d4d4; padding:10px; font-family: 'Courier New', monospace; height: 800px; overflow: auto; white-space: pre-wrap; border-radius: 4px;"></div>`, outputElementID)

	// 注意：onMount 脚本中的字符串需要进行转义，以适应 JSON
	onMountScript := fmt.Sprintf(`
        (function() {
            var ws = new WebSocket('%s');
            var outputDiv = document.getElementById('%s');
            if (!outputDiv) {
                console.error('输出元素未找到');
                return;
            }
            ws.binaryType = 'arraybuffer';

            ws.onopen = function() {
                outputDiv.innerText += '[已连接到服务器]\n';
            };

            ws.onmessage = function(event) {
                if (typeof event.data === 'string') {
                    outputDiv.innerText += event.data;
                } else {
                    var decoder = new TextDecoder('utf-8');
                    var text = decoder.decode(event.data);
                    outputDiv.innerText += text;
                }
                outputDiv.scrollTop = outputDiv.scrollHeight;
            };

            ws.onerror = function(error) {
                outputDiv.innerText += '[连接错误: ' + error.type + ']\n';
                console.error('WebSocket Error:', error);
            };

            ws.onclose = function(event) {
                var reason = event.reason ? event.reason : '无原因';
                outputDiv.innerText += '[连接已关闭，代码: ' + event.code + ', 原因: ' + reason + ']\n';
            };
        })();
    `, wsURL, outputElementID)

	return NewCustom(html).SetOnMount(onMountScript)
}
