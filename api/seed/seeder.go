// package seed

// import (
// 	"log"

// 	"github.com/jinzhu/gorm"
// 	"github.com/victorsteven/fullstack/api/models"
// )

// var users = []models.User{
// 	models.User{
// 		Nickname:  "Mohammad Rifqi Sucahyo",
// 		Email:     "m.r.sucahyo@gmail.com",
// 		Password:  "password",
// 		Phone:     "085730098098",
// 		PackageID: 1,
// 	},
// 	models.User{
// 		Nickname:  "Afwa",
// 		Email:     "afwa@gmail.com",
// 		Password:  "password",
// 		Phone:     "084778776837",
// 		PackageID: 2,
// 	},
// 	models.User{
// 		Nickname:  "ninta",
// 		Email:     "ninta@gmail.com",
// 		Password:  "password",
// 		Phone:     "084773376837",
// 		PackageID: 3,
// 	},
// }

// var posts = []models.Post{
// 	models.Post{
// 		Title:   "Title 1",
// 		Content: "Hello world 1",
// 	},
// 	models.Post{
// 		Title:   "Title 2",
// 		Content: "Hello world 2",
// 	},
// }

// var packages = []models.Package{
// 	models.Package{
// 		Title:   "Paket 500kb",
// 		Content: "Paket yang ditawarkan untuk bisa digunakan pada 5 device dengan kecepatan hingga 500kb.",
// 		Price:   "50000",
// 		Speed:   "500Kb",
// 	},
// 	models.Package{
// 		Title:   "Paket 1Mb",
// 		Content: "Paket yang ditawarkan untuk bisa digunakan pada 5 device dengan kecepatan hingga 1Mb.",
// 		Price:   "100000",
// 		Speed:   "1Mb",
// 	},
// 	models.Package{
// 		Title:   "Paket 1.5Mb",
// 		Content: "Paket yang ditawarkan untuk bisa digunakan pada 5 device dengan kecepatan hingga 1.5Mb.",
// 		Price:   "150000",
// 		Speed:   "1.5Mb",
// 	},
// 	models.Package{
// 		Title:   "Paket 2Mb",
// 		Content: "Paket yang ditawarkan untuk bisa digunakan pada 5 device dengan kecepatan hingga 2Mb.",
// 		Price:   "200000",
// 		Speed:   "2Mb",
// 	},
// }

// var devices = []models.Device{
// 	models.Device{
// 		Title:      "Sucahyo@Hp",
// 		MacAddress: "03:C0:DE:EE:FE:09",
// 		UserID:     1,
// 		CreatedBy:  "system",
// 		UpdatedBy:  "system",
// 	},
// 	models.Device{
// 		Title:      "Sucahyo@Laptop",
// 		MacAddress: "03:C0:DE:EE:FE:19",
// 		UserID:     1,
// 		CreatedBy:  "system",
// 		UpdatedBy:  "system",
// 	},
// 	models.Device{
// 		Title:      "Sucahyo@Tabled",
// 		MacAddress: "03:C0:DE:EE:FE:29",
// 		UserID:     1,
// 		CreatedBy:  "system",
// 		UpdatedBy:  "system",
// 	},
// 	models.Device{
// 		Title:      "Sucahyo@Komputer",
// 		MacAddress: "03:C0:DE:EE:FE:39",
// 		UserID:     1,
// 		CreatedBy:  "system",
// 		UpdatedBy:  "system",
// 	},
// }

// var bills = []models.Bill{
// 	models.Bill{
// 		UserID:     1,
// 		Price:      70000,
// 		BillMonth:  5,
// 		BillYear:   2020,
// 		BillStatus: "COMPLETE",
// 		CreatedBy:  "system",
// 		UpdatedBy:  "system",
// 	},
// 	models.Bill{
// 		UserID:     1,
// 		Price:      70000,
// 		BillMonth:  6,
// 		BillYear:   2020,
// 		BillStatus: "COMPLETE",
// 		CreatedBy:  "system",
// 		UpdatedBy:  "system",
// 	},
// 	models.Bill{
// 		UserID:     1,
// 		Price:      70000,
// 		BillMonth:  7,
// 		BillYear:   2020,
// 		BillStatus: "COMPLETE",
// 		CreatedBy:  "system",
// 		UpdatedBy:  "system",
// 	},
// }

// func Load(db *gorm.DB) {

// 	err := db.Debug().DropTableIfExists(&models.Post{}, &models.User{}, &models.Package{}, &models.Device{}, &models.Bill{}).Error
// 	if err != nil {
// 		log.Fatalf("cannot drop table: %v", err)
// 	}
// 	err = db.Debug().AutoMigrate(&models.User{}, &models.Post{}, &models.Package{}, &models.Device{}, &models.Bill{}).Error
// 	if err != nil {
// 		log.Fatalf("cannot migrate table: %v", err)
// 	}

// 	for i, _ := range users {
// 		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
// 		if err != nil {
// 			log.Fatalf("cannot seed users table: %v", err)
// 		}
// 	}

// 	for i, _ := range packages {
// 		err = db.Debug().Model(&models.Package{}).Create(&packages[i]).Error
// 		if err != nil {
// 			log.Fatalf("cannot seed package table: %v", err)
// 		}
// 	}

// 	for i, _ := range devices {
// 		err = db.Debug().Model(&models.Device{}).Create(&devices[i]).Error
// 		if err != nil {
// 			log.Fatalf("cannot seed device table: %v", err)
// 		}
// 	}
// 	for i, _ := range bills {
// 		err = db.Debug().Model(&models.Bill{}).Create(&bills[i]).Error
// 		if err != nil {
// 			log.Fatalf("cannot seed device table: %v", err)
// 		}
// 	}

// }
