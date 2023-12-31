package validator

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/url"
	"regexp"
	"strings"
)

type SubQuery struct {
	Sub           string               `form:"sub" binding:""`
	Subs          []string             `form:"-" binding:""`
	Proxy         string               `form:"proxy" binding:""`
	Proxies       []string             `form:"-" binding:""`
	Refresh       bool                 `form:"refresh,default=false" binding:""`
	Template      string               `form:"template" binding:""`
	RuleProvider  string               `form:"ruleProvider" binding:""`
	RuleProviders []RuleProviderStruct `form:"-" binding:""`
	Rule          string               `form:"rule" binding:""`
	Rules         []RuleStruct         `form:"-" binding:""`
	AutoTest      bool                 `form:"autoTest,default=false" binding:""`
	Lazy          bool                 `form:"lazy,default=false" binding:""`
}

type RuleProviderStruct struct {
	Behavior string
	Url      string
	Group    string
	Prepend  bool
}

type RuleStruct struct {
	Rule    string
	Prepend bool
}

func ParseQuery(c *gin.Context) (SubQuery, error) {
	var query SubQuery
	if err := c.ShouldBind(&query); err != nil {
		return SubQuery{}, errors.New("参数错误: " + err.Error())
	}
	if query.Sub == "" && query.Proxy == "" {
		return SubQuery{}, errors.New("参数错误: sub 和 proxy 不能同时为空")
	}
	if query.Sub != "" {
		query.Subs = strings.Split(query.Sub, ",")
		for i := range query.Subs {
			query.Subs[i], _ = url.QueryUnescape(query.Subs[i])
			if _, err := url.ParseRequestURI(query.Subs[i]); err != nil {
				return SubQuery{}, errors.New("参数错误: " + err.Error())
			}
		}
	} else {
		query.Subs = nil
	}
	if query.Proxy != "" {
		query.Proxies = strings.Split(query.Proxy, ",")
		for i := range query.Proxies {
			query.Proxies[i], _ = url.QueryUnescape(query.Proxies[i])
			if _, err := url.ParseRequestURI(query.Proxies[i]); err != nil {
				return SubQuery{}, errors.New("参数错误: " + err.Error())
			}
		}
	} else {
		query.Proxies = nil
	}
	if query.Template != "" {
		unescape, err := url.QueryUnescape(query.Template)
		if err != nil {
			return SubQuery{}, errors.New("参数错误: " + err.Error())
		}
		uri, err := url.ParseRequestURI(unescape)
		query.Template = uri.String()
		if err != nil {
			return SubQuery{}, errors.New("参数错误: " + err.Error())
		}
	}
	if query.RuleProvider != "" {
		var err error
		query.RuleProvider, err = url.QueryUnescape(query.RuleProvider)
		if err != nil {
			return SubQuery{}, errors.New("参数错误: " + err.Error())
		}
		reg := regexp.MustCompile(`\[(.*?)\]`)
		ruleProviders := reg.FindAllStringSubmatch(query.RuleProvider, -1)
		for i := range ruleProviders {
			length := len(ruleProviders)
			parts := strings.Split(ruleProviders[length-i-1][1], ",")
			if len(parts) != 4 {
				return SubQuery{}, errors.New("参数错误: ruleProvider 格式错误")
			}
			u := parts[1]
			u, err = url.QueryUnescape(u)
			if err != nil {
				return SubQuery{}, errors.New("参数错误: " + err.Error())
			}
			uri, err := url.ParseRequestURI(u)
			u = uri.String()
			if err != nil {
				return SubQuery{}, errors.New("参数错误: " + err.Error())
			}
			query.RuleProviders = append(
				query.RuleProviders, RuleProviderStruct{
					Behavior: parts[0],
					Url:      u,
					Group:    parts[2],
					Prepend:  parts[3] == "true",
				},
			)
		}
	} else {
		query.RuleProviders = nil
	}
	if query.Rule != "" {
		reg := regexp.MustCompile(`\[(.*?)\]`)
		rules := reg.FindAllStringSubmatch(query.Rule, -1)
		for i := range rules {
			length := len(rules)
			r := rules[length-1-i][1]
			strings.LastIndex(r, ",")
			parts := [2]string{}
			parts[0] = r[:strings.LastIndex(r, ",")]
			parts[1] = r[strings.LastIndex(r, ",")+1:]
			query.Rules = append(
				query.Rules, RuleStruct{
					Rule:    parts[0],
					Prepend: parts[1] == "true",
				},
			)
		}
	} else {
		query.Rules = nil
	}
	return query, nil
}
