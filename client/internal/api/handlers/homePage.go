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

func GetHomePage() (bytes.Buffer, error) {
	var b bytes.Buffer

	p := &Shortener{
		Title: "Shortener",
	}

	err := t.ExecuteTemplate(&b, "homePage.html", p)
	if err != nil {
		return b, fmt.Errorf("A error occured: %s", err)
	}

	return b, nil
}

func PostHomePage(strJSON []byte) (bytes.Buffer, error) {
	var b bytes.Buffer

	p := &Shortener{
		Title: "Shortener",
	}

	err := t.ExecuteTemplate(&b, "homePage.html", p)
	if err != nil {
		return b, fmt.Errorf("A error occured: %s", err)
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
		log.Println("A error occured client Do: ", err)
		p.Error = "Не удалось получить ответ от сервера."
		b.Reset()
		err = t.ExecuteTemplate(&b, "homePage.html", p)
		if err != nil {
			return b, fmt.Errorf("an error redirect homePage.html: %s", err)
		}
		return b, nil
	}
	defer res.Body.Close()

	shortDB := &Shortener{}
	if err = json.NewDecoder(res.Body).Decode(&shortDB); err != nil {
		return b, fmt.Errorf("error unmarshal request: %s, status code: %d", err, http.StatusInternalServerError)
	}

	cliHost := os.Getenv("CLI_HOST")
	if cliHost == "" {
		log.Fatal("unknown CLI_HOST = ", cliHost)
	}

	p.ShortLink = fmt.Sprintf("http://%s/%s", cliHost, shortDB.ShortLink)
	p.FullLink = shortDB.FullLink
	p.CreatedAt = shortDB.CreatedAt
	p.StatLink = fmt.Sprintf("http://%s/stat/%s", cliHost, shortDB.StatLink)

	b.Reset()

	err = t.ExecuteTemplate(&b, "homePage.html", p)
	if err != nil {
		return b, fmt.Errorf("an error occured rendering home page: %s", err)
	}

	return b, nil
}
