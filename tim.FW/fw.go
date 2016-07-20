/**
 * donnie4w@gmail.com  tim server
 */
/**流程*/
package FW

type FLOW int

const (
	/**已连接*/
	CONNECT FLOW = 1
	/**已认证*/
	AUTH FLOW = 2
	/**已关闭*/
	CLOSE FLOW = 3
)
