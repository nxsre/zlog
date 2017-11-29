package file

type FileConfig struct {
	Enable        bool   `yaml:"enable"`
	Path          string `yaml:"path"`
	Filename      string `yaml:"filename"`
	RotateEveryKB int    `yaml:"rotate_every_kb"`
	NumberOfFiles int    `yaml:"number_of_files"`
}
