

公共结构

```
type MutiConf struct {
	Driver Db `yaml:"db"`
}
type Db struct {
	DBinfo Mysql `yaml:"mysql"`
}
type Mysql struct {
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Server   string `yaml:"server"`
	Port     int    `yaml:"port"`
	Info     Info   `yaml:"info"`
}
type Info struct {
	ID      int    `yaml:"id"`
	Postion string `yaml:"postion"`
}
```



> yaml配置文件解析demo

```
// 序列化文件数据
mutiConf := MutiConf{
		Driver: Db{
			DBinfo: Mysql{
				Name: "mysql",
				User: "admin",
				Port: 999,
				Info: Info{
					ID:      1111111,
					Postion: "china",
				},
			},
		},
	}
	// marshal data to a file
	//parser := config.NewYamlParser("test01.yaml")
	parser := config.NewConfParser("yaml")
	parser.Marshal(mutiConf, "test02.yaml")
	// 函数方式
	config.YamlMarshalAndWriteFile(mutiConf, "testFun02.yaml")
	
	// 反序列化数据到文件
		var conf MutiConf
	//parser := config.NewYamlParser("test01.yaml")
	parser := config.NewConfParser("yaml")
	if err := parser.Unmarshal(&conf, "test02.yaml"); err != nil {
		log.Println("unmarshal error:", err)
	}
	log.Println("\n", conf)
	// 函数方式
	var conf2 MutiConf
	config.YamlReadFileAndUnmarshal(&conf2, "testFun02.yaml")
	log.Println("\n", conf2)

```



说明：

有两种方式对序列化yaml文件数据和反序列化yaml文件数据

第一种是解析器对象，第二种是调用函数