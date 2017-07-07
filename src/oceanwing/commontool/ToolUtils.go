package commontool

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/cihub/seelog"
)

// CA hh..
var (
	CA       *tls.Config
	SubSinal = make(chan struct{})
)

// GenerateClientID when connect to a mqtt broker, need a clientid(for ankerbox.)
func GenerateClientID() string {
	return strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
}

// GetCurrentTime hh..
func GetCurrentTime() string {
	return time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
}

// GetTimeAsFileName hh.
func GetTimeAsFileName() string {
	return time.Unix(time.Now().Unix(), 0).Format("2006-01-02-15-04-05")
}

// RandInt64 取值范围：大于等于 min, 小于 max
func RandInt64(min, max int64) int64 {
	maxBigInt := big.NewInt(max)
	i, _ := rand.Int(rand.Reader, maxBigInt)
	if i.Int64() < min {
		RandInt64(min, max)
	}
	return i.Int64()
}

// ReadFileContent  hh.
func ReadFileContent(fpath string) ([]string, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// Buildbytes hh.
func Buildbytes(command string) []byte {
	//command := "A50101000000010101015A0000003604FFFF98FA"
	command = strings.ToUpper(command)
	//debug
	// fmt.Println(command)
	strlen := len(command) / 2
	b := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		pos := i * 2
		b[i] = byte((findByteIndex(command[pos]) << 4) | findByteIndex(command[pos+1]))
	}
	return b
}

func findByteIndex(c byte) byte {
	return byte(strings.IndexByte("0123456789ABCDEF", c))
}

// BuildTlSConfig hh.
func BuildTlSConfig(caPath string) {
	if caPath != "" {
		certPool := x509.NewCertPool()
		if pemCerts, err := ioutil.ReadFile(caPath); err == nil {
			certPool.AppendCertsFromPEM(pemCerts)
		} else {
			panic(err)
		}
		CA = &tls.Config{
			RootCAs:            certPool,
			ClientAuth:         tls.NoClientCert,
			ClientCAs:          nil,
			InsecureSkipVerify: false,
		}
	}
}

// InitLogInstance hh.
func InitLogInstance(level string) {
	logformat := `
		<seelog minlevel="%s">
			<outputs formatid="main">
			    <console />
				<buffered size="10000" flushperiod="1000">  
					<file path="./%s"/>
        		</buffered>
			</outputs>
			<formats>
				<format id="main" format="%Date %Time [%LEVEL] %Msg%n"/>
			</formats>
		</seelog>`
	logConfig := fmt.Sprintf(logformat, level, GetTimeAsFileName()+"-log.log")
	logger, _ := log.LoggerFromConfigAsBytes([]byte(logConfig))
	log.ReplaceLogger(logger)
}
