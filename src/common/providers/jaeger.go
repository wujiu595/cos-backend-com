package providers

import (
	"fmt"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"

	"cos-backend-com/src/common/proto"
	"cos-backend-com/src/common/util"
)

func JaegerTracing(appName string, cfg proto.JaegerConfig) interface{} {
	var tags []opentracing.Tag
	if cfg.Tags != "" {
		tags = util.ParseOpenTracingTags(cfg.Tags)
	}

	jaegerCfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: cfg.AgentAddr,
			User:               cfg.User,
			Password:           cfg.Password,
		},
		Tags: tags,
	}
	tracer, _, err := jaegerCfg.New(appName)
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	return func() opentracing.Tracer {
		return tracer
	}
}
