package main

import (
	"encoding/json"
	"github.com/karashiiro/godestone/search"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/karashiiro/godestone"
)

func main() {
	s := godestone.NewScraper(godestone.EN)

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		buf := make([]byte, 500)
		_, err := req.Body.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}

		var id uint64
		args := strings.Split(string(buf), " ")
		if len(args) == 1 {
			id, err = strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				log.Println(err)
				return
			}
		} else if len(args) == 3 {
			name := args[0] + " " + args[1]
			for result := range s.SearchCharacters(search.CharacterOptions{Name: name, World: args[2]}) {
				if result.Error != nil {
					log.Println(err)
					return
				}

				if strings.ToLower(result.Name) == strings.ToLower(name) &&
					strings.ToLower(result.World) == strings.ToLower(args[2]) {
					id = uint64(result.ID)
				}
			}
		} else {
			log.Println("Not enough args")
			return
		}

		c, err := s.FetchCharacter(uint32(id))
		if err != nil {
			log.Println(err)
			return
		}

		data, _ := json.Marshal(c)
		_, err = res.Write(data)
		if err != nil {
			log.Println(err)
		}
	})

	log.Fatalln(http.ListenAndServe(":5059", nil))
}
