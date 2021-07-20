package core

import (
	"VGO/pi/internal/orm"
	"strings"
	"time"
)

var LogCall = new(Call)

type Call struct {
	TimeStart time.Time
	TimeEnd   time.Time
	Intention []string
	Policy    string
	CallID    int //数据库ID
}

func (c *Call) Start(telfrom string) {
	c.TimeStart = time.Now()
	call := orm.LogsCall{}.Add(&orm.LogsCall{TelFrom: telfrom, TimeStart: int(c.TimeStart.Unix())})
	c.CallID = call.Id
}
func (c *Call) AddIntention(intention string) {
	c.Intention = append(c.Intention, intention)
	intentionStr := strings.Join(c.Intention, "|")
	orm.LogsCall{}.Updates(orm.LogsCall{Id: c.CallID, Intention: intentionStr})
}

func (c *Call) End(text, content, recording string) {
	if c.CallID > 0 {
		c.TimeEnd = time.Now()
		minute := int((c.TimeEnd.Sub(c.TimeStart).Seconds() + 30) / 60)
		orm.LogsCall{}.Updates(orm.LogsCall{Id: c.CallID, Text: text, Content: content, Recording: recording, TimeEnd: int(c.TimeEnd.Unix()), Minute: minute})
		c.CallID = 0
	}
}
func (c *Call) EndUpdate(text, content, recording string) {
	if c.CallID > 0 {
		orm.LogsCall{}.Updates(orm.LogsCall{Id: c.CallID, Text: text, Content: content, Recording: recording})
		c.CallID = 0
	}
}
