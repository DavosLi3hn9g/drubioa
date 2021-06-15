package aichat

import (
	"VGO/pi/internal/cache"
	"VGO/pi/internal/orm"
	"VGO/pi/internal/pkg/file"
	"VGO/pkg/fun"
	"regexp"
	"strconv"
	"strings"
)

var logIO *file.Log

type IQiar struct {
	User            *orm.Users
	Policies        *orm.Policies
	Scores          map[int]int
	cacheUsers      *cache.Users
	cacheIntentions *cache.Intentions
	answerKey       map[int]int
	IsEND           bool
}

func (iq *IQiar) Out(outChat, inputReply string) string {
	var out string
	if inputReply == "turing" {
		chat := Turing{}.Chat(outChat)
		if chat.Intent.Code == 10004 {
			out = chat.Results[0].Values["text"]
		} else {
			logIO.Error("turing api故障！")
		}
	}
	if inputReply == "custom" || inputReply == "" {
		chat := iq.SortOut(outChat)
		if chat != "" {
			out = chat
		}
	}
	return out
}

func (iq *IQiar) SortOut(in string) string {
	var chat string

	if iq.cacheUsers == nil {
		iq.cacheUsers = cache.UsersCache.New(false)
	}
	if iq.cacheIntentions == nil {
		iq.cacheIntentions = cache.IntentionsCache.New(false)
	}
	if iq.Policies == nil {
		iq.Policies = &orm.Policies{
			Checked: "sms", //默认只发送短信给主人
		}
	}
	if iq.Scores == nil {
		iq.Scores = make(map[int]int)
	}
	for k, user := range iq.cacheUsers.List {
		if k == 0 { //取第一个为默认用户
			iq.User = user.Users
		}
		for _, call := range user.Calls {
			r, _ := regexp.Compile(call.Call)
			if r.FindString(in) != "" {
				user.Users.IncHit(1)
				iq.User = user.Users
			}
		}
	}
	var ScoresHello int
	var query string
	for _, iv := range iq.cacheIntentions.List {
		for _, qv := range iv.Queries {
			if qv.Mode == 1 {
				query = "(" + strings.Replace(qv.Query, ",", "|", -1) + ")"
			} else {
				query = qv.Query
			}
			r, _ := regexp.Compile(query)
			if r.FindString(in) != "" {
				if iv.Hello {
					ScoresHello = qv.Scores
				} else {
					iq.Scores[iv.Sid] = iq.Scores[iv.Sid] + qv.Scores + ScoresHello
				}
				if iq.Scores[iv.Sid] >= 100 {
					iq.Scores[iv.Sid] = 0
					if iv.End != "" && fun.IsNumeric(iv.End) {
						policyId, _ := strconv.Atoi(iv.End)
						ormPolicies := new(orm.Policies)
						policy := ormPolicies.Get(policyId)
						policy.IncHit(1)
						iq.Policies = policy
						if !policy.Silent && policy.Title != "" {
							iq.IsEND = true
							return policy.Title
						} else {
							return ""
						}
					} else if iv.End != "" {
						return iv.End
					}
				}
				if qv.Answer != "" {
					if iq.answerKey == nil {
						iq.answerKey = make(map[int]int)
					}
					answerSlice := strings.Split(qv.Answer, ",")
					answer := strings.Join(answerSlice[iq.answerKey[qv.Id]:iq.answerKey[qv.Id]+1], "")
					if iq.answerKey[qv.Id] >= len(answerSlice)-1 {
						iq.answerKey[qv.Id] = 0
					} else {
						iq.answerKey[qv.Id]++
					}
					return answer
				}
			}
		}
	}
	return chat
}
