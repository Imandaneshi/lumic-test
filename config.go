package main

type Config struct {
	debug bool
	mysqlUri string
	testMysqlUri string
	redisUri string
	testRedisUri string
	secret string
}

var config = Config{}