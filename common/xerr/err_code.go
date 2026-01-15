package xerr

// 成功返回
const OK uint32 = 200

// 全局错误码
const SERVER_COMMON_ERROR uint32 = 100001
const REUQEST_PARAM_ERROR uint32 = 100002
const TOKEN_EXPIRE_ERROR uint32 = 100003
const DB_ERROR uint32 = 100004

// 业务错误码可以往下扩展
const USER_NOT_EXIST uint32 = 200001
