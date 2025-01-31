# εncooη

PROOF OF CONCEPT

Currently, εncooη is an experimental web application focused on data management.

The research behind εncooη aims to develop a robust engine for data structuring, relationship management, presentation, and, most importantly, intuitive data navigation.

This innovative application requires no database knowledge or skills, as all aspects are managed in real time through a simple, web-based user interface.

Primarily designed for small businesses, the application enables the management of diverse and extensive organizational information, facilitating easy and natural sharing among stakeholders using only a web browser.

εncooη also serves as the foundation for developing a comprehensive business-oriented application development toolkit.

Copyright David Lambert 2025

# Techstack

1. Docker: https://www.docker.com/
1. Postgresql: https://www.postgresql.org/ 
1. Go: https://go.dev/
1. Kafka: https://kafka.apache.org/ (https://www.confluent.io/ at the moment)
1. Svelte: https://svelte.dev/
1. Flowbite Svelte: https://flowbite-svelte.com/
1. Tailwindcss: https://tailwindcss.com/

# History

2007 - Project factory (https://sourceforge.net/projects/projectfactory/). As an open source project for project management, Factory lets you organize actors in teams, define projects, create version-based plans, generate forecast calendars and track statuses. Small and stand alone, it runs on every system with Java.
2012 - Prototype using Ruby on Rails and sqlite3 . First attenpt to make somthing entirely generic and dynamic. The implementation of row-level security came very late and turned to be impossible to make.
2022 - Prototype using Docker, Go, Go-gin, React and Postgreqsl. Adopting Go was great, but React wasn't. The code for the UI was a mess. The traditional approach as a monolith was a mistake.
2025 - Prototype using Docker, Go, Kafka, Postgreqsl and Svelte. Adoption an event-based architecture using Kafka. Svelte is sweet.

# Notes and references

Go
    https://go.dev/
    https://github.com/uber-go/guide/blob/master/style.md

    Contexts: https://www.digitalocean.com/community/tutorials/how-to-use-contexts-in-go

    Docker: https://registry.hub.docker.com/_/golang

PostgreSQL
    https://www.postgresql.org/ 
    https://postgresapp.com/

    https://www.meetspaceapp.com/2016/04/12/passwords-postgresql-pgcrypto.html

    Docker: https://registry.hub.docker.com/_/postgres

    Passwords:
        create an "encoded64 password"
        run CREATE EXTENSION pgcrypto;
        run
            select crypt('encoded64 password', gen_salt('bf', 8));
        on the database

Authentication
    https://blog.logrocket.com/jwt-authentication-go/
    https://jinhoyoo.github.io/Build-a-Go-lang-back-end-with-the-Gin-framework-an/

SAML
    https://medium.com/@arpitkh96/adding-saml-sso-in-your-golang-service-in-20-minutes-e35a30f52abd
    https://pkg.go.dev/github.com/crewjam/saml

Confluent (zookeeper, kafka)
    https://www.confluent.io/ 

Kafka
    Go library: https://pkg.go.dev/github.com/segmentio/kafka-go#section-readme 
    Svelte library: https://kafka.js.org/

    https://www.redpanda.com/guides/kafka-cloud-kafka-headers

    Commit messages: see https://pkg.go.dev/github.com/segmentio/kafka-go#section-readme.

    Streaming:
        https://svelte.dev/docs/kit/load#Streaming-with-promises
        https://khromov.se/sveltekit-streaming-the-complete-guide/
        https://joyofcode.xyz/using-websockets-with-sveltekit 
        https://medium.com/version-1/websockets-in-sveltekit-28e91eec9245
        https://stackoverflow.com/questions/74330190/how-to-respond-with-a-stream-in-a-sveltekit-server-load-function

    Graceful shutdown:
        Backend: https://withcodeexample.com/a-practical-guide-to-using-golang-with-apache-kafka/
        Frontend: https://medium.com/@curtis.porter/graceful-termination-of-kafkajs-client-processes-b05dd185759d 

Tailwind Css
    https://tailwindcss.com/
    https://play.tailwindcss.com/Ow445YYOoI?layout=horizontal    

Future
    Kafka ACL
    Logging thru Kafka
    Create topics for locators, auth and logging

Gen AI
    https://developers.googleblog.com/en/introducing-genkit-for-go-build-scalable-ai-powered-apps-in-go/ 