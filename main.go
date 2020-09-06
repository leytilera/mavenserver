package main

import (
	"gopkg.in/yaml.v2"
	"mavenserver/maven"
	"os"
)

func main() {
	var auth maven.AuthManager
	var conf string
	if len(os.Args) > 1 {
		conf = os.Args[1]
	} else {
		conf = "config.yaml"
	}
	config := readconf(conf)
	auth = maven.Create(config.Database)
	maven.StartServer(config.Server, &auth)
}

func readconf(loc string) maven.Configuration {
	if _, err := os.Stat(loc); os.IsNotExist(err) {
		conf := `
server:
  port: 8080
  mavenpath: "maven"
  basepath: "/maven"
database:
  host: "localhost"
  database: "mavenusers"
  username: "root"
  password: "password"`
		file, err := os.Create(loc)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		_, err = file.WriteString(conf)
		if err != nil {
			panic(err)
		}
	}
	file, err := os.Open(loc)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var conf maven.Configuration
	dec := yaml.NewDecoder(file)
	err = dec.Decode(&conf)
	if err != nil {
		panic(err)
	}
	return conf
}
