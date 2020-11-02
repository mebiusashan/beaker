package beaker

const HOST string = "http://localhost:9092"

const SUCC int = 99
const FAIL int = 1
const CONFIG_NOT_FOUND int = 2

const TAG_ARCHIVE string = "arc_"
const TAG_PAGE string = "pag_"
const TAG_HOME string = "home_"
const TAG_CAT string = "cat_"
const TAG_NOTFOUND string = "404_"

//数据库最大空闲数
const def_Database_MAX_IDLE_NUM = 20

//数据库最大连接数
const def_Database_MAX_OPEN_NUM = 100

//默认Redis过期时间
const def_Redis_EXPIRE_TIME = 3600

//默认网站名称
const def_Website_SITE_NAME = "Beaker"

//默认网站描述
const def_Website_SITE_DES = "Beaker is a simple blog system."

//默认footer内容
const def_Website_SITE_FOOTER = "自豪的采用Golang"

//默认首页显示记录数
const def_Website_INDEX_LIST_NUM = 10

//默认tweet每页显示数
const def_Website_TWEET_NUM_ONE_PAGE = 10

const PAGE string = "post.jet"

const ARCHIVE string = "post.jet"

const HOME string = "home.jet"

const CAT string = "home.jet"

const TWEET string = "tweet.jet"

const NotFound string = "404.jet"
