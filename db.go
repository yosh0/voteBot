package main

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/Masterminds/squirrel"
	"github.com/Masterminds/structable"
)

func dbInit() string {
	return fmt.Sprintf("host=%s " +
		"port=%s " +
		"user=%s " +
		"password=%s " +
		"dbname=%s " +
		"sslmode=%s",
		C.DB.Host,
		C.DB.Port,
		C.DB.User,
		C.DB.Pass,
		C.DB.Name,
		C.DB.SSL,
	)
}

func (SU SavedUser) dbInsert() {
	con, err := sql.Open(C.DB.Type, dbInit())
	if (err != nil) {
		LogFuncStr(fName(), err.Error())
	}
	cache := squirrel.NewStmtCacheProxy(con)

	voteBot := NewVoteBotTable(cache)
	voteBot.TgID = SU.TgID
	voteBot.UserName = SU.UserName
	voteBot.VoteVariant = SU.VoteVariant
	voteBot.Category = SU.Category
	voteBot.UpdatedAt = SU.UpdatedAt

	if err := voteBot.Insert(); err != nil {
		LogFuncStr(fName(), err.Error())
	}
	con.Close()
}

func (VB *VoteBotTable) Insert() error {
	return VB.rec.Insert()
}

func NewVoteBotTable(db squirrel.DBProxyBeginner) *VoteBotTable {
	d := new(VoteBotTable)
	d.builder = squirrel.StatementBuilder.RunWith(db)
	if C.DB.Type == DB_DRIVER {
		d.builder = d.builder.PlaceholderFormat(squirrel.Dollar)
	}
	d.rec = structable.New(db, C.DB.Type).Bind(VOTE_BOT_TABLE, d)
	return d
}

const (
	DB_DRIVER		= "postgres"
	VOTE_BOT_TABLE	 	= "vote_bot"
)

type VoteBotTable struct {
	rec		structable.Recorder
	builder 	squirrel.StatementBuilderType

	Id		int		`stbl:"id,PRIMARY_KEY,SERIAL"`
	TgID		int64		`stbl:"tg_id,UNIQUE"`
	UserName	string		`stbl:"user_name"`
	VoteVariant	int		`stbl:"vote_variant"`
	Category	string		`stbl:"category"`
	UpdatedAt	int64		`stbl:"updated_at"`
}
