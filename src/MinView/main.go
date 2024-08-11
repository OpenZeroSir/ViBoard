package main

import (
	"MinView/Lib"
	"database/sql"
	"embed"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	systemruntime "runtime"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

// timingMiddleware 记录请求处理时间的中间件
func timingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println("----", r.RemoteAddr, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

type FileLoader struct {
	http.Handler
}

func NewFileLoader() *FileLoader {
	return &FileLoader{}
}

// func NewFileLoader() http.Handler {

// 	assetsPath := filepath.Join(".", "assets")

// 	// 创建处理程序
// 	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// 获取请求的文件路径
// 		filePath := filepath.Join(assetsPath, r.URL.Path)

// 		// 加载静态文件并提供给应用程序
// 		http.ServeFile(w, r, filePath)
// 	})
// 	return handler

// }

func (h *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	//处理其他请求
	newurl := req.URL.Path
	// fmt.Println(req.URL.String())
	if newurl == "/baidu/api" {
		newurl = "https://api.map.baidu.com/api?" + req.URL.RawQuery
		resp, err := http.Get(newurl)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("Error:" + err.Error()))
			return
		}
		defer resp.Body.Close()

		// 读取API响应的内容
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("Error:" + err.Error()))
			return
		}
		res.WriteHeader(http.StatusOK)
		// data := strings.ReplaceAll(string(body), "https://", "wails://")
		var data string
		switch systemruntime.GOOS {
		case "darwin":
			data = strings.ReplaceAll(string(body), "https://", "wails://")
		case "linux":
			data = strings.ReplaceAll(string(body), "https://", "wails://")
		case "windows":
			data = strings.ReplaceAll(string(body), "https://", "http://wails.localhost/")
		default:
			// return "未知系统类型"
		}
		// data := strings.ReplaceAll(string(body), "https://", "wails://")
		res.Write([]byte(data))
		return
	}
	newurl = req.URL.String()
	if strings.Contains(newurl, "baidu.com") || strings.Contains(newurl, "bdimg.com") {
		newurl = strings.ReplaceAll(newurl, "wails://", "https://")
		newurl = strings.ReplaceAll(newurl, "http://wails.localhost/", "https://")
		resp, err := http.Get(newurl)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("Error:" + err.Error()))
			return
		}
		defer resp.Body.Close()

		// 读取API响应的内容
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("Error:" + err.Error()))
			return
		}
		res.WriteHeader(http.StatusOK)
		// data := strings.ReplaceAll(string(body), "https://", "wails://")
		var data string
		switch systemruntime.GOOS {
		case "darwin":
			data = strings.ReplaceAll(string(body), "https://", "wails://")
		case "linux":
			data = strings.ReplaceAll(string(body), "https://", "wails://")
		case "windows":
			data = strings.ReplaceAll(string(body), "https://", "http://wails.localhost/")
		default:
			// return "未知系统类型"
		}
		res.Write([]byte(data))
		return
	}
	newurl = strings.ReplaceAll(newurl, "wails://wails/", "/")
	newurl = strings.ReplaceAll(newurl, "http://wails.localhost/", "/")
	if newurl[:8] == "/NIBoard" {
		newurl = "/" + newurl[17:]
	}
	urls := strings.Split(newurl, "NNI")
	if len(urls) < 2 {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf("FileLoader ServeHTTP len(urls)<0 %s", req.URL.Path)))
		return
	}
	comid_ext := urls[len(urls)-1]
	coms := strings.Split(comid_ext, ".")
	comid := coms[0]
	if comid == "-" {
		comid = ""
	}
	vieid := urls[len(urls)-2][1:]
	// 查询数据
	selectData := "SELECT Name,Data FROM MEDIA WHERE ViewID='" + vieid + "' AND ComID='" + comid + "' LIMIT 1;"
	// fmt.Println(selectData)
	var rows *sql.Rows
	var err error
	if app.appmode == "/" {
		rows, err = app.memsql.GetDB().Query(selectData)
	} else {
		rows, err = app.displaysql.GetDB().Query(selectData)
	}
	if err != nil {
		if app.appmode == "/" {
			app.memsql.Mutex.Unlock()
		} else {
			app.displaysql.Mutex.Unlock()
		}
		// fmt.Println("无法查询数据:", err)
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf("Error:ServeHTTP Select %s failed", req.URL.Path)))
		return
	}
	// defer rows.Close()

	state := rows.Next()
	if !state {
		if app.appmode == "/" {
			if e := rows.Close(); e != nil {
				app.memsql.Mutex.Unlock()
				res.WriteHeader(http.StatusBadRequest)
				res.Write([]byte(fmt.Sprintf("Error:ServeHTTP %s %s", req.URL.Path, e.Error())))
				return
			}
			app.memsql.Mutex.Unlock()
		} else {
			if e := rows.Close(); e != nil {
				app.displaysql.Mutex.Unlock()
				res.WriteHeader(http.StatusBadRequest)
				res.Write([]byte(fmt.Sprintf("Error:ServeHTTP %s %s", req.URL.Path, e.Error())))
				return
			}
			app.displaysql.Mutex.Unlock()
		}
		// fmt.Println("查询不到数据")
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf("Error:ServeHTTP No %s file", req.URL.Path)))
		return
	}
	var name string
	var data []byte
	err = rows.Scan(&name, &data)
	if err != nil {
		if app.appmode == "/" {
			if e := rows.Close(); e != nil {
				app.memsql.Mutex.Unlock()
				res.WriteHeader(http.StatusBadRequest)
				res.Write([]byte(fmt.Sprintf("Error:ServeHTTP %s %s", req.URL.Path, e.Error())))
				return
			}
			app.memsql.Mutex.Unlock()
		} else {
			if e := rows.Close(); e != nil {
				app.displaysql.Mutex.Unlock()
				res.WriteHeader(http.StatusBadRequest)
				res.Write([]byte(fmt.Sprintf("Error:ServeHTTP %s %s", req.URL.Path, e.Error())))
				return
			}
			app.displaysql.Mutex.Unlock()
		}
		// fmt.Println("Error:无法读取数据", err)
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf("Error:ServeHTTP Could not read file %s", req.URL.Path)))
		return
	}
	err = rows.Close()
	if err != nil {
		if app.appmode == "/" {
			app.memsql.Mutex.Unlock()
		} else {
			app.displaysql.Mutex.Unlock()
		}
		// fmt.Println("Error:无法读取数据", err)
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf("Error:ServeHTTP Could not read file %s", req.URL.Path)))
		return
	}
	if app.appmode == "/" {
		app.memsql.Mutex.Unlock()
	} else {
		app.displaysql.Mutex.Unlock()
	}
	// fileData, err := os.ReadFile("/tmp/MinView" + req.URL.Path)
	// if err != nil {
	// 	fmt.Println("FileLoader ServeHTTP ReadFile 文件读取失败", err)
	// 	res.WriteHeader(http.StatusBadRequest)
	// 	res.Write([]byte(fmt.Sprintf("FileLoader ServeHTTP Could not load file %s", req.URL.Path)))
	// 	return
	// }
	// fmt.Println(fileData)
	mimeType, err := Lib.GetMIMEType(name)
	if err != nil {
		// fmt.Println("FileLoader ServeHTTP GetMIMEType", err)
		res.Header().Set("Content-Type", "")
	} else {
		// fmt.Println("FileLoader ServeHTTP GetMIMEType", mimeType)
		res.Header().Add("Content-Type", mimeType)
		res.Header().Add("Access-Control-Expose-Headers", "Content-Type")
	}
	// fmt.Println(res.Header())
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(data)
	if err != nil {
		fmt.Println("FileLoader ServeHTTP Write 失败", err)
	}
	// }
}

// 支持跨域请求
// func basicCors() *cors.Cors {
// 	return cors.New(cors.Options{
// 		AllowedOrigins:   []string{"https://miao.baidu.com", "https://api.map.baidu.com"},
// 		AllowCredentials: true,
// 		Debug:            true,
// 	})
// }

// func basicCors() *cors.Cors {
// 	corsInstance := cors.New(cors.Options{
// 		AllowedOrigins:   []string{"https://miao.baidu.com", "https://api.map.baidu.com"},
// 		AllowCredentials: true,
// 		Debug:            true,
// 	})
// 	fmt.Println("CORS Middleware Applied Successfully")
// 	return corsInstance
// }

func main() {

	// Create an instance of the app structure
	app := NewApp()
	// Create application with options
	err := wails.Run(&options.App{
		Title:  "MinView",
		Width:  1280,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets:     assets,
			Handler:    NewFileLoader(),
			Middleware: timingMiddleware,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnBeforeClose:    app.beforeClose,
		Fullscreen:       false,
		Bind: []interface{}{
			app,
			// basicCors(),
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
