package cmd

import (
  cobra "github.com/spf13/cobra"
  fmt   "fmt"
  gorm  "gorm.io/gorm"
  log   "github.com/sirupsen/logrus"
  mysql "gorm.io/driver/mysql"
)

var warnOpeningOrderCmd = &cobra.Command{
	Use:   "warn-opening-order",
	Short: "send mails to warn for opening orders",
	Long:  `every users wich checks the corresponding box in their interfaces will receive a notification
  when the order starts`,
	Run:   warnOpeningOrder,
}

//adds serverCmd to rootCmd
func init() {
	rootCmd.AddCommand(warnOpeningOrderCmd)
}

var db *gorm.DB
var err error

type User struct{
  Id         int    `json:"id"`
  FirstName  string `json:"firstName"`
  LastName   string `json:"lastName"`
  Email	     string `json:"email"`
  FirstName2 string `json:"firstName2"`
  LastName2  string `json:"lastName2"`
  Email2	   string `json:"email2"`
  Flags      int    `json:"flags"`
}

func (User) TableName() string {
    return "User"
}

//runs the server and also does the calculations and send result to client
func warnOpeningOrder(cmd *cobra.Command, args []string) {
  dsn := "docker:docker@tcp(127.0.0.1:3306)/db?charset=utf8mb4&parseTime=True&loc=Local"
  db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

  if err != nil {
    log.Println("Connection Failed to Open")
  } else {
    log.Println("Connection Established")
  }

  // db.AutoMigrate(&User{})

  var user []User
  db.Table("User").Select("Email", "Email2", "Flags").Scan(&user)

  // 7: 4h,24h,Ouverture
  // 6: 24h,Ouverture
  // 4: Ouverture

  for _, u := range user {
    fmt.Printf("%s - %s - %d\n", u.Email, u.Email2, u.Flags)
  }

  log.Info("warn-opening-order has sent mail !")
}
