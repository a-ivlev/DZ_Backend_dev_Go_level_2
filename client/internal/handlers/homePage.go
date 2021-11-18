package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
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
	CreatedAt time.Time `json:"created_at"`
	Error     string
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
		client := &http.Client{Timeout: time.Second * 2}
		req, err := http.NewRequest(http.MethodPost, "http://localhost:8035/create", bytes.NewBuffer(strJSON))
		if err != nil {
			fmt.Fprintln(os.Stdout, "A error occured NewRequest.")
		}
		req.Header.Set("Content-Type", "application/json")
		//req, err := http.Post("https://reqbin.com/echo/post/json", "application/json", bytes.NewBuffer(strJSON))
		fmt.Println("json: ", string(strJSON))
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
		fmt.Printf("%s", res.Status)

		defer res.Body.Close()

		shortDB := &Shortener{}
		//if err := json.NewDecoder(res.Body).Decode(&shortDB); err != nil {
		//	http.Error(w, "server error", http.StatusInternalServerError)
		//	return
		//}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("An error occurred while reading the response body: %s", err)
		}
		fmt.Println("body ", string(body))
		err = json.Unmarshal(body, &shortDB)
		if err != nil {
			fmt.Fprintf(w, "Error unmarshal request.")
			return
		}
		fmt.Println("shortDB ", shortDB)
		p.ShortLink = shortDB.ShortLink
		p.FullLink = shortDB.FullLink
		p.CreatedAt = shortDB.CreatedAt

		err = t.ExecuteTemplate(&b, "homePage.html", p)
		if err != nil {
			fmt.Fprintf(w, "A error occured.")
			return
		}
		_, err = b.WriteTo(w)
		if err != nil {
			log.Println("write render home page error: ", err)
		}
	}
}
