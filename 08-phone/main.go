/*
1.Add entries to db
2. Iterate over the entries to retrieve phone numbers
3. Normalize phone numbers
4. Update entries / Delete entry

*/

package main

import (
	"fmt"

	phonedb "github.com/yousseffarkhani/gophercises/08-phone/db"

	_ "github.com/lib/pq"
)

const (
	host      = "localhost"
	port      = 5432
	user      = "postgres"
	password  = "CASAblanca1"
	dbname    = "gophercises_phone"
	tableName = "phone_numbers"
)

func main() {
	psqlInfo := fmt.Sprintf("host= %s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	must(phonedb.Reset("postgres", psqlInfo, dbname))

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	must(phonedb.Migrate("postgres", psqlInfo, tableName))

	db, err := phonedb.Open("postgres", psqlInfo)
	must(err)
	defer db.Close() // Ici on met un defer car on souhaite fermer la DB après la transaction.

	err = db.Seed()
	must(err)

	phones, err := db.GetAllPhones()
	must(err)
	for _, p := range phones {
		fmt.Printf("working on ... %+v\n", p)
		normalizedNumber := normalize(p.Number)
		if normalizedNumber == p.Number {
			fmt.Println("No changes required")
		} else {
			fmt.Println("Updating or removing ...", normalizedNumber)
			existing, err := db.FindPhone(normalizedNumber)
			must(err)
			if existing != nil {
				fmt.Println("Removing ...", normalizedNumber)
				must(db.DeletePhone(p.ID))
			} else {
				p.Number = normalizedNumber
				fmt.Println("Updating ...", normalizedNumber)
				must(db.UpdatePhone(&p))
			}
		}
	}
}

// must(db.Ping()) // permet d'être sûr que nous sommes connectés à la DB car sql.Open n'effectue pas cette vérification

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func normalize(phone string) string {
	var phoneNormalized []rune
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			phoneNormalized = append(phoneNormalized, ch)
		}
	}
	return string(phoneNormalized)
}
