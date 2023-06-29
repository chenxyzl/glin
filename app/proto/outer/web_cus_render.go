package outer

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"laiya/proto/code"
	"net/http"
)

var PbJson = protojson.MarshalOptions{EmitUnpopulated: true, UseEnumNumbers: true}

type CusRender struct {
	Code code.Code
	Data proto.Message
}

func (r *CusRender) writeContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

// Render 实现gin.Render接口的方法
func (r *CusRender) Render(w http.ResponseWriter) error {
	// 在这里实现自定义的渲染逻辑
	r.writeContentType(w)
	msg, err := PbJson.Marshal(r.Data)
	if err != nil {
		return err
	}

	v := gin.H{"code": r.Code, "msg": r.Code.GetCodeMsg(), "data": string(msg)}
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}

// WriteContentType 实现gin.Render接口的方法
func (r *CusRender) WriteContentType(w http.ResponseWriter) {
	r.writeContentType(w)
}
