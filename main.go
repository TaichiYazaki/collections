package main

import (
	"collections/awesome"
	"collections/getreadmeurl"
)

func main() {
	r :=getreadmeurl.ReadmeElement{TopPage: awesome.TopPage, Target: awesome.Target}
	getreadmeurl.SetReadmeURL(r)
	arg:= awesome.InputToSection()
	awesome.CountSection(arg)
}