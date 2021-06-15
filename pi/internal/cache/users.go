package cache

import (
	"VGO/pi/internal/orm"
)

var UsersCache *Users

type User struct {
	*orm.Users
	Calls []orm.CallName `form:"call" xml:"call" json:"call"`
}

type Users struct {
	List []User
}

func (s *Users) Default() User {
	if s == nil || len(s.List) == 0 {
		return s.New(false).List[0]
	}
	return s.List[0]
}
func (s *Users) New(fromSql bool) *Users {
	if s == nil || fromSql {
		var list = make([]User, 0)
		var userCall = make(map[int][]orm.CallName)
		callList := orm.CallName{}.All(&orm.CallName{})
		for _, c := range callList {
			userCall[c.Uid] = append(userCall[c.Uid], c)
		}
		userList := orm.Users{}.All(&orm.Users{})
		for _, v := range userList {
			if userCall[v.Uid] == nil {
				userCall[v.Uid] = []orm.CallName{}
			}
			vStr := v
			list = append(list, User{&vStr, userCall[v.Uid]})
		}
		UsersCache = &Users{List: list}
	}
	return UsersCache
}

func (s *Users) Update() *Users {
	s.Clear()
	return s.New(false)
}

func (s *Users) Clear() {
	UsersCache = nil
}
