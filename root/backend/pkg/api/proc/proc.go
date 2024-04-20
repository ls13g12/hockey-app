package proc

import "github.com/ls13g12/hockey-app/root/backend/pkg/api/routers"

func StartApi() {
	routers.InitApiServer()
}
