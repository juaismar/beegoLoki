# grafanaLoki

_Send logs to grafana Loki_

### Installation 🔧

_Install with the next command:_

```
go get github.com/juaismar/beegoLoki
```

_and import the package with:_

```
import ("github.com/juaismar/beegoLoki")
```
## Working example 🚀

-This is a simple code that sets beego logs to loki
-Can use basic auth (user and password)
```
package initializer

import (
	"log"
	"os"

	"github.com/beego/beego/v2/core/logs"
	"github.com/juaismar/beegoLoki"
)

func ConfigureLogs() {

	logs.Register("loki", func() logs.Logger {
		return &beegoLoki.LokiAdapter{}
	})

	err := logs.SetLogger("loki",
		`{"endpoint":"http://mylokiserver.com:3100/loki/api/v1/push", `+
			`"user":"optional_user", `+
			`"pass":"optional_password", `+
			`"labels":{"job":"my_project","env":"production"}}`,
	)

	if err != nil {
		log.Fatal("failed to set logger:", err)
	}

	logs.SetLevel(6)
	if os.Getenv("ENABLE_DEBUG_LOG") == "1" {
		logs.SetLevel(7)
	}

	logs.Emergency("Log de Emergency")
	logs.Alert("Log de Alert")
	logs.Critical("Log de Critical")
	logs.Error("Log de Error")
	logs.Warning("Log de Warning")
	logs.Notice("Log de Notice")
	logs.Informational("Log de Informational")
	logs.Debug("Log de Debug")
}
```


## Author ✒️

* **Juan Iscar** - (https://github.com/juaismar)

## Thanks 🎁
* All my friends at work.
* Thanks ChatGPT


_Readme.md based in https://gist.github.com/Villanuevand/6386899f70346d4580c723232524d35a_
