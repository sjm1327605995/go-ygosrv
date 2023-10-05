package utils

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
)

func UTF16ToStr(arr []byte) ([]byte, string, error) {
	if len(arr)%2 != 0 {
		return nil, "", errors.New("not utf16")
	}
	if len(arr) == 0 {
		return nil, "", nil
	}
	var last int
	for ; last < len(arr); last += 2 {
		//有效数据 2个位置都不能等于0
		if arr[last] == 0 && arr[last+1] == 0 {
			break
		}
	}

	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	// 将二进制数组转换为字符串
	str, err := io.ReadAll(transform.NewReader(bytes.NewReader(arr[:last]), decoder))
	if err != nil {
		fmt.Println("解码失败:", err)
		return nil, "", err
	}
	return arr[:last], string(str), nil
}
func UTF16ToStrArr(arr []byte) ([]byte, error) {
	if len(arr)%2 != 0 {
		return nil, errors.New("not utf16")
	}
	if len(arr) == 0 {
		return nil, nil
	}
	var last int
	for ; last < len(arr); last += 2 {
		//有效数据 2个位置都不能等于0
		if arr[last] == 0 && arr[last+1] == 0 {
			break
		}
	}

	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	// 将二进制数组转换为字符串
	str, err := io.ReadAll(transform.NewReader(bytes.NewReader(arr[:last]), decoder))
	if err != nil {
		fmt.Println("解码失败:", err)
		return nil, err
	}
	return str, nil
}
