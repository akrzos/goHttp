package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

var livez bool = false
var readyz bool = false
var livenessCount int = 0

func livenessDelay(delay int) {
	log.Print("Starting livez delay...")
	time.Sleep(time.Duration(delay) * time.Second)
	livez = true
	log.Print("Completed livez delay")
}

func readinessDelay(delay int) {
	log.Print("Starting readyz delay...")
	time.Sleep(time.Duration(delay) * time.Second)
	readyz = true
	log.Print("Completed readyz delay")
}

func main() {
	log.Print("Starting the server...")

	port := os.Getenv("PORT")
	listenDelay := os.Getenv("LISTEN_DELAY_SECONDS")
	livezDelay := os.Getenv("LIVENESS_DELAY_SECONDS")
	readyzDelay := os.Getenv("READINESS_DELAY_SECONDS")
	responseDelay := os.Getenv("RESPONSE_DELAY_MILLISECONDS")
	livenessSuccessMax := os.Getenv("LIVENESS_SUCCESS_MAX")

	if port == "" {
		port = "8000"
		log.Print("Using default port 8000")
	} else {
		log.Print("Using port " + port)
	}

	if listenDelay == "" {
		listenDelay = "10"
		log.Print("Using listen delay default of 10s")
	} else {
		log.Print("Using listen delay " + listenDelay + "s")
	}
	listenDelaySeconds, err := strconv.Atoi(listenDelay)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to convert listenDelay"))
	}

	if livezDelay == "" {
		livezDelay = "2"
		log.Print("Using live delay default of 2s")
	} else {
		log.Print("Using live delay " + listenDelay + "s")
	}
	livezDelaySeconds, err := strconv.Atoi(livezDelay)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to convert listenDelay"))
	}

	if readyzDelay == "" {
		readyzDelay = "10"
		log.Print("Using readiness delay default of 10s")
	} else {
		log.Print("Using readiness delay " + readyzDelay + "s")
	}
	readyzDelaySeconds, err := strconv.Atoi(readyzDelay)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to convert readyzDelay"))
	}

	if responseDelay == "" {
		responseDelay = "0"
		log.Print("Using response delay default of 0")
	} else {
		log.Print("Using response delay " + responseDelay + "ms")
	}
	responseDelayMilliSeconds, err := strconv.Atoi(responseDelay)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to convert responseDelay"))
	}

	if livenessSuccessMax == "" {
		livenessSuccessMax = "0"
		log.Print("Using liveness success max default of 0")
	} else {
		log.Print("Using liveness success max " + livenessSuccessMax)
	}
	livenessCountMax, err := strconv.Atoi(livenessSuccessMax)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to convert livenessSuccessMax"))
	}

	if listenDelaySeconds > 0 {
		log.Print("Starting listen delay...")
		time.Sleep(time.Duration(listenDelaySeconds) * time.Second)
		log.Print("Completed listen delay")
	} else {
		log.Print("No listen delay")
	}

	go livenessDelay(livezDelaySeconds)
	go readinessDelay(readyzDelaySeconds)

	http.HandleFunc("/home", func(w http.ResponseWriter, _ *http.Request) {
		if responseDelayMilliSeconds != 0 {
			time.Sleep(time.Duration(responseDelayMilliSeconds) * time.Millisecond)
		}
		if readyz {
			log.Print("/home request when ready")
			fmt.Fprint(w, "/home request processed")
		} else {
			log.Print("/home request when not ready")
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	},
	)

	http.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		if responseDelayMilliSeconds != 0 {
			time.Sleep(time.Duration(responseDelayMilliSeconds) * time.Millisecond)
		}
		if readyz {
			log.Print("/readyz request when ready")
			fmt.Fprint(w, "/readyz request processed")
		} else {
			log.Print("/readyz request when not ready")
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	},
	)

	http.HandleFunc("/livez", func(w http.ResponseWriter, _ *http.Request) {
		if responseDelayMilliSeconds != 0 {
			time.Sleep(time.Duration(responseDelayMilliSeconds) * time.Millisecond)
		}
		if livez {
			if livenessCountMax != 0 {
				livenessCount++
				if livenessCount > livenessCountMax {
					log.Print("/livez request after liveness success count exceeded " + strconv.Itoa(livenessCount) + "/" + strconv.Itoa(livenessCountMax))
					w.WriteHeader(http.StatusServiceUnavailable)
				} else {
					log.Print("/livez request when ready " + strconv.Itoa(livenessCount) + "/" + strconv.Itoa(livenessCountMax))
					fmt.Fprint(w, "/livez request processed")
				}
			} else {
				log.Print("/livez request when ready")
				fmt.Fprint(w, "/livez request processed")
			}
		} else {
			log.Print("/livez request when not ready")
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	},
	)

	http.HandleFunc("/crash", func(w http.ResponseWriter, _ *http.Request) {
		if responseDelayMilliSeconds != 0 {
			time.Sleep(time.Duration(responseDelayMilliSeconds) * time.Millisecond)
		}
		if readyz {
			log.Print("/crash request when ready")
			fmt.Fprint(w, "/crash request processed")
		} else {
			log.Print("/crash request when not ready")
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		log.Fatal("/crash endpoint received a request, crashing...")
	},
	)

	log.Print("The service is listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
