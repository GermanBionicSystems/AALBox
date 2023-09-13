package main

import "AALBox/src/controller"

func main() {
	//TODO: initialize the docker container (postgres database). maybe do docker file. for now start it manually

	control := controller.NewController()
	control.Start()
}
