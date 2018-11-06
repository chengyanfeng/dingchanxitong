package util

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"egg_backend/redis"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"os"
	"sort"
	"strconv"
	"time"

	//	"github.com/gin-gonic/gin"
	//	"qixijie/util"

	//	"Egg/util"
	"strings"

	"github.com/muesli/cache2go"
	"gopkg.in/mgo.v2/bson"
)

type P map[string]interface{}

func Md5(s ...interface{}) (r string) {
	return Hash("md5", s...)
}
func ToFloat(s interface{}, default_v ...float64) float64 {
	f64, e := strconv.ParseFloat(ToString(s), 64)
	if e != nil && len(default_v) > 0 {
		return default_v[0]
	}
	return f64
}
func Hash(algorithm string, s ...interface{}) (r string) {
	var h hash.Hash
	switch algorithm {
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	case "sha2", "sha256":
		h = sha256.New()
	}
	for _, value := range s {
		switch value.(type) {
		case []byte:
			h.Write(value.([]byte))
		default:
			h.Write([]byte(ToString(value)))
		}
	}
	r = hex.EncodeToString(h.Sum(nil))
	return
}
func ToString(v interface{}) string {
	if v != nil {
		switch v.(type) {
		case bson.ObjectId:
			return v.(bson.ObjectId).Hex()
		case []byte:
			return string(v.([]byte))
		case *P, P:
			var p P
			switch v.(type) {
			case *P:
				if v.(*P) != nil {
					p = *v.(*P)
				}
			case P:
				p = v.(P)
			}
			var keys []string
			for k := range p {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			r := "P{"
			for _, k := range keys {
				r = JoinStr(r, k, ":", p[k], " ")
			}
			r = JoinStr(r, "}")
			return r
		case int64:
			return strconv.FormatInt(v.(int64), 10)
		default:
			return fmt.Sprintf("%v", v)
		}
	}
	return ""
}
func JoinStr(val ...interface{}) (r string) {
	for _, v := range val {
		r += ToString(v)
	}
	return
}

//string 转P
func JsonDecode(b []byte) (p *map[string]interface{}) {
	p = &map[string]interface{}{}
	err := json.Unmarshal(b, p)
	if err != nil {
		fmt.Print(err)
	}
	return
}

func IsEmpty(v interface{}) bool {
	if v == nil {
		return true
	}
	switch v.(type) {
	case P:
		return len(v.(P)) == 0
	}
	return ToString(v) == ""
}
func (p *P) ToInt(s ...string) {
	for _, k := range s {
		v := ToString((*p)[k])
		(*p)[k] = ToInt(v)
	}
}
func ToInt(s interface{}, default_v ...int) int {
	i, e := strconv.Atoi(ToString(s))
	if e != nil && len(default_v) > 0 {
		return default_v[0]
	}
	return i
}

func StringToListString(OldHenIdList, addOrDeleHenId string, flag bool) (newHenId string) {
	if flag == true {
		if len(OldHenIdList) == 0 {
			newHenId = addOrDeleHenId
			return
		} else {
			newHenId = OldHenIdList + "," + addOrDeleHenId
		}
		return newHenId
	} else {
		if len(OldHenIdList) == 0 {
			return OldHenIdList
		} else {
			listString := strings.Split(OldHenIdList, ",")
			listDele := []string{}
			for _, v := range listString {
				if ToString(v) != addOrDeleHenId {
					listDele = append(listDele, v)
				}
			}
			newHenId = strings.Join(listDele, ",")
			return
		}
	}

}

func StringToStringList(stringList string) []string {
	listString := []string{}
	listString = strings.Split(stringList, ",")
	return listString
}
func WriteFile(url string, body []byte) bool {
	f, err := os.Create(url)
	defer f.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else {
		_, err = f.Write(body)
		if err != nil {
			return false
		} else {
			return true
		}
	}
}

/****************************------------------以下方法为缓存-------------------------************************************/
type Cacha struct {
	value    string
	moreData []byte
}

//添加缓存，这是go 的自带缓存，好处不用安装redis ,坏处，每次重启都必须重新登陆
func AddCacheGoSelf(token, openId string) bool {
	//创建缓存表,有则忽略，无则创建
	cache := cache2go.Cache("Cache")

	val := Cacha{openId, []byte{}}
	cache.Add(token, 120*time.Minute, &val)

	// 验证是否存在
	res, err := cache.Value(token)
	if err == nil {
		fmt.Print(token)
		fmt.Print(res.Data().(*Cacha).value)
		return true
	} else {
		return false
	}
}

//添加缓存,redis
func AddCache(token, openId string) bool {
	//查看能否为11位的手机号码
	if len(token) < 14 {
		//像redis，添加缓存，时间为120秒
		redis.ReAdd(token, openId, 120)
	} else {
		//像redis，添加缓存，时间为7200秒
		redis.ReAdd(token, openId, 7200)
	}

	//查询是存储成功
	flag := redis.ReIsEx(token)
	return flag
}

//获取缓存,这是go 的自带缓存，好处不用安装redis ,坏处，每次重启都必须重新登陆,暂时无用
func GetCacheGoSelf(token string) string {
	//创建缓存表,有则获取Cache表，无则创建
	cache := cache2go.Cache("Cache")
	res, err := cache.Value(token)
	if err == nil {
		return res.Data().(*Cacha).value
	} else {
		return ""
	}
}

//获取value
func GetCache(token string) string {
	value, flag := redis.ReGet(token)
	if !flag {
		return "redis database is fail"
	} else {
		//重新设置key的时间,单位为小时
		redis.ReExpr(token, 2)
		return value
	}
}

//获取所有缓存,这是go 的自带缓存，好处不用安装redis ,坏处，每次重启都必须重新登陆，暂时无用
func GetAllCacheGoSelf() (listvalue []int) {
	//创建缓存表,有则获取Cache表，无则创建
	cache := cache2go.Cache("Cache")
	listvalue = []int{}
	//获取数据
	trans := func(key interface{}, item *cache2go.CacheItem) {
		value := item.Data().(*Cacha).value
		listvalue = append(listvalue, ToInt(value))
	}
	//获取所有的数据
	cache.Foreach(trans)
	return
}

//获取所有缓存
func GetAllCache() (listvalue []int) {
	stringList, flag := redis.ReGetAllKey("*")
	if !flag {
		return
	}
	for _, v := range stringList {
		value, err := redis.ReGet(v)
		if !err {
			listvalue = append(listvalue, ToInt(value))
		}
	}

	return
}

func GetCurDayTime() int {
	the_time, _ := time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Local)
	return ToInt(the_time.Unix())
}

func GetYesDayTime() int {
	the_time, _ := time.ParseInLocation("2006-01-02", time.Now().Format("2006-01-02"), time.Local)
	return ToInt(the_time.Unix()) - 86400
}

func DelMember(list []int, m int) (ret []int) {
	l := len(list)
	for i := 0; i < l; i++ {
		if list[i-1] != m {
			ret = append(ret, list[i])
		}
	}
	return
}
