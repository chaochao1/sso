package models

import (
	"github.com/cicdi-go/sso/src/utils"
	"time"
	"github.com/cicdi-go/jwt"
)

type User struct {
	*Base              `xorm:"-"`
	Id                 int64  `json:"id"`
	Username           string `xorm:"varchar(100) notnull index default ''" json:"username"`
	RealName           string `xorm:"varchar(100) default ''" json:"real_name"`
	Email              string `xorm:"varchar(50) default ''" json:"email"`
	Status             int    `xorm:"SMALLINT default 1" json:"status"`
	AuthKey            string `xorm:"varchar(32) default ''" json:"-"`
	PasswordHash       string `xorm:"varchar(255) default ''" json:"-"`
	PasswordResetToken string `xorm:"varchar(255) default ''" json:"-"`
	password           string `xorm:"-"`
	CreatedAt          int    `xorm:"created" json:"created_at"`
	UpdatedAt          int    `xorm:"updated" json:"updated_at"`
}

func (u *User) TableName() string {
	return utils.Config.TablePrefix + "user"
}

//func init() {
//	u := new(User)
//	if e, err := u.GetDb(); err != nil {
//		log.Println(err)
//	} else {
//		err := e.Sync2(u)
//		if err != nil {
//			log.Println(err)
//		}
//	}
//}

func (u *User) Insert() (err error) {
	engine, err := u.GetDb()
	if err != nil {
		return
	}
	id, err := engine.Insert(u)
	if err != nil {
		return err
	}
	u.Id = id
	return
}

func (u *User) SetPassword(value string) {
	u.password = value
	u.generateAuthKey()
	u.PasswordHash, _ = utils.SetPassword(u.password, u.AuthKey)
}

func (u *User) GetPasswordHash(p string) string {
	passwordHash, err := utils.SetPassword(p, u.AuthKey)
	if err != nil {
		return ""
	}
	return passwordHash
}

func (u *User) GetPassword() string {
	return u.password
}

func (u *User) generateAuthKey() {
	u.AuthKey = utils.GenerateRandomKey()
}

func (u *User) Verify(p string) bool {
	engine, err := u.GetDb()
	if err != nil {
		return false
	}
	engine.Where("username = ?", u.Username).Get(u)
	return u.GetPasswordHash(p) == u.PasswordHash
}

// 生成jwt
func (u *User) GenerateToken() (token string, expire time.Time, err error) {
	algorithm :=  jwt.HmacSha256("cicdi")
	claims := jwt.NewClaim()
	expire = time.Now().Add(time.Second*time.Duration(utils.Config.Expire))
	claims.Set("Username", u.Username)
	claims.Set("Role", "Admin")
	claims.SetTime("exp", expire)
	token, err = algorithm.Encode(claims)
	return token, expire, err
}