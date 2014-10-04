package main

import (
    "log"
    "builtins"
    "github.com/docker/docker/engine"
    "os/exec"
)

func main() {
    host := "127.0.0.1:4000"

    eng := engine.New()
    if err := builtins.Register(eng); err != nil {
            log.Fatal(err)
    }

    // 註冊自定義延伸指令
    eng.Register("info2", func(job *engine.Job) engine.Status {
            v := &engine.Env{}

            out, err := exec.Command("uname","-a").Output()
	    if( err == nil ){
		    v.Set( "uname", string(out) )
	    }

            if _, err := v.WriteTo(job.Stdout); err != nil {
                    return job.Error(err)
            }
            return engine.StatusOK
    })

    go func() {
            if err := eng.Job("acceptconnections").Run(); err != nil {
                    log.Fatal(err)
            }
    }()

    job := eng.Job("serveapi", "tcp://"+host)
    if err := job.Run(); err != nil {
	    log.Fatal(err)
    }
}
