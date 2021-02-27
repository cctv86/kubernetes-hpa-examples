package main

import (
	"github.com/kataras/iris/v12"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/prodanlabs/kubernetes-hpa-examples/config"
	"github.com/prodanlabs/kubernetes-hpa-examples/handle"
	prometheusMiddleware "github.com/prodanlabs/kubernetes-hpa-examples/prometheus"
)

func main() {
	app := config.NewAPP()

	m := prometheusMiddleware.New("iris")

	app.Use(m.ServeHTTP)

	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		// error code handlers are not sharing the same middleware as other routes, so we have
		// to call them inside their body.
		m.ServeHTTP(ctx)

		ctx.Writef("Not Found")
	})

	//app.Get("/metrics", iris.FromStd(promhttp.Handler()))
	app.Get("/metrics", iris.FromStd(promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	)))

	resAPI := app.Party("/api/v1")
	resAPI.Get("/namespaces", handle.GetNameSpace)
	resAPI.Get("/ip", handle.GetIP)
	resAPI.Get("/hostname", handle.GetHostname)

	app.Run(iris.Addr("0.0.0.0:8080"), iris.WithoutServerError(iris.ErrServerClosed), iris.WithOptimizations)
}
