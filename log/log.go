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
	"strings"
	"time"
)

var log = logger.NewLogger()

func init() {
	SetFile("tim.log")
}

func SetFile(filepath string) {
	log.SetOption(&logger.Option{CallDepth: 1, FileOption: &logger.FileSizeMode{Filename: filepath, Maxsize: 1 << 30, Maxbuckup: 30, IsCompress: true}})
	log.SetLevelOption(logger.LEVEL_WARN, &logger.LevelOption{Format: logger.FORMAT_NANO})
	log.SetConsole(true).SetFormat(logger.FORMAT_DATE | logger.FORMAT_TIME)
}

func SetConsole(ok bool) {
	log.SetConsole(ok)
}

func SetLevel(level logger.LEVELTYPE) {
	log.SetLevel(level)
}

func Debug(v ...interface{}) {
	log.Debug(v...)
}

func FmtPrint(v ...interface{}) {
	info := fmt.Sprint(v...)
	const ll = 80
	if len(info) >= ll {
		log.Info(info)
		return
	}
	padLen := (ll - len(info)) / 2
	padding := strings.Repeat("═", padLen)
	extraPadding := ""
	if len(info)+padLen*2 < ll {
		extraPadding = " "
	}
	log.Infof("%s %s%s %s", padding, extraPadding, info, padding)
}

func DelayPrint(v ...any) {
	go func() {
		time.Sleep(2 * time.Second)
		log.Warn(fmt.Sprint(v...))
	}()
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
