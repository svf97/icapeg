package blockexe

import (
	"bytes"
	utils "icapeg/consts"
	"icapeg/logging"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Processing is a func used for to processing the http message
func (be *blockexe) Processing(partial bool) (int, interface{}, map[string]string) {
	logging.Logger.Info(be.serviceName + " service has started processing")
	serviceHeaders := make(map[string]string)
	// no need to scan part of the file, this service needs all the file at ine time
	if partial {
		logging.Logger.Info(be.serviceName + " service has stopped processing partially")
		return utils.Continue, nil, nil
	}
	isGzip := false

	file, reqContentType, err := be.generalFunc.CopyingFileToTheBuffer(be.methodName)
	if err != nil {
		logging.Logger.Error(be.serviceName + " error: " + err.Error())
		logging.Logger.Info(be.serviceName + " service has stopped processing")
		return utils.InternalServerErrStatusCodeStr, nil, serviceHeaders
	}

	//getting the extension of the file
	var contentType []string
	if len(contentType) == 0 {
		contentType = append(contentType, "")
	}
	var fileName string
	if be.methodName == utils.ICAPModeReq {
		contentType = be.httpMsg.Request.Header["Content-Type"]
		fileName = be.generalFunc.GetFileName()
	} else {
		contentType = be.httpMsg.Response.Header["Content-Type"]
		fileName = be.generalFunc.GetFileName()
	}
	if len(contentType) == 0 {
		contentType = append(contentType, "")
	}
	fileExtension := be.generalFunc.GetMimeExtension(file.Bytes(), contentType[0], fileName)

	//check if the file extension is a bypass extension
	//if yes we will not modify the file, and we will return 204 No modifications
	isProcess, icapStatus, httpMsg := be.generalFunc.CheckTheExtension(fileExtension, be.extArrs,
		be.processExts, be.rejectExts, be.bypassExts, be.return400IfFileExtRejected, isGzip,
		be.serviceName, be.methodName, BEIdentifier, be.httpMsg.Request.RequestURI, reqContentType, file)
	if !isProcess {
		logging.Logger.Info(be.serviceName + " service has stopped processing")
		return icapStatus, httpMsg, serviceHeaders
	}

	//check if the file size is greater than max file size of the service
	//if yes we will return 200 ok or 204 no modification, it depends on the configuration of the service
	if be.maxFileSize != 0 && be.maxFileSize < file.Len() {
		status, file, httpMsg := be.generalFunc.IfMaxFileSizeExc(be.returnOrigIfMaxSizeExc, be.serviceName, be.methodName, file, be.maxFileSize)
		fileAfterPrep, httpMsg := be.generalFunc.IfStatusIs204WithFile(be.methodName, status, file, isGzip, reqContentType, httpMsg, true)
		if fileAfterPrep == nil && httpMsg == nil {
			logging.Logger.Info(be.serviceName + " service has stopped processing")
			return utils.InternalServerErrStatusCodeStr, nil, serviceHeaders
		}
		switch msg := httpMsg.(type) {
		case *http.Request:
			msg.Body = io.NopCloser(bytes.NewBuffer(fileAfterPrep))
			logging.Logger.Info(be.serviceName + " service has stopped processing")
			return status, msg, nil
		case *http.Response:
			msg.Body = io.NopCloser(bytes.NewBuffer(fileAfterPrep))
			logging.Logger.Info(be.serviceName + " service has stopped processing")
			return status, msg, nil
		}
		return status, nil, nil
	}

	scannedFile := file.Bytes()

	//returning the scanned file if everything is ok
	scannedFile = be.generalFunc.PreparingFileAfterScanning(scannedFile, reqContentType, be.methodName)
	logging.Logger.Info(be.serviceName + " service has stopped processing")
	return utils.OkStatusCodeStr, be.generalFunc.ReturningHttpMessageWithFile(be.methodName, scannedFile), serviceHeaders

}

func (be *blockexe) ISTagValue() string {
	epochTime := strconv.FormatInt(time.Now().Unix(), 10)
	return "epoch-" + epochTime
}
