package beegoLoki

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/beego/beego/v2/core/logs"
)

var levelName = [8]string{"emerg", "alert", "critical", "error", "warning", "notice", "info", "debug"}

type LokiAdapter struct {
	endpoint  string
	user      string
	pass      string
	labels    map[string]interface{}
	formatter logs.LogFormatter
}

func (l *LokiAdapter) Init(config string) error {
	// Aquí puedes deserializar la configuración si es necesario
	var cfg map[string]interface{}
	if err := json.Unmarshal([]byte(config), &cfg); err != nil {
		return err
	}

	l.endpoint = cfg["endpoint"].(string)
	l.user = cfg["user"].(string)
	l.pass = cfg["pass"].(string)
	l.labels = cfg["labels"].(map[string]interface{})
	return nil
}

func (l *LokiAdapter) WriteMsg(lm *logs.LogMsg) error {
	// Construye la carga útil para Loki

	levelLabel := map[string]interface{}{
		"level": levelName[lm.Level],
	}

	// Fusionar las etiquetas existentes con el nivel de gravedad
	for k, v := range l.labels {
		levelLabel[k] = v
	}

	payload := struct {
		Streams []struct {
			Stream map[string]interface{} `json:"stream"`
			Values [][2]string            `json:"values"`
		} `json:"streams"`
	}{
		Streams: []struct {
			Stream map[string]interface{} `json:"stream"`
			Values [][2]string            `json:"values"`
		}{
			{
				Stream: levelLabel,
				Values: [][2]string{
					{fmt.Sprintf("%d", time.Now().UnixNano()), lm.Msg},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Envía la carga útil a Loki
	client := &http.Client{}

	// Crear solicitud HTTP
	req, err := http.NewRequest("POST", l.endpoint, bytes.NewBuffer(body))
	if err != nil {
		// Manejar error
		panic(err)
	}

	// Configurar encabezados necesarios
	req.Header.Set("Content-Type", "application/json")

	// Si estás usando autenticación básica
	basicAuth := base64.StdEncoding.EncodeToString([]byte(l.user + ":" + l.pass))
	req.Header.Set("Authorization", "Basic "+basicAuth)

	// Realizar la solicitud
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (l *LokiAdapter) Destroy() {}
func (l *LokiAdapter) Flush()   {}

func (l *LokiAdapter) SetFormatter(f logs.LogFormatter) {
	l.formatter = f
}
