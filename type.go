package onebot

import "fmt"

type Role string

const (
	PermOwner  Role = "owner"  // 群主
	PermAdmin  Role = "admin"  // 管理员
	PermMember Role = "member" // 群成员
)

// Friend 好友
type Friend struct {
	UserId   int64  `json:"user_id"`  // QQ号
	Nickname string `json:"nickname"` // 昵称
	Remark   string `json:"remark"`   // 备注
}

func (f *Friend) String() string {
	return fmt.Sprintf("%s(%d)", f.Nickname, f.UserId)
}

// Group 群
type Group struct {
	Id         int64  `json:"id"`         // 群号
	Name       string `json:"name"`       // 群名称
	Permission Role   `json:"permission"` // Bot在群中的权限
}

func (g *Group) String() string {
	return fmt.Sprintf("%s(%d)", g.Name, g.Id)
}

// Member 群成员
type Member struct {
	UserId   int64  `json:"user_id"`  // QQ号
	Nickname string `json:"nickname"` // 昵称
	Card     string `json:"card"`     // 群名片／备注
	Sex      Sex    `json:"sex"`      // 性别
	Age      int32  `json:"age"`      // 年龄
	Area     string `json:"area"`     // 地区
	Level    string `json:"level"`    // 成员等级
	Role     Role   `json:"role"`     // 角色
}

func (m *Member) String() string {
	return fmt.Sprintf("%s(%d)", m.Nickname, m.UserId)
}

type AnonymousMember struct {
	Id   int64  `json:"id"`   // 匿名用户 ID
	Name string `json:"name"` // 匿名用户名称
	Flag string `json:"flag"` // 匿名用户 flag，在调用禁言 API 时需要传入
}

func (a *AnonymousMember) String() string {
	return fmt.Sprintf("%s(%d)", a.Name, a.Id)
}

type Sex string

const (
	SexUnknown Sex = "unknown" // 未知
	SexMale    Sex = "male"    // 男
	SexFemale  Sex = "female"  // 女
)

// Profile 用户资料
type Profile struct {
	UserId   int64  `json:"user_id"`  // QQ号
	Nickname string `json:"nickname"` // 昵称
	Sex      Sex    `json:"sex"`      // 性别
	Age      int32  `json:"age"`      // 年龄
}
