package constants

import "net/http"

const (
	URI_V1                         = "/v1"
	URI_V1_REGISTER_FUNCTION_CHAIN = "/registerFunctionChain"
)

const (
	REQUEST_METHOD_POST = http.MethodPost
)

const (
	APPLICATION_CONFIG_NAME = "application"
	RESOURCES_PATH          = "./resources"
	CONFIG_FILE_TYPE        = "yaml"
	PORT_CONFIG_NAME        = "web.port"
	WEB_FRAMEWORK           = "web.framework"
	GORILLA                 = "gorilla"
	ECHO                    = "echo"
	CONTANT_TYPE_JSON       = "application/json"

	REQUEST_TIME    = "test.performance.request_time"
	REQUEST_THREADS = "test.performance.request_threads"
)

const (
	FUNCTION_EXCUTION_LOG_PREFIX = "**********Executing Function**********"
	FUNCTION_EXCUTION_LOG_SURFIX = "**********Done**********"
)
