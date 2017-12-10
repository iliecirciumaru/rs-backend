package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/iliecirciumaru/rs-backend/db"
	"github.com/iliecirciumaru/rs-backend/model"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	db, err := db.GetDb("root", "password", "rs")
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = createDb(db)
	if err != nil {
		log.Fatal(err)
	}

	//err = migrateUsers(db)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//err = migrateMovies(db)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//err = migrateRatings(db)
	//if err != nil {
	//	log.Fatal(err)
	//}

	err = migratePosters(2550, 3450)
	if err != nil {
		log.Fatal(err)
	}
}

func migrateUsers(db *sql.DB) error {
	_, err := db.Exec("TRUNCATE TABLE `users`")
	if err != nil {
		return err
	}

	file, err := os.Open("./data/smallset/ratings.csv")
	if err != nil {
		return err
	}

	r := csv.NewReader(bufio.NewReader(file))
	prevUser := 0
	prepareQuery, err := db.Prepare("INSERT INTO users VALUES (?, ?, ?, ?)")

	logins := []string{"log1", "log2", "log3", "log4", "log5"}
	passwords := []string{"pass1", "pass2", "pass3", "pass4", "pass5"}
	names := []string{"name1", "name2", "name3", "name4", "name5"}

	if err != nil {
		return err
	}

	i := -1

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if i == -1 {
			i++
			continue
		}

		userId, err := strconv.Atoi(record[0])
		if err != nil {
			return err
		}
		if userId == prevUser {
			continue
		}
		prevUser = userId
		sf := strconv.Itoa(i)
		i++
		_, err = prepareQuery.Exec(userId, logins[i%len(logins)]+sf, passwords[i%len(passwords)]+sf, names[i%len(names)]+sf)
		if err != nil {
			return err
		}
	}

	prepareQuery.Close()
	fmt.Println("Migration of users successfully completed")

	return nil
}

func migrateRatings(db *sql.DB) error {
	_, err := db.Exec("TRUNCATE TABLE `ratings`")
	if err != nil {
		return err
	}

	file, err := os.Open("./data/smallset/ratings.csv")
	if err != nil {
		return err
	}

	r := csv.NewReader(bufio.NewReader(file))

	prepareQuery, err := db.Prepare("INSERT INTO ratings VALUES (?, ?, ?, ?)")

	if err != nil {
		return err
	}

	i := -1
	var userID, movieID, ts int
	var rating float64
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if i == -1 {
			i++
			continue
		}

		userID, err = strconv.Atoi(record[0])
		if err != nil {
			return err
		}
		movieID, err = strconv.Atoi(record[1])
		if err != nil {
			return err
		}
		rating, err = strconv.ParseFloat(record[2], 64)
		if err != nil {
			return err
		}
		ts, err = strconv.Atoi(record[3])
		if err != nil {
			return err
		}

		//i++
		//if i == 10 {
		//	break
		//}

		_, err = prepareQuery.Exec(userID, movieID, rating, ts)
		if err != nil {
			return err
		}
	}

	prepareQuery.Close()
	fmt.Println("Migration of ratings successfully completed")

	return nil
}

func migrateMovies(db *sql.DB) error {
	_, err := db.Exec("TRUNCATE TABLE `movies`")
	if err != nil {
		return err
	}

	file, err := os.Open("./data/smallset/movies.csv")
	if err != nil {
		return err
	}

	r := csv.NewReader(bufio.NewReader(file))

	prepareQuery, err := db.Prepare("INSERT INTO movies VALUES (?, ?, ?)")

	if err != nil {
		return err
	}

	i := -1

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if i == -1 {
			i++
			continue
		}

		_, err = prepareQuery.Exec(record[0], record[1], record[2])
		if err != nil {
			return err
		}
	}

	prepareQuery.Close()
	fmt.Println("Migration of movies successfully completed")

	return nil
}

func createDb(db *sql.DB) error {
	userTable := `
		CREATE TABLE IF NOT EXISTS users (
			id int(11) NOT NULL AUTO_INCREMENT,
			login varchar(150) DEFAULT NULL,
			password varchar(150) DEFAULT NULL,
			name varchar(150) DEFAULT NULL,
			PRIMARY KEY (id),
			UNIQUE KEY login_UNIQUE (login)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`
	_, err := db.Exec(userTable)
	if err != nil {
		return err
	}

	ratingTable := `
		CREATE TABLE IF NOT EXISTS ratings (
			iduser int(11) NOT NULL,
			idmovie int(11) NOT NULL,
			rating float DEFAULT NULL,
			timestamp int(11) DEFAULT NULL,
			PRIMARY KEY (iduser,idmovie)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`
	_, err = db.Exec(ratingTable)
	if err != nil {
		return err
	}

	movieTable := `
		CREATE TABLE IF NOT EXISTS movies (
		id int(11) NOT NULL AUTO_INCREMENT,
		title varchar(200) NOT NULL,
		information varchar(200) DEFAULT NULL,
		PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`
	_, err = db.Exec(movieTable)
	if err != nil {
		return err
	}

	return nil
}

func migratePosters(fromMovie, toMovie int) error {
	upper, err := db.GetUpperDB("root", "password", "127.0.0.1", "rs")
	if err != nil {
		return err
	}

	var movies []model.Movie

	err = upper.SelectFrom("movies").OrderBy("id").All(&movies)
	if err != nil {
		return err
	}

	if toMovie > len(movies) {
		toMovie = len(movies)
	}

	movies = movies[fromMovie:toMovie]

	links, err := readLinks()
	if err != nil {
		return err
	}

	var imdbID string
	var jsonMovie map[string]interface{}

	for _, m := range movies {
		// fetch from imdb movie
		// http://www.omdbapi.com/?i=tt0114709&apikey=3e4a893b
		imdbID = links[m.ID]
		res, err := http.Get(fmt.Sprintf("http://www.omdbapi.com/?i=tt%s&apikey=3e4a893b", imdbID))
		if err != nil {
			fmt.Printf("Problem during fetch occured: %d - %s\n", m.ID, imdbID)
			return err
		}

		// extract url
		data, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return err
		}

		err = json.Unmarshal(data, &jsonMovie)
		if err != nil {
			fmt.Printf("Problem during unmarshal occured: %d - %s\n", m.ID, imdbID)
			continue
		}

		poster, ok := jsonMovie["Poster"].(string)
		if !ok {
			fmt.Printf("No poster: %d - %s\n", m.ID, imdbID)
			continue
		}

		// save it in db
		_, err = upper.Exec(`UPDATE movies SET poster_image_url = ? WHERE id = ?`, poster, m.ID)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	fmt.Println("Fetching of poster movies successfully completed")

	return nil
}

// return map: movieLensID => imdbID
func readLinks() (map[int64]string, error) {
	file, err := os.Open("./data/smallset/links.csv")
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(bufio.NewReader(file))

	links := make(map[int64]string)

	i := -1

	var movieID int
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}

		if i == -1 {
			i++
			continue
		}

		movieID, _ = strconv.Atoi(record[0])

		links[int64(movieID)] = record[1]

		if err != nil {
			return nil, err
		}
	}

	return links, nil
}
