package maps

import (
	"reflect"
	"strings"
)

// MapToStruct 是一个函数，它将一个 map 的值复制到一个结构体实例中。
// 它使用 Go 的反射机制来实现这个功能。
// data 是我们要复制的 map，obj 是我们要复制到的结构体实例的指针。
func MapToStruct(data map[string]interface{}, obj interface{}) {
	// 获取 obj 的类型和值的反射对象
	t := reflect.TypeOf(obj).Elem()
	v := reflect.ValueOf(obj).Elem()
	// 遍历结构体的所有字段
	for i := 0; i < t.NumField(); i++ {
		// 获取当前字段的反射对象
		field := t.Field(i)
		// 获取字段的 json 标签
		tag := field.Tag.Get("json")
		// 如果标签为空或为"-"，则跳过此字段
		if tag == "" || tag == "-" {
			continue
		}
		// 使用 json 标签作为键
		// 如果标签有额外的选项，如 `json:"field,optional"`，则只使用逗号前的部分作为键
		tag = strings.Split(tag, ",")[0]
		// 从 map 中获取值
		if value, ok := data[tag]; ok {
			// 获取结构体字段的反射对象
			val := v.Field(i)
			// 如果字段是一个指向字符串的指针
			if val.Kind() == reflect.Ptr && val.Type().Elem().Kind() == reflect.String {
				// 将值转换为字符串
				str := value.(string)
				// 将字符串的地址设置到字段中
				val.Set(reflect.ValueOf(&str))
			} else {
				// 否则，直接将值设置到字段中
				val.Set(reflect.ValueOf(value))
			}
			// 从 map 中删除已经设置的项
			delete(data, tag)
		}
	}
}
