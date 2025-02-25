package dal

import (
	"github.com/jazwu/tiktok-ecommerce/app/cart/biz/dal/mysql"
	"github.com/jazwu/tiktok-ecommerce/app/cart/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
