package static

import (
	"net/http"
	"gobible/logmanager/cli/http/controllers"
)

func SearchdirRouter(ZipResultDir string)  {

	controllers.GetGlobalRouter().ServeFiles("/log/*filepath",http.Dir(ZipResultDir))
}
