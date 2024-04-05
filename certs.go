//go:generate go run certs_gen.go

package main

var certBytes = "-----BEGIN CERTIFICATE-----\nMIIF+jCCBOKgAwIBAgISA6wnZKCYQhtTFYg+2Cxpsm6UMA0GCSqGSIb3DQEBCwUA\nMDIxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MQswCQYDVQQD\nEwJSMzAeFw0yNDAyMjkyMDMzMTVaFw0yNDA1MjkyMDMzMTRaMCExHzAdBgNVBAMT\nFnNjb3JlYm9hcmQubmV0a290aC5vcmcwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAw\nggIKAoICAQDsIsq5WggRvXgS+7DBHcQ77Qjht8gr1ZFFW88CHK4JV8c+86KctqiS\n7bbOEZ1dtCFtnhXXfT8WRaK+lnoHHWKBTC8IW1H5cg8xDXJK8uIX7glhwiUiCLcC\ntXKoohFgEDVbfmHEWy7p54P77mlBiE5cIZfnP5hMVoCKYRuoJH1wKw6TbB3qUkwU\nuBX8Ur4RzIMqaN9wHDdKqm3Lst0zqX3XSF6CswR5WnYe99xO5fliLBsPKlnpQrCE\n6XMv8+WRZfjhGB9gmIGfI8VGwJa5/r8HBGlRzH/TzSLPkYuYFdpsWIn6NmYqU73D\nnPLg3eL3bfPjxHHFXvGEzIUxprStfuLz2a0I389ib4C+hYT+LheY4XseRULj0cse\nQilkyphdhWPwpO1Vfp0vYn0Q1HRmTmLxDntlHnniZzIhLKlD1sSBjZQ6ryhAzpsB\neWzC73/6XSlmeI85vJJPQct4FyPQ5h1C3gFzvZClHOB6HrvvwOmxKhkWrtOvlxcF\nuWjcJ8kRMKwFMttXTpw3v/vq6LXGFcHQZRPnaO0dkRbmacpiu+GUfFu9JmAR6YS3\nFsCx6wrAiTsgWMMW8LurDnxxQAX8v8YcU+tEqYLpzyWBHv6gQb4uz5zTfTZvKScA\nsLUDwlmN47HZ2Q+CGVSg5BHyKGMhAqsiIubw020/CXx0ODjKRHTSeQIDAQABo4IC\nGTCCAhUwDgYDVR0PAQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEF\nBQcDAjAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBQneSNY36K/7g4/WUYoYrBAaZr4\nOzAfBgNVHSMEGDAWgBQULrMXt1hWy65QCUDmH6+dixTCxjBVBggrBgEFBQcBAQRJ\nMEcwIQYIKwYBBQUHMAGGFWh0dHA6Ly9yMy5vLmxlbmNyLm9yZzAiBggrBgEFBQcw\nAoYWaHR0cDovL3IzLmkubGVuY3Iub3JnLzAhBgNVHREEGjAYghZzY29yZWJvYXJk\nLm5ldGtvdGgub3JnMBMGA1UdIAQMMAowCAYGZ4EMAQIBMIIBBQYKKwYBBAHWeQIE\nAgSB9gSB8wDxAHYAO1N3dT4tuYBOizBbBv5AO2fYT8P0x70ADS1yb+H61BcAAAGN\n9sktnAAABAMARzBFAiEA1FLn/36Aen7mbgv4upzK/qEmL8E9CobPfh1GNWJiuqAC\nIE1YPvfWugnOTaacN5K0AdyXKnW+nw5WpptzacuErA76AHcAouK/1h7eLy8HoNZO\nbTen3GVDsMa1LqLat4r4mm31F9gAAAGN9skt2wAABAMASDBGAiEAu6AlOZbDMtm4\nMhAE2OeYTE3CaXZ6SKYpWnOu9+kphs4CIQCl7URwEGduL1XX8a/rAAQ8GK2VrjhJ\n03CkhA/S9czsVzANBgkqhkiG9w0BAQsFAAOCAQEAYuyGCqVWLfd8FOGT3D6QdkTV\nZL0SlntnY+xzqlxJ/cpqjmqSaL3VWDTyJup9Vyx73bLXtjfvhvgO8oO81gSdkLmn\nKGsDSVrFz/lF80EETnzqF2A+34B9OD7nA6Kb6DdyxkNS72lC223UEvkKSzObfMNS\nDEJrMSqqMJSwlbuvgFUV+zg9dFqEbJlVOA4OR/28HI0pv83AELK/RKRxan8qNjqo\n4+cuIpy2EcmHShU8n+5dzZk5wKtE6Fh+ODo5zXkcp4gDLOJFfEI/UikvaVPD46TH\nPBh01PSpH7hQngLtBUDwqoprwpxStfr9AuLH+2/oe3xkE9PArV+ULqzepwXM6g==\n-----END CERTIFICATE-----\n\n-----BEGIN CERTIFICATE-----\nMIIFFjCCAv6gAwIBAgIRAJErCErPDBinU/bWLiWnX1owDQYJKoZIhvcNAQELBQAw\nTzELMAkGA1UEBhMCVVMxKTAnBgNVBAoTIEludGVybmV0IFNlY3VyaXR5IFJlc2Vh\ncmNoIEdyb3VwMRUwEwYDVQQDEwxJU1JHIFJvb3QgWDEwHhcNMjAwOTA0MDAwMDAw\nWhcNMjUwOTE1MTYwMDAwWjAyMQswCQYDVQQGEwJVUzEWMBQGA1UEChMNTGV0J3Mg\nRW5jcnlwdDELMAkGA1UEAxMCUjMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK\nAoIBAQC7AhUozPaglNMPEuyNVZLD+ILxmaZ6QoinXSaqtSu5xUyxr45r+XXIo9cP\nR5QUVTVXjJ6oojkZ9YI8QqlObvU7wy7bjcCwXPNZOOftz2nwWgsbvsCUJCWH+jdx\nsxPnHKzhm+/b5DtFUkWWqcFTzjTIUu61ru2P3mBw4qVUq7ZtDpelQDRrK9O8Zutm\nNHz6a4uPVymZ+DAXXbpyb/uBxa3Shlg9F8fnCbvxK/eG3MHacV3URuPMrSXBiLxg\nZ3Vms/EY96Jc5lP/Ooi2R6X/ExjqmAl3P51T+c8B5fWmcBcUr2Ok/5mzk53cU6cG\n/kiFHaFpriV1uxPMUgP17VGhi9sVAgMBAAGjggEIMIIBBDAOBgNVHQ8BAf8EBAMC\nAYYwHQYDVR0lBBYwFAYIKwYBBQUHAwIGCCsGAQUFBwMBMBIGA1UdEwEB/wQIMAYB\nAf8CAQAwHQYDVR0OBBYEFBQusxe3WFbLrlAJQOYfr52LFMLGMB8GA1UdIwQYMBaA\nFHm0WeZ7tuXkAXOACIjIGlj26ZtuMDIGCCsGAQUFBwEBBCYwJDAiBggrBgEFBQcw\nAoYWaHR0cDovL3gxLmkubGVuY3Iub3JnLzAnBgNVHR8EIDAeMBygGqAYhhZodHRw\nOi8veDEuYy5sZW5jci5vcmcvMCIGA1UdIAQbMBkwCAYGZ4EMAQIBMA0GCysGAQQB\ngt8TAQEBMA0GCSqGSIb3DQEBCwUAA4ICAQCFyk5HPqP3hUSFvNVneLKYY611TR6W\nPTNlclQtgaDqw+34IL9fzLdwALduO/ZelN7kIJ+m74uyA+eitRY8kc607TkC53wl\nikfmZW4/RvTZ8M6UK+5UzhK8jCdLuMGYL6KvzXGRSgi3yLgjewQtCPkIVz6D2QQz\nCkcheAmCJ8MqyJu5zlzyZMjAvnnAT45tRAxekrsu94sQ4egdRCnbWSDtY7kh+BIm\nlJNXoB1lBMEKIq4QDUOXoRgffuDghje1WrG9ML+Hbisq/yFOGwXD9RiX8F6sw6W4\navAuvDszue5L3sz85K+EC4Y/wFVDNvZo4TYXao6Z0f+lQKc0t8DQYzk1OXVu8rp2\nyJMC6alLbBfODALZvYH7n7do1AZls4I9d1P4jnkDrQoxB3UqQ9hVl3LEKQ73xF1O\nyK5GhDDX8oVfGKF5u+decIsH4YaTw7mP3GFxJSqv3+0lUFJoi5Lc5da149p90Ids\nhCExroL1+7mryIkXPeFM5TgO9r0rvZaBFOvV2z0gp35Z0+L4WPlbuEjN/lxPFin+\nHlUjr8gRsI3qfJOQFy/9rKIJR0Y/8Omwt/8oTWgy1mdeHmmjk7j1nYsvC9JSQ6Zv\nMldlTTKB3zhThV1+XWYp6rjd5JW1zbVWEkLNxE7GJThEUG3szgBVGP7pSWTUTsqX\nnLRbwHOoq7hHwg==\n-----END CERTIFICATE-----\n"
var keyBytes = "-----BEGIN RSA PRIVATE KEY-----\n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\n-----END RSA PRIVATE KEY-----\n"
