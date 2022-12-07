package module

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"strings"
	"teamide/internal/base"
	"teamide/internal/context"
	"teamide/internal/module/module_file_manager"
	"teamide/internal/module/module_log"
	"teamide/internal/module/module_login"
	"teamide/internal/module/module_node"
	"teamide/internal/module/module_power"
	"teamide/internal/module/module_register"
	"teamide/internal/module/module_terminal"
	"teamide/internal/module/module_toolbox"
	"teamide/internal/module/module_user"
	"teamide/pkg/util"
)

func NewApi(ServerContext *context.ServerContext) (api *Api, err error) {

	api = &Api{
		ServerContext:     ServerContext,
		userService:       module_user.NewUserService(ServerContext),
		registerService:   module_register.NewRegisterService(ServerContext),
		loginService:      module_login.NewLoginService(ServerContext),
		installService:    NewInstallService(ServerContext),
		toolboxService:    module_toolbox.NewToolboxService(ServerContext),
		nodeService:       module_node.NewNodeService(ServerContext),
		powerRoleService:  module_power.NewPowerRoleService(ServerContext),
		powerRouteService: module_power.NewPowerRouteService(ServerContext),
		powerUserService:  module_power.NewPowerUserService(ServerContext),
		logService:        module_log.NewLogService(ServerContext),
		apiCache:          make(map[string]*base.ApiWorker),
	}
	var apis []*base.ApiWorker
	apis, err = api.GetApis()
	if err != nil {
		return
	}
	for _, one := range apis {
		err = api.appendApi(one)
		if err != nil {
			return
		}
	}

	var apiPowerMap = make(map[string]bool)
	for _, one := range api.apiCache {
		apiPowerMap[one.Power.Action] = true
	}
	ps := base.GetPowers()
	for _, one := range ps {
		if !ServerContext.IsServer {
			if !one.StandAlone {
				continue
			}
		}
		_, ok := apiPowerMap[one.Action]
		if !ok {
			ServerContext.Logger.Warn("权限[" + one.Action + "]未配置动作")
		}
	}

	err = api.installService.Install()
	if err != nil {
		return
	}
	if ServerContext.IsServer {
		err = api.initServer()
		if err != nil {
			return
		}
	} else {
		err = api.initStandAlone()
		if err != nil {
			return
		}
	}
	go api.nodeService.InitContext()
	return
}

// Api ID服务
type Api struct {
	*context.ServerContext
	toolboxService    *module_toolbox.ToolboxService
	nodeService       *module_node.NodeService
	userService       *module_user.UserService
	registerService   *module_register.RegisterService
	loginService      *module_login.LoginService
	powerRoleService  *module_power.PowerRoleService
	powerRouteService *module_power.PowerRouteService
	powerUserService  *module_power.PowerUserService
	logService        *module_log.LogService
	installService    *InstallService
	apiCache          map[string]*base.ApiWorker
}

var (

	//PowerRegister 基础权限
	PowerRegister    = base.AppendPower(&base.PowerAction{Action: "register", Text: "注册", StandAlone: false})
	PowerData        = base.AppendPower(&base.PowerAction{Action: "data", Text: "数据", StandAlone: true})
	PowerSession     = base.AppendPower(&base.PowerAction{Action: "session", Text: "会话", StandAlone: true})
	PowerLogin       = base.AppendPower(&base.PowerAction{Action: "login", Text: "登录", StandAlone: false})
	PowerLogout      = base.AppendPower(&base.PowerAction{Action: "logout", Text: "登出", StandAlone: false})
	PowerAutoLogin   = base.AppendPower(&base.PowerAction{Action: "auto_login", Text: "自动登录", StandAlone: false})
	PowerUpload      = base.AppendPower(&base.PowerAction{Action: "upload", Text: "上传", StandAlone: true})
	PowerUpdateCheck = base.AppendPower(&base.PowerAction{Action: "update_check", Text: "更新检测", ShouldPower: true, ShouldLogin: true, StandAlone: true})
	PowerWebsocket   = base.AppendPower(&base.PowerAction{Action: "websocket", Text: "WebSocket", StandAlone: true})
)

func (this_ *Api) GetApis() (apis []*base.ApiWorker, err error) {
	apis = append(apis, &base.ApiWorker{Apis: []string{"data"}, Power: PowerData, Do: this_.apiData})
	apis = append(apis, &base.ApiWorker{Apis: []string{"login"}, Power: PowerLogin, Do: this_.apiLogin})
	apis = append(apis, &base.ApiWorker{Apis: []string{"autoLogin"}, Power: PowerAutoLogin, Do: this_.apiLogin})
	apis = append(apis, &base.ApiWorker{Apis: []string{"logout"}, Power: PowerLogout, Do: this_.apiLogout})
	apis = append(apis, &base.ApiWorker{Apis: []string{"register"}, Power: PowerRegister, Do: this_.apiRegister})
	apis = append(apis, &base.ApiWorker{Apis: []string{"session"}, Power: PowerSession, Do: this_.apiSession})
	apis = append(apis, &base.ApiWorker{Apis: []string{"upload"}, Power: PowerUpload, Do: this_.apiUpload, IsUpload: true})
	apis = append(apis, &base.ApiWorker{Apis: []string{"updateCheck"}, Power: PowerUpdateCheck, Do: this_.apiUpdateCheck})
	apis = append(apis, &base.ApiWorker{Apis: []string{"websocket"}, Power: PowerWebsocket, Do: this_.apiWebsocket, IsWebSocket: true})

	apis = append(apis, module_toolbox.NewToolboxApi(this_.toolboxService).GetApis()...)
	apis = append(apis, module_node.NewNodeApi(this_.nodeService).GetApis()...)
	apis = append(apis, module_file_manager.NewApi(this_.toolboxService, this_.nodeService).GetApis()...)
	apis = append(apis, module_terminal.NewApi(this_.toolboxService, this_.nodeService).GetApis()...)
	apis = append(apis, module_user.NewApi(this_.userService).GetApis()...)

	return
}

func (this_ *Api) appendApi(apis ...*base.ApiWorker) (err error) {
	if len(apis) == 0 {
		return
	}
	for _, api := range apis {
		if api.Power == nil {
			err = errors.New(fmt.Sprint("API未设置权限!", api))
			return
		}
		if len(api.Apis) == 0 {
			err = errors.New(fmt.Sprint("API未设置映射路径!", api))
			return
		}

		if !this_.IsServer {
			if !api.Power.StandAlone {
				continue
			}
		}
		for _, apiName := range api.Apis {

			_, find := this_.apiCache[apiName]
			if find {
				err = errors.New(fmt.Sprint("API映射路径[", apiName, "]已存在!", api))
				return
			}
			// println("add api path :" + apiName + ",action:" + api.Power.Action)
			this_.apiCache[apiName] = api
		}
	}
	return
}

func (this_ *Api) getRequestBean(c *gin.Context) (request *base.RequestBean) {
	request = &base.RequestBean{}
	request.JWT = this_.getJWT(c)
	return
}

func (this_ *Api) DoApi(path string, c *gin.Context) bool {

	index := strings.LastIndex(path, "api/")
	if index < 0 {
		return false
	}
	action := path[index+len("api/"):]

	api := this_.apiCache[action]
	if api == nil {
		return false
	}
	if api.IsGet && !strings.EqualFold(c.Request.Method, "get") {
		return false
	}
	if api.IsWebSocket && !strings.EqualFold(c.Request.Method, "get") {
		return false
	}
	requestBean := this_.getRequestBean(c)
	requestBean.Path = path
	if !this_.checkPower(api, requestBean.JWT, c) {
		return true
	}
	if api.Do != nil {
		var err error
		var startTime = util.Now()
		userAgentStr := c.Request.UserAgent()
		logRecode := &module_log.LogModel{
			Action:     action,
			Method:     c.Request.Method,
			StartTime:  startTime,
			CreateTime: startTime,
			Ip:         c.ClientIP(),
			UserAgent:  userAgentStr,
		}
		var param = make(map[string]interface{})
		_ = c.Request.ParseForm()
		f := c.Request.Form
		for k, v := range f {
			param[k] = v
		}
		f = c.Request.PostForm
		for k, v := range f {
			param[k] = v
		}
		if len(param) > 0 {
			bs, _ := json.Marshal(param)
			logRecode.Param = string(bs)
		}
		if !api.IsUpload {
			var data = make(map[string]interface{})
			_ = c.ShouldBindBodyWith(&data, binding.JSON)
			if len(data) > 0 {
				bs, _ := json.Marshal(data)
				logRecode.Data = string(bs)
			}
		}
		if requestBean.JWT != nil {
			logRecode.UserId = requestBean.JWT.UserId
		}
		_ = this_.logService.Start(logRecode)

		defer func() {

			_ = this_.logService.End(logRecode.LogId, startTime, err)
		}()

		this_.Logger.Info("处理操作", zap.String("action", action))
		res, err := api.Do(requestBean, c)
		if err != nil {
			this_.Logger.Error("操作异常", zap.String("action", action), zap.Any("error", err.Error()))
		}
		if res == base.HttpNotResponse {
			return true
		}
		base.ResponseJSON(res, err, c)
	}
	return true
}
