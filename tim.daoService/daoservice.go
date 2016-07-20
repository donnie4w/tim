/**
 * donnie4w@gmail.com  tim server
 */
package daoService

import (
	"fmt"
	"runtime/debug"
	"strconv"
	//	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/go-logger/logger"
	//	"tim.DB"
	"tim.base64Util"
	. "tim.common"
	"tim.dao"
	. "tim.protocol"
	"tim.utils"
)

/*保存离线信息*/
func SaveOfflineMBean(mbean *TimMBean) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("SaveOfflineMBean,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	if ConfBean.Db_Exsit != 1 {
		return
	}
	tim_offline := dao.NewTim_offline()
	mid, _ := strconv.Atoi(mbean.GetMid())
	tim_offline.SetMid(int64(mid))
	tim_offline.SetDomain(mbean.FromTid.GetDomain())
	tim_offline.SetFromuser(mbean.GetFromTid().GetName())
	tim_offline.SetCreatetime(utils.NowTime())
	tim_offline.SetUsername(mbean.GetToTid().GetName())
	tim_offline.SetStamp(utils.TimeMills())
	mbean.Offline = NewTimTime()
	mbean.Offline.Timestamp = mbean.Timestamp
	stanza, _ := thrift.NewTSerializer().Write(mbean)
	base64string := base64Util.Base64Encode(stanza)
	length := len([]byte(base64string))
	tim_offline.SetStanza(base64string)
	tim_offline.SetMessage_size(int64(length))
	tim_offline.Insert()
	go UpdateOffMessage(mbean, 0)
}

/*load 离线信息*/
func LoadOfflineMBean(tid *Tid) (mbeans []*TimMBean) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("LoadOfflineMBean,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	if ConfBean.Db_Exsit != 1 {
		return
	}
	tim_offline := dao.NewTim_offline()
	tim_offline.Where(tim_offline.Domain.EQ(tid.GetDomain()), tim_offline.Username.EQ(tid.GetName()))
	tim_offline.OrderBy(tim_offline.Id.Asc())
	offlines, err := tim_offline.Selects()
	if err == nil {
		mbeans = make([]*TimMBean, 0)
		for _, of := range offlines {
			var timmbean *TimMBean = NewTimMBean()
			bb, er := base64Util.Base64Decode(of.GetStanza())
			if er == nil {
				thrift.NewTDeserializer().Read(timmbean, []byte(bb))
				mbeans = append(mbeans, timmbean)
			} else {
				logger.Error("Base64Decode:", er)
			}
		}
	}
	return
}

/*删除指定信息*/
func DelteOfflineMBean(mid *string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("DelteOfflineMBean,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	if ConfBean.Db_Exsit != 1 {
		return
	}
	tim_offline := dao.NewTim_offline()
	tim_offline.Where(tim_offline.Mid.EQ(mid))
	tim_offline.Delete()
}

/*保存信息*/
func SaveMBean(mbean *TimMBean) (id int32, timestamp string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("SaveMBean,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	if ConfBean.Db_Exsit != 1 {
		return
	}
	return _saveMBean(mbean, 1, 1)
}

/*保存信息*/
func SaveSingleMBean(mbean *TimMBean) (id int32, timestamp string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("SaveMBean,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	if ConfBean.Db_Exsit != 1 {
		return
	}
	fromname := mbean.FromTid.GetName()
	toname := mbean.ToTid.GetName()
	small, large := 0, 0
	if toname > fromname {
		large = 1
	} else {
		small = 1
	}
	return _saveMBean(mbean, small, large)
}

/*保存信息*/
func _saveMBean(mbean *TimMBean, small, large int) (id int32, timestamp string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("_saveMBean,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	logger.Debug("mbean====>", mbean)
	if ConfBean.Db_Exsit != 1 {
		return
	}
	domain := mbean.GetFromTid().GetDomain()
	fromname := mbean.GetFromTid().GetName()
	toname := mbean.GetToTid().GetName()
	message := dao.NewTim_message()
	chatid := utils.Chatid(fromname, toname, domain)
	message.SetChatid(chatid)
	timestamp = mbean.GetTimestamp()
	message.SetStamp(timestamp)
	message.SetCreatetime(utils.NowTime())
	message.SetFromuser(fromname)
	message.SetTouser(toname)
	message.SetSmall(int64(small))
	message.SetLarge(int64(large))
	stanza, _ := thrift.NewTSerializer().Write(mbean)
	stanzastr := string(base64Util.Base64Encode(stanza))
	message.SetStanza(stanzastr)
	message.Insert()

	mess := dao.NewTim_message()
	mess.Where(mess.Stamp.EQ(timestamp), mess.Chatid.EQ(chatid))
	var err error
	mess, err = mess.Select()
	if err == nil {
		id = mess.GetId()
	}
	return
}

/**
  离线信息发送成功后 更新 small或large 状态
*/
func UpdateOffMessage(mbean *TimMBean, status int) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("UpdateOffMessage", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	if ConfBean.Db_Exsit != 1 {
		return
	}
	//	domain := mbean.GetFromTid().GetDomain()
	fromname := mbean.GetFromTid().GetName()
	toname := mbean.GetToTid().GetName()
	//	chatid := utils.Chatid(fromname, toname, domain)
	message := dao.NewTim_message()
	if toname < fromname {
		message.SetSmall(int64(status))
	} else {
		message.SetLarge(int64(status))
	}
	message.Where(message.Id.EQ(mbean.GetMid()))
	message.Update()
}

/***/
func LoadMBean(fidname, tidname, domain string, fromstamp, tostamp *string, limitcount int32) (tms []*TimMBean) {
	logger.Debug("LoadMBean:", fidname, " ", tidname, " ", domain, " ", fromstamp, " ", tostamp, " ", limitcount)
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	chatid := utils.Chatid(fidname, tidname, domain)
	isLarge := fidname > tidname
	timMessage := dao.NewTim_message()
	wheres := make([]*gdao.Where, 0)

	if fromstamp != nil && tostamp != nil {
		wheres = append(wheres, timMessage.Stamp.Between(*fromstamp, *tostamp))
	} else if fromstamp != nil {
		wheres = append(wheres, timMessage.Stamp.GT(*fromstamp))
	} else if tostamp != nil {
		wheres = append(wheres, timMessage.Stamp.LT(*tostamp))
	}
	wheres = append(wheres, timMessage.Chatid.EQ(chatid))
	if isLarge {
		wheres = append(wheres, timMessage.Large.EQ(1))
	} else {
		wheres = append(wheres, timMessage.Small.EQ(1))
	}
	timMessage.Where(wheres...)
	timMessage.OrderBy(timMessage.Id.Desc())
	if limitcount > 0 {
		timMessage.Limit(0, limitcount)
	}
	timMessages, err := timMessage.Selects()
	if err == nil && timMessages != nil {
		tms = make([]*TimMBean, 0)
		for _, msg := range timMessages {
			tm := new(TimMBean)
			bb, er := base64Util.Base64Decode(msg.GetStanza())
			if er == nil {
				thrift.NewTDeserializer().Read(tm, bb)
				mid := fmt.Sprint(msg.GetId())
				tm.Mid = &mid
				tms = append(tms, tm)
			} else {
				logger.Error("Base64Decode:", er)
			}
		}
	}
	return
}

func DelMBean(fidname, tidname, domain, mid string) {
	logger.Debug("DelMBean:", fidname, " ", tidname, " ", domain, " ", mid)
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	chatid := utils.Chatid(fidname, tidname, domain)
	isLarge := fidname > tidname
	timMessage := dao.NewTim_message()
	if isLarge {
		timMessage.SetLarge(0)
	} else {
		timMessage.SetSmall(0)
	}
	timMessage.Where(timMessage.Chatid.EQ(chatid), timMessage.Id.EQ(mid))
	timMessage.Update()
}

func DelAllMBean(fidname, tidname, domain string) {
	logger.Debug("DelAllMBean:", fidname, " ", tidname, " ", domain)
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	chatid := utils.Chatid(fidname, tidname, domain)
	isLarge := fidname > tidname
	timMessage := dao.NewTim_message()
	if isLarge {
		timMessage.SetLarge(0)
	} else {
		timMessage.SetSmall(0)
	}
	timMessage.Where(timMessage.Chatid.EQ(chatid))
	timMessage.Update()
}

///*lastTime 时间之后的消息*/
//func LoadMBean(fid, tid *Tid, lastTime time.Time) (mbeans []*TimMBean) {
//	return
//}

/**ip地址是否被限制*/
func AllowHttpIp(ip string) bool {
	return true
}

func IsTidExist(tid *Tid) bool {
	return true
}

func Auth(domain, loginName, pwd string) bool {
	return true
}

func CheckDomain(domain string) bool {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	if ConfBean.Db_Exsit != 1 {
		return true
	}
	tim_domain := dao.NewTim_domain()
	tim_domain.Where(tim_domain.Domain.EQ(domain))
	var err error
	tim_domain, err = tim_domain.Select()
	if err == nil && tim_domain != nil && tim_domain.GetId() > 0 {
		return true
	}
	return false
}

func AddConf() {
	logger.Debug("Addconf ok")
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	if ConfBean.Db_Exsit != 1 {
		return
	}
	tim_config := dao.NewTim_config()
	confs, err := tim_config.Selects()
	if err == nil && confs != nil {
		for _, conf := range confs {
			if conf.GetKeyword() != "" && conf.GetValuestr() != "" {
				ConfBean.KV[conf.GetKeyword()] = conf.GetValuestr()
			}
		}
	}
}

//
func GetOnlineRoser(fromtid *Tid) []*Tid {
	domain := fromtid.GetDomain()
	fromname := fromtid.GetName()
	logger.Debug(domain, " ", fromname)
	return nil
}
