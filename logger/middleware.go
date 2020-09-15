package logger

import (
	"context"
	"go.uber.org/zap"
	"net/http"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l, _ := zap.NewDevelopment()
		l = l.With(zap.Namespace("hometic"), zap.String("I'm", "gopher"))
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "logger", l)))
	})
}

func L(ctx context.Context) *zap.Logger {
	v := ctx.Value("logger")
	if v == nil{
		return zap.NewExample()
	}

	l, ok := v.(*zap.Logger)
	if ok {
		return l
	}

	return zap.NewExample()
}
