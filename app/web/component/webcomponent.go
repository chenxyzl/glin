package component

import (
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/chenxyzl/glin/grain"
	"github.com/chenxyzl/glin/slog"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slices"
	"laiya/config"
	"laiya/model/home_model"
	"laiya/model/sub_model"
	"laiya/model/web_model"
	"laiya/model/web_model/web_old_model"
	"laiya/proto/code"
	"laiya/proto/inner"
	"laiya/proto/outer"
	"laiya/share/call"
	"laiya/share/oss"
	"laiya/share/sms"
	token2 "laiya/share/token"
	"laiya/share/utils"
	"laiya/web/iface"
	"mime/multipart"
	"regexp"
	"time"
)

var phPattern = regexp.MustCompile(`^1[3456789]\d{9}$`)

type WebComponent struct {
	grain.BaseComponent
	iface.IWebEntity
}

func (w *WebComponent) BeforeInit() error {
	w.IWebEntity = w.GetEntity().(iface.IWebEntity)
	return nil
}
func (w *WebComponent) AfterInit() error {
	return nil
}                                        //actor Init完成前
func (w *WebComponent) BeforeTerminate() {} //actor Terminate完成前
func (w *WebComponent) AfterTerminate()  {} //actor Terminate完成前
func (w *WebComponent) Tick()            {}
func (w *WebComponent) HandleGetCaptcha(req *outer.GetCaptcha_Request, ctx *gin.Context) (*outer.GetCaptcha_Reply, code.Code) {
	if slices.Contains(config.Get().ConstConfig.GodAcc, req.GetPh()) {
		slog.Infof("god account login, ph:%v", req.GetPh())
		return &outer.GetCaptcha_Reply{}, code.Code_Ok
	}
	//手机号合法检查
	if !phPattern.MatchString(req.GetPh()) {
		slog.Errorf("ph illegal, ph:%v", req.GetPh())
		return nil, code.Code_PhNotIllegal
	}
	//加载db
	mod := &web_model.Captcha{
		Ph: req.GetPh(),
	}
	err := mod.Load()
	if err != nil {
		slog.Errorf("load captcha err, ph:%v|err:%v", req.GetPh(), err)
		return nil, code.Code_InnerError
	}
	//是否过滤频繁
	now := time.Now()
	if mod.LastReqTime != 0 && mod.LastReqTime > now.UnixNano()-int64(config.Get().WebConfig.CaptchaReqIntervalTime) {
		return nil, code.Code_CaptchaTooFrequently
	}

	//生成code
	number := int64(now.UnixNano()) / 100 % 10000
	captcha := fmt.Sprintf("%04d", number)

	//发送验证码
	cod := sms.SendCode(req.GetPh(), captcha)
	if cod != code.Code_Ok {
		return nil, cod
	}

	//更新db
	mod.LastReqTime = now.UnixNano()
	mod.Captcha = captcha
	mod.CaptchaExpTime = now.UnixNano() + int64(config.Get().WebConfig.CaptchaExpireTime)
	err = mod.Save()
	if err != nil {
		slog.Errorf("save captcha err, err:%v", err)
		return nil, code.Code_InnerError
	}

	slog.Infof("get captcha ok, ph:%v|cod:%v", req.GetPh(), captcha)
	return &outer.GetCaptcha_Reply{}, code.Code_Ok
}

// 检查验证码
func (w *WebComponent) checkCaptcha(ph, captcha string) (code.Code, error) {
	mod := &web_model.Captcha{Ph: ph}
	err := mod.Load()
	if err != nil {
		return code.Code_InnerError, fmt.Errorf("load captcha err, ph:%v|err:%v", ph, err)
	}
	//验证码存在--
	if mod.Captcha == "" {
		return code.Code_NotHasCaptcha, fmt.Errorf("captcha not found, ph:%v", ph)
	}
	//验证码一致--
	if mod.Captcha != captcha {
		return code.Code_CaptchaNotEeq, fmt.Errorf("captcha not eq, ph:%v|db:%v|req:%v", ph, mod.Captcha, captcha)
	}
	//验证码过期--
	now := time.Now()
	if now.UnixNano() > mod.CaptchaExpTime {
		return code.Code_CaptchaExpire, fmt.Errorf("captcha expire, ph:%v|req:%v|expTime:%v|now:%v", ph, mod.Captcha, mod.CaptchaExpTime, now)
	}
	//使用验证码--
	mod.Captcha = ""
	mod.CaptchaExpTime = 0
	mod.LoginTimes++
	err = mod.Save()
	if err != nil {
		return code.Code_InnerError, fmt.Errorf("save captcha date err, ph:%v|err:%v", ph, err)
	}
	return code.Code_Ok, nil
}

func (w *WebComponent) HandleCheckCaptcha(req *outer.CheckCaptcha_Request, ctx *gin.Context) (*outer.CheckCaptcha_Reply, code.Code) {
	//校验验证码
	if slices.Contains(config.Get().ConstConfig.GodAcc, req.GetPh()) && req.GetCaptcha() == config.Get().ConstConfig.GodCaptcha {
		slog.Infof("god account check ok, ph:%v|acc:%v", req.GetPh(), req.GetCaptcha())
	} else {
		cod, err := w.checkCaptcha(req.GetPh(), req.GetCaptcha())
		if cod != code.Code_Ok {
			slog.Error(err)
			return nil, cod
		}
	}
	//获取账号数据
	account, isNew, err := web_model.GetAccountData(req.GetPh())
	if err != nil {
		slog.Errorf("get account data err, ph:%v|err:%V", req.GetPh(), err)
		return nil, code.Code_InnerError
	}
	//
	state := outer.WebLoginUseInfo_Registered
	//直接从db获取数据
	playerHead := &sub_model.PlayerHead{Uid: account.Uid}
	err = playerHead.Load()
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		slog.Errorf("get account data err, ph:%v|err:%V", req.GetPh(), err)
		return nil, code.Code_InnerError
	}

	//如果名字为空证明用户不存在--不存在则加载老数据
	if errors.Is(err, mongo.ErrNoDocuments) {
		//加载老数据
		oldPlayer := &web_old_model.OldPlayer{}
		err = oldPlayer.LoadOldPlayerByAccount(account.Account)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			slog.Errorf("get old data err, ph:%v|err:%V", req.GetPh(), err)
			return nil, code.Code_InnerError
		}
		//构造数据
		if err != nil { //新用户注册
			playerHead = &sub_model.PlayerHead{Uid: account.Uid, Name: utils.Uint64ToHashedBase52(account.Uid), Icon: config.Get().ConstConfig.UserDefaultIcon, Des: ""}
			state = outer.WebLoginUseInfo_New
		} else { //老数据则继承
			playerHead = &sub_model.PlayerHead{Uid: account.Uid, Name: oldPlayer.Name, Icon: oldPlayer.Icon, Des: oldPlayer.Des}
			state = outer.WebLoginUseInfo_Old
		}
		_, err := home_model.CreatePlayer(playerHead.Uid, account.Account, playerHead.Name, playerHead.Icon, playerHead.Des)
		if err != nil {
			slog.Errorf("create player err, account:%v|uid:%v|err:%v", account.Account, account.Uid, err)
			return nil, code.Code_InnerError
		}
	}
	// 生成token
	tokenString, err := token2.BuildToken(account.Uid, config.Get().AConfig.TokenExpireTime, config.Get().AConfig.AppKey)
	if err != nil {
		slog.Errorf("token error, ph:%v|cod:%s", req.GetPh(), err)
		return nil, code.Code_InnerError
	}
	//
	userInfo := &outer.WebLoginUseInfo_UserInfo{
		Uid:   playerHead.Uid,
		Name:  playerHead.Name,
		Icon:  playerHead.Icon,
		Des:   playerHead.Des,
		State: state,
	}

	slog.Infof("check captcha ok, ph:%v|cod:%v|isNew:%v|state:%v|token:%v|userInfo:%v", req.GetPh(), req.GetCaptcha(), isNew, state, tokenString, userInfo)
	return &outer.CheckCaptcha_Reply{Token: tokenString, UserInfo: userInfo}, code.Code_Ok
}

func (w *WebComponent) HandleCheckToken(req *outer.CheckToken_Request, ctx *gin.Context) (*outer.CheckToken_Reply, code.Code) {
	//检查uid
	uid := ctx.MustGet("uid").(uint64)
	if uid == 0 {
		slog.Errorf("not set uid, uid:%v", uid)
		return nil, code.Code_UidNotFound
	}
	//直接从db获取数据
	playerCheck := sub_model.PlayerDeleteCheck{Uid: uid}
	err := playerCheck.Load()
	//是否存在
	if errors.Is(err, mongo.ErrNoDocuments) {
		slog.Errorf("playerCheck not exist, uid:%v|err:%V", uid, err)
		//return nil, code.Code_UserNotExist
		return nil, code.Code_TokenErr
	}
	//其他错误
	if err != nil {
		slog.Errorf("get playerCheck data err, uid:%v|err:%V", uid, err)
		return nil, code.Code_InnerError
	}
	if playerCheck.IsDeleted {
		slog.Errorf("playerCheck deleted, uid:%v|err:%V", uid, err)
		//return nil, code.Code_UserNotExist
		return nil, code.Code_TokenErr
	}
	//
	slog.Infof("check token ok, uid:%v|token:%v", uid, ctx.MustGet("token"))
	return &outer.CheckToken_Reply{}, code.Code_Ok
}

func (w *WebComponent) HandleUploadIcon(req *outer.UploadIcon_Request, ctx *gin.Context) (*outer.UploadIcon_Reply, code.Code) {
	//检查uid
	uid := ctx.MustGet("uid").(uint64)
	if uid == 0 {
		slog.Errorf("not set uid, uid:%v", uid)
		return nil, code.Code_UidNotFound
	}

	//获取type
	iconTypeInt, ok := outer.UploadIcon_IconType_value[ctx.PostForm("type")]
	if !ok {
		slog.Error("not find form-data, name:type")
		return nil, code.Code_NotFoundUpdateFileNamedIcon
	}
	iconType := outer.UploadIcon_IconType(iconTypeInt)
	//获取icon
	files, err := ctx.FormFile("icon")
	if err != nil {
		slog.Error("not find form-data, name:icon")
		return nil, code.Code_NotFoundUpdateFileNamedIcon
	}
	//open
	file, err := files.Open()
	if err != nil {
		slog.Errorf("file open error, err:%v", err)
		return nil, code.Code_Error
	}
	//defer close
	defer func(file multipart.File) {
		err = file.Close()
		if err != nil {
			slog.Errorf("file close error, err:%v", err)
		}
	}(file)

	var client *oss.Bucket
	var url string
	var name = fmt.Sprintf("%v_%v_%v.jpg", iconType.String(), uid, time.Now().UnixMilli())
	//不同类型的上传到不同仓库
	if iconType == outer.UploadIcon_ChatPic {
		client, err = oss_helper.GetPicBuket()
		url = oss_helper.GetPicUrl(name)
	} else {
		client, err = oss_helper.GetHeadBuket()
		url = oss_helper.GetHeadUrl(name)
	}
	if err != nil {
		slog.Errorf("buket get err, err:%v", err)
		return nil, code.Code_InnerError
	}
	err = client.PutObject(name, file)
	if err != nil {
		slog.Errorf("put file to buket err, err:%v", err)
		return nil, code.Code_InnerError
	}
	return &outer.UploadIcon_Reply{Url: url}, code.Code_Ok
}

func (w *WebComponent) HandleGetGameTypes(req *outer.GetGameTypes_Request, ctx *gin.Context) (*outer.GetGameTypes_Reply, code.Code) {
	return &outer.GetGameTypes_Reply{GameTypes: config.Get().GameTypesConfig.GameTypes}, code.Code_Ok
}

func (w *WebComponent) HandleGetThirdPlatformConfig(req *outer.GetThirdPlatformConfig_Request, ctx *gin.Context) (*outer.GetThirdPlatformConfig_Reply, code.Code) {
	conf := config.Get().GamePlatformsConfig.PlatformTypes
	items := make(map[int32]*outer.GetThirdPlatformConfig_ThirdPlatform)
	for id, item := range conf {
		items[id] = &outer.GetThirdPlatformConfig_ThirdPlatform{
			Id:      id,
			Name:    item.Name,
			Icon:    item.Icon,
			SortIdx: item.Idx,
		}
	}
	return &outer.GetThirdPlatformConfig_Reply{ThirdPlatformConfig: items}, code.Code_Ok
}

func (w *WebComponent) HandleGetHallBannersConfig(req *outer.GetHallBannersConfig_Request, ctx *gin.Context) (*outer.GetHallBannersConfig_Reply, code.Code) {
	conf := config.Get().HallBannersConfig.Banners
	var items []*outer.GetHallBannersConfig_Banner
	for _, item := range conf {
		items = append(items, &outer.GetHallBannersConfig_Banner{
			Name: item.Name,
			Des:  item.Des,
			Icon: item.Icon,
			Link: item.Link,
		})
	}
	return &outer.GetHallBannersConfig_Reply{Banners: items}, code.Code_Ok
}

func (w *WebComponent) HandleGetClientDynamicConfig(req *outer.GetClientDynamicConfig_Request, ctx *gin.Context) (*outer.GetClientDynamicConfig_Reply, code.Code) {
	return &outer.GetClientDynamicConfig_Reply{Value: config.Get().ClientDynamic.GetConfig(req.GetKey())}, code.Code_Ok
}

func (w *WebComponent) HandleDeleteAccount(req *outer.DeleteAccount_Request, ctx *gin.Context) (*outer.DeleteAccount_Reply, code.Code) {
	//检查uid
	senderUid := ctx.MustGet("uid").(uint64)
	if senderUid == 0 {
		slog.Errorf("not set uid, uid:%v", senderUid)
		return nil, code.Code_UidNotFound
	}
	//god
	god := &sub_model.PlayerAccountCheck{Uid: senderUid}
	if err := god.Load(); err != nil {
		slog.Errorf("load god data err, uid:%v|err:%v", senderUid, err)
		return nil, code.Code_UidNotFound
	}

	if !slices.Contains(config.Get().GmConfig.Gods, god.Account) {
		slog.Errorf("not god account, opAccount:%v", god.Account)
		return nil, code.Code_NoAccess
	}

	ret, cod := call.RequestRemote[*inner.Web2HoDeletePlayer_Reply](req.GetUid(), &inner.Web2HoDeletePlayer_Request{})
	if cod != code.Code_Ok && cod != code.Code_UserNotExist {
		slog.Errorf("delete player err, cod:%v", cod)
		return nil, cod
	}
	account := &web_model.Account{Account: ret.GetAccount()}
	err := account.Load()
	if errors.Is(err, mongo.ErrNoDocuments) {
		slog.Infof("not found account data, must check by gm, account:%v", account.Account)
		return &outer.DeleteAccount_Reply{}, code.Code_Ok
	} else if err != nil {
		slog.Errorf("delete account load date err, must check by gm, account:%v", account.Account)
		return nil, code.Code_InnerError
	}
	if account.Uid != req.GetUid() {
		slog.Errorf("delete account err, uid not equal, must check by gm, account:%v|hasUid%v|reqUid:%v",
			account.Account, account.Uid, req.GetUid())
		return nil, code.Code_InnerError
	}
	err = account.Delete()
	if err != nil {
		slog.Errorf("delete account err,  must check by gm, account:%v|uid%v",
			account.Account, account.Uid)
		return nil, code.Code_InnerError
	}
	slog.Infof("delete player data success, uid:%v|account:%V|god:%v", req.GetUid(), ret.GetAccount(), god.Account)
	return &outer.DeleteAccount_Reply{}, code.Code_Ok
}

func (w *WebComponent) HandleCheckClientVersion(req *outer.CheckClientVersion_Request, ctx *gin.Context) (*outer.CheckClientVersion_Reply, code.Code) {
	ret := &outer.CheckClientVersion_Reply{}
	clientVersionConf := config.Get().ClientVersion
	ret.VersionAvailable = clientVersionConf.CheckClientVersion(req.GetClientVersion())
	if !ret.VersionAvailable {
		slog.Warningf("client version check, app need update, version:%v", req.GetClientVersion())
	}
	lastVersion := clientVersionConf.LastVersion
	if lastVersion.Version != "" && req.ClientVersion != lastVersion.Version {
		ret.LastVersion = &outer.CheckClientVersion_Version{
			Version: lastVersion.Version,
			Title:   lastVersion.Title,
			Detail:  lastVersion.Detail,
		}
	}
	return ret, code.Code_Ok
}
