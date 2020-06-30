# ginger
Package of gin framework

> 开源库使用申明：

1.gin

    github.com/gin-gonic/gin

2.gzip压缩

    https://github.com/gin-contrib/gzip

3.swagger

    github.com/swaggo/gin-swagger

> Wish list

1.gzip yes

2.上传，图像裁剪，真实URL屏蔽 no

3.gin常用功能封装 yes

4.gin命令行工具 yes

5.数据库工具(数据库迁移，数据库填充，数据库默认创建，自动字段生成) partial yes https://github.com/gobkc/md

6.日志处理 ES，FILE支持 yes 全局注入 https://github.com/gobkc/logger

7.swagger集成 yes

8.socket和ws协议集成 no

9.grpc支持 yes

10.test单元测试，基准测试快速创建工具 no

11.curl库支持

>使用方式

    export GOPROXY=https://mirrors.aliyun.com/goproxy
    go mod edit -require github.com/gobkc/ginger@v0.0.8
    go mod tidy
    go mod vendor

## DEMO LIST

> curl

1.GET请求 不携带参数，并且将结果反馈给result
    
    var result string
	var err = ginger.NewCurl("http://gobk.cn:9999/").
		Token("57aaec9cbd72").
		Get().BindString(&result)

2.GET请求 携带参数(param为结构体)

    var result string
    var param = Param{}
    param.Username = "xiong"
	var err = ginger.NewCurl("http://gobk.cn:9999/").
		Token("57aaec9cbd72").
		Get(&param).BindString(&result)

3.POST请求 携带参数，并将结果返回到一个结构体

    var result = UserOk{}
    var param = Param{}
    param.Username = "xiong"
    param.Password = "123456"
	var err = ginger.NewCurl("http://gobk.cn:9999/").
		Token("57aaec9cbd72").
		POST(&param).Bind(&result)

4.其余请求参考上述用法，可用BindString获取返回的字符串自行处理，也可以用Bind获取结构体，格式化为你需要的数据

5.Token方法和Auth方法同时存在只有一个起作用

6.可以用ContentType方法设置请求类型,也可以多次调用Header方法来设置多个header
