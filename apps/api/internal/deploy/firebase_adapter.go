package deploy

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"cms-builder/api/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

const (
	firebaseAPIBaseURL    = "https://firebasehosting.googleapis.com/v1beta1"
	firebaseUploadBaseURL = "https://upload-firebasehosting.googleapis.com/upload"
	firebaseScope         = "https://www.googleapis.com/auth/firebase.hosting"
)

type FirebaseAdapter struct {
	APIBaseURL        string
	UploadBaseURL     string
	HTTPClient        *http.Client
	SecretResolver    func(string) (string, error)
	TokenFetcher      func(context.Context, string) (string, error)
	PreviewChannelTTL time.Duration
}

type firebaseDeployConfig struct {
	ProjectID               string
	SiteID                  string
	ServiceAccountSecretRef string
	PreviewChannelID        string
	PublicURL               string
}

type firebaseVersionResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type firebasePopulateFilesRequest struct {
	Files map[string]string `json:"files"`
}

type firebasePopulateFilesResponse struct {
	UploadRequiredHashes []string `json:"uploadRequiredHashes"`
	UploadURL            string   `json:"uploadUrl"`
}

type firebaseChannelResponse struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type firebaseAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type firebaseServiceAccount struct {
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
	TokenURI    string `json:"token_uri"`
}

func NewFirebaseAdapter() FirebaseAdapter {
	return FirebaseAdapter{
		APIBaseURL:        firebaseAPIBaseURL,
		UploadBaseURL:     firebaseUploadBaseURL,
		HTTPClient:        &http.Client{Timeout: 120 * time.Second},
		SecretResolver:    resolveSecretFromEnv,
		TokenFetcher:      fetchFirebaseAccessToken,
		PreviewChannelTTL: 7 * 24 * time.Hour,
	}
}

func (a FirebaseAdapter) Deploy(ctx context.Context, site models.Site, build models.Build, outputPath string) (*DeployResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	provider, config := providerConfigForBuild(site, build)
	if !strings.EqualFold(provider, "firebase") {
		return nil, fmt.Errorf("firebase adapter received unsupported provider %q", provider)
	}

	deployConfig, err := parseFirebaseDeployConfig(site, config, build)
	if err != nil {
		return nil, err
	}

	secretRef := deployConfig.ServiceAccountSecretRef
	if secretRef == "" {
		return nil, errors.New("firebase service account secret reference is required")
	}

	secretJSON, err := a.SecretResolver(secretRef)
	if err != nil {
		return nil, err
	}
	accessToken, err := a.TokenFetcher(ctx, secretJSON)
	if err != nil {
		return nil, err
	}

	files, err := collectStaticFiles(outputPath)
	if err != nil {
		return nil, err
	}
	if len(files.items) == 0 {
		return nil, errors.New("no files found to deploy")
	}

	version, err := a.createVersion(ctx, accessToken, deployConfig.SiteID)
	if err != nil {
		return nil, err
	}
	populated, err := a.populateFiles(ctx, accessToken, deployConfig.SiteID, version.Name, files.manifest)
	if err != nil {
		return nil, err
	}
	if err := a.uploadFiles(ctx, accessToken, populated.UploadURL, files, populated.UploadRequiredHashes); err != nil {
		return nil, err
	}
	if err := a.finalizeVersion(ctx, accessToken, deployConfig.SiteID, version.Name); err != nil {
		return nil, err
	}

	if strings.EqualFold(strings.TrimSpace(build.BuildType), "preview") {
		channelID := deployConfig.PreviewChannelID
		if strings.TrimSpace(channelID) == "" {
			channelID = previewChannelID(site)
		}

		channel, err := a.ensurePreviewChannel(ctx, accessToken, deployConfig.SiteID, channelID)
		if err != nil {
			return nil, err
		}
		if _, err := a.releaseToChannel(ctx, accessToken, deployConfig.SiteID, channelID, version.Name); err != nil {
			return nil, err
		}
		return &DeployResult{
			Provider: "firebase",
			URL:      channel.URL,
			Message:  fmt.Sprintf("deployed preview channel %s", channelID),
		}, nil
	}

	if _, err := a.releaseLive(ctx, accessToken, deployConfig.SiteID, version.Name); err != nil {
		return nil, err
	}

	return &DeployResult{
		Provider: "firebase",
		URL:      firebaseLiveURL(deployConfig),
		Message:  fmt.Sprintf("deployed live site %s", deployConfig.SiteID),
	}, nil
}

func parseFirebaseDeployConfig(site models.Site, config map[string]any, build models.Build) (firebaseDeployConfig, error) {
	deployConfig := firebaseDeployConfig{
		ProjectID:               configString(config, "projectId"),
		SiteID:                  configString(config, "siteId"),
		ServiceAccountSecretRef: configString(config, "serviceAccountSecretRef"),
		PublicURL:               configString(config, "publicUrl"),
		PreviewChannelID:        configString(config, "channelId"),
	}

	if deployConfig.ServiceAccountSecretRef == "" {
		deployConfig.ServiceAccountSecretRef = configString(config, "tokenSecretRef")
	}
	if deployConfig.SiteID == "" {
		deployConfig.SiteID = deployConfig.ProjectID
	}
	if deployConfig.SiteID == "" {
		return firebaseDeployConfig{}, errors.New("firebase siteId or projectId is required")
	}
	if deployConfig.PublicURL == "" && !strings.EqualFold(strings.TrimSpace(build.BuildType), "preview") {
		deployConfig.PublicURL = fmt.Sprintf("https://%s.web.app", deployConfig.SiteID)
	}
	if deployConfig.PreviewChannelID == "" {
		deployConfig.PreviewChannelID = previewChannelID(site)
	}
	return deployConfig, nil
}

func (a FirebaseAdapter) createVersion(ctx context.Context, accessToken, siteID string) (firebaseVersionResponse, error) {
	var version firebaseVersionResponse
	if err := a.doJSON(ctx, accessToken, http.MethodPost, fmt.Sprintf("%s/sites/%s/versions", a.baseAPIURL(), url.PathEscape(siteID)), nil, &version); err != nil {
		return firebaseVersionResponse{}, err
	}
	if strings.TrimSpace(version.Name) == "" {
		return firebaseVersionResponse{}, errors.New("firebase version name was not returned")
	}
	return version, nil
}

func (a FirebaseAdapter) populateFiles(ctx context.Context, accessToken, siteID, versionName string, manifest map[string]string) (firebasePopulateFilesResponse, error) {
	request := firebasePopulateFilesRequest{Files: manifest}
	var response firebasePopulateFilesResponse
	if err := a.doJSON(ctx, accessToken, http.MethodPost, fmt.Sprintf("%s/sites/%s/versions/%s:populateFiles", a.baseAPIURL(), url.PathEscape(siteID), url.PathEscape(versionIDFromName(versionName))), request, &response); err != nil {
		return firebasePopulateFilesResponse{}, err
	}
	return response, nil
}

func (a FirebaseAdapter) uploadFiles(ctx context.Context, accessToken, uploadURL string, files staticFiles, requiredHashes []string) error {
	required := make(map[string]struct{}, len(requiredHashes))
	for _, hash := range requiredHashes {
		required[strings.TrimSpace(hash)] = struct{}{}
	}

	for _, file := range files.sorted() {
		if len(required) > 0 {
			if _, ok := required[file.Hash]; !ok {
				continue
			}
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(uploadURL, "/")+"/"+file.Hash, bytes.NewReader(file.Gzipped))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("Content-Type", "application/octet-stream")

		resp, err := a.httpClient().Do(req)
		if err != nil {
			return err
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("firebase upload failed: %s", strings.TrimSpace(string(body)))
		}
	}

	return nil
}

func (a FirebaseAdapter) finalizeVersion(ctx context.Context, accessToken, siteID, versionName string) error {
	payload := map[string]string{"status": "FINALIZED"}
	return a.doJSON(ctx, accessToken, http.MethodPatch, fmt.Sprintf("%s/sites/%s/versions/%s?update_mask=status", a.baseAPIURL(), url.PathEscape(siteID), url.PathEscape(versionIDFromName(versionName))), payload, nil)
}

func (a FirebaseAdapter) ensurePreviewChannel(ctx context.Context, accessToken, siteID, channelID string) (firebaseChannelResponse, error) {
	ttl := a.PreviewChannelTTL
	if ttl <= 0 {
		ttl = 7 * 24 * time.Hour
	}
	payload := map[string]string{"ttl": fmt.Sprintf("%ds", int64(ttl.Seconds()))}
	var channel firebaseChannelResponse
	err := a.doJSON(ctx, accessToken, http.MethodPost, fmt.Sprintf("%s/sites/%s/channels?channelId=%s", a.baseAPIURL(), url.PathEscape(siteID), url.QueryEscape(channelID)), payload, &channel)
	if err == nil {
		return channel, nil
	}
	if !strings.Contains(strings.ToLower(err.Error()), "409") {
		return firebaseChannelResponse{}, err
	}
	return a.getChannel(ctx, accessToken, siteID, channelID)
}

func (a FirebaseAdapter) getChannel(ctx context.Context, accessToken, siteID, channelID string) (firebaseChannelResponse, error) {
	var channel firebaseChannelResponse
	if err := a.doJSON(ctx, accessToken, http.MethodGet, fmt.Sprintf("%s/sites/%s/channels/%s", a.baseAPIURL(), url.PathEscape(siteID), url.PathEscape(channelID)), nil, &channel); err != nil {
		return firebaseChannelResponse{}, err
	}
	return channel, nil
}

func (a FirebaseAdapter) releaseLive(ctx context.Context, accessToken, siteID, versionName string) (map[string]any, error) {
	var release map[string]any
	if err := a.doJSON(ctx, accessToken, http.MethodPost, fmt.Sprintf("%s/sites/%s/releases?versionName=%s", a.baseAPIURL(), url.PathEscape(siteID), url.QueryEscape(versionName)), nil, &release); err != nil {
		return nil, err
	}
	return release, nil
}

func (a FirebaseAdapter) releaseToChannel(ctx context.Context, accessToken, siteID, channelID, versionName string) (map[string]any, error) {
	var release map[string]any
	if err := a.doJSON(ctx, accessToken, http.MethodPost, fmt.Sprintf("%s/sites/%s/channels/%s/releases?versionName=%s", a.baseAPIURL(), url.PathEscape(siteID), url.PathEscape(channelID), url.QueryEscape(versionName)), nil, &release); err != nil {
		return nil, err
	}
	return release, nil
}

func (a FirebaseAdapter) doJSON(ctx context.Context, accessToken, method, endpoint string, requestBody any, responseBody any) error {
	var body io.Reader
	if requestBody != nil {
		raw, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}
		body = bytes.NewReader(raw)
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := a.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	payload, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("firebase request failed (%d): %s", resp.StatusCode, strings.TrimSpace(string(payload)))
	}
	if responseBody != nil {
		if len(payload) == 0 {
			return nil
		}
		if err := json.Unmarshal(payload, responseBody); err != nil {
			return err
		}
	}
	return nil
}

func (a FirebaseAdapter) httpClient() *http.Client {
	if a.HTTPClient != nil {
		return a.HTTPClient
	}
	return &http.Client{Timeout: 120 * time.Second}
}

func (a FirebaseAdapter) baseAPIURL() string {
	if strings.TrimSpace(a.APIBaseURL) == "" {
		return firebaseAPIBaseURL
	}
	return strings.TrimRight(a.APIBaseURL, "/")
}

func (a FirebaseAdapter) uploadBaseURL() string {
	if strings.TrimSpace(a.UploadBaseURL) == "" {
		return firebaseUploadBaseURL
	}
	return strings.TrimRight(a.UploadBaseURL, "/")
}

func firebaseLiveURL(config firebaseDeployConfig) string {
	if strings.TrimSpace(config.PublicURL) != "" {
		return strings.TrimRight(strings.TrimSpace(config.PublicURL), "/") + "/"
	}
	return fmt.Sprintf("https://%s.web.app/", strings.TrimSpace(config.SiteID))
}

func previewChannelID(site models.Site) string {
	slug := safePathSegment(site.Slug)
	if slug == "" {
		slug = "site"
	}
	return "preview-" + slug
}

func versionIDFromName(name string) string {
	parts := strings.Split(strings.TrimSpace(name), "/")
	if len(parts) == 0 {
		return strings.TrimSpace(name)
	}
	return parts[len(parts)-1]
}

type staticFile struct {
	Path    string
	Hash    string
	Gzipped []byte
}

type staticFiles struct {
	manifest map[string]string
	items    []staticFile
}

func (s staticFiles) sorted() []staticFile {
	items := append([]staticFile(nil), s.items...)
	sort.Slice(items, func(i, j int) bool {
		return items[i].Path < items[j].Path
	})
	return items
}

func collectStaticFiles(root string) (staticFiles, error) {
	items := make([]staticFile, 0)
	manifest := make(map[string]string)

	err := filepath.WalkDir(root, func(current string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}

		contents, err := os.ReadFile(current)
		if err != nil {
			return err
		}

		gzipped, err := gzipBytes(contents)
		if err != nil {
			return err
		}
		hash := sha256Hex(gzipped)

		relative, err := filepath.Rel(root, current)
		if err != nil {
			return err
		}
		relative = "/" + filepath.ToSlash(relative)
		manifest[relative] = hash
		items = append(items, staticFile{Path: relative, Hash: hash, Gzipped: gzipped})
		return nil
	})
	if err != nil {
		return staticFiles{}, err
	}

	return staticFiles{manifest: manifest, items: items}, nil
}

func gzipBytes(contents []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	writer.Header.ModTime = time.Unix(0, 0)
	writer.Header.OS = 255
	if _, err := writer.Write(contents); err != nil {
		_ = writer.Close()
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func sha256Hex(contents []byte) string {
	sum := sha256.Sum256(contents)
	return hex.EncodeToString(sum[:])
}

func resolveSecretFromEnv(ref string) (string, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return "", errors.New("secret reference is required")
	}
	value := strings.TrimSpace(os.Getenv(ref))
	if value == "" {
		return "", fmt.Errorf("secret %s is not set", ref)
	}
	return value, nil
}

func fetchFirebaseAccessToken(ctx context.Context, serviceAccountJSON string) (string, error) {
	var account firebaseServiceAccount
	if err := json.Unmarshal([]byte(serviceAccountJSON), &account); err != nil {
		return "", fmt.Errorf("invalid Firebase service account JSON: %w", err)
	}
	if strings.TrimSpace(account.ClientEmail) == "" || strings.TrimSpace(account.PrivateKey) == "" {
		return "", errors.New("Firebase service account JSON must include client_email and private_key")
	}
	tokenURI := strings.TrimSpace(account.TokenURI)
	if tokenURI == "" {
		tokenURI = "https://oauth2.googleapis.com/token"
	}

	privateKey, err := parseRSAPrivateKey([]byte(account.PrivateKey))
	if err != nil {
		return "", err
	}
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"iss":   account.ClientEmail,
		"scope": firebaseScope,
		"aud":   tokenURI,
		"iat":   now.Unix(),
		"exp":   now.Add(55 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	assertion, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	form := url.Values{}
	form.Set("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	form.Set("assertion", assertion)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURI, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	payload, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("firebase auth token exchange failed (%d): %s", resp.StatusCode, strings.TrimSpace(string(payload)))
	}

	var result firebaseAccessTokenResponse
	if err := json.Unmarshal(payload, &result); err != nil {
		return "", err
	}
	if strings.TrimSpace(result.AccessToken) == "" {
		return "", errors.New("firebase auth token exchange returned an empty access token")
	}
	return result.AccessToken, nil
}

func parseRSAPrivateKey(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("invalid RSA private key PEM")
	}

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not RSA")
	}
	return rsaKey, nil
}
