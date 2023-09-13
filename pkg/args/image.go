package args

import (
	"github.com/AlecAivazis/survey/v2"
)

type ImageOptions struct {
}

func NewImageOptions() *ImageOptions {
	return &ImageOptions{}
}

func (p *ImageOptions) Ask() (result interface{}, err error) {
	var img string
	err = survey.AskOne(&survey.Input{
		Message: "Which image should be used by the ephemeral container?",
		Suggest: func(toComplete string) []string {
			return nil
		},
		Default: "nicolaka/netshoot",
	}, &img, survey.WithShowCursor(true))
	if err != nil {
		return nil, err
	}
	return img, nil
}
