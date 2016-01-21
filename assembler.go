package main

import (
  "os"
  "fmt"
  "io/ioutil"

  "gopkg.in/yaml.v2"
  "github.com/gliderlabs/sigil"
  _ "github.com/gliderlabs/sigil/builtin"
)

var config = make(map[string]string)
var services = make(map[string]interface{})

func main() {

  data, err := ioutil.ReadFile(os.Args[1])
  fatal(err)
  var obj map[string]interface{}
	fatal(yaml.Unmarshal(data, &obj))
  if value, ok := obj["export"]; ok {
    exports := value.(map[interface{}]interface{})
    for k,v := range exports {
      config[k.(string)] = v.(string)
    }
  }
  if value, ok := obj["services"]; ok {
    svcs := value.([]interface{})
    for _, v := range svcs {
      load(v.(string))
    }
  }
  yml, err := yaml.Marshal(services)
  fatal(err)
  fmt.Println(string(yml))
}

func load(name string) {
  data, err := ioutil.ReadFile(name + ".yml")
  fatal(err)
  rendered, err := sigil.Execute(data, config, name)
  fatal(err)
  var service map[string]interface{}
	fatal(yaml.Unmarshal(rendered.Bytes(), &service))
  if value, ok := service["export"]; ok {
    exports := value.(map[interface{}]interface{})
    for k,v := range exports {
      config[k.(string)] = v.(string)
    }
    delete(service, "export")
  }
  services[name] = service
}

func fatal(err error) {
  if (err != nil) {
    panic(err)
  }
}
