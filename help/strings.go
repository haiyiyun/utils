package help

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

	"go.haiyiyun.org/uuid"
)

type Strings string

func NewString(i interface{}) Strings {
	s := Strings(fmt.Sprintf("%v", i))
	return s
}

func (s Strings) String() string {
	return string(s)
}

func (s Strings) Md5() string {
	m := md5.New()
	io.WriteString(m, s.String())

	return fmt.Sprintf("%x", m.Sum(nil))
}

func (s Strings) RandString(width int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	uORl := []int{65, 97}
	bs := make([]byte, width)
	for i := 0; i < width; i++ {
		b := r.Intn(25) + uORl[r.Intn(2)]
		bs[i] = byte(b)
	}

	return string(bs)
}

func (s Strings) RandNumber(width int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	numeric := []byte("0123456789")
	lnum := len(numeric)
	bs := make([]byte, width)

	for i := 0; i < width; i++ {
		bs[i] = numeric[r.Intn(lnum)]
	}

	return string(bs)
}

func (s Strings) RandMixed(width int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	numeric := []byte("0123456789")
	lnum := len(numeric)
	rs := s.RandString(width)
	rsb := []byte(rs)
	lrsb := len(rsb)
	n := r.Intn(lrsb)
	for i := 0; i < n; i++ {
		nk := r.Intn(lrsb)
		nn := r.Intn(lnum)
		rsb[nk] = numeric[nn]
	}

	return string(rsb)
}

// delimiterStyle: '-'
// convert like this: "HelloWorld" to "hello-world"
func (s Strings) SnakeCasedNameByDelimiterStyle(delimiterStyle rune) string {
	newstr := make([]rune, 0)
	firstTime := true

	for _, chr := range string(s) {
		if isUpper := 'A' <= chr && chr <= 'Z'; isUpper {
			if firstTime == true {
				firstTime = false
			} else {
				newstr = append(newstr, delimiterStyle)
			}
			chr -= ('A' - 'a')
		}
		newstr = append(newstr, chr)
	}

	return string(newstr)
}

// convert like this: "HelloWorld" to "hello_world"
func (s Strings) SnakeCasedName() string {
	return s.SnakeCasedNameByDelimiterStyle('_')
}

// delimiterStyle: '-'
// convert like this: "hello-world" to "HelloWorld"
func (s Strings) TitleCasedNameByDelimiterStyle(delimiterStyle rune) string {
	newstr := make([]rune, 0)
	upNextChar := true

	for _, chr := range string(s) {
		switch {
		case upNextChar:
			upNextChar = false
			chr -= ('a' - 'A')
		case chr == delimiterStyle:
			upNextChar = true
			continue
		}

		newstr = append(newstr, chr)
	}

	return string(newstr)
}

// convert like this: "hello_world" to "HelloWorld"
func (s Strings) TitleCasedName() string {
	return s.TitleCasedNameByDelimiterStyle('_')
}

func (s Strings) PluralizeString() string {
	str := string(s)
	if strings.HasSuffix(str, "y") {
		str = str[:len(str)-1] + "ie"
	}
	return str + "s"
}

func (s Strings) GenerateSecret(secret string) string {
	h := sha512.New()
	h.Write([]byte(string(s)))
	newUUID := uuid.NewMD5(uuid.Must(uuid.NewRandom()), []byte(secret))
	newSecret := base64.URLEncoding.EncodeToString(h.Sum(newUUID.Bytes()))
	newSecret = strings.TrimRight(newSecret, "=")
	return newSecret
}

func (s Strings) TrimSide(cutset string) string {
	str := string(s)
	str = strings.TrimLeft(str, cutset)
	return strings.TrimRight(str, cutset)
}
