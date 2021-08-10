package log

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strconv"
	"sync"
)

type logger struct {
	out io.Writer

	// App Specific Data for the logger
	app           appInfo
	correlationID string

	isTerminal bool
}

type appInfo struct {
	Name      string                 `json:"name"`
	Version   string                 `json:"version"`
	Framework string                 `json:"framework"`
	Data      map[string]interface{} `json:"data,omitempty"`
	syncData  *sync.Map
}

func (a *appInfo) getAppData() appInfo {
	res := appInfo{}

	res.Name = a.Name
	res.Version = a.Version
	res.Framework = a.Framework
	res.Data = make(map[string]interface{})

	if a.syncData != nil {
		a.syncData.Range(func(key, value interface{}) bool {
			res.Data[key.(string)] = value
			return true
		})
	}

	return res
}

// log does the actual logging. This function creates the entry message and outputs it in color format
// in terminal context and gives out json in non terminal context. Also, sends to echo if client is present.
func (k *logger) log(level level, format string, args ...interface{}) {
	if rls.level < level {
		return // No need to do anything if we are not going to log it.
	}

	e, data, isPerformanceLog := entryFromInputs(format, args...)

	e.Level = level
	e.System = fetchSystemStats()
	e.App = k.app.getAppData()

	if isPerformanceLog {
		// in performance log, app data is under the key `appData` instead of e.App.Data.
		// This is done to ensure that appData is consistent for concurrent requests
		appData, _ := e.Data["appData"].(map[string]interface{})
		e.App.Data = appData
		delete(e.Data, "appData")
	}

	// logging struct/map in app data
	for key, val := range data {
		e.App.Data[key] = val
	}

	if k.correlationID != "" { // CorrelationID from Application Log
		e.CorrelationID = k.correlationID
	} else if correlationID, ok := e.App.Data["correlationID"]; ok {
		/*CorrelationID for middleware apart from Performance log.
		For performance log the correlationID comes from Logline struct defined in logging.go.*/
		e.CorrelationID, _ = correlationID.(string)
	}

	id, err := strconv.Atoi(e.CorrelationID)
	if err == nil && id == 0 {
		e.CorrelationID = ""
	}

	// Deleting the correlationId in case of any duplication.
	delete(e.App.Data, "correlationID")

	if k.isTerminal {
		fmt.Fprint(k.out, e.TerminalOutput())
	} else {
		_ = json.NewEncoder(k.out).Encode(e)
	}

}

func isJSON(s interface{}) (ok bool, hashmap map[string]interface{}) {
	var js map[string]interface{}

	sBytes, _ := json.Marshal(s)

	return json.Unmarshal(sBytes, &js) == nil, js
}

func (k *logger) Log(args ...interface{}) {
	k.log(Info, "", args...)
}

func (k *logger) Logf(format string, args ...interface{}) {
	k.log(Info, format, args...)
}

func (k *logger) Info(args ...interface{}) {
	k.log(Info, "", args...)
}

func (k *logger) Infof(format string, args ...interface{}) {
	k.log(Info, format, args...)
}

func (k *logger) Debug(args ...interface{}) {
	k.log(Debug, "", args...)
}

func (k *logger) Debugf(format string, args ...interface{}) {
	k.log(Debug, format, args...)
}

func (k *logger) Warn(args ...interface{}) {
	k.log(Warn, "", args...)
}

func (k *logger) Warnf(format string, args ...interface{}) {
	k.log(Warn, format, args...)
}

func (k *logger) Error(args ...interface{}) {
	k.AddData("StackTrace", string(debug.Stack()))
	k.log(Error, "", args...)
	k.removeData("StackTrace")
}

func (k *logger) Errorf(format string, args ...interface{}) {
	k.AddData("StackTrace", string(debug.Stack()))
	k.log(Error, format, args...)
	k.removeData("StackTrace")
}

func (k *logger) Fatal(args ...interface{}) {
	k.AddData("StackTrace", string(debug.Stack()))
	k.log(Fatal, "", args...)
	os.Exit(1)
}

func (k *logger) Fatalf(format string, args ...interface{}) {
	k.AddData("StackTrace", string(debug.Stack()))
	k.log(Fatal, format, args...)
	os.Exit(1)
}

func (k *logger) AddData(key string, value interface{}) {
	k.app.syncData.Store(key, value)
}

func (k *logger) removeData(key string) {
	k.app.syncData.Delete(key)
}
