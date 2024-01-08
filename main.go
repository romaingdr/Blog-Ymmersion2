package main

import (
	"ymmersion2/routeur"
	"ymmersion2/templates"
)

func main() {
	templates.InitTemplate()
	routeur.Initserv()
}
