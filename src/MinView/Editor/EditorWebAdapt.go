package Editor

import (
	"MinView/Lib"
	"net/http"
	"os"
	"strconv"

	// "github.com/gorilla/websocket"
	// "github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// 编辑器WEB服务，用于上传附件并生成链接给视图访问
type EditorWeb struct {
	//上传文件保存路径
	path string
	//端口，用于视图访问的端口，
	//可能打开多个实例的时候，端口打开失败，需要端口++
	//默认从61234开始，port++最大10次
	port int
	//打开服务的IP,默认127.0.0.1
	host string
	//需要验证token正确才能收发信息
	token string
}

// 保存token给uploadHandler访问
// var token string

func NewEditorWeb(port int) *EditorWeb {
	var web EditorWeb
	web.path = "/tmp/ViBoard/"
	web.port = port
	web.host = "127.0.0.1"
	web.token = uuid.NewV4().String()
	// token = web.token
	return &web
}

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin: func(r *http.Request) bool {
// 		// 在此处添加你允许的Origin逻辑
// 		return true
// 	},
// }

// func uploadHandler(w http.ResponseWriter, r *http.Request) {
// 	// 升级连接为 WebSocket 连接
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		fmt.Println("无法升级WebSocket连接", err)
// 		return
// 	}
// 	defer conn.Close()

// 	// 从客户端接收文件
// 	for {
// 		// 读取消息
// 		_, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			fmt.Println(err)
// 			break
// 		}

// 		// 解析 JSON 数据
// 		var data map[string]interface{}
// 		err = json.Unmarshal(msg, &data)
// 		if err != nil {
// 			fmt.Println(err)
// 			break
// 		}
// 		fmt.Println(data)
// 		// 处理接收到的 JSON 数据
// 		if data["token"] != token {
// 			continue
// 		}

// 	}
// }

// func handleWebSocket(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	defer conn.Close()

// 	for {
// 		// 读取 WebSocket 消息
// 		_, message, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("read:", err)
// 			break
// 		}

// 		// 打印收到的消息
// 		fmt.Printf("received: %s\n", message)

// 		// 回复消息
// 		err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Hello %s, It's show time!", message)))
// 		if err != nil {
// 			log.Println("write:", err)
// 			break
// 		}
// 	}
// }

func (w *EditorWeb) Init() error {
	result, err := Lib.PathExists(w.path)
	if err != nil {
		return err
	}
	if !result {
		if err := os.MkdirAll(w.path, 0755); err != nil {
			return err
		}
	}
	// // 处理 WebSocket 连接
	// http.HandleFunc("/CeilAI", uploadHandler)
	// // 提供静态文件服务
	// fs := http.FileServer(http.Dir(w.path))
	// http.Handle("/CeilAI/", http.StripPrefix("/CeilAI", fs))
	// http.HandleFunc("/ws", handleWebSocket)
	// http.ListenAndServe(":8080", nil)
	//尝试十次启动，可能端口被占用
	var loop_result error
	start_count := 10
	for start_count >= 0 {
		port := strconv.Itoa(w.port)
		err := http.ListenAndServe(w.host+":"+port, nil)
		if err == nil {
			loop_result = nil
			break
		} else {
			loop_result = err
		}
		w.port++
		start_count--
	}
	return loop_result
}

func (w *EditorWeb) GetPort() int {
	return w.port
}

func (w *EditorWeb) GetToken() string {
	return w.token
}
func (w *EditorWeb) GetPath() string {
	return w.path
}
