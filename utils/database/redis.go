// NewRedisPool creates a new pool and a default connection.
// It is applicable for multiple connection scenarios.
package database

import (
	"time"

	"github.com/TianQinS/commhttp/config"
	"github.com/gomodule/redigo/redis"
)

var (
	Conf = &config.Conf.Redis
)

type RDB struct {
	// db calls the Get method to get a connection from the pool and
	// the Close method to return the connection's resources to the pool.
	db *redis.Pool
}

// Check the health of an idle connection before the connection is used again by the application.
func testOnBorrow(con redis.Conn, t time.Time) error {
	if time.Since(t) < time.Minute {
		return nil
	}
	_, err := con.Do("PING")
	return err
}

// Notice that timeouts will not be set up when the readTimeout or writeTimeout parameter is zero.
func NewRedisPool(url, pwd string, index, readTimeout, writeTimeout int) (rdb *RDB, err error) {
	rdb = &RDB{
		db: &redis.Pool{
			MaxActive: Conf.MaxActive,
			MaxIdle:   Conf.MaxIdle,
			// IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				var con redis.Conn
				if con, err = redis.Dial("tcp", url, redis.DialPassword(pwd),
					redis.DialConnectTimeout(3*time.Second),
					redis.DialReadTimeout(time.Duration(readTimeout)*time.Millisecond),
					redis.DialWriteTimeout(time.Duration(writeTimeout)*time.Millisecond),
					redis.DialKeepAlive(30*time.Second),
				); err != nil {
					return nil, err
				}
				if _, err = con.Do("SELECT", index); err != nil {
					con.Close()
					return nil, err
				}
				return con, nil
			},
			TestOnBorrow: testOnBorrow,
		},
	}
	return
}

// Gets a connection. The application must close the returned connection.
// If there is an error, then the connection Err, Do, Send, Flush
// and Receive methods return that error.
func (this *RDB) Get() redis.Conn {
	return this.db.Get()
}

func (this *RDB) GetConn() (redis.Conn, error) {
	con := this.db.Get()
	err := con.Err()
	if err != nil {
		con.Close()
		return nil, err
	}
	return con, nil
}

func (this *RDB) GetKey(key string) (res string, err error) {
	con := this.db.Get()
	defer con.Close()
	res, err = redis.String(con.Do("GET", key))
	return
}

func (this *RDB) SetKey(key, value string) error {
	con := this.db.Get()
	defer con.Close()
	_, err := con.Do("SET", key, value)
	return err
}

func (this *RDB) SetEx(key, value string, timeout int) error {
	// Con.DoWithTimeout(time.Second, "SET", key, value, "EX", timeout)
	con := this.db.Get()
	defer con.Close()
	_, err := con.Do("SET", key, value, "EX", timeout)
	return err
}

// Close releases the resources used by the pool.
func (this *RDB) Close() {
	this.db.Close()
}
