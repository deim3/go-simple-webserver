package main

import (
  "net/http"
  "log"
  "os"
  "flag"
  "github.com/BurntSushi/toml"
)

var config ConfigType

type ConfigType struct {
  Global struct {
    IP              string
    WebDir          string
    DirectoryIndex  string
    Port            string
    Logfile         string
  }
}

func init() {
  configFile := flag.String("config", "./webserver.conf", "Config path")
  flag.Parse()
  if _, err := os.Stat(*configFile); os.IsNotExist(err) {
    log.Fatal("Config file not found!\n")
  } else {
    if _, err := toml.DecodeFile(*configFile, &config); err != nil {
      log.Fatalln(err)
    }
  }
}

func main() {
  if len(config.Global.IP) == 0 {
    log.Fatal("Server IP is not specified.")
  } else if len(config.Global.WebDir) == 0 {
    log.Fatal("Webdir is not specified.")
  } else if len(config.Global.DirectoryIndex) == 0 {
    log.Fatal("DirectoryIndex is not specified.")
  } else if len(config.Global.Port) == 0 {
    log.Fatal("Port is not specified.")
  }
  fs := http.FileServer(http.Dir(config.Global.WebDir))
  http.Handle("/", fs)
  http.HandleFunc("/filter/", func(w http.ResponseWriter, r *http.Request){
    http.ServeFile(w, r, config.Global.DirectoryIndex);
  });
  listen := config.Global.IP + ":" + config.Global.Port
  err := http.ListenAndServe(listen, logRequest(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}

func logRequest(handler http.Handler) http.Handler {
  if len(config.Global.Logfile) == 0 {
    log.Fatal("Logfile is not specified.")
  }
  f, err := os.OpenFile(config.Global.Logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
  if err != nil {
    log.Fatal(err)
  }
  log.SetOutput(f)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
