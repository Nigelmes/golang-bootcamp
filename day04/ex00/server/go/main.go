/*
 * Candy Server
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	"log"
	sw "server/go"
)

func main() {
	log.Printf("Server started")

	router := sw.NewRouter()

	log.Fatal(router.Run(":8080"))
}
