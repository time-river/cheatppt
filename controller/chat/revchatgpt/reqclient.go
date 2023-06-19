package revchatgpt3

import (
	tls_client "github.com/bogdanfinn/tls-client"
)

type reqClient struct {
	jar    tls_client.CookieJar
	client tls_client.HttpClient
}

var client *reqClient

func newReqClient() *reqClient {
	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(360),
		tls_client.WithClientProfile(tls_client.Safari_IOS_16_0),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
	}
	client, _ := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	client.SetProxy(getHttpProxy())

	return &reqClient{
		jar:    jar,
		client: client,
	}
}

func getReqClient() *reqClient {
	if client == nil {
		client = newReqClient()
	}

	return client
}
