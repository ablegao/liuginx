package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cgi"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

//用于解析json的数据类型
type ConfFileType struct {
	File   map[string]ConfCgi
	Server map[string]interface{} //服务器配置
}

type ConfCgi struct {
	Script string
	Argv   []string
	Env    []string
}

var ConfList ConfFileType

//基本配置
var cgipath = flag.String("path", ".", "Setting website room path. Default '.' ")
var port = flag.String("port", ":8080", "Set server port!")
var route_index = flag.String("route_index", "", "Route url to index! set index.php or index.py !")
var route_all = flag.Bool("route_all", true, "All file to route_index")

func LoadConf() {
	//获取脚本路径
	log.Println(os.Args[0])
	file, _ := exec.LookPath(os.Args[0])
	dir, _ := path.Split(file)
	//(path)
	conf := dir + "liuginx.conf"

	//获取配置文件信息
	if str, err := ioutil.ReadFile(conf); err == nil { //os.OpenFile(conf, os.O_RDWR, 0777)
		//str = []byte(strings.Replace(string(str), "$file", cgifile, -1))
		json.Unmarshal(str, &ConfList)
		//log.Println(ConfList)
	} else {
		log.Println(err)
	}

}

//起否启动url定向
var route_index_run bool = false

//路由器
type MyMux struct {
	cgi_room_path   string
	route_index_run bool
}

/*
func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.CgiRun(w, r)
	http.NotFound(w, r)
	return
}
*/
//cgi 脚本处理
func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := new(cgi.Handler)

	cgi_script_name := p.cgi_room_path + r.URL.Path

	fileStat, CgiFileError := os.Stat(cgi_script_name)

	//index 定向规则验证。 
	if true == p.route_index_run && nil != CgiFileError { //存在index定向 ， 同事cgi 脚本不存在时。 使用route_index指定脚本运行。 
		cgi_script_name = p.cgi_room_path + "/" + *route_index
	} else if false == p.route_index_run && nil != CgiFileError { //不存在index定向 ， 同时cgi脚本不存在时 ， 抛出错误提示。 
		http.NotFound(w, r)
		log.Printf("file not exists!%s.", cgi_script_name)
		return
	} else if true == p.route_index_run && r.URL.Path == "/" { //启动index时， 跟目录定向到index
		cgi_script_name = p.cgi_room_path + "/" + *route_index
	}
	//全部定向到 route_index
	if true == p.route_index_run && true == *route_all {
		cgi_script_name = p.cgi_room_path + "/" + *route_index
	}

	env := []string{}
	env = append(env, "SCRIPT_FILENAME="+cgi_script_name)

	/*
		if _, err := os.Stat(cgi_script_name); err != nil {
			//log.Println(p.route_index_run)
			http.NotFound(w, r)
			log.Printf("file not exists!%s.", cgi_script_name)
			return
		} else {
			log.Println(cgi_script_name)
		}*/

	//被解析路径
	handler.Dir = p.cgi_room_path
	ext := filepath.Ext(cgi_script_name)
	if info, ok := ConfList.File[ext]; ok {
		handler.Path = info.Script
		for x := 0; x < len(info.Argv); x++ {
			handler.Args = append(handler.Args, strings.Replace(info.Argv[x], "$file", cgi_script_name, -1))
		}
		for x := 0; x < len(info.Env); x++ {
			handler.Env = append(handler.Env, strings.Replace(info.Env[x], "$file", cgi_script_name, -1))
		}

		//替换SCRIPT_FILENAME
		handler.Env = append(handler.Env, env...)
		/*
			env := []string{
				"PATH_INFO=",
			}*/
		//w.Header().Add("Content-Type", "text/plain;charset=utf-8")

		handler.ServeHTTP(w, r)

	} else {
		if fileStat.IsDir() {
			w.Write([]byte("<div>Liugnix 是一个开发环境下的多平台CGI 服务器. 文档地址：<a href='https://github.com/ablegao/liugnix'/>github.com/ablegao/liugnix</a></div>"))
		}
		http.ServeFile(w, r, cgi_script_name)
		//http.Handle(r.URL.Path, http.StripPrefix("/", http.FileServer(http.Dir(cgi_room_path))))

	}

}

func main() {

	//载入配置文件
	LoadConf()
	runtime.GOMAXPROCS(int(ConfList.Server["process"].(float64)))
	flag.Parse()

	//判断是否设定index
	if len(*route_index) > 0 {
		route_index_run = true
	} else {
		route_index_run = false
	}

	if cgi_room_path, err := filepath.Abs(*cgipath); err == nil {
		mux := &MyMux{
			cgi_room_path:   cgi_room_path,
			route_index_run: route_index_run,
		}
		log.Fatal(http.ListenAndServe(*port, mux))
	} else {
		log.Printf("error : %v", err)
	}
}
