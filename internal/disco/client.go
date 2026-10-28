package disco

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/unstoppable-tanmay/jvm/internal/platform"
)

const BaseAPI = "https://api.foojay.io/disco/v3.0"

// CommonSupportedVendors lists the primary supported JDK distribution identifiers
var CommonSupportedVendors = []string{
	"temurin",
	"corretto",
	"zulu",
	"graalvm",
	"liberica",
	"microsoft",
	"oracle",
	"dragonwell",
	"sap_machine",
	"semeru",
	"jetbrains",
}

// Client interacts with Foojay Disco API
type Client struct {
	httpClient *http.Client
}

// NewClient returns a new Foojay Disco API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// PackageItem represents a JDK distribution release returned by Disco API
type PackageItem struct {
	ID                  string `json:"id"`
	Distribution        string `json:"distribution"`
	MajorVersion        int    `json:"major_version"`
	JavaVersion         string `json:"java_version"`
	DistributionVersion string `json:"distribution_version"`
	OperatingSystem     string `json:"operating_system"`
	Architecture        string `json:"architecture"`
	ArchiveType         string `json:"archive_type"`
	PackageType         string `json:"package_type"`
	Filename            string `json:"filename"`
	Size                int64  `json:"size"`
	TermOfSupport       string `json:"term_of_support"`
	Links               struct {
		PkgInfoURI             string `json:"pkg_info_uri"`
		PkgDownloadRedirectURI string `json:"pkg_download_redirect"`
	} `json:"links"`
}

type PackagesResponse struct {
	Result  []PackageItem `json:"result"`
	Message string        `json:"message"`
}

type PackageDetails struct {
	Filename          string `json:"filename"`
	DirectDownloadURI string `json:"direct_download_uri"`
	Checksum          string `json:"checksum"`
	ChecksumType      string `json:"checksum_type"`
}

type PackageDetailsResponse struct {
	Result  []PackageDetails `json:"result"`
	Message string           `json:"message"`
}

// QueryPackages queries available packages for the given version and distribution
func (c *Client) QueryPackages(version string, dist string, p platform.Info) ([]PackageItem, error) {
	v := url.Values{}
	if version != "" && version != "latest" {
		v.Set("version", version)
	}
	if dist != "" {
		v.Set("distribution", strings.ToLower(dist))
	}
	v.Set("operating_system", p.OS)
	v.Set("architecture", p.Arch)
	v.Set("package_type", "jdk")
	v.Set("archive_type", p.ArchiveType)
	v.Set("release_status", "ga")

	endpoint := fmt.Sprintf("%s/packages?%s", BaseAPI, v.Encode())

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "jvm-cli/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to contact discovery service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest && dist != "" {
			return nil, fmt.Errorf("unknown JDK vendor/distribution '%s'.\nSupported vendors include: %s",
				dist, strings.Join(CommonSupportedVendors, ", "))
		}
		return nil, fmt.Errorf("discovery service returned status: %s", resp.Status)
	}

	var res PackagesResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return res.Result, nil
}

// GetPackageDetails fetches direct download URL and SHA256 checksum for a package ID
func (c *Client) GetPackageDetails(pkgID string) (*PackageDetails, error) {
	endpoint := fmt.Sprintf("%s/ids/%s", BaseAPI, pkgID)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "jvm-cli/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch package details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch package details, status: %s", resp.Status)
	}

	var res PackageDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode package details: %w", err)
	}

	if len(res.Result) == 0 {
		return nil, fmt.Errorf("package details not found for ID %s", pkgID)
	}

	return &res.Result[0], nil
}
