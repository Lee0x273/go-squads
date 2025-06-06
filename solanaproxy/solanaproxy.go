package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

var dst = flag.String("t", "https://api.mainnet-beta.solana.com", "target URL, -t https://api.mainnet-beta.solana.com")
var wsDst = flag.String("ws", "wss://api.mainnet-beta.solana.com", "WebSocket target URL")
var bind = flag.String("bind", ":8080", "bind host and port,like 127.0.0.1:8080")

var keyFile = flag.String("keyfile", "", "")
var certFile = flag.String("certfile", "", "")

func main() {
	flag.Parse()
	targetURL, err := url.Parse(*dst)
	if err != nil {
		log.Fatal(err)
	}

	wsURL, err := url.Parse(*wsDst)
	if err != nil {
		log.Fatal(err)
	}

	// WebSocket upgrader
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Connection"), "Upgrade") &&
			strings.Contains(r.Header.Get("Upgrade"), "websocket") {
			// WebSocket 处理
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Printf("WebSocket upgrade error: %v", err)
				return
			}
			defer conn.Close()

			// 连接到目标 WebSocket 服务器
			wsConn, _, err := websocket.DefaultDialer.Dial(wsURL.String(), nil)
			if err != nil {
				log.Printf("WebSocket dial error: %v", err)
				return
			}
			defer wsConn.Close()

			// 双向转发 WebSocket 消息
			go func() {
				for {
					messageType, message, err := conn.ReadMessage()
					if err != nil {
						return
					}
					if err := wsConn.WriteMessage(messageType, message); err != nil {
						return
					}
				}
			}()

			for {
				messageType, message, err := wsConn.ReadMessage()
				if err != nil {
					return
				}
				if err := conn.WriteMessage(messageType, message); err != nil {
					return
				}
			}
		} else {
			// HTTP 处理
			log.Printf("handle %s from %s\n", r.RequestURI, r.RemoteAddr)
			r.Host = targetURL.Host
			r.URL.Scheme = targetURL.Scheme
			r.URL.Host = targetURL.Host
			proxy.ServeHTTP(w, r)
		}
	})

	// 启动代理服务器
	log.Println("bind", *bind, "proxy to ", (*targetURL).String())
	if *keyFile != "" && *certFile != "" {
		log.Fatal(http.ListenAndServeTLS(*bind, *certFile, *keyFile, nil))
	} else {
		log.Fatal(http.ListenAndServe(*bind, nil))
	}

}
