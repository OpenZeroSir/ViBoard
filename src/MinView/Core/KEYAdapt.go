package Core

import (
	"MinView/Lib"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/gob"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	mrand "math/rand"
	"os"
	"strconv"
	"time"

	version "github.com/hashicorp/go-version"
	"github.com/ulikunitz/xz"
)

/*
密钥结构
*/
type KEY struct {
	//服务KEY
	serverKey string
	//本地私有KEY
	localPrivateKey string
	//本地共有KEY
	localPublicKey string
	//文档key,这个key用来根密码一起加密数据
	documentKey string
	//显示端数据加密key,公钥在编辑器端，私钥在显示端
	//显示私钥
	displayPrivateKey string
	//显示公钥
	displayPublicKey string
	//程序对称加密密钥，每个版本都要发生变化,不要保存到文档或项目
	programDesKey []byte
}

func NewKEY() (*KEY, error) {
	var key KEY
	err := key.initAppKeys()
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// 初始化KEY
func (a *KEY) initAppKeys() error {
	a.serverKey = ""
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	publicKey := &privateKey.PublicKey
	privateKeyPEM := a.encodePrivateKeyToPEM(privateKey)
	a.localPrivateKey = string(privateKeyPEM)
	// 将公钥以 PEM 格式编码
	publicKeyPEM, err := a.encodePublicKeyToPEM(publicKey)
	if err != nil {
		return err
	}
	a.localPublicKey = string(publicKeyPEM)
	docKey, err := Lib.GenerateRandomByte(16)
	if err != nil {
		return err
	}
	a.documentKey = string(docKey)
	displayPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	displayPublicKey := &displayPrivateKey.PublicKey
	displayPrivateKeyPEM := a.encodePrivateKeyToPEM(displayPrivateKey)
	a.displayPrivateKey = string(displayPrivateKeyPEM)
	displayPublicKeyPEM, err := a.encodePublicKeyToPEM(displayPublicKey)
	if err != nil {
		return err
	}
	a.displayPublicKey = string(displayPublicKeyPEM)
	a.programDesKey = []byte(Lib.VI_APP_KEYS[Lib.VI_VERSION])
	return nil
}

// 将私钥编码为 PEM 格式
func (a KEY) encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	return privateKeyPEM
}

// 将公钥编码为 PEM 格式
func (a KEY) encodePublicKeyToPEM(publicKey *rsa.PublicKey) ([]byte, error) {
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	return publicKeyPEM, nil
}

// 获取服务器公钥
func (a KEY) GetServerKey() string {
	return a.serverKey
}

// 获取私钥PEM
func (a KEY) GetPrivateKeyPEM() string {
	return a.localPrivateKey
}

// 获取公钥PEM
func (a KEY) GetPublicKeyPEM() string {
	return a.localPublicKey
}

// 从PEM获取私钥
func (a *KEY) GetPrivateKey() (*rsa.PrivateKey, error) {
	// 从PEM格式还原私钥
	block, _ := pem.Decode([]byte(a.localPrivateKey))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("无效的PEM格式私钥")
	}
	parsedPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return parsedPrivateKey, nil
}

// 从PEM获取公钥
func (a KEY) GetPublicKey() (*rsa.PublicKey, error) {
	// 解码PEM格式的公钥
	block, _ := pem.Decode([]byte(a.localPublicKey))
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("无效的PEM格式公钥")
	}
	// 解析RSA编码的公钥
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

// 获取文档key
func (a KEY) GetDocumentKey() []byte {
	return []byte(a.documentKey)
}

// 获取显示私钥PEM
func (a KEY) GetDisplayPrivateKeyPEM() string {
	return a.displayPrivateKey
}

// 获取显示公钥PEM
func (a KEY) GetDisplayPublicKeyPEM() string {
	return a.displayPublicKey
}

// 从显示PEM获取私钥
func (a KEY) GetDisplayPrivateKey() (*rsa.PrivateKey, error) {
	// 从PEM格式还原私钥
	block, _ := pem.Decode([]byte(a.displayPrivateKey))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("无效的PEM格式私钥")
	}
	parsedPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return parsedPrivateKey, nil
}

// 从显示PEM获取公钥
func (a KEY) GetDisplayPublicKey() (*rsa.PublicKey, error) {
	// 解码PEM格式的公钥
	block, _ := pem.Decode([]byte(a.displayPublicKey))
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, errors.New("无效的PEM格式公钥")
	}
	// 解析RSA编码的公钥
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

// 获取程序密钥
func (a KEY) GetProgramDesKey() []byte {
	return a.programDesKey
}

// 初始化服务器KEY,服务器密钥应当保存一份KEY在程序中，
// 用于方便服务器和程序通信加密
// 返回结果中，用string就是方便前端调用，
// 因为有时候可能，程序没有联网，初始化肯能会失败，
// 失败之后应当不影响程序启动，前端检测到联网后应当跟服务器取得联系，
// 并读取到服务器密钥，方便记载更多组件数据到程序中。
// 放回结果中不是空字符串就是有错误的
// 从服务器获取到的KEY应该是PEM格式的，并保存到serverKey就可以
func (a *KEY) InitServerKey() {
	a.serverKey = "This is test key! need implement"
}

// 生成的KEY保存在string中，在需要的时候转成[]byte数组，
// 方便用来加密解码,该函数只是用来提供方便站换
func (a KEY) StringToBytes(str string) []byte {
	if len(str) < 1 {
		return nil
	}
	return []byte(str)
}

// 加密项目数据保存到file.,
// dbPath是数据库Backup得到的路径，
// obj_data是视图对象数据
func (a KEY) EncryptProjectData(file *os.File, dbPath string, obj_data string) error {
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	privateKeyPEM := a.GetPrivateKeyPEM()
	programDesKey := a.GetProgramDesKey()
	privateKeyResult, err := Lib.EncryptData(programDesKey, []byte(privateKeyPEM))
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(file)
	patch, err := Lib.GenerateRandomByte(16)
	if err != nil {
		return err
	}
	mrand.Seed(time.Now().UnixNano())
	rsa_index := mrand.Intn(Lib.GetAPPRsaSize())
	app_rsaPEM, err := Lib.GetAppRsaKey(rsa_index)
	if err != nil {
		return err
	}
	desKey := Lib.CryptPasswd(patch, string(a.GetDocumentKey()), -1)
	//文件类型
	data, err := Lib.EncryptData(desKey, []byte("NNIE"))
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//程序名称
	data, err = Lib.EncryptData(desKey, []byte(Lib.VI_APP_NAME))
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//程序版本
	data, err = Lib.EncryptData([]byte(Lib.VI_APP_KEYS[Lib.VI_INIT_VERSION]), []byte(Lib.VI_VERSION))
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//程序版权
	data, err = Lib.EncryptData(desKey, []byte(Lib.VI_COPYRIGHT))
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//程序所有
	data, err = Lib.EncryptData(desKey, []byte(Lib.VI_EMAIL))
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//程序官网
	data, err = Lib.EncryptData(desKey, []byte(Lib.VI_URL))
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//问题报告
	data, err = Lib.EncryptData(desKey, []byte(Lib.VI_REPORT))
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//写入程序rsa索引
	err = encoder.Encode(rsa_index)
	if err != nil {
		return err
	}
	//写入私钥
	err = encoder.Encode(privateKeyResult)
	if err != nil {
		return err
	}
	//获取公钥
	publicKey, err := a.GetPublicKey()
	if err != nil {
		return err
	}
	//保存文档key
	documentKeyReuslt, err := Lib.EncryptRSAData(publicKey, a.GetDocumentKey())
	if err != nil {
		return err
	}
	err = encoder.Encode(documentKeyReuslt)
	if err != nil {
		return err
	}
	//保存文档patch
	// patchResult, err := EncryptRSAData(publicKey, patch)
	// if err != nil {
	// 	return err
	// }
	app_privateKey, err := Lib.PEM2PrivateKey(app_rsaPEM)
	if err != nil {
		return err
	}
	patch_init, err := Lib.EncryptRSAData(&app_privateKey.PublicKey, patch)
	if err != nil {
		return err
	}
	err = encoder.Encode(patch_init)
	if err != nil {
		return err
	}
	//保存服务器公钥
	data, err = Lib.EncryptData(desKey, []byte(a.GetServerKey()))
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//保存显示私钥
	data, err = Lib.EncryptData(desKey, []byte(a.GetDisplayPrivateKeyPEM()))
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//保存视图对象文件
	data, err = Lib.EncryptData(desKey, []byte(obj_data))
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//读取数据库文件
	data, err = os.ReadFile(dbPath)
	if err != nil {
		return err
	}
	// // 创建LZMA写入器，并设置压缩级别
	// var buf bytes.Buffer
	// writer, err := xz.NewWriter(&buf)
	// if err != nil {
	// 	return err
	// }
	// _, err = io.WriteString(writer, string(data))
	// if err != nil {
	// 	return err
	// }
	// if err := writer.Close(); err != nil {
	// 	return err
	// }
	//获取私钥
	privateKey, err := a.GetPrivateKey()
	if err != nil {
		return err
	}
	//秘密文件内容
	file_data, err := Lib.EncryptData(desKey, data)
	if err != nil {
		return err
	}
	//签名文件内容
	sig_result := Lib.SignData(file_data, privateKey)
	if sig_result.Err != nil {
		return err
	}
	//写入文件签名，签名需要加密，避免签名被修改
	data, err = Lib.EncryptData(desKey, sig_result.Signature)
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//写入文件哈希，哈希没加密，因为签名已经加密了
	err = encoder.Encode(sig_result.Hashed)
	if err != nil {
		return err
	}
	//写入文件内容
	err = encoder.Encode(file_data)
	if err != nil {
		return err
	}
	return nil
}

type DecryptProjectResult struct {
	//视图对象数据
	ViewData string
	//数据库数据
	DbData []byte
	//处理结果
	Err error
}

// 解密项目文件
func (a *KEY) DecryptProjectData(file_path string) DecryptProjectResult {
	var result DecryptProjectResult
	result.Err = nil
	file, err := os.Open(file_path)
	if err != nil {
		result.Err = err
		return result
	}
	defer file.Close()
	var keyStruct KEY
	var data []byte
	decoder := gob.NewDecoder(file)
	//文件类型
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	file_type := data
	//程序名称
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	app_name := data
	//程序版本
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	data, err = Lib.DecryptData([]byte(Lib.VI_APP_KEYS[Lib.VI_INIT_VERSION]), data)
	if err != nil {
		result.Err = err
		return result
	}
	app_version := string(data)
	v1, err := version.NewVersion(app_version)
	if err != nil {
		result.Err = err
		return result
	}
	v2, err := version.NewVersion(Lib.VI_VERSION)
	if err != nil {
		result.Err = err
		return result
	}
	if v2.LessThan(v1) {
		result.Err = errors.New("程序版本太低，请升级到最新版本！")
		return result
	}
	//程序版权
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	app_copyright := data
	//程序所有
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	app_email := data
	//程序官网
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	app_url := data
	//问题报告
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	app_report := data
	//读取程序rsa索引
	var rsa_index int
	err = decoder.Decode(&rsa_index)
	if err != nil {
		result.Err = err
		return result
	}
	//解密私钥PEM
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	data, err = Lib.DecryptData([]byte(Lib.VI_APP_KEYS[app_version]), data)
	if err != nil {
		result.Err = err
		return result
	}
	//修改本地私钥
	keyStruct.localPrivateKey = string(data)
	privateKey, err := keyStruct.GetPrivateKey()
	if err != nil {
		result.Err = err
		return result
	}
	//读取公钥PEM
	publicKey, err := a.encodePublicKeyToPEM(&privateKey.PublicKey)
	if err != nil {
		result.Err = err
		return result
	}
	keyStruct.localPublicKey = string(publicKey)
	//读取文档KEY
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	data, err = Lib.DecryptRSAData(privateKey, data)
	if err != nil {
		result.Err = err
		return result
	}
	keyStruct.documentKey = string(data)
	//读取文档patch
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	app_rsaPEM, err := Lib.GetAppRsaKey(rsa_index)
	if err != nil {
		result.Err = err
		return result
	}
	app_privateKey, err := Lib.PEM2PrivateKey(app_rsaPEM)
	if err != nil {
		result.Err = err
		return result
	}
	data, err = Lib.DecryptRSAData(app_privateKey, data)
	if err != nil {
		result.Err = err
		return result
	}
	//生成DES-265密钥
	desKey := Lib.CryptPasswd(data, keyStruct.documentKey, -1)
	//基础信息检查
	info, err := Lib.DecryptData(desKey, file_type)
	if err != nil {
		result.Err = err
		return result
	}
	if string(info) != "NNIE" {
		result.Err = errors.New("未知文件类型")
		return result
	}
	info, err = Lib.DecryptData(desKey, app_name)
	if err != nil {
		result.Err = err
		return result
	}
	if string(info) != Lib.VI_APP_NAME {
		result.Err = errors.New("未知程序文件")
		return result
	}
	info, err = Lib.DecryptData(desKey, app_copyright)
	if err != nil {
		result.Err = err
		return result
	}
	if string(info) != Lib.VI_COPYRIGHT {
		result.Err = errors.New("未知程序版权")
		return result
	}
	info, err = Lib.DecryptData(desKey, app_email)
	if err != nil {
		result.Err = err
		return result
	}
	if string(info) != Lib.VI_EMAIL {
		result.Err = errors.New("未知程序所有")
		return result
	}
	info, err = Lib.DecryptData(desKey, app_url)
	if err != nil {
		result.Err = err
		return result
	}
	if string(info) != Lib.VI_URL {
		result.Err = errors.New("未知程序网站")
		return result
	}
	info, err = Lib.DecryptData(desKey, app_report)
	if err != nil {
		result.Err = err
		return result
	}
	if string(info) != Lib.VI_REPORT {
		result.Err = errors.New("未知程序报告邮箱")
		return result
	}
	//服务器密钥
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	data, err = Lib.DecryptData(desKey, data)
	if err != nil {
		result.Err = err
		return result
	}
	keyStruct.serverKey = string(data)
	//读取显示私钥
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	data, err = Lib.DecryptData(desKey, data)
	if err != nil {
		result.Err = err
		return result
	}
	keyStruct.displayPrivateKey = string(data)
	//读取显示公钥PEM
	displayPrivateKey, err := Lib.PEM2PrivateKey([]byte(keyStruct.displayPrivateKey))
	if err != nil {
		result.Err = err
		return result
	}
	dispalyPublicKey, err := a.encodePublicKeyToPEM(&displayPrivateKey.PublicKey)
	if err != nil {
		result.Err = err
		return result
	}
	keyStruct.displayPublicKey = string(dispalyPublicKey)
	//读取视图对象内容
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	data, err = Lib.DecryptData(desKey, data)
	if err != nil {
		result.Err = err
		return result
	}
	result.ViewData = string(data)
	//读取数据库签名
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	data, err = Lib.DecryptData(desKey, data)
	if err != nil {
		result.Err = err
		return result
	}
	var sig = data
	//读取数据库hash
	var hashed [32]byte
	err = decoder.Decode(&hashed)
	if err != nil {
		result.Err = err
		return result
	}

	//读取数据库
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	//签名文件内容
	sig_result := Lib.SignData(data, privateKey)
	if sig_result.Err != nil {
		result.Err = err
		return result
	}
	//校验签名
	err = Lib.CheckSign(sig_result.Hashed, sig, &privateKey.PublicKey)
	if err != nil {
		result.Err = err
		return result
	}
	err = Lib.CheckSign(hashed, sig_result.Signature, &privateKey.PublicKey)
	if err != nil {
		result.Err = err
		return result
	}
	data, err = Lib.DecryptData(desKey, data)
	if err != nil {
		result.Err = err
		return result
	}
	// //解压
	// reader, err := xz.NewReader(bytes.NewBuffer(data))
	// if err != nil {
	// 	result.Err = err
	// 	return result
	// }
	// var buf bytes.Buffer
	// if _, err = io.Copy(&buf, reader); err != nil {
	// 	result.Err = err
	// 	return result
	// }
	result.DbData = data
	//刷新KEY
	err = a.initAppKeys()
	if err != nil {
		result.Err = err
		return result
	}
	if keyStruct.serverKey != "" {
		a.serverKey = keyStruct.serverKey
	} else {
		a.InitServerKey()
	}
	return result
}

// 导出视图编码，obj_data是视图对象信息，docinfo是文档信息
func (a KEY) ExportViewCode(file *os.File, dbPath string, obj_data string, docinfo string, passwd string) error {
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	mrand.Seed(time.Now().UnixNano())
	//首次密码迭代索引
	init_index := mrand.Intn(10000) + 10000
	appkey := Lib.VI_APP_KEYS[Lib.VI_INIT_VERSION]
	new_passwd := Lib.CryptPasswd([]byte(appkey), passwd, init_index)
	passwdHash := sha256.Sum256(new_passwd)
	passwordHash := passwdHash[:]
	encoder := gob.NewEncoder(file)
	//写入文档信息,这个是需要在验证密码之前显示的
	data, err := Lib.EncryptData([]byte(Lib.VI_APP_KEYS[Lib.VI_INIT_VERSION]), []byte(docinfo))
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//加密和写入初始迭代索引，采用第一个RSA加密
	init_PEM, err := Lib.GetAppRsaKey(0)
	if err != nil {
		return err
	}
	initPrivateKey, err := Lib.PEM2PrivateKey(init_PEM)
	if err != nil {
		return err
	}
	init_index_crypto, err := Lib.EncryptRSAData(&initPrivateKey.PublicKey, []byte(strconv.Itoa(init_index)))
	if err != nil {
		return err
	}
	err = encoder.Encode(init_index_crypto)
	if err != nil {
		return err
	}
	//程序公共RSA KEY
	rsa_index := mrand.Intn(Lib.GetAPPRsaSize())
	app_rsaPEM, err := Lib.GetAppRsaKey(rsa_index)
	if err != nil {
		return err
	}
	//写入程序rsa索引
	data, err = Lib.EncryptData([]byte(passwordHash), []byte(strconv.Itoa(rsa_index)))
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	appPrivateKey, err := Lib.PEM2PrivateKey(app_rsaPEM)
	if err != nil {
		return err
	}
	//保存文档key
	data, err = Lib.EncryptRSAData(&appPrivateKey.PublicKey, a.GetDocumentKey())
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//文档新密码
	docPasswd := Lib.CryptPasswd(a.GetDocumentKey(), string(passwordHash), -1)
	//保存显示私钥
	data, err = Lib.EncryptData(docPasswd, []byte(a.GetDisplayPrivateKeyPEM()))
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//写入文档信息,这个是需要在验证密码之后显示的
	data, err = Lib.EncryptData(docPasswd, []byte(docinfo))
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//patch新密码
	patch, err := Lib.GenerateRandomByte(16)
	if err != nil {
		return err
	}
	displayPrivateKey, err := a.GetDisplayPrivateKey()
	if err != nil {
		return err
	}
	//保存patch key
	data, err = Lib.EncryptRSAData(&displayPrivateKey.PublicKey, patch)
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	patchPasswd := Lib.CryptPasswd(patch, string(passwordHash), -1)

	//文件类型
	data, err = Lib.EncryptData(docPasswd, []byte("NNIC"))
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//程序名称
	data, err = Lib.EncryptData(docPasswd, []byte(Lib.VI_APP_NAME))
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//程序版本
	data, err = Lib.EncryptData(docPasswd, []byte(Lib.VI_VERSION))
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//程序版权
	data, err = Lib.EncryptData(docPasswd, []byte(Lib.VI_COPYRIGHT))
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//程序所有
	data, err = Lib.EncryptData(docPasswd, []byte(Lib.VI_EMAIL))
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//程序官网
	data, err = Lib.EncryptData(docPasswd, []byte(Lib.VI_URL))
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//问题报告
	data, err = Lib.EncryptData(docPasswd, []byte(Lib.VI_REPORT))
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//保存服务器公钥
	data, err = Lib.EncryptData(patchPasswd, []byte(a.GetServerKey()))
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//保存视图对象
	data, err = Lib.EncryptData(patchPasswd, []byte(obj_data))
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//签名视图对象
	sig_result := Lib.SignData(data, appPrivateKey)
	if sig_result.Err != nil {
		return err
	}
	//写入视图对象签名，签名需要加密，避免签名被修改
	data, err = Lib.EncryptData(docPasswd, sig_result.Signature)
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//写入文件哈希，哈希没加密，因为签名已经加密了
	err = encoder.Encode(sig_result.Hashed)
	if err != nil {
		return err
	}
	//读取数据库文件
	data, err = os.ReadFile(dbPath)
	if err != nil {
		return err
	}
	// 创建LZMA写入器，并设置压缩级别
	var buf bytes.Buffer
	writer, err := xz.NewWriter(&buf)
	if err != nil {
		return err
	}
	_, err = io.WriteString(writer, string(data))
	if err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	//合并密钥hash
	all_hash := sha256.Sum256([]byte(string(patchPasswd) + string(docPasswd)))
	//加密文件内容
	file_data, err := Lib.EncryptData(all_hash[:], buf.Bytes())
	if err != nil {
		return err
	}
	//需要生成第二次密码
	left_index := ((int(patchPasswd[rsa_index%10]) + int(all_hash[rsa_index%10])) % 10) * 5000
	patchPasswd2 := Lib.CryptPasswd(patch, string(passwordHash), left_index)
	//二次加密文件内容
	file_data, err = Lib.EncryptData(patchPasswd2, file_data)
	if err != nil {
		return err
	}
	//签名文件内容
	sig_result = Lib.SignData(file_data, displayPrivateKey)
	if sig_result.Err != nil {
		return err
	}
	//写入文件签名，签名需要加密，避免签名被修改
	data, err = Lib.EncryptData(docPasswd, sig_result.Signature)
	if err != nil {
		return err
	}
	err = encoder.Encode(data)
	if err != nil {
		return err
	}
	//写入文件哈希，哈希没加密，因为签名已经加密了
	err = encoder.Encode(sig_result.Hashed)
	if err != nil {
		return err
	}
	//写入文件内容
	err = encoder.Encode(file_data)
	if err != nil {
		return err
	}
	return nil
}

type DecryptViewCodeResult struct {
	//服务器密钥
	ServerKey string
	//视图对象数据
	ViewData string
	//数据库数据
	DbData []byte
	//视图编码信息
	ViewInfo string
	//处理结果
	Err error
}

// 导入视图编码
func (a KEY) ImportViewCode(file_path string, passwd string) DecryptViewCodeResult {
	var result DecryptViewCodeResult
	file, err := os.Open(file_path)
	if err != nil {
		result.Err = err
		return result
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	//这里读取到的文档信息，可能是被伪造过的，所以直接跳过就行
	var docinfo []byte
	err = decoder.Decode(&docinfo)
	if err != nil {
		result.Err = err
		return result
	}
	//读取首次密码迭代索引
	var init_index_crypto []byte
	err = decoder.Decode(&init_index_crypto)
	if err != nil {
		result.Err = err
		return result
	}
	//读取和解密初始迭代索引，采用第一个RSA解密
	init_PEM, err := Lib.GetAppRsaKey(0)
	if err != nil {
		result.Err = err
		return result
	}
	initPrivateKey, err := Lib.PEM2PrivateKey(init_PEM)
	if err != nil {
		result.Err = err
		return result
	}
	init_index, err := Lib.DecryptRSAData(initPrivateKey, init_index_crypto)
	if err != nil {
		result.Err = err
		return result
	}
	init_index_num, err := strconv.Atoi(string(init_index))
	if err != nil {
		result.Err = err
		return result
	}
	appkey := Lib.VI_APP_KEYS[Lib.VI_INIT_VERSION]
	new_passwd := Lib.CryptPasswd([]byte(appkey), passwd, init_index_num)
	//密码hash
	passwdHash := sha256.Sum256(new_passwd)
	passwordHash := passwdHash[:]
	//读取程序rsa索引
	var rsa_index []byte
	err = decoder.Decode(&rsa_index)
	if err != nil {
		result.Err = err
		return result
	}
	data, err := Lib.DecryptData([]byte(passwordHash), []byte(rsa_index))
	if err != nil {
		result.Err = err
		return result
	}
	index, err := strconv.Atoi(string(data))
	if err != nil {
		result.Err = err
		return result
	}
	//程序私有RSA KEY
	app_rsaPEM, err := Lib.GetAppRsaKey(index)
	if err != nil {
		result.Err = err
		return result
	}
	appPrivateKey, err := Lib.PEM2PrivateKey(app_rsaPEM)
	if err != nil {
		result.Err = err
		return result
	}
	//读取文档key
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	docKey, err := Lib.DecryptRSAData(appPrivateKey, data)
	if err != nil {
		result.Err = err
		return result
	}
	//文档新密码
	docPasswd := Lib.CryptPasswd(docKey, string(passwordHash), -1)
	//读取文档私钥，也就是显示私钥
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	//显示私钥PEM
	data, err = Lib.DecryptData(docPasswd, data)
	if err != nil {
		result.Err = err
		return result
	}
	displayPrivateKey, err := Lib.PEM2PrivateKey(data)
	if err != nil {
		result.Err = err
		return result
	}
	//读取文档信息,这个是需要在验证密码之后显示的
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	docInfo, err := Lib.DecryptData(docPasswd, data)
	if err != nil {
		result.Err = err
		return result
	}
	result.ViewInfo = string(docInfo)
	//读取patch KEY
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	patchKey, err := Lib.DecryptRSAData(displayPrivateKey, data)
	if err != nil {
		result.Err = err
		return result
	}
	patchPasswd := Lib.CryptPasswd(patchKey, string(passwordHash), -1)
	//读取文件类型
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	data, err = Lib.DecryptData(docPasswd, data)
	if err != nil {
		result.Err = err
		return result
	}
	if string(data) != "NNIC" {
		result.Err = errors.New("未知文件类型")
		return result
	}
	//读取程序名称
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	data, err = Lib.DecryptData(docPasswd, data)
	if err != nil {
		result.Err = err
		return result
	}
	if string(data) != Lib.VI_APP_NAME {
		result.Err = errors.New("未知程序文件")
		return result
	}
	//读取程序版本
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	data, err = Lib.DecryptData(docPasswd, data)
	if err != nil {
		result.Err = err
		return result
	}
	app_version := string(data)
	v1, err := version.NewVersion(app_version)
	if err != nil {
		result.Err = err
		return result
	}
	v2, err := version.NewVersion(Lib.VI_VERSION)
	if err != nil {
		result.Err = err
		return result
	}
	if v2.LessThan(v1) {
		result.Err = errors.New("程序版本太低，请升级到最新版本！")
		return result
	}
	//读取程序版权
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	data, err = Lib.DecryptData(docPasswd, data)
	if err != nil {
		result.Err = err
		return result
	}
	if string(data) != Lib.VI_COPYRIGHT {
		result.Err = errors.New("未知程序所有邮箱")
		return result
	}
	//读取程序所有
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	data, err = Lib.DecryptData(docPasswd, data)
	if err != nil {
		result.Err = err
		return result
	}
	if string(data) != Lib.VI_EMAIL {
		result.Err = errors.New("未知程序所有邮箱")
		return result
	}
	//读取程序官网
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	data, err = Lib.DecryptData(docPasswd, data)
	if err != nil {
		result.Err = err
		return result
	}
	if string(data) != Lib.VI_URL {
		result.Err = errors.New("未知程序官网地址")
		return result
	}
	//读取问题报告
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	data, err = Lib.DecryptData(docPasswd, data)
	if err != nil {
		result.Err = err
		return result
	}
	if string(data) != Lib.VI_REPORT {
		result.Err = errors.New("未知问题报告地址")
		return result
	}
	//读取服务器公钥，暂时先不用
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	data, err = Lib.DecryptData(patchPasswd, data)
	if err != nil {
		result.Err = err
		return result
	}
	result.ServerKey = string(data)
	//读取视图对象
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	obj_data := data
	//签名视图对象
	sig_result := Lib.SignData(data, appPrivateKey)
	if sig_result.Err != nil {
		result.Err = err
		return result
	}
	//读取视图对象签名
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	signature, err := Lib.DecryptData(docPasswd, data)
	if err != nil {
		result.Err = err
		return result
	}
	//读取视图对象hash
	var data1 [32]byte
	err = decoder.Decode(&data1)
	if err != nil {
		result.Err = err
		return result
	}
	err = Lib.CheckSign(sig_result.Hashed, signature, &appPrivateKey.PublicKey)
	if err != nil {
		result.Err = err
		return result
	}
	err = Lib.CheckSign(data1, sig_result.Signature, &appPrivateKey.PublicKey)
	if err != nil {
		result.Err = err
		return result
	}

	data, err = Lib.DecryptData(patchPasswd, obj_data)
	if err != nil {
		result.Err = err
		return result
	}
	result.ViewData = string(data)
	//读取文件签名
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	file_signature, err := Lib.DecryptData(docPasswd, data)
	if err != nil {
		result.Err = err
		return result
	}
	//读取文件hash
	data1 = [32]byte{}
	err = decoder.Decode(&data1)
	if err != nil {
		result.Err = err
		return result
	}
	//读取文件数据
	data = []byte{}
	err = decoder.Decode(&data)
	if err != nil {
		result.Err = err
		return result
	}
	//签名文件内容
	sig_result = Lib.SignData(data, displayPrivateKey)
	if sig_result.Err != nil {
		result.Err = err
		return result
	}
	err = Lib.CheckSign(sig_result.Hashed, file_signature, &displayPrivateKey.PublicKey)
	if err != nil {
		result.Err = err
		return result
	}
	err = Lib.CheckSign(data1, sig_result.Signature, &displayPrivateKey.PublicKey)
	if err != nil {
		result.Err = err
		return result
	}
	//合并密钥hash
	all_hash := sha256.Sum256([]byte(string(patchPasswd) + string(docPasswd)))
	//解码文件
	//需要生成第二次密码
	left_index := ((int(patchPasswd[index%10]) + int(all_hash[index%10])) % 10) * 5000
	patchPasswd2 := Lib.CryptPasswd(patchKey, string(passwordHash), left_index)
	//二次加密文件内容
	file_data, err := Lib.DecryptData(patchPasswd2, data)
	if err != nil {
		result.Err = err
		return result
	}
	//加密文件内容
	file_data, err = Lib.DecryptData(all_hash[:], file_data)
	if err != nil {
		result.Err = err
		return result
	}
	//解压
	reader, err := xz.NewReader(bytes.NewBuffer(file_data))
	if err != nil {
		result.Err = err
		return result
	}
	var buf bytes.Buffer
	if _, err = io.Copy(&buf, reader); err != nil {
		result.Err = err
		return result
	}
	result.DbData = buf.Bytes()
	//刷新KEY
	err = a.initAppKeys()
	if err != nil {
		result.Err = err
		return result
	}
	if result.ServerKey != "" {
		a.serverKey = result.ServerKey
	} else {
		a.InitServerKey()
	}
	return result
}
