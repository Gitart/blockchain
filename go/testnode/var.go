package main  
import "github.com/jinzhu/gorm"
import "time"
import "net/http"

// Public variable
var Cnt int = 0
var _ http.RoundTripper = &transport{}
type Mst map[string]interface{}


// Log
type logs   struct {
     Note   string 
     Acc    string  
     Node   string 
     Type   string
}


// Node 
type node   struct {
     Title   string 
     Ip      string  
}

// Account 
type account struct {
    Title     string 
    Sum       string 
    Dat       string 
    Note      string 
    Tx        string 
    Descript  string
    Createat  string
    Node      string
}

// Network
type Network  struct {
     Id        string   `json:"id"`
}

// Node 
type Node  struct {
     Ip        string   `json:"Ip"`
     Port      string   `json:"Port"`
     Note      string   `json:"Note"`
     Status    string   `json:"Status"`
     Disabled  string   `json:"Disabled"`
     Datetime  string   `json:"Datetime"`
     Network   string
}

// Sdetting structure
type Sett struct {
     Mainport string `json:"Mainport"`
     Version  string `json:"Version"`
     Nodes   []Node
}

// Transport structure
type transport struct {
     http.RoundTripper
}

// type
type Result struct{
    result string
}

type User struct {
    // gorm.Model
    Age          int            `gorm:"Age"`
    Name         string         `gorm:"Name"`
    Num          int  
}

type User1 struct {
    gorm.Model
    Birthday     time.Time
    Age          int
    Name         string  `gorm:"size:255"` // Default size for string is 255, reset it with this tag
    Num          int     `gorm:"AUTO_INCREMENT"`
    CreditCard        CreditCard      // One-To-One relationship (has one - use CreditCard's UserID as foreign key)
    Emails            []Email         // One-To-Many relationship (has many - use Email's UserID as foreign key)
    BillingAddress    Address         // One-To-One relationship (belongs to - use BillingAddressID as foreign key)
    // BillingAddressID  sql.NullInt64
    ShippingAddress   Address         // One-To-One relationship (belongs to - use ShippingAddressID as foreign key)
    ShippingAddressID int
    IgnoreMe          int `gorm:"-"`   // Ignore this field
    Languages         []Language `gorm:"many2many:user_languages;"` // Many-To-Many relationship, 'user_languages' is join table
}

type Email struct {
    ID      int
    UserID  int     `gorm:"index"` // Foreign key (belongs to), tag `index` will create index for this column
    Email   string  `gorm:"type:varchar(100);unique_index"` // `type` set sql type, `unique_index` will create unique index for this column
    Subscribed bool
}

type Address struct {
    ID       int
    Address1 string         `gorm:"not null;unique"` // Set field as not nullable and unique
    Address2 string         `gorm:"type:varchar(100);unique"`
    // Post     sql.NullString `gorm:"not null"`
}

type Language struct {
    ID   int
    Name string `gorm:"index:idx_name_code"` // Create index with name, and will create combined index if find other fields defined same name
    Code string `gorm:"index:idx_name_code"` // `unique_index` also works
}

type CreditCard struct {
    gorm.Model
    UserID  uint
    Number  string
}
