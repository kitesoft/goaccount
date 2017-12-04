package main

import (
	"context"
	"flag"
	"fmt"

	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"

	"goaccount/config"
	"goaccount/log"
	"goaccount/model"
	"goaccount/server"
	"goaccount/util"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:  %s [options] ...\n", os.Args[0])
	flag.PrintDefaults()
}

var (
	// Flags
	helpShort  = flag.Bool("h", false, "Show usage text (same as --help).")
	helpLong   = flag.Bool("help", false, "Show usage text (same as -h).")
	serverIp   = flag.String("ip", "127.0.0.1", "the server ip")
	serverPort = flag.Int("p", 7200, "the server port")
	configFile = flag.String("c", "config.yml", "config file")
)

var srvConfig = &config.Config{}

func init() {
	if err := srvConfig.Load(*configFile); err != nil {
		logrus.Error(err.Error())
		panic(err.Error())
	}
	log.InitLog(srvConfig.Server.Name)
	logrus.Debugf("%v", srvConfig)

	var err error
	logrus.Debugf("%s", srvConfig.MysqlConfig().String())
	model.DB, err = gorm.Open("mysql", srvConfig.MysqlConfig().String())
	if err != nil {
		logrus.Debugf("error %s", err.Error())
		panic("failed to connect database")
	}

	//	model.DB.AutoMigrate()
}

func main() {

	flag.Usage = usage
	flag.Parse()
	if *helpShort || *helpLong {
		flag.Usage()
		return
	}
	rand.Seed(time.Now().UnixNano())

	pearCors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "HEAD", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	accountService := NewAccountService(
		new(util.DecodeAndValidator),
		PubJWT(srvConfig.Server.Name),
	)
	router := mux.NewRouter()

	accountService.RegisterRoutes(router, srvConfig.Server.UrlPrefix)

	handler := negroni.New(negroni.NewRecovery(),
		server.NewLogger(),
		server.NewUberRatelimit(srvConfig.Server.ReteLimit),
		server.NewParseForm())
	handler.Use(pearCors)
	handler.UseHandler(router)

	srv := &http.Server{
		Handler:        handler,
		ReadTimeout:    2 * time.Minute,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go server.StartHttpServer(srv, fmt.Sprintf(":%d", *serverPort), srvConfig.Server.ConnLimit)

	// subscribe to SIGINT signals
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)
	<-stopChan // wait for SIGINT
	logrus.Info("Shutting down server...")

	// shut down gracefully, but wait no longer than 5 seconds before halting
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
	logrus.Info("Server gracefully stopped")
}
