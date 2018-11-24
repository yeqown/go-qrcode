## go-qrcode

### 示例

link to [CODE](./example/main.go)
```go
package main

import (
	"fmt"

	qrcode "github.com/yeqown/go-qrcode"
)

func main() {
	// 配置文件，默认在repo的根路径下
	qrcode.SetVersionCfgFile("../versionCfg.json")

	qrc, err := qrcode.New("https://github.com/yeqown/go-qrcode")
	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err)
	}

	// save file
	if err := qrc.Save("../testdata/repo-qrcode.jpeg"); err != nil {
		fmt.Printf("could not save image: %v", err)
	}
}
```

<img src="./testdata/repo-qrcode.jpeg" width="200px"/>

### QRCode原理

#### 1 数据分析（data analysis）：

分析输入数据，根据数据决定要使用的QR码版本、容错级别和编码模式。低版本的QR码无法编码过长的数据，含有非数字字母字符的数据要使用扩展字符编码模式。

#### 2 编码数据（data encoding）：

根据选择的编码模式，将输入的字符串转换成比特流，插入模式标识码（mode indicator）和终止标识符（terminator），把比特流切分成八比特的字节，加入填充字节来满足标准的数据字码数要求。

#### 3 计算容错码（error correction coding）：

对步骤二产生的比特流计算容错码，附在比特流之后。高版本的编码方式可能需要将数据流切分成块（block）再分别进行容错码计算。

#### 4 组织数据（structure final message）：

根据结构图把步骤三得到的有容错的数据流切分，准备填充。

#### 5 填充（module placement in matrix）：

把数据和功能性图样根据标准填充到矩阵中。

#### 6 应用数据掩码（data masking）：

应用标准中的八个数据掩码来变换编码区域的数据，选择最优的掩码应用。讲到再展开。

#### 7 填充格式和版本信息（format and version information）：

计算格式和版本信息填入矩阵，完成QR码。

### 参考文献

* [QR Code tutori](https://www.thonky.com/qr-code-tutorial/)
* [QRCode Wiki](https://en.wikipedia.org/wiki/QR_code)
* [二维码详解（QR Code）](https://zhuanlan.zhihu.com/p/21463650)
* [数据编码](https://zhuanlan.zhihu.com/p/25432676)