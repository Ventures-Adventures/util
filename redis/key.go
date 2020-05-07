package redis


const(
	KeyFormat              = "abctime:frequency:%d:%s"
	LockKeyFormat          = "abctime:frequency:lock:%d:%s"

	CaptchaFormat          = "abctime:captcha:login:%s"

	LogUserInfoHashMap     = "abctime:logsystem:hashmap:userinfo_"
	LogTopics   		   = "abctime:logsystem:topics"
)