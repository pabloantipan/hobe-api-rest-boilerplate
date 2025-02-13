package exceptions

type PictureError struct {
	_code    int
	_message string
}

func (e PictureError) Status() int {
	return e._code
}

func (e PictureError) Message() string {
	return e._message
}

func (e PictureError) IsPictureError() bool {
	return e.Status() != NoException
}

func (e PictureError) IsOk() bool {
	return e.Status() == NoException
}

type PictureException interface {
	Status() int
	Message() string
	IsPictureError() bool
	IsOk() bool
}

func NewPictureException(code int) PictureError {
	return PictureError{
		_code:    code,
		_message: getPictureErrorMessage(code),
	}
}

var picture_error_message = map[int]string{
	20000: "no exception",
	50000: "no picture file provided",
	50010: "picture file size exceeds maximum allowed size",
	50020: "picture file extension not allowed",
	50021: "unknown picture file extension",
	50030: "mime type not allowed",
	50090: "internal server error",
}

var (
	NoException         = 20000
	NoFileProvided      = 50000
	ExceedsMaxSize      = 50010
	ExtensionNotAllowed = 50020
	UnknownExtension    = 50021
	MimeTypeNotAllowed  = 50030
	ErrServerError      = 50090
)

func getPictureErrorMessage(code int) string {
	message := picture_error_message[code]
	if message == "" {
		return "unknown error"
	}
	return message
}
