package ginger

import (
	"flag"
	"fmt"
	"os"
	"sync"
)

var CliIns *Cli
var cliOnce sync.Once

//单例模式获取实例
func GetCli() *Cli {
	cliOnce.Do(func() {
		CliIns = new(Cli)
	})
	return CliIns
}

const (
	CliString = iota
	CliBool
	CliInt
)

const (
	ExplainMain = iota
	ExplainItem
)

type CliRow struct {
	Name    string
	Type    int
	Default interface{}
	UseAge  string
}

type Explain struct {
	Info string
	Type int
}

type Cli struct {
	Row        CliRow
	ExplainRow Explain
	Explains   []Explain
}

type CliInput interface {
	GetInt()int
	GetString()string
	GetBool()bool
}

func (c *Cli) Item(name string) *Cli {
	c.Row.Name = name
	return c
}

func (c *Cli) SetDefault(d interface{}) *Cli {
	c.Row.Default = d
	c.Row.Type = CliString
	return c
}

func (c *Cli) SetIntDefault(d interface{}) *Cli {
	c.Row.Default = d
	c.Row.Type = CliInt
	return c
}

func (c *Cli) SetBoolDefault(d interface{}) *Cli {
	c.Row.Default = d
	c.Row.Type = CliBool
	return c
}

func (c *Cli) SetUsage(usage string) *Cli {
	c.Row.UseAge = usage
	return c
}

func (c *Cli) Save() CliInput {
	var result interface{}
	if c.Row.Type == CliString{
		result = flag.String(c.Row.Name, c.Row.Default.(string), c.Row.UseAge)
	}
	if c.Row.Type == CliBool{
		result = flag.Bool(c.Row.Name, c.Row.Default.(bool), c.Row.UseAge)
	}
	if c.Row.Type == CliInt{
		result = flag.Int(c.Row.Name, c.Row.Default.(int), c.Row.UseAge)
	}
	c.Row = CliRow{}
	cliData := new(CliData)
	cliData.data = result
	return cliData
}

type CliData struct {
	data interface{}
}

func (c *CliData)GetInt() int {
	return *c.data.(*int)
}

func (c *CliData)GetString() string {
	return *c.data.(*string)
}

func (c *CliData)GetBool() bool {
	return *c.data.(*bool)
}

func (c *Cli) Explain(info string) *Cli {
	c.ExplainRow.Type = ExplainMain
	c.ExplainRow.Info = info
	c.Explains = append(c.Explains, c.ExplainRow)
	c.ExplainRow = Explain{}
	return c
}

func (c *Cli) SetExplainItem(info string) *Cli {
	c.ExplainRow.Type = ExplainItem
	c.ExplainRow.Info = info
	c.Explains = append(c.Explains, c.ExplainRow)
	c.ExplainRow = Explain{}
	return c
}

//解析CLI命令并返回是否有说明
func (c *Cli) SaveExplain() bool {
	var hasExplain = false
	flag.Usage = func() {
		flag.PrintDefaults()
		var eNum int
		if len(c.Explains) > 0 {
			hasExplain = true
		}
		for _, v := range c.Explains {
			if v.Type == ExplainMain {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("%c[1;32;40m\n%s：%c[0m\n", 0x1B, v.Info, 0x1B))
				eNum = 1
			} else {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("%c[1;33;40m (%v) %s%c[0m\n", 0x1B, eNum, v.Info, 0x1B))
				eNum++
			}
		}
	}
	flag.Parse()
	return hasExplain
}
