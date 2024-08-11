package Lib

import (
	"bytes"
	"context"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"math/big"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/pbkdf2"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func ModifyFileName(oldPath, newFileName string) string {
	dir := filepath.Dir(oldPath)           // 获取旧路径的目录部分
	ext := filepath.Ext(oldPath)           // 获取旧路径的扩展名部分
	newBase := newFileName + ext           // 构建新的文件名
	newPath := filepath.Join(dir, newBase) // 组合新的路径

	return newPath
}

func CopyFile(sourceFile, destinationFile string) error {
	// 打开源文件
	source, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer source.Close()

	// 创建目标文件
	destination, err := os.Create(destinationFile)
	if err != nil {
		return err
	}
	defer destination.Close()

	// 复制文件内容
	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	return nil
}

// 前端可能有一部分内容无法跨域加载图片
// 需要用后端加载并转base64
func FetchImgToBase64(url string) string {
	// 下载图片
	response, err := http.Get(url)
	if err != nil {
		return "ERROR " + err.Error()
	}
	defer response.Body.Close()

	// 读取图片内容
	imageBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "ERROR " + err.Error()
	}

	// 将图片转换为 Base64
	base64String := base64.StdEncoding.EncodeToString(imageBytes)
	return base64String
}

func FetchImgToBase64Wails(url string) string {
	name := strings.TrimPrefix(url, "wails://wails")
	fileData, err := os.ReadFile("/tmp/MinView" + name)
	if err != nil {
		return ""
	}
	// 将图片转换为 Base64
	base64String := base64.StdEncoding.EncodeToString(fileData)
	return base64String
}

func NewImg(width, height int) string {
	// 创建一个 RGBA 图片对象，宽度为 500，高度为 500
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 将图片中的所有像素都设置为透明
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.SetRGBA(x, y, color.RGBA{0, 0, 0, 0})
		}
	}

	// 将图片保存到文件
	file_name := ReplaceFirstNChars(uuid.NewV4().String(), "e", 1)
	file, err := os.Create("/tmp/MinView/" + file_name + ".png")
	if err != nil {
		return err.Error()
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return err.Error()
	}
	return "wails://wails/" + file_name + ".png"
}

func WriteToCache(data_url string) string {
	file_name := ReplaceFirstNChars(uuid.NewV4().String(), "e", 1)
	parts := strings.Split(data_url, ",")
	if len(parts) != 2 {
		return "Invalid data URL"
	}

	imageData, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return err.Error()
	}

	// 解码 JPG 格式的图像数据
	img, err := jpeg.Decode(bytes.NewReader(imageData))
	if err != nil {
		return err.Error()
	}

	// 缩放图像
	// newImg := resize(img, 640, 480)

	// 将图像保存为 JPG 格式的文件
	outputFile, err := os.Create("/tmp/MinView/" + file_name + ".jpg")
	if err != nil {
		return err.Error()
	}
	defer outputFile.Close()

	// 将图像写入文件
	err = jpeg.Encode(outputFile, img, &jpeg.Options{Quality: 100})
	if err != nil {
		return err.Error()
	}
	return "wails://wails/" + file_name + ".jpg"
}

// // 缩放图像函数
// func resize(img image.Image, width, height int) image.Image {
// 	// 使用 Lanczos3 滤波器进行图像缩放
// 	newImg := image.NewRGBA(image.Rect(0, 0, width, height))
// 	imdraw.CatmullRom.Scale(newImg, newImg.Bounds(), img, img.Bounds(), draw.Src, nil)
// 	return newImg
// }

// 获取文件类型
func GetMIMEType(path string) (string, error) {
	ext := filepath.Ext(path)
	if ext == "" {
		return "", errors.New("未知文件扩展名")
	}
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		return "", errors.New("未知mime类型")
	}
	return mimeType, nil
}

// 拼接创建表的sql语句
func GenerateCreateTableSQL(tableName string, fieldNames []string, fieldTypes []string) string {
	var sb strings.Builder

	sb.WriteString("CREATE TABLE IF NOT EXISTS ")
	sb.WriteString("'" + tableName + "'")
	sb.WriteString(" (")

	for i := 0; i < len(fieldNames); i++ {
		sb.WriteString("'" + fieldNames[i] + "'")
		sb.WriteString(" ")
		sb.WriteString(fieldTypes[i])

		if i < len(fieldNames)-1 {
			sb.WriteString(", ")
		}
	}

	sb.WriteString(")")

	return sb.String()
}

// 拼接插入数据sql语句
func GenerateInsertSQL(tableName string, fieldNames []string, fieldTypes []string, values []string) string {
	var sb strings.Builder

	sb.WriteString("INSERT INTO ")
	sb.WriteString("'" + tableName + "'")
	sb.WriteString(" (")

	sb.WriteString("'" + strings.Join(fieldNames, "', '") + "'")

	sb.WriteString(") VALUES (")

	var valuePlaceholders []string
	for i := 0; i < len(fieldNames); i++ {
		var placeholder string

		if fieldTypes[i] == "TEXT" {
			placeholder = "'" + values[i] + "'"
		} else if fieldTypes[i] == "REAL" {
			if _, err := strconv.ParseFloat(values[i], 64); err == nil {
				placeholder = values[i] //strconv.FormatFloat(num, 'f', -1, 64)
			} else {
				placeholder = "NULL"
			}
		} else {
			placeholder = "NULL"
		}

		valuePlaceholders = append(valuePlaceholders, placeholder)
	}
	sb.WriteString(strings.Join(valuePlaceholders, ", "))

	sb.WriteString(")")

	return sb.String()
}

// 生成删除表
func GenerateDropSQL(tableName string) string {
	return "DROP TABLE IF EXISTS '" + tableName + "'"
}

// 备份数据库
func BackupDatabase(db *sql.DB) (string, error) {
	// 备份到临时文件
	tmpFile, err := ioutil.TempFile(``, `NNI-*`)
	if err != nil {
		return "", err
	}
	tmpFile.Close()
	ctx := context.Background()
	// 目的数据库
	dstDB, err := sql.Open(`sqlite3`, tmpFile.Name())
	if err != nil {
		return "", err
	}
	defer dstDB.Close()

	dstConn, err := dstDB.Conn(ctx)
	if err != nil {
		return "", err
	}
	defer dstConn.Close()
	if err := dstConn.Raw(func(dstDC interface{}) error {
		rawDstConn := dstDC.(*sqlite3.SQLiteConn)
		// 源数据库
		srcConn, err := db.Conn(ctx)
		if err != nil {
			return err
		}
		defer srcConn.Close()

		if err := srcConn.Raw(func(srcDC interface{}) error {
			rawSrcConn := srcDC.(*sqlite3.SQLiteConn)

			// 备份函数调用
			backup, err := rawDstConn.Backup(`main`, rawSrcConn, `main`)
			if err != nil {
				return err
			}

			// errors can be safely ignored.
			_, _ = backup.Step(-1)

			if err := backup.Close(); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return "", err
	}
	_, err = dstDB.Exec("VACUUM")
	if err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}

// 还原数据库,从file_path还原到db
func RestoreDatabase(file_path string, db *sql.DB) error {
	ctx := context.Background()
	// 目的数据库
	sourceDB, err := sql.Open(`sqlite3`, file_path)
	if err != nil {
		return err
	}
	defer sourceDB.Close()

	sourceConn, err := sourceDB.Conn(ctx)
	if err != nil {
		return err
	}
	defer sourceConn.Close()
	if err := sourceConn.Raw(func(sourceDC interface{}) error {
		rawSourceConn := sourceDC.(*sqlite3.SQLiteConn)

		// 源数据库
		distConn, err := db.Conn(ctx)
		if err != nil {
			return err
		}
		defer distConn.Close()

		if err := distConn.Raw(func(distDC interface{}) error {
			rawDistConn := distDC.(*sqlite3.SQLiteConn)

			// 备份函数调用
			backup, err := rawDistConn.Backup(`main`, rawSourceConn, `main`)
			if err != nil {
				return err
			}

			// errors can be safely ignored.
			_, _ = backup.Step(-1)

			if err := backup.Close(); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}

// 给数据签名
type DataSign struct {
	Hashed    [32]byte
	Signature []byte
	Err       error
}

// 使用私钥进行签名
func SignData(data []byte, privateKey *rsa.PrivateKey) DataSign {
	// 使用私钥进行签名
	var sig DataSign
	sig.Err = nil
	hashed := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		sig.Err = err
		return sig
	}
	sig.Hashed = hashed
	sig.Signature = signature
	return sig
}

// 检查数据签名
func CheckSign(hashed [32]byte, signature []byte, publicKey *rsa.PublicKey) error {
	// 使用公钥进行验证
	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
	return err
}

// 对称密钥加密数据 AES-256
func EncryptData(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

// 对称密钥解密数据 AES-256
func DecryptData(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("密文长度不足")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

// 生成大写字符串
func GenerateRandomString(length int) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		result[i] = letters[n.Int64()]
	}
	return string(result)
}

// 生成大写字符串和数字
func GenerateRandomStringN(length int) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		result[i] = letters[n.Int64()]
	}
	return string(result)
}

// 生成密码字符串
func GenerateRandomPasswd(length int) string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_+=~!@#$&.?"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		result[i] = letters[n.Int64()]
	}
	return string(result)
}

// 生成特定长度的随机字节数据
func GenerateRandomByte(length int) ([]byte, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	return randomBytes, nil
}

// 密码加密,iterations迭代次数
func CryptPasswd(salt []byte, pwd string, iterations int) []byte {
	keyLength := 32 // 密钥长度为 32 字节
	iter := iterations
	// 迭代次数
	if iter < 10000 {
		iter = 10000
	}
	// 使用 PBKDF2 算法生成密钥
	key := pbkdf2.Key([]byte(pwd), salt, iter, keyLength, sha256.New)
	return key
}

// RSA公钥加密
func EncryptRSAData(publicKey *rsa.PublicKey, originalMessage []byte) ([]byte, error) {
	// 使用公钥进行加密
	// encryptedMessage, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, originalMessage)
	// if err != nil {
	// 	return nil, err
	// }
	// 使用公钥进行加密
	encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, originalMessage, nil)
	if err != nil {
		return nil, err
	}
	return encryptedData, nil
}

// RSA私钥解密
func DecryptRSAData(privateKey *rsa.PrivateKey, encryptedMessage []byte) ([]byte, error) {
	// 使用私钥进行解密
	// decryptedMessage, err := rsa.DecryptPKCS1v15(nil, publicKey, encryptedMessage)
	// if err != nil {
	// 	return nil, err
	// }
	// 使用私钥进行解密
	decryptedData, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedMessage, nil)
	if err != nil {
		return nil, err
	}
	return decryptedData, nil
}

// 私钥转换器，转成PEM
func PEM2PrivateKey(privatePEM []byte) (*rsa.PrivateKey, error) {
	// 从PEM格式还原私钥
	block, _ := pem.Decode([]byte(privatePEM))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("无效的PEM格式私钥")
	}
	parsedPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return parsedPrivateKey, nil
}

// base64加密
func EncryptBase64(plantmsg []byte) string {
	return base64.StdEncoding.EncodeToString(plantmsg)
}

// base64解密
func DecryptBase64(cryptmsg string) ([]byte, error) {
	msg, err := base64.StdEncoding.DecodeString(cryptmsg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// 把前n个字符给替换
func ReplaceFirstNChars(input, replacement string, n int) string {
	// 获取原始字符串的前n个字符
	prefix := input[:n]
	// 替换前n个字符
	newString := strings.Replace(input, prefix, replacement, 1)
	return newString
}
