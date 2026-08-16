package articles

import (
	"encoding/json"
	"net/http"
	"strings"
)

func NewHandler(service *Service, static http.Handler) http.Handler {
	mux := http.NewServeMux()
	articlesHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "只支持 GET")
			return
		}
		path := strings.TrimPrefix(r.URL.Path, "/api/articles")
		if path != "" && path != "/" {
			id := strings.Trim(path, "/")
			article, ok := service.Find(id)
			if !ok {
				writeError(w, http.StatusNotFound, "文章不存在")
				return
			}
			writeJSON(w, http.StatusOK, article)
			return
		}
		writeJSON(w, http.StatusOK, service.Search(r.URL.Query().Get("q")))
	}
	mux.HandleFunc("/api/articles", articlesHandler)
	mux.HandleFunc("/api/articles/", articlesHandler)
	mux.HandleFunc("/api/tools", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "只支持 GET")
			return
		}
		writeJSON(w, http.StatusOK, []map[string]string{
			{"name": "网线压线钳", "use": "压接 RJ45 水晶头"},
			{"name": "网络测线仪", "use": "检查八芯通断与线序"},
			{"name": "打线刀", "use": "将模块线缆压入 IDC 端子"},
			{"name": "标签机", "use": "给两端和弱电箱端口编号"},
		})
	})
	mux.Handle("/", static)
	return cors(mux)
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
