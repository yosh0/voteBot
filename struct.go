package main

type Config struct {
	DB DB
	Log Log
	Serv Serv
	Vote Vote
	StartMsg StartMsg
	KbBtnText KbBtnText
	BotMessage BotMessage
}

type DB struct {
	Host	string
	Port	string
	Name	string
	User	string
	Pass	string
	SSL	string
	Type	string
}

type Log struct {
	Dir	string
	File	string
}

type Serv struct {
	Token	string
	Debug	bool
}

type StartMsg struct {
	Start		string
	Command		string
}

type KbBtnText struct {
	Start		string
}

type BotMessage struct {
	VoteVariantSelect	string
	VoteEnd			string
	VoteLimit		string
	VoteStarted		string
}

type Vote struct {
	Photos		[]string
	Buttons		int
	MimeType	string
	Category	string
}

type SavedUser struct {
	TgID		int64
	UserName	string
	StartVote	bool
	EndVote		bool
	VoteVariant	int
	Category	string
	UpdatedAt	int64
}

type VotePictureSet struct {
	Element		int
	MimeType	string
	PictureUrl	string
}
