package config

import (
	"fmt"
	"strings"
	"testing"

	"github.com/m1k1o/neko/server/pkg/types/codec"
)

func stubCheckPlugins(t *testing.T, fn func([]string) error) {
	t.Helper()

	previous := checkPlugins
	checkPlugins = fn
	t.Cleanup(func() {
		checkPlugins = previous
	})
}

func TestNewVideoPipelineVP8UsesVBR(t *testing.T) {
	stubCheckPlugins(t, func([]string) error { return nil })

	pipeline, err := NewVideoPipeline(codec.VP8(), ":99", "", 25, 2048, 4096, HwEncUnset)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(pipeline, "vp8enc") {
		t.Fatalf("expected vp8enc in pipeline, got: %s", pipeline)
	}
	if !strings.Contains(pipeline, "end-usage=vbr") {
		t.Fatalf("expected VBR rate control in pipeline, got: %s", pipeline)
	}
	if strings.Contains(pipeline, "end-usage=cbr") {
		t.Fatalf("did not expect CBR rate control in pipeline, got: %s", pipeline)
	}
}

func TestNewVideoPipelineOpenH264UsesConfiguredMaxBitrate(t *testing.T) {
	stubCheckPlugins(t, func([]string) error { return nil })

	pipeline, err := NewVideoPipeline(codec.H264(), ":99", "", 25, 2048, 4096, HwEncUnset)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(pipeline, "openh264enc") {
		t.Fatalf("expected openh264enc in pipeline, got: %s", pipeline)
	}
	if !strings.Contains(pipeline, "bitrate=2048000") {
		t.Fatalf("expected openh264 bitrate, got: %s", pipeline)
	}
	if !strings.Contains(pipeline, "max-bitrate=4096000") {
		t.Fatalf("expected configured max bitrate, got: %s", pipeline)
	}
}

func TestNewVideoPipelineOpenH264UsesDefaultMaxBitrateWhenUnset(t *testing.T) {
	stubCheckPlugins(t, func([]string) error { return nil })

	pipeline, err := NewVideoPipeline(codec.H264(), ":99", "", 25, 2048, 0, HwEncUnset)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(pipeline, "max-bitrate=3072000") {
		t.Fatalf("expected default max bitrate, got: %s", pipeline)
	}
}

func TestNewVideoPipelineX264UsesConfiguredMaxBitrate(t *testing.T) {
	stubCheckPlugins(t, func(plugins []string) error {
		for _, plugin := range plugins {
			if plugin == "openh264" {
				return fmt.Errorf("required gstreamer plugin %s not found", plugin)
			}
		}
		return nil
	})

	pipeline, err := NewVideoPipeline(codec.H264(), ":99", "", 25, 2048, 4096, HwEncUnset)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(pipeline, "x264enc") {
		t.Fatalf("expected x264enc in pipeline, got: %s", pipeline)
	}
	if !strings.Contains(pipeline, "option-string=vbv-maxrate=4096") {
		t.Fatalf("expected configured x264 max bitrate option, got: %s", pipeline)
	}
}
