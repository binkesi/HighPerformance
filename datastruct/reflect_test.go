// Elem and Indirect:
/*
If a reflect.Value is a pointer, then v.Elem() is equivalent to reflect.Indirect(v). If it is not a pointer, then they are not equivalent:

If the value is an interface then reflect.Indirect(v) will return the same value, while v.Elem() will return the contained dynamic value.
If the value is something else, then v.Elem() will panic.
The reflect.Indirect helper is intended for cases where you want to accept either a particular type, or a pointer to that type.
One example is the database/sql conversion routines: by using reflect.Indirect, it can use the same code paths to handle the various types and pointers to those types.
*/

// 使用反射赋值，效率非常低下，如果有替代方案，尽可能避免使用反射，特别是会被反复调用的热点代码
// 例如 RPC 协议中，需要对结构体进行序列化和反序列化，这个时候避免使用 Go 语言自带的 json 的 Marshal 和 Unmarshal 方法，
// 因为标准库中的 json 序列化和反序列化是利用反射实现的。可选的替代方案有 easyjson，在大部分场景下，相比标准库，有 5 倍左右的性能提升。

package datastruct

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

type Config struct {
	Name    string `json:"server-name"` // CONFIG_SERVER_NAME
	IP      string `json:"server-ip"`   // CONFIG_SERVER_IP
	URL     string `json:"server-url"`  // CONFIG_SERVER_URL
	Timeout string `json:"timeout"`     // CONFIG_TIMEOUT
}

func readconfig() *Config {
	config := Config{}
	typ := reflect.TypeOf(config)
	val := reflect.Indirect(reflect.ValueOf(&config))
	for i := 0; i < typ.NumField(); i++ {
		ft := typ.Field(i)
		if v, ok := ft.Tag.Lookup("json"); ok {
			key := fmt.Sprintf("CONFIG_%s", strings.ReplaceAll(strings.ToUpper(v), "-", "_"))
			if env, exist := os.LookupEnv(key); exist {
				val.FieldByName(ft.Name).Set(reflect.ValueOf(env))
			}
		}
	}
	return &config
}

func TestConfig(t *testing.T) {
	os.Setenv("CONFIG_SERVER_NAME", "global_server")
	os.Setenv("CONFIG_SERVER_IP", "10.0.0.1")
	os.Setenv("CONFIG_SERVER_URL", "sungn.com")
	rc := readconfig()
	fmt.Printf("%+v", rc)
}

func BenchmarkNew(b *testing.B) {
	var config *Config
	for i := 0; i < b.N; i++ {
		config = new(Config)
	}
	_ = config
}

func BenchmarkReflectNew(b *testing.B) {
	var config *Config
	typ := reflect.TypeOf(Config{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config, _ = reflect.New(typ).Interface().(*Config)
	}
	_ = config
}

func BenchmarkSet(b *testing.B) {
	config := new(Config)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config.Name = "name"
		config.IP = "ip"
		config.URL = "url"
		config.Timeout = "timeout"
	}
}

func BenchmarkReflect_FieldSet(b *testing.B) {
	typ := reflect.TypeOf(Config{})
	ins := reflect.New(typ).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ins.Field(0).SetString("name")
		ins.Field(1).SetString("ip")
		ins.Field(2).SetString("url")
		ins.Field(3).SetString("timeout")
	}
}

func BenchmarkReflect_FieldByNameSet(b *testing.B) {
	typ := reflect.TypeOf(Config{})
	ins := reflect.New(typ).Elem()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ins.FieldByName("Name").SetString("name")
		ins.FieldByName("IP").SetString("ip")
		ins.FieldByName("URL").SetString("url")
		ins.FieldByName("Timeout").SetString("timeout")
	}
}

func BenchmarkReflect_FieldByNameCacheSet(b *testing.B) {
	typ := reflect.TypeOf(Config{})
	ins := reflect.New(typ).Elem()
	cache := make(map[string]int)
	for i := 0; i < typ.NumField(); i++ {
		cache[typ.Field(i).Name] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ins.Field(cache["Name"]).SetString("name")
		ins.Field(cache["IP"]).SetString("ip")
		ins.Field(cache["URL"]).SetString("url")
		ins.Field(cache["Timeout"]).SetString("timeout")
	}
}
