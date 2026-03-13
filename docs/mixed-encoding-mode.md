# QR Code Mixed Encoding Mode

## Overview

QR Code supports multiple encoding modes, and the specification allows switching between modes within a single QR code. This enables more efficient encoding of mixed content (e.g., Kanji + numbers + ASCII).

### Benefits

- **Optimal Efficiency**: Use the most compact encoding for each character type
- **Space Savings**: Reduce overall QR Code size for mixed content
- **Flexibility**: Handle diverse content without sacrificing efficiency

## Mode Indicators

Each encoding segment starts with a 4-bit mode indicator:

| Mode | Binary | Hex | Description |
|------|--------|-----|-------------|
| Numeric | `0001` | 0x1 | Digits 0-9 |
| Alphanumeric | `0010` | 0x2 | 0-9, A-Z, and $%*+-./: |
| Byte | `0100` | 0x4 | ISO-8859-1 or UTF-8 |
| Kanji | `1000` | 0x8 | Shift JIS Kanji characters |
| ECI | `0111` | 0x7 | Extended Channel Interpretation |

## Data Structure

Each segment in the data stream has the following structure:

```
┌──────────────┬──────────────┬─────────────────┐
│ Mode         │ Character    │ Data            │
│ Indicator    │ Count        │                 │
│ (4 bits)     │ (varies)     │ (varies)        │
└──────────────┴──────────────┴─────────────────┘
```

### Character Count Indicator Bits

The bit width of the character count indicator depends on the QR version and mode:

| Mode | Version 1-9 | Version 10-26 | Version 27-40 |
|------|-------------|---------------|---------------|
| Numeric | 10 | 12 | 14 |
| Alphanumeric | 9 | 11 | 13 |
| Byte | 8 | 16 | 16 |
| Kanji | 8 | 10 | 12 |

## Encoding Example

### Example: Encoding `日本語123`

#### Without Mixed Mode (Byte Mode Only)

```
Mode: Byte (0100)
Character Count: 9 bytes (UTF-8)
Data: 日本語123 = 9 bytes = 72 bits
Total: 4 + 8 + 72 = 84 bits
```

#### With Mixed Mode (Kanji + Numeric)

```
Segment 1: Kanji Mode
┌────────┬────────┬─────────────────────────────┐
│ Mode   │ Count  │ Data                        │
│ 1000   │ 00000011 │ [日本語 encoded as 13 bits each] │
│ 4 bits │ 8 bits  │ 39 bits (3 × 13)            │
└────────┴────────┴─────────────────────────────┘

Segment 2: Numeric Mode
┌────────┬────────┬─────────────────┐
│ Mode   │ Count  │ Data            │
│ 0001   │ 00000011 │ 123            │
│ 4 bits │ 8 bits  │ 10 bits         │
└────────┴────────┴─────────────────┘

Total: (4 + 8 + 39) + (4 + 8 + 10) = 73 bits
```

**Savings: 84 - 73 = 11 bits (13% more efficient)**

## Encoding Algorithm

### Step 1: Segment the Input

Identify consecutive characters that can be encoded in the same mode:

```go
type Segment struct {
    Mode  encMode
    Chars []rune
}

func segmentInput(raw string) []Segment {
    segments := []Segment{}
    currentMode := detectMode(raw[0])
    currentChars := []rune{}

    for _, r := range raw {
        mode := detectMode(r)
        if mode != currentMode {
            segments = append(segments, Segment{Mode: currentMode, Chars: currentChars})
            currentMode = mode
            currentChars = []rune{}
        }
        currentChars = append(currentChars, r)
    }

    segments = append(segments, Segment{Mode: currentMode, Chars: currentChars})
    return segments
}
```

### Step 2: Encode Each Segment

```go
func encodeSegments(segments []Segment, version int) *binary.Binary {
    dst := binary.New()

    for _, seg := range segments {
        // Mode indicator (4 bits)
        dst.AppendUint32(modeIndicator(seg.Mode), 4)

        // Character count (varies by version and mode)
        charCountBits := getCharCountBits(seg.Mode, version)
        dst.AppendUint32(uint32(len(seg.Chars)), charCountBits)

        // Data
        encodeData(dst, seg.Mode, seg.Chars)
    }

    return dst
}
```

### Step 3: Add Terminator and Padding

After encoding all segments:
1. Add 4-bit terminator (`0000`) if space permits
2. Pad with zeros to byte boundary
3. Add pad bytes (`0xEC`, `0x11` alternating) to fill capacity

## Optimization Considerations

### When to Switch Modes

Switching modes adds overhead (4 bits for mode indicator + character count bits). Consider:

```
Switch cost = 4 (mode) + charCountBits

Example for Version 1:
- Switch to Numeric: 4 + 10 = 14 bits overhead
- Switch to Kanji: 4 + 8 = 12 bits overhead
```

### Mode Selection Strategy

1. **Greedy**: Use the most efficient mode for each character
2. **Look-ahead**: Consider if switching saves bits overall
3. **ECI**: For non-UTF8 character sets

### Example Trade-off

For `A日本語` in Version 1:

| Approach | Bits Used |
|----------|-----------|
| Byte mode only | 4 + 8 + (4×8) = 44 bits |
| Alphanumeric + Kanji | (4+9+5.5) + (4+8+39) = 70.5 bits |

**Result**: Don't switch for just 1 character - overhead exceeds savings!

## Current Implementation Status

The current go-qrcode implementation uses **single-mode encoding** for simplicity:

```go
// Current: Single mode for entire input
func (e *encoder) Encode(raw string) (*binary.Binary, error) {
    switch e.mode {
    case EncModeKanji:
        data = toShiftJIS(raw)  // All characters as Kanji
    case EncModeByte:
        data = []byte(raw)      // All characters as bytes
    // ...
    }
}
```

## Future Implementation

To support mixed mode encoding:

1. Add `EncModeMixed` constant
2. Implement segment detection
3. Implement per-segment encoding
4. Add optimization logic for mode switching decisions
5. Update version analysis to consider mixed mode capacity

## References

- [ISO/IEC 18004:2015](https://www.iso.org/standard/62021.html) - QR Code specification
- [Thonky QR Code Tutorial](https://www.thonky.com/qr-code-tutorial/) - Detailed encoding guide
- [Wikipedia: QR Code](https://en.wikipedia.org/wiki/QR_code) - Overview and encoding modes

---

# QR Code 混合编码模式

## 概述

QR 码支持多种编码模式，规范允许在单个 QR 码中切换模式。这使得混合内容（如汉字 + 数字 + ASCII）能够更高效地编码。

### 优势

- **最优效率**: 为每种字符类型使用最紧凑的编码
- **节省空间**: 减少混合内容的 QR 码总体大小
- **灵活性**: 处理多样化内容而不牺牲效率

## 模式指示器

每个编码段以 4 位模式指示器开头：

| 模式 | 二进制 | 十六进制 | 描述 |
|------|--------|----------|------|
| Numeric | `0001` | 0x1 | 数字 0-9 |
| Alphanumeric | `0010` | 0x2 | 0-9, A-Z 和 $%*+-./: |
| Byte | `0100` | 0x4 | ISO-8859-1 或 UTF-8 |
| Kanji | `1000` | 0x8 | Shift JIS 汉字字符 |
| ECI | `0111` | 0x7 | 扩展通道解释 |

## 数据结构

数据流中每个段的结构如下：

```
┌──────────────┬──────────────┬─────────────────┐
│ 模式         │ 字符计数     │ 数据            │
│ 指示器       │              │                 │
│ (4 位)       │ (可变)       │ (可变)          │
└──────────────┴──────────────┴─────────────────┘
```

### 字符计数指示器位数

字符计数指示器的位宽取决于 QR 版本和模式：

| 模式 | 版本 1-9 | 版本 10-26 | 版本 27-40 |
|------|---------|------------|------------|
| Numeric | 10 | 12 | 14 |
| Alphanumeric | 9 | 11 | 13 |
| Byte | 8 | 16 | 16 |
| Kanji | 8 | 10 | 12 |

## 编码示例

### 示例：编码 `日本語123`

#### 不使用混合模式（仅字节模式）

```
模式: Byte (0100)
字符计数: 9 字节 (UTF-8)
数据: 日本語123 = 9 字节 = 72 位
总计: 4 + 8 + 72 = 84 位
```

#### 使用混合模式（Kanji + Numeric）

```
段 1: Kanji 模式
┌────────┬────────┬─────────────────────────────┐
│ 模式   │ 计数   │ 数据                        │
│ 1000   │ 00000011 │ [日本語 每字符编码为 13 位]     │
│ 4 位   │ 8 位   │ 39 位 (3 × 13)              │
└────────┴────────┴─────────────────────────────┘

段 2: Numeric 模式
┌────────┬────────┬─────────────────┐
│ 模式   │ 计数   │ 数据            │
│ 0001   │ 00000011 │ 123            │
│ 4 位   │ 8 位   │ 10 位           │
└────────┴────────┴─────────────────┘

总计: (4 + 8 + 39) + (4 + 8 + 10) = 73 位
```

**节省: 84 - 73 = 11 位（效率提升 13%）**

## 编码算法

### 步骤 1: 分段输入

识别可以用相同模式编码的连续字符：

```go
type Segment struct {
    Mode  encMode
    Chars []rune
}

func segmentInput(raw string) []Segment {
    segments := []Segment{}
    currentMode := detectMode(raw[0])
    currentChars := []rune{}

    for _, r := range raw {
        mode := detectMode(r)
        if mode != currentMode {
            segments = append(segments, Segment{Mode: currentMode, Chars: currentChars})
            currentMode = mode
            currentChars = []rune{}
        }
        currentChars = append(currentChars, r)
    }

    segments = append(segments, Segment{Mode: currentMode, Chars: currentChars})
    return segments
}
```

### 步骤 2: 编码每个段

```go
func encodeSegments(segments []Segment, version int) *binary.Binary {
    dst := binary.New()

    for _, seg := range segments {
        // 模式指示器 (4 位)
        dst.AppendUint32(modeIndicator(seg.Mode), 4)

        // 字符计数（根据版本和模式变化）
        charCountBits := getCharCountBits(seg.Mode, version)
        dst.AppendUint32(uint32(len(seg.Chars)), charCountBits)

        // 数据
        encodeData(dst, seg.Mode, seg.Chars)
    }

    return dst
}
```

### 步骤 3: 添加终止符和填充

编码所有段后：
1. 如果空间允许，添加 4 位终止符 (`0000`)
2. 用零填充到字节边界
3. 添加填充字节（`0xEC`、`0x11` 交替）以填满容量

## 优化考量

### 何时切换模式

切换模式会增加开销（4 位模式指示器 + 字符计数位）。需要考虑：

```
切换成本 = 4 (模式) + 字符计数位

版本 1 示例:
- 切换到 Numeric: 4 + 10 = 14 位开销
- 切换到 Kanji: 4 + 8 = 12 位开销
```

### 模式选择策略

1. **贪心算法**: 为每个字符使用最高效的模式
2. **前瞻优化**: 考虑切换是否整体节省位数
3. **ECI**: 用于非 UTF-8 字符集

### 权衡示例

版本 1 中编码 `A日本語`：

| 方法 | 使用位数 |
|------|---------|
| 仅字节模式 | 4 + 8 + (4×8) = 44 位 |
| Alphanumeric + Kanji | (4+9+5.5) + (4+8+39) = 70.5 位 |

**结论**: 仅 1 个字符时不要切换 - 开销超过节省！

## 当前实现状态

当前 go-qrcode 实现为了简洁使用**单模式编码**：

```go
// 当前: 整个输入使用单一模式
func (e *encoder) Encode(raw string) (*binary.Binary, error) {
    switch e.mode {
    case EncModeKanji:
        data = toShiftJIS(raw)  // 所有字符作为 Kanji
    case EncModeByte:
        data = []byte(raw)      // 所有字符作为字节
    // ...
    }
}
```

## 未来实现

支持混合模式编码需要：

1. 添加 `EncModeMixed` 常量
2. 实现分段检测
3. 实现每段编码
4. 添加模式切换决策的优化逻辑
5. 更新版本分析以考虑混合模式容量

## 参考资料

- [ISO/IEC 18004:2015](https://www.iso.org/standard/62021.html) - QR 码规范
- [Thonky QR Code Tutorial](https://www.thonky.com/qr-code-tutorial/) - 详细编码指南
- [Wikipedia: QR Code](https://en.wikipedia.org/wiki/QR_code) - 概述和编码模式
