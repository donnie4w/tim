package serializer

import (
	"context"
	"fmt"
	"testing"
	"tim/base64Util"
	"tim/protocol"

	"github.com/apache/thrift/lib/go/thrift"
)

func Test_ser(t *testing.T) {
	mbean := protocol.NewTimMBean()
	body := "wuxiaodong"
	mbean.Body = &body
	b, _ := thrift.NewTSerializer().Write(context.Background(), mbean)
	base64str := string(base64Util.Base64Encode(b))
	fmt.Println(">>>>>>>>", base64str)
	var mbean2 *protocol.TimMBean = protocol.NewTimMBean()
	bb, _ := base64Util.Base64Decode(base64str)
	thrift.NewTDeserializer().Read(context.Background(), mbean2, bb)
	fmt.Println(mbean2)
	fmt.Println(*mbean2.Body)
}
