package models

type Config struct {
	CommafeedSetup []Category `yaml:"commafeedSetup"`
}

type Category struct {
	ParentName *string     `yaml:"parentName,omitempty"`
	Name       string      `yaml:"category"`
	Feeds      []Feed      `yaml:"feeds"`
	Children   []*Category `yaml:"children"`
}

type Feed struct {
	Name     string `yaml:"name"`
	Url      string `yaml:"url"`
	Disabled bool   `yaml:"disabled,omitempty"`
}
