package utils

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"reflect"
)

type ProblemDetail struct {
	Status     int    `json:"status,omitempty"`
	Title      string `json:"title,omitempty"`
	Detail     string `json:"detail,omitempty"`
	Type       string `json:"type,omitempty"`
	Instance   string `json:"instance,omitempty"`
	StackTrace string `json:"stackTrace,omitempty"`
}

var mappers = map[reflect.Type]func() ProblemDetailErr{}
var mapperStatus = map[int]func() ProblemDetailErr{}

// ProblemDetailErr ProblemDetail error interface
type ProblemDetailErr interface {
	SetStatus(status int) ProblemDetailErr
	GetStatus() int
	SetTitle(title string) ProblemDetailErr
	GetTitle() string
	SetDetail(detail string) ProblemDetailErr
	GetDetails() string
	SetType(typ string) ProblemDetailErr
	GetType() string
	SetInstance(instance string) ProblemDetailErr
	GetInstance() string
	SetStackTrace(stackTrace string) ProblemDetailErr
	GetStackTrace() string
}

func (p *ProblemDetail) SetDetail(detail string) ProblemDetailErr {
	p.Detail = detail

	return p
}

func (p *ProblemDetail) GetDetails() string {
	return p.Detail
}

func (p *ProblemDetail) SetStatus(status int) ProblemDetailErr {
	p.Status = status

	return p
}

func (p *ProblemDetail) GetStatus() int {
	return p.Status
}

func (p *ProblemDetail) SetTitle(title string) ProblemDetailErr {
	p.Title = title

	return p
}

func (p *ProblemDetail) GetTitle() string {
	return p.Title
}

func (p *ProblemDetail) SetType(typ string) ProblemDetailErr {
	p.Type = typ

	return p
}

func (p *ProblemDetail) GetType() string {
	return p.Type
}

func (p *ProblemDetail) SetInstance(instance string) ProblemDetailErr {
	p.Instance = instance

	return p
}

func (p *ProblemDetail) GetInstance() string {
	return p.Instance
}

func (p *ProblemDetail) SetStackTrace(stackTrace string) ProblemDetailErr {
	p.StackTrace = stackTrace

	return p
}

func (p *ProblemDetail) GetStackTrace() string {
	return p.StackTrace
}

func writeTo(w http.ResponseWriter, p ProblemDetailErr) (int, error) {

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(p.GetStatus())

	val, err := json.Marshal(p)
	if err != nil {
		return 0, err
	}
	return w.Write(val)
}

// ResolveProblemDetails retrieve and resolve error with format problem details error
func ResolveProblemDetails(w http.ResponseWriter, r *http.Request, err error) (ProblemDetailErr, error) {

	var statusCode int = http.StatusInternalServerError
	var echoError *echo.HTTPError

	if errors.As(err, &echoError) {
		var HTTPError *echo.HTTPError
		errors.As(err, &HTTPError)
		statusCode = HTTPError.Code
		err = err.(*echo.HTTPError).Message.(error)
	}

	var mapCustomType, mapCustomTypeErr = setMapCustomType(w, r, err)
	if mapCustomType != nil {
		return mapCustomType, mapCustomTypeErr
	}

	var mapStatus, mapStatusErr = setMapStatusCode(w, r, err, statusCode)
	if mapStatus != nil {
		return mapStatus, mapStatusErr
	}

	var p, errr = setDefaultProblemDetails(w, r, err, statusCode)
	if errr != nil {
		return nil, err
	}
	return p, errr
}

func setMapCustomType(w http.ResponseWriter, r *http.Request, err error) (ProblemDetailErr, error) {

	problemCustomType := mappers[reflect.TypeOf(err)]
	if problemCustomType != nil {
		prob := problemCustomType()

		validationProblems(prob, err, r)

		for k, v := range mapperStatus {
			if k == prob.GetStatus() {
				_, err = writeTo(w, v())
				if err != nil {
					return nil, err
				}
				return prob, err
			}
		}

		_, err = writeTo(w, prob)
		if err != nil {
			return nil, err
		}
		return prob, err
	}
	return nil, err
}

func setMapStatusCode(w http.ResponseWriter, r *http.Request, err error, statusCode int) (ProblemDetailErr, error) {
	problemStatus := mapperStatus[statusCode]
	if problemStatus != nil {
		prob := problemStatus()
		validationProblems(prob, err, r)
		_, err = writeTo(w, prob)
		if err != nil {
			return nil, err
		}
		return prob, err
	}
	return nil, err
}

func setDefaultProblemDetails(w http.ResponseWriter, r *http.Request, err error, statusCode int) (ProblemDetailErr, error) {
	defaultProblem := func() ProblemDetailErr {
		return &ProblemDetail{
			Type:       getDefaultType(statusCode),
			Status:     statusCode,
			Detail:     err.Error(),
			Title:      http.StatusText(statusCode),
			Instance:   r.URL.RequestURI(),
			StackTrace: errorsWithStack(err),
		}
	}
	prob := defaultProblem()
	_, err = writeTo(w, prob)
	if err != nil {
		return nil, err
	}
	return prob, err
}

func validationProblems(problem ProblemDetailErr, err error, r *http.Request) {
	problem.SetDetail(err.Error())

	if problem.GetStatus() == 0 {
		problem.SetStatus(http.StatusInternalServerError)
	}
	if problem.GetInstance() == "" {
		problem.SetInstance(r.URL.RequestURI())
	}
	if problem.GetType() == "" {
		problem.SetType(getDefaultType(problem.GetStatus()))
	}
	if problem.GetTitle() == "" {
		problem.SetTitle(http.StatusText(problem.GetStatus()))
	}
	if problem.GetStackTrace() == "" {
		problem.SetStackTrace(errorsWithStack(err))
	}
}

func getDefaultType(statusCode int) string {
	return fmt.Sprintf("https://httpstatuses.io/%d", statusCode)
}

func errorsWithStack(err error) string {
	res := fmt.Sprintf("%+v", err)
	return res
}
