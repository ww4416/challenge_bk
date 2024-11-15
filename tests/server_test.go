package main

import (
	"crypto/tls"
	"net"
	"net/http"
	"testing"
	"os/exec"
	"strings"
)

// runCurlCommand helper function executes a curl command and returns the output and error
func runCurlCommand(args ...string) (string, error) {
	cmd := exec.Command("curl", args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// TestHTTPToHTTPSRedirect checks if the server correctly redirects HTTP to HTTPS using curl
func TestHTTPToHTTPSRedirect(t *testing.T) {
	// Run curl command to check HTTP response
	output, err := runCurlCommand("-i", "http://localhost")
	if err != nil {
		t.Fatalf("Curl command failed: %v", err)
	}

	// Check if the status code is 301 (Moved Permanently) for redirection
	if !strings.Contains(output, "301 Moved Permanently") {
		t.Errorf("Expected status code 301, got output: %s", output)
	}

	// Check if the 'Location' header points to the HTTPS version
	if !strings.Contains(output, "Location: https://localhost/") {
		t.Errorf("Expected redirect to https://localhost/, got output: %s", output)
	}
}

// TestHTTPSContent checks if the server serves the correct content over HTTPS using curl
func TestHTTPSContent(t *testing.T) {
	// Run curl command to check HTTPS content
	output, err := runCurlCommand("-k", "https://localhost")
	if err != nil {
		t.Fatalf("Curl command failed: %v", err)
	}

	expectedContent := "<body><h1>Hello World!</h1></body>"

	// Ensure output contains the expected content
	if !strings.Contains(output, expectedContent) {
		t.Errorf("Expected content '%s' not found in output: %s", expectedContent, output)
	}
}


// TestServerPorts ensures server is listening on ports 80 and 443 
func TestServerPorts(t *testing.T) {
	// Test port 80 
	conn, err := net.Dial("tcp", "localhost:80")
	if err != nil {
		t.Fatalf("Unable to connect to port 80 (HTTP): %v", err)
	}
	conn.Close()

	// Test port 443 
	conn, err = net.Dial("tcp", "localhost:443")
	if err != nil {
		t.Fatalf("Unable to connect to port 443 (HTTPS): %v", err)
	}
	conn.Close()
}

// TestSSLValidation ensures SSL certificate is valid
func TestSSLValidation(t *testing.T) {
	client := &http.Client{
		Transport: &http.Transport{
			// Skip SSL verification for self-signed certificates
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Sends a request to the HTTPS server
	resp, err := client.Get("https://localhost")
	if err != nil {
		t.Fatalf("Failed to make HTTPS request: %v", err)
	}
	defer resp.Body.Close()

	// Get the TLS connection state to inspect certificate
	connState := resp.TLS
	if connState == nil {
		t.Fatalf("Failed to retrieve TLS connection state")
	}

	// Check if the certificate is available
	if len(connState.PeerCertificates) == 0 {
		t.Fatalf("No certificate provided by server")
	}

	// Validate the certificate common name (CN) or subject
	cert := connState.PeerCertificates[0]
	expectedCN := "localhost"
	if cert.Subject.CommonName != expectedCN {
		t.Errorf("Expected CN %s, but got %s", expectedCN, cert.Subject.CommonName)
	}

	// Check if the certificate expired
	if cert.NotAfter.Before(cert.NotBefore) {
		t.Errorf("SSL certificate is expired. Validity ends at: %v", cert.NotAfter)
	}
}

// TestInvalidRequest checks if the server returns 404 for non-existent resources
func TestInvalidRequest(t *testing.T) {
	client := &http.Client{
		Transport: &http.Transport{
			// Skips SSL verification for self-signed certificates
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Send an HTTPS request for a non-existent resource
	resp, err := client.Get("https://localhost/nonexistent")
	if err != nil {
		t.Fatalf("Failed HTTPS request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the status code is 404 
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %d", resp.StatusCode)
	}
}

// TestBadRequest checks if the server handles bad requests correctly
func TestBadRequest(t *testing.T) {
	// Curl command to simulate a bad request
	output, err := runCurlCommand("-k", "https://localhost/invalid-endpoint")
	if err != nil {
		t.Fatalf("Curl command failed: %v", err)
	}

	expectedStatusCode := "404 Not Found" // or "400 Bad Request", depending on your server setup

	// Check if the output contains the expected status code
	if !strings.Contains(output, expectedStatusCode) {
		t.Errorf("Expected status code '%s' not found in output: %s", expectedStatusCode, output)
	}
}
