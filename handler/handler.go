package handler

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/arkavidia-hmif/deployment/handler/events"
)

// Webhook secret
const secret = "shhhhh!!"

func signBody(secret, body []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	computed.Write(body)
	return []byte(computed.Sum(nil))
}

func verifySignature(secret []byte, signature string, body []byte) bool {

	const signaturePrefix = "sha1="
	const signatureLength = 45 // len(SignaturePrefix) + len(hex(sha1))

	if len(signature) != signatureLength || !strings.HasPrefix(signature, signaturePrefix) {
		return false
	}

	actual := make([]byte, 20)
	hex.Decode(actual, []byte(signature[5:]))

	return hmac.Equal(signBody(secret, body), actual)
}

// HookContext is the context of the event
type HookContext struct {
	Signature string
	Event     string
	ID        string
	Payload   []byte
}

func parseHook(secret []byte, req *http.Request) (*HookContext, error) {
	hc := HookContext{}

	if hc.Signature = req.Header.Get("x-hub-signature"); len(hc.Signature) == 0 {
		return nil, errors.New("no signature")
	}

	if hc.Event = req.Header.Get("x-github-event"); len(hc.Event) == 0 {
		return nil, errors.New("no event")
	}

	if hc.ID = req.Header.Get("x-github-delivery"); len(hc.ID) == 0 {
		return nil, errors.New("no event id")
	}

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		return nil, err
	}

	if !verifySignature(secret, hc.Signature, body) {
		return nil, errors.New("Invalid signature")
	}

	hc.Payload = body

	return &hc, nil
}

// RedeployStaging calls the bash script for deploying staging
func RedeployStaging() {
	cmd := exec.Command("/bin/sh", "deploy-staging.sh")
	fmt.Println(cmd.Output())
}

// RedeployProduction calls the bash script for deploying production
func RedeployProduction() {
	cmd := exec.Command("/bin/sh", "deploy-production.sh")
	fmt.Println(cmd.Output())
}

// Handler handle the webhook request
func Handler(w http.ResponseWriter, r *http.Request) {
	hc, err := parseHook([]byte(secret), r)

	w.Header().Set("Content-type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Failed processing hook! ('%s')", err)
		io.WriteString(w, "{}")
		return
	}

	log.Printf("Received %s", hc.Event)

	if hc.Event == "package" {
		var event events.Event
		json.Unmarshal(hc.Payload, &event)
		log.Printf("Receiving %s from ID %d and version %s", event.Action, event.Package.ID, event.Package.PackageVersion.Version)
		if event.Package.ID == 47982 && event.Package.PackageVersion.Version == "latest" {
			RedeployStaging()
		}
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "{}")
	return
}
