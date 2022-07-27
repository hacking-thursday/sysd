module github.com/hacking-thursday/sysd

go 1.18

replace github.com/hacking-thursday/sysd/api/server2 => ../api/server2

replace github.com/hacking-thursday/sysd/mods/loader => ../mods/loader

replace github.com/hacking-thursday/sysd/mods => ../mods

require github.com/hacking-thursday/sysd/api/server2 v0.0.0-00010101000000-000000000000

require (
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/hacking-thursday/sysd/mods v0.0.0-00010101000000-000000000000 // indirect
	github.com/hacking-thursday/sysd/mods/loader v0.0.0-00010101000000-000000000000 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/tsaikd/KDGoLib v0.0.0-20211113074651-c6ea6ab4ee08 // indirect
	golang.org/x/crypto v0.0.0-20190605123033-f99c8df09eb5 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)
