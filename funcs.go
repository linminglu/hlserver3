package main

import (
	"bufio"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// load referer table from general.db and mirror to Referer syncmap
// any error will make it exit immediately
func loadallreferers() {
	dbgeneral, err := sql.Open("sqlite3", DirDB+"general.db") // Apertura de la dateDayly.db antigua para lectura del pico/hora
	if err != nil {
		log.Fatalln("Fails openning general.db:", err)
	}
	defer dbgeneral.Close()
	dbgen_mu.RLock()
	query, err := dbgeneral.Query("SELECT username, streamname, referrers FROM referer")
	dbgen_mu.RUnlock()
	if err != nil {
		log.Fatalln("Fails querying general.db:", err)
		return
	}
	defer query.Close()
	for query.Next() {
		var user, stream, referer string
		err = query.Scan(&user, &stream, &referer)
		if err != nil {
			log.Fatalln("Fails scanning general.db:", err)
		}
		Referer.Store(user+"-"+stream, referer)
	}
}

// splits the IPv4/6 from the port used
func getip(pseudoip string) string {
	var res string
	if strings.Contains(pseudoip, "]:") {
		part := strings.Split(pseudoip, "]:")
		res = part[0]
		res = res[1:]
	} else {
		part := strings.Split(pseudoip, ":")
		res = part[0]
	}
	return res
}

// converts a string to a numerical integer
func toInt(cant string) (res int) {
	res, _ = strconv.Atoi(cant)
	return
}

func random(min, max int) int { // [min,max)
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}

func geoIP(ipaddr string) (country, isocode, city string) {
	// If you are using strings that may be invalid, check that ip is not nil
	ip := net.ParseIP(ipaddr)
	mu_dbgeoip.Lock()
	record, err := dbgeoip.City(ip)
	mu_dbgeoip.Unlock()
	if err != nil {
		return
	}
	city = record.City.Names["en"]
	country = record.Country.Names["en"]
	isocode = record.Country.IsoCode

	return country, isocode, city
}

// from a complete url
func getdomain(url string) string {
	var domain string

	p := strings.Split(url, "/")
	if len(p) > 2 {
		domain = p[2]
	}

	return domain
}

// get os from useragent
func getos(agent string) string {
	os := "other"

	for key, value := range userAgent {
		if strings.Contains(agent, value) {
			os = key
			break
		}
	}

	return os
}

func loadSettings(filename string) {
	fr, err := os.Open(filename)
	defer fr.Close()
	if err == nil {
		reader := bufio.NewReader(fr)
		for {
			linea, rerr := reader.ReadString('\n')
			if rerr != nil {
				break
			}
			linea = strings.TrimRight(linea, "\n")
			item := strings.Split(linea, " = ")
			mu_cloud.Lock()
			if len(item) == 2 {
				cloud[item[0]] = item[1]
			}
			mu_cloud.Unlock()
		}
	}
}

// clean old registers of more than 1 day
func clean(key, val interface{}) bool { // ["near_proxy=rawstream"] = UNIXtimestamp_int64
	var k string
	var v int

	k = key.(string)
	v = val.(int)
	limit_time := time.Now().Unix() - 86400
	if int64(v) < limit_time {
		Forecaster.Delete(k)
	}

	return true
}
