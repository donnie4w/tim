/**
 * donnie4w@gmail.com  tim server
 */
package impl

import (
	. "tim.protocol"
)

func newTid(name string, domain, resource *string) *Tid {
	tid := NewTid()
	tid.Domain = domain
	tid.Name = name
	tid.Resource = resource
	return tid
}
