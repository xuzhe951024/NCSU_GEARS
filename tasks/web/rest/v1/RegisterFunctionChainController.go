package v1

import (
	"NCSU_Gears/common/constants"
	"NCSU_Gears/models"
	"NCSU_Gears/tasks/web/service"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

var RegisteredFunctionChainsMap map[string]map[string]models.Function

func registerGorillaFunctionChainHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		logrus.Error("Failed to read body")
		return
	}

	var result models.RegisterFunctionChainVO
	err = json.Unmarshal(body, &result)
	if err != nil {
		http.Error(w, "Failed to parse body", http.StatusBadRequest)
		logrus.Error("Failed to parse body")
		return
	}

	fnMappings, funcs := service.ParseJsonToMaps(result)

	if nil == RegisteredFunctionChainsMap {
		initMap := make(map[string]map[string]models.Function)
		initMap[result.Identifier] = fnMappings
		RegisteredFunctionChainsMap = initMap
	} else {
		RegisteredFunctionChainsMap[result.Identifier] = fnMappings
	}

	if nil != fnMappings && nil != funcs {
		w.WriteHeader(http.StatusOK)
		response, _ := json.Marshal(result.Identifier)
		w.Write(response)
	}
}

func RegisterGorillaFunctionChainController(port string) {
	router := mux.NewRouter()
	router.HandleFunc(constants.URI_V1_REGISTER_FUNCTION_CHAIN, registerGorillaFunctionChainHandler).Methods(constants.REQUEST_METHOD_POST)
	logrus.Info("starting serving on port: " + port)
	//err := http.ListenAndServe(port, router)
	err := http.ListenAndServeTLS(port, "server.crt", "server.key", router)
	if err != nil {
		logrus.Error(err)
		return
	}
}

func registerEchoFunctionChainHandler(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)

	if err != nil {
		logrus.Error("Failed to read body")
		return c.String(http.StatusBadRequest, "Failed to read body")
	}

	var result models.RegisterFunctionChainVO
	err = json.Unmarshal(body, &result)
	if err != nil {
		logrus.Error("Failed to parse body")
		return c.String(http.StatusBadRequest, "Failed to parse body")
	}

	fnMappings, funcs := service.ParseJsonToMaps(result)

	if nil == RegisteredFunctionChainsMap {
		initMap := make(map[string]map[string]models.Function)
		initMap[result.Identifier] = fnMappings
		RegisteredFunctionChainsMap = initMap
	} else {
		RegisteredFunctionChainsMap[result.Identifier] = fnMappings
	}

	if nil != fnMappings && nil != funcs {
		response, _ := json.Marshal(result.Identifier)
		return c.String(http.StatusOK, string(response))
	}
	return nil
}

func RegisterEchoFunctionChainController(port string) {
	e := echo.New()

	e.POST(constants.URI_V1_REGISTER_FUNCTION_CHAIN, registerEchoFunctionChainHandler)

	logrus.Info("starting serving on port: " + port)
	//err := e.Start(port)
	err := e.StartTLS(port, "server.crt", "server.key")
	if err != nil {
		logrus.Error(err)
		return
	}
}
