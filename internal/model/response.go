package model

// Response 统一API响应结构
type Response struct {
    Data interface{} `json:"data"`
    Msg  string     `json:"msg"`
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(data interface{}) Response {
    return Response{
        Data: data,
        Msg:  "",
    }
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(msg string) Response {
    return Response{
        Data: nil,
        Msg:  msg,
    }
}

// OCRResponse OCR识别结果
type OCRResponse struct {
    Texts []string `json:"texts"`
} 