package main

import (  
  "database/sql"
  "fmt"
  "strings"

  _ "github.com/lib/pq"
)



func checkTemp(temp int) error {
  if temp < -273 || temp > 100 {
    return fmt.Errorf("Temperature %d is out of range.", temp)
  }
  return nil
}



func printAll(db *sql.DB) ([]Temp, error) {
	var data []Temp
  rows, err := db.Query("SELECT * FROM temperatures")
	if err != nil {
    return data, err  
  }
  defer rows.Close()
  for rows.Next() {
    var temp Temp
    err = rows.Scan(&temp.ID, &temp.Temperature, &temp.Location)
    if err != nil {
      return data, err
    }
    data = append(data, temp)
    fmt.Println(temp)
  }
  // get any error encountered during iteration
  err = rows.Err()
  if err != nil {
    return data, err
  }
  return data, nil  
}

func insertTemp(newTemp int, newLoc string, db *sql.DB) error {
  err := checkTemp(newTemp)
  if err != nil {
    return err
  }
  sqlStatement := "INSERT INTO temperatures (temperature, location) VALUES (" + fmt.Sprintf("%v", newTemp) + ", '"+ escapeInject(newLoc) + "');"  
  fmt.Println(sqlStatement)
  _, err = db.Exec(sqlStatement)
  if err != nil {
     fmt.Println(err.Error()) 
     return err
  }
  fmt.Println("Insertion was successful!")
  return nil
}

func updateTemp(ID int, newTemp int, db *sql.DB) error {
  err := checkTemp(newTemp)
  if err != nil {
    return err
  }
  sqlStatement := `UPDATE temperatures SET temperature = $2 WHERE id = $1 RETURNING id, temperature`  
  _, err = db.Exec(sqlStatement, ID, newTemp)  
  if err != nil {  
      return err
  }
  fmt.Println(fmt.Sprintf("Id %d was successfully updated to %d", ID, newTemp) )
  return nil
}

func deleteID(ID int, db *sql.DB) error {
  sqlStatement := "DELETE FROM temperatures WHERE id = $1;"  
  _, err := db.Exec(sqlStatement, ID)  
  if err != nil {  
    return err
  }
  fmt.Println(fmt.Sprintf("Id %d was deleted", ID))
  return nil
}

func deleteTemp(temp int, db *sql.DB) error {
  sqlStatement := `DELETE FROM temperatures WHERE temperature = $1;`  
  _, err := db.Exec(sqlStatement, temp)  
  if err != nil {  
    return err
  }
  fmt.Println(fmt.Sprintf("Id's with temperature %d were deleted", temp))
  return nil
}

func deleteLoc(loc string, db *sql.DB) error {
  sqlStatement := `DELETE FROM temperatures WHERE location = $1;`  
  _, err := db.Exec(sqlStatement, loc)  
  if err != nil {  
    return err
  }
  fmt.Println(fmt.Sprintf("Id's with location %d were deleted", loc))
  return nil
}


func escapeInject(unsafe string) string {
  halfSafe := strings.Replace(unsafe, "'", "''", -1)
  return strings.Replace(halfSafe, "\"", "\"\"", -1)
}