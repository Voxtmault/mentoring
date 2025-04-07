package http_utility

import (
	"net/http"

	"github.com/rotisserie/eris"
	"github.com/voxtmault/mentoring/library-project/pkg/config"
)

type HTTPResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
	Metadata   interface{} `json:"metadata,omitempty"`

	// For internal use only, not to be exposed in the API response
	// These fields are going to be translataed into the appropriate language and message
	InternalMessage      Message           `json:"-"`
	InternalErrorMessage Message           `json:"-"`
	FreeMessage          string            `json:"-"` // For messages that are not registered in the dictionary
	SelectedLanguage     SupportedLanguage `json:"-"` // For the selected language
	ErrorStack           error             `json:"-"` // For the error that is returned from the service
}

func New(cfg *config.AppConfig) *HTTPResponse {
	return &HTTPResponse{
		StatusCode:       http.StatusOK,
		Message:          Messages.Successes[English][Success],
		SelectedLanguage: SupportedLanguage(cfg.AppLanguage),
	}
}

// Translate translates the internal message and error message into the appropriate language
// and sets the translated message and error in the HTTP response.
// It uses the selected language from the HTTP response to determine the language to use.
// It also sets the error stack to the error message if it is not nil.
//
// Call this method when you are about to send the response to the client
func (r *HTTPResponse) Translate() {
	switch r.SelectedLanguage {
	case English:
		if r.ErrorStack == nil {
			r.Message = Messages.Successes[English][r.InternalMessage]
		} else {
			r.Message = Messages.Errors[English][r.InternalErrorMessage]
			r.Error = eris.ToString(r.ErrorStack, true)
		}
	case Indonesian:
		if r.ErrorStack == nil {
			r.Message = Messages.Successes[Indonesian][r.InternalMessage]
		} else {
			r.Message = Messages.Errors[Indonesian][r.InternalErrorMessage]
			r.Error = eris.ToString(r.ErrorStack, true)
		}
	default:
		// Default to English
		if r.ErrorStack == nil {
			r.Message = Messages.Successes[English][r.InternalMessage]
		} else {
			r.Error = Messages.Errors[English][r.InternalErrorMessage]
		}
	}
}
