# metrics
baa middleware for collector api  info for  prometheus.

## Use：
```
// init
b.Use(metrics.Metrics(metrics.Config{
	MetricsPath: "/metrics",
	Namespace:   "trend",
}))
```

## router
```
import "github.com/prometheus/client_golang/prometheus/promhttp"

// metrics
b.Get("/metrics", func(c *baa.Context) {
	promhttp.Handler().ServeHTTP(c.Resp, c.Req)
})
```

## Config

### Namespace 
the Namespace is the namespace of api metrics

### MetricsPath
the path of prometheus collect metrics, default "/metics"

