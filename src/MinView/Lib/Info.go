package Lib

import (
	"errors"
)

const VI_APP_NAME string = "ViBoard"
const VI_VERSION string = "1.0.0"

// 这个是固定的，用来获取第一个版本的APPKEY
const VI_INIT_VERSION string = "1.0.0"

var VI_APP_KEYS = map[string]string{
	"1.0.0": "NTestProgamKeys:open@zerosir.com",
}

const VI_COPYRIGHT string = "ZEROSIR"
const VI_URL string = "https://zerosir.com"

// root
const VI_EMAIL string = "root@zerosir.com"

// public
const VI_REPORT string = "public@zerosir.com"

var ras_APP_KEYS = map[int]string{
	0: "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",
	1: "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",
}

// 生成公共程序密钥 ras_APP_KEYS
func GenerateAppKeys(privateKeyPEM string) (string, error) {
	keyss, err := EncryptData([]byte(VI_APP_KEYS[VI_INIT_VERSION]), []byte(privateKeyPEM))
	if err != nil {
		return "", err
	}
	return EncryptBase64(keyss), nil
}

// 获取程序RSA数量
func GetAPPRsaSize() int {
	return len(ras_APP_KEYS)
}

// 获取程序RSA密钥PEM
func GetAppRsaKey(index int) ([]byte, error) {
	if index >= len(ras_APP_KEYS) {
		return nil, errors.New("RSA数组索引越界")
	}
	key := ras_APP_KEYS[index]
	// rsaBase64, err := base64.StdEncoding.DecodeString(string(key))
	rsaBase64, err := DecryptBase64(string(key))
	if err != nil {
		return nil, err
	}
	rsaPEM, err := DecryptData([]byte(VI_APP_KEYS[VI_INIT_VERSION]), []byte(rsaBase64))
	if err != nil {
		return nil, err
	}
	return rsaPEM, nil
}

// 前端web ui post信息
const (
	//未定义
	PostTypeUndefined = "Undefined"
	//webui获取显示服务表信息，包含 文件Key 文件名称 表名 表key
	PostTypeDisplayGetVolume = "PostTypeDisplayGetVolume"
	//webui获取显示服务组件信息
	PostTypeDisplayGetComponent = "PostTypeDisplayGetComponent"
	//webui 更新表数据,可能通过文件来更新表
	PostTypeImportSheet = "Sheet"
	//webui 更新表数据,可能通过JSON来更新表，暂时未启用
	// PostTypeImporData = "Data"
	//webui 更新表数据,可能通过JSON来更新
	PostTypeUpdateComponent = "Component"
	//前端UI上传一些文件到数据库，用来远程更新文件
	PostTypeUploadFile = "File"
)

// 界面触发信息给后端
const (
	//未定义
	UIEmitUndefined = "Undefined"
	//远程更新组件信息后，显示器emit给后端
	UIEmitComponentsData = "UIEmitComponentsData"
	//远程上传表格文件到显示器数据库后，emit表key给显示器，用于刷新echart数据
	UIEmitSheetData = "UIEmitSheetData"
)
