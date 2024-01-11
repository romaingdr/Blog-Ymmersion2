package main

import (
	"fmt"
	"ymmersion2/backend"
	"ymmersion2/routeur"
	"ymmersion2/templates"
)

func main() {
	active, user := backend.CheckRememberStatus("rememberSession.json")
	if active {
		fmt.Println("Une session a été sauvegardée")
		backend.GlobalSession = backend.Session{Username: user, State: backend.GetAccountState(user), Mail: backend.GetAccountMail(user)}
		fmt.Println("Session initialisée")
	}
	templates.InitTemplate()
	routeur.Initserv()
}
