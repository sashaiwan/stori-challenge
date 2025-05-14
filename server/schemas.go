package server

type RequestSchema struct {
	BucketName string `json:bucket`
	ObjectKey  string `json:key`
	Email      string `json:email`
}

type ResponseSchema struct {
	Status  string `json:status`
	Message string `json:message`
}
