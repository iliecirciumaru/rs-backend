package structs

type JsonResponse struct {
	message string `json:"message"`
}

func SuccessResponse(message string) JsonResponse {
	return JsonResponse{message:message}
}