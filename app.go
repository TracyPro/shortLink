package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"gopkg.in/validator-2"
	"log"
	"net/http"
)

// router and middleware

type App struct {
	Router      *mux.Router
	Middlewares *Middleware
	config      *Env
}

// request

type shortenReq struct {
	URL                 string `json:"url" validate:"nonzero"`                 // 验证非0
	ExpirationInMinutes int64  `json:"expiration_in_minutes" validate:"min=0"` // 过期时间 不能小于0
}

// response
type shortLinkResp struct {
	ShortLink string `json:"short_link"`
}

// init app struct
func (a *App) Initialize(e *Env) {
	// 定义log格式
	// 日志发生时间和日期 | 日志行号和文件名
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 初始化app中router
	a.Router = mux.NewRouter()

	// 初始化中间件
	a.Middlewares = &Middleware{}

	// 绑定router和handler之间关系
	a.initializeRouters()

	// 初始化config
	a.config = e
}

func (a *App) initializeRouters() {
	// 在路由匹配之前执行middleware
	// 使用alice包
	m := alice.New(a.Middlewares.LoggingHandler, a.Middlewares.RecoverHandler)
	// 获取短链接
	a.Router.Handle("/api/shorten", m.ThenFunc(a.createShortLink)).Methods("POST")
	// 短链接详细信息
	a.Router.Handle("/api/info", m.ThenFunc(a.getShortInfo)).Methods("GET")
	// 只能是字符或数字 位数1-11位
	a.Router.Handle("/{shortlink:[a-zA-Z0-9]{1,11}}", m.ThenFunc(a.redirect)).Methods("GET")
}

func (a *App) createShortLink(w http.ResponseWriter, r *http.Request) {
	var req shortenReq
	// 用户传来的请求是json格式 解析body中json信息 解析道req
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// json解析错误
		// 参数1：response，
		// 参数2：定义的结构体，实现自定义错误接口，错误字符串描述
		respondWithError(w, StatusError{http.StatusBadRequest,
			fmt.Errorf("Parse parameters failed %v", r.Body)})
		return
	}
	// 验证解析内容是否满足定义的非零值和最小值为0
	if err := validator.Validate(req); err != nil {
		respondWithError(w, StatusError{http.StatusBadRequest,
			fmt.Errorf("validate parameters failed %v", req)})
		return
	}

	defer r.Body.Close()
	s, err := a.config.S.Shorten(req.URL, req.ExpirationInMinutes)
	if err != nil {
		respondWithError(w, err)
	} else {
		respondWithJSON(w, http.StatusCreated, shortLinkResp{ShortLink: s})
	}
}

func (a *App) getShortInfo(w http.ResponseWriter, r *http.Request) {
	// 得到URL参数
	vals := r.URL.Query()
	s := vals.Get("shortlink")

	d, err := a.config.S.ShortLinkInfo(s)
	if err != nil {
		respondWithError(w, err)
	} else {
		respondWithJSON(w, http.StatusOK, d)
	}
	// 测试 recover middleware
	//panic(s)
}

func (a *App) redirect(w http.ResponseWriter, r *http.Request) {
	// 重定向函数 shortlink是从变量中取的 返回字典类型
	vars := mux.Vars(r)
	u, err := a.config.S.Unshort(vars["shortlink"])
	if err != nil {
		respondWithError(w, err)
	} else {
		http.Redirect(w, r, u, http.StatusTemporaryRedirect)
	}
}

// 启动监听和服务
// 参数：地址
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
