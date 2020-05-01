package handlers

import (
	"fmt"
	"go-module/helpers"
	"go-module/models"
	"html/template"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	uuid "github.com/satori/go.uuid"
)

var dbSessions = map[string]string{}   // SessionID, Username
var dbUsers = map[string]models.User{} // Username, User

func Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	name := r.FormValue("Username")
	mail := r.FormValue("Email")
	phone := r.FormValue("PhoneNumber")
	pass := r.FormValue("Password")
	cfpass := r.FormValue("Confirm")

	if !helpers.ValidRegisterUsername(name) || !helpers.ValidRegisterEmail(mail) || !helpers.ValidRegisterPhoneNumber(phone) || pass != cfpass {
		tmp, err := template.ParseFiles("templates/register.html")
		if err != nil {
			panic("Could not parse file to register")
		}
		tmp.Execute(w, nil)
	} else {
		db, err := gorm.Open("postgres", address)
		if err != nil {
			panic("Could not open database to register")
		}

		hashedPass := helpers.HashAndSalt(pass)
		user := models.User{Username: name, Email: mail, PhoneNumber: phone, Password: hashedPass}
		db.Create(&user).Save(&user)
		db.Close()

		fmt.Println("User Created Successfully")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	name := r.FormValue("Username")
	pass := r.FormValue("Password")

	if !helpers.ValidLogin(name, pass) {
		tmp, err := template.ParseFiles("templates/login.html")
		if err != nil {
			panic("Could not parse file to login")
		}
		tmp.Execute(w, nil)
	} else {
		// create cookie
		c, err := r.Cookie("session")
		if err != nil {
			SessionID := uuid.NewV4()
			c = &http.Cookie{
				Name:  "session",
				Value: SessionID.String(),
			}
			http.SetCookie(w, c)
		}

		// save the connection
		db, er := gorm.Open("postgres", address)
		if er != nil {
			panic("Could not open database to login")
		}
		var user models.User
		db.Table("users").Where("username = ?", name).Find(&user)
		db.Close()

		dbSessions[c.Value] = name
		dbUsers[name] = user

		fmt.Println("Login Succeeded")
		http.Redirect(w, r, "/welcome", http.StatusSeeOther)
	}
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	// get cookie
	c, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	name, ok := dbSessions[c.Value]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	user := dbUsers[name]

	// process
	tmp, err := template.ParseFiles("templates/welcome.html")
	if err != nil {
		panic("Could not parse file to welcome")
	}
	tmp.Execute(w, user)
}

func GetUpdate(w http.ResponseWriter, r *http.Request) {
	// get user
	c, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	name, ok := dbSessions[c.Value]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	user := dbUsers[name]

	// process
	tmp, er := template.ParseFiles("templates/update.html")
	if er != nil {
		panic("Could not parse files to update")
	}
	tmp.Execute(w, user)
}

func PostUpdate(w http.ResponseWriter, r *http.Request) {
	// get username
	c, _ := r.Cookie("session")
	name, _ := dbSessions[c.Value]

	// process
	r.ParseForm()
	npass := r.FormValue("New Password")
	nphone := r.FormValue("New PhoneNumber")

	db, err := gorm.Open("postgres", address)
	if err != nil {
		panic("Could not open database to update")
	}

	if len(npass) > 0 {
		hashedPass := helpers.HashAndSalt(npass)
		db.Table("users").Where("username = ?", name).Update("password", hashedPass)
	}
	if len(nphone) > 0 {
		db.Table("users").Where("username = ?", name).Update("phone_number", nphone)
	}
	db.Close()

	fmt.Println("Account Updated")
	http.Redirect(w, r, "/welcome", http.StatusSeeOther)
}

func Release(w http.ResponseWriter, r *http.Request) {
	// delete username
	c, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	name, _ := dbSessions[c.Value]
	c.MaxAge = -1
	http.SetCookie(w, c)

	//process
	db, err := gorm.Open("postgres", address)
	if err != nil {
		panic("Could not open database to release")
	}
	db.Table("users").Unscoped().Where("username = ?", name).Delete(&models.User{})

	fmt.Println("Account Released")

	tmp, er := template.ParseFiles("templates/release.html")
	if er != nil {
		panic("Could not parse file to release")
	}
	tmp.Execute(w, nil)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// delete cookie
	c, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	c.MaxAge = -1
	http.SetCookie(w, c)

	fmt.Println("Account Logged Out")

	tmp, err := template.ParseFiles("templates/logout.html")
	if err != nil {
		panic("Could not parse file to logout")
	}
	tmp.Execute(w, nil)
}
