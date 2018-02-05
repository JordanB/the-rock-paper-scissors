package main

import (
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "os"
    "time"
)

const rock string = "rock"
const paper string = "paper"
const scissors string = "scissors"

var moves []string = []string{rock, paper, scissors}
var base int = 300
var leftThreshold int = base / 3
var rightThreshold int = base * 2 / 3

func main() {
    http.HandleFunc("/health", healthHandler)
    http.HandleFunc("/newgame", newGameHandler)
    http.HandleFunc("/", moveHandler)

    port := os.Getenv("PORT")
    if port == "" {
        port = "5000"
    }
    
    log.Print("Listening on port " + port)
    log.Fatal(http.ListenAndServe(":" + port, nil))
}

func healthHandler(writer http.ResponseWriter, request *http.Request) {
    fmt.Fprint(writer, "ok")
}

func newGameHandler(writer http.ResponseWriter, request *http.Request) {
    seed := rand.NewSource(time.Now().UnixNano())
    random := rand.New(seed)
    move := moves[random.Intn(3)]
    // [ ... rock ... | ... paper ... | ... scissors ... ]
    //                ^               ^
    //            leftThres       rightThres
    if move == rock {
        leftThreshold, rightThreshold = 250, 275
    } else if move == paper {
        leftThreshold, rightThreshold = 25, 275
    } else if move == scissors {
        leftThreshold, rightThreshold = 25, 50
    }
    fmt.Fprint(writer, move)
}

func moveHandler(writer http.ResponseWriter, request *http.Request) {
    if request.URL.Path != "/" {
        http.NotFound(writer, request)
        return
    }
    
    request.ParseForm()

    if len(request.Form["move"]) == 0 {
        http.Error(writer, "A 'move' must be specified", http.StatusBadRequest)
        return
    }

    userMove := request.Form["move"][0]
    botMove := generateMove()
    score := computeScore(userMove, botMove)

    fmt.Fprint(writer, score)
}

func generateMove() string {
    seed := rand.NewSource(time.Now().UnixNano())
    random := rand.New(seed)
    value := random.Intn(base)
    if value >= 0 && value < leftThreshold {
        return rock
    }
    if value >= leftThreshold && value < rightThreshold {
        return paper
    }
    return scissors
}

func computeScore(this string, other string) int {
    if this == other {
        return 0
    }
    if this == rock && other == scissors || this == paper && other == rock || this == scissors && other == paper {
        return 1
    }
    return -1
}
