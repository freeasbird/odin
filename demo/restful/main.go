package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/httplib"
	"github.com/offer365/example/endecrypt"
	pb "github.com/offer365/odin/demo/proto"
	"github.com/offer365/odin/utils"

	"strconv"
	"time"
)

const (
	Server_crt = `-----BEGIN CERTIFICATE-----
MIIDAzCCAesCCQDqO1aVDNi/IzANBgkqhkiG9w0BAQsFADBDMQswCQYDVQQGEwJH
QjEOMAwGA1UEBwwFQ2hpbmExDzANBgNVBAoMBmdvYm9vazETMBEGA1UEAwwKZ2l0
aHViLmNvbTAgFw0xOTEwMjAxMTAyMzhaGA8yMTE5MDkyNjExMDIzOFowQjELMAkG
A1UEBhMCR0IxDjAMBgNVBAcMBUNoaW5hMQ8wDQYDVQQKDAZzZXJ2ZXIxEjAQBgNV
BAMMCXNlcnZlci5pbzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMRF
1mgbKYNO2X0iqX89Rrzc+xftqegQ+7V0n9Sa1HE07xQVcgu05faGB4B/29HPQ/gh
JMt1IxkXlISNuQwIDM5XVSahkH1OhQmtQnTKjYXFgboRFHMQUk26lKoIZ3o9AJ8s
QTPCLBw7a9StBpeWhBzEumDymP60hmGhTft4tbY85MrmObfTZ8KbQiHIy22jXNGV
N5ok61q4tlMV8HYK89q4WX7IcQusdK9NNwL1jZNQ4+WICEe2/zs8xY9r4REONKoM
HOME5aS+EvQSVwh5LyvNuPxa8io83EjokT3yRqZllvmXD/hVS/BCM927fgsiDfm0
ezuE5+AGiMR1N0agv1cCAwEAATANBgkqhkiG9w0BAQsFAAOCAQEAHxXVd/v7noVZ
LJ8IsLty3BjMX7ZjSvkyyrchxdQQIfCoMc/UGkDZ5TgvdPkE8eAfdSVwwrcpGf8V
C4ccB9flekd6HrO7Uo9mWrKcjyQn2MjwAZNDhcs5Sxrz8TusJEQk4iQYSq0oc4Nr
qGrR/2kXEirwXi/xQ0saVXalfhkK5W+rO/YWTc8K3leARQ6BDjGbF2BHRtj6HEZL
RnhJEbx+BvplXMlWQ5CBBYt/NQa/MKJDd2stT70Si8E1lIGIGaVQAy43uT7xy8XW
jSrruOAv1SVLovhSxjsMiu/jXwZsVAtaFAuT4ajiWQHzbNqUjVnt7dJIWJPCnL6h
lhaV0MUy7Q==
-----END CERTIFICATE-----

`
	Server_key = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAxEXWaBspg07ZfSKpfz1GvNz7F+2p6BD7tXSf1JrUcTTvFBVy
C7Tl9oYHgH/b0c9D+CEky3UjGReUhI25DAgMzldVJqGQfU6FCa1CdMqNhcWBuhEU
cxBSTbqUqghnej0AnyxBM8IsHDtr1K0Gl5aEHMS6YPKY/rSGYaFN+3i1tjzkyuY5
t9NnwptCIcjLbaNc0ZU3miTrWri2UxXwdgrz2rhZfshxC6x0r003AvWNk1Dj5YgI
R7b/OzzFj2vhEQ40qgwc4wTlpL4S9BJXCHkvK824/FryKjzcSOiRPfJGpmWW+ZcP
+FVL8EIz3bt+CyIN+bR7O4Tn4AaIxHU3RqC/VwIDAQABAoIBAAQbHeghoVWw4ZXf
ksIpqwAqc0pF24cSS+G45dsRvh38KIA4DqG2EBV/KksC4bta5aYcM2PaOHi+6Il5
WYSp6nKqmwpq2NX2PYw9RqWg0yMYRaV50/6wObiMja2c7WU+P3QU/ewyRK/2gkP5
tqiXKn5bkzaR/KdfaWxDbpkzJkIArLAELqEBuS0noxikrfypPanGnXk7IDhYo+rZ
cE0UHOhpkeo7gXeVc9tU/cjTRwBK7awKLIDWyknHGrL28nxMqKf06SzxG2oz6Hn3
twOtwAUS7tjophOZ6WCStgCOVFf0Ue6yJmja9xgWy/r2jJsH5/VV0xJZvmWGxr8T
IQh4oskCgYEA92Katy0Cvl1kS1/cf0ExMtOzXIwtCDu35axGl1FR3VMcoboPmH2h
HrRxSpcIgkRXz7wxsj3zttBXu8assjmwtCWzbDIE0YGYQ3v1CwDITihAyhevhW4b
UxN181RhMo1qHIcgULsVR5+P857FAHRSWWewh77ZK7x17fdQJshZujMCgYEAyxuT
R1CthfC7rbIX359tD9jb1XtG+XCgygZYv+6uoknmWMMmUqgDmQ3u8p4kuHudB6gm
/kZkxrluwJM5B8UKC1NRkejHP2ZO8ygpEGQp7t1H3BBFSfUVlu+YmfD5SjHhK9U5
2t+hfyuO8m0r+XdYk6lliEYufVlPMzJffT3rSk0CgYEAhs+jRGMw9ZBrUXAB9w8N
wob/XVW+TJhOlMiXB2r3U8cw+SktyonbvaHTgzRfHK4ltDz4UAvWvi83QEr6XX12
wBUze6ieW5Vl5pCsbryUa5MgC4Fw0yO3nEQkqN+4wBW0V6uDfrsU050ukzJYZPD+
113cI31rV5wyH+YANcJEs2UCgYEAmh0SY8qT4E4KGoJIGyadWqjyJcqk0CDl4GVw
cjJp0DrCzhdFvPI/yKMJ7I6Szmj9fhHZhJdlYGTT5MvROlQIiw9tlYlLpo+62EZg
4k8egmDlZdXyvWt6Nk0XPbfbcLDoapogjDOkFxq2HL054NDuJR0kLYMTQ4nAztgq
HJ4fKwECgYEAinsJM6lw9m3eyRRuPRFE4jNwg5KmZRjVuZ06+nPW/Sb7GXdN+5e6
62y87e63MRTm1r2C4g/esnqAOcS6iRHQtdTFrG8DU/j9F+uaB5TWZTroxqQ6h2F0
OjGZcdCMohluWRztbas01OZKSoDx1pEfP+H4kKFJ8LhWQXLU0lWibEw=
-----END RSA PRIVATE KEY-----

`
	Client_crt = `-----BEGIN CERTIFICATE-----
MIIDAzCCAesCCQDqO1aVDNi/JDANBgkqhkiG9w0BAQsFADBDMQswCQYDVQQGEwJH
QjEOMAwGA1UEBwwFQ2hpbmExDzANBgNVBAoMBmdvYm9vazETMBEGA1UEAwwKZ2l0
aHViLmNvbTAgFw0xOTEwMjAxMTAyMzhaGA8yMTE5MDkyNjExMDIzOFowQjELMAkG
A1UEBhMCR0IxDjAMBgNVBAcMBUNoaW5hMQ8wDQYDVQQKDAZjbGllbnQxEjAQBgNV
BAMMCWNsaWVudC5pbzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAObU
dbPa2MMnp7X0P6TeUM9+gJgRgVdrOm05EPnf4p1xEmFq09bGupZpD+pVoU/yH/oi
wA4gwYtgk5ETtTfTbF8YUma6LYDye2m98zXiyVWpTs9pmxVRUTcpnjmyIS7mXSNE
hShN26OCTk8DtlL9STFnFWQY2Sb9PVjwDWTrXkHalQU3PFEmoQ/QPbTbBN2gydDn
WkK6LxgTaSA9xMw/j5upZh58aoLVwd8IevzKn/YnwQBEC0ptVQGl6B5EUKabhTWh
q6c4gDAhcqhdRFZa4UMcOZnzgwEuR7XzJlTL3AanBXJu5sUjDPTweOhENcdSHBQ/
sX6Cr9NFRm6bQqOrjmsCAwEAATANBgkqhkiG9w0BAQsFAAOCAQEAD4IIkxuNITUM
bHU2ebLEMPq8Udhcl9s3mBlaWf3ecDi4Yu+MBTy+ggcRhnq7zqaRVaRxdhtyVVIA
3hFwrWZK38jPGKrI9qZLoqQJu3RFq241jOjVol6zAkYhuqvO8n9AKhShjUFHPfA+
TN1BC8qb30lwYZnELaHdKFM16f7uska2lMY6N8uYqySNWFz/B77zIqUACRnvyGfS
gJ8QRDcGAjA0+SEMKtI0tB0qsL4c+de8uPaUjyO5uzWLXJap50gBi5Zx17YE8aMk
wntulWmvYwJSokZIOVi+3PDSc8Zh2ukhddA716NF2U8c7N++BZFBGcGLZK86Hh76
HayLryJQmA==
-----END CERTIFICATE-----

`
	Client_key = `-----BEGIN RSA PRIVATE KEY-----
MIIEpgIBAAKCAQEA5tR1s9rYwyentfQ/pN5Qz36AmBGBV2s6bTkQ+d/inXESYWrT
1sa6lmkP6lWhT/If+iLADiDBi2CTkRO1N9NsXxhSZrotgPJ7ab3zNeLJValOz2mb
FVFRNymeObIhLuZdI0SFKE3bo4JOTwO2Uv1JMWcVZBjZJv09WPANZOteQdqVBTc8
USahD9A9tNsE3aDJ0OdaQrovGBNpID3EzD+Pm6lmHnxqgtXB3wh6/Mqf9ifBAEQL
Sm1VAaXoHkRQppuFNaGrpziAMCFyqF1EVlrhQxw5mfODAS5HtfMmVMvcBqcFcm7m
xSMM9PB46EQ1x1IcFD+xfoKv00VGbptCo6uOawIDAQABAoIBAQCTFY5qrGiy8fHL
33cudvrHPLR0MbNZINp5/oLytdaQvBwaNxgFI1yBuzCJAUdoyb/Wg44dcoHhbgiZ
yRUQHYhQkA7xpnCYWeqJ1p/DFl90Vg4B3CkVzFsT61EHMpoyaFewwViX9gSei8ma
T6M9/mdFM4pN3geA8JzGry/ZvqCxFID3Sz4/08zq9UjS54GiZgJb3lyGazdDk3Gf
h2NukbBRtvdh8iILjEM38czgqTBrDqXlFa5q0p9oq+UPn9twaVZcJ9t4IrcWIgaD
l9cYRE/agXj0cRO/IVOi/RB0e/NLiR0XqXSo4Rx7uGcSJys1yuPt96OGMIh2+c99
VGJbzBsBAoGBAP8qagFe5kJNrjweo9yNhs0H/TFx1mhCqQKPNFouDtttaCDcNvXx
d3B5KYKgWpTJPaZ1eGfPeA4OTLhKCLVG7EVQXUUsztyDS1JpuUJkm1texA30g0sw
UWhLfQfFEgWCaIQkbqZv1J5OYrc2xvPqjHfP+NG1CAte1w5QQ7FA541fAoGBAOeV
rO0yF30sDOUXlixfKN35j2FIgVB0DOT6nkpPyh1OYcdcshu3utGqmOiN7twqwyiL
m3Uucix/JbTb2m+HSAX9/s/SHHOoXeUp21wVSGYesknrBEZt3VifINzu/OFCjLk1
Plx4F0am0WrsDnAtQwgpCV29lgQjmFsXQZlUW051AoGBAJdvpbAgkUmCbsixapCn
0fv3JNZmeFgyT7n8IZbvxNOHkAgIifnXEArJbdBfuMKa2KLlDsuVfuvgormw/pAP
goP0mRZH7JFEvrwvkMqNiQJmMLcTiaRjDb13J8InvHVWmw7pzF2s+yPk44NW2CbE
6g7leAeFiDuvUrTk//e/zGzDAoGBAJx4TLaWubghIzVGkni4cuxHydB5JKYvQucT
Tg/3iR/z7ay9vLltkhRHp7i47UJkwieK7CZok0vtPJTOVvAz/z3NN3VDCWY7w/Uq
KsQ0vQ4Cf4Ph/ql3Ya6XFaUw9Dtes6YPi2r+2PsriyMrCzZP3pKM538msU1qn24s
cG4gyPBhAoGBAL+VTkIaLK07qChlT0Y2hwbmfLwAlOrPguJps7D0C0aBUDPXylOO
S5myV8jp+htbP6Mn5MEzZHhvoVSEe9GiCv9E5KMisJjPtQRRRKGNPAnTt9KJVQ1U
BCggzbZzimK/EFR72woV0071B93C4jO07jEmvkCb3gzmyWkgjREZQusj
-----END RSA PRIVATE KEY-----

`
	Ca_crt = `-----BEGIN CERTIFICATE-----
MIIDWzCCAkOgAwIBAgIJAKst9d2m1o1TMA0GCSqGSIb3DQEBCwUAMEMxCzAJBgNV
BAYTAkdCMQ4wDAYDVQQHDAVDaGluYTEPMA0GA1UECgwGZ29ib29rMRMwEQYDVQQD
DApnaXRodWIuY29tMCAXDTE5MTAyMDExMDIzOFoYDzIxMTkwOTI2MTEwMjM4WjBD
MQswCQYDVQQGEwJHQjEOMAwGA1UEBwwFQ2hpbmExDzANBgNVBAoMBmdvYm9vazET
MBEGA1UEAwwKZ2l0aHViLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC
ggEBALh9e8OuCRP0zMYjbCqUk5b+d6J3tC9INL3P0VwcmWx5jCpUQLz0SGafnIL8
LworJfQkbDdOKNol9Zt4vzsxdV1k2VaZuAY0qWG5Kg+n1tCml46By9mQH3B6ngKe
cNDdBmRGYYDkuqI9g8UBgRYT4TbIQJ1Ns4wuKQR02/kCUfWypvE+8bEQEXTRKcHo
inILmcO7RvhWkfwWVbTpUv7M7K8wwIGKawDgl3DeW5g+tss0PD/iCdMo0DMRHykx
4KeTsrPYdxpxgf42LwG0aJ+/28GzYCQ4mYJaTADr5pp+vlUZWtYK8m7fFXbpGlrU
5aLTA5aEPdIuyTa2/DZXl4JBxTkCAwEAAaNQME4wHQYDVR0OBBYEFCikHb0Ms/7f
jci0C5Amwvf7cFmYMB8GA1UdIwQYMBaAFCikHb0Ms/7fjci0C5Amwvf7cFmYMAwG
A1UdEwQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEBAG+jH3wHkTqx97/9voaftE/b
0tkbV+9b3SxPv5KoW0fm24x6UNrMPE9APt0J00Vlv20LNc/tOWquyKGDIhhe29/x
ehte/l7doGVW0Wg3xQtiIT5aJdMHNy+bSLogzV5D5sbHcPStKNj3M1wwhMj03YZ7
Nj5ua/c5aUU+MBMv0C/FNPnB+GSeRO2MxYHsZP2mBEJaLhPZ+K29kwGPCVWIESCH
IOS/jew/kfpPLavuvyPqoGAfc1xpe6QQXZUEGCtzTDU/rl/hQWMxCJg85E1S5Usx
gahmAgIzeyFCjb2txOo65VtLM0DfzzkIX2PrLz7CyiXP40m8uBMtCDG+IZS0arQ=
-----END CERTIFICATE-----
`
)

var tlsConifg *tls.Config

func init() {
	tlsConifg = Tls([]byte(pb.Client_crt), []byte(pb.Client_key), []byte(pb.Ca_crt))
}

func NewApp(name, id, token string, servers []string) (app *Application) {
	app = &Application{
		Name:    name,
		ID:      id,
		Servers: servers,
		Token:   token,
		Cycle:   new(Cycle),
	}
	return
}

type Application struct {
	Name    string
	ID      string
	Servers []string
	Token   string
	*Cycle
	tlsConfig *tls.Config
}

// 一次认证请求
type Cycle struct {
	Urls     []string
	Uuid     string
	AuthInfo string
	Lease    int64
}

func Tls(crt, key, ca []byte) *tls.Config {
	certificate, err := tls.X509KeyPair(crt, key)
	if err != nil {
		return nil
	}
	certPool := x509.NewCertPool()

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		// err = errors.New("failed to append ca certs")
		return nil
	}
	return &tls.Config{
		Certificates:       []tls.Certificate{certificate},
		ClientAuth:         tls.NoClientCert, // NOTE: 这是可选的!
		ClientCAs:          certPool,
		InsecureSkipVerify: true,
		Rand:               rand.Reader,
		Time:               time.Now,
		// NextProtos:         []string{"http/1.1", http2.NextProtoTLS},
	}

}

// auth
func (app *Application) Active() {
	verify := make(map[string]interface{})
	verify["app"] = app.Name
	verify["id"] = app.ID
	verify["date"] = time.Now().Unix()
	verify["token"] = app.Token
	byt, err := json.Marshal(verify)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(byt))
	byt, err = endecrypt.Encrypt(endecrypt.Aes2key32, byt)
	if err != nil {
		fmt.Println(err)
		return
	}
	req := &pb.Request{
		App:    app.Name,
		Id:     app.ID,
		Date:   time.Now().Unix(),
		Verify: base64.StdEncoding.EncodeToString(byt),
	}
	byt, _ = json.Marshal(req)
	fmt.Println(string(byt))
	result, err := httplib.Post(app.Servers[0]).SetTimeout(2*time.Second, 3*time.Second).Debug(true).SetBasicAuth("admin", "123").SetTLSClientConfig(tlsConifg).Header("Content-Type", "application/json; charset=utf-8").JSONBody(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp := &pb.Response{}
	fmt.Println(result.ToJSON(resp))
	fmt.Println(resp.Code, resp.Msg)
	byt, err = endecrypt.Decrypt(endecrypt.Pub2Rsa1024, resp.Data.Cipher) // uuid
	app.Uuid = string(byt)
	byt, err = endecrypt.Decrypt(endecrypt.Pub2Rsa2048, resp.Data.Auth) // {"attrs":[{"Name":"热词","Key":"hotword","Value":1000},{"Name":"类热词","Key":"classword","Value":1000}],"time":1571909203232224000}
	app.AuthInfo = string(byt) // eg: {"attrs":[{"Name":"热词","Key":"hotword","Value":1000},{"Name":"类热词","Key":"classword","Value":1000}],"time":1571909203232224000}
	fmt.Println(string(byt), err)
	app.Lease = resp.Data.Lease
}

func (app *Application) KeepLine() {
	req := &pb.Request{
		App:    app.Name,
		Id:     app.ID,
		Date:   time.Now().Unix(),
		Verify: "",
		Umd5:   utils.Md5sum([]byte(app.Uuid), nil),
		Lease:  app.Lease,
	}

	for range time.Tick(time.Second * 6) {
		byt, _ := json.Marshal(req)
		fmt.Println(string(byt))
		result, _ := httplib.Put(app.Servers[0]).SetTimeout(2*time.Second, 3*time.Second).Debug(true).SetBasicAuth("admin", "123").SetTLSClientConfig(tlsConifg).Header("Content-Type", "application/json; charset=utf-8").JSONBody(req)
		resp := &pb.Response{}
		result.ToJSON(resp)
		fmt.Println(result.String())
		fmt.Println(resp.Data.Lease)
	}
}

func (app *Application) OffLine() {
	req := &pb.Request{
		App:    app.Name,
		Id:     app.ID,
		Date:   time.Now().Unix(),
		Verify: "",
		Umd5:   utils.Md5sum([]byte(app.Uuid), nil),
		Lease:  app.Lease,
	}

	result, _ := httplib.Delete(app.Servers[0]).SetTimeout(2*time.Second, 3*time.Second).Debug(true).SetBasicAuth("admin", "123").SetTLSClientConfig(tlsConifg).Header("Content-Type", "application/json; charset=utf-8").JSONBody(req)
	resp := &pb.Response{}
	result.ToJSON(resp)
	fmt.Println(resp.Data.Lease)

}

func main() {
	ManyApp(1000)
	// SingleAPP()
	select {}

}

func ManyApp(ins int) {
	apps := make([]*Application, 0)
	for i := 0; i < ins; i++ {
		token := "app" + strconv.Itoa(i)
		app := NewApp("nlp", "app"+strconv.Itoa(i), token, pb.Member)
		apps = append(apps, app)
	}
	for _, ap := range apps {
		time.Sleep(time.Millisecond * 500)
		go func(app *Application) {
			time.Sleep(time.Second)
			app.Active()
			app.KeepLine()
		}(ap)
	}
}

func SingleAPP() {
	app := NewApp("hotword", "app"+strconv.Itoa(99), "app"+strconv.Itoa(99), pb.Member, )
	app.Active()
	app.KeepLine()
	app.OffLine()
	fmt.Println(111)
}
