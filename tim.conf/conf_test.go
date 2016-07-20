package conf

import (
	"fmt"
	"testing"
)

func Test_init(t *testing.T) {
	cf := new(ConfBean)
	cf.LoadFromXml(`D:\liteIDEspace\tim\src\tim.xml`)
	fmt.Println(cf)
}
