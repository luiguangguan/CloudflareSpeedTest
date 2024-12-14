package utils

import (
	"bufio"
	"fmt"
	"os"
)

// WriteToFile 将内容写入指定文件
// path: 文件路径
// content: 要写入的内容
// encoding: 文件编码方式，如 "utf-8", "gbk" 等
// append: 是否追加内容，true 表示追加，false 表示覆盖
func WriteToFile(path string, content string, encoding string, append bool) error {
	// 打开文件，如果文件不存在则创建
	var file *os.File
	var err error

	// 判断是否为追加模式
	if append {
		file, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	} else {
		file, err = os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
	}

	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// 如果编码是 "utf-8" 或其他默认编码，直接写入
	if encoding == "utf-8" || encoding == "" {
		writer := bufio.NewWriter(file)
		_, err := writer.WriteString(content)
		if err != nil {
			return fmt.Errorf("failed to write content to file: %v", err)
		}
		writer.Flush() // 确保内容写入文件
	} else {
		// 如果是其他编码方式的需求，你可以在这里添加相关的编码处理
		// 目前只支持直接写入 UTF-8 编码内容，若需要其他编码处理需要使用外部库
		return fmt.Errorf("unsupported encoding: %s", encoding)
	}

	return nil
}

// 判斷文件是否存在
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// 读取文件内容（UTF-8）
func ReadFileUTF8(filepath string) (text string, err error) {
	// 读取文件
	content, err := os.ReadFile(filepath) // 使用传入的 filepath 参数
	if err != nil {
		// 错误处理
		fmt.Println("读取文件失败:", err)
		return "", err // 如果出错，返回空字符串和错误信息
	}

	// 成功读取文件，返回文件内容
	return string(content), nil
}
