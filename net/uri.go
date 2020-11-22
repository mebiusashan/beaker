package net

//------------- cli -------------
//user
const CLI_CHECK string = "/user/check"
const CLI_LOGIN string = "/user/login"
const CLI_PING string = "/user/ping"

//article
const CLI_ART_LIST string = "/admin/art/list"
const CLI_ART_ADD string = "/admin/art/add"
const CLI_ART_RM string = "/admin/art/del"
const CLI_ART_DOWN string = "/admin/art/down"
const CLI_ART_MODIFY string = "/admin/art/modify"

//page
const CLI_PAGE_LIST string = "/admin/pag/list"
const CLI_PAGE_ADD string = "/admin/pag/add"
const CLI_PAGE_RM string = "/admin/pag/del"
const CLI_PAGE_DOWN string = "/admin/pag/down"
const CLI_PAGE_MODIFY string = "/admin/pag/modify"

//category
const CLI_CAT_LIST string = "/admin/cat/list"
const CLI_CAT_ADD string = "/admin/cat/add"
const CLI_CAT_RM string = "/admin/cat/del"
const CLI_CAT_MODIFY string = "/admin/cat/change"

//tweet
const CLI_TWEET_LIST string = "/admin/twe/list"
const CLI_TWEET_ADD string = "/admin/twe/add"
const CLI_TWEET_RM string = "/admin/twe/del"

//other
const CLI_CLEAN string = "/admin/clr/cache"
const CLI_OPTION string = "/admin/opt"

//------------- admin -------------

//user
const ADMIN_GROUP_USER string = "user"
const ADMIN_PING string = "ping"
const ADMIN_LOGIN string = "login"
const ADMIN_CHECK string = "check"

const ADMIN_GROUP_ADMIN string = "admin"

//article
const ADMIN_ART_ADD string = "art/add"
const ADMIN_ART_RM string = "art/del"
const ADMIN_ART_LIST string = "art/list"
const ADMIN_ART_DOWNLOAD string = "art/down"
const ADMIN_ART_MODIFY string = "art/modify"

//page
const ADMIN_PAGE_ADD string = "pag/add"
const ADMIN_PAGE_RM string = "pag/del"
const ADMIN_PAGE_LIST string = "pag/list"
const ADMIN_PAGE_DOWNLOAD string = "pag/down"
const ADMIN_PAGE_MODIFY string = "/pag/modify"

//tweet
const ADMIN_TWEET_ADD string = "twe/add"
const ADMIN_TWEET_RM string = "/twe/del"
const ADMIN_TWEET_LIST string = "/twe/list"

//category
const ADMIN_CAT_ADD string = "cat/add"
const ADMIN_CAT_RM string = "cat/del"
const ADMIN_CAT_LIST string = "cat/list"
const ADMIN_CAT_MODIFY string = "cat/change"

//other
const ADMIN_OPT string = "opt"
const ADMIN_CLEAN string = "clr/cache"

//------------- server -------------
const SERVER_INDEX string = "/"
const SERVER_ART string = "/archives/:id"
const SERVER_PAGE string = "/page/:id"
const SERVER_CAT string = "/category/:alias"
const SERVER_TWEET string = "/tweet/*page"
