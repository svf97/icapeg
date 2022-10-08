package blockexe

import (
	http_message "icapeg/http-message"
	"icapeg/logging"
	"icapeg/readValues"
	services_utilities "icapeg/service/services-utilities"
	general_functions "icapeg/service/services-utilities/general-functions"
	"sync"
	"time"
)

var doOnce sync.Once
var blockexeConfig *blockexe

const BEIdentifier = "BlockExe ID"

type blockexe struct {
	httpMsg                    *http_message.HttpMsg
	elapsed                    time.Duration
	serviceName                string
	methodName                 string
	maxFileSize                int
	bypassExts                 []string
	processExts                []string
	rejectExts                 []string
	extArrs                    []services_utilities.Extension
	SocketPath                 string
	Timeout                    time.Duration
	badFileStatus              []string
	okFileStatus               []string
	returnOrigIfMaxSizeExc     bool
	return400IfFileExtRejected bool
	generalFunc                *general_functions.GeneralFunc

	//optional, it's up to you and to optional variables have been added in the service section in config.toml file (you should map them with these struct fields)
	// generalFunc            *general_functions.GeneralFunc     //optional helper field
	// base_url = "echo"
	// scan_endpoint = "echo"
	// api_key = "<api key>"
	// timeout  = 300 #seconds , ICAP will return 408 - Request timeout
	// fail_threshold = 2
	// max_filesize = 0 #bytes
	// return_original_if_max_file_size_exceeded=false
}

func InitBlockexeConfig(serviceName string) {
	logging.Logger.Debug("loading " + serviceName + " service configurations")
	doOnce.Do(func() {
		blockexeConfig = &blockexe{
			maxFileSize:                readValues.ReadValuesInt(serviceName + ".max_filesize"),
			bypassExts:                 readValues.ReadValuesSlice(serviceName + ".bypass_extensions"),
			processExts:                readValues.ReadValuesSlice(serviceName + ".process_extensions"),
			rejectExts:                 readValues.ReadValuesSlice(serviceName + ".reject_extensions"),
			Timeout:                    readValues.ReadValuesDuration(serviceName+".timeout") * time.Second,
			returnOrigIfMaxSizeExc:     readValues.ReadValuesBool(serviceName + ".return_original_if_max_file_size_exceeded"),
			return400IfFileExtRejected: readValues.ReadValuesBool(serviceName + ".return_400_if_file_ext_rejected"),
		}
		blockexeConfig.extArrs = services_utilities.InitExtsArr(blockexeConfig.processExts, blockexeConfig.rejectExts, blockexeConfig.bypassExts)
	})
}

func NewBlockexeService(serviceName, methodName string, httpMsg *http_message.HttpMsg) *blockexe {
	return &blockexe{
		httpMsg:                    httpMsg,
		serviceName:                serviceName,
		methodName:                 methodName,
		generalFunc:                general_functions.NewGeneralFunc(httpMsg),
		maxFileSize:                blockexeConfig.maxFileSize,
		bypassExts:                 blockexeConfig.bypassExts,
		processExts:                blockexeConfig.processExts,
		rejectExts:                 blockexeConfig.rejectExts,
		extArrs:                    blockexeConfig.extArrs,
		returnOrigIfMaxSizeExc:     blockexeConfig.returnOrigIfMaxSizeExc,
		return400IfFileExtRejected: blockexeConfig.return400IfFileExtRejected,
	}
}
