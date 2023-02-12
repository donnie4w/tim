/**
 * donnie4w@gmail.com  tim server
 */
package impl

import (
	"runtime/debug"

	"github.com/donnie4w/go-logger/logger"
)

func CatchErr() {
	if err := recover(); err != nil {
		logger.Error(err, "\n", string(debug.Stack()))
	}
}
