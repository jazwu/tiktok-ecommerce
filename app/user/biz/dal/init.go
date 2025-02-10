package dal

import (
	"github.com/jazwu/tiktok-ecommerce/app/user/biz/dal/mysql"
)

func Init() {
	// redis.Init()
	mysql.Init()
}
