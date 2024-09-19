package network

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"time"
	"github.com/gsterjov/go-libsecret"
	"golang.org/x/crypto/pbkdf2"
	"crypto/sha1"
	"crypto/aes"
	"crypto/cipher"
	"gotube/youtube"
)

const CHROMIUM_COOKIES_LOCATION string = "/.config/chromium/Default/Cookies"

func get_keys() [][]byte {

	service, err := libsecret.NewService()
	if err != nil {
		panic(err)
	}
	
	session, err := service.Open()
	if err != nil {
		panic(err)
	}
	
	collections, err := service.Collections()
	if err != nil {
		panic(err)
	}
	
	items, err := collections[0].Items()
	if err != nil {
		panic(err)
	}
	
	if len(items) == 0 {
		items, err = collections[1].Items()
		if err != nil {
			panic(err)
		}
	}
	
	toReturn := [][]byte{}
	
	for _, item := range items {
		label, err := item.Label()
		if err != nil {
			panic(err)
		}
		
		if label == "Chromium Safe Storage" {
			secret, err := item.GetSecret(session)
			if err != nil {
				panic(err)
			}
			
			toReturn = append(toReturn, secret.Value)
		}
		
	}
	return toReturn
}

func clean(input []byte) []byte {
	var paddingLen int = int(input[len(input)-1])
	var lastRealByte int = len(input) - paddingLen
	return input[:lastRealByte]
}

func decryptCookie(password []byte, encryptedCookie []byte) string {
	encryptedCookie = encryptedCookie[3:]
	
	var salt []byte = []byte("saltysalt")
	var iv []byte = []byte("                ")
	var iterations int = 1
	var length int = 16
	
	key := pbkdf2.Key(password, salt, iterations, length, sha1.New)
	// https://medium.com/insiderengineering/aes-encryption-and-decryption-in-golang-php-and-both-with-full-codes-ceb598a34f41
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	
	mode := cipher.NewCBCDecrypter(block, iv)
	
	mode.CryptBlocks(encryptedCookie, encryptedCookie)
	return string(clean(encryptedCookie))
}

func GetCookiesChromium() *cookiejar.Jar {
	var keys [][]byte = get_keys()
	
	db, err := sql.Open("sqlite3", youtube.HOME_DIR + CHROMIUM_COOKIES_LOCATION)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//row, err := db.Query("SELECT name, value, host, path, expiry, isHttpOnly FROM moz_cookies where host= \".youtube.com\"")
	
	row, err := db.Query("SELECT host_key, name, encrypted_value, path, expires_utc, is_httponly FROM cookies WHERE host_key= \".youtube.com\"")
	
	//row, err := db.Query("SELECT name, value, host_key, path, expires_utc, is_httponly FROM Cookies where host_key= \".youtube.com\"")
	
	if err != nil {
		panic(err)
	}
	defer row.Close()
	
	var cookiesList []*http.Cookie

	for row.Next() {
		var domain, name, path, expires_utc, is_httponly string
		var encrypted_value []byte
		err := row.Scan(&domain, &name, &encrypted_value, &path, &expires_utc, &is_httponly)
		if err != nil {
			panic(err)
		}
		
		var httpOnlyBool bool = false
		if is_httponly == "1" {
			httpOnlyBool = true
		}
		
		expiresAt64, err := strconv.ParseInt(expires_utc, 10, 64)
		if err != nil {
			panic(err)
		}
		
		tm := time.Unix(expiresAt64, 0)

		if name == "__Secure-1PAPISID" || name == "__Secure-1PSID" || name == "__Secure-1PSIDCC" || name == "__Secure-1PSIDTS" || name == "__Secure-3PAPISID" || name == "__Secure-3PSID" || name == "__Secure-3PSIDCC" || name == "__Secure-3PSIDTS" || name == "APISID" || name == "CONSENT" || name == "CONSISTENCY" || name == "HSID" || name == "LOGIN_INFO" || name == "PREF" || name == "SAPISID" || name == "SID" || name == "SIDCC" || name == "SOCS" || name == "SSID" || name == "VISITOR_INFO1_LIVE" || name == "VISITOR_PRIVACY_METADATA" || name == "YSC" {
			cookie := &http.Cookie{
				Name:     name,
				Value:    decryptCookie(keys[0], encrypted_value),
				Path:     path,
				Domain:   domain,
				Expires:  tm,
				HttpOnly: httpOnlyBool,
			}
			cookiesList = append(cookiesList, cookie)
			err = cookie.Valid()
			if err != nil {
				panic(err)
			}
		}
	}
	//os.Exit(0)

	jar, err := cookiejar.New(nil)
	u, _ := url.Parse("https://www.youtube.com/")
	jar.SetCookies(u, cookiesList)
	row.Close()
	db.Close()
	return jar
}
