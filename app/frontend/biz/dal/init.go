package dal

import (
	"github.com/jazwu/tiktok-ecommerce/app/frontend/biz/dal/mysql"
	"github.com/jazwu/tiktok-ecommerce/app/frontend/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
