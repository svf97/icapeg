package service

import (
	http_message "icapeg/http-message"
	"icapeg/logging"
	"icapeg/service/services/blockexe"
	"icapeg/service/services/clamav"
	"icapeg/service/services/cloudmersive"
	"icapeg/service/services/echo"
	"icapeg/service/services/grayimages"
	hashlookuppackage "icapeg/service/services/hashlookup"
	"icapeg/service/services/virustotal"
)

// Vendors names
const (
	VendorEcho         = "echo"
	VendorClamav       = "clamav"
	VendorVirustotal   = "virustotal"
	VendorCloudMersive = "cloudmersive"
	VendorGrayimages   = "grayimages"
	VendorHashlookup   = "hashlookup"
	VendorBlockexe     = "blockexe"
)

type (
	// Service holds the info to distinguish a service
	Service interface {
		Processing(bool) (int, interface{}, map[string]string)
		ISTagValue() string
	}
)

// GetService returns a service based on the service name
// change name to vendor and add parameter service name
func GetService(vendor, serviceName, methodName string, httpMsg *http_message.HttpMsg) Service {
	logging.Logger.Info("getting instance from " + serviceName + " struct")
	switch vendor {
	case VendorEcho:
		return echo.NewEchoService(serviceName, methodName, httpMsg)
	case VendorVirustotal:
		return virustotal.NewVirustotalService(serviceName, methodName, httpMsg)
	case VendorClamav:
		return clamav.NewClamavService(serviceName, methodName, httpMsg)
	case VendorCloudMersive:
		return cloudmersive.NewCloudMersiveService(serviceName, methodName, httpMsg)
	case VendorGrayimages:
		return grayimages.NewGrayimagesService(serviceName, methodName, httpMsg)
	case VendorHashlookup:
		return hashlookuppackage.NewHashlookupService(serviceName, methodName, httpMsg)
	case VendorBlockexe:
		return blockexe.NewBlockexeService(serviceName, methodName, httpMsg)
	}

	return nil
}

// InitServiceConfig is used to load the services configuration
func InitServiceConfig(vendor, serviceName string) {
	logging.Logger.Info("loading all the services configuration")
	switch vendor {
	case VendorEcho:
		echo.InitEchoConfig(serviceName)
	case VendorClamav:
		clamav.InitClamavConfig(serviceName)
	case VendorVirustotal:
		virustotal.InitVirustotalConfig(serviceName)
	case VendorCloudMersive:
		cloudmersive.InitCloudMersiveConfig(serviceName)
	case VendorGrayimages:
		grayimages.InitGrayimagesConfig(serviceName)
	case VendorHashlookup:
		hashlookuppackage.InitHashlookupConfig(serviceName)
	case VendorBlockexe:
		blockexe.InitBlockexeConfig(serviceName)
	}
}
