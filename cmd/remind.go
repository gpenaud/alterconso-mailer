package cmd

import (
  bytes "bytes"
  cobra "github.com/spf13/cobra"
	viper "github.com/spf13/viper"
  fmt   "fmt"
  gorm  "gorm.io/gorm"
  log   "github.com/sirupsen/logrus"
  mysql "gorm.io/driver/mysql"
  template "html/template"
  json "encoding/json"
  http "net/http"
  ioutil "io/ioutil"
)

var remindCmd = &cobra.Command{
	Use:   "remind",
	Short: "send mails to remind for opening & closing orders",
	Long:  `every users wich checks the corresponding box in their interfaces will receive a notification
  when the order opens or closes`,
	Run:   remind,
}

//adds serverCmd to rootCmd
func init() {
  remindCmd.Flags().StringP("subject", "s", "", "The mail subject (EXAMPLE: \"Orders are opened !\")")
	viper.BindPFlag("subject", remindCmd.Flags().Lookup("subject"))

  remindCmd.Flags().StringP("template-name", "n", "", "The HTML template to use (EXAMPLE: opening_orders.html)")
  viper.BindPFlag("template_name", remindCmd.Flags().Lookup("template-name"))

  remindCmd.Flags().StringP("template-address", "m", "", "Website address to dsplay in the HTML template (EXAMPLE: http://alterconso.leportail.org)")
	viper.BindPFlag("template_address", remindCmd.Flags().Lookup("template-address"))

  remindCmd.Flags().StringP("group-name", "g", "", "Alterconso group name (EXAMPLE: Alterconso du Val de Brenne)")
	viper.BindPFlag("group_name", remindCmd.Flags().Lookup("group-name"))

  remindCmd.Flags().StringP("sender-name", "", "", "Mail sender name (EXAMPLE: Administrateur du groupe Alterconso du Val de Brenne)")
	viper.BindPFlag("sender_name", remindCmd.Flags().Lookup("sender-name"))

  remindCmd.Flags().StringP("sender-mail", "", "", "Mail sender address (EXAMPLE: alterconso@leportail.org)")
	viper.BindPFlag("sender_mail", remindCmd.Flags().Lookup("sender-mail"))

  rootCmd.AddCommand(remindCmd)
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

type MailTemplate struct {
  GroupName       string
  TemplateAddress string
}

type MailAddress struct {
  Email string `json:"email"`
}

type MailObject struct {
  Subject string   `json:"subject"`
  Html    string   `json:"html"`
  FromName string  `json:"from_name"`
  FromEmail string `json:"from_email"`
  To []MailAddress `json:"to"`
}

func send(json []byte) {
  url := "http://0.0.0.0:5000/send"
  req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
  req.Header.Set("Content-Type", "application/json")

  fmt.Println(req)

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    panic(err)
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)

  fmt.Println(string(body))
}

//runs the server and also does the calculations and send result to client
func remind(cmd *cobra.Command, args []string) {
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
  var to []MailAddress

  for _, u := range user {
    if u.Flags >= 4 {
      to = append(to, MailAddress{ Email: u.Email })
      if u.Email2 != "" {
        to = append(to, MailAddress{ Email: u.Email2 })
      }
    }
  }

  t := template.New("opening_order.tmpl")
  body_template_name := fmt.Sprintf("templates/%s.tmpl", viper.GetString("template_name"))
  parsedTemplate, err := t.ParseFiles("templates/_before.tmpl", body_template_name, "templates/_after.tmpl")

  if err != nil {
    panic(err)
  }

  m := MailTemplate{
    GroupName: viper.GetString("group_name"),
    TemplateAddress : viper.GetString("template_address"),
  }

  // for testing purposes
  to = []MailAddress{ { Email: "guillaume.penaud@gmail.com" } }

  o := MailObject{
    Subject: viper.GetString("subject"),
    FromName: fmt.Sprintf("Administrateur du groupe \"%s\"", viper.GetString("group_name")),
    FromEmail: viper.GetString("sender_name"),
    To: to,
  }

  var buffer bytes.Buffer

  err = parsedTemplate.Execute(&buffer, m)
  if err != nil {
  	panic(err)
	}

  o.Html = buffer.String()

  json, err := json.Marshal(o)
  if err != nil {
    panic(err)
  }

  send(json)
}
