#liugnix 
一个 cgi模式 web服务器。 

##支持语言
因使用的为CGI 模式，支持所有cgi模式运行的程序。 

php

python

go

perl







#如何运行

cd mywebsite

liugnix_ENV 


	liuginx --help
	2013/05/07 01:10:43 liuginx
	Usage of liuginx:
	  -path=".": Setting website room path. Default '.'
	  -port=":8080": Set server port!
	  -route_all=true: All file to route_index
	  -route_index="": Route url to index! set index.php or index.py !


#安装方法？

	已经有编译好的程序， 可以根据自己的需要下载

	windows 32位系统 :

	https://github.com/ablegao/liugnix/liugnix_win32.exe

	windows 64位系统:
	https://github.com/ablegao/liugnix/liugnix_win64.exe

	linux 64位:
	https://github.com/ablegao/liugnix/liuginx_linux64

	Mac 10.8.3
	https://github.com/ablegao/liugnix/liuginx_mac


#自己编译

1. 确保自己的GOLang环境可以运行。 

2. go get github.com/ablegao/liuginx

3. go install github.com/ablegao/liuginx

#如何配置语言解析？

你需要一个名为liuginx.conf文件 ， 这个文件需要和liuginx在同一个目录下。 （源码中有附带。 ）

配置文件内容：
	{
		"server":{
			"process":10240
		},
		"file":{
			".php":{
				"script":"/usr/local/bin/php-cgi",
				"argv":["$file"],
				"env":[]
			},

			".phps":{
				"script":"/usr/local/bin/php-cgi",
				"argv":["-s","$file"],
				"env":[]
			},


			".go":{
				"script":"/usr/local/go11/bin/go",
				"argv":["run" , "$file"],
				"env":[
					"GOPATH=~/code/mygo",
					"GOROOT=/usr/local/go11"
				]
			},
			".py":{
				"script":"/opt/python/py27/bin/python",
				"argv":["-u", "$file"],
				"env":[]
			}
		}
		
	}

##配置文件结构


server:{
	process: 最大进程数， 如果配置的很高， 能带来很高的并发。 
}

file下的配置：
  如.php ， 则是配置php文件的解析。 
  其中， script 用来指定php解析器。 这里使用的 php-cgi ,windows 下是php-cgi.exe
  argv 中$file 为一个必填选项。 实际执行时相当于  /usr/local/bin/php-cgi $file

  看.phps的解析 ， 就是再 php-cgi 基础上加入了一个 -s 参数， 同时， 可以查看php-cgi的其他命令参数， 加入到argv中。 

  env 用于扩展其他命令集参数。 

##如果我的使用的是thinkphp cakephp等， 怎么隐藏index.php文件， 实现路由功能？

liuginx.exe --route_index=index.php 

##有的文件， 已经存在， 我并不想经过路由怎么办?
liuginx.exe --route_index=index.php  --route_all=false 


