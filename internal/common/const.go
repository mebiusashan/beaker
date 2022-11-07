package common

const SUCC int = 99
const FAIL int = 1
const CONFIG_NOT_FOUND int = 2

const TAG_ARCHIVE string = "arc_"
const TAG_PAGE string = "pag_"
const TAG_HOME string = "home_"
const TAG_CAT string = "cat_"
const TAG_NOTFOUND string = "404_"

//数据库最大空闲数
const Def_Database_MAX_IDLE_NUM = 20

//数据库最大连接数
const Def_Database_MAX_OPEN_NUM = 100

//默认Redis过期时间
const Def_Redis_EXPIRE_TIME = 3600

//默认网站名称
const Def_Website_SITE_NAME = "Beaker"

//默认网站描述
const Def_Website_SITE_DES = "Beaker is a simple blog system."

//默认footer内容
const Def_Website_SITE_FOOTER = "Beaker is a simple blog system."

//默认首页显示记录数
const Def_Website_INDEX_LIST_NUM = 10

//默认tweet每页显示数
const Def_Website_TWEET_NUM_ONE_PAGE = 10

const DefaultEditor string = "nano"

const SERVER_ENV_KEY string = "BEAKERPATH"
const ADMIN_ENV_KEY string = "BEAKERADMINPATH"

const SERVER_PUBLIC_KEY = "pub.pem"
const SERVER_PRIVATE_KEY = "pri.pem"
