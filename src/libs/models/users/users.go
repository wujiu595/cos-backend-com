package users

import (
	"context"
	"cos-backend-com/src/common/dbconn"
	"cos-backend-com/src/common/flake"
	"cos-backend-com/src/common/util"
	"cos-backend-com/src/libs/models"
	"cos-backend-com/src/libs/sdk/account"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/wujiu2020/strip/utils"
)

const (
	accessKeyPrefix = "cn"
)

var (
	keyAlphabets = []byte("0ictnbprkfs21zeovqj6wuldhymag748359x")
)

var Users = &users{
	Connector: models.DefaultConnector,
}

type users struct {
	dbconn.Connector
}

func (p *users) FindOrCreate(ctx context.Context, walletAddr string, output *account.LoginUserResult) (err error) {
	if err := p.GetByWalletAddr(ctx, walletAddr, output); err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}
	if output.Id != flake.ID(0) {
		return
	}
	return p.Create(ctx, walletAddr, output)
}

func (p *users) Get(ctx context.Context, id flake.ID, output interface{}) (err error) {
	stmt := `
		SELECT *
		FROM users
		WHERE id = ${id};
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{id}": id,
	})

	return p.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, output, query, args...)
	})
}

func (p *users) GetByWalletAddr(ctx context.Context, walletAddr string, output interface{}) (err error) {
	stmt := `
		SELECT *
		FROM users
		WHERE wallet_addr = ${walletAddr};
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{walletAddr}": walletAddr,
	})

	return p.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, output, query, args...)
	})
}

func (p *users) Create(ctx context.Context, walletAddr string, output interface{}) (err error) {
	stmt := `
		INSERT INTO users(wallet_addr)
		VALUES (${walletAddr})
		RETURNING id;
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{walletAddr}": walletAddr,
	})

	return p.Invoke(ctx, func(db *sqlx.Tx) error {
		newCtx := dbconn.WithDB(ctx, db)
		var uid flake.ID
		if err := db.GetContext(newCtx, &uid, query, args...); err != nil {
			return err
		}
		if err := p.UpdateSecret(newCtx, uid, output); err != nil {
			return err
		}
		return err
	})
}

func (p *users) UpdateSecret(ctx context.Context, uid flake.ID, output interface{}) (err error) {
	accessKey := CreateAccessKey(uid)
	accessSecret := CreateSecretKey()
	stmt := `
		UPDATE users
		SET (
		     private_secret,
		     public_secret
		)= (
		     ${private_secret},
		     ${public_secret}
		) WHERE id =${id} RETURNING *;
	`
	query, args := util.PgMapQuery(stmt, map[string]interface{}{
		"{private_secret}": accessKey,
		"{public_secret}":  accessSecret,
		"{id}":             uid,
	})

	return p.Invoke(ctx, func(db dbconn.Q) error {
		return db.GetContext(ctx, output, query, args...)
	})
}

func CreateAccessKey(uid flake.ID) string {
	s, err := utils.RandomCreateString(6, keyAlphabets...)
	if err != nil {
		panic(err)
	}
	encoded := utils.NumberEncode(utils.ToStr(uid), keyAlphabets)
	return accessKeyPrefix + encoded + s
}

func CreateSecretKey() string {
	s, err := utils.RandomCreateString(32, keyAlphabets...)
	if err != nil {
		panic(err)
	}
	return s
}
