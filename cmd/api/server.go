package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	mw "restapi/internal/api/middlewares"
	"restapi/internal/api/routers"
	"restapi/internal/repository/sqlconnect"
	"restapi/pkg/utils"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}

	_, err = sqlconnect.ConnectDb()
	if err != nil {
		panic(err) //panic because without db, whole api wont work
	}

	PORT := os.Getenv("API_PORT")

	cert := "cert.pem"
	key := "key.pem"

	mux := routers.MainRouter()
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	rl := mw.NewRateLimiter(6, time.Minute)

	hpp := mw.HPPOptions{
		CheckQuerry:             true,
		CheckBody:               true,
		CheckBodyForContentType: "application/x-www-form-urlencoded",
		Whitelist:               []string{"sortby", "sortorder", "name", "age", "class"},
	}

	//secureMux := mw.Cors(rl.MiddleWare(mw.ResponseTime(mw.Compression(mw.SecurityHeader(mw.Hpp(hpp)(mux))))))
	secureMux := utils.ApplyMiddlewares(
		mux,
		mw.Hpp(hpp),
		mw.Compression,
		mw.SecurityHeader,
		mw.ResponseTime,
		rl.MiddleWare,
		mw.Cors,
	)

	server := &http.Server{
		Addr:      fmt.Sprintf(":%s", PORT),
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("server started at port:", PORT)
	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalf("error starting server: %s", err)
		return
	}
}
