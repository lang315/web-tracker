package main


type User struct {
	ID         string `sql:",pk,varchar(26)" json:"id"`
	Email      string `sql:",varchar(100),notnull,unique" json:"email"`
	Password   string `sql:",varchar(255),notnull" json:"-"`
	Membership string `sql:",varchar(10),notnull" json:"membership"`
	CreatedAt  int64  `sql:",notnull" json:"created_at"`
}
