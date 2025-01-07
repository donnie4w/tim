// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/tim

package log

import (
	"fmt"
	"github.com/donnie4w/go-logger/logger"
)

var log = logger.NewLogger()

func SetFile(filepath string) {
	log.SetOption(&logger.Option{FileOption: &logger.FileMixedMode{Filename: filepath, Maxsize: 1 << 20, Maxbuckup: 30, IsCompress: true, Timemode: logger.MODE_DAY}})
	log.SetConsole(true)
}

func SetConsole(ok bool) {
	log.SetConsole(ok)
}

func SetLevel(level logger.LEVELTYPE) {
	log.SetLevel(level)
}

func Write(bs []byte) {
	log.Write(bs)
}

func Debug(v ...interface{}) {
	log.Debug(v...)
}

func FmtPrint(v ...interface{}) {
	fmtPrint(log, v...)
}

func Info(v ...interface{}) {
	log.Info(v...)
}

func Warn(v ...interface{}) {
	log.Warn(v...)
}

func Error(v ...interface{}) {
	log.Error(v...)
}

func Debugf(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

func Infof(format string, v ...interface{}) {
	log.Infof(format, v...)
}

func Warnf(format string, v ...interface{}) {
	log.Warnf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	log.Errorf(format, v...)
}

func fmtPrint(log *logger.Logging, v ...any) {
	info := fmt.Sprint(v...)
	a, b := "", ""
	ll := 80
	if ll >= len(info) {
		for i := 0; i < (ll-len(info))/2; i++ {
			a = a + "="
		}
		b = a
		if ll > len(info)+len(a)*2 {
			b = a + "="
		}
	}
	log.Info(a, info, b)
}
