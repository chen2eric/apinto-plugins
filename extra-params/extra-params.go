package extra_params

import (
	"encoding/json"
	"fmt"
	"github.com/eolinker/eosc"
	http_service "github.com/eolinker/eosc/http-service"
	"strconv"
	"strings"
)

var _ http_service.IFilter = (*ExtraParams)(nil)

type ExtraParams struct {
	*Driver
	id        string
	name      string
	params    []*ExtraParam
	errorType string
}

func (e *ExtraParams) DoFilter(ctx http_service.IHttpContext, next http_service.IChain) error {
	statusCode, err := e.access(ctx)
	if err != nil {
		ctx.Response().SetBody([]byte(err.Error()))
		ctx.Response().SetStatus(statusCode, strconv.Itoa(statusCode))
		return err
	}

	if next != nil {
		return next.DoChain(ctx)
	}
	return nil
}

func (e *ExtraParams) access(ctx http_service.IHttpContext) (int, error) {
	// 判断请求携带的content-type
	contentType := ctx.Proxy().Header().GetHeader("Content-Type")

	body, _ := ctx.Proxy().Body().RawBody()
	bodyParams, formParams, err := parseBodyParams(ctx, body, contentType)
	if err != nil {
		errInfo := fmt.Sprintf(parseBodyErrInfo, err.Error())
		err = encodeErr(e.errorType, errInfo, serverErrStatusCode)
		return serverErrStatusCode, err
	}

	headers := ctx.Proxy().Header().Headers()
	// 先判断参数类型
	for _, param := range e.params {
		switch param.Position {
		case "query":
			{
				value, err := getQueryValue(ctx, param)
				if err != nil {
					err = encodeErr(e.errorType, err.Error(), clientErrStatusCode)
					return clientErrStatusCode, err
				}
				ctx.Proxy().URI().SetQuery(param.Name, value)
			}
		case "header":
			{
				value, err := getHeaderValue(headers, param)
				if err != nil {
					err = encodeErr(e.errorType, err.Error(), clientErrStatusCode)
					return clientErrStatusCode, err
				}
				ctx.Proxy().Header().SetHeader(param.Name, value)
			}
		case "body":
			{
				value, err := getBodyValue(bodyParams, formParams, param, contentType)
				if err != nil {
					err = encodeErr(e.errorType, err.Error(), clientErrStatusCode)
					return clientErrStatusCode, err
				}
				if strings.Contains(contentType, FormParamType) {
					err = ctx.Proxy().Body().SetToForm(param.Name, value.(string))
					if err != nil {
						err = encodeErr(e.errorType, err.Error(), clientErrStatusCode)
						return clientErrStatusCode, err
					}
				} else if strings.Contains(contentType, JsonType) {
					bodyParams[param.Name] = value
				}
			}
		}
	}
	if strings.Contains(contentType, JsonType) {
		b, _ := json.Marshal(bodyParams)
		ctx.Proxy().Body().SetRaw(contentType, b)
	}

	return successStatusCode, nil
}

func (e *ExtraParams) Id() string {
	return e.id
}

func (e *ExtraParams) Start() error {
	return nil
}

func (e *ExtraParams) Reset(conf interface{}, workers map[eosc.RequireId]interface{}) error {
	confObj, err := e.check(conf)
	if err != nil {
		return err
	}

	e.params = confObj.Params
	e.errorType = confObj.ErrorType

	return nil
}

func (e *ExtraParams) Stop() error {
	return nil
}

func (e *ExtraParams) Destroy() {
	e.params = nil
	e.errorType = ""
}

func (e *ExtraParams) CheckSkill(skill string) bool {
	return http_service.FilterSkillName == skill
}
