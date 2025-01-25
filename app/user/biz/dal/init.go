package dal

import (
	"github.com/jazwu/tiktok-ecommerce/app/user/biz/dal/mysql"
	"github.com/jazwu/tiktok-ecommerce/app/user/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
