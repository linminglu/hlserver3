package main

import (
	"database/sql"
	"github.com/oschwald/geoip2-golang"
	"github.com/todostreaming/gohw"
	"golang.org/x/sync/syncmap"
	"log"
	"sync"
	"time"
)

var (
	ident    int64      // identifier for every streaming session openned by an individual player
	mu_ident sync.Mutex // exclusive mutex for the identifier
	numgo    int        //number of goroutines working
	Hardw    *gohw.GoHw
	// DB Live vars
	dblive    *sql.DB    // db only with live players raw info
	mu_dblive sync.Mutex // also exclusive mutex for
	// DB mutexes
	dbday_mu sync.RWMutex
	dbmon_mu sync.RWMutex
	dbgen_mu sync.RWMutex
	// GeoIP2 vars
	dbgeoip    *geoip2.Reader
	mu_dbgeoip sync.Mutex
	// error loggers
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	// bandwidths map of the rawstreams (encoders) in bps
	// dont forget to empty it at the end of the day after dayly resume
	Bw_int *syncmap.Map
	// referer map ( ["rawstream"] = "domain1.com;domain2.com" )
	Referer *syncmap.Map
	// forecasters map ( ["near_proxy=rawstream"] = UNIXtimestamp_int64 )
	Forecaster *syncmap.Map
	// internal session maps (id, username, timestamp, type of user)
	id_     map[string]int       = make(map[string]int)
	user_   map[string]string    = make(map[string]string)
	time_   map[string]time.Time = make(map[string]time.Time)
	type_   map[string]int       = make(map[string]int)
	mu_user sync.RWMutex
	// user agents for OS's
	userAgent = map[string]string{"win": "Windows", "mac": "Mac OS X", "and": "Android", "lin": "Linux"}
	// settings.reg file
	cloud    map[string]string = make(map[string]string)
	mu_cloud sync.RWMutex
	// year months
	YearMonths = []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
)
