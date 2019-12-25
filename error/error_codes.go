package error

type PipaError int

const (
	Ok PipaError = iota
	ErrInvalidTaskString
	ErrDownloadCode
	StatusRequestEntityTooLarge
	StatusUnsupportedMediaType
	ErrNotFoundOssProcess
	ErrInvalidParameter
	ErrInvalidParameterFormat
	ErrInvalidWatermarkProcess
	ErrInvalidWatermarkPicture
	ErrPictureWidthOrHeightTooLong
	ErrWatermarkCanNotProcess
	ErrInvalidParameterTransparency
	ErrInvalidParameterPosition
	ErrInvalidParameterXMargin
	ErrInvalidParameterYMargin
	ErrInvalidParameterVoffset
	ErrInvalidParameterText
	ErrInvalidParameterTextSize
	ErrInvalidParameterRotate
	ErrInvalidParameterFill
	ErrInvalidParameterLimit
	ErrInvalidParameterMode
	ErrInvalidParameterProportion
	ErrInvalidParameterBorder
)

type ErrorStruct struct {
	ErrorCode    int
	ErrorMessage string
}

var ErrorCodeResponse = map[PipaError]ErrorStruct{
	Ok: {
		ErrorCode:    200,
		ErrorMessage: "ok",
	},
	ErrInvalidTaskString: {
		ErrorCode:    400,
		ErrorMessage: "Invalid task string from request.",
	},
	ErrDownloadCode: {
		ErrorCode:    401,
		ErrorMessage: "Download response code is not 200",
	},
	StatusRequestEntityTooLarge: {
		ErrorCode:    413,
		ErrorMessage: "Picture too large",
	},
	StatusUnsupportedMediaType: {
		ErrorCode:    415,
		ErrorMessage: "Unsupported Media Type",
	},
	ErrNotFoundOssProcess: {
		ErrorCode:    402,
		ErrorMessage: "Can not parameter x-oss-process.",
	},
	ErrInvalidParameter: {
		ErrorCode:    403,
		ErrorMessage: "Invalid parameter: param operation type wrong",
	},
	ErrInvalidParameterFormat: {
		ErrorCode:    405,
		ErrorMessage: "Invalid parameter format.",
	},
	ErrInvalidWatermarkProcess: {
		ErrorCode:    406,
		ErrorMessage: "Invalid watermark parameter.",
	},
	ErrInvalidWatermarkPicture: {
		ErrorCode:    406,
		ErrorMessage: "Invalid watermark picture.",
	},
	ErrPictureWidthOrHeightTooLong: {
		ErrorCode:    407,
		ErrorMessage: "Picture Width or Height too long",
	},
	ErrWatermarkCanNotProcess: {
		ErrorCode:    407,
		ErrorMessage: "Watermark can not process",
	},
	ErrInvalidParameterTransparency: {
		ErrorCode:    403,
		ErrorMessage: "Invalid parameter: transparency wrong.",
	},
	ErrInvalidParameterPosition: {
		ErrorCode:    403,
		ErrorMessage: "Invalid parameter: position wrong.",
	},
	ErrInvalidParameterXMargin: {
		ErrorCode:    403,
		ErrorMessage: "Invalid parameter: XMargin wrong.",
	},
	ErrInvalidParameterYMargin: {
		ErrorCode:    403,
		ErrorMessage: "Invalid parameter: YMargin wrong.",
	},
	ErrInvalidParameterVoffset: {
		ErrorCode:    403,
		ErrorMessage: "Invalid parameter: voffset wrong.",
	},
	ErrInvalidParameterText: {
		ErrorCode:    403,
		ErrorMessage: "Invalid parameter: text wrong.",
	},
	ErrInvalidParameterTextSize: {
		ErrorCode:    403,
		ErrorMessage: "Invalid parameter: text size wrong.",
	},
	ErrInvalidParameterRotate: {
		ErrorCode:    403,
		ErrorMessage: "Invalid parameter: rotate wrong.",
	},
	ErrInvalidParameterLimit: {
		ErrorCode:    403,
		ErrorMessage: "Invalid parameter: limit wrong.",
	},
	ErrInvalidParameterFill: {
		ErrorCode:    403,
		ErrorMessage: "Invalid parameter: fill wrong.",
	},
	ErrInvalidParameterMode: {
		ErrorCode:    403,
		ErrorMessage: "Invalid parameter: mode wrong.",
	},
	ErrInvalidParameterProportion: {
		ErrorCode:    403,
		ErrorMessage: "Invalid parameter: proportion wrong.",
	},
	ErrInvalidParameterBorder: {
		ErrorCode:    403,
		ErrorMessage: "Invalid parameter: params for image border are wrong.",
	},
}

func (e PipaError) ErrorCode() (int, string) {
	err, ok := ErrorCodeResponse[e]
	if !ok {
		return 400, "No error has found"
	}
	return err.ErrorCode, err.ErrorMessage
}

func (e PipaError) Error() string {
	err, ok := ErrorCodeResponse[e]
	if !ok {
		return "We encountered an internal error, please try again."
	}
	return err.ErrorMessage
}
