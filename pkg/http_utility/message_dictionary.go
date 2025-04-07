package http_utility

type SupportedLanguage string

const (
	English    SupportedLanguage = "en"
	Indonesian SupportedLanguage = "id"
)

type Message uint

const (
	// For success
	Success Message = iota
	Created

	// For errors
	GeneralError
	ResourceNotFound
	Unauthorized
	Unprocessable
	BadRequest

	// For SQL related errors
)

type MessageDictionary struct {
	Errors    map[SupportedLanguage]map[Message]string
	Successes map[SupportedLanguage]map[Message]string
}

var Messages = MessageDictionary{
	Errors: map[SupportedLanguage]map[Message]string{
		English: {
			ResourceNotFound: "Resource not found",
			Unauthorized:     "Unauthorized",
			GeneralError:     "Internal server error, please try again",
			Unprocessable:    "Unprocessable entity",
			BadRequest:       "Bad request",
		},
		Indonesian: {
			ResourceNotFound: "Data tidak ditemukan",
			Unauthorized:     "Tidak terotorisasi",
			GeneralError:     "Terjadi kesalahan pada server, silakan coba lagi",
			Unprocessable:    "Entitas tidak dapat di proses",
			BadRequest:       "Input tidak valid",
		},
	},
	Successes: map[SupportedLanguage]map[Message]string{
		English: {
			Success: "Success",
			Created: "Created",
		},
		Indonesian: {
			Success: "Sukses",
			Created: "Berhasil dibuat",
		},
	},
}
