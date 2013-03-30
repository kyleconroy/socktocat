package main

import (
	"code.google.com/p/go.net/websocket"
	"github.com/tuxychandru/pubsub"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	ps := pubsub.New(1)

	index, err := template.ParseFiles("index.html")

	if err != nil {
		log.Fatal(err)
	}

	incomingSubs := make(chan string)

	go func() {
		repos := map[string]bool{}

		for {
			repo, ok := <-incomingSubs

			if !ok {
				continue
			}

			_, seen := repos[repo]

			if seen {
				continue
			}

			repos[repo] = true

			go func() {
				var etag string
				client := &http.Client{}

				for {

					url := "https://api.github.com/repos/" + repo + "/events"
					req, err := http.NewRequest("GET", url, nil)

					if err != nil {
						break
					}

					if len(etag) != 0 {
						req.Header.Add("If-None-Match", etag)
					}

					resp, err := client.Do(req)

					if err != nil {
						break
					}

					log.Printf("%s %d", url, resp.StatusCode)

					wait := resp.Header.Get("X-Poll-Interval") + "s"
					sleep, err := time.ParseDuration(wait)
					log.Print(sleep)

					if err != nil {
						break
					}

					if resp.StatusCode == 304 || resp.StatusCode == 404 {
						time.Sleep(sleep)
						continue
					}

					defer resp.Body.Close()

					body, err := ioutil.ReadAll(resp.Body)

					if err != nil {
						break
					}

					ps.Pub(string(body), repo)

					etag = resp.Header.Get("Etag")

					time.Sleep(sleep)
				}
			}()
		}
	}()

	http.Handle("/hooks", websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		request := ws.Request()
		request.ParseForm()

		for _, r := range request.Form["repo"] {
			incomingSubs <- r
		}

		channel := ps.Sub(request.Form["repo"]...)

		for {
			msg, ok := <-channel

			if ok {
				websocket.Message.Send(ws, msg.(string))
			} else {
				break
			}
		}
	}))

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static/"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		index.Execute(w, nil)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
