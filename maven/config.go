package maven

type ServerConf struct {
	Port      int    `yaml:"port"`
	MavenPath string `yaml:"mavenpath"`
	BasePath  string `yaml:"basepath"`
}

type DSN struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
}

type Configuration struct {
	Database DSN        `yaml:"database"`
	Server   ServerConf `yaml:"server"`
}
