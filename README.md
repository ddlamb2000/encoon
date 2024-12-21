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

# Notes and references

Go
    https://go.dev/
    https://github.com/uber-go/guide/blob/master/style.md

    Contexts: https://www.digitalocean.com/community/tutorials/how-to-use-contexts-in-go

    Docker: https://registry.hub.docker.com/_/golang

Gin
    https://gin-gonic.com/ 
    Templates: https://pkg.go.dev/text/template, https://gohugo.io/templates/
    Custom middlewares: https://gin-gonic.com/docs/examples/custom-middleware/

Minification of javascripts
    https://gist.github.com/gaearon/42a2ffa41b8319948f9be4076286e1f3

React and JSX
    <!-- Don't use this in production: -->
    <script src="https://unpkg.com/@babel/standalone/babel.min.js"></script>

    Read this section for a production-ready setup with JSX:
    https://reactjs.org/docs/add-react-to-a-website.html#add-jsx-to-a-project

    In a larger project, you can use an integrated toolchain that includes JSX instead:
    https://reactjs.org/docs/create-a-new-react-app.html

Babel
    https://babeljs.io/

Bootstrap
    https://getbootstrap.com/

Icons
    https://icons.getbootstrap.com/

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

Tabulator
    https://tabulator.info/

Rich text editors
    https://github.com/ianstormtaylor/slate
    https://quilljs.com


Test automation
    Pactum: https://pactumjs.github.io
    Mocha: https://mochajs.org
    Nodes.js: https://nodejs.org/
    Cucumber: https://cucumber.io

    Godog: https://github.com/cucumber/godog

    Testing APIs with gin: https://circleci.com/blog/gin-gonic-testing/

npm
    https://docs.npmjs.com/about-npm

nvm
    https://github.com/nvm-sh/nvm


Confluent (zookeeper, kafka)
    https://www.confluent.io/ 