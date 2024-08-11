package main

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"MinView/Core"
	"MinView/Editor"
	"MinView/Lib"
	systemruntime "runtime"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
	//前端显示的是编辑模式还是显示模式,"/"是编辑模式，这是默认程序打开进入编辑模式，/display是显示模式
	appmode string
	//程序静态资源端口
	port int
	//web服务状态
	webstatus bool
	//加密密钥
	key *Core.KEY
	//编辑器内存数据库
	memsql *Editor.EditorSQL
	//远程数据库
	displaysql *Editor.EditorSQL
	//显示器加载视图时，保存视图信息，当另存视图或保存视图时，用这个保存，最大可能避免视图信息被修改。
	displayinfo string
	//远程登陆信息
	rs *Core.RSession
	//编辑器web服务，提交的附件需要保存到本地和数据库，并生成一个链接
	//供视图加载
	editorweb *Editor.EditorWeb
	//远程web服务
	uiserver *http.Server
	//显示服务的组件对象,给web ui加载组件用
	displaycompo string
	//是否已经加载地图api,一般是没有网络的时候加载不了百度地图api
	mapapiloaded bool
}

var app App

// NewApp creates a new App application struct
func NewApp() *App {
	return &app
}
func GetApp() *App {
	return &app
}

//go:embed templates/static
var tmpl embed.FS

// 设置程序显示模式
func (a *App) SetupAppMode(mode string) {
	a.appmode = mode
}

// 检查远程cookie中间件
func (a *App) checkCookieMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")                   // 允许来自任何源的请求
		rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // 允许的HTTP方法
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")       // 允许的请求头
		cookie, err := r.Cookie("NNI")
		if err != nil || cookie.Value == "" {
			http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
			return
		}

		ok, err := a.rs.CheckCookie(cookie.Value)
		if err != nil || !ok {
			http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
			return
		}
		next(rw, r)
	}
}

// 检查远程token中间件
func (a *App) checkTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")                   // 允许来自任何源的请求
		rw.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // 允许的HTTP方法
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type")       // 允许的请求头

		token := r.Header.Get("token")
		if token == "" {
			rw.Write([]byte("请配置token属性"))
			return
		}
		ok := a.rs.CheckToken(token)
		if !ok {
			rw.Write([]byte("token检查失败"))
			return
		}
		next(rw, r)
	}
}

// ui登陆成功跳转
func (a *App) nniHandler(rw http.ResponseWriter, r *http.Request) {
	http.Redirect(rw, r, "/NNI/templates/static", http.StatusTemporaryRedirect)
}

// ui get提交
// func (a *App) nniGetHandler(rw http.ResponseWriter, r *http.Request) {
// 	decoder := json.NewDecoder(r.Body)
// 	var params map[string]string
// 	err := decoder.Decode(&params)
// 	if err != nil {
// 		rw.WriteHeader(401)
// 		return
// 	}
// 	fmt.Println(params)
// 	rw.Write([]byte("ok!"))
// }

// 开发者提交 post
func (a *App) nniDevPostHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rw.WriteHeader(405) // 返回错误代码
		return
	}
	//解析 FormData 数据
	if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		err := r.ParseMultipartForm(10 << 20) // 设置合适的内存限制
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
	}
	post_key := r.FormValue("key")
	switch len(post_key) {
	case 0:
		//未配置key
		rw.Write([]byte("请正确配置key字段"))
	case 36: //表的key长度
		replace := r.FormValue("replace")
		if replace == "" {
			replace = "false"
		}
		replace = strings.ToUpper(replace)
		is_replace := false
		if replace == "TRUE" {
			is_replace = true
		}
		data := r.FormValue("data")
		if len(data) == 0 {
			rw.Write([]byte("请求应当包含data的内容"))
			return
		}
		err := a.displaysql.RemotePostUpdateDB(post_key, is_replace, data)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		items := make(map[string]string)
		//因为没有表名，所以随便给一个表名，在前端识别到表名就更新数据
		items[post_key] = "UpdateFromDev"
		runtime.EventsEmit(a.ctx, Lib.UIEmitSheetData, items)
		rw.Write([]byte("ok"))
	case 21: //组件key的长度,nni开头，例如nni3SPMaDm5kz8ID4Aftz16X
		data := r.FormValue("data")
		if len(data) == 0 {
			rw.Write([]byte("请求应当包含data的内容"))
			return
		}
		// fmt.Println(data)
		//发送数据到显示器，显示器监听事件获得数据会刷新组件
		//数据应该是数组类数据的string,例如："[选项1,选项2]"
		//组件id都是nni开头
		//需要把key跟data一起传送给，显示器接收到数据能分辨出是哪个组件
		map_data := make(map[string]string)
		map_data["key"] = "nni" + post_key
		map_data["data"] = data
		runtime.EventsEmit(a.ctx, "nni"+post_key, map_data)
		rw.Write([]byte("ok"))
	default:
		rw.Write([]byte("未知操作，程序未支持。"))
	}
}

// ui提交post
func (a *App) nniPostHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rw.WriteHeader(405) // 返回错误代码
		return
	}
	//解析 FormData 数据
	if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		err := r.ParseMultipartForm(10 << 20) // 设置合适的内存限制
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
	}
	// err := r.ParseForm()
	// if err != nil {
	// 	http.Error(rw, "Failed to parse form data", http.StatusBadRequest)
	// 	return
	// }
	// post_type := r.Form.Get("type")
	// decoder := json.NewDecoder(r.Body)
	// var params map[string]string
	// err := decoder.Decode(&params)
	// if err != nil {
	// 	rw.WriteHeader(401)
	// 	return
	// }
	// rw.Header().Set("Content-Type", "application/json")
	post_type := r.FormValue("type")
	// post_type := r.MultipartForm.Value["type"][0]
	switch post_type {
	case Lib.PostTypeDisplayGetVolume:
		lst, err := a.displaysql.GetVOLUME()
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		jsonData, err := json.Marshal(lst)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write(jsonData)
	case Lib.PostTypeDisplayGetComponent:

		rw.Write([]byte(a.displaycompo))
	// case Lib.PostTypeImporData:
	// 	data := r.FormValue("Data")
	// 	var new_data []interface{}

	// 	// 解析收到的字符串为JSON格式
	// 	err := json.Unmarshal([]byte(data), &new_data)
	// 	if err != nil {
	// 		fmt.Println("Error parsing JSON:", err)
	// 		return
	// 	}
	// 	fmt.Println(new_data)

	case Lib.PostTypeImportSheet:
		file, handler, err := r.FormFile("File")
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		tmp_path := filepath.Join("/tmp", handler.Filename)
		tempFile, err := os.Create(tmp_path)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		// defer tempFile.Close()
		// 将上传的文件内容复制到临时文件
		_, err = io.Copy(tempFile, file)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		err = tempFile.Close()
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		key := r.FormValue("key")
		if key == "" {
			rw.Write([]byte("数据表KEY不能为空"))
			return
		}
		replace := r.FormValue("replace")
		if replace == "" {
			replace = "false"
		}
		replace = strings.ToUpper(replace)
		is_replace := false
		if replace == "TRUE" {
			is_replace = true
		}
		xlsx := Editor.NewEditorXlsx()
		err = xlsx.Open(tmp_path)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		items, err := a.displaysql.GetVOLUMEInfoByKey(key)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		err = xlsx.UpdateDataToDB(items, a.displaysql, is_replace)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		//emit表名和key给display前端用来刷新echarts,
		//items map[49584f95-e6ec-4a86-82ab-9d84b6ebba82:Sheet1]
		//前端只用到key,但value也一起传过去了
		runtime.EventsEmit(a.ctx, Lib.UIEmitSheetData, items)
		rw.Write([]byte("ok"))
	case Lib.PostTypeUpdateComponent:
		data := r.FormValue("data")
		if len(data) == 0 {
			rw.Write([]byte("请求应当包含data的内容"))
			return
		}
		//发送数据到显示器，显示器监听事件获得数据会刷新组件
		//数据应该是数组类数据的string,例如："[选项1,选项2]"
		//组件id都是nni开头
		//需要把key跟data一起传送给，显示器接收到数据能分辨出是哪个组件
		map_data := make(map[string]string)
		compo_key := r.FormValue("key")
		map_data["key"] = compo_key
		map_data["data"] = data
		compo_type := r.FormValue("comtype")
		map_data["comtype"] = compo_type
		runtime.EventsEmit(a.ctx, compo_key, map_data)
		rw.Write([]byte("ok"))
	case Lib.PostTypeUploadFile:
		file, handler, err := r.FormFile("File")
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		defer file.Close()
		// fmt.Println(handler.Filename)
		// 读取文件内容
		fileContent, err := io.ReadAll(file)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		db := app.displaysql.GetDB()
		// 将文件内容作为Blob类型保存到SQLite3数据库中
		stmt, err := db.Prepare("INSERT OR REPLACE INTO MEDIA values(?, ?, ?, ?)")
		if err != nil {
			app.displaysql.Mutex.Unlock()
			rw.Write([]byte(err.Error()))
			return
		}
		// defer stmt.Close()

		viewid := r.FormValue("viewid")
		compid := r.FormValue("compid")
		// tx, err := db.Begin()
		// if err != nil {
		// 	// 处理错误
		// 	if err != nil {
		// 		rw.Write([]byte(err.Error()))
		// 		return
		// 	}
		// }

		_, err = stmt.Exec(viewid, compid, handler.Filename, fileContent)
		if err != nil {
			// tx.Rollback()
			e := stmt.Close()
			if e != nil {
				app.displaysql.Mutex.Unlock()
				rw.Write([]byte(e.Error()))
				return
			}
			app.displaysql.Mutex.Unlock()
			rw.Write([]byte(err.Error()))
			return
		}
		// // 提交事务
		// err = tx.Commit()
		// if err != nil {
		// 	// 处理错误
		// 	// 注意：如果提交事务失败，数据库可能处于不一致状态，需要谨慎处理
		// 	tx.Rollback()
		// 	rw.Write([]byte(err.Error()))
		// 	return
		// }
		//base := "wails://wails/"
		// mimeType, err := Lib.GetMIMEType(handler.Filename)
		// if err != nil {
		// 	mimeType = ""
		// }
		e := stmt.Close()
		if e != nil {
			app.displaysql.Mutex.Unlock()
			rw.Write([]byte(e.Error()))
			return
		}
		app.displaysql.Mutex.Unlock()
		result := viewid + "NNI" + compid + "." + handler.Filename
		rw.Write([]byte(result))
		return
	default:
		rw.Write([]byte("未知操作,程序不支持。"))
	}
}

// 启动远程web服务
func (a *App) StartRemoteWeb(ip string) string {
	//初始化一个端口用于前端访问，需要循环检测端口是否被占用，如果占用了port++,默认50520，打开多实例可能会被占用。
	a.port = 50520
	count := 100
	for count > 0 {
		listener, err := net.Listen("tcp", ip+":"+strconv.Itoa(a.port))
		if err != nil {
			fmt.Printf("Port %d is already in use\n", a.port)
			a.port++
		} else {
			listener.Close()
			fmt.Printf("Port %d is available\n", a.port)
			break
		}
		count--
	}
	a.uiserver = &http.Server{Addr: ip + ":" + strconv.Itoa(a.port)}
	result := "启动成功:" + strconv.Itoa(a.port)
	go func() {
		a.webstatus = true
		if err := a.uiserver.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			result = err.Error()
			a.webstatus = false
			return
		}
	}()
	return result
}

// 停止远程web服务
func (a *App) StopRemoteWeb() string {
	result := "停止成功"
	if a.uiserver != nil {
		// 创建一个超时上下文
		// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// defer cancel()
		// // 关闭服务器
		// if err := a.uiserver.Shutdown(ctx); err != nil {
		// 	fmt.Printf("Server shutdown failed: %v\n", err)
		// } else {
		// 	fmt.Println("Server stopped.")
		// }
		if err := a.uiserver.Close(); err != nil {
			result = err.Error()
		}
	}
	a.webstatus = false
	return result
}

func (a *App) StatusRemoteWeb() int {
	if a.webstatus {
		return a.port
	}
	return -1
}

//坚持远程web服务状态

// 注册前端UI服务
func (a *App) RegisterUI() {
	t, _ := template.ParseFS(tmpl, "templates/static/nnilogin.html")
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		//首页跳转到nnilogin.html，
		t.ExecuteTemplate(rw, "nnilogin.html", nil)
	})

	//nnilogin.html提交登陆，检查账号密码，成功的花跳转/nni
	http.HandleFunc("/checklogin", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			rw.WriteHeader(405) // 返回错误代码
			return
		}
		decoder := json.NewDecoder(r.Body)
		var params map[string]string
		err := decoder.Decode(&params)
		if err != nil {
			rw.WriteHeader(401)
			return
		}
		token, err := a.rs.GetCookie(params["username"], params["password"])
		if err != nil {
			rw.WriteHeader(401)
			return
		}
		cookie1 := &http.Cookie{Name: "NNI", Value: token, HttpOnly: true}
		r.AddCookie(cookie1)
		http.SetCookie(rw, cookie1)
		http.Redirect(rw, r, "/nni", http.StatusTemporaryRedirect)
	})
	//登陆成功后跳转页面
	http.HandleFunc("/nni", a.checkCookieMiddleware(a.nniHandler))
	// http.HandleFunc("/nniget", a.checkCookieMiddleware(a.nniGetHandler))
	http.HandleFunc("/nnipost", a.checkCookieMiddleware(a.nniPostHandler))
	//程序远程提交数据接口
	http.HandleFunc("/ai", a.checkTokenMiddleware(a.nniDevPostHandler))
	fs := http.FileServer(http.FS(tmpl))
	hd := http.StripPrefix("/NNI/", fs)
	http.Handle("/NNI/", hd)
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	err := a.init_app(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// a.SendFileToFront()
}

// shutdown
// func (a *App) shutdown(ctx context.Context) {
// 	fmt.Println("shutdown")
// }

func (a *App) UIEmitComponentsData(optionalData ...interface{}) {
	data := fmt.Sprintf("%v", optionalData[0])
	a.displaycompo = data
}

// 发送文件给前端，
func (a *App) SendFileToFront(name string) string {
	//templates/static/chartoptions/brush.json
	data, err := tmpl.ReadFile(name)
	if err != nil {
		return "Error:" + err.Error()
	}
	return string(data)
}

func (a *App) init_app(ctx context.Context) error {
	a.ctx = ctx
	a.appmode = "/"
	a.webstatus = false
	a.mapapiloaded = false
	key, err := Core.NewKEY()
	if err != nil {
		return err
	}
	a.key = key
	memsql, err := Editor.NewSQL()
	if err != nil {
		return err
	}
	a.memsql = memsql
	displaysql, err := Editor.NewSQL()
	if err != nil {
		return err
	}
	a.displaysql = displaysql
	a.port = 50520
	a.rs = Core.NewRSession()

	a.uiserver = nil
	a.RegisterUI()
	a.displaycompo = ""
	// a.StartRemoteWeb()
	// a.editorweb = Editor.NewEditorWeb(a.port)
	// err = a.editorweb.Init()
	// if err != nil {
	// 	return err
	// }
	runtime.EventsOn(a.ctx, "UIEmitComponentsData", a.UIEmitComponentsData)
	// a.ReadKmlFile()
	return nil
}

func (a App) GetSQL() *Editor.EditorSQL {
	return a.memsql
}

// 前端需要获取到websocket端口,
// 因为端口可能被占用，会用其它端口打开
func (a App) GetSocketPort() int {
	return a.editorweb.GetPort()
}

// 前端该函数获取websocket token,
// 然后通过websocket发送token给后端验证是否是合法
func (a App) GetSocketToken() string {
	return a.editorweb.GetToken()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// type DisplayVOLUMEResultStruct struct {
// 	//表名和Key信息
// 	InfoList []string `json:"InfoList"`
// 	//错误信息
// 	Error string `json:"Error"`
// }

// 获取显示服务的表信息
// func (a *App) GetDisplayVOLUME() DisplayVOLUMEResultStruct {
// 	var v DisplayVOLUMEResultStruct
// 	v.Error = ""
// 	res, err := a.displaysql.GetVOLUME()
// 	if err != nil {
// 		v.Error = err.Error()
// 	}
// 	v.InfoList = res
// 	return v
// }

// 获取VOLUME的表信息
// @key是PriKey
// 如果key是空就获取PriKey和FileName
// 如果key不为空，获取SecKey和SheetName,
// 如果差不到，就可能是差key表的字段
func (a App) GetVOLUME(key string) map[string]string {
	// fmt.Println("GetVOLUME(key string)", key)
	res := make(map[string]string)
	res["Error"] = ""
	sql := ""
	if key == "" {
		sql = "SELECT DISTINCT PriKey,FileName FROM 'VOLUME'"
	} else {
		sql = "SELECT DISTINCT SecKey,SheetName FROM 'VOLUME' WHERE PriKey='" + key + "'"
	}
	rows, err := a.memsql.GetDB().Query(sql)
	if err != nil {
		a.memsql.Mutex.Unlock()
		res["Error"] = err.Error()
		return res
	}
	// defer rows.Close()
	// a.memsql.Mutex.Unlock()
	hasData := false
	// 遍历结果
	for rows.Next() {
		var (
			key  string
			name string
			// ...
		)
		err := rows.Scan(&key, &name /* ... */)
		if err != nil {
			a.memsql.Mutex.Unlock()
			res["Error"] = err.Error()
			return res
		}
		// 处理查询结果
		res[key] = name
		hasData = true
	}
	if err = rows.Err(); err != nil {
		if e := rows.Close(); e != nil {
			a.memsql.Mutex.Unlock()
			res["Error"] = err.Error()
			return res
		}
		a.memsql.Mutex.Unlock()
		res["Error"] = err.Error()
		return res
	}
	if e := rows.Close(); e != nil {
		res["Error"] = e.Error()
		a.memsql.Mutex.Unlock()
		return res
	}
	a.memsql.Mutex.Unlock()
	if !hasData && len(key) > 0 {
		rows, err := a.memsql.GetDB().Query("PRAGMA table_info('" + key + "')")
		if err != nil {
			a.memsql.Mutex.Unlock()
			res["Error"] = err.Error()
			return res
		}
		// defer rows.Close()
		// a.memsql.Mutex.Unlock()
		for rows.Next() {
			var (
				columnID   int
				columnName string
				columnType string
				va1        interface{}
				va2        interface{}
				va3        int
			)

			err := rows.Scan(&columnID, &columnName, &columnType, &va1, &va2, &va3)
			if err != nil {
				if e := rows.Close(); e != nil {
					res["Error"] = err.Error()
					a.memsql.Mutex.Unlock()
					return res
				}
				res["Error"] = err.Error()
				a.memsql.Mutex.Unlock()
				return res
			}
			res[columnName] = columnName + " " + columnType
		}
		if err = rows.Err(); err != nil {
			a.memsql.Mutex.Unlock()
			res["Error"] = err.Error()
			return res
		}
		a.memsql.Mutex.Unlock()
	}

	return res
}

// 根据PriKey删除表，因为一个PriKey可能对应多个SecKey,一个SecKey对应一个表，
// 所以需要遍历到PriKey对应的SecKey删除对应的表，再删除VOLUME表中PriKey是key的内容

func (a *App) DeleteVOLUME(key string) string {
	// 查询所有表名
	rr, xx := a.memsql.GetDB().Query("SELECT name FROM sqlite_master WHERE type='table';")
	if xx != nil {
		a.memsql.Mutex.Unlock()
		return xx.Error()
	}
	// defer rr.Close()
	// a.memsql.Mutex.Unlock()
	// 遍历查询结果
	for rr.Next() {
		var tableName string
		xx := rr.Scan(&tableName)
		if xx != nil {
			if e := rr.Close(); e != nil {
				a.memsql.Mutex.Unlock()
				return e.Error()
			}
			a.memsql.Mutex.Unlock()
			return xx.Error()
		}
		// fmt.Println(tableName)
	}
	if e := rr.Close(); e != nil {
		a.memsql.Mutex.Unlock()
		return e.Error()
	}
	a.memsql.Mutex.Unlock()
	// sql1 := "DROP TABLE 'MEDIA'"
	// aa, err1 := a.memsql.GetDB().Exec(sql1)
	// fmt.Println(err1, aa)
	sql := "SELECT PriKey,SecKey FROM 'VOLUME' WHERE PriKey='" + key + "'"
	rows, err := a.memsql.GetDB().Query(sql)
	if err != nil {
		a.memsql.Mutex.Unlock()
		return err.Error()
	}
	// defer rows.Close()
	// a.memsql.Mutex.Unlock()
	var drop_list []string
	for rows.Next() {
		var (
			priKey string
			secKey string
		)
		err = rows.Scan(&priKey, &secKey)
		if err != nil {
			if e := rows.Close(); e != nil {
				a.memsql.Mutex.Unlock()
				return e.Error()
			}
			a.memsql.Mutex.Unlock()
			return err.Error()
		}
		drop_list = append(drop_list, secKey)
		//不能在这里删除是因为可能是bug,rows没读完又删除表，提示找不到表。
	}
	if err = rows.Err(); err != nil {
		if e := rows.Close(); e != nil {
			a.memsql.Mutex.Unlock()
			return e.Error()
		}
		a.memsql.Mutex.Unlock()
		return err.Error()
	}
	if e := rows.Close(); e != nil {
		a.memsql.Mutex.Unlock()
		return e.Error()
	}
	a.memsql.Mutex.Unlock()
	for _, name := range drop_list {
		sql = "DROP TABLE `" + name + "`;"
		_, err = a.memsql.GetDB().Exec(sql)
		if err != nil {
			a.memsql.Mutex.Unlock()
			return err.Error()
		}
		a.memsql.Mutex.Unlock()
	}

	sql = "DELETE FROM 'VOLUME' WHERE PriKey='" + key + "'"
	_, err = a.memsql.GetDB().Exec(sql)
	if err != nil {
		a.memsql.Mutex.Unlock()
		return err.Error()
	}
	a.memsql.Mutex.Unlock()
	return ""
}

// 提供一个接口，用来获取文件路径，给别的接口处理
// 根据fileExt类型显示文件
// "*.jpg;*.png"
// 前端识别带有“Error:”的结果就认为是出错。
func (a *App) SelectFileDialog(fileExt string) string {
	// runtime.EventsEmit(a.ctx, "testEvent", "xxx")
	ext := make([]runtime.FileFilter, 1)
	ext[0].DisplayName = fileExt
	if strings.Contains(fileExt, "NNIE") {
		ext[0].DisplayName = "编辑视图(*.NNIE);编码视图(*.NNIC)"
	} else if strings.Contains(fileExt, "NNIC") {
		ext[0].DisplayName = "编码视图(*.NNIC)"
	}
	ext[0].Pattern = fileExt
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   "选择文件",
		Filters: ext,
	})
	if err != nil {
		return "Error:" + err.Error()
	}
	if len(path) == 0 {
		return "Error:没有选择文件"
	}
	return path
}

// 导入文件数据到数据库
// @replace确定是导入还是替换
func (a *App) ImportDataToDB(path string, replace bool) string {
	xlsx := Editor.NewEditorXlsx()
	err := xlsx.Open(path)
	if err != nil {
		return "Error:" + err.Error()
	}
	err = xlsx.CopyDataToDB(nil, a.memsql.GetDB(), replace)
	a.memsql.Mutex.Unlock()
	if err != nil {
		return "Error:" + err.Error()
	}
	return ""
}

// // 提供视图提交附件到数据库和Editor WEB服务，并获得下载地址。
// func (a *App) UploadFile(file os.File) string {
// 	fmt.Println("go UploadFile", file)
// 	return ""
// }

// 根据fileExt类型显示文件
// "*.jpg;*.png"
// 因为前端需要的是文件名的时候，byName是true,否则为false,前端在请求如视频、图片等，需要后端提供的mime类型，
// 所以false会返回携带mime类型
func (a *App) OpenFileDialog(fileExt string, viewid string, comid string, byName bool) string {

	ext := make([]runtime.FileFilter, 1)
	ext[0].DisplayName = fileExt
	ext[0].Pattern = fileExt
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   "选择文件",
		Filters: ext,
	})
	if err != nil {
		return "Error:" + err.Error()
	}
	if len(path) == 0 {
		return "Error:没有选择文件"
	}
	// 读取文件内容
	data, err := ioutil.ReadFile(path)
	if err != nil {
		// fmt.Println(err)
		return "Error:" + err.Error()
	}
	db := app.memsql.GetDB()
	// 将文件内容作为Blob类型保存到SQLite3数据库中
	stmt, err := db.Prepare("INSERT OR REPLACE INTO MEDIA values(?, ?, ?, ?)")
	if err != nil {
		a.memsql.Mutex.Unlock()
		return "Error:" + err.Error()
	}
	// defer stmt.Close()
	_, err = stmt.Exec(viewid, comid, filepath.Base(path), data)

	if err != nil {
		if e := stmt.Close(); e != nil {
			a.memsql.Mutex.Unlock()
			return err.Error()
		}
		a.memsql.Mutex.Unlock()
		return "Error:" + err.Error()
	}
	a.memsql.Mutex.Unlock()
	//base := "wails://wails/"
	mimeType, err := Lib.GetMIMEType(path)
	if err != nil {
		mimeType = ""
	}
	id := comid
	if id == "" {
		id = "-"
	}
	if byName {
		mimeType = filepath.Base(path)
	}
	return viewid + "NNI" + id + "." + mimeType
}

//需要提供一个函数用来给前端调用系统打开某个文件，比如双击附件可以打开文件。

func (a *App) OpenFile(viewid string, comid string) string {
	selectData := "SELECT Name,Data FROM MEDIA WHERE ViewID='" + viewid + "' AND ComID='" + comid + "' LIMIT 1;"
	// fmt.Println(selectData)
	var rows *sql.Rows
	var err error
	if a.appmode == "/" {
		rows, err = a.memsql.GetDB().Query(selectData)
		if err != nil {
			a.memsql.Mutex.Unlock()
			return err.Error()
		}
		// a.memsql.Mutex.Unlock()
	} else {
		rows, err = a.displaysql.GetDB().Query(selectData)
		if err != nil {
			a.displaysql.Mutex.Unlock()
			return err.Error()
		}
		// a.displaysql.Mutex.Unlock()
	}
	// if err != nil {
	// 	return "无法查询数据:" + err.Error()
	// }
	// defer rows.Close()
	state := rows.Next()
	if !state {
		if e := rows.Close(); e != nil {
			if a.appmode == "/" {
				a.memsql.Mutex.Unlock()
			} else {
				a.displaysql.Mutex.Unlock()
			}
			return e.Error()
		}
		if a.appmode == "/" {
			a.memsql.Mutex.Unlock()
		} else {
			a.displaysql.Mutex.Unlock()
		}
		return "查询不到数据"
	}
	var name string
	var data []byte
	err = rows.Scan(&name, &data)
	if err != nil {
		if e := rows.Close(); e != nil {
			if a.appmode == "/" {
				a.memsql.Mutex.Unlock()
			} else {
				a.displaysql.Mutex.Unlock()
			}
			return err.Error()
		}
		if a.appmode == "/" {
			a.memsql.Mutex.Unlock()
		} else {
			a.displaysql.Mutex.Unlock()
		}
		return "无法读取数据:" + err.Error()
	}
	if e := rows.Close(); e != nil {
		if a.appmode == "/" {
			a.memsql.Mutex.Unlock()
		} else {
			a.displaysql.Mutex.Unlock()
		}
		return e.Error()
	}
	if a.appmode == "/" {
		a.memsql.Mutex.Unlock()
	} else {
		a.displaysql.Mutex.Unlock()
	}
	file_path := os.TempDir() + "/" + name
	err = os.WriteFile(file_path, data, fs.ModePerm)
	if err != err {
		return "写入失败：" + err.Error()
	}
	command := "start"
	switch systemruntime.GOOS {
	case "darwin":
		command = "xdg-open"
	case "linux":
		command = "xdg-open"
	case "windows":
		command = "start"
	default:
		return "未知系统类型"
	}
	// fmt.Println("打开文件：" + file_path)
	cmd := exec.Command(command, file_path)
	err = cmd.Start()
	if err != nil {
		return "打开失败：" + err.Error()
	}
	return ""
}

// 请求保存文件，可能有路径，如果没有的选择文件
// ext是又前端传过来的，因为可能是工程文件，也有可能是视图文件
func (a *App) RequestSavePath(save_path string, file_ext string) string {
	exts := strings.Split(file_ext, ".")
	file_path := save_path
	if file_path == "" {
		display_name := "编辑视图(" + file_ext + ")"
		if exts[1] == "NNIC" {
			display_name = "编码视图(" + file_ext + ")"
		}
		ext := make([]runtime.FileFilter, 1)
		ext[0].DisplayName = display_name
		ext[0].Pattern = file_ext
		path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
			Title:   "保存文件",
			Filters: ext,
			// DefaultFilename: exts[1],
		})
		if err != nil {
			return "Error:" + err.Error()
		}
		if len(path) == 0 {
			return "Error:没有选择文件"
		}
		file_path = path
	}
	save_ext := filepath.Ext(file_path)
	if save_ext == "" {
		file_path = file_path + "." + exts[1]
	}
	return file_path
}

// 前端请求保存资源文件
func (a *App) RequestSaveProject(file_path string, data string) string {
	if len(file_path) < 3 {
		return "Error:文件路径有误"
	}
	if len(data) < 1 {
		return "Error:数据大小异常"
	}

	db_path, err := a.memsql.Backup()
	if err != nil {
		return "Error:" + err.Error()
	}
	file, err := os.Create(file_path)
	if err != nil {
		return "Error:" + err.Error()
	}
	err = a.key.EncryptProjectData(file, db_path, data)
	if err != nil {
		return "Error:" + err.Error()
	}
	return "保存完成"
}

// 打开项目文件
func (a *App) ReqestOpenProjectData(path string) string {
	res := a.key.DecryptProjectData(path)
	if res.Err != nil {
		return "Error:" + res.Err.Error()
	}
	err := a.memsql.Restore(res.DbData)
	if err != nil {
		return "Error:" + res.Err.Error()
	}
	return res.ViewData
}

// 新建视图工程
func (a *App) ReqestNewProject() string {
	//前端清理数据，后端清理数据库
	err := a.memsql.Clear()
	if err != nil {
		return err.Error()
	}
	return ""
}

// 导出视图编码，obj_data是视图对象信息，docinfo是文档信息
func (a *App) ReqestExportViewCode(file_path string, obj_data string, docinfo string, passwd string) string {
	if len(file_path) < 3 {
		return "Error:文件路径有误"
	}
	if len(obj_data) < 1 || len(docinfo) < 1 {
		return "Error:数据大小异常"
	}

	db_path, err := a.memsql.Backup()
	if err != nil {
		return "Error:" + err.Error()
	}
	file, err := os.Create(file_path)
	if err != nil {
		return "Error:" + err.Error()
	}
	pwd := passwd
	if pwd == "" {
		pwd = "123456"
	}
	err = a.key.ExportViewCode(file, db_path, obj_data, docinfo, pwd)
	if err != nil {
		return "Error:" + err.Error()
	}
	return "导出完成"
}

// 读取视图编码信息，因为在没有解码之前需要显示一些联系人之类的信息
func (a *App) RequestViewCodeInfo(file_path string) string {
	if len(file_path) < 3 {
		return "Error:" + "文件路径有误"
	}
	file, err := os.Open(file_path)
	if err != nil {
		return "Error:" + err.Error()
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	var docinfo []byte
	err = decoder.Decode(&docinfo)
	if err != nil {
		return "Error:" + err.Error()
	}
	data, err := Lib.DecryptData([]byte(Lib.VI_APP_KEYS[Lib.VI_INIT_VERSION]), []byte(docinfo))
	if err != nil {
		return "Error:" + err.Error()
	}
	return string(data)
}

type ImportViewCodeResult struct {
	//视图对象
	ObjData string `json:"ObjData"`
	//文档信息
	DocInfo string `json:"DocInfo"`
	//错误信息
	Error string `json:"Error"`
}

// 显示器导入视图编码
func (a *App) ReqestImportViewCode(path string, passwd string) ImportViewCodeResult {
	var result ImportViewCodeResult
	pwd := passwd
	if pwd == "" {
		pwd = "123456"
	}
	res := a.key.ImportViewCode(path, pwd)
	if res.Err != nil {
		result.Error = res.Err.Error()
		return result
	}
	result.DocInfo = res.ViewInfo
	a.displayinfo = res.ViewInfo
	result.ObjData = res.ViewData
	err := a.displaysql.Restore(res.DbData)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	// sss, err := a.displaysql.GetFieldAndTypeByKey("bb496da1-f87f-456c-b060-5c6b2305d242")
	// if err != nil {
	// 	result.Error = err.Error()
	// 	return result
	// }
	// fmt.Println(sss)
	return result
}

// 显示器保存或另存视图编码，obj_data是视图对象信息,和导出视图编码函数是一样流程的，只是数据库不一样
func (a *App) ReqestSaveViewCode(file_path string, obj_data string, passwd string) string {
	if len(file_path) < 3 {
		return "Error:文件路径有误"
	}
	if len(obj_data) < 1 || len(a.displayinfo) < 1 {
		return "Error:数据大小异常"
	}

	db_path, err := a.displaysql.Backup()
	if err != nil {
		return "Error:" + err.Error()
	}
	file, err := os.Create(file_path)
	if err != nil {
		return "Error:" + err.Error()
	}
	pwd := passwd
	if pwd == "" {
		pwd = "123456"
	}
	err = a.key.ExportViewCode(file, db_path, obj_data, a.displayinfo, pwd)
	if err != nil {
		return "Error:" + err.Error()
	}
	return "导出完成"
}

// 获取本机IP地址,可能是多个IP,主要用于远程访问,
// 格式 127.0.0.1：50520 。。。
func (a App) RequestIPAddress() []string {
	n := Core.NewNetInfo()
	ips, err := n.GetNetIPs()
	if err != nil {
		ips[0] = "Error:" + err.Error()
	}
	return ips
}

/*
远程登陆信息
*/

// 获取远程密码
func (a App) RequestGetRemotePasswd() string {
	return a.rs.GetPasswd()
}

// 获取远程token
func (a App) RequestGetRemoteToken() string {
	return a.rs.GetToken()
}

// 刷新远程信息，显示器加载一个文件刷新一次
func (a *App) ReqestResetRemoteSession() Core.RemoteSessionInfo {
	return a.rs.ResetRSession()
}

// 刷新远程登陆密码
func (a *App) ReqestResetRemoteSessionPasswd() string {
	return a.rs.ResetRSessionPasswd()
}

// 刷新远程token
func (a *App) ReqestResetRemoteSessionToken() string {
	return a.rs.ResetRSessionToken()
}

/*
下面内容可能是用不到的，暂时没有删除。
*/

// 前端可能有一部分内容无法跨域加载图片 需要用后端加载并转base64
func (a *App) FetchImgToBase64(url string) string {
	return Lib.FetchImgToBase64(url)
}
func (a *App) FetchImgToBase64Wails(url string) string {
	return Lib.FetchImgToBase64Wails(url)
}

// 由于缩略图刷新会导致程序变慢，只能先用生成一张图用来当缩略图
func (a *App) NewImg(width, height float64) string {
	return Lib.NewImg(int(width), int(height))
}

// 缩略图截图之后上传到缓存，并获得地址
func (a *App) WriteToCache(data_url string) string {
	return Lib.WriteToCache(data_url)
}

type QueryBySheetAndFieldResult struct {
	Error string `json:"error"`
	Data  []any  `json:"data"`
}

// 给前端编辑器发送数据
// 一般是通过表key和字段名称查询,
// 有些内容需要去重，需要通过duplicate指定
func (a *App) QueryBySheetAndField(sheet_key string, field_key []string, displaymode bool) QueryBySheetAndFieldResult {
	var db *sql.DB
	if displaymode {
		db = a.displaysql.GetDB()
	} else {
		db = a.memsql.GetDB()
	}
	var res QueryBySheetAndFieldResult
	if len(field_key) == 0 {
		if displaymode {
			a.displaysql.Mutex.Unlock()
		} else {
			a.memsql.Mutex.Unlock()
		}
		res.Error = "field_key empty"
		return res
	}
	str := ""
	for _, v := range field_key {
		str = str + "[" + v + "],"
	}

	str = str[:len(str)-1]

	sql_str := "select " + str + " from \"" + sheet_key + "\" "
	// fmt.Println(sql_str)
	rows, err := db.Query(sql_str)
	if err != nil {
		if displaymode {
			a.displaysql.Mutex.Unlock()
		} else {
			a.memsql.Mutex.Unlock()
		}
		res.Error = err.Error()
		return res
	}
	// defer rows.Close()
	var infos []any
	for rows.Next() {
		values := make([]interface{}, len(field_key))
		for i := range field_key {
			var v interface{}
			values[i] = &v
		}

		err := rows.Scan(values...)
		if err != nil {
			if e := rows.Close(); e != nil {
				if displaymode {
					a.displaysql.Mutex.Unlock()
				} else {
					a.memsql.Mutex.Unlock()
				}
				res.Error = err.Error()
				return res
			}
			if displaymode {
				a.displaysql.Mutex.Unlock()
			} else {
				a.memsql.Mutex.Unlock()
			}
			res.Error = err.Error()
			return res
		}
		// if len(values) == 1 {
		// 	infos = append(infos, values[0])
		// } else {
		// 	infos = append(infos, values)
		// }
		infos = append(infos, values)
	}
	res.Data = infos
	if e := rows.Close(); e != nil {
		if displaymode {
			a.displaysql.Mutex.Unlock()
		} else {
			a.memsql.Mutex.Unlock()
		}
		res.Error = e.Error()
		return res
	}
	if displaymode {
		a.displaysql.Mutex.Unlock()
	} else {
		a.memsql.Mutex.Unlock()
	}
	// fmt.Println("res", res)
	return res
}

// 避免前端跨域求情
func (a *App) RequestBaiduAPI(baidu_key string) string {
	resp, err := http.Get("https://api.map.baidu.com/api?v=3.0&ak=" + baidu_key)
	if err != nil {
		return "Error:" + err.Error()
	}
	defer resp.Body.Close()

	// 读取API响应的内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Error:" + err.Error()
	}
	//定义正则表达式模式
	pattern := `src="([^"]+)"`
	// 编译正则表达式
	regExp := regexp.MustCompile(pattern)
	// 使用正则表达式查找匹配的部分
	matches := regExp.FindStringSubmatch(string(body))
	// 输出匹配到的 URL
	if len(matches) < 1 {
		return "Error:No URL match found."
	}
	url := matches[1]
	// newresp, err := http.Get(url)
	// if err != nil {
	// 	return "Error:" + err.Error()
	// }
	// defer newresp.Body.Close()

	// // 读取API响应的内容
	// newbody, err := ioutil.ReadAll(newresp.Body)
	// if err != nil {
	// 	return "Error:" + err.Error()
	// }
	// data := strings.ReplaceAll(string(url), "https://", "wails://")
	// data := strings.ReplaceAll(string(url), "https://", "wails://")
	data := ""
	switch systemruntime.GOOS {
	case "darwin":
		data = strings.ReplaceAll(string(url), "https://", "wails://")
	case "linux":
		data = strings.ReplaceAll(string(url), "https://", "wails://")
	case "windows":
		data = strings.ReplaceAll(string(url), "https://", "http://wails.localhost/")
	default:
		// return "未知系统类型"
	}
	a.mapapiloaded = true
	return data
}

// 检查地图api是否已经加载
func (a *App) CheckMapAPILoaded() bool {
	return a.mapapiloaded
}

func NewCoordinates(longitude, latitude, altitude float64) string {
	return fmt.Sprintf("%f,%f,%f", longitude, latitude, altitude)
}
func NewCoordinatesOffset(longitude, latitude float64) string {
	return fmt.Sprintf("%f,%f,(AltitudeAboveSeaLevelOffset+1)", longitude, latitude)
}

type Geo struct {
	Name        string `json:"name"`
	CoordString string `json:"CoordString"`
	Color       string `json:"Color"`
	Fill        string `json:"Fill"`
}

// 只解析这几种类型，MultiGeometry暂时不解析
type GeoMessage struct {
	//点
	Points []Geo `json:"Points"`
	//线
	LineStrings []Geo `json:"LineStrings"`
	//环
	LinearRings []Geo `json:"LinearRings"`
	//错误信息
	Error string `json:"Error"`
}

func (a *App) SetGeoStyle(g *Geo, value Lib.Placemark, styles []Lib.MainStyle) {
	for _, style := range styles {
		if style.Id == Lib.KMLReplaceString(strings.Replace(value.StyleUrl, "#", "", 1)) {
			if style.LineStyle.Color != "" {
				g.Color = Lib.KMLReplaceString(style.LineStyle.Color)
			}
			if style.LineStyle.Fill != "" {
				g.Fill = Lib.KMLReplaceString(style.LineStyle.Fill)
			}
			if style.PolyStyle.Color != "" {
				g.Fill = Lib.KMLReplaceString(style.PolyStyle.Color)
			}
			if style.PolyStyle.Fill != "" {
				g.Fill = Lib.KMLReplaceString(style.PolyStyle.Fill)
			}
			if style.PointStyle.Color != "" {
				g.Color = Lib.KMLReplaceString(style.PointStyle.Color)
			}
			if style.PointStyle.Fill != "" {
				g.Fill = Lib.KMLReplaceString(style.PointStyle.Fill)
			}
			if value.Style.LineStyle.Color != "" {
				g.Color = Lib.KMLReplaceString(value.Style.LineStyle.Color)
			}
			if value.Style.LineStyle.Fill != "" {
				g.Fill = Lib.KMLReplaceString(value.Style.LineStyle.Fill)
			}
			if value.Style.PolyStyle.Color != "" {
				g.Fill = Lib.KMLReplaceString(value.Style.PolyStyle.Color)
			}
			if value.Style.PolyStyle.Fill != "" {
				g.Fill = Lib.KMLReplaceString(value.Style.PolyStyle.Fill)
			}
			if value.Style.PointStyle.Color != "" {
				g.Color = Lib.KMLReplaceString(value.Style.PointStyle.Color)
			}
			if value.Style.PointStyle.Fill != "" {
				g.Fill = Lib.KMLReplaceString(value.Style.PointStyle.Fill)
			}
			break
		}
	}
}

// 读取kml
func (a *App) ReadKmlFile(viewid, compid string) GeoMessage {
	ext := make([]runtime.FileFilter, 1)
	ext[0].DisplayName = "KML文件(*.kml)"
	ext[0].Pattern = "*.kml"
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title:   "选择文件",
		Filters: ext,
	})
	var geo GeoMessage
	if err != nil {
		geo.Error = err.Error()
		return geo
	}
	file, err := os.Open(path)
	if err != nil {
		geo.Error = err.Error()
		return geo
	}
	defer file.Close()
	// 解码 XML
	kmlbyte, _ := io.ReadAll(file)
	kmlbuf := bytes.NewBuffer(kmlbyte)
	var kml Lib.KML
	d := xml.NewDecoder(kmlbuf)
	err = d.Decode(&kml)
	if err != nil {
		geo.Error = err.Error()
		return geo
	}
	for _, value := range kml.Document.Folder.Placemarks {
		var g Geo
		g.Name = Lib.KMLReplaceString(value.Name)
		if len(value.Point.StringCoords) > 0 {
			g.CoordString = Lib.KMLReplaceString(value.Point.StringCoords)
			if value.StyleUrl != "" {
				a.SetGeoStyle(&g, value, kml.Document.Styles)
			}
			geo.Points = append(geo.Points, g)
		}
		if len(value.Polygon.OuterBoundaryIs.LinearRing.StringCoords) > 0 {
			g.CoordString = Lib.KMLReplaceString(value.Polygon.OuterBoundaryIs.LinearRing.StringCoords)
			if value.StyleUrl != "" {
				a.SetGeoStyle(&g, value, kml.Document.Styles)
			}
			geo.LinearRings = append(geo.LinearRings, g)
		}
		if len(value.LineString.StringCoords) > 0 {
			g.CoordString = Lib.KMLReplaceString(value.LineString.StringCoords)
			if value.StyleUrl != "" {
				a.SetGeoStyle(&g, value, kml.Document.Styles)
			}
			geo.LineStrings = append(geo.LineStrings, g)
		}
	}
	if len(geo.LineStrings) == 0 && len(geo.LinearRings) == 0 && len(geo.Points) == 0 {
		kmlbuf := bytes.NewBuffer(kmlbyte)
		var kml Lib.KML2
		d := xml.NewDecoder(kmlbuf)
		err = d.Decode(&kml)
		if err != nil {
			geo.Error = err.Error()
			return geo
		}
		for _, value := range kml.Document.Placemarks {
			var g Geo
			g.Name = Lib.KMLReplaceString(value.Name)
			if len(value.Point.StringCoords) > 0 {
				g.CoordString = Lib.KMLReplaceString(value.Point.StringCoords)
				if value.StyleUrl != "" {
					a.SetGeoStyle(&g, value, kml.Document.Styles)
				}
				geo.Points = append(geo.Points, g)
			}
			if len(value.Polygon.OuterBoundaryIs.LinearRing.StringCoords) > 0 {
				g.CoordString = Lib.KMLReplaceString(value.Polygon.OuterBoundaryIs.LinearRing.StringCoords)
				if value.StyleUrl != "" {
					a.SetGeoStyle(&g, value, kml.Document.Styles)
				}
				geo.LinearRings = append(geo.LinearRings, g)
			}
			if len(value.LineString.StringCoords) > 0 {
				g.CoordString = Lib.KMLReplaceString(value.LineString.StringCoords)
				if value.StyleUrl != "" {
					a.SetGeoStyle(&g, value, kml.Document.Styles)
				}
				geo.LineStrings = append(geo.LineStrings, g)
			}
		}
		if len(geo.LineStrings) == 0 && len(geo.LinearRings) == 0 && len(geo.Points) == 0 {
			geo.Error = "没有解析到数据，请确保KML内容结构符合标准。"
			return geo
		}
	}
	return geo
}

//程序Ui写文件到数据库

func (a *App) AppUIWriteFile(viewid string, comid string, data string) string {
	db := app.memsql.GetDB()
	// 将文件内容作为Blob类型保存到SQLite3数据库中
	stmt, err := db.Prepare("INSERT OR REPLACE INTO MEDIA values(?, ?, ?, ?)")
	if err != nil {
		a.memsql.Mutex.Unlock()
		return "Error:" + err.Error()
	}
	// defer stmt.Close()
	str := strings.Trim(data, "[]")
	arrStr := strings.Split(str, ",")
	var arr []byte
	for _, s := range arrStr {
		num, _ := strconv.Atoi(s)
		arr = append(arr, byte(num))
	}
	_, err = stmt.Exec(viewid, comid, comid+".png", arr)
	if err != nil {
		if e := stmt.Close(); e != nil {
			a.memsql.Mutex.Unlock()
			return "Error:" + e.Error()
		}
		a.memsql.Mutex.Unlock()
		return "Error:" + err.Error()
	}
	err = stmt.Close()
	if err != nil {
		a.memsql.Mutex.Unlock()
		return "Error:" + err.Error()
	}
	a.memsql.Mutex.Unlock()
	//base := "wails://wails/"
	// mimeType, err := Lib.GetMIMEType(comid + ".png")
	// if err != nil {
	// 	mimeType = ""
	// }
	return viewid + "NNI" + comid + ".png" //+ mimeType
}

// func (a *App) beforeClose(ctx context.Context) (prevent bool) {
// 	dialog, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
// 		Type:    runtime.QuestionDialog,
// 		Title:   "Quit?",
// 		Message: "Are you sure you want to quit?",
// 	})

// 	if err != nil {
// 		return false
// 	}
// 	return dialog != "Yes"
// }

func (b *App) beforeClose(ctx context.Context) (prevent bool) {
	dialog, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
		Type:    runtime.QuestionDialog,
		Title:   "退出程序?",
		Message: "视图数据可能没有保存，确定要退出程序吗?",
	})

	if err != nil {
		return false
	}
	return dialog != "Yes"
}

//检查程序更新

func (a *App) CheckUpdate() string {
	resp, err := http.Get("https://zerosir.com/viboard_update")
	if err != nil {
		return "Error:" + err.Error()
	}
	defer resp.Body.Close()

	// 读取API响应的内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Error:" + err.Error()
	}
	return string(body)
}

//读取License文件

func (a App) ReadLicense() string {
	f, err := tmpl.Open("templates/static/license.txt")
	if err != nil {
		return err.Error()
	}
	defer f.Close()
	var data []byte
	buffer := make([]byte, 1024)
	// _, err = f.Read(buffer)
	// if err != nil {
	// 	return err.Error()
	// }
	for {
		n, err := f.Read(buffer)
		if n == 0 {
			break // End of file reached
		}
		if err != nil {
			return err.Error()
		}
		for i := 0; i < n; i++ {
			data = append(data, buffer[i])
		}
	}
	return string(data)
}
