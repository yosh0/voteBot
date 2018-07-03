package main

import (
	"os"
	"fmt"
	"log"
	"path"
	"runtime"
	"strconv"
)

func defaultLog() string {
	_, file, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(file), C.Log.File)
}

func LogFuncStr(fName, text string) {
	f, err := os.OpenFile(fmt.Sprintf("%s%s", C.Log.Dir, C.Log.File), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		f, err = os.OpenFile(defaultLog(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	}
	log.SetOutput(f)
	log.Println(fName, text)
	f.Close()
}

func LogSavedUsers(fName string, SU SavedUser) {
	f, err := os.OpenFile(fmt.Sprintf("%s%s", C.Log.Dir, C.Log.File), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		f, err = os.OpenFile(defaultLog(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	}
	log.SetOutput(f)
	log.Println(fName, "{\n" +
		"	TgID: "+strconv.FormatInt(SU.TgID, 10)+"\n"+
		"	UserName: "+SU.UserName+"\n"+
		"	StartVote: "+strconv.FormatBool(SU.StartVote)+"\n"+
		"	EndVote: "+strconv.FormatBool(SU.EndVote)+"\n"+
		"	Variant: "+strconv.Itoa(SU.VoteVariant)+"\n"+
		"	Category: "+SU.Category+"\n"+
		"	UpdatedAt: "+strconv.FormatInt(SU.UpdatedAt, 10)+"\n"+
		"}",
	)
	f.Close()
}
