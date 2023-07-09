package nutshttp

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xujiajun/nutsdb"
)

type NutsHTTPServer struct {
	core *core
	r    *gin.Engine
}

func NewNutsHTTPServer(db *nutsdb.DB) (*NutsHTTPServer, error) {
	c := &core{db}

	r := gin.Default()

	s := &NutsHTTPServer{
		core: c,
		r:    r,
	}

	err := s.InitConfig()
	if err != nil {
		return nil, err
	}

	s.initRouter()

	return s, nil
}

func (s *NutsHTTPServer) InitConfig() error {
	jwtSetting = JWTSetting{
		Secret: "nutsdb",
		Issuer: "nuts-http",
		Expire: 2880,
	}
	jwtSetting.Expire *= time.Minute
	return nil
}

func (s *NutsHTTPServer) Run(addr string) error {
	return http.ListenAndServe(addr, s.r)
}

func (s *NutsHTTPServer) initRouter() {

	s.r.Use(Cors(), gin.Recovery())

	s.initXanaduRouter()

	s.initSetRouter()

	s.initListRouter()

	s.initStringRouter()

	s.initZSetRouter()

	s.initLoginRouter()

	s.initCommonRouter()

}
