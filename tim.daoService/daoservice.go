/**
 * donnie4w@gmail.com  tim server
 */
package daoService

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"

	"database/sql"
	"sync"

	"github.com/donnie4w/gdao"
	"github.com/donnie4w/go-logger/logger"
	"tim.DB"
	. "tim.Map"
	"tim.base64Util"
	. "tim.common"
	"tim.connect"
	"tim.dao"
	. "tim.protocol"
	"tim.utils"
)

var authProviderDB *sql.DB
var once sync.Once
var domainmap *HashTable = NewHashTable()

func initAuthProviderDB() {
	logger.Info("initAuthProviderDB")
	authProviderDB, _ = DB.GetDB(CF.GetKV("tim.mysql.connection", ""), 100, 10)
}

func InitDaoservice() {
	AddConf()
	updateVersion()
}

/*保存离线信息列表*/
func SaveOfflineMBeanList(mbeans []*TimMBean) {
	if mbeans != nil && len(mbeans) > 0 {
		for _, mbean := range mbeans {
			SaveOfflineMBean(mbean)
		}
	}
}

/*保存离线信息*/
func SaveOfflineMBean(mbean *TimMBean) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("SaveOfflineMBean,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	if CF.Db_Exsit == 0 {
		return
	}
	if mbean.GetType() == "groupchat" {
		_saveOfflineMucBean(mbean)
	} else {
		_saveOfflineMBean(mbean)
	}
}

func _saveOfflineMBean(mbean *TimMBean) {
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

func _saveOfflineMucBean(mbean *TimMBean) {
	tim_mucoffline := dao.NewTim_mucoffline()
	tim_mucoffline.SetCreatetime(utils.NowTime())
	tim_mucoffline.SetMid(utils.Atoi64(mbean.GetMid()))
	tim_mucoffline.SetDomain(mbean.GetFromTid().GetDomain())
	tim_mucoffline.SetUsername(mbean.GetToTid().GetName())
	tim_mucoffline.SetStamp(mbean.GetTimestamp())
	tim_mucoffline.SetRoomid(mbean.GetFromTid().GetName())
	tim_mucoffline.SetMsgtype(int64(mbean.GetMsgType()))
	tim_mucoffline.Insert()
}

/*load 离线信息*/
func LoadOfflineMBean(tid *Tid) (mbeans []*TimMBean) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("LoadOfflineMBean,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	if CF.Db_Exsit == 0 {
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

func LoadOfflineMucMBean(tid *Tid) (mbeans []*TimMBean) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("LoadOfflineMucMBean,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	if CF.Db_Exsit == 0 {
		return
	}
	tim_mucoffline := dao.NewTim_mucoffline()
	tim_mucoffline.Where(tim_mucoffline.Domain.EQ(tid.GetDomain()), tim_mucoffline.Username.EQ(tid.GetName()))
	tim_mucoffline.OrderBy(tim_mucoffline.Id.Desc())
	mucofflines, err := tim_mucoffline.Selects()
	if err == nil && mucofflines != nil && len(mucofflines) > 0 {
		mids := make([]interface{}, 0)
		for _, mucoffline := range mucofflines {
			mids = append(mids, mucoffline.GetMid())
		}
		tim_mucmessage := dao.NewTim_mucmessage()
		tim_mucmessage.Where(tim_mucmessage.Id.IN(mids...))
		mucmessages, err := tim_mucmessage.Selects()
		if err == nil && mucmessages != nil && len(mucmessages) > 0 {
			mbeans := make([]*TimMBean, 0)
			for _, mucmsg := range mucmessages {
				var timmbean *TimMBean = NewTimMBean()
				bb, er := base64Util.Base64Decode(mucmsg.GetStanza())
				if er == nil {
					thrift.NewTDeserializer().Read(timmbean, []byte(bb))
					mbeans = append(mbeans, timmbean)
				} else {
					logger.Error("Base64Decode:", er)
				}
			}
		}
	}
	return
}

func LoadMucmember(roomid *Tid) (tids []*Tid) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	if CF.Db_Exsit == 0 {
		return nil
	}
	mucRoomSQL := CF.GetKV("tim.mysql.mucRoomSQL", "")
	if mucRoomSQL == "" {
		tim_mucmember := dao.NewTim_mucmember()
		tim_mucmember.Where(tim_mucmember.Domain.EQ(roomid.GetDomain()), tim_mucmember.Roomtid.EQ(roomid.GetName()))
		tim_mucmembers, err := tim_mucmember.Selects()
		if err == nil && tim_mucmembers != nil && len(tim_mucmembers) > 0 {
			tids = make([]*Tid, 0)
			for _, r := range tim_mucmembers {
				tid := NewTid()
				domain := roomid.GetDomain()
				tid.Domain = &domain
				tid.Name = r.GetTidname()
				tids = append(tids, tid)
			}
		}
	} else {
		provider()
		if authProviderDB == nil {
			logger.Error("authProviderDB is nil")
			return nil
		}
		gbbeans, err := gdao.Query(authProviderDB, mucRoomSQL, roomid.GetName())
		if err == nil && gbbeans != nil && len(gbbeans) > 0 {
			for _, gbbean := range gbbeans {
				uname := gbbean.FieldMapName["username"].ValueString()
				tid := NewTid()
				domain := roomid.GetDomain()
				tid.Domain = &domain
				tid.Name = uname
				tids = append(tids, tid)
			}
		}
	}
	return
}

func AuthMucmember(roomid, tid *Tid) (b bool) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	if CF.Db_Exsit == 0 {
		return true
	}
	mucAuthSQL := CF.GetKV("tim.mysql.mucAuthSQL", "")
	if mucAuthSQL == "" {
		tim_mucmember := dao.NewTim_mucmember()
		tim_mucmember.Where(tim_mucmember.Domain.EQ(roomid.GetDomain()), tim_mucmember.Roomtid.EQ(roomid.GetName()), tim_mucmember.Tidname.EQ(tid.GetName()))
		tim_mucmember.Limit(0, 1)
		gbbeans, err := tim_mucmember.QueryBeen(tim_mucmember.Id.Count())
		if err == nil && gbbeans != nil && len(gbbeans) > 0 {
			b = true
		}
	} else {
		provider()
		if authProviderDB == nil {
			logger.Error("authProviderDB is nil")
			return
		}
		gbbeans, err := gdao.Query(authProviderDB, mucAuthSQL, roomid.GetName(), tid.GetName())
		if err == nil && gbbeans != nil && len(gbbeans) > 0 {
			b = true
		}
	}
	return
}

/*删除指定信息*/
func DelOfflineMBean(mid *string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("DelOfflineMBean,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	if CF.Db_Exsit == 0 {
		return
	}
	tim_offline := dao.NewTim_offline()
	tim_offline.Where(tim_offline.Mid.EQ(mid))
	tim_offline.Delete()
}

func DelOfflineMucMBean(mid *string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("DelOfflineMBean,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	if CF.Db_Exsit == 0 {
		return
	}
	tim_mucoffline := dao.NewTim_mucoffline()
	tim_mucoffline.Where(tim_mucoffline.Mid.EQ(mid))
	tim_mucoffline.Delete()
}

/*删除指定信息列表*/
func DelOfflineMBeanList(mids ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("DelOfflineMBeanList,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	if CF.Db_Exsit == 0 {
		return
	}
	tim_offline := dao.NewTim_offline()
	tim_offline.Where(tim_offline.Mid.IN(mids...))
	tim_offline.Delete()
}

func DelOfflineMucMBeanList(mids ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("DelOfflineMucMBeanList,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	if CF.Db_Exsit == 0 {
		return
	}
	tim_mucoffline := dao.NewTim_mucoffline()
	tim_mucoffline.Where(tim_mucoffline.Mid.IN(mids...))
	tim_mucoffline.Delete()
}

/*保存信息*/
func SaveMBean(mbean *TimMBean) (mid string, timestamp string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("SaveMBean,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	if CF.Db_Exsit == 0 {
		return
	}
	return _saveMBean(mbean, 1, 1)
}

/*保存信息*/
func SaveSingleMBean(mbean *TimMBean) (mid string, timestamp string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("SaveMBean,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	if CF.Db_Exsit == 0 {
		if mbean.GetMid() == "" {
			mid = fmt.Sprint(utils.GetRand(100000000))
			mbean.Mid = &mid
			timestamp = mbean.GetTimestamp()
		}
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
func _saveMBean(mbean *TimMBean, small, large int) (mid string, timestamp string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("_saveMBean,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	if CF.Db_Exsit == 0 {
		if mbean.GetMid() == "" {
			mid := fmt.Sprint(utils.GetRand(100000000))
			mbean.Mid = &mid
		}
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
		mid = fmt.Sprint(mess.GetId())
		mbean.Mid = &mid
	}
	return
}

func SaveMucMBean(mbean *TimMBean) (mid string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("SaveMucMBean,", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	tim_mucmessage := dao.NewTim_mucmessage()
	tim_mucmessage.SetStamp(mbean.GetTimestamp())
	tim_mucmessage.SetFromuser(mbean.GetLeaguerTid().GetName())
	tim_mucmessage.SetRoomtidname(mbean.GetFromTid().GetName())
	tim_mucmessage.SetDomain(mbean.GetLeaguerTid().GetDomain())
	tim_mucmessage.SetMsgtype(int64(mbean.GetMsgType()))
	stanza, _ := thrift.NewTSerializer().Write(mbean)
	stanzastr := string(base64Util.Base64Encode(stanza))
	tim_mucmessage.SetStanza(stanzastr)
	tim_mucmessage.SetCreatetime(utils.NowTime())
	tim_mucmessage.Insert()

	mucmessage := dao.NewTim_mucmessage()
	mucmessage.Where(mucmessage.Stamp.EQ(mbean.GetTimestamp()), mucmessage.Fromuser.EQ(mbean.LeaguerTid.GetName()), mucmessage.Domain.EQ(mbean.LeaguerTid.GetDomain()), mucmessage.Roomtidname.EQ(mbean.GetFromTid().GetName()))
	var err error
	mucmessage, err = mucmessage.Select(mucmessage.Id)
	if err == nil {
		mid = fmt.Sprint(mucmessage.GetId())
		mbean.Mid = &mid
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
	if CF.Db_Exsit == 0 {
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

/**
  离线信息发送成功后 更新 small或large 状态  列表
*/
func UpdateOffMessageList(mbeans []*TimMBean, status int) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("UpdateOffMessageList", err)
			logger.Error(string(debug.Stack()))
		}
	}()
	if CF.Db_Exsit == 0 {
		return
	}
	if len(mbeans) == 0 {
		return
	}
	fromname := mbeans[0].GetFromTid().GetName()
	toname := mbeans[0].GetToTid().GetName()
	message := dao.NewTim_message()
	if toname < fromname {
		message.SetSmall(int64(status))
	} else {
		message.SetLarge(int64(status))
	}
	mids := make([]interface{}, 0)
	for _, mbean := range mbeans {
		mids = append(mids, mbean.GetMid())
	}
	message.Where(message.Id.IN(mids...))
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
	if CF.Db_Exsit == 0 {
		return nil
	}
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
	if CF.Db_Exsit == 0 {
		return
	}
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
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	if CF.Db_Exsit == 0 {
		return
	}
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

func Auth(tid *Tid, pwd string) (b bool) {
	if CF.MustAuth == 0 {
		return true
	}
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	authProvider_passwordSQL := CF.GetKV("tim.mysql.passwordSQL", "")
	loginname, _ := connect.GetLoginName(tid)
	if authProvider_passwordSQL == "" {
		tim_user := dao.NewTim_user()
		tim_user.Where(tim_user.Loginname.EQ(loginname))
		user, err := tim_user.Select()
		if err == nil && user != nil {
			switch CF.GetKV("authProvider.passwordType", "") {
			case "plain":
				b = eqString(user.GetEncryptedpassword(), pwd)
			case "md5":
				b = eqString(user.GetEncryptedpassword(), utils.MD5(pwd))
			case "sha1":
				b = eqString(user.GetEncryptedpassword(), utils.Sha1(pwd))
			default:
				b = eqString(user.GetEncryptedpassword(), pwd)
			}
		}
	} else {
		provider()
		if authProviderDB == nil {
			logger.Error("authProviderDB is nil")
			return false
		}
		b = _Auth4Sql(authProvider_passwordSQL, tid, pwd)
		if !b {
			authProvider_passwordSQL2 := CF.GetKV("tim.mysql.passwordSQL2", "")
			if authProvider_passwordSQL2 != "" {
				b = _Auth4Sql(authProvider_passwordSQL2, tid, pwd)
			}
		}
	}
	return
}

func _Auth4Sql(authProvider_passwordSQL string, tid *Tid, pwd string) (b bool) {
	provider()
	if authProviderDB == nil {
		logger.Error("authProviderDB is nil")
		return false
	}
	gbbean, err := gdao.Query(authProviderDB, authProvider_passwordSQL, tid.GetName())
	if err == nil && gbbean != nil && len(gbbean) == 1 {
		if bean, ok := gbbean[0].FieldMapName["password"]; ok {
			switch CF.GetKV("authProvider.passwordType", "") {
			case "plain":
				b = eqString(bean.ValueString(), pwd)
			case "md5":
				b = eqString(bean.ValueString(), utils.MD5(pwd))
			case "sha1":
				b = eqString(bean.ValueString(), utils.Sha1(pwd))
			default:
				b = eqString(bean.ValueString(), pwd)
			}
		}
	}
	return
}

func eqString(s1, s2 string) bool {
	return strings.ToUpper(s1) == strings.ToUpper(s2)
}

func provider() {
	if authProviderDB == nil && CF.GetKV("tim.mysql.connection", "") != "" {
		once.Do(initAuthProviderDB)
	}
}

func CheckDomain(domain string) bool {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	if CF.Db_Exsit == 0 {
		return true
	}
	d := domainmap.Get(domain)
	if d != nil {
		if (time.Now().UnixNano()/1000000000 - d.(int64)) < 10*60 {
			return true
		} else {
			domainmap.Del(domain)
		}
	}
	tim_domain := dao.NewTim_domain()
	tim_domain.Where(tim_domain.Domain.EQ(domain))
	var err error
	tim_domain, err = tim_domain.Select()
	if err == nil && tim_domain != nil && tim_domain.GetId() > 0 {
		domainmap.Put(domain, time.Now().UnixNano()/1000000000)
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
	if CF.Db_Exsit == 0 {
		return
	}
	tim_config := dao.NewTim_config()
	confs, err := tim_config.Selects()
	if err == nil && confs != nil && len(confs) > 0 {
		for _, conf := range confs {
			if conf.GetKeyword() != "" && conf.GetValuestr() != "" {
				CF.KV[conf.GetKeyword()] = conf.GetValuestr()
			}
		}
	}
	tim_property := dao.NewTim_property()
	propertys, err := tim_property.Selects()
	if err == nil && propertys != nil && len(propertys) > 0 {
		for _, property := range propertys {
			if property.GetKeyword() != "" && (property.GetValueint() > 0 || property.GetValuestr() != "") {
				if property.GetValuestr() != "" {
					CF.KV[property.GetKeyword()] = property.GetValuestr()
				} else if property.GetValueint() > 0 {
					CF.KV[property.GetKeyword()] = fmt.Sprint(property.GetValueint())
				}
			}
		}
	}
}

//
func GetOnlineRoser(fromtid *Tid) (tids []*Tid) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	if CF.Db_Exsit == 0 {
		return nil
	}
	domain := fromtid.GetDomain()
	fromname := fromtid.GetName()
	logger.Debug(domain, " ", fromname)
	authProvider_rosterSql := CF.GetKV("tim.mysql.rosterSQL", "")
	loginname, _ := connect.GetLoginName(fromtid)
	if authProvider_rosterSql == "" {
		tim_roster := dao.NewTim_roster()
		tim_roster.Where(tim_roster.Loginname.EQ(loginname))
		rosters, err := tim_roster.Selects()
		if err == nil && rosters != nil && len(rosters) > 0 {
			tids = make([]*Tid, 0)
			for _, r := range rosters {
				tid := NewTid()
				domain := fromtid.GetDomain()
				tid.Domain = &domain
				tid.Name = r.GetRostername()
				tids = append(tids, tid)
			}
		}
	} else {
		provider()
		if authProviderDB == nil {
			logger.Error("authProviderDB is nil")
			return nil
		}
		gbbeans, err := gdao.Query(authProviderDB, authProvider_rosterSql, fromname)
		if err == nil && gbbeans != nil && len(gbbeans) > 0 {
			for _, gbbean := range gbbeans {
				uname := gbbean.FieldMapName["roster"].ValueString()
				tid := NewTid()
				domain := fromtid.GetDomain()
				tid.Domain = &domain
				tid.Name = uname
				tids = append(tids, tid)
			}
		}
	}
	return
}

func updateVersion() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	if CF.Db_Exsit == 0 {
		return
	}
	timDomain := dao.NewTim_config()
	timDomain.Where(timDomain.Keyword.EQ("version"))
	td, err := timDomain.Select()
	if err == nil && td != nil && td.GetId() > 0 {
		timDomain = dao.NewTim_config()
		timDomain.SetValuestr(fmt.Sprint(VersionCode))
		timDomain.SetRemark(fmt.Sprint(VersionName, " | ", VersionCode, " | ", utils.NowTime()))
		timDomain.Where(timDomain.Id.EQ(td.GetId()))
		timDomain.Update()
	} else {
		timDomain = dao.NewTim_config()
		timDomain.SetValuestr(fmt.Sprint(VersionCode))
		timDomain.SetRemark(fmt.Sprint(VersionName, " | ", VersionCode, " | ", utils.NowTime()))
		timDomain.SetCreatetime(utils.NowTime())
		timDomain.SetKeyword("version")
		timDomain.Insert()
	}
}
