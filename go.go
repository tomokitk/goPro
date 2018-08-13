package main

import ("fmt"
		"bufio"
		"os"
		"database/sql"
		_ "github.com/go-sql-driver/mysql"
	)

func main() {
	fmt.Printf("welcome\n if you want connect database please type connectDB")
	if Question(""){
		fmt.Printf("hello, world\nselect DB name")
		// DB選択
		DBname := inputData()
		db, err := sql.Open("mysql", "root:J02M05A004@tcp(127.0.0.1:13306)/" + DBname)
		if err != nil {
			panic(err.Error())
		}

		defer db.Close()

		for{
			fmt.Printf("input any name")
			data := inputData()
			stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO sample (name) VALUES (?)"))
			res, err := stmt.Exec(data)
			lastId, err := res.LastInsertId()
			rowCnt, err := res.RowsAffected()

			fmt.Printf("ID = %d, affected = %d\n", lastId, rowCnt)

			rows, err := db.Query("SELECT * FROM sample")
			defer rows.Close()
			if err != nil {
				panic(err.Error())
			}

			for rows.Next() {
				var id int
				var name string
				if err := rows.Scan(&id, &name); err != nil {
					panic(err.Error())
				}
				fmt.Println(id, name)
			}

			fmt.Printf("successfully\n")

			if Question("connectDB or disconnect"){
				continue;
			}else{
				db.Close()
				fmt.Printf("see you \n")
				break;
			}
		}
	}
}

func Question(q string) bool {
	result := true
	fmt.Print(q)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		i := scanner.Text()

		if i == "connectDB" || i == "CONNECTDB" {
			break
		} else if i == "disconnect" || i == "DISCONNECT" {
			result = false
			break
		} else {
			fmt.Println("unexpected answer. please input correct one")
			fmt.Print(q)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return result
}

func inputData() string{
	data := ""
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		data = scanner.Text()
		break
	}
	return data
}