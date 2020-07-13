package config

type ConfParser interface {
	Marshal(conf interface{}, toFile string) error
	Unmarshal(conf interface{}, fromFile string) error
}

// NewConfParser 创建不同解析器的工厂函数
func NewConfParser(kind string) ConfParser {
	switch kind {
	case "yaml":
		return newYamlParser()
	case "json":
		return newJSONParser()
	}
	return nil
}
