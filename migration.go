package main

import (
	"os"
	"log"
	"encoding/csv"
	"bufio"
	"io"
	"github.com/iliecirciumaru/rs-backend/db"
	"database/sql"
	"strconv"
	"fmt"
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
	prevUser := 0;
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
		_, err = prepareQuery.Exec(userId, logins[i%len(logins)] + sf, passwords[i%len(passwords)] + sf, names[i%len(names)] + sf)
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
