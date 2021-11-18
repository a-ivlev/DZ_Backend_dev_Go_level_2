package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"

	"log"
	"net/http"
	"os"
	"time"
)

var t = template.Must(template.New("homePage.html").Funcs(funcMap).ParseFiles(
	"client/web/homePage.html",
	"client/web/head.html",
	"client/web/footer.html",
))

type FullLink struct {
	FullLink string `json:"full_link"`
}

type Shortener struct {
	Title     string
	FullLink  string    `json:"full_link"`
	ShortLink string    `json:"short_link"`
	StatLink  string    `json:"stat_link"`
	CreatedAt time.Time `json:"created_at"`
	Error     string
}
func (Shortener) Bind(r *http.Request) error {
	return nil
}
func (Shortener) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

var funcMap = template.FuncMap{
	"dateFormat": dateTimeFormat,
}

func dateTimeFormat(layout string, d time.Time) string {
	return d.Format(layout)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Error panic: %s (%T)\n", err, err)
		}
	}()
	p := &Shortener{
		Title: "Shortener",
	}
	if r.Method == http.MethodGet {
		var b bytes.Buffer
		err := t.ExecuteTemplate(&b, "homePage.html", p)
		if err != nil {
			fmt.Fprintf(w, "A error occured.")
			return
		}
		_, err = b.WriteTo(w)
		if err != nil {
			log.Println("error rendering home page: ", err)
		}
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
		}

		var b bytes.Buffer

		fullLink := &FullLink{
			FullLink: r.PostFormValue("fullLink"),
		}

		strJSON, err := json.Marshal(&fullLink)
		if err != nil {
			fmt.Fprintf(w, "A error occured json.NewEncoder(&b).Encode(p).")
		}

		srvHost := os.Getenv("SRV_HOST")
		if srvHost == "" {
			log.Fatal("unknown SRV_HOST = ", srvHost)
		}

		srvPort := os.Getenv("SRV_PORT")
		if srvPort == "" {
			log.Fatal("unknown SRV_PORT = ", srvPort)
		}

		srv := fmt.Sprintf("http://%s:%s/create", srvHost, srvPort)

		client := &http.Client{Timeout: time.Second * 2}
		req, err := http.NewRequest(http.MethodPost, srv, bytes.NewBuffer(strJSON))
		if err != nil {
			fmt.Fprintln(os.Stdout, "A error occured NewRequest.")
		}
		req.Header.Set("Content-Type", "application/json")
		res, err := client.Do(req)
		if err != nil {
			fmt.Fprintln(os.Stdout, "A error occured client Do.")
			p.Error = "Не удалось получить ответ от сервера."
			b.Reset()
			err = t.ExecuteTemplate(&b, "homePage.html", p)
			if err != nil {
				fmt.Fprintf(w, "A error occured.")
				return
			}
			_, err = b.WriteTo(w)
			if err != nil {
				log.Println("error writing error home page: ", err)
			}
		}

		defer res.Body.Close()

		shortDB := &Shortener{}
		if err = json.NewDecoder(res.Body).Decode(&shortDB); err != nil {
			http.Error(w, "error unmarshal request", http.StatusInternalServerError)
			return
		}

		cliHost := os.Getenv("CLI_HOST")
		if cliHost == "" {
			log.Fatal("unknown CLI_HOST = ", cliHost)
		}

		p.ShortLink = fmt.Sprintf("http://%s/%s", cliHost, shortDB.ShortLink)
		p.FullLink = shortDB.FullLink
		p.CreatedAt = shortDB.CreatedAt
		p.StatLink = fmt.Sprintf("http://%s/stat/%s", cliHost, shortDB.StatLink)

		err = t.ExecuteTemplate(&b, "homePage.html", p)
		if err != nil {
			fmt.Fprintf(w, "an error occured rendering home page")
			return
		}

		_, err = b.WriteTo(w)
		if err != nil {
			log.Println("write render home page error: ", err)
		}
	}
}
