# Kanji Encoding Mode

## Overview

Kanji mode is a specialized encoding mode in QR Code designed to efficiently encode Japanese Kanji characters. It provides significant space savings compared to byte mode by leveraging the structure of Shift JIS (Japanese Industrial Standards) character encoding.

### Benefits

- **Compact Encoding**: Each Kanji character is encoded in 13 bits (vs. 16 bits in UTF-16 or 8 bytes per character in UTF-8)
- **Efficient Storage**: Reduces QR Code size for Japanese text by approximately 50% compared to byte mode
- **Optimized for Japanese**: Specifically designed for the Japanese writing system

## Specifications

| Parameter | Value |
|-----------|-------|
| Mode Indicator | `1000` (4 bits) |
| Bits per Character | 13 bits |
| Character Encoding | Shift JIS (JIS X 0208) |

### Character Count Indicator Bits

The number of bits used to store the character count varies by QR Code version:

| QR Code Version | Character Count Bits |
|-----------------|---------------------|
| 1-9 | 8 bits |
| 10-26 | 10 bits |
| 27-40 | 12 bits |

## Character Set

Kanji mode supports characters from the JIS X 0208 character set, encoded using Shift JIS. The valid Kanji characters fall within two ranges:

### Shift JIS Ranges

| Range | Start | End | Description |
|-------|-------|-----|-------------|
| Range 1 | 0x8140 | 0x9FFC | First Kanji block |
| Range 2 | 0xE040 | 0xEBBF | Second Kanji block |

### Unicode to Shift JIS Mapping

Modern applications typically work with Unicode characters. To use Kanji mode, Unicode characters must first be converted to their Shift JIS byte representation.

Example mappings:
- Unicode `U+4E16` (世) → Shift JIS `0x90 0xB6`
- Unicode `U+754C` (界) → Shift JIS `0x8A 0x79`

## Encoding Algorithm

### Step-by-Step Process

1. **Input**: Unicode Kanji character(s)
2. **Convert**: Transform each Unicode character to its Shift JIS 2-byte representation
3. **Adjust**: Apply the adjustment formula based on the Shift JIS value
4. **Encode**: Compress to 13-bit representation

### Mathematical Formula

Given a Shift JIS code `code` (2 bytes):

```
// Step 1: Adjust the base
if (code >= 0x8140 && code <= 0x9FFC) {
    adjusted = code - 0x8140
} else if (code >= 0xE040 && code <= 0xEBBF) {
    adjusted = code - 0xC140
}

// Step 2: Split into high and low bytes
high = adjusted >> 8      // Upper byte
low = adjusted & 0xFF     // Lower byte

// Step 3: Calculate encoded value (13-bit result)
encoded = (high × 0xC0) + low
```

### Why 0xC0?

The multiplier `0xC0` (192 in decimal) is derived from the Shift JIS encoding structure:
- In the valid ranges, the lower byte can be `0x40-0xFC` (except `0x7F`)
- This gives 188 possible values per high byte
- The encoding packs these efficiently: `high × 192 + low` results in at most 13 bits

### Encoding Example

Let's encode the Kanji character "世" (Unicode `U+4E16`):

1. **Convert to Shift JIS**: `0x90B6`
2. **Check range**: `0x90B6` is in range 1 (0x8140-0x9FFC)
3. **Adjust**: `0x90B6 - 0x8140 = 0x0F76`
4. **Split**: `high = 0x0F`, `low = 0x76`
5. **Encode**: `(0x0F × 0xC0) + 0x76 = 0x1176`
6. **Binary**: `1000101110110` (13 bits)

## Character Detection

To determine if a character is eligible for Kanji mode encoding:

### Detection Algorithm

1. **Check if character is Kanji**: The character must be in the Japanese Kanji Unicode ranges (primarily U+4E00-U+9FFF for CJK Unified Ideographs)
2. **Convert to Shift JIS**: Attempt conversion from Unicode to Shift JIS
3. **Validate range**: The resulting Shift JIS value must be in:
   - `0x8140` to `0x9FFC`, OR
   - `0xE040` to `0xEBBF`
4. **Check byte length**: Each character must encode to exactly 2 bytes in Shift JIS

### Detection Criteria Summary

```
IsKanji(character) {
    shiftJIS = UnicodeToShiftJIS(character)

    if (shiftJIS.length != 2) {
        return false
    }

    code = (shiftJIS[0] << 8) | shiftJIS[1]

    return (code >= 0x8140 && code <= 0x9FFC) ||
           (code >= 0xE040 && code <= 0xEBBF)
}
```

### Automatic Mode Selection

When encoding mixed content, use Kanji mode only when:
- ALL characters in the data are valid Kanji characters
- Each character successfully converts to a valid Shift JIS value in the allowed ranges
- The content is primarily Japanese text

If any character fails validation, fall back to a compatible mode (typically byte mode with UTF-8).

## Practical Considerations

### When to Use Kanji Mode

- Japanese text containing primarily Kanji characters
- When minimizing QR Code size is critical
- When target scanners support Kanji mode decoding

### Limitations

- Only supports JIS X 0208 characters
- Hiragana and Katakana are NOT supported (use byte mode)
- Some rare Kanji characters outside the ranges cannot be encoded
- Requires proper Shift JIS conversion capability

### Compatibility

Most modern QR Code scanners support Kanji mode, but for maximum compatibility with older scanners, consider using byte mode with UTF-8 encoding, especially for international applications.

---

# Kanji 编码模式

## 概述

Kanji 模式是 QR 码中专门设计的一种编码模式，用于高效编码日文汉字字符。通过利用 Shift JIS（日本工业标准）字符编码的结构，它相比字节模式能显著节省空间。

### 优势

- **紧凑编码**: 每个汉字字符编码为 13 位（相比 UTF-16 的 16 位或 UTF-8 的每字符 8 字节）
- **高效存储**: 相比字节模式，可将日文文本的 QR 码大小减少约 50%
- **专为日文优化**: 专门针对日文字符系统设计

## 规范说明

| 参数 | 值 |
|------|-----|
| 模式指示器 | `1000` (4 位) |
| 每字符位数 | 13 位 |
| 字符编码 | Shift JIS (JIS X 0208) |

### 字符计数指示器位数

用于存储字符计数的位数随 QR 码版本变化：

| QR 码版本 | 字符计数位数 |
|-----------|-------------|
| 1-9 | 8 位 |
| 10-26 | 10 位 |
| 27-40 | 12 位 |

## 字符集

Kanji 模式支持 JIS X 0208 字符集中的字符，使用 Shift JIS 编码。有效的汉字字符落在两个范围内：

### Shift JIS 范围

| 范围 | 起始 | 结束 | 描述 |
|------|------|------|------|
| 范围 1 | 0x8140 | 0x9FFC | 第一汉字块 |
| 范围 2 | 0xE040 | 0xEBBF | 第二汉字块 |

### Unicode 到 Shift JIS 映射

现代应用通常使用 Unicode 字符。要使用 Kanji 模式，Unicode 字符必须首先转换为其 Shift JIS 双字节表示。

映射示例：
- Unicode `U+4E16` (世) → Shift JIS `0x90 0xB6`
- Unicode `U+754C` (界) → Shift JIS `0x8A 0x79`

## 编码算法

### 逐步流程

1. **输入**: Unicode 汉字字符
2. **转换**: 将每个 Unicode 字符转换为其 Shift JIS 双字节表示
3. **调整**: 根据 Shift JIS 值应用调整公式
4. **编码**: 压缩为 13 位表示

### 数学公式

给定 Shift JIS 代码 `code`（2 字节）：

```
// 步骤 1: 调整基数
if (code >= 0x8140 && code <= 0x9FFC) {
    adjusted = code - 0x8140
} else if (code >= 0xE040 && code <= 0xEBBF) {
    adjusted = code - 0xC140
}

// 步骤 2: 拆分为高位和低位字节
high = adjusted >> 8      // 高位字节
low = adjusted & 0xFF     // 低位字节

// 步骤 3: 计算编码值（13 位结果）
encoded = (high × 0xC0) + low
```

### 为什么是 0xC0？

乘数 `0xC0`（十进制 192）源自 Shift JIS 编码结构：
- 在有效范围内，低位字节可以是 `0x40-0xFC`（除 `0x7F` 外）
- 这为每个高位字节提供 188 个可能值
- 编码将其高效打包：`high × 192 + low` 结果最多为 13 位

### 编码示例

让我们编码汉字 "世"（Unicode `U+4E16`）：

1. **转换为 Shift JIS**: `0x90B6`
2. **检查范围**: `0x90B6` 在范围 1 内 (0x8140-0x9FFC)
3. **调整**: `0x90B6 - 0x8140 = 0x0F76`
4. **拆分**: `high = 0x0F`, `low = 0x76`
5. **编码**: `(0x0F × 0xC0) + 0x76 = 0x1176`
6. **二进制**: `1000101110110` (13 位)

## 字符检测

判断字符是否符合 Kanji 模式编码条件：

### 检测算法

1. **检查是否为汉字**: 字符必须在日文汉字 Unicode 范围内（主要是 U+4E00-U+9FFF 的 CJK 统一表意文字）
2. **转换为 Shift JIS**: 尝试从 Unicode 转换为 Shift JIS
3. **验证范围**: 结果 Shift JIS 值必须在：
   - `0x8140` 到 `0x9FFC`，或
   - `0xE040` 到 `0xEBBF`
4. **检查字节长度**: 每个字符在 Shift JIS 中必须编码为恰好 2 字节

### 检测标准总结

```
IsKanji(character) {
    shiftJIS = UnicodeToShiftJIS(character)

    if (shiftJIS.length != 2) {
        return false
    }

    code = (shiftJIS[0] << 8) | shiftJIS[1]

    return (code >= 0x8140 && code <= 0x9FFC) ||
           (code >= 0xE040 && code <= 0xEBBF)
}
```

### 自动模式选择

编码混合内容时，仅当满足以下条件时使用 Kanji 模式：
- 数据中的所有字符都是有效的汉字字符
- 每个字符都能成功转换为允许范围内的有效 Shift JIS 值
- 内容主要是日文文本

如果任何字符验证失败，则回退到兼容模式（通常为 UTF-8 字节模式）。

## 实际考虑

### 何时使用 Kanji 模式

- 主要包含汉字字符的日文文本
- 最小化 QR 码大小至关重要时
- 目标扫描器支持 Kanji 模式解码

### 限制

- 仅支持 JIS X 0208 字符
- 不支持平假名和片假名（使用字节模式）
- 范围外的某些罕见汉字无法编码
- 需要正确的 Shift JIS 转换能力

### 兼容性

大多数现代 QR 码扫描器支持 Kanji 模式，但为了与旧扫描器实现最大兼容性，对于国际应用可考虑使用带 UTF-8 编码的字节模式。
