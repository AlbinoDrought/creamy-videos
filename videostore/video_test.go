package videostore

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVideoDeserialization(t *testing.T) {
	inputJSON := []byte(`{
		"id": 69,
		"title": "foo",
		"description": "bar",
		"thumbnail": "file:///dev/null",
		"source": "file:///dev/null",
		"original_file_name": "foo.mp4",
		"time_created": "2018-12-25T00:00:00Z",
		"time_updated": "2018-12-25T00:00:00Z",
		"tags": [
			"barfoo",
			"foobar"
		]
	}`)

	video := Video{}
	err := json.Unmarshal(inputJSON, &video)

	assert.Nil(t, err)

	assert.Equal(t, uint(69), video.ID)
	assert.Equal(t, "foo", video.Title)
	assert.Equal(t, "bar", video.Description)
	assert.Equal(t, "file:///dev/null", video.Thumbnail)
	assert.Equal(t, "file:///dev/null", video.Source)
	assert.Equal(t, "foo.mp4", video.OriginalFileName)
	assert.Equal(t, "2018-12-25T00:00:00Z", video.TimeCreated)
	assert.Equal(t, "2018-12-25T00:00:00Z", video.TimeUpdated)
	assert.Equal(t, 2, len(video.Tags))
	assert.Equal(t, "barfoo", video.Tags[0])
	assert.Equal(t, "foobar", video.Tags[1])
}

func TestVideoSerialization(t *testing.T) {
	video := Video{
		Title:            "foo",
		Description:      "bar",
		Thumbnail:        "file:///dev/null",
		Source:           "file:///dev/null",
		OriginalFileName: "foo.mp4",
		TimeCreated:      "2018-12-25T00:00:00Z",
		TimeUpdated:      "2018-12-25T00:00:00Z",
		Tags: []string{
			"barfoo",
			"foobar",
		},
	}
	video.ID = 69

	expectedJSON := []byte(`{"id":69,"title":"foo","description":"bar","thumbnail":"file:///dev/null","source":"file:///dev/null","original_file_name":"foo.mp4","time_created":"2018-12-25T00:00:00Z","time_updated":"2018-12-25T00:00:00Z","tags":["barfoo","foobar"]}`)
	actualJSON, err := json.Marshal(video)

	assert.Nil(t, err)
	assert.Equal(t, string(expectedJSON), string(actualJSON))
}

func TestVideoFilter_Empty(t *testing.T) {
	type fields struct {
		Title string
		Tags  []string
		Any   string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "empty",
			fields: fields{
				Title: "",
				Tags:  []string{},
				Any:   "",
			},
			want: true,
		},
		{
			name: "title",
			fields: fields{
				Title: "hi",
				Tags:  []string{},
				Any:   "",
			},
			want: false,
		},
		{
			name: "tags",
			fields: fields{
				Title: "",
				Tags:  []string{"a", "b"},
				Any:   "",
			},
			want: false,
		},
		{
			name: "any",
			fields: fields{
				Title: "",
				Tags:  []string{},
				Any:   "foo",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := VideoFilter{
				Title: tt.fields.Title,
				Tags:  tt.fields.Tags,
				Any:   tt.fields.Any,
			}
			if got := filter.Empty(); got != tt.want {
				t.Errorf("VideoFilter.Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVideo_Exists(t *testing.T) {
	type fields struct {
		ID               uint
		Title            string
		Description      string
		Thumbnail        string
		Source           string
		OriginalFileName string
		TimeCreated      string
		TimeUpdated      string
		Tags             []string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "empty",
			fields: fields{},
			want:   false,
		},
		{
			name: "0 ID",
			fields: fields{
				ID:    0,
				Title: "foo",
			},
			want: false,
		},
		{
			name: "non-zero ID",
			fields: fields{
				ID:    69,
				Title: "foo",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			video := Video{
				ID:               tt.fields.ID,
				Title:            tt.fields.Title,
				Description:      tt.fields.Description,
				Thumbnail:        tt.fields.Thumbnail,
				Source:           tt.fields.Source,
				OriginalFileName: tt.fields.OriginalFileName,
				TimeCreated:      tt.fields.TimeCreated,
				TimeUpdated:      tt.fields.TimeUpdated,
				Tags:             tt.fields.Tags,
			}
			if got := video.Exists(); got != tt.want {
				t.Errorf("Video.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}
