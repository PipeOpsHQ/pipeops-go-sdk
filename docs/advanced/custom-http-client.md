# Custom HTTP Client

Use a custom HTTP client for advanced configurations.

## Basic Custom Client

```go
import (
    "net/http"
    "time"
)

customClient := &http.Client{
    Timeout: 60 * time.Second,
}

client, _ := pipeops.NewClient("",
    pipeops.WithHTTPClient(customClient),
)
```

## Custom Transport

```go
transport := &http.Transport{
    MaxIdleConns:        100,
    MaxIdleConnsPerHost: 10,
    IdleConnTimeout:     90 * time.Second,
}

httpClient := &http.Client{
    Transport: transport,
    Timeout:   60 * time.Second,
}

client, _ := pipeops.NewClient("",
    pipeops.WithHTTPClient(httpClient),
)
```

## TLS Configuration

```go
import "crypto/tls"

tlsConfig := &tls.Config{
    MinVersion: tls.VersionTLS12,
}

transport := &http.Transport{
    TLSClientConfig: tlsConfig,
}

httpClient := &http.Client{
    Transport: transport,
}

client, _ := pipeops.NewClient("",
    pipeops.WithHTTPClient(httpClient),
)
```

## Proxy Support

```go
import "net/url"

proxyURL, _ := url.Parse("http://proxy.example.com:8080")

transport := &http.Transport{
    Proxy: http.ProxyURL(proxyURL),
}

httpClient := &http.Client{
    Transport: transport,
}

client, _ := pipeops.NewClient("",
    pipeops.WithHTTPClient(httpClient),
)
```

## See Also

- [Configuration](../getting-started/configuration.md)
