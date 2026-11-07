
package main

import (
	"fmt"
	"net/url"
	"strings"
)

func main() {

	// 1. Parse URL
	fmt.Println("\n--- 1. URL Parsing ---")

	rawURL := "https://prem@example.com:8080/products/123?name=laptop&price=50000#details"

	u, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Original URL:", u.String())
	fmt.Println("Scheme:", u.Scheme)
	fmt.Println("Host:", u.Host)
	fmt.Println("Hostname:", u.Hostname())
	fmt.Println("Port:", u.Port())
	fmt.Println("User:", u.User.Username())
	fmt.Println("Path:", u.Path)
	fmt.Println("Raw Query:", u.RawQuery)
	fmt.Println("Fragment:", u.Fragment)

	// 2. Read query parameters
	fmt.Println("\n--- 2. Query Parameters ---")

	queryParams := u.Query()

	fmt.Println("Product Name:", queryParams.Get("name"))
	fmt.Println("Price:", queryParams.Get("price"))

	// 3. Add query parameters
	fmt.Println("\n--- 3. Add Query Parameters ---")

	queryParams.Set("category", "electronics")
	queryParams.Set("brand", "Dell")

	u.RawQuery = queryParams.Encode()

	fmt.Println("Updated URL:", u.String())

	// 4. Create URL from scratch
	fmt.Println("\n--- 4. Create URL ---")

	newURL := url.URL{
		Scheme: "https",
		Host:   "example.com",
		Path:   "/api/employees",
	}

	params := url.Values{}
	params.Set("department", "IT")
	params.Set("city", "Faridabad")
	params.Set("active", "true")

	newURL.RawQuery = params.Encode()

	fmt.Println("Created URL:", newURL.String())

	// 5. URL Encoding
	fmt.Println("\n--- 5. URL Encoding ---")

	text := "Go Programming Language & Java"

	encoded := url.QueryEscape(text)
	decoded, _ := url.QueryUnescape(encoded)

	fmt.Println("Original:", text)
	fmt.Println("Encoded:", encoded)
	fmt.Println("Decoded:", decoded)

	// 6. Resolve relative URL
	fmt.Println("\n--- 6. Resolve Relative URL ---")

	baseURL, _ := url.Parse("https://example.com/api/")
	relativeURL, _ := url.Parse("employees/101")

	resolved := baseURL.ResolveReference(relativeURL)

	fmt.Println("Resolved URL:", resolved)

	// 7. URL validation
	fmt.Println("\n--- 7. URL Validation ---")

	testURL := "https://www.google.com"

	parsed, err := url.ParseRequestURI(testURL)

	if err == nil &&
		(parsed.Scheme == "http" || parsed.Scheme == "https") &&
		parsed.Host != "" &&
		!strings.ContainsAny(parsed.Host, " \t\n") {
		fmt.Println("Valid HTTP URL")
	} else {
		fmt.Println("Invalid HTTP URL")
	}
}
