package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
)

func pointer[T any](t T) *T { return &t }

func serverPurge(running []configServer, add int, maxSlots int) []configServer {
	items := make(map[string]int)
	for _, v := range running {
		items[v.ServerID] = v.ResourceSlots
	}

	totalSlots := 0
	for _, v := range items {
		totalSlots += v
	}

	log.Debug("Purging servers to free slots",
		"occupied", totalSlots, "free", maxSlots-totalSlots,
		"total", maxSlots, "required", add)

	if add > maxSlots {
		// TODO: this can naturally occur as this is ran twice, so ye idk if we should log
		// log.Error("This server will exceed the upper slot cap, aborting", "total", maxSlots, "required", add)
		return nil
	}

	if totalSlots+add <= maxSlots {
		log.Debug("Enough slots available, no need to purge", "free", maxSlots-totalSlots, "required", add)
		return []configServer{}
	}

	// Sort items by value
	type kv struct {
		Key   string
		Value int
	}
	var sortedItems []kv
	for k, v := range items {
		sortedItems = append(sortedItems, kv{k, v})
	}
	sort.Slice(sortedItems, func(i, j int) bool {
		return sortedItems[i].Value < sortedItems[j].Value
	})

	requiredSpace := add // FIXME: totalSlots + add - maxSlots

	var itemsToRemove []string
	for i := len(sortedItems) - 1; i >= 0 && requiredSpace > 0; i-- {
		itemsToRemove = append(itemsToRemove, sortedItems[i].Key)
		requiredSpace -= sortedItems[i].Value
	}

	log.Debug("Executing the great purge", "sorted", sortedItems, "toRemove", itemsToRemove)

	var servers []configServer
	for _, v := range itemsToRemove {
		servers = append(servers,
			running[slices.IndexFunc(running, func(w configServer) bool {
				return w.ServerID == v
			})])
	}
	return servers
}

func removeTrailingSlash(s string) string {
	if s[len(s)-1] == '/' {
		return s[:len(s)-1]
	}
	return s
}

func httpGet(url, token string) (map[string]any, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	reader, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	raw, err := io.ReadAll(reader.Body)
	if err != nil {
		return nil, err
	}

	var jsonData map[string]any
	err = json.Unmarshal(raw, &jsonData)
	if err != nil {
		return nil, err
	}

	if jsonData["status"] != nil && jsonData["status"].(string) != "ok" {
		return nil, fmt.Errorf("%s: %s", jsonData["error"], jsonData["error_data"])
	}

	return jsonData, nil
}

func httpPost(url string, data map[string]any, token string) (map[string]any, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	reader, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	raw, err = io.ReadAll(reader.Body)
	if err != nil {
		return nil, err
	}

	var jsonData map[string]any
	err = json.Unmarshal(raw, &jsonData)
	if err != nil {
		return nil, err
	}

	if jsonData["status"] != nil && jsonData["status"].(string) != "ok" {
		return nil, fmt.Errorf("%s: %s", jsonData["error"], jsonData["error_data"])
	}

	return jsonData, nil
}

func hexToInt(hex string) int {
	hex = strings.TrimPrefix(hex, "#")
	hex = strings.TrimPrefix(hex, "0x")

	i, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		return 0
	}
	return int(i)
}

func errWarn(err error) {
	if err != nil {
		log.Warn("Something went wrong", "err", err)
	}
}

func errError(err error) {
	if err != nil {
		log.Error("Something went wrong", "err", err)
	}
}

func errFatal(err error) {
	if err != nil {
		log.Fatal("Something went wrong", "err", err)
	}
}
