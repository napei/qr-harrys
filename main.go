package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
)

type QRData struct {
	Name string
	URL  string
	ID   int
}

const QRPath string = "qr.csv"
const rootURL string = `http://qr.harrythehirer.com.au/Productions`

var mu = &sync.Mutex{}
var data []QRData

func loadCSV() error {
	f, err := os.Open("qr.csv")
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.Read() // skip header

	d := []QRData{}
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			data = d
			return err
		}

		q := QRData{}

		if q.ID, err = strconv.Atoi(row[0]); err != nil {
			return err
		}
		q.Name = row[1]
		q.URL = row[2]
		d = append(d, q)
	}

}

func searchForID(id int) *QRData {
	i := sort.Search(len(data), func(i int) bool { return data[i].ID >= id })
	if i < len(data) && data[i].ID == id {
		return &data[i]
	} else {
		return nil
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			mu.Lock()
			if err := loadCSV(); err != nil {
				log.Printf("error loading CSV: %v", err)
			}
			sort.SliceStable(data, func(i, j int) bool {
				return data[i].ID < data[j].ID
			})
			mu.Unlock()

			log.Println("successfully reloaded CSV")
			time.Sleep(time.Minute * 1)
		}
	}()

	app := fiber.New()
	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect(rootURL, fiber.StatusFound)
	})

	app.Get("/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id", -1)
		if err != nil {
			return c.Status(500).SendString("error: " + err.Error())
		}
		mu.Lock()
		record := searchForID(id)
		mu.Unlock()
		if record == nil {
			return c.Redirect(rootURL, fiber.StatusFound)
		}

		return c.Redirect(record.URL, fiber.StatusFound)
	})

	app.Listen(":3000")
}
