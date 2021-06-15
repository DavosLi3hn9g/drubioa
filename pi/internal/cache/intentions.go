package cache

import (
	"VGO/pi/internal/orm"
	"VGO/pkg/fun"
	"strconv"
)

var IntentionsCache *Intentions

type Intention struct {
	*orm.Intentions
	Queries []*orm.Queries `form:"queries" xml:"queries" json:"queries"`
	Policy  *orm.Policies  `form:"policy" xml:"policy" json:"policy"`
}

type Intentions struct {
	List []Intention
}

func (s *Intentions) New(fromSql bool) *Intentions {
	if IntentionsCache == nil || fromSql {

		var list = make([]Intention, 0)
		var intention = make(map[int][]*orm.Queries)

		pList := orm.Queries{}.All(&orm.Queries{})
		for _, pv := range pList {
			intention[pv.Sid] = append(intention[pv.Sid], pv)
		}
		cList := orm.Intentions{}.All(&orm.Intentions{})
		for _, cv := range cList {
			if intention[cv.Sid] == nil {
				intention[cv.Sid] = []*orm.Queries{}
			}
			var policy = new(orm.Policies)
			if cv.End != "" && fun.IsNumeric(cv.End) {
				policyId, _ := strconv.Atoi(cv.End)
				policy = policy.Get(policyId)
			}

			vStr := cv
			list = append(list, Intention{vStr, intention[cv.Sid], policy})
		}

		IntentionsCache = &Intentions{List: list}
	}
	return IntentionsCache
}

func (s *Intentions) Update() *Intentions {
	s.Clear()
	return s.New(false)
}

func (s *Intentions) Clear() {
	IntentionsCache = nil
}
