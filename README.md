# εncooη

For the time being, εncooη is an experimental web-based application 
oriented toward data management.

The object of the reasearches εncooη is based on is to propose 
a true engine that allows data structuration, data relationship 
management, data presentation, and, mainly, to make obvious data navigation.

This brand new application doesn't require database system 
learning or skills, because everything in the application 
is managed through a very simple web-based user interface, 
in real time. 

Primary dedicated to small business structures, 
the application lets manage a large amount of multi-purposes information 
within the organization, and share this information across business stakeholders 
in an easy and natural way, using one Internet browser only.

εncooη is also the foundation for the development of a true business-oriented 
application development toolkit.

Copyright David Lambert 2024

# Techstack

Docker: https://www.docker.com/
    Postgresql: https://www.postgresql.org/ 
        Kafka: https://kafka.apache.org/
            Go: https://go.dev/
    Svelte: https://svelte.dev/
        Flowbite Svelte: https://flowbite-svelte.com/
            Tailwindcss: https://tailwindcss.com/

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