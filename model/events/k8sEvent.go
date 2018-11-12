/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package events

type K8sWatchEvent struct {
	Time    string `json:"time"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Message string `json:"message"`
}
