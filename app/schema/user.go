package schema

// UserCurrent 用户基本信息
type UserCurrent struct {
	UserGpaCurrent
	UnreadCount int           `json:"unreadCount,omitempty"` // 未读消息
	System      string        `json:"system,omitempty"`      // 前端给出, 不同system带来的access也是不同的
	Access      interface{}   `json:"access,omitempty"`      // 返回权限列表,注意,其返回的权限只是部分确认的权限,而不是全部权限,并且,返回的权限是跟当前系统相关的
	Role        UserGpaRole   `json:"role,omitempty"`        // 当前用户使用的角色
	Menus       []UserGpaMenu `json:"menus,omitempty"`       // 当前用户菜单
	CreateAt    uint64        `json:"createAt,omitempty"`    // 获取时间
}

// UserGpaCurrent 用户基本信息
type UserGpaCurrent struct {
	ID     int    `db:"id" json:"-"`                    // 用户id
	UID    string `db:"uid" json:"userid"`              // 用户uid
	Avatar string `db:"avatar" json:"avatar,omitempty"` // 头像
	Name   string `db:"name" json:"name"`               // 姓名
}

// UserGpaRole 用户角色信息
type UserGpaRole struct {
	ID   int    `db:"id" json:"-"`      // 角色id
	UID  string `db:"uid" json:"id"`    // 角色uid
	Name string `db:"name" json:"name"` // 角色name
}

// UserGpaMenu 用户菜单
type UserGpaMenu struct {
	ID       int           `db:"id" json:"-"`                   // 菜单ID
	Key      string        `db:"uid" json:"key"`                // 菜单KEY, 全局唯一标识符,在当前系统中, 层级结构,每层使用3个字符
	Locale   string        `db:"locale" json:"local,omitempty"` // 可以直接抽取前端i18n对应的内容,如果不配置,可以通过name抽取
	Name     string        `db:"name" json:"name"`              // 菜单内容, 必要字段
	Icon     string        `db:"icon" json:"icon1"`             // 兼容kratos系统,json字段为icon1
	Path     string        `db:"router" json:"path"`            // 访问地址
	Children []UserGpaMenu `json:"children,omitempty"`          // 子菜单
	Parent   *UserGpaMenu  `json:"-"`                           // 父菜单
}
