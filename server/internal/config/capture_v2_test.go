package config

import (
	"strings"
	"testing"

	"github.com/m1k1o/neko/server/pkg/types/codec"
	"github.com/spf13/viper"
)

func TestCaptureSetV2CreatesMainPipelineWithoutCursor(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	stubCheckPlugins(t, func([]string) error { return nil })

	viper.Set("video_codec", "h264")
	viper.Set("video_bitrate", 2048)
	viper.Set("video_max_bitrate", 4096)

	cfg := Capture{
		Display:    ":99",
		VideoCodec: codec.VP8(),
	}
	cfg.SetV2()

	main, ok := cfg.VideoPipelines["main"]
	if !ok {
		t.Fatalf("expected main pipeline to be generated")
	}
	legacy, ok := cfg.VideoPipelines["legacy"]
	if !ok {
		t.Fatalf("expected legacy pipeline to be generated")
	}

	if !strings.Contains(main.GstPipeline, "show-pointer=false") {
		t.Fatalf("expected main pipeline cursor hidden, got: %s", main.GstPipeline)
	}
	if strings.Contains(main.GstPipeline, "show-pointer=true") {
		t.Fatalf("main pipeline should not contain show-pointer=true, got: %s", main.GstPipeline)
	}

	if !strings.Contains(legacy.GstPipeline, "show-pointer=true") {
		t.Fatalf("expected legacy pipeline cursor visible, got: %s", legacy.GstPipeline)
	}
	if !strings.Contains(legacy.GstPipeline, "max-bitrate=4096000") {
		t.Fatalf("expected legacy pipeline max bitrate, got: %s", legacy.GstPipeline)
	}
}
