package ginger

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //导入MYSQL
	"log"
	"reflect"
	"strings"
	"time"
)

//Orm mysql的orm包拓展
type Orm struct {
	user     string
	password string
	server   string
	port     int
	dbname   string
	*gorm.DB
}

//SetDbname 设置数据库名称
func (o *Orm) SetDbname(dbname string) *Orm {
	o.dbname = dbname
	return o
}

//SetPort 设置端口
func (o *Orm) SetPort(port int) *Orm {
	o.port = port
	return o
}

//SetServer 设置服务器
func (o *Orm) SetServer(server string) *Orm {
	o.server = server
	return o
}

//SetPassword 设置密码
func (o *Orm) SetPassword(password string) *Orm {
	o.password = password
	return o
}

//SetUser 设置用户名
func (o *Orm) SetUser(user string) *Orm {
	o.user = user
	return o
}

//NewOrm 返回Orm实例
func NewOrm() *Orm {
	return &Orm{}
}

//Register 注册orm
func (o *Orm) Register() *Orm {
	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local",
		o.user, o.password, o.server, o.port, o.dbname,
	)
	var err error
	o.DB, err = gorm.Open("mysql", conn)
	if err != nil {
		log.Fatalln("连接数据库时出错：", err.Error())
	}
	/*打印日志*/
	o.DB.LogMode(true)
	return o
}

//AddAll 添加多行数据
func (o *Orm) SaveAll(data interface{}) error {
	isPtr := reflect.TypeOf(data).Kind().String() == "ptr"
	if !isPtr {
		return errors.New("必须是一个地址")
	}
	isSlice := reflect.TypeOf(data).Elem().Kind().String() == "slice"
	if !isSlice {
		return errors.New("不是一个slice")
	}
	//获取interface的元素,因为上面已经判定它必定是一个地址了
	element := reflect.ValueOf(data).Elem()
	eleLen := element.Len()
	if eleLen<1{
		return errors.New("没有需要添加的数据")
	}
	//取出第一行数据，用来获取表名
	tableName := o.snakeString(reflect.TypeOf(element.Index(0).Interface()).Name())
	//定义字段名 数组
	var fieldsArr []string
	//对应的值
	var allRows []string
	//如果插入不了，则更新的字段数据
	var updateFieldsArr []string
	//每一行的数据
	var rowValue []string
	var fName string
	dataType := reflect.TypeOf(element.Index(0).Interface())
	var fnTypes []string
	var gormTags []string
	for i := 0; i < element.Index(0).NumField(); i++ {
		fName = dataType.Field(i).Name
		defaultTag := dataType.Field(i).Tag.Get("json")
		gormTag := dataType.Field(i).Tag.Get("gorm")
		if gormTag=="-"{
			continue
		}
		gormTags = append(gormTags,gormTag)
		fName = dataType.Field(i).Name
		if defaultTag != ""{
			fName = defaultTag
		}else{
			fName = o.snakeString(fName)
		}
		fnType := dataType.Field(i).Type.String()
		fnTypes = append(fnTypes,fnType)
		fieldsArr = append(fieldsArr, fmt.Sprintf("`%s`", fName))
		updateFieldsArr = append(updateFieldsArr, fmt.Sprintf("%s = VALUES(%s)", fName, fName))
	}
	for i := 0; i < eleLen; i++ {
		rowKeyLen := element.Index(i).NumField()
		rowValue = []string{}
		for rowKey := 0; rowKey < rowKeyLen; rowKey++ {
			gormTag := gormTags[rowKey]
			if gormTag=="-"{
				continue
			}
			fnType := fnTypes[rowKey]
			var rowV interface{}
			if fnType=="time.Time"{
				rowV = time.Now().Format("2006-01-02 15:04:05")
			} else if fnType=="*time.Time"{
				rowV = "NULL"
				rowValue = append(rowValue, fmt.Sprintf("%v", rowV))
				continue
			}else{
				rowV = element.Index(i).Field(rowKey).Interface()
			}
			rowValue = append(rowValue, fmt.Sprintf("'%v'", rowV))
		}
		allRows = append(allRows, fmt.Sprintf("(%s)", strings.Join(rowValue, ",")))
	}
	sql := fmt.Sprintf("INSERT INTO `%ss`(%s) VALUES %s ON DUPLICATE KEY UPDATE %s",
		tableName,
		strings.Join(fieldsArr, ","),
		strings.Join(allRows, ","),
		strings.Join(updateFieldsArr, ","),
	)
	return o.DB.Exec(sql).Error
}

//snakeString 驼峰命名转蛇形
func (o *Orm) snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}
