package redis


const(
	KeyFormat              = ":frequency:%d:%s"
	LockKeyFormat          = ":frequency:lock:%d:%s"

	CaptchaFormat          = ":captcha:login:%s"

	LogUserInfoHashMap     = ":logsystem:hashmap:userinfo_"
	LogTopics   		   = ":logsystem:topics"
)
