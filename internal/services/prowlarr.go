package services

var Prowlarr = &prowlarr{}

type prowlarr struct{}

func (*prowlarr) Check(url string) bool {
	return true
}
