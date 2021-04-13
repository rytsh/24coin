package common

var Version string = "v0.0.0"
var Settings settings

type settings struct {
	UI struct {
		Port int
		Host string
	}
}
